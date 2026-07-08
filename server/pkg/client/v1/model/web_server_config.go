package model

import (
	"encoding/json"

	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
	utilsV3 "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/utils"
)

// WebServerConfig is the client-side representation of web server config.
// It uses json.RawMessage for Config to avoid unmarshaling issues with interface types.
type WebServerConfig struct {
	Config               json.RawMessage            `json:"config"`
	OriginalFingerprints utilsV3.ConfigFingerprints `json:"original-fingerprints"`
}

func (w *WebServerConfig) GetConfig() (configuration.NginxConfig, error) {
	return configuration.NewNginxConfigFromJsonBytes(w.Config)
}
