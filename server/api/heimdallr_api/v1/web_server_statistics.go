package v1

import (
	metav1 "github.com/tanganyu1114/heimdallr-reborn/internal/pkg/meta/v1"

	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
)

type ProxyServiceInfo struct {
	ServerName          string                 `json:"server-name"`
	ServerPort          int                    `json:"server-port"`
	Location            string                 `json:"location"`
	IfCondition         string                 `json:"if-condition"`
	ProxyOriginalURL    string                 `json:"proxy-original-url"`
	ProxyURI            string                 `json:"proxy-uri"`
	ProxyProtocol       string                 `json:"proxy-protocol"`
	ProxyAddress        []local.ProxiedAddress `json:"proxy-address"`
	ProxyServiceComment string                 `json:"proxy-service-comment"`
	TmpComments         []string               `json:"-"`

	ContextPos metav1.ConfigContextPos `json:"context-pos"`
}
