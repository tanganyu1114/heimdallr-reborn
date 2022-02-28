package v1

type ListMeta struct {
	TotalCount int64 `json:"totalCount,omitempty"`
}

// ListOptions is the query options to a standard REST list call.
type ListOptions struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

type IDOptions struct {
	ID float64 `json:"id" form:"id"`
}

type IDsOptions struct {
	IDs []int `json:"ids" form:"ids"`
}

type UintObjectMeta struct {
	Label string `json:"label"`
	Value uint   `json:"value"`
}

type StringObjectMeta struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type WebServerOptions struct {
	GroupID    uint   `json:"group_id"`
	HostID     uint   `json:"host_id"`
	ServerName string `json:"srv_name"`
}

type WebServerLogOptions struct {
	WebServerOptions    `json:",inline"`
	LogName             string `json:"log_name"`
	FilteringRegexpRule string `json:"filtering_regexp_rule"`
}
