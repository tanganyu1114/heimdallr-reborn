package bifrosts

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	storev1 "gin-vue-admin/internal/hmdr_api/store/v1"
	storev1utils "gin-vue-admin/internal/hmdr_api/store/v1/utils"
	"gin-vue-admin/internal/pkg/bifrosts"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	"strings"
	"sync"

	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
	nginx_context "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"github.com/marmotedu/errors"
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
						} else {
							return srvPos, snPos.Target().Error()
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
		return nil, err
	}
	conf, _, err := bc.WebServerConfig().Get(opts.ServerName)
	if err != nil {
		return nil, err
	}

	//var proxySvcInfos = make([]v1.ProxyServiceInfo, 0)
	pipeline := httpAndStreamProxyBriefHandlerPipe()
	infos, err := pipeline.Handle(conf.Main(), v1.ProxyServiceInfo{})
	if err != nil {
		return nil, err
	}

	return infos, nil
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
