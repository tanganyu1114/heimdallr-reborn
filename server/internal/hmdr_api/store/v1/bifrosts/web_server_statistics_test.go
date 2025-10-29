package bifrosts

import (
	"context"
	"encoding/json"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	"gin-vue-admin/internal/pkg/bifrosts"
	"gin-vue-admin/internal/pkg/bifrosts/fake"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	"reflect"
	"slices"
	"strings"
	"testing"

	bifrostv1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	bifrostclinetv1 "github.com/ClessLi/bifrost/pkg/client/bifrost/v1"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
	"go.uber.org/mock/gomock"
)

func Test_webServerStatisticsStore_GetProxyServiceInfo(t *testing.T) {
	webSrvOpts := metav1.WebServerOptions{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBifrostsManager := bifrosts.NewMockManager(ctrl)
	mockBifrostsManager.EXPECT().GetBifrostClient(webSrvOpts).AnyTimes().Return(&bifrostclinetv1.Client{Factory: new(fake.ServiceClient)}, nil)
	type fields struct {
		bm bifrosts.Manager
	}
	type args struct {
		in0  context.Context
		opts metav1.WebServerOptions
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
					metav1.ConfigContextPos{"C:\\config_test\\conf.d\\location.conf", []int{0, 1}}},
				{"test1.com", 80, "/test1-location", "($http_api_name != '')", "http://wrong_proxy", "", "HTTP",
					[]local.ProxiedAddress{{"wrong_proxy", 80, []*local.Socket{}, local.ToJSONError(nil)}},
					"", []string{},
					metav1.ConfigContextPos{"C:\\config_test\\conf.d\\location.conf", []int{0, 0, 0}}},
				{"test1.com", 80, "/test-to-baidu", "", "https://www.baidu.com", "", "HTTPS",
					[]local.ProxiedAddress{{"www.baidu.com", 443, []*local.Socket{{"183.240.99.58", 443, bifrostv1.NetUnknown, bifrostv1.NetUnknown},
						{"183.240.99.169", 443, bifrostv1.NetUnknown, bifrostv1.NetUnknown}}, local.ToJSONError(nil)}},
					"", []string{},
					metav1.ConfigContextPos{"C:\\config_test\\conf.d\\location.conf", []int{1, 0}}},
				{"test1.com", 8080, "/test2", "", "http://test2.test.com", "", "HTTP",
					[]local.ProxiedAddress{{"test2.test.com", 80, []*local.Socket{{"69.167.164.199", 80, bifrostv1.NetUnknown, bifrostv1.NetUnknown}},
						local.ToJSONError(nil)}},
					"", []string{},
					metav1.ConfigContextPos{"C:\\config_test\\conf.d\\location1.conf", []int{0, 0}}},
				{"test1.com", 8080, "/test1-location", "", "http://right_proxy", "", "HTTP",
					[]local.ProxiedAddress{{"right_proxy", 80, []*local.Socket{}, local.ToJSONError(nil)}},
					"test inline comments", []string{"test inline comments"},
					metav1.ConfigContextPos{"C:\\config_test\\conf.d\\location.conf", []int{0, 1}}},
				{"test1.com", 8080, "/test1-location", "($http_api_name != '')", "http://wrong_proxy", "", "HTTP",
					[]local.ProxiedAddress{{"wrong_proxy", 80, []*local.Socket{}, local.ToJSONError(nil)}},
					"", []string{},
					metav1.ConfigContextPos{"C:\\config_test\\conf.d\\location.conf", []int{0, 0, 0}}},
				{"test1.com", 8080, "/test-to-baidu", "", "https://www.baidu.com", "", "HTTPS",
					[]local.ProxiedAddress{{"www.baidu.com", 443, []*local.Socket{{"183.240.99.58", 443, bifrostv1.NetUnknown, bifrostv1.NetUnknown},
						{"183.240.99.169", 443, bifrostv1.NetUnknown, bifrostv1.NetUnknown}}, local.ToJSONError(nil)}},
					"", []string{},
					metav1.ConfigContextPos{"C:\\config_test\\conf.d\\location.conf", []int{1, 0}}},
				{"test2.com", 8080, "/test1-location", "", "http://right_proxy", "", "HTTP",
					[]local.ProxiedAddress{{"right_proxy", 80, []*local.Socket{}, local.ToJSONError(nil)}},
					"test inline comments", []string{"test inline comments"},
					metav1.ConfigContextPos{"C:\\config_test\\conf.d\\location.conf", []int{0, 1}}},
				{"test2.com", 8080, "/test1-location", "($http_api_name != '')", "http://wrong_proxy", "", "HTTP",
					[]local.ProxiedAddress{{"wrong_proxy", 80, []*local.Socket{}, local.ToJSONError(nil)}},
					"", []string{},
					metav1.ConfigContextPos{"C:\\config_test\\conf.d\\location.conf", []int{0, 0, 0}}},
				{"test2.com", 8080, "/test-to-baidu", "", "https://www.baidu.com", "", "HTTPS",
					[]local.ProxiedAddress{{"www.baidu.com", 443, []*local.Socket{{"183.240.99.58", 443, bifrostv1.NetUnknown, bifrostv1.NetUnknown},
						{"183.240.99.169", 443, bifrostv1.NetUnknown, bifrostv1.NetUnknown}}, local.ToJSONError(nil)}},
					"", []string{},
					metav1.ConfigContextPos{"C:\\config_test\\conf.d\\location.conf", []int{1, 0}}},
				{"test2.com", 8080, "/test2", "", "http://test2.test.com", "", "HTTP",
					[]local.ProxiedAddress{{"test2.test.com", 80, []*local.Socket{{"69.167.164.199", 80, bifrostv1.NetUnknown, bifrostv1.NetUnknown}},
						local.ToJSONError(nil)}},
					"", []string{},
					metav1.ConfigContextPos{"C:\\config_test\\conf.d\\location1.conf", []int{0, 0}}},
				{"", 8888, "", "", "baidu.com:22", "", "TCP/UDP",
					[]local.ProxiedAddress{{"baidu.com", 22, []*local.Socket{{"39.156.70.37", 22, 0, 0},
						{"220.181.7.203", 22, 0, 0}}, local.ToJSONError(nil)}},
					"", []string{},
					metav1.ConfigContextPos{"C:\\config_test\\nginx.conf", []int{12, 0, 1}}},
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
			gotJson, err := json.Marshal(got)
			if err != nil {
				t.Errorf("json.Marshal() error = %v", err)
			}
			wantJson, err := json.Marshal(tt.want)
			if err != nil {
				t.Errorf("json.Marshal() error = %v", err)
			}
			if string(gotJson) != string(wantJson) {
				t.Errorf("GetProxyServiceInfo() got = %s, want %s", gotJson, wantJson)
			}
		})
	}
}

