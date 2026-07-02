package transport

import (
	"net/http"
	"testing"
)

func TestNewTransport(t *testing.T) {
	type args struct {
		client  *http.Client
		baseURL string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "creates transport with valid client and base URL",
			args: args{
				client:  &http.Client{},
				baseURL: "http://localhost:8080",
			},
		},
		{
			name: "creates transport with nil client",
			args: args{
				client:  nil,
				baseURL: "http://localhost:8080",
			},
		},
		{
			name: "creates transport with empty base URL",
			args: args{
				client:  &http.Client{},
				baseURL: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTransport(tt.args.client, tt.args.baseURL)
			if tt.args.baseURL == "" {
				if err == nil {
					t.Errorf("NewTransport() expected error for empty baseURL, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("NewTransport() unexpected error: %v", err)
			}
			if got == nil {
				t.Errorf("NewTransport() = nil, want non-nil")
			}
		})
	}
}

func Test_transport_AgentInfos(t1 *testing.T) {
	tests := []struct {
		name    string
		baseURL string
	}{
		{
			name:    "returns agent info transport singleton",
			baseURL: "http://localhost:8080",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &transport{
				Client:  &http.Client{},
				baseURL: tt.baseURL,
			}
			got := t.AgentInfos()
			if got == nil {
				t1.Errorf("AgentInfos() = nil, want non-nil")
			}
		})
	}
}

func Test_transport_Groups(t1 *testing.T) {
	tests := []struct {
		name    string
		baseURL string
	}{
		{
			name:    "returns group transport singleton",
			baseURL: "http://localhost:8080",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &transport{
				Client:  &http.Client{},
				baseURL: tt.baseURL,
			}
			got := t.Groups()
			if got == nil {
				t1.Errorf("Groups() = nil, want non-nil")
			}
		})
	}
}

func Test_transport_Hosts(t1 *testing.T) {
	tests := []struct {
		name    string
		baseURL string
	}{
		{
			name:    "returns host transport singleton",
			baseURL: "http://localhost:8080",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &transport{
				Client:  &http.Client{},
				baseURL: tt.baseURL,
			}
			got := t.Hosts()
			if got == nil {
				t1.Errorf("Hosts() = nil, want non-nil")
			}
		})
	}
}

func Test_transport_SysUsers(t1 *testing.T) {
	tests := []struct {
		name    string
		baseURL string
	}{
		{
			name:    "returns sys user transport singleton",
			baseURL: "http://localhost:8080",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &transport{
				Client:  &http.Client{},
				baseURL: tt.baseURL,
			}
			got := t.SysUsers()
			if got == nil {
				t1.Errorf("SysUsers() = nil, want non-nil")
			}
		})
	}
}

func Test_transport_WebServerBinCMDs(t1 *testing.T) {
	tests := []struct {
		name    string
		baseURL string
	}{
		{
			name:    "returns web server bin cmd transport singleton",
			baseURL: "http://localhost:8080",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &transport{
				Client:  &http.Client{},
				baseURL: tt.baseURL,
			}
			got := t.WebServerBinCMDs()
			if got == nil {
				t1.Errorf("WebServerBinCMDs() = nil, want non-nil")
			}
		})
	}
}

func Test_transport_WebServerConfigs(t1 *testing.T) {
	tests := []struct {
		name    string
		baseURL string
	}{
		{
			name:    "returns web server config transport singleton",
			baseURL: "http://localhost:8080",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &transport{
				Client:  &http.Client{},
				baseURL: tt.baseURL,
			}
			got := t.WebServerConfigs()
			if got == nil {
				t1.Errorf("WebServerConfigs() = nil, want non-nil")
			}
		})
	}
}

func Test_transport_WebServerStatistics(t1 *testing.T) {
	tests := []struct {
		name    string
		baseURL string
	}{
		{
			name:    "returns web server statistics transport singleton",
			baseURL: "http://localhost:8080",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &transport{
				Client:  &http.Client{},
				baseURL: tt.baseURL,
			}
			got := t.WebServerStatistics()
			if got == nil {
				t1.Errorf("WebServerStatistics() = nil, want non-nil")
			}
		})
	}
}
