package bifrosts

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	"gin-vue-admin/global"
	storev1 "gin-vue-admin/internal/hmdr_api/store/v1"
	"gin-vue-admin/internal/pkg/bifrosts"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration/parser"
	"github.com/marmotedu/errors"
	"go.uber.org/zap"
	"regexp"
	"sync"
)

var (
	wssOnce           = new(sync.Once)
	singletonWSSStore *webServerStatisticsStore

	regProxyPassValue = regexp.MustCompile(`^proxy_pass\s+(.+)$`)
	regHttpURL        = regexp.MustCompile(`^http://([^/]+)(/\S*)?$`)
	regHttpsURL       = regexp.MustCompile(`^https://([^/]+)(/\S*)?$`)
	regOtherURL       = regexp.MustCompile(`^\S+://([^/]+)(/\S*)?$`)
	regUpstream       = regexp.MustCompile(`^([^/:]+)(/\S*)?$`)
	regAddrWithPort   = regexp.MustCompile(`^(\S+):(\d+)`)
	regIPAddr         = regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)

	regUpstreamSrvKey = regexp.MustCompile(`^server\s+(\S+)`)
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
	conf, err := configuration.NewConfigurationFromJsonBytes(confData)
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

func getStreamSrvsProxyBrief(conf configuration.Configuration) ([]v1.ProxyServiceInfo, error) {
	var streamSrvsProxyBrief []v1.ProxyServiceInfo
	stream, err := conf.Query("stream")
	if err != nil {
		return nil, err
	}
	streamsrvs, err := stream.QueryAll("server")
	if err != nil {
		return nil, err
	}

	for _, streamsrv := range streamsrvs {
		// get listening port
		port := configuration.Port(streamsrv)

		proxyKey, err := streamsrv.Query(`key:sep: :reg: ^proxy_pass\s+`)
		if err != nil {
			continue
		}
		proxyAddrs := getProxyAddresses(stream, proxyKey)
		if len(proxyAddrs) == 0 {
			global.GVA_LOG.Warn("获取steam server反向代理配置异常", zap.Any("listen_port", port), zap.Any("proxy_pass_value", proxyKey.Self().GetValue()))
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

func getHttpSrvsProxyBrief(conf configuration.Configuration) ([]v1.ProxyServiceInfo, error) {
	var httpSrvsProxyBrief []v1.ProxyServiceInfo
	http, err := conf.Query("http")
	if err != nil {
		return nil, err
	}
	httpsrvs, err := http.QueryAll("server")
	if err != nil {
		return nil, err
	}
	for _, httpsrv := range httpsrvs {
		// get server name
		var srvname string
		srvnameKey, err := httpsrv.Query(`key:sep: :reg: ^server_name\s+`)
		if err != nil {
			if !errors.IsCode(err, 110007) {
				global.GVA_LOG.Warn("获取http server server_name key失败", zap.Any("err", err))
				continue
			}
		} else if configuration.RegServerNameValue.MatchString(srvnameKey.Self().GetValue()) {
			srvname = configuration.RegServerNameValue.FindStringSubmatch(srvnameKey.Self().GetValue())[1]
		} else {
			global.GVA_LOG.Warn("解析http server server_name key失败", zap.Any("server_name_value", srvnameKey.Self().GetValue()))
			continue
		}

		// get listening port
		port := configuration.Port(httpsrv)

		locations, err := httpsrv.QueryAll(`location:sep: :reg: .*`)
		if err != nil {
			global.GVA_LOG.Warn("获取http server location失败", zap.Any("err", err))
			continue
		}
		for _, location := range locations {
			proxyKey, err := location.Query(`key:sep: :reg: ^proxy_pass\s+`)
			if err != nil {
				continue
			}

			proxyAddrs := getProxyAddresses(http, proxyKey)
			if len(proxyAddrs) == 0 {
				global.GVA_LOG.Warn("获取location反向代理配置异常", zap.Any("location", location.Self().GetValue()), zap.Any("proxy_pass_value", proxyKey.Self().GetValue()))
				continue
			}

			for _, proxyAddr := range proxyAddrs {
				httpSrvsProxyBrief = append(httpSrvsProxyBrief, v1.ProxyServiceInfo{
					ProxyType:    "HTTP",
					ServerName:   srvname,
					Port:         port,
					Location:     location.Self().GetValue(),
					ProxyAddress: proxyAddr,
				})

			}
		}

	}
	return httpSrvsProxyBrief, nil
}

func getProxyAddresses(httpOrStream, proxyKey configuration.Querier) []string {
	var url string
	if regProxyPassValue.MatchString(proxyKey.Self().GetValue()) {
		url = regProxyPassValue.FindStringSubmatch(proxyKey.Self().GetValue())[1]
	} else {
		return nil
	}

	return parseAddresses(httpOrStream, url)
}

func parseAddresses(httpOrStream configuration.Querier, url string) []string {
	switch httpOrStream.Self().(type) {
	case *parser.Http, *parser.Stream:
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
	case regHttpsURL.MatchString(url):
		defaultPort = "443"
		address = regHttpsURL.FindStringSubmatch(url)[1]
	case regHttpURL.MatchString(url):
		defaultPort = "80"
		address = regHttpURL.FindStringSubmatch(url)[1]
	case regOtherURL.MatchString(url):
		defaultPort = "unknown_port"
		address = regOtherURL.FindStringSubmatch(url)[1]
	case regAddrWithPort.MatchString(url):
		address = regAddrWithPort.FindStringSubmatch(url)[1] + ":" + regAddrWithPort.FindStringSubmatch(url)[2]
	case regUpstream.MatchString(url):
		upstreamName = regUpstream.FindStringSubmatch(url)[1]
	default:
		return nil
	}

	if upstreamName == "" {
		if regAddrWithPort.MatchString(address) {
			port = regAddrWithPort.FindStringSubmatch(address)[2]
			address = regAddrWithPort.FindStringSubmatch(address)[1]
		} else {
			port = defaultPort
		}

		if regIPAddr.MatchString(address) {
			addrs = append(addrs, address+":"+port)
			return addrs
		}
		upstreamName = address
	}

	upstream, err := httpOrStream.Query(`upstream:sep: ` + upstreamName)
	if err != nil || upstream == nil {
		addrs = append(addrs, upstreamName+":"+port)
		return addrs
	}

	upstreamSrvKeys, err := upstream.QueryAll(`key:sep: :reg: ^server\s+`)
	if err != nil {
		return nil
	}

	for _, srvKey := range upstreamSrvKeys {
		if regUpstreamSrvKey.MatchString(srvKey.Self().GetValue()) {
			srvAddr := regUpstreamSrvKey.FindStringSubmatch(srvKey.Self().GetValue())[1]
			//addrs = append(addrs, parseAddresses(httpOrStream, srvAddr)...)  // upstream 不能嵌套调用
			addrs = append(addrs, srvAddr)
		}
	}
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
