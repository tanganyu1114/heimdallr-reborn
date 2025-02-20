package bifrosts

import (
	"context"
	"encoding/json"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	"gin-vue-admin/internal/pkg/bifrosts"
	"gin-vue-admin/internal/pkg/bifrosts/fake"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	bifrostclinetv1 "github.com/ClessLi/bifrost/pkg/client/bifrost/v1"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
	nginx_context "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	"go.uber.org/mock/gomock"
	"reflect"
	"testing"
)

func Test_getStreamSrvsProxyBrief(t *testing.T) {
	webSrvOpts := metav1.WebServerOptions{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBifrostsManager := bifrosts.NewMockManager(ctrl)
	mockBifrostsManager.EXPECT().GetBifrostClient(webSrvOpts).AnyTimes().Return(&bifrostclinetv1.Client{Factory: new(fake.ServiceClient)}, nil)
	client, err := mockBifrostsManager.GetBifrostClient(webSrvOpts)
	if err != nil {
		t.Fatal(err)
	}
	data, err := client.WebServerConfig().Get(webSrvOpts.ServerName)
	if err != nil {
		t.Fatal(err)
	}
	conf, err := configuration.NewNginxConfigFromJsonBytes(data)
	if err != nil {
		t.Fatal(err)
	}
	proxyB, err := getStreamSrvsProxyBrief(conf)
	if err != nil {
		t.Fatal(err)
	}
	jsonData, err := json.Marshal(proxyB)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(jsonData))
}

func Test_getHttpSrvsProxyBrief(t *testing.T) {
	webSrvOpts := metav1.WebServerOptions{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBifrostsManager := bifrosts.NewMockManager(ctrl)
	mockBifrostsManager.EXPECT().GetBifrostClient(webSrvOpts).AnyTimes().Return(&bifrostclinetv1.Client{Factory: new(fake.ServiceClient)}, nil)
	client, err := mockBifrostsManager.GetBifrostClient(webSrvOpts)
	if err != nil {
		t.Fatal(err)
	}
	data, err := client.WebServerConfig().Get(webSrvOpts.ServerName)
	if err != nil {
		t.Fatal(err)
	}
	conf, err := configuration.NewNginxConfigFromJsonBytes(data)
	if err != nil {
		t.Fatal(err)
	}
	proxyB, err := getHttpSrvsProxyBrief(conf)
	if err != nil {
		t.Fatal(err)
	}
	jsonData, err := json.Marshal(proxyB)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(jsonData))
}

