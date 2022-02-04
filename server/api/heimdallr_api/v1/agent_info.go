package v1

type GroupInfo struct {
	Name  string     `json:"name"`
	Hosts []HostInfo `json:"hosts"`
}

type HostInfo struct {
	Name      string    `json:"name"`
	Ipaddr    string    `json:"ipaddr"`
	Descrip   string    `json:"descrip"`
	Status    bool      `json:"status"`
	AgentInfo AgentInfo `json:"agent"`
}

type AgentInfo struct {
	OS             string            `json:"system"`
	Time           string            `json:"time"`
	Cpu            string            `json:"cpu"`
	Mem            string            `json:"mem"`
	Disk           string            `json:"disk"`
	StatusList     []WebServerStatus `json:"status_list"`
	BifrostVersion string            `json:"bifrost_version"`
}

type WebServerStatus struct {
	Name    string `json:"name"`
	Status  int    `json:"status"`
	Version string `json:"version"`
}