func Test_webServerStatisticsStore_ConnectivityCheckOfProxyService(t *testing.T) {
	webSrvOpts := metav1.WebServerOptions{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBifrostsManager := bifrosts.NewMockManager(ctrl)
	mockBifrostsManager.EXPECT().GetBifrostClient(webSrvOpts).AnyTimes().Return(&bifrostclinetv1.Client{Factory: new(fake.ServiceClient)}, nil)
	httpPos := metav1.ConfigContextPos{Config: "C:\\config_test\\conf.d\\location.conf", ContextPosPath: []int{1, 0}}
	streamPos := metav1.ConfigContextPos{Config: "C:\\config_test\\nginx.conf", ContextPosPath: []int{12, 0, 1}}
	wrongPos := metav1.ConfigContextPos{Config: "C:\\config_test\\nginx.conf", ContextPosPath: []int{11, 0, 0}}
	type fields struct {
		bm bifrosts.Manager
	}
	type args struct {
		ctx          context.Context
		opts         metav1.WebServerOptions
		proxyPassPos metav1.ConfigContextPos
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
				ProxyAddress: []local.ProxiedAddress{{DomainName: "www.baidu.com", Port: 443, Sockets: []*local.Socket{{"183.240.99.58", 443, 1, 0},
					{"183.240.99.169", 443, 1, 0}}, ResolveErr: local.ToJSONError(nil)}},
				ContextPos: metav1.ConfigContextPos{Config: "C:\\config_test\\conf.d\\location.conf", ContextPosPath: []int{1, 0}},
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
					Sockets: []*local.Socket{{"39.156.70.37", 22, 2, 0},
						{"220.181.7.203", 22, 2, 0}}, ResolveErr: local.ToJSONError(nil)}},
				ContextPos: metav1.ConfigContextPos{Config: "C:\\config_test\\nginx.conf", ContextPosPath: []int{12, 0, 1}},
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
			if !reflect.DeepEqual(gotInfo, tt.wantInfo) {
				t.Errorf("ConnectivityCheckOfProxyService() gotInfo = %v, want %v", gotInfo, tt.wantInfo)
			}
		})
	}
}
