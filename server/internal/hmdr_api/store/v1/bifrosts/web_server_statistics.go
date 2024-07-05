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
	"sync"
)

var (
	wssOnce           = new(sync.Once)
	singletonWSSStore *webServerStatisticsStore

	regexpProxyPassValue = regexp.MustCompile(`^proxy_pass\s+(.+)$`)
	regexpHttpURL        = regexp.MustCompile(`^http://([^/]+)(/\S*)?$`)
	regexpHttpsURL       = regexp.MustCompile(`^https://([^/]+)(/\S*)?$`)
	regexpOtherURL       = regexp.MustCompile(`^\S+://([^/]+)(/\S*)?$`)
	regexpUpstream       = regexp.MustCompile(`^([^/:]+)(/\S*)?$`)
	regexpAddrWithPort   = regexp.MustCompile(`^(\S+):(\d+)`)
	regexpIPAddr         = regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)

	regexpUpstreamSrvDirective = regexp.MustCompile(`^server\s+(\S+)`)
)

type webServerStatisticsStore struct {
	bm bifrosts.Manager
}

func (w *webServerStatisticsStore) GetProxyServiceInfo(ctx context.Context, opts metav1.WebServerOptions) ([]v1.ProxyServiceInfo, error) {
	bc, err := w.bm.GetBifrostClient(opts)
	if err != nil {
		return nil, err
	}
	confData, err := bc.WebServerConfig().Get(opts.ServerName)
	if err != nil {
		return nil, err
	}
	conf, err := configuration.NewNginxConfigFromJsonBytes(confData)
	if err != nil {
		return nil, err
	}

	var proxySvcInfos []v1.ProxyServiceInfo
	// Get Http Server Proxy Brief
	httpSrvsProxyBrief, err := getHttpSrvsProxyBrief(conf)
	if err != nil && !errors.IsCode(err, 110007) {
		return nil, err
	}
	proxySvcInfos = append(proxySvcInfos, httpSrvsProxyBrief...)

	// Get Stream Server Proxy Brief
	streamSrvsProxyBrief, err := getStreamSrvsProxyBrief(conf)
	if err != nil && !errors.IsCode(err, 110007) {
		return nil, err
	}
	proxySvcInfos = append(proxySvcInfos, streamSrvsProxyBrief...)

	return proxySvcInfos, nil
}

func getStreamSrvsProxyBrief(conf configuration.NginxConfig) ([]v1.ProxyServiceInfo, error) {
	var streamSrvsProxyBrief []v1.ProxyServiceInfo
	stream := conf.Main().QueryByKeyWords(nginx_context.NewKeyWords(context_type.TypeStream)).Target()

	for _, pos := range stream.QueryAllByKeyWords(nginx_context.NewKeyWords(context_type.TypeServer)) {
		// get listening port
		port := configuration.Port(pos.Target())

		proxyPass := pos.Target().QueryByKeyWords(
			nginx_context.NewKeyWords(context_type.TypeDirective).
				SetRegexpMatchingValue(`^proxy_pass\s+`).
				SetCascaded(false),
		).Target()

		if err := proxyPass.Error(); err != nil {
			if !errors.Is(err, nginx_context.NullPos().Target().Error()) {
				global.GVA_LOG.Warn("检索proxy_pass配置项异常", zap.Any("err", err))
			}
			continue
		}

		proxyAddrs := getProxyAddresses(stream, proxyPass)
		if len(proxyAddrs) == 0 {
			global.GVA_LOG.Warn("获取Stream Server反向代理配置异常", zap.Any("listen_port", port), zap.Any("proxy_pass_value", proxyPass.Value()))
			continue
		}

		for _, proxyAddr := range proxyAddrs {
			streamSrvsProxyBrief = append(streamSrvsProxyBrief, v1.ProxyServiceInfo{
				ProxyType:    "TCP/UDP",
				Port:         port,
				ProxyAddress: proxyAddr,
			})
		}

	}

	return streamSrvsProxyBrief, nil
}

