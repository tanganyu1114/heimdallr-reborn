package bifrosts

import (
	"bytes"
	"context"
	"fmt"
	v1 "github.com/tanganyu1114/heimdallr-reborn/api/heimdallr_api/v1"
	"github.com/tanganyu1114/heimdallr-reborn/global"
	storev1 "github.com/tanganyu1114/heimdallr-reborn/internal/hmdr_api/store/v1"
	storev1utils "github.com/tanganyu1114/heimdallr-reborn/internal/hmdr_api/store/v1/utils"
	"github.com/tanganyu1114/heimdallr-reborn/internal/pkg/bifrosts"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/internal/pkg/meta/v1"
	"strings"
	"sync"

	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
	nginx_context "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"github.com/marmotedu/errors"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

var (
	wssOnce           = new(sync.Once)
	singletonWSSStore *webServerStatisticsStore
)

type proxyServiceInfoHandlerPipe struct {
	handlerPipe func(target nginx_context.Context, in v1.ProxyServiceInfo) (next nginx_context.PosSet, out []v1.ProxyServiceInfo, err error)
	nextPipe    v1.HandlerPipe[v1.ProxyServiceInfo, []v1.ProxyServiceInfo]
}

func (p *proxyServiceInfoHandlerPipe) Handle(target nginx_context.Context, in v1.ProxyServiceInfo) (out []v1.ProxyServiceInfo, err error) {
	next, tmpOut, err := p.PosMapHandle(target, in)
	if err != nil {
		return nil, err
	}

	if p.nextPipe == nil {
		return tmpOut, nil
	}

	setIdx := 0
	return out, next.Map(func(pos nginx_context.Pos) (nginx_context.Pos, error) {
		o, err := p.nextPipe.Handle(pos.Target(), tmpOut[setIdx])
		if err != nil {
			return pos, err
		}
		out = append(out, o...)
		setIdx++
		return pos, nil
	}).Error()
}

func (p *proxyServiceInfoHandlerPipe) PosMapHandle(target nginx_context.Context, in v1.ProxyServiceInfo) (next nginx_context.PosSet, out []v1.ProxyServiceInfo, err error) {
	return p.handlerPipe(target, in)
}

type proxyServiceInfoHandlerMultiPipe struct {
	handlerPipe func(target nginx_context.Context, in v1.ProxyServiceInfo) (next nginx_context.PosSet, out []v1.ProxyServiceInfo, err error)
	selectPipe  func(target nginx_context.Context) string
	nextPipes   map[string]v1.HandlerPipe[v1.ProxyServiceInfo, []v1.ProxyServiceInfo]
}

func (p *proxyServiceInfoHandlerMultiPipe) Handle(target nginx_context.Context, in v1.ProxyServiceInfo) (out []v1.ProxyServiceInfo, err error) {
	next, tmpOut, err := p.PosMapHandle(target, in)
	if err != nil {
		global.GVA_LOG.Debug("坐标获取异常", zap.Any("err", err))
		return nil, err
	}

	if len(p.nextPipes) == 0 {
		return tmpOut, nil
	}

	setIdx := 0
	return out, next.Map(func(pos nginx_context.Pos) (nginx_context.Pos, error) {
		nextPipe, has := p.nextPipes[p.selectPipe(pos.Target())]
		if !has {
			return pos, errors.Errorf("cannot find next pipe for %s", p.selectPipe(pos.Target()))
		}
		o, err := nextPipe.Handle(pos.Target(), tmpOut[setIdx])
		if err != nil {
			global.GVA_LOG.Debug("坐标获取异常", zap.Any("err", err))
			return pos, err
		}
		out = append(out, o...)
		setIdx++
		return pos, nil
	}).Error()
}

func (p *proxyServiceInfoHandlerMultiPipe) PosMapHandle(target nginx_context.Context, in v1.ProxyServiceInfo) (next nginx_context.PosSet, out []v1.ProxyServiceInfo, err error) {
	return p.handlerPipe(target, in)
}

func newProxyServiceInfoHandlerSinglePipeline(singlePipes ...*proxyServiceInfoHandlerPipe) *proxyServiceInfoHandlerPipe {
	pipe := new(proxyServiceInfoHandlerPipe)
	next := pipe
	for i, singlePipe := range singlePipes {
		next.handlerPipe = singlePipe.handlerPipe
		if i == len(singlePipes)-1 {
			break
		}
		next.nextPipe = singlePipes[i+1]

		next = singlePipes[i+1]
	}

	return pipe
}

func httpAndStreamProxyBriefHandlerPipe() *proxyServiceInfoHandlerMultiPipe {
	return &proxyServiceInfoHandlerMultiPipe{
		handlerPipe: func(config nginx_context.Context, in v1.ProxyServiceInfo) (next nginx_context.PosSet, out []v1.ProxyServiceInfo, err error) {
			next = config.ChildrenPosSet().QueryAll(nginx_context.NewKeyWords(func(targetCtx nginx_context.Context) bool {
				return targetCtx.Type() == context_type.TypeHttp || targetCtx.Type() == context_type.TypeStream
			})).Map(func(pos nginx_context.Pos) (nginx_context.Pos, error) {
				out = append(out, v1.ProxyServiceInfo{})
				return pos, nil
			})
			err = next.Error()
			global.GVA_LOG.Debug("http stream查询", zap.Any("err", err))

			return
		},
		selectPipe: func(target nginx_context.Context) string {
			switch target.Type() {
			case context_type.TypeHttp:
				return "http"
			case context_type.TypeStream:
				return "stream"
			default:
				return "unknown"
			}
		},
		nextPipes: map[string]v1.HandlerPipe[v1.ProxyServiceInfo, []v1.ProxyServiceInfo]{
			"http": newProxyServiceInfoHandlerSinglePipeline(
				httpServerProxyBriefHandlerPipe(),
				httpLocProxyBriefHandlerPipe(),
				httpIfProxyBriefHandlerPipe(),
				httpPPProxyBriefHandlerPipe(),
			),
			"stream": newProxyServiceInfoHandlerSinglePipeline(
				streamServerProxyBriefHandlerPipe(),
				streamPPProxyBriefHandlerPipe(),
			),
		},
	}
}

func httpServerProxyBriefHandlerPipe() *proxyServiceInfoHandlerPipe {
	return &proxyServiceInfoHandlerPipe{
		handlerPipe: func(httpCtx nginx_context.Context, in v1.ProxyServiceInfo) (next nginx_context.PosSet, out []v1.ProxyServiceInfo, err error) {
			next = httpCtx.ChildrenPosSet().QueryAll(nginx_context.NewKeyWordsByType(context_type.TypeServer).SetSkipQueryFilter(nginx_context.SkipDisabledCtxFilterFunc)).
				Map(
					func(srvPos nginx_context.Pos) (nginx_context.Pos, error) {
						proxySvcInfo := in
						srvComment := getCommentInfo(srvPos)
						if snPos := srvPos.QueryOne(nginx_context.NewKeyWordsByType(context_type.TypeDirective).
							SetRegexpMatchingValue(nginx_context.RegexpMatchingServerNameValue).
							SetSkipQueryFilter(nginx_context.SkipDisabledCtxFilterFunc)); snPos.Target().Error() == nil {
							proxySvcInfo.ServerName = snPos.Target().(*local.Directive).Params
						}

						// get listening port
						proxySvcInfo.ServerPort = configuration.Port(srvPos.Target())
						if len(srvComment) > 0 {
							proxySvcInfo.TmpComments = append(proxySvcInfo.TmpComments, srvComment)
						}

						out = append(out, proxySvcInfo)

						return srvPos, nil
					},
				)
			err = next.Error()
			global.GVA_LOG.Debug("httpserver查询", zap.Any("err", err))
			return
		},
	}
}

func httpLocProxyBriefHandlerPipe() *proxyServiceInfoHandlerPipe {
	return &proxyServiceInfoHandlerPipe{
		handlerPipe: func(srvCtx nginx_context.Context, in v1.ProxyServiceInfo) (next nginx_context.PosSet, out []v1.ProxyServiceInfo, err error) {
			next = srvCtx.ChildrenPosSet().QueryAll(nginx_context.NewKeyWordsByType(context_type.TypeLocation).SetSkipQueryFilter(nginx_context.SkipDisabledCtxFilterFunc)).
				Map(
					func(locPos nginx_context.Pos) (nginx_context.Pos, error) {
						proxySvcInfo := in
						locComment := getCommentInfo(locPos)

						proxySvcInfo.Location = locPos.Target().Value()

						if len(locComment) > 0 {
							proxySvcInfo.TmpComments = append(proxySvcInfo.TmpComments, locComment)
						}

						out = append(out, proxySvcInfo)

						return locPos, nil
					},
				)
			err = next.Error()
			global.GVA_LOG.Debug("location查询", zap.Any("err", err))
			return
		},
	}
}

func httpIfProxyBriefHandlerPipe() *proxyServiceInfoHandlerPipe {
	return &proxyServiceInfoHandlerPipe{
		handlerPipe: func(locCtx nginx_context.Context, in v1.ProxyServiceInfo) (next nginx_context.PosSet, out []v1.ProxyServiceInfo, err error) {
			out = append(out, in) // append the info of the locPos
			next = nginx_context.NewPosSet().Append(nginx_context.GetPos(locCtx)).
				AppendWithPosSet(locCtx.ChildrenPosSet().QueryAll(nginx_context.NewKeyWordsByType(context_type.TypeIf).SetSkipQueryFilter(nginx_context.SkipDisabledCtxFilterFunc)).
					Map(
						func(ifPos nginx_context.Pos) (nginx_context.Pos, error) {
							proxySvcInfo := in
							ifComment := getCommentInfo(ifPos)

							proxySvcInfo.IfCondition = ifPos.Target().Value()

							if len(ifComment) > 0 {
								proxySvcInfo.TmpComments = append(proxySvcInfo.TmpComments, ifComment)
							}

							out = append(out, proxySvcInfo)

							return ifPos, nil
						},
					))
			err = next.Error()
			global.GVA_LOG.Debug("http的if查询", zap.Any("err", err))
			return
		},
	}
}

func httpPPProxyBriefHandlerPipe() *proxyServiceInfoHandlerPipe {
	return &proxyServiceInfoHandlerPipe{
		handlerPipe: func(ppFatherCtx nginx_context.Context, in v1.ProxyServiceInfo) (next nginx_context.PosSet, out []v1.ProxyServiceInfo, err error) {
			next = ppFatherCtx.ChildrenPosSet().QueryAll(nginx_context.NewKeyWordsByType(context_type.TypeDirHTTPProxyPass).
				SetSkipQueryFilter(nginx_context.SkipDisabledCtxFilterFunc).
				SetSkipQueryFilter( // filter out redundant search position
					func(targetCtx nginx_context.Context) bool {
						return ppFatherCtx.Type() == context_type.TypeLocation && targetCtx.Type() == context_type.TypeIf
					},
				),
			).
				Map(
					func(ppPos nginx_context.Pos) (nginx_context.Pos, error) {
						proxySvcInfo := in
						proxyPass := ppPos.Target().(*local.HTTPProxyPass)
						ctxPos, err := local.PosBasedOnConfig(proxyPass)
						if err != nil {
							global.GVA_LOG.Debug("解析坐标异常", zap.Any("err", err))
							return ppPos, err
						}
						proxySvcInfo.ContextPos = storev1utils.ConvertToConfigContextPos(ctxPos)
						ppComment := getCommentInfo(ppPos)
						if len(ppComment) > 0 {
							proxySvcInfo.TmpComments = append(proxySvcInfo.TmpComments, ppComment)
						}
						proxySvcInfo.ProxyOriginalURL = proxyPass.OriginalURL
						proxySvcInfo.ProxyURI = proxyPass.URI
						proxySvcInfo.ProxyProtocol = strings.ToUpper(proxyPass.Protocol)
						proxySvcInfo.ProxyAddress = proxyPass.Addresses
						proxySvcInfo.ProxyServiceComment = strings.Join(proxySvcInfo.TmpComments, "\t|\t")

						out = append(out, proxySvcInfo)

						return ppPos, nil
					},
				)
			err = next.Error()
			global.GVA_LOG.Debug("http proxy_pass解析", zap.Any("err", err))
			return
		},
	}
}

func streamServerProxyBriefHandlerPipe() *proxyServiceInfoHandlerPipe {
	return &proxyServiceInfoHandlerPipe{
		handlerPipe: func(streamCtx nginx_context.Context, in v1.ProxyServiceInfo) (next nginx_context.PosSet, out []v1.ProxyServiceInfo, err error) {
			next = streamCtx.ChildrenPosSet().QueryAll(nginx_context.NewKeyWordsByType(context_type.TypeServer).SetSkipQueryFilter(nginx_context.SkipDisabledCtxFilterFunc)).
				Map(
					func(srvPos nginx_context.Pos) (nginx_context.Pos, error) {
						proxySvcInfo := in
						srvComment := getCommentInfo(srvPos)

						// get listening port
						proxySvcInfo.ServerPort = configuration.Port(srvPos.Target())
						if len(srvComment) > 0 {
							proxySvcInfo.TmpComments = append(proxySvcInfo.TmpComments, srvComment)
						}

						out = append(out, proxySvcInfo)

						return srvPos, nil
					},
				)
			err = next.Error()
			global.GVA_LOG.Debug("streamserver查询", zap.Any("err", err))
			return
		},
	}
}

func streamPPProxyBriefHandlerPipe() *proxyServiceInfoHandlerPipe {
	return &proxyServiceInfoHandlerPipe{
		handlerPipe: func(ppFatherCtx nginx_context.Context, in v1.ProxyServiceInfo) (next nginx_context.PosSet, out []v1.ProxyServiceInfo, err error) {
			next = ppFatherCtx.ChildrenPosSet().QueryAll(nginx_context.NewKeyWordsByType(context_type.TypeDirStreamProxyPass).SetSkipQueryFilter(nginx_context.SkipDisabledCtxFilterFunc)).
				Map(
					func(ppPos nginx_context.Pos) (nginx_context.Pos, error) {
						proxySvcInfo := in
						proxyPass := ppPos.Target().(*local.StreamProxyPass)
						ctxPos, err := local.PosBasedOnConfig(proxyPass)
						if err != nil {
							global.GVA_LOG.Debug("解析坐标异常", zap.Any("err", err))
							return ppPos, err
						}
						proxySvcInfo.ContextPos = storev1utils.ConvertToConfigContextPos(ctxPos)
						ppComment := getCommentInfo(ppPos)
						if len(ppComment) > 0 {
							proxySvcInfo.TmpComments = append(proxySvcInfo.TmpComments, ppComment)
						}
						proxySvcInfo.ProxyOriginalURL = proxyPass.OriginalAddress
						proxySvcInfo.ProxyProtocol = "TCP/UDP"
						proxySvcInfo.ProxyAddress = proxyPass.Addresses
						proxySvcInfo.ProxyServiceComment = strings.Join(proxySvcInfo.TmpComments, "\t|\t")

						out = append(out, proxySvcInfo)

						return ppPos, nil
					},
				)
			err = next.Error()
			global.GVA_LOG.Debug("stream prox_pass解析", zap.Any("err", err))
			return
		},
	}
}

type webServerStatisticsStore struct {
	bm bifrosts.Manager
}

func (w *webServerStatisticsStore) ConnectivityCheckOfProxyService(ctx context.Context, opts metav1.WebServerOptions, proxyPassPos metav1.ConfigContextPos) (info v1.ProxyServiceInfo, err error) {
	bc, err := w.bm.GetBifrostClient(opts)
	if err != nil {
		return
	}
	conf, ofp, err := bc.WebServerConfig().Get(opts.ServerName)
	if err != nil {
		return
	}
	proxyPassCtx, err := storev1utils.ParseContext(conf, proxyPassPos.Config, proxyPassPos.ContextPosPath)
	if err != nil {
		return
	}

	proxyPass, ok := proxyPassCtx.(local.ProxyPass)
	if !ok {
		err = errors.Errorf("invalid proxy pass context type: %s, context: %v", proxyPassCtx.Type(), proxyPassCtx)
		return
	}

	resp, err := bc.WebServerConfig().ConnectivityCheckOfProxiedServers(opts.ServerName, proxyPass, ofp.Fingerprints())
	if err != nil {
		return
	}

	switch resp.Type() {
	case context_type.TypeDirHTTPProxyPass:
		httpProxyPass, ok := resp.(*local.HTTPProxyPass)
		if !ok {
			err = errors.Errorf("invalid http proxy pass response: %v", resp)
			return
		}
		info.ProxyOriginalURL = httpProxyPass.OriginalURL
		info.ProxyURI = httpProxyPass.URI
		info.ProxyProtocol = strings.ToUpper(httpProxyPass.Protocol)
		info.ProxyAddress = httpProxyPass.Addresses
	case context_type.TypeDirStreamProxyPass:
		streamProxyPass, ok := resp.(*local.StreamProxyPass)
		if !ok {
			err = errors.Errorf("invalid stream proxy pass response: %v", resp)
			return
		}
		info.ProxyOriginalURL = streamProxyPass.OriginalAddress
		info.ProxyProtocol = "TCP/UDP"
		info.ProxyAddress = streamProxyPass.Addresses
	default:
		err = errors.Errorf("invalid proxy pass response: %v", resp)
	}

	info.ContextPos = proxyPassPos

	return
}

func (w *webServerStatisticsStore) GetProxyServiceInfo(_ context.Context, opts metav1.WebServerOptions) ([]v1.ProxyServiceInfo, error) {
	bc, err := w.bm.GetBifrostClient(opts)
	if err != nil {
		global.GVA_LOG.Debug("获取bifrost managers客户端失败", zap.Any("err", err))
		return nil, err
	}
	conf, _, err := bc.WebServerConfig().Get(opts.ServerName)
	if err != nil {
		global.GVA_LOG.Debug("获取配置异常", zap.Any("err", err))
		return nil, err
	}

	//var proxySvcInfos = make([]v1.ProxyServiceInfo, 0)
	pipeline := httpAndStreamProxyBriefHandlerPipe()
	infos, err := pipeline.Handle(conf.Main(), v1.ProxyServiceInfo{})
	if err != nil {
		global.GVA_LOG.Debug("查询代理配置异常", zap.Any("err", err))
		return nil, err
	}

	return infos, nil
}

func (w *webServerStatisticsStore) ExportProxyServiceInfoToExcel(ctx context.Context, opts metav1.WebServerOptions) ([]byte, error) {
	proxies, err := w.GetProxyServiceInfo(ctx, opts)
	if err != nil {
		return nil, err
	}

	// 创建 Excel 文件
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			return
		}
	}()

	// 设置工作表名称
	sheetName := "代理服务信息"
	f.SetSheetName("Sheet1", sheetName)

	// 设置表头
	headers := []string{
		"反向代理服务描述",
		"反向代理服务名称",
		"反向代理服务侦听端口",
		"反向代理路由路径",
		"特殊判断条件",
		"被反向代理原始 URL",
		"被反向代理协议",
		"被反向代理地址集",
	}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// 创建自动换行样式（包含垂直居中和文本缩进）
	wrapStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			WrapText: true,     // 自动换行
			Vertical: "center", // 垂直居中
			Indent:   1,        // 缩进
		},
	})
	if err != nil {
		return nil, err
	}

	// 填充数据，并记录每列的最大长度
	colMaxLen := make([]int, len(headers)) // 记录每列的最大字符数
	for rowIdx, info := range proxies {
		rowNum := rowIdx + 2 // 从第 2 行开始（第 1 行是表头）

		// 代理服务描述
		cell, _ := excelize.CoordinatesToCellName(1, rowNum)
		f.SetCellValue(sheetName, cell, info.ProxyServiceComment)
		f.SetCellStyle(sheetName, cell, cell, wrapStyle)
		if len(info.ProxyServiceComment) > colMaxLen[0] {
			colMaxLen[0] = len(info.ProxyServiceComment)
		}

		// 服务器名称
		cell, _ = excelize.CoordinatesToCellName(2, rowNum)
		f.SetCellValue(sheetName, cell, info.ServerName)
		f.SetCellStyle(sheetName, cell, cell, wrapStyle)
		if len(info.ServerName) > colMaxLen[1] {
			colMaxLen[1] = len(info.ServerName)
		}

		// 服务器端口
		portStr := fmt.Sprintf("%d", info.ServerPort)
		cell, _ = excelize.CoordinatesToCellName(3, rowNum)
		f.SetCellValue(sheetName, cell, info.ServerPort)
		f.SetCellStyle(sheetName, cell, cell, wrapStyle)
		if len(portStr) > colMaxLen[2] {
			colMaxLen[2] = len(portStr)
		}

		// 位置
		cell, _ = excelize.CoordinatesToCellName(4, rowNum)
		f.SetCellValue(sheetName, cell, info.Location)
		f.SetCellStyle(sheetName, cell, cell, wrapStyle)
		if len(info.Location) > colMaxLen[3] {
			colMaxLen[3] = len(info.Location)
		}

		// IF 条件
		cell, _ = excelize.CoordinatesToCellName(5, rowNum)
		f.SetCellValue(sheetName, cell, info.IfCondition)
		f.SetCellStyle(sheetName, cell, cell, wrapStyle)
		if len(info.IfCondition) > colMaxLen[4] {
			colMaxLen[4] = len(info.IfCondition)
		}

		// 代理原始 URL
		cell, _ = excelize.CoordinatesToCellName(6, rowNum)
		f.SetCellValue(sheetName, cell, info.ProxyOriginalURL)
		f.SetCellStyle(sheetName, cell, cell, wrapStyle)
		if len(info.ProxyOriginalURL) > colMaxLen[5] {
			colMaxLen[5] = len(info.ProxyOriginalURL)
		}

		// 代理协议
		cell, _ = excelize.CoordinatesToCellName(7, rowNum)
		f.SetCellValue(sheetName, cell, info.ProxyProtocol)
		f.SetCellStyle(sheetName, cell, cell, wrapStyle)
		if len(info.ProxyProtocol) > colMaxLen[6] {
			colMaxLen[6] = len(info.ProxyProtocol)
		}

		// 代理地址
		// 预分配切片容量，避免多次内存重新分配
		addresses := make([]string, 0, len(info.ProxyAddress))
		for _, addr := range info.ProxyAddress {
			sockets := make([]string, 0, len(addr.Sockets))
			for _, socket := range addr.Sockets {
				sockets = append(sockets, fmt.Sprintf("%s:%d", socket.IPv4, socket.Port))
			}
			addresses = append(addresses, fmt.Sprintf("%s:%d(%s)", addr.DomainName, addr.Port, strings.Join(sockets, ", ")))
		}
		addressStr := strings.Join(addresses, "\n")
		cell, _ = excelize.CoordinatesToCellName(8, rowNum)
		f.SetCellValue(sheetName, cell, addressStr)
		f.SetCellStyle(sheetName, cell, cell, wrapStyle)
		if len(addressStr) > colMaxLen[7] {
			colMaxLen[7] = len(addressStr)
		}
	}

	// 设置自适应列宽（根据实际数据长度 + 表头长度）
	colLetters := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	for i, colLetter := range colLetters {
		// 计算该列的最大宽度：取表头和数据的最大长度
		maxLen := len(headers[i])
		if colMaxLen[i] > maxLen {
			maxLen = colMaxLen[i]
		}

		// 设置列宽：基础宽度 + 额外缓冲空间（每个字符约 0.7 个 Excel 宽度单位）
		// 最小宽度 20，最大宽度 50
		// 设置上限确保换行效果不会导致列过窄
		width := float64(maxLen)*0.7 + 3
		if width < 20 {
			width = 20
		}
		if width > 50 {
			width = 50
		}

		f.SetColWidth(sheetName, colLetter, colLetter, width)
	}

	// 将 Excel 文件写入字节数组
	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func getCommentInfo(pos nginx_context.Pos) string {
	father, targetIdx := pos.Position()
	if c := pos.Target().Child(0); c.Type() == context_type.TypeInlineComment {
		return c.Value()
	}
	if n := father.Child(targetIdx + 1); n.Type() == context_type.TypeInlineComment {
		return n.Value()
	}
	if targetIdx > 0 {
		b := father.Child(targetIdx - 1)
		if b.Type() == context_type.TypeComment {
			return b.Value()
		}
	}
	return ""
}

func newWebServerStatisticsStore(store *bifrostsStore) storev1.WebServerStatisticsStore {
	wssOnce.Do(func() {
		if singletonWSSStore == nil {
			singletonWSSStore = &webServerStatisticsStore{
				bm: store.bm,
			}
		}
	})
	if singletonWSSStore == nil {
		panic(errors.New("singleton web server statistics store is nil"))
	}
	return singletonWSSStore
}
