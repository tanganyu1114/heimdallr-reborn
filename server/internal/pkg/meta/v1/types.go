package v1

import (
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	utilsV3 "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/utils"
)

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
	ServerName string `json:"srv_name" binding:"required"`
}

type WebServerLogOptions struct {
	WebServerOptions    `json:",inline"`
	LogName             string `json:"log_name" binding:"required"`
	FilteringRegexpRule string `json:"filtering_regexp_rule"`
}

type WebServerConfig struct {
	Config               local.MainContext          `json:"config"`
	OriginalFingerprints utilsV3.ConfigFingerprints `json:"original-fingerprints"`
}

type WebServerConfigContextPosSearchOptions struct {
	WebServerOptions     `json:"web-server-options" binding:"required"`
	SearchKeywordsMeta   `json:",inline"`
	OriginalFingerprints utilsV3.ConfigFingerprints `json:"original-fingerprints" binding:"required"`
}

type WebServerConfigTargetContextOptions struct {
	WebServerOptions     `json:",inline"`
	ConfigContextPos     `json:",inline"`
	OriginalFingerprints utilsV3.ConfigFingerprints `json:"original-fingerprints" binding:"required"`
}

type ConnectivityCheckOfProxiedServersRequestOptions struct {
	WebServerOptions `json:",inline"`
	ConfigContextPos `json:",inline"`
}

type WebServerConfigUpdateOptions struct {
	WebServerOptions     `json:"web-server-options" binding:"required"`
	Data                 []byte                     `json:"data" binding:"required"`
	OriginalFingerprints utilsV3.ConfigFingerprints `json:"original-fingerprints" binding:"required"`
}

type WebServerConfigContextUpdateOptions[TargetContextMeta CloneConfigContextMeta | NewConfigContextMeta | ConfigContextEnabledStateMeta] struct {
	WebServerOptions                              `json:"web-server-options" binding:"required"`
	TargetConfigContextOptions[TargetContextMeta] `json:"target-config-context-options" binding:"required"`
	DisableTheTarget                              bool                       `json:"disable-the-target"`
	OriginalFingerprints                          utilsV3.ConfigFingerprints `json:"original-fingerprints" binding:"required"`
}

type TargetConfigContextOptions[TargetContextMeta CloneConfigContextMeta | NewConfigContextMeta | ConfigContextEnabledStateMeta] struct {
	Position      ConfigContextPos  `json:"position" binding:"required"`
	TargetContext TargetContextMeta `json:"target-context" binding:"required"`
}

type ConfigContextPos struct {
	Config         string `json:"config" binding:"required"`
	ContextPosPath []int  `json:"context-pos-path"`
}

type SearchKeywordsMeta struct {
	StartingPositionList []ConfigContextPos `json:"starting-position-list" binding:"required"`
	Keywords             string             `json:"keywords" binding:"required"`
	IsRegexpRule         bool               `json:"is-regexp-rule"`
	IsOnlyInCurrent      bool               `json:"is-only-in-current"`
}

type NewConfigContextMeta struct {
	Enabled             bool                     `json:"enabled"`
	ContextType         context_type.ContextType `json:"context-type" binding:"required"`
	ContextValue        string                   `json:"context-value"`
	ChildrenContextMeta []NewConfigContextMeta   `json:"children-context-meta"`
}

type CloneConfigContextMeta struct {
	ConfigContextPos `json:"clone-context-pos" binding:"required"`
}

type ConfigContextEnabledStateMeta struct {
	Enabled bool `json:"enabled"`
}

type WebServerBinCMDExecRequest struct {
	WebServerOptions `json:"web-server-options" binding:"required"`
	Args             []string `json:"args"`
}

type WebServerBinCMDExecResponse struct {
	Successful bool   `json:"successful"`
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
}