func getHttpSrvsProxyBrief(conf configuration.NginxConfig) ([]v1.ProxyServiceInfo, error) {
	var httpSrvsProxyBrief []v1.ProxyServiceInfo
	http := conf.Main().QueryByKeyWords(nginx_context.NewKeyWords(context_type.TypeHttp)).Target()
	for _, serverPos := range http.QueryAllByKeyWords(nginx_context.NewKeyWords(context_type.TypeServer)) {
		// get server name
		srvnameDirective := serverPos.Target().QueryByKeyWords(
			nginx_context.NewKeyWords(context_type.TypeDirective).
				SetRegexpMatchingValue(nginx_context.RegexpMatchingServerNameValue).
				SetCascaded(false),
		).Target()
		if err := srvnameDirective.Error(); err != nil {
			if !errors.Is(err, nginx_context.NullPos().Target().Error()) {
				global.GVA_LOG.Warn("未检索到Http Server的server_name配置项")
			} else {
				global.GVA_LOG.Warn("检索Http Server的server_name配置项失败", zap.Any("err", err))
			}
			continue
		}
		srvname := srvnameDirective.(*local.Directive).Params

		// get listening port
		port := configuration.Port(serverPos.Target())

		for _, locationPos := range serverPos.Target().QueryAllByKeyWords(nginx_context.NewKeyWords(context_type.TypeLocation)) {
			proxyPass := locationPos.Target().QueryByKeyWords(
				nginx_context.NewKeyWords(context_type.TypeDirective).
					SetRegexpMatchingValue(`^proxy_pass\s+`).
					SetCascaded(false),
			).Target()

			if err := proxyPass.Error(); err != nil {
				if a, ok := err.(errors.Aggregate); ok {
					if len(a.Errors()) == 1 {
						err = a.Errors()[0]
					}
				}
				if !errors.IsCode(err, 110017) {
					global.GVA_LOG.Warn("获取proxy_pass配置项异常", zap.Any("err", err))
				}
				continue
			}

			proxyAddrs := getProxyAddresses(http, proxyPass)
			if len(proxyAddrs) == 0 {
				global.GVA_LOG.Warn("获取location反向代理配置异常", zap.Any("location", locationPos.Target().Value()), zap.Any("proxy_pass_value", proxyPass.Value()))
				continue
			}

			for _, proxyAddr := range proxyAddrs {
				httpSrvsProxyBrief = append(httpSrvsProxyBrief, v1.ProxyServiceInfo{
					ProxyType:    "HTTP",
					ServerName:   srvname,
					Port:         port,
					Location:     locationPos.Target().Value(),
					ProxyAddress: proxyAddr,
				})

			}
		}

	}
	return httpSrvsProxyBrief, nil
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

	var addrs []string

	var port string
	var defaultPort string
	var address string
	var upstreamName string
	// get default port, address and upstream name
	switch {
	case regexpHttpsURL.MatchString(url):
		defaultPort = "443"
		address = regexpHttpsURL.FindStringSubmatch(url)[1]
	case regexpHttpURL.MatchString(url):
		defaultPort = "80"
		address = regexpHttpURL.FindStringSubmatch(url)[1]
	case regexpOtherURL.MatchString(url):
		defaultPort = "unknown_port"
		address = regexpOtherURL.FindStringSubmatch(url)[1]
	case regexpAddrWithPort.MatchString(url):
		address = regexpAddrWithPort.FindStringSubmatch(url)[1] + ":" + regexpAddrWithPort.FindStringSubmatch(url)[2]
	case regexpUpstream.MatchString(url):
		upstreamName = regexpUpstream.FindStringSubmatch(url)[1]
	default:
		return nil
	}

	if upstreamName == "" {
		if regexpAddrWithPort.MatchString(address) {
			port = regexpAddrWithPort.FindStringSubmatch(address)[2]
			address = regexpAddrWithPort.FindStringSubmatch(address)[1]
		} else {
			port = defaultPort
		}

		if regexpIPAddr.MatchString(address) {
			addrs = append(addrs, address+":"+port)
			return addrs
		}
		upstreamName = address
	}

	upstream := httpOrStream.QueryByKeyWords(
		nginx_context.NewKeyWords(context_type.TypeUpstream).
			SetStringMatchingValue(upstreamName),
	).Target()
	if upstream == nginx_context.NullPos().Target() {
		addrs = append(addrs, upstreamName+":"+port)
		return addrs
	}

	for _, pos := range upstream.QueryAllByKeyWords(
		nginx_context.NewKeyWords(context_type.TypeDirective).
			SetRegexpMatchingValue(`^server\s+`),
	) {
		if regexpUpstreamSrvDirective.MatchString(pos.Target().Value()) {
			srvAddr := regexpUpstreamSrvDirective.FindStringSubmatch(pos.Target().Value())[1]
			//addrs = append(addrs, parseAddresses(httpOrStream, srvAddr)...)  // upstream 不能嵌套调用
			addrs = append(addrs, srvAddr)
			return addrs
		}
	}

	addrs = append(addrs, address)
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
