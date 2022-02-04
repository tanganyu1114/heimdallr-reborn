package v1

type Service interface {
	AgentInfo() AgentInfoSrv
	WebServerConfig() WebServerConfigSrv
	Group() GroupSrv
	Host() HostSrv
	Websocket() WebsocketSrv
}
