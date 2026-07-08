package v1

import (
	mwclientv1 "github.com/tanganyu1114/heimdallr-reborn/pkg/client/v1/middleware"
	mwauthv1 "github.com/tanganyu1114/heimdallr-reborn/pkg/client/v1/middleware/auth"
	"net/http"
	"time"
)

// Config holds the configuration for the Heimdallr client
type Config struct {
	ServerAddress  string
	Timeout        time.Duration
	APIKey         string
	APISecret      string
	AuthBufferTime int64
}

// CompletedConfig holds the completed configuration
type CompletedConfig struct {
	*Config
}

// Complete completes the configuration
func (c *Config) Complete() CompletedConfig {
	return CompletedConfig{c}
}

// NewClient creates a new client from completed configuration
func (cc CompletedConfig) NewClient() (*Client, error) {
	if cc.APIKey != "" && cc.APISecret != "" {
		bt := cc.AuthBufferTime
		if bt <= 0 {
			bt = 60
		}
		mwclientv1.Register(mwauthv1.MiddlewareName, mwauthv1.AuthMiddleware(cc.APIKey, cc.APISecret, bt))
	}
	return NewClient(cc.ServerAddress, &http.Client{Timeout: cc.Timeout})
}

// NewConfig creates a new Config with default values
func NewConfig() *Config {
	return &Config{
		Timeout: 30 * time.Second,
	}
}
