package bifrosts

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	"gin-vue-admin/global"
	storev1 "gin-vue-admin/internal/hmdr_api/store/v1"
	"gin-vue-admin/internal/pkg/bifrosts"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
	nginx_context "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"github.com/marmotedu/errors"
	"go.uber.org/zap"
	"regexp"
	"strings"
	"sync"
)

var (
	wssOnce           = new(sync.Once)
	singletonWSSStore *webServerStatisticsStore

	regexpProxyPassValue = regexp.MustCompile(`^proxy_pass\s+(.+)$`)
	regexpHttpURL        = regexp.MustCompile(`^http://([^/]+)(/\S*)?$`)
	regexpHttpsURL       = regexp.MustCompile(`^https://([^/]+)(/\S*)?$`)
	regexpOtherURL       = regexp.MustCompile(`^(?:[^/:]+?://)?([^/\s]+)(/\S*)?$`)
	regexpAddrWithPort   = regexp.MustCompile(`^(\S+):(\d+)`)
	regexpIPAddr         = regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)

	regexpUpstreamSrvDirective = regexp.MustCompile(`^server\s+(\S+)`)
)

type webServerStatisticsStore struct {
	bm bifrosts.Manager
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

	var proxySvcInfos = make([]v1.ProxyServiceInfo, 0)
	// Get Http Server Proxy Brief
	httpSrvsProxyBrief, err := getHttpSrvsProxyBrief(conf)
	if err != nil && !errors.Is(err, nginx_context.NotFoundPos().Target().Error().(errors.Aggregate).Errors()[0]) {
		return proxySvcInfos, err
	}
	// Get Stream Server Proxy Brief
	streamSrvsProxyBrief, err := getStreamSrvsProxyBrief(conf)
	if err != nil && !errors.Is(err, nginx_context.NotFoundPos().Target().Error().(errors.Aggregate).Errors()[0]) {
		return proxySvcInfos, err
	}

	proxySvcInfos = append(proxySvcInfos, httpSrvsProxyBrief...)
	proxySvcInfos = append(proxySvcInfos, streamSrvsProxyBrief...)

	return proxySvcInfos, nil
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

func getStreamSrvsProxyBrief(conf configuration.NginxConfig) ([]v1.ProxyServiceInfo, error) {
	var streamSrvsProxyBrief []v1.ProxyServiceInfo
	streamPos := conf.Main().ChildrenPosSet().
		QueryOne(nginx_context.NewKeyWords(context_type.TypeStream).
			SetSkipQueryFilter(nginx_context.SkipDisabledCtxFilterFunc))
	proxyPassKW := nginx_context.NewKeyWords(context_type.TypeDirective).
		SetSkipQueryFilter(nginx_context.SkipDisabledCtxFilterFunc).
		SetRegexpMatchingValue(`^proxy_pass\s+`)
	return streamSrvsProxyBrief, streamPos.QueryAll(nginx_context.NewKeyWords(context_type.TypeServer).
		SetSkipQueryFilter(nginx_context.SkipDisabledCtxFilterFunc)).
		Filter( // filter out the proxy server
			func(srvPos nginx_context.Pos) bool { return srvPos.QueryOne(proxyPassKW).Target().Error() == nil },
		).
		Map(
			func(srvPos nginx_context.Pos) (nginx_context.Pos, error) {
				srvComment := getCommentInfo(srvPos)
				// get listening port
				port := configuration.Port(srvPos.Target())
				return srvPos, srvPos.QueryAll(proxyPassKW).
					Filter( // filter out the normal "proxy_pass" directive
						func(directivePos nginx_context.Pos) bool {
							if directivePos.Target().Error() != nil {
								global.GVA_LOG.Warn(
									"检索proxy_pass配置项异常",
									zap.Any("err", directivePos.Target().Error()),
								)
								return false
							}
							return true
						},
					).
					Map(
						func(directivePos nginx_context.Pos) (nginx_context.Pos, error) {
							proxyAddrs := getProxyAddresses(streamPos.Target(), directivePos.Target())
							if len(proxyAddrs) == 0 {
								global.GVA_LOG.Warn(
									"获取Stream Server反向代理配置异常",
									zap.Any("listen_port", port),
									zap.Any("proxy_pass_value", directivePos.Target().Value()),
								)
								return directivePos, nil
							}
							proxyComment := getCommentInfo(directivePos)
							cmts := make([]string, 0)
							if len(srvComment) > 0 {
								cmts = append(cmts, srvComment)
							}
							if len(proxyComment) > 0 {
								cmts = append(cmts, proxyComment)
							}
							for _, proxyAddr := range proxyAddrs {
								streamSrvsProxyBrief = append(streamSrvsProxyBrief, v1.ProxyServiceInfo{
									ProxyType:           "TCP/UDP",
									Port:                port,
									ProxyAddress:        proxyAddr,
									ProxyServiceComment: strings.Join(cmts, "\t|\t"),
								})
							}
							return directivePos, nil
						},
					).
					Error()
			},
		).
		Error()
}

func getHttpSrvsProxyBrief(conf configuration.NginxConfig) ([]v1.ProxyServiceInfo, error) {
	var httpSrvsProxyBrief []v1.ProxyServiceInfo
	httpPos := conf.Main().ChildrenPosSet().
		QueryOne(nginx_context.NewKeyWords(context_type.TypeHttp).
			SetSkipQueryFilter(nginx_context.SkipDisabledCtxFilterFunc))
	proxyPassKW := nginx_context.NewKeyWords(context_type.TypeDirective).
		SetSkipQueryFilter(nginx_context.SkipDisabledCtxFilterFunc).
		SetRegexpMatchingValue(`^proxy_pass\s+`)
	return httpSrvsProxyBrief, httpPos.QueryAll(nginx_context.NewKeyWords(context_type.TypeServer).
		SetSkipQueryFilter(nginx_context.SkipDisabledCtxFilterFunc)).
		Filter( // filter out the proxy server
			func(srvPos nginx_context.Pos) bool { return srvPos.QueryOne(proxyPassKW).Target().Error() == nil },
		).
		Map(
			func(srvPos nginx_context.Pos) (nginx_context.Pos, error) {
				// get server name
				var srvname string
				srvComment := getCommentInfo(srvPos)
				if p := srvPos.QueryOne(nginx_context.NewKeyWords(context_type.TypeDirective).
					SetRegexpMatchingValue(nginx_context.RegexpMatchingServerNameValue).
					SetSkipQueryFilter(nginx_context.SkipDisabledCtxFilterFunc)); p.Target().Error() == nil {
					srvname = p.Target().(*local.Directive).Params
				}

				// get listening port
				port := configuration.Port(srvPos.Target())

				return srvPos, srvPos.QueryAll(nginx_context.NewKeyWords(context_type.TypeLocation).
					SetSkipQueryFilter(nginx_context.SkipDisabledCtxFilterFunc)).
					Filter( // filter out the proxy location
						func(locPos nginx_context.Pos) bool { return locPos.QueryOne(proxyPassKW).Target().Error() == nil },
					).
					Map(
						func(locPos nginx_context.Pos) (nginx_context.Pos, error) {
							locComment := getCommentInfo(locPos)
							return locPos, locPos.QueryAll(proxyPassKW).
								Filter( // filter out the normal "proxy_pass" directive
									func(directivePos nginx_context.Pos) bool {
										if directivePos.Target().Error() != nil {
											global.GVA_LOG.Warn(
												"检索proxy_pass配置项异常",
												zap.Any("err", directivePos.Target().Error()),
											)
											return false
										}
										return true
									},
								).
								Map(
									func(directivePos nginx_context.Pos) (nginx_context.Pos, error) {
										proxyAddrs := getProxyAddresses(httpPos.Target(), directivePos.Target())
										if len(proxyAddrs) == 0 {
											global.GVA_LOG.Warn(
												"获取location反向代理配置异常",
												zap.Any("location", locPos.Target().Value()),
												zap.Any("proxy_pass_value", directivePos.Target().Value()),
											)
											return directivePos, nil
										}
										proxyComment := getCommentInfo(directivePos)
										cmts := make([]string, 0)
										if len(srvComment) > 0 {
											cmts = append(cmts, srvComment)
										}
										if len(locComment) > 0 {
											cmts = append(cmts, locComment)
										}
										if len(proxyComment) > 0 {
											cmts = append(cmts, proxyComment)
										}
										for _, proxyAddr := range proxyAddrs {
											httpSrvsProxyBrief = append(httpSrvsProxyBrief, v1.ProxyServiceInfo{
												ProxyType:           "HTTP",
												ServerName:          srvname,
												Port:                port,
												Location:            locPos.Target().Value(),
												ProxyAddress:        proxyAddr,
												ProxyServiceComment: strings.Join(cmts, "\t|\t"),
											})
										}
										return directivePos, nil
									},
								).
								Error()
						},
					).
					Error()
			},
		).
		Error()
}

func getProxyAddresses(httpOrStream, proxypass nginx_context.Context) []string {
	var url string
	if regexpProxyPassValue.MatchString(proxypass.Value()) {
		url = regexpProxyPassValue.FindStringSubmatch(proxypass.Value())[1]
	} else {
		return nil
	}

	return parseAddresses(httpOrStream, url)
}

func parseAddresses(httpOrStream nginx_context.Context, url string) []string {
	switch httpOrStream.Type() {
	case context_type.TypeHttp, context_type.TypeStream:
	default:
		global.GVA_LOG.Warn("解析代理地址异常！入参httpOrStream非http或stream上下文")
		return nil
	}

	var host string
	var port string
	var defaultProxyPassPort string
	var address string
	// get default port, address and upstream name
	switch {
	case regexpHttpsURL.MatchString(url):
		defaultProxyPassPort = "443"
		address = regexpHttpsURL.FindStringSubmatch(url)[1]
	case regexpHttpURL.MatchString(url):
		defaultProxyPassPort = "80"
		address = regexpHttpURL.FindStringSubmatch(url)[1]
	case regexpOtherURL.MatchString(url):
		defaultProxyPassPort = "unknown_port"
		address = regexpOtherURL.FindStringSubmatch(url)[1]
	default:
		return nil
	}
	// parse address to host and port
	if regexpAddrWithPort.MatchString(address) {
		host = regexpAddrWithPort.FindStringSubmatch(address)[1]
		port = regexpAddrWithPort.FindStringSubmatch(address)[2]
	} else {
		host = address
	}

	if regexpIPAddr.MatchString(host) {
		if port == "" {
			port = defaultProxyPassPort
		}
		return []string{host + ":" + port}
	}

	upstreamPos := httpOrStream.ChildrenPosSet().QueryOne(
		nginx_context.NewKeyWords(context_type.TypeUpstream).
			SetSkipQueryFilter(nginx_context.SkipDisabledCtxFilterFunc).
			SetStringMatchingValue(host),
	)
	if upstreamPos.Target().Error() != nil {
		if port == "" {
			port = defaultProxyPassPort
		}
		return []string{host + ":" + port}
	}
	var addrs []string
	existAddrMap := make(map[string]bool)
	upstreamPos.QueryAll(nginx_context.NewKeyWords(context_type.TypeDirective).
		SetRegexpMatchingValue(`^server\s+`).
		SetSkipQueryFilter(nginx_context.SkipDisabledCtxFilterFunc)).
		Map(
			func(directivePos nginx_context.Pos) (nginx_context.Pos, error) {
				if regexpUpstreamSrvDirective.MatchString(directivePos.Target().Value()) {
					srvAddr := regexpUpstreamSrvDirective.FindStringSubmatch(directivePos.Target().Value())[1]
					var srvHost, srvPort string
					if regexpAddrWithPort.MatchString(srvAddr) {
						srvHost = regexpAddrWithPort.FindStringSubmatch(srvAddr)[1]
						srvPort = regexpAddrWithPort.FindStringSubmatch(srvAddr)[2]
					} else {
						srvHost = srvAddr
						if httpOrStream.Type() == context_type.TypeHttp {
							srvPort = "80"
						}
					}
					if port != "" {
						srvPort = port
					}
					if srvPort == "" {
						srvPort = defaultProxyPassPort
					}
					addr := srvHost + ":" + srvPort
					if !existAddrMap[addr] {
						addrs = append(addrs, addr)
						existAddrMap[addr] = true
					}
				}
				return directivePos, nil
			},
		)
	return addrs
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