func Test_parseAddresses(t *testing.T) {
	/*  HTTP - Upstream:
	Syntax:	server address [parameters];
	Default:	—
	Context:	upstream
	Defines the address and other parameters of a server. The address can be specified as a domain name or IP address,
	with an optional port, or as a UNIX-domain socket path specified after the “unix:” prefix. If a port is not specified,
	the port 80 is used. A domain name that resolves to several IP addresses defines multiple servers at once.
	*/
	upstream := local.NewContext(context_type.TypeUpstream, "proxy_upstream"). // proxy upstream
											Insert(local.NewContext(context_type.TypeDirective, "server 127.0.0.1"), 0).       // proxy to 127.0.0.1:80
											Insert(local.NewContext(context_type.TypeDirective, "server example.com"), 1).     // proxy to example.com:80
											Insert(local.NewContext(context_type.TypeDirective, "server 127.0.0.1:8080"), 2).  // proxy to 127.0.0.1:8080
											Insert(local.NewContext(context_type.TypeDirective, "server example.com:8080"), 3) // proxy to example.com:8080
	http := local.NewContext(context_type.TypeHttp, "").Insert(upstream, 0)
	stream := local.NewContext(context_type.TypeStream, "").Insert(upstream, 0)
	type args struct {
		httpOrStream nginx_context.Context
		url          string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "proxy to IP and specified port",
			args: args{
				httpOrStream: http,
				url:          "http://127.0.0.1:8080/xxxx",
			},
			want: []string{"127.0.0.1:8080"},
		},
		{
			name: "http proxy to IP without port",
			args: args{
				httpOrStream: http,
				url:          "http://127.0.0.1/xxx/xxx",
			},
			want: []string{"127.0.0.1:80"},
		},
		{
			name: "https proxy to IP without port",
			args: args{
				httpOrStream: http,
				url:          "https://127.0.0.1/ssss",
			},
			want: []string{"127.0.0.1:443"},
		},
		{
			name: "proxy to domain and specified port",
			args: args{
				httpOrStream: http,
				url:          "https://example.com:8080/asldkfj",
			},
			want: []string{"example.com:8080"},
		},
		{
			name: "http proxy to domain whitout port",
			args: args{
				httpOrStream: http,
				url:          "http://example.com/asldkfj",
			},
			want: []string{"example.com:80"},
		},
		{
			name: "https proxy to domain without port",
			args: args{
				httpOrStream: http,
				url:          "https://example.com/asldkfj",
			},
			want: []string{"example.com:443"},
		},
		{
			name: "http proxy to upstream without specified port",
			args: args{
				httpOrStream: http,
				url:          "http://" + upstream.Value() + "/xxxx",
			},
			want: []string{
				"127.0.0.1:80",
				"example.com:80",
				"127.0.0.1:8080",
				"example.com:8080",
			},
		},
		{
			name: "https proxy to upstream without specified port",
			args: args{
				httpOrStream: http,
				url:          "https://" + upstream.Value() + "/xxxx",
			},
			want: []string{
				"127.0.0.1:80",
				"example.com:80",
				"127.0.0.1:8080",
				"example.com:8080",
			},
		},
		{
			name: "http proxy to upstream with specified port",
			args: args{
				httpOrStream: http,
				url:          "http://" + upstream.Value() + ":9090/xxxx",
			},
			want: []string{
				"127.0.0.1:9090",
				"example.com:9090",
			},
		},
		{
			name: "https proxy to upstream with specified port",
			args: args{
				httpOrStream: http,
				url:          "https://" + upstream.Value() + ":9090/xxxx",
			},
			want: []string{
				"127.0.0.1:9090",
				"example.com:9090",
			},
		},
		{
			name: "stream proxy to IP and specified port, with protocol",
			args: args{
				httpOrStream: stream,
				url:          "sftp://127.0.0.1:2022",
			},
			want: []string{"127.0.0.1:2022"},
		},
		{
			name: "stream proxy to IP and specified port, without protocol",
			args: args{
				httpOrStream: stream,
				url:          "127.0.0.1:8088",
			},
			want: []string{"127.0.0.1:8088"},
		},
		{
			name: "stream proxy to IP, with protocol",
			args: args{
				httpOrStream: stream,
				url:          "sftp://127.0.0.1",
			},
			want: []string{"127.0.0.1:unknown_port"},
		},
		{
			name: "stream proxy to IP, without protocol",
			args: args{
				httpOrStream: stream,
				url:          "127.0.0.1",
			},
			want: []string{"127.0.0.1:unknown_port"},
		},
		{
			name: "stream proxy to domain and specified port, with protocol",
			args: args{
				httpOrStream: stream,
				url:          "sftp://example.com:2022",
			},
			want: []string{"example.com:2022"},
		},
		{
			name: "stream proxy to domain and specified port, without protocol",
			args: args{
				httpOrStream: stream,
				url:          "example.com:8088",
			},
			want: []string{"example.com:8088"},
		},
		{
			name: "stream proxy to domain, with protocol",
			args: args{
				httpOrStream: stream,
				url:          "sftp://example.com",
			},
			want: []string{"example.com:unknown_port"},
		},
		{
			name: "stream proxy to domain, without protocol",
			args: args{
				httpOrStream: stream,
				url:          "example.com",
			},
			want: []string{"example.com:unknown_port"},
		},
		{
			name: "stream proxy to upstream and specified port, with protocol",
			args: args{
				httpOrStream: stream,
				url:          "sftp://" + upstream.Value() + ":2022",
			},
			want: []string{
				"127.0.0.1:2022",
				"example.com:2022",
			},
		},
		{
			name: "stream proxy to upstream and specified port, without protocol",
			args: args{
				httpOrStream: stream,
				url:          upstream.Value() + ":8088",
			},
			want: []string{
				"127.0.0.1:8088",
				"example.com:8088",
			},
		},
		{
			name: "stream proxy to upstream, with protocol",
			args: args{
				httpOrStream: stream,
				url:          "sftp_specified://" + upstream.Value(),
			},
			want: []string{
				"127.0.0.1:unknown_port",
				"example.com:unknown_port",
				"127.0.0.1:8080",
				"example.com:8080",
			},
		},
		{
			name: "stream proxy to upstream, without protocol",
			args: args{
				httpOrStream: stream,
				url:          upstream.Value(),
			},
			want: []string{
				"127.0.0.1:unknown_port",
				"example.com:unknown_port",
				"127.0.0.1:8080",
				"example.com:8080",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseAddresses(tt.args.httpOrStream, tt.args.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseAddresses() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
				{"HTTP", "localhost", 80, "/test_proxy", "10.1.1.1:443"},
				{"HTTP", "localhost", 80, "/test_proxy", "10.1.1.2:443"},
				{"HTTP", "localhost", -1, "/test_proxy", "10.1.1.1:443"},
				{"HTTP", "localhost", -1, "/test_proxy", "10.1.1.2:443"},
				{"HTTP", "test1.com", 80, "/test_proxy", "10.1.1.1:443"},
				{"HTTP", "test1.com", 80, "/test_proxy", "10.1.1.2:443"},
				{"HTTP", "test2.com", 80, "/test_proxy", "10.1.1.1:443"},
				{"HTTP", "test2.com", 80, "/test_proxy", "10.1.1.2:443"},
				{"HTTP", "test1.com", 80, "/test_proxy", "10.1.1.1:443"},
				{"HTTP", "test1.com", 80, "/test_proxy", "10.1.1.2:443"},
				{"HTTP", "test1.com", 80, "/test_proxy2", "baidu.com:443"},
				{"HTTP", "test1.com", 80, "/test_proxy3", "baidu2.com:80"},
				{"HTTP", "test1.com", 80, "/test_proxy4", "10.1.1.2:333"},
				{"HTTP", "test1.com", 8080, "/test_proxy", "10.1.1.1:443"},
				{"HTTP", "test1.com", 8080, "/test_proxy", "10.1.1.2:443"},
				{"HTTP", "test1.com", 8080, "/test_proxy2", "baidu.com:443"},
				{"HTTP", "test1.com", 8080, "/test_proxy3", "baidu2.com:80"},
				{"HTTP", "test1.com", 8080, "/test_proxy4", "10.1.1.2:333"},
				{"HTTP", "test2.com", 8080, "/test_proxy", "10.1.1.1:443"},
				{"HTTP", "test2.com", 8080, "/test_proxy", "10.1.1.2:443"},
				{"HTTP", "test2.com", 8080, "/test_proxy2", "baidu.com:443"},
				{"HTTP", "test2.com", 8080, "/test_proxy3", "baidu2.com:80"},
				{"HTTP", "test2.com", 8080, "/test_proxy4", "10.1.1.2:333"},
				{"TCP/UDP", "", 5510, "", "www.baidu.com:443"},
				{"TCP/UDP", "", 5520, "", "test.1.cn:22"},
				{"TCP/UDP", "", 5530, "", "vvvvv1:55"},
				{"TCP/UDP", "", 5530, "", "test.vv2.3.cn:44"},
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetProxyServiceInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}
