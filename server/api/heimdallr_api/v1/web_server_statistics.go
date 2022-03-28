package v1

//type LocationProxyInfo struct {
//	Location     string   `json:"location"`
//	ProxyAddress []string `json:"proxy-address"`
//}
//
//type HttpSrvProxyBrief struct {
//	ServerName string              `json:"server-name"`
//	Port       int                 `json:"port"`
//	ProxyInfos []LocationProxyInfo `json:"proxy-infos"`
//}
//
//type FourthLevelProxyInfo struct {
//	ProxyAddress []string `json:"proxy-address"`
//}
//
//type StreamSrvProxyBrief struct {
//	Port                 int `json:"port"`
//	FourthLevelProxyInfo `json:",inline"`
//}
//
//type ProxyServiceInfoStatistics struct {
//	//bifrostapiv1.Statistics `json:",inline"`
//	HttpSrvsProxyBrief   []HttpSrvProxyBrief   `json:"http-srvs-proxy-brief"`
//	StreamSrvsProxyBrief []StreamSrvProxyBrief `json:"stream-srvs-proxy-brief"`
//}

type ProxyServiceInfo struct {
	ProxyType    string `json:"proxy-type"`
	ServerName   string `json:"server-name"`
	Port         int    `json:"port"`
	Location     string `json:"location"`
	ProxyAddress string `json:"proxy-address"`
}
