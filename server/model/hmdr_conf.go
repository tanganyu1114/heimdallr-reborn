package model

type SearchConf struct {
	GroupId uint   `json:"group_id"`
	HostId  uint   `json:"host_id"`
	SrvName string `json:"srv_name"`
}
