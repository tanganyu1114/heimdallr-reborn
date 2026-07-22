package bifrosts

import (
	"context"
	"encoding/json"
	"slices"
	"strings"
	"testing"

	v1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	"github.com/tanganyu1114/heimdallr-reborn/server/internal/pkg/bifrosts"
	"github.com/tanganyu1114/heimdallr-reborn/server/internal/pkg/bifrosts/fake"

	bifrostv1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	bifrostclinetv1 "github.com/ClessLi/bifrost/pkg/client/bifrost/v1"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
	"go.uber.org/mock/gomock"
)

func Test_webServerStatisticsStore_GetProxyServiceInfo(t *testing.T) {
	webSrvOpts := v1.WebServerOptions{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBifrostsManager := bifrosts.NewMockManager(ctrl)
	mockBifrostsManager.EXPECT().GetBifrostClient(webSrvOpts).AnyTimes().Return(&bifrostclinetv1.Client{Factory: new(fake.ServiceClient)}, nil)
	type fields struct {
		bm bifrosts.Manager
	}
	type args struct {
		in0  context.Context
		opts v1.WebServerOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []v1.ProxyServiceInfo
		wantErr bool
	}{
		{
			name:   "normal test",
			fields: fields{bm: mockBifrostsManager},
			args:   args{opts: webSrvOpts},
			want: []v1.ProxyServiceInfo{
				{"test1.com", 80, "/test1-location", "", "http://right_proxy", "", "HTTP",
					[]local.ProxiedAddress{{"right_proxy", 80, []*local.Socket{}, local.ToJSONError(nil)}},
					"test inline comments", []string{"test inline comments"},
					v1.ConfigContextPos{"C:\\config_test\\conf.d\\location.conf", []int{0, 1}}},
				{"test1.com", 80, "/test1-location", "($http_api_name != '')", "http://wrong_proxy", "", "HTTP",
					[]local.ProxiedAddress{{"wrong_proxy", 80, []*local.Socket{}, local.ToJSONError(nil)}},
					"", []string{},
					v1.ConfigContextPos{"C:\\config_test\\conf.d\\location.conf", []int{0, 0, 0}}},
				{"test1.com", 80, "/test-to-baidu", "", "https://www.baidu.com", "", "HTTPS",
					[]local.ProxiedAddress{{"www.baidu.com", 443, []*local.Socket{{"183.240.99.58", 443, bifrostv1.NetUnknown, bifrostv1.NetUnknown},
						{"183.240.99.169", 443, bifrostv1.NetUnknown, bifrostv1.NetUnknown}}, local.ToJSONError(nil)}},
					"", []string{},
					v1.ConfigContextPos{"C:\\config_test\\conf.d\\location.conf", []int{1, 0}}},
				{"test1.com", 8080, "/test2", "", "http://test2.test.com", "", "HTTP",
					[]local.ProxiedAddress{{"test2.test.com", 80, []*local.Socket{{"69.167.164.199", 80, bifrostv1.NetUnknown, bifrostv1.NetUnknown}},
						local.ToJSONError(nil)}},
					"", []string{},
					v1.ConfigContextPos{"C:\\config_test\\conf.d\\location1.conf", []int{0, 0}}},
				{"test1.com", 8080, "/test1-location", "", "http://right_proxy", "", "HTTP",
					[]local.ProxiedAddress{{"right_proxy", 80, []*local.Socket{}, local.ToJSONError(nil)}},
					"test inline comments", []string{"test inline comments"},
					v1.ConfigContextPos{"C:\\config_test\\conf.d\\location.conf", []int{0, 1}}},
				{"test1.com", 8080, "/test1-location", "($http_api_name != '')", "http://wrong_proxy", "", "HTTP",
					[]local.ProxiedAddress{{"wrong_proxy", 80, []*local.Socket{}, local.ToJSONError(nil)}},
					"", []string{},
					v1.ConfigContextPos{"C:\\config_test\\conf.d\\location.conf", []int{0, 0, 0}}},
				{"test1.com", 8080, "/test-to-baidu", "", "https://www.baidu.com", "", "HTTPS",
					[]local.ProxiedAddress{{"www.baidu.com", 443, []*local.Socket{{"183.240.99.58", 443, bifrostv1.NetUnknown, bifrostv1.NetUnknown},
						{"183.240.99.169", 443, bifrostv1.NetUnknown, bifrostv1.NetUnknown}}, local.ToJSONError(nil)}},
					"", []string{},
					v1.ConfigContextPos{"C:\\config_test\\conf.d\\location.conf", []int{1, 0}}},
				{"test2.com", 8080, "/test1-location", "", "http://right_proxy", "", "HTTP",
					[]local.ProxiedAddress{{"right_proxy", 80, []*local.Socket{}, local.ToJSONError(nil)}},
					"test inline comments", []string{"test inline comments"},
					v1.ConfigContextPos{"C:\\config_test\\conf.d\\location.conf", []int{0, 1}}},
				{"test2.com", 8080, "/test1-location", "($http_api_name != '')", "http://wrong_proxy", "", "HTTP",
					[]local.ProxiedAddress{{"wrong_proxy", 80, []*local.Socket{}, local.ToJSONError(nil)}},
					"", []string{},
					v1.ConfigContextPos{"C:\\config_test\\conf.d\\location.conf", []int{0, 0, 0}}},
				{"test2.com", 8080, "/test-to-baidu", "", "https://www.baidu.com", "", "HTTPS",
					[]local.ProxiedAddress{{"www.baidu.com", 443, []*local.Socket{{"183.240.99.58", 443, bifrostv1.NetUnknown, bifrostv1.NetUnknown},
						{"183.240.99.169", 443, bifrostv1.NetUnknown, bifrostv1.NetUnknown}}, local.ToJSONError(nil)}},
					"", []string{},
					v1.ConfigContextPos{"C:\\config_test\\conf.d\\location.conf", []int{1, 0}}},
				{"test2.com", 8080, "/test2", "", "http://test2.test.com", "", "HTTP",
					[]local.ProxiedAddress{{"test2.test.com", 80, []*local.Socket{{"69.167.164.199", 80, bifrostv1.NetUnknown, bifrostv1.NetUnknown}},
						local.ToJSONError(nil)}},
					"", []string{},
					v1.ConfigContextPos{"C:\\config_test\\conf.d\\location1.conf", []int{0, 0}}},
				{"", 8888, "", "", "baidu.com:22", "", "TCP/UDP",
					[]local.ProxiedAddress{{"baidu.com", 22, []*local.Socket{{"124.237.177.164", 22, 0, 0},
						{"110.242.74.102", 22, 0, 0},
						{"111.63.65.103", 22, 0, 0}}, local.ToJSONError(nil)}},
					"", []string{},
					v1.ConfigContextPos{"C:\\config_test\\nginx.conf", []int{12, 0, 1}}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerStatisticsStore{
				bm: tt.fields.bm,
			}
			got, err := w.GetProxyServiceInfo(tt.args.in0, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProxyServiceInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			slices.SortFunc(got, func(a, b v1.ProxyServiceInfo) int {
				return strings.Compare(a.ProxyProtocol, b.ProxyProtocol)
			})
			slices.SortFunc(got, func(a, b v1.ProxyServiceInfo) int {
				return strings.Compare(a.IfCondition, b.IfCondition)
			})
			slices.SortFunc(got, func(a, b v1.ProxyServiceInfo) int {
				return strings.Compare(a.Location, b.Location)
			})
			slices.SortFunc(got, func(a, b v1.ProxyServiceInfo) int {
				return strings.Compare(a.ServerName, b.ServerName)
			})
			slices.SortFunc(got, func(a, b v1.ProxyServiceInfo) int {
				return a.ServerPort - b.ServerPort
			})
			slices.SortFunc(tt.want, func(a, b v1.ProxyServiceInfo) int {
				return strings.Compare(a.ProxyProtocol, b.ProxyProtocol)
			})
			slices.SortFunc(tt.want, func(a, b v1.ProxyServiceInfo) int {
				return strings.Compare(a.IfCondition, b.IfCondition)
			})
			slices.SortFunc(tt.want, func(a, b v1.ProxyServiceInfo) int {
				return strings.Compare(a.Location, b.Location)
			})
			slices.SortFunc(tt.want, func(a, b v1.ProxyServiceInfo) int {
				return strings.Compare(a.ServerName, b.ServerName)
			})
			slices.SortFunc(tt.want, func(a, b v1.ProxyServiceInfo) int {
				return a.ServerPort - b.ServerPort
			})
			// Sort sockets within each ProxyServiceInfo to handle DNS resolution order changes
			for i := range got {
				slices.SortFunc(got[i].ProxyAddress, func(a, b local.ProxiedAddress) int {
					return strings.Compare(a.DomainName, b.DomainName)
				})
				for j := range got[i].ProxyAddress {
					slices.SortFunc(got[i].ProxyAddress[j].Sockets, func(a, b *local.Socket) int {
						return strings.Compare(a.IPv4, b.IPv4)
					})
				}
			}
			for i := range tt.want {
				slices.SortFunc(tt.want[i].ProxyAddress, func(a, b local.ProxiedAddress) int {
					return strings.Compare(a.DomainName, b.DomainName)
				})
				for j := range tt.want[i].ProxyAddress {
					slices.SortFunc(tt.want[i].ProxyAddress[j].Sockets, func(a, b *local.Socket) int {
						return strings.Compare(a.IPv4, b.IPv4)
					})
				}
			}
			// Compare using a custom comparator that ignores dynamic DNS IPs
			if len(got) != len(tt.want) {
				gotJson, _ := json.Marshal(got)
				wantJson, _ := json.Marshal(tt.want)
				t.Errorf("GetProxyServiceInfo() length mismatch: got %d, want %d\ngot = %s\nwant = %s", len(got), len(tt.want), gotJson, wantJson)
				return
			}
			for i := range got {
				if got[i].ServerName != tt.want[i].ServerName ||
					got[i].ServerPort != tt.want[i].ServerPort ||
					got[i].Location != tt.want[i].Location ||
					got[i].IfCondition != tt.want[i].IfCondition ||
					got[i].ProxyOriginalURL != tt.want[i].ProxyOriginalURL ||
					got[i].ProxyURI != tt.want[i].ProxyURI ||
					got[i].ProxyProtocol != tt.want[i].ProxyProtocol ||
					got[i].ProxyServiceComment != tt.want[i].ProxyServiceComment ||
					got[i].ContextPos.Config != tt.want[i].ContextPos.Config ||
					!slices.Equal(got[i].ContextPos.ContextPosPath, tt.want[i].ContextPos.ContextPosPath) {
					gotJson, _ := json.Marshal(got)
					wantJson, _ := json.Marshal(tt.want)
					t.Errorf("GetProxyServiceInfo() mismatch at index %d:\ngot = %s\nwant = %s", i, gotJson, wantJson)
					return
				}
				// Compare ProxyAddress: check domain, port, and socket count (ignore specific IPs)
				if len(got[i].ProxyAddress) != len(tt.want[i].ProxyAddress) {
					gotJson, _ := json.Marshal(got)
					wantJson, _ := json.Marshal(tt.want)
					t.Errorf("GetProxyServiceInfo() ProxyAddress length mismatch at index %d:\ngot = %s\nwant = %s", i, gotJson, wantJson)
					return
				}
				for j := range got[i].ProxyAddress {
					if got[i].ProxyAddress[j].DomainName != tt.want[i].ProxyAddress[j].DomainName ||
						got[i].ProxyAddress[j].Port != tt.want[i].ProxyAddress[j].Port ||
						len(got[i].ProxyAddress[j].Sockets) != len(tt.want[i].ProxyAddress[j].Sockets) {
						gotJson, _ := json.Marshal(got)
						wantJson, _ := json.Marshal(tt.want)
						t.Errorf("GetProxyServiceInfo() ProxyAddress mismatch at index %d.%d:\ngot = %s\nwant = %s", i, j, gotJson, wantJson)
						return
					}
				}
			}
		})
	}
}

func Test_webServerStatisticsStore_ConnectivityCheckOfProxyService(t *testing.T) {
	webSrvOpts := v1.WebServerOptions{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBifrostsManager := bifrosts.NewMockManager(ctrl)
	mockBifrostsManager.EXPECT().GetBifrostClient(webSrvOpts).AnyTimes().Return(&bifrostclinetv1.Client{Factory: new(fake.ServiceClient)}, nil)
	httpPos := v1.ConfigContextPos{Config: "C:\\config_test\\conf.d\\location.conf", ContextPosPath: []int{1, 0}}
	streamPos := v1.ConfigContextPos{Config: "C:\\config_test\\nginx.conf", ContextPosPath: []int{12, 0, 1}}
	wrongPos := v1.ConfigContextPos{Config: "C:\\config_test\\nginx.conf", ContextPosPath: []int{11, 0, 0}}
	type fields struct {
		bm bifrosts.Manager
	}
	type args struct {
		ctx          context.Context
		opts         v1.WebServerOptions
		proxyPassPos v1.ConfigContextPos
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantInfo v1.ProxyServiceInfo
		wantErr  bool
	}{
		{
			name:   "http proxy pass check test",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				ctx:          context.Background(),
				opts:         webSrvOpts,
				proxyPassPos: httpPos,
			},
			wantInfo: v1.ProxyServiceInfo{
				ProxyOriginalURL: "https://www.baidu.com",
				ProxyURI:         "",
				ProxyProtocol:    "HTTPS",
				ProxyAddress: []local.ProxiedAddress{{DomainName: "www.baidu.com", Port: 443, Sockets: []*local.Socket{{"183.240.99.58", 443, bifrostv1.NetUnknown, bifrostv1.NetUnknown},
					{"183.240.99.169", 443, bifrostv1.NetUnknown, bifrostv1.NetUnknown}}, ResolveErr: local.ToJSONError(nil)}},
				ContextPos: v1.ConfigContextPos{Config: "C:\\config_test\\conf.d\\location.conf", ContextPosPath: []int{1, 0}},
			},
		},
		{
			name:   "stream proxy pass check test",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				ctx:          context.Background(),
				opts:         webSrvOpts,
				proxyPassPos: streamPos,
			},
			wantInfo: v1.ProxyServiceInfo{
				ProxyOriginalURL: "baidu.com:22",
				ProxyURI:         "",
				ProxyProtocol:    "TCP/UDP",
				ProxyAddress: []local.ProxiedAddress{{DomainName: "baidu.com", Port: 22,
					Sockets: []*local.Socket{{"124.237.177.164", 22, bifrostv1.NetUnknown, bifrostv1.NetUnknown},
						{"110.242.74.102", 22, bifrostv1.NetUnknown, bifrostv1.NetUnknown},
						{"111.63.65.103", 22, bifrostv1.NetUnknown, bifrostv1.NetUnknown}}, ResolveErr: local.ToJSONError(nil)}},
				ContextPos: v1.ConfigContextPos{Config: "C:\\config_test\\nginx.conf", ContextPosPath: []int{12, 0, 1}},
			},
		},
		{
			name:   "wrong proxy pass check test",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				ctx:          context.Background(),
				opts:         webSrvOpts,
				proxyPassPos: wrongPos,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerStatisticsStore{
				bm: tt.fields.bm,
			}
			gotInfo, err := w.ConnectivityCheckOfProxyService(tt.args.ctx, tt.args.opts, tt.args.proxyPassPos)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConnectivityCheckOfProxyService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Sort sockets within ProxyAddress to handle DNS resolution order changes
			for i := range gotInfo.ProxyAddress {
				slices.SortFunc(gotInfo.ProxyAddress[i].Sockets, func(a, b *local.Socket) int {
					return strings.Compare(a.IPv4, b.IPv4)
				})
			}
			for i := range tt.wantInfo.ProxyAddress {
				slices.SortFunc(tt.wantInfo.ProxyAddress[i].Sockets, func(a, b *local.Socket) int {
					return strings.Compare(a.IPv4, b.IPv4)
				})
			}
			// Compare using a custom comparator that ignores dynamic DNS IPs
			if gotInfo.ProxyOriginalURL != tt.wantInfo.ProxyOriginalURL ||
				gotInfo.ProxyURI != tt.wantInfo.ProxyURI ||
				gotInfo.ProxyProtocol != tt.wantInfo.ProxyProtocol ||
				gotInfo.ProxyServiceComment != tt.wantInfo.ProxyServiceComment ||
				gotInfo.ContextPos.Config != tt.wantInfo.ContextPos.Config ||
				!slices.Equal(gotInfo.ContextPos.ContextPosPath, tt.wantInfo.ContextPos.ContextPosPath) {
				gotJson, _ := json.Marshal(gotInfo)
				wantJson, _ := json.Marshal(tt.wantInfo)
				t.Errorf("ConnectivityCheckOfProxyService() mismatch:\ngot = %s\nwant = %s", gotJson, wantJson)
				return
			}
			// Compare ProxyAddress: check domain, port, and socket count (ignore specific IPs)
			if len(gotInfo.ProxyAddress) != len(tt.wantInfo.ProxyAddress) {
				gotJson, _ := json.Marshal(gotInfo)
				wantJson, _ := json.Marshal(tt.wantInfo)
				t.Errorf("ConnectivityCheckOfProxyService() ProxyAddress length mismatch:\ngot = %s\nwant = %s", gotJson, wantJson)
				return
			}
			for i := range gotInfo.ProxyAddress {
				if gotInfo.ProxyAddress[i].DomainName != tt.wantInfo.ProxyAddress[i].DomainName ||
					gotInfo.ProxyAddress[i].Port != tt.wantInfo.ProxyAddress[i].Port ||
					len(gotInfo.ProxyAddress[i].Sockets) != len(tt.wantInfo.ProxyAddress[i].Sockets) {
					gotJson, _ := json.Marshal(gotInfo)
					wantJson, _ := json.Marshal(tt.wantInfo)
					t.Errorf("ConnectivityCheckOfProxyService() ProxyAddress mismatch at index %d:\ngot = %s\nwant = %s", i, gotJson, wantJson)
					return
				}
			}
		})
	}
}
