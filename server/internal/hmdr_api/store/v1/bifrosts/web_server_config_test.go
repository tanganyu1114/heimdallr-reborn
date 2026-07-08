package bifrosts

import (
	"context"
	"encoding/json"
	v1 "github.com/tanganyu1114/heimdallr-reborn/api/heimdallr_api/v1"
	"github.com/tanganyu1114/heimdallr-reborn/global"
	"github.com/tanganyu1114/heimdallr-reborn/internal/pkg/bifrosts"
	"github.com/tanganyu1114/heimdallr-reborn/internal/pkg/bifrosts/fake"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/internal/pkg/meta/v1"
	"github.com/tanganyu1114/heimdallr-reborn/pkg/sort_map"
	"reflect"
	"testing"

	bifrostclinetv1 "github.com/ClessLi/bifrost/pkg/client/bifrost/v1"
	utilsV3 "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/utils"
	logV1 "github.com/ClessLi/component-base/pkg/log/v1"
	"github.com/marmotedu/errors"
	"go.uber.org/mock/gomock"
)

func Test_newWebServerConfigStore(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBifrostsManager := bifrosts.NewMockManager(ctrl)
	store := &bifrostsStore{bm: mockBifrostsManager}
	type args struct {
		store *bifrostsStore
	}
	tests := []struct {
		name string
		args args
		want *webServerConfigStore
	}{
		{
			name: "normal_test",
			args: args{store: store},
			want: &webServerConfigStore{bm: mockBifrostsManager},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newWebServerConfigStore(tt.args.store); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newWebServerConfigStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webServerConfigStore_ChangeContextEnabledState(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	webSrvOpts := metav1.WebServerOptions{
		GroupID:    1,
		HostID:     1,
		ServerName: "test-bifrost",
	}

	// Use fake data to get config and fingerprints
	fakeSvc := fake.WebServerConfigService{}
	_, fakeFP, _ := fakeSvc.Get("test-bifrost")
	ofp := fakeFP.Fingerprints()

	// Create diff fingerprint for error test case
	difffp := utilsV3.ConfigFingerprints{
		"C:\\config_test\\nginx.conf": "1111111111111111111111111111111111111111111111111111111111111111",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBifrostsManager := bifrosts.NewMockManager(ctrl)
	mockBifrostsManager.EXPECT().GetBifrostClient(webSrvOpts).AnyTimes().Return(&bifrostclinetv1.Client{Factory: new(fake.ServiceClient)}, nil)
	type fields struct {
		bm bifrosts.Manager
	}
	type args struct {
		in0     context.Context
		opts    metav1.WebServerOptions
		ofp     utilsV3.ConfigFingerprints
		ctxmeta metav1.TargetConfigContextOptions[metav1.ConfigContextEnabledStateMeta]
	}
	tests := []struct {
		name                    string
		fields                  fields
		args                    args
		wantErr                 bool
		wantErrIsInconsistentFP bool
	}{
		{
			name:   "enable context",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				opts: webSrvOpts,
				ofp:  ofp,
				ctxmeta: metav1.TargetConfigContextOptions[metav1.ConfigContextEnabledStateMeta]{
					Position: metav1.ConfigContextPos{
						Config:         "C:\\config_test\\nginx.conf",
						ContextPosPath: []int{0},
					},
					TargetContext: metav1.ConfigContextEnabledStateMeta{Enabled: true},
				},
			},
			wantErr: false,
		},
		{
			name:   "disable context",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				opts: webSrvOpts,
				ofp:  ofp,
				ctxmeta: metav1.TargetConfigContextOptions[metav1.ConfigContextEnabledStateMeta]{
					Position: metav1.ConfigContextPos{
						Config:         "C:\\config_test\\nginx.conf",
						ContextPosPath: []int{1},
					},
					TargetContext: metav1.ConfigContextEnabledStateMeta{Enabled: false},
				},
			},
			wantErr: false,
		},
		{
			name:   "wrong position",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				opts: webSrvOpts,
				ofp:  ofp,
				ctxmeta: metav1.TargetConfigContextOptions[metav1.ConfigContextEnabledStateMeta]{
					Position: metav1.ConfigContextPos{
						Config:         "C:\\config_test\\nginx.conf",
						ContextPosPath: []int{1, 2, 3, 4, 5, 6, 7},
					},
					TargetContext: metav1.ConfigContextEnabledStateMeta{Enabled: true},
				},
			},
			wantErr: true,
		},
		{
			name:   "inconsistent fingerprints",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				opts: webSrvOpts,
				ofp:  difffp,
				ctxmeta: metav1.TargetConfigContextOptions[metav1.ConfigContextEnabledStateMeta]{
					Position: metav1.ConfigContextPos{
						Config:         "C:\\config_test\\nginx.conf",
						ContextPosPath: []int{0},
					},
					TargetContext: metav1.ConfigContextEnabledStateMeta{Enabled: true},
				},
			},
			wantErr:                 true,
			wantErrIsInconsistentFP: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigStore{
				bm: tt.fields.bm,
			}
			err := w.ChangeContextEnabledState(tt.args.in0, tt.args.opts, tt.args.ofp, tt.args.ctxmeta)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChangeContextEnabledState() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && (errors.Is(err, metav1.ErrInconsistentFingerprints) || errors.IsCode(err, 110010)) != tt.wantErrIsInconsistentFP {
				t.Errorf("ChangeContextEnabledState() error = %v, wantErrIsInconsistentFP %v", err, tt.wantErrIsInconsistentFP)
			}
		})
	}
}

func Test_webServerConfigStore_GetConfig(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	webSrvOpts := metav1.WebServerOptions{
		GroupID:    1,
		HostID:     1,
		ServerName: "test-bifrost",
	}

	// Use fake data to get config and fingerprints
	fakeSvc := fake.WebServerConfigService{}
	_, fakeFP, _ := fakeSvc.Get("test-bifrost")
	ofp := fakeFP.Fingerprints()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBifrostsManager := bifrosts.NewMockManager(ctrl)
	mockBifrostsManager.EXPECT().GetBifrostClient(webSrvOpts).AnyTimes().Return(&bifrostclinetv1.Client{Factory: new(fake.ServiceClient)}, nil)
	type fields struct {
		bm bifrosts.Manager
	}
	type args struct {
		ctx  context.Context
		opts metav1.WebServerOptions
	}
	tests := []struct {
		name                     string
		fields                   fields
		args                     args
		wantConfigTextLines      []string
		wantOriginalFingerprints utilsV3.ConfigFingerprints
		wantErr                  bool
	}{
		{
			name:   "normal test",
			fields: fields{bm: mockBifrostsManager},
			args:   args{opts: webSrvOpts},
			wantConfigTextLines: []string{
				"# testConfigManager",
				"# testConfigManager",
				"# testConfigManager_update...",
				"# testConfigManager_update...",
				"# user nobody;",
				"worker_processes 1;    # test inline comments",
				"# error_log logs/error.log;",
				"# error_log logs/error.log  notice;",
				"# error_log logs/error.log  info;",
				"# pid logs/nginx.pid;",
				"events {",
				"    worker_connections 1024;",
				"}",
				"stream {",
				"    server {",
				"        listen 8888;",
				"        proxy_pass baidu.com:22;",
				"    }",
				"}",
				"http {",
				"    # include <== mime.types",
				"    default_type application/octet-stream;",
				"    # log_format  main  '$remote_addr - $remote_user [$time_local] \"$request\" '",
				"    # '$status $body_bytes_sent \"$http_referer\" '",
				"    # '\"$http_user_agent\" \"$http_x_forwarded_for\"';",
				"    # access_log logs/access.log  main;",
				"    sendfile on;",
				"    # tcp_nopush on;",
				"    # keepalive_timeout 0;",
				"    keepalive_timeout 65;",
				"    # gzip on;",
				"    # include <== ./conf.d/test*.com.conf",
				"    server {",
				"        listen 80;",
				"        server_name test1.com;",
				"        # include <== ./conf.d/location.conf",
				"        location /test1-location {",
				"            if ($http_api_name != '') {",
				"                proxy_pass http://wrong_proxy;",
				"                break;",
				"            }",
				"            proxy_pass http://right_proxy;    # test inline comments",
				"        }",
				"        location /test-to-baidu {",
				"            proxy_pass https://www.baidu.com;",
				"        }",
				"    }",
				"    server {",
				"        listen 8080;",
				"        server_name test1.com;",
				"        # include <== ./conf.d/location*.conf",
				"        location /test1-location {",
				"            if ($http_api_name != '') {",
				"                proxy_pass http://wrong_proxy;",
				"                break;",
				"            }",
				"            proxy_pass http://right_proxy;    # test inline comments",
				"        }",
				"        location /test-to-baidu {",
				"            proxy_pass https://www.baidu.com;",
				"        }",
				"        location /test2 {",
				"            proxy_pass http://test2.test.com;",
				"        }",
				"        # # include <== ./conf.d/location.conf",
				"    }",
				"    server {",
				"        listen 8080;",
				"        server_name test2.com;",
				"        # include <== ./conf.d/location*.conf",
				"        location /test1-location {",
				"            if ($http_api_name != '') {",
				"                proxy_pass http://wrong_proxy;",
				"                break;",
				"            }",
				"            proxy_pass http://right_proxy;    # test inline comments",
				"        }",
				"        location /test-to-baidu {",
				"            proxy_pass https://www.baidu.com;",
				"        }",
				"        location /test2 {",
				"            proxy_pass http://test2.test.com;",
				"        }",
				"        # # include <== ./conf.d/location.conf",
				"    }",
				"    server {",
				"        listen 80;",
				"        server_name localhost;",
				"        # charset koi8-r;",
				"        # access_log logs/host.access.log  main;",
				"        location / {",
				"            root html;",
				"            index index.html index.htm;",
				"        }",
				"        # error_page 404              /404.html;",
				"        # redirect server error pages to the static page /50x.html",
				"        #",
				"        error_page 500 502 503 504  /50x.html;",
				"        location = /50x.html {",
				"            root html;",
				"        }",
				"        # proxy the PHP scripts to Apache listening on 127.0.0.1:80",
				"        #",
				"        # location ~ \\.php$ {",
				"        #     proxy_pass http://127.0.0.1;",
				"        # }",
				"        # pass the PHP scripts to FastCGI server listening on 127.0.0.1:9000",
				"        #",
				"        # location ~ \\.php$ {",
				"        #     root html;",
				"        #     fastcgi_pass 127.0.0.1:9000;",
				"        #     fastcgi_index index.php;",
				"        #     fastcgi_param SCRIPT_FILENAME  /scripts$fastcgi_script_name;",
				"        #     # include <== fastcgi_params",
				"        # }",
				"        # deny access to .htaccess files, if Apache's document root",
				"        # concurs with nginx's one",
				"        #",
				"        # location ~ /\\.ht {",
				"        #     deny all;",
				"        # }",
				"    }",
				"    # another virtual host using mix of IP-, name-, and port-based configuration",
				"    #",
				"    # server {",
				"    #     listen 8000;",
				"    #     listen somename:8080;",
				"    #     server_name somename  alias  another.alias;",
				"    #     location / {",
				"    #         root html;",
				"    #         index index.html index.htm;",
				"    #     }",
				"    # }",
				"    # HTTPS server",
				"    #",
				"    # server {",
				"    #     listen 443 ssl;",
				"    #     server_name localhost;",
				"    #     ssl_certificate cert.pem;",
				"    #     ssl_certificate_key cert.key;",
				"    #     ssl_session_cache shared:SSL:1m;",
				"    #     ssl_session_timeout 5m;",
				"    #     ssl_ciphers HIGH:!aNULL:!MD5;",
				"    #     ssl_prefer_server_ciphers on;",
				"    #     location / {",
				"    #         root html;",
				"    #         index index.html index.htm;",
				"    #     }",
				"    # }",
				"}",
			},
			wantOriginalFingerprints: ofp,
			wantErr:                  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigStore{
				bm: tt.fields.bm,
			}
			gotConfig, gotofp, err := w.GetConfig(tt.args.ctx, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotConfig.TextLines(), tt.wantConfigTextLines) {
				t.Errorf("GetConfig() got.TextLines() = %v, want %v", gotConfig.TextLines(), tt.wantConfigTextLines)
			}
			if !reflect.DeepEqual(gotofp.Fingerprints(), tt.wantOriginalFingerprints) {
				t.Errorf("GetConfig() got.OriginalFingerprints = %v, want %v", gotofp.Fingerprints(), tt.wantOriginalFingerprints)
			}
		})
	}
}

func Test_webServerConfigStore_GetContext(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	webSrvOpts := metav1.WebServerOptions{
		GroupID:    1,
		HostID:     1,
		ServerName: "test-bifrost",
	}

	// Use fake data to get config and fingerprints
	fakeSvc := fake.WebServerConfigService{}
	_, fakeFP, _ := fakeSvc.Get("test-bifrost")
	ofp := fakeFP.Fingerprints()

	// Create diff fingerprint for error test case
	difffp := utilsV3.ConfigFingerprints{
		"C:\\config_test\\nginx.conf": "1111111111111111111111111111111111111111111111111111111111111111",
	}

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
		fp   utilsV3.ConfigFingerprints
		pos  metav1.ConfigContextPos
	}
	tests := []struct {
		name                    string
		fields                  fields
		args                    args
		wantConfigLines         []string
		wantErr                 bool
		wantErrIsInconsistentFP bool
	}{
		{
			name:   "one level pos path",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				opts: webSrvOpts,
				fp:   ofp,
				pos: metav1.ConfigContextPos{
					Config:         "C:\\config_test\\conf.d\\location.conf",
					ContextPosPath: []int{0},
				},
			},
			wantConfigLines: []string{
				"location /test1-location {",
				"    if ($http_api_name != '') {",
				"        proxy_pass http://wrong_proxy;",
				"        break;",
				"    }",
				"    proxy_pass http://right_proxy;    # test inline comments",
				"}",
			},
			wantErr: false,
		},
		{
			name:   "two level pos path",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				opts: webSrvOpts,
				fp:   ofp,
				pos: metav1.ConfigContextPos{
					Config:         "C:\\config_test\\conf.d\\location.conf",
					ContextPosPath: []int{0, 1},
				},
			},
			wantConfigLines: []string{
				"proxy_pass http://right_proxy;",
			},
			wantErr: false,
		},
		{
			name:   "null pos path",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				opts: webSrvOpts,
				fp:   ofp,
				pos: metav1.ConfigContextPos{
					Config:         "C:\\config_test\\nginx.conf",
					ContextPosPath: []int{},
				},
			},
			wantConfigLines: []string{
				"# testConfigManager",
				"# testConfigManager",
				"# testConfigManager_update...",
				"# testConfigManager_update...",
				"# user nobody;",
				"worker_processes 1;    # test inline comments",
				"# error_log logs/error.log;",
				"# error_log logs/error.log  notice;",
				"# error_log logs/error.log  info;",
				"# pid logs/nginx.pid;",
				"events {",
				"    worker_connections 1024;",
				"}",
				"stream {",
				"    server {",
				"        listen 8888;",
				"        proxy_pass baidu.com:22;",
				"    }",
				"}",
				"http {",
				"    # include <== mime.types",
				"    default_type application/octet-stream;",
				"    # log_format  main  '$remote_addr - $remote_user [$time_local] \"$request\" '",
				"    # '$status $body_bytes_sent \"$http_referer\" '",
				"    # '\"$http_user_agent\" \"$http_x_forwarded_for\"';",
				"    # access_log logs/access.log  main;",
				"    sendfile on;",
				"    # tcp_nopush on;",
				"    # keepalive_timeout 0;",
				"    keepalive_timeout 65;",
				"    # gzip on;",
				"    # include <== ./conf.d/test*.com.conf",
				"    server {",
				"        listen 80;",
				"        server_name test1.com;",
				"        # include <== ./conf.d/location.conf",
				"        location /test1-location {",
				"            if ($http_api_name != '') {",
				"                proxy_pass http://wrong_proxy;",
				"                break;",
				"            }",
				"            proxy_pass http://right_proxy;    # test inline comments",
				"        }",
				"        location /test-to-baidu {",
				"            proxy_pass https://www.baidu.com;",
				"        }",
				"    }",
				"    server {",
				"        listen 8080;",
				"        server_name test1.com;",
				"        # include <== ./conf.d/location*.conf",
				"        location /test1-location {",
				"            if ($http_api_name != '') {",
				"                proxy_pass http://wrong_proxy;",
				"                break;",
				"            }",
				"            proxy_pass http://right_proxy;    # test inline comments",
				"        }",
				"        location /test-to-baidu {",
				"            proxy_pass https://www.baidu.com;",
				"        }",
				"        location /test2 {",
				"            proxy_pass http://test2.test.com;",
				"        }",
				"        # # include <== ./conf.d/location.conf",
				"    }",
				"    server {",
				"        listen 8080;",
				"        server_name test2.com;",
				"        # include <== ./conf.d/location*.conf",
				"        location /test1-location {",
				"            if ($http_api_name != '') {",
				"                proxy_pass http://wrong_proxy;",
				"                break;",
				"            }",
				"            proxy_pass http://right_proxy;    # test inline comments",
				"        }",
				"        location /test-to-baidu {",
				"            proxy_pass https://www.baidu.com;",
				"        }",
				"        location /test2 {",
				"            proxy_pass http://test2.test.com;",
				"        }",
				"        # # include <== ./conf.d/location.conf",
				"    }",
				"    server {",
				"        listen 80;",
				"        server_name localhost;",
				"        # charset koi8-r;",
				"        # access_log logs/host.access.log  main;",
				"        location / {",
				"            root html;",
				"            index index.html index.htm;",
				"        }",
				"        # error_page 404              /404.html;",
				"        # redirect server error pages to the static page /50x.html",
				"        #",
				"        error_page 500 502 503 504  /50x.html;",
				"        location = /50x.html {",
				"            root html;",
				"        }",
				"        # proxy the PHP scripts to Apache listening on 127.0.0.1:80",
				"        #",
				"        # location ~ \\.php$ {",
				"        #     proxy_pass http://127.0.0.1;",
				"        # }",
				"        # pass the PHP scripts to FastCGI server listening on 127.0.0.1:9000",
				"        #",
				"        # location ~ \\.php$ {",
				"        #     root html;",
				"        #     fastcgi_pass 127.0.0.1:9000;",
				"        #     fastcgi_index index.php;",
				"        #     fastcgi_param SCRIPT_FILENAME  /scripts$fastcgi_script_name;",
				"        #     # include <== fastcgi_params",
				"        # }",
				"        # deny access to .htaccess files, if Apache's document root",
				"        # concurs with nginx's one",
				"        #",
				"        # location ~ /\\.ht {",
				"        #     deny all;",
				"        # }",
				"    }",
				"    # another virtual host using mix of IP-, name-, and port-based configuration",
				"    #",
				"    # server {",
				"    #     listen 8000;",
				"    #     listen somename:8080;",
				"    #     server_name somename  alias  another.alias;",
				"    #     location / {",
				"    #         root html;",
				"    #         index index.html index.htm;",
				"    #     }",
				"    # }",
				"    # HTTPS server",
				"    #",
				"    # server {",
				"    #     listen 443 ssl;",
				"    #     server_name localhost;",
				"    #     ssl_certificate cert.pem;",
				"    #     ssl_certificate_key cert.key;",
				"    #     ssl_session_cache shared:SSL:1m;",
				"    #     ssl_session_timeout 5m;",
				"    #     ssl_ciphers HIGH:!aNULL:!MD5;",
				"    #     ssl_prefer_server_ciphers on;",
				"    #     location / {",
				"    #         root html;",
				"    #         index index.html index.htm;",
				"    #     }",
				"    # }",
				"}",
			},
		},
		{
			name:   "nil pos path",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				opts: webSrvOpts,
				fp:   ofp,
				pos: metav1.ConfigContextPos{
					Config:         "C:\\config_test\\nginx.conf",
					ContextPosPath: nil,
				},
			},
			wantConfigLines: []string{
				"# testConfigManager",
				"# testConfigManager",
				"# testConfigManager_update...",
				"# testConfigManager_update...",
				"# user nobody;",
				"worker_processes 1;    # test inline comments",
				"# error_log logs/error.log;",
				"# error_log logs/error.log  notice;",
				"# error_log logs/error.log  info;",
				"# pid logs/nginx.pid;",
				"events {",
				"    worker_connections 1024;",
				"}",
				"stream {",
				"    server {",
				"        listen 8888;",
				"        proxy_pass baidu.com:22;",
				"    }",
				"}",
				"http {",
				"    # include <== mime.types",
				"    default_type application/octet-stream;",
				"    # log_format  main  '$remote_addr - $remote_user [$time_local] \"$request\" '",
				"    # '$status $body_bytes_sent \"$http_referer\" '",
				"    # '\"$http_user_agent\" \"$http_x_forwarded_for\"';",
				"    # access_log logs/access.log  main;",
				"    sendfile on;",
				"    # tcp_nopush on;",
				"    # keepalive_timeout 0;",
				"    keepalive_timeout 65;",
				"    # gzip on;",
				"    # include <== ./conf.d/test*.com.conf",
				"    server {",
				"        listen 80;",
				"        server_name test1.com;",
				"        # include <== ./conf.d/location.conf",
				"        location /test1-location {",
				"            if ($http_api_name != '') {",
				"                proxy_pass http://wrong_proxy;",
				"                break;",
				"            }",
				"            proxy_pass http://right_proxy;    # test inline comments",
				"        }",
				"        location /test-to-baidu {",
				"            proxy_pass https://www.baidu.com;",
				"        }",
				"    }",
				"    server {",
				"        listen 8080;",
				"        server_name test1.com;",
				"        # include <== ./conf.d/location*.conf",
				"        location /test1-location {",
				"            if ($http_api_name != '') {",
				"                proxy_pass http://wrong_proxy;",
				"                break;",
				"            }",
				"            proxy_pass http://right_proxy;    # test inline comments",
				"        }",
				"        location /test-to-baidu {",
				"            proxy_pass https://www.baidu.com;",
				"        }",
				"        location /test2 {",
				"            proxy_pass http://test2.test.com;",
				"        }",
				"        # # include <== ./conf.d/location.conf",
				"    }",
				"    server {",
				"        listen 8080;",
				"        server_name test2.com;",
				"        # include <== ./conf.d/location*.conf",
				"        location /test1-location {",
				"            if ($http_api_name != '') {",
				"                proxy_pass http://wrong_proxy;",
				"                break;",
				"            }",
				"            proxy_pass http://right_proxy;    # test inline comments",
				"        }",
				"        location /test-to-baidu {",
				"            proxy_pass https://www.baidu.com;",
				"        }",
				"        location /test2 {",
				"            proxy_pass http://test2.test.com;",
				"        }",
				"        # # include <== ./conf.d/location.conf",
				"    }",
				"    server {",
				"        listen 80;",
				"        server_name localhost;",
				"        # charset koi8-r;",
				"        # access_log logs/host.access.log  main;",
				"        location / {",
				"            root html;",
				"            index index.html index.htm;",
				"        }",
				"        # error_page 404              /404.html;",
				"        # redirect server error pages to the static page /50x.html",
				"        #",
				"        error_page 500 502 503 504  /50x.html;",
				"        location = /50x.html {",
				"            root html;",
				"        }",
				"        # proxy the PHP scripts to Apache listening on 127.0.0.1:80",
				"        #",
				"        # location ~ \\.php$ {",
				"        #     proxy_pass http://127.0.0.1;",
				"        # }",
				"        # pass the PHP scripts to FastCGI server listening on 127.0.0.1:9000",
				"        #",
				"        # location ~ \\.php$ {",
				"        #     root html;",
				"        #     fastcgi_pass 127.0.0.1:9000;",
				"        #     fastcgi_index index.php;",
				"        #     fastcgi_param SCRIPT_FILENAME  /scripts$fastcgi_script_name;",
				"        #     # include <== fastcgi_params",
				"        # }",
				"        # deny access to .htaccess files, if Apache's document root",
				"        # concurs with nginx's one",
				"        #",
				"        # location ~ /\\.ht {",
				"        #     deny all;",
				"        # }",
				"    }",
				"    # another virtual host using mix of IP-, name-, and port-based configuration",
				"    #",
				"    # server {",
				"    #     listen 8000;",
				"    #     listen somename:8080;",
				"    #     server_name somename  alias  another.alias;",
				"    #     location / {",
				"    #         root html;",
				"    #         index index.html index.htm;",
				"    #     }",
				"    # }",
				"    # HTTPS server",
				"    #",
				"    # server {",
				"    #     listen 443 ssl;",
				"    #     server_name localhost;",
				"    #     ssl_certificate cert.pem;",
				"    #     ssl_certificate_key cert.key;",
				"    #     ssl_session_cache shared:SSL:1m;",
				"    #     ssl_session_timeout 5m;",
				"    #     ssl_ciphers HIGH:!aNULL:!MD5;",
				"    #     ssl_prefer_server_ciphers on;",
				"    #     location / {",
				"    #         root html;",
				"    #         index index.html index.htm;",
				"    #     }",
				"    # }",
				"}",
			},
		},
		{
			name:   "wrong pos path",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				opts: webSrvOpts,
				fp:   ofp,
				pos: metav1.ConfigContextPos{
					Config:         "C:\\config_test\\nginx.conf",
					ContextPosPath: []int{1, 2, 3, 4},
				},
			},
			wantErr: true,
		},
		{
			name:   "inconsistent fingerprints",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				opts: webSrvOpts,
				fp:   difffp,
				pos: metav1.ConfigContextPos{
					Config:         "C:\\config_test\\conf.d\\location.conf",
					ContextPosPath: []int{0, 1},
				},
			},
			wantErr:                 true,
			wantErrIsInconsistentFP: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigStore{
				bm: tt.fields.bm,
			}
			got, err := w.GetContext(tt.args.in0, tt.args.opts, tt.args.fp, tt.args.pos)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err != nil && (errors.Is(err, metav1.ErrInconsistentFingerprints) || errors.IsCode(err, 110010)) != tt.wantErrIsInconsistentFP {
				t.Errorf("GetContext() error = %v, wantErrIsInconsistentFP %v", err, tt.wantErrIsInconsistentFP)
				return
			}
			gotConfigLines, _ := got.ConfigLines(false)
			if !reflect.DeepEqual(gotConfigLines, tt.wantConfigLines) {
				t.Errorf("GetContext() got.ConfigLines( false ) = %v, want %v", gotConfigLines, tt.wantConfigLines)
			}
		})
	}
}

func Test_webServerConfigStore_GetIncludedConfigs(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	webSrvOpts := metav1.WebServerOptions{
		GroupID:    1,
		HostID:     1,
		ServerName: "test-bifrost",
	}

	// Use fake data to get config and fingerprints
	fakeSvc := fake.WebServerConfigService{}
	_, fakeFP, _ := fakeSvc.Get("test-bifrost")
	ofp := fakeFP.Fingerprints()

	// Create diff fingerprint for error test case
	difffp := utilsV3.ConfigFingerprints{
		"C:\\config_test\\nginx.conf": "1111111111111111111111111111111111111111111111111111111111111111",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBifrostsManager := bifrosts.NewMockManager(ctrl)
	mockBifrostsManager.EXPECT().GetBifrostClient(webSrvOpts).AnyTimes().Return(&bifrostclinetv1.Client{Factory: new(fake.ServiceClient)}, nil)
	type fields struct {
		bm bifrosts.Manager
	}
	type args struct {
		ctx  context.Context
		opts metav1.WebServerOptions
		fp   utilsV3.ConfigFingerprints
		pos  metav1.ConfigContextPos
	}
	tests := []struct {
		name                    string
		fields                  fields
		args                    args
		want                    []string
		wantErr                 bool
		wantErrIsInconsistentFP bool
	}{
		{
			name:   "normal test",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				opts: webSrvOpts,
				fp:   ofp,
				pos: metav1.ConfigContextPos{
					Config:         "C:\\config_test\\nginx.conf",
					ContextPosPath: []int{13, 11},
				},
			},
			want: []string{
				"C:\\config_test\\conf.d\\test1.com.conf",
				"C:\\config_test\\conf.d\\test2.com.conf",
			},
			wantErr: false,
		},
		{
			name:   "wrong position",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				opts: webSrvOpts,
				fp:   ofp,
				pos: metav1.ConfigContextPos{
					Config:         "C:\\config_test\\conf.d\\server_test1.con",
					ContextPosPath: nil,
				},
			},
			wantErr: true,
		},
		{
			name:   "the target is not an `include` context",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				opts: webSrvOpts,
				fp:   ofp,
				pos: metav1.ConfigContextPos{
					Config:         "C:\\config_test\\nginx.conf",
					ContextPosPath: []int{8, 1},
				},
			},
			wantErr: true,
		},
		{
			name:   "inconsistent fingerprints",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				opts: webSrvOpts,
				fp:   difffp,
				pos: metav1.ConfigContextPos{
					Config:         "C:\\config_test\\nginx.conf",
					ContextPosPath: []int{8, 2},
				},
			},
			wantErr:                 true,
			wantErrIsInconsistentFP: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigStore{
				bm: tt.fields.bm,
			}
			got, err := w.GetIncludedConfigs(tt.args.ctx, tt.args.opts, tt.args.fp, tt.args.pos)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIncludedConfigs() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err != nil && (errors.Is(err, metav1.ErrInconsistentFingerprints) || errors.IsCode(err, 110010)) != tt.wantErrIsInconsistentFP {
				t.Errorf("GetIncludedConfigs() error = %v, wantErrIsInconsistentFP %v", err, tt.wantErrIsInconsistentFP)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetIncludedConfigs() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webServerConfigStore_GetOptions(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBifrostsManager := bifrosts.NewMockManager(ctrl)

	// Test case 1: empty result
	mockBifrostsManager.EXPECT().Range(gomock.Any()).DoAndReturn(func(f func(keyer sort_map.Keyer, v interface{}) bool) {
		// No groups, return immediately
	}).AnyTimes()

	type fields struct {
		bm bifrosts.Manager
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []v1.BifrostGroupMeta
		wantErr bool
	}{
		{
			name:    "empty_result",
			fields:  fields{bm: mockBifrostsManager},
			args:    args{ctx: context.Background()},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigStore{
				bm: tt.fields.bm,
			}
			got, err := w.GetOptions(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				gotJson, _ := json.Marshal(got)
				wantJson, _ := json.Marshal(tt.want)
				t.Errorf("GetOptions() got = %s, want %s", gotJson, wantJson)
			}
		})
	}
}

func Test_webServerConfigStore_InsertWithClone(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	webSrvOpts := metav1.WebServerOptions{
		GroupID:    1,
		HostID:     1,
		ServerName: "test-bifrost",
	}

	// Use fake data to get config and fingerprints
	fakeSvc := fake.WebServerConfigService{}
	_, fakeFP, _ := fakeSvc.Get("test-bifrost")
	ofp := fakeFP.Fingerprints()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBifrostsManager := bifrosts.NewMockManager(ctrl)
	mockBifrostsManager.EXPECT().GetBifrostClient(webSrvOpts).AnyTimes().Return(&bifrostclinetv1.Client{Factory: new(fake.ServiceClient)}, nil)
	type fields struct {
		bm bifrosts.Manager
	}
	type args struct {
		in0            context.Context
		opts           metav1.WebServerOptions
		ofp            utilsV3.ConfigFingerprints
		ctxmeta        metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]
		disabledTarget bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "normal test",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				opts: webSrvOpts,
				ofp:  ofp,
				ctxmeta: metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]{
					Position: metav1.ConfigContextPos{
						Config:         "C:\\config_test\\nginx.conf",
						ContextPosPath: []int{13, 15, 3},
					},
					TargetContext: metav1.CloneConfigContextMeta{ConfigContextPos: metav1.ConfigContextPos{
						Config:         "C:\\config_test\\conf.d\\location1.conf",
						ContextPosPath: []int{0},
					}},
				},
				disabledTarget: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigStore{
				bm: tt.fields.bm,
			}
			if err := w.InsertWithClone(tt.args.in0, tt.args.opts, tt.args.ofp, tt.args.ctxmeta, tt.args.disabledTarget); (err != nil) != tt.wantErr {
				t.Errorf("InsertWithClone() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigStore_InsertWithNew(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	webSrvOpts := metav1.WebServerOptions{
		GroupID:    1,
		HostID:     1,
		ServerName: "test-bifrost",
	}

	// Use fake data to get config and fingerprints
	fakeSvc := fake.WebServerConfigService{}
	_, fakeFP, _ := fakeSvc.Get("test-bifrost")
	ofp := fakeFP.Fingerprints()

	// Create diff fingerprint for error test case
	difffp := utilsV3.ConfigFingerprints{
		"C:\\config_test\\nginx.conf": "1111111111111111111111111111111111111111111111111111111111111111",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBifrostsManager := bifrosts.NewMockManager(ctrl)
	mockBifrostsManager.EXPECT().GetBifrostClient(webSrvOpts).AnyTimes().Return(&bifrostclinetv1.Client{Factory: new(fake.ServiceClient)}, nil)
	type fields struct {
		bm bifrosts.Manager
	}
	type args struct {
		in0            context.Context
		opts           metav1.WebServerOptions
		ofp            utilsV3.ConfigFingerprints
		ctxmeta        metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]
		disabledTarget bool
	}
	tests := []struct {
		name                    string
		fields                  fields
		args                    args
		wantErr                 bool
		wantErrIsInconsistentFP bool
	}{
		{
			name:   "normal test",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				opts: webSrvOpts,
				ofp:  ofp,
				ctxmeta: metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]{
					Position: metav1.ConfigContextPos{
						Config:         "C:\\config_test\\nginx.conf",
						ContextPosPath: []int{13, 15, 3},
					},
					TargetContext: metav1.NewConfigContextMeta{
						ContextType:  "location",
						ContextValue: "~ /normal-test",
					},
				},
				disabledTarget: true,
			},
			wantErr: false,
		},
		{
			name:   "inconsistent fingerprints",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				opts: webSrvOpts,
				ofp:  difffp,
				ctxmeta: metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]{
					Position: metav1.ConfigContextPos{
						Config:         "C:\\config_test\\conf.d\\location2.conf",
						ContextPosPath: []int{2},
					},
					TargetContext: metav1.NewConfigContextMeta{
						ContextType:  "location",
						ContextValue: "~ /normal-test",
					},
				},
				disabledTarget: true,
			},
			wantErr:                 true,
			wantErrIsInconsistentFP: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigStore{
				bm: tt.fields.bm,
			}
			if err := w.InsertWithNew(tt.args.in0, tt.args.opts, tt.args.ofp, tt.args.ctxmeta, tt.args.disabledTarget); (err != nil) != tt.wantErr {
				t.Errorf("InsertWithNew() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && (errors.Is(err, metav1.ErrInconsistentFingerprints) || errors.IsCode(err, 110010)) != tt.wantErrIsInconsistentFP {
				t.Errorf("InsertWithNew() error = %v, wantErrIsInconsistentFP %v", err, tt.wantErrIsInconsistentFP)
			}
		})
	}
}

func Test_webServerConfigStore_ModifyContextValue(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	webSrvOpts := metav1.WebServerOptions{
		GroupID:    1,
		HostID:     1,
		ServerName: "test-bifrost",
	}

	// Use fake data to get config and fingerprints
	fakeSvc := fake.WebServerConfigService{}
	_, fakeFP, _ := fakeSvc.Get("test-bifrost")
	ofp := fakeFP.Fingerprints()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBifrostsManager := bifrosts.NewMockManager(ctrl)
	mockBifrostsManager.EXPECT().GetBifrostClient(webSrvOpts).AnyTimes().Return(&bifrostclinetv1.Client{Factory: new(fake.ServiceClient)}, nil)
	type fields struct {
		bm bifrosts.Manager
	}
	type args struct {
		ctx     context.Context
		opts    metav1.WebServerOptions
		ofp     utilsV3.ConfigFingerprints
		ctxmeta metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "normal test",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				opts: webSrvOpts,
				ofp:  ofp,
				ctxmeta: metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]{
					Position: metav1.ConfigContextPos{
						Config:         "C:\\config_test\\nginx.conf",
						ContextPosPath: []int{13, 15, 3},
					},
					TargetContext: metav1.NewConfigContextMeta{
						ContextType:  "location",
						ContextValue: "~ /normal-test",
					},
				},
			},
			wantErr: false,
		},
		{
			name:   "unmatched context type",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				opts: webSrvOpts,
				ofp:  ofp,
				ctxmeta: metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]{
					Position: metav1.ConfigContextPos{
						Config:         "C:\\config_test\\conf.d\\location2.conf",
						ContextPosPath: []int{2},
					},
					TargetContext: metav1.NewConfigContextMeta{
						ContextType:  "location",
						ContextValue: "~ /normal-test",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigStore{
				bm: tt.fields.bm,
			}
			if err := w.ModifyContextValue(tt.args.ctx, tt.args.opts, tt.args.ofp, tt.args.ctxmeta); (err != nil) != tt.wantErr {
				t.Errorf("ModifyContextValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigStore_ModifyWithClone(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	webSrvOpts := metav1.WebServerOptions{
		GroupID:    1,
		HostID:     1,
		ServerName: "test-bifrost",
	}

	// Use fake data to get config and fingerprints
	fakeSvc := fake.WebServerConfigService{}
	_, fakeFP, _ := fakeSvc.Get("test-bifrost")
	ofp := fakeFP.Fingerprints()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBifrostsManager := bifrosts.NewMockManager(ctrl)
	mockBifrostsManager.EXPECT().GetBifrostClient(webSrvOpts).AnyTimes().Return(&bifrostclinetv1.Client{Factory: new(fake.ServiceClient)}, nil)
	type fields struct {
		bm bifrosts.Manager
	}
	type args struct {
		in0     context.Context
		opts    metav1.WebServerOptions
		ofp     utilsV3.ConfigFingerprints
		ctxmeta metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "normal test",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				opts: webSrvOpts,
				ofp:  ofp,
				ctxmeta: metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]{
					Position: metav1.ConfigContextPos{
						Config:         "C:\\config_test\\nginx.conf",
						ContextPosPath: []int{13, 15, 3},
					},
					TargetContext: metav1.CloneConfigContextMeta{ConfigContextPos: metav1.ConfigContextPos{
						Config:         "C:\\config_test\\conf.d\\location1.conf",
						ContextPosPath: []int{0},
					}},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigStore{
				bm: tt.fields.bm,
			}
			if err := w.ModifyWithClone(tt.args.in0, tt.args.opts, tt.args.ofp, tt.args.ctxmeta); (err != nil) != tt.wantErr {
				t.Errorf("ModifyWithClone() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigStore_ModifyWithNew(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	webSrvOpts := metav1.WebServerOptions{
		GroupID:    1,
		HostID:     1,
		ServerName: "test-bifrost",
	}

	// Use fake data to get config and fingerprints
	fakeSvc := fake.WebServerConfigService{}
	_, fakeFP, _ := fakeSvc.Get("test-bifrost")
	ofp := fakeFP.Fingerprints()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBifrostsManager := bifrosts.NewMockManager(ctrl)
	mockBifrostsManager.EXPECT().GetBifrostClient(webSrvOpts).AnyTimes().Return(&bifrostclinetv1.Client{Factory: new(fake.ServiceClient)}, nil)
	type fields struct {
		bm bifrosts.Manager
	}
	type args struct {
		in0     context.Context
		opts    metav1.WebServerOptions
		ofp     utilsV3.ConfigFingerprints
		ctxmeta metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "normal test",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				opts: webSrvOpts,
				ofp:  ofp,
				ctxmeta: metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]{
					Position: metav1.ConfigContextPos{
						Config:         "C:\\config_test\\nginx.conf",
						ContextPosPath: []int{13, 15, 3},
					},
					TargetContext: metav1.NewConfigContextMeta{
						ContextType:  "location",
						ContextValue: "~ /normal-test",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigStore{
				bm: tt.fields.bm,
			}
			if err := w.ModifyWithNew(tt.args.in0, tt.args.opts, tt.args.ofp, tt.args.ctxmeta); (err != nil) != tt.wantErr {
				t.Errorf("ModifyWithNew() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigStore_Remove(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	webSrvOpts := metav1.WebServerOptions{
		GroupID:    1,
		HostID:     1,
		ServerName: "test-bifrost",
	}

	// Use fake data to get config and fingerprints
	fakeSvc := fake.WebServerConfigService{}
	_, fakeFP, _ := fakeSvc.Get("test-bifrost")
	ofp := fakeFP.Fingerprints()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBifrostsManager := bifrosts.NewMockManager(ctrl)
	mockBifrostsManager.EXPECT().GetBifrostClient(webSrvOpts).AnyTimes().Return(&bifrostclinetv1.Client{Factory: new(fake.ServiceClient)}, nil)
	type fields struct {
		bm bifrosts.Manager
	}
	type args struct {
		ctx  context.Context
		ofp  utilsV3.ConfigFingerprints
		opts metav1.WebServerOptions
		pos  metav1.ConfigContextPos
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "normal test",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				opts: webSrvOpts,
				ofp:  ofp,
				pos: metav1.ConfigContextPos{
					Config:         "C:\\config_test\\conf.d\\location.conf",
					ContextPosPath: []int{0},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigStore{
				bm: tt.fields.bm,
			}
			if err := w.Remove(tt.args.ctx, tt.args.opts, tt.args.ofp, tt.args.pos); (err != nil) != tt.wantErr {
				t.Errorf("Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigStore_Move(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	webSrvOpts := metav1.WebServerOptions{
		GroupID:    1,
		HostID:     1,
		ServerName: "test-bifrost",
	}

	// Use fake data to get config and fingerprints
	fakeSvc := fake.WebServerConfigService{}
	_, fakeFP, _ := fakeSvc.Get("test-bifrost")
	ofp := fakeFP.Fingerprints()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBifrostsManager := bifrosts.NewMockManager(ctrl)
	mockBifrostsManager.EXPECT().GetBifrostClient(webSrvOpts).AnyTimes().Return(&bifrostclinetv1.Client{Factory: new(fake.ServiceClient)}, nil)
	type fields struct {
		bm bifrosts.Manager
	}
	type args struct {
		in0            context.Context
		opts           metav1.WebServerOptions
		ofp            utilsV3.ConfigFingerprints
		ctxmeta        metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]
		disabledTarget bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "normal test",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				opts: webSrvOpts,
				ofp:  ofp,
				ctxmeta: metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]{
					Position: metav1.ConfigContextPos{
						Config:         "C:\\config_test\\nginx.conf",
						ContextPosPath: []int{13, 15, 3},
					},
					TargetContext: metav1.CloneConfigContextMeta{ConfigContextPos: metav1.ConfigContextPos{
						Config:         "C:\\config_test\\conf.d\\location1.conf",
						ContextPosPath: []int{0},
					}},
				},
				disabledTarget: true,
			},
			wantErr: false,
		},
		{
			name:   "index of context to be moved, is out of range",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				in0:  nil,
				opts: webSrvOpts,
				ofp:  ofp,
				ctxmeta: metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]{
					Position: metav1.ConfigContextPos{
						Config:         "C:\\config_test\\nginx.conf",
						ContextPosPath: []int{8, 13, 4},
					},
					TargetContext: metav1.CloneConfigContextMeta{ConfigContextPos: metav1.ConfigContextPos{
						Config:         "C:\\config_test\\conf.d\\location2.conf",
						ContextPosPath: []int{8},
					}},
				},
				disabledTarget: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigStore{
				bm: tt.fields.bm,
			}
			if err := w.Move(tt.args.in0, tt.args.opts, tt.args.ofp, tt.args.ctxmeta, tt.args.disabledTarget); (err != nil) != tt.wantErr {
				t.Errorf("Move() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigStore_SearchContextPositions(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	webSrvOpts := metav1.WebServerOptions{
		GroupID:    1,
		HostID:     1,
		ServerName: "test-bifrost",
	}

	// Use fake data to get config and fingerprints
	fakeSvc := fake.WebServerConfigService{}
	_, fakeFP, _ := fakeSvc.Get("test-bifrost")
	ofp := fakeFP.Fingerprints()

	// Create diff fingerprint for error test case
	difffp := utilsV3.ConfigFingerprints{
		"C:\\config_test\\nginx.conf": "1111111111111111111111111111111111111111111111111111111111111111",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBifrostsManager := bifrosts.NewMockManager(ctrl)
	mockBifrostsManager.EXPECT().GetBifrostClient(webSrvOpts).AnyTimes().Return(&bifrostclinetv1.Client{Factory: new(fake.ServiceClient)}, nil)
	type fields struct {
		bm bifrosts.Manager
	}
	type args struct {
		ctx    context.Context
		opts   metav1.WebServerOptions
		fp     utilsV3.ConfigFingerprints
		kwmeta metav1.SearchKeywordsMeta
	}
	tests := []struct {
		name                    string
		fields                  fields
		args                    args
		want                    []metav1.ConfigContextPos
		wantErr                 bool
		wantErrIsInconsistentFP bool
	}{
		{
			name:   "string match rule, only in current config",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				ctx:  nil,
				opts: webSrvOpts,
				fp:   ofp,
				kwmeta: metav1.SearchKeywordsMeta{
					StartingPositionList: []metav1.ConfigContextPos{{Config: "C:\\config_test\\nginx.conf"}},
					Keywords:             "listen 80",
					IsRegexpRule:         false,
					IsOnlyInCurrent:      true,
				},
			},
			want: []metav1.ConfigContextPos{
				{"C:\\config_test\\nginx.conf", []int{13, 12, 0}},
				{"C:\\config_test\\nginx.conf", []int{13, 15, 0}},
			},
		},
		{
			name:   "regexp match rule, only in current config",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				ctx:  nil,
				opts: webSrvOpts,
				fp:   ofp,
				kwmeta: metav1.SearchKeywordsMeta{
					StartingPositionList: []metav1.ConfigContextPos{{Config: "C:\\config_test\\nginx.conf"}},
					Keywords:             `^server_name\s+local.*$`,
					IsRegexpRule:         true,
					IsOnlyInCurrent:      true,
				},
			},
			want: []metav1.ConfigContextPos{
				{"C:\\config_test\\nginx.conf", []int{13, 12, 1}},
				{"C:\\config_test\\nginx.conf", []int{13, 18, 1}},
			},
		},
		{
			name:   "string match rule, not only in current config",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				ctx:  nil,
				opts: webSrvOpts,
				fp:   ofp,
				kwmeta: metav1.SearchKeywordsMeta{
					StartingPositionList: []metav1.ConfigContextPos{{Config: "C:\\config_test\\nginx.conf", ContextPosPath: []int{13}}},
					Keywords:             "listen 80",
					IsRegexpRule:         false,
					IsOnlyInCurrent:      false,
				},
			},
			want: []metav1.ConfigContextPos{
				{"C:\\config_test\\conf.d\\test1.com.conf", []int{0, 0}},
				{"C:\\config_test\\conf.d\\test2.com.conf", []int{0, 0}},
				{"C:\\config_test\\conf.d\\test2.com.conf", []int{1, 0}},
				{"C:\\config_test\\nginx.conf", []int{13, 12, 0}},
				{"C:\\config_test\\nginx.conf", []int{13, 15, 0}},
			},
		},
		{
			name:   "regexp match rule, not only in current config",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				ctx:  nil,
				opts: webSrvOpts,
				fp:   ofp,
				kwmeta: metav1.SearchKeywordsMeta{
					StartingPositionList: []metav1.ConfigContextPos{{Config: "C:\\config_test\\nginx.conf", ContextPosPath: []int{13}}},
					Keywords:             `^server_name\s+test.*$`,
					IsRegexpRule:         true,
					IsOnlyInCurrent:      false,
				},
			},
			want: []metav1.ConfigContextPos{
				{"C:\\config_test\\conf.d\\test1.com.conf", []int{0, 1}},
				{"C:\\config_test\\conf.d\\test2.com.conf", []int{0, 1}},
				{"C:\\config_test\\conf.d\\test2.com.conf", []int{1, 1}},
			},
		},
		{
			name:   "context not found",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				ctx:  nil,
				opts: webSrvOpts,
				fp:   ofp,
				kwmeta: metav1.SearchKeywordsMeta{
					StartingPositionList: []metav1.ConfigContextPos{{Config: "C:\\config_test\\nginx.conf", ContextPosPath: []int{8}}},
					Keywords:             ".*not found.*",
					IsRegexpRule:         true,
					IsOnlyInCurrent:      false,
				},
			},
			want: []metav1.ConfigContextPos{},
		},
		{
			name:   "config not found",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				ctx:  nil,
				opts: webSrvOpts,
				fp:   ofp,
				kwmeta: metav1.SearchKeywordsMeta{
					StartingPositionList: []metav1.ConfigContextPos{{Config: "C:\\config_test\\unknown.conf", ContextPosPath: []int{8}}},
					Keywords:             ".*not found.*",
					IsRegexpRule:         true,
					IsOnlyInCurrent:      false,
				},
			},
			wantErr: true,
		},
		{
			name:   "inconsistent fingerprints",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				ctx:  nil,
				opts: webSrvOpts,
				fp:   difffp,
				kwmeta: metav1.SearchKeywordsMeta{
					StartingPositionList: []metav1.ConfigContextPos{{Config: "C:\\config_test\\nginx.conf", ContextPosPath: []int{8}}},
					Keywords:             "listen 80",
					IsRegexpRule:         false,
					IsOnlyInCurrent:      false,
				},
			},
			wantErr:                 true,
			wantErrIsInconsistentFP: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigStore{
				bm: tt.fields.bm,
			}
			got, err := w.SearchContextPositions(tt.args.ctx, tt.args.opts, tt.args.fp, tt.args.kwmeta)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchContextPositions() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err != nil && (errors.Is(err, metav1.ErrInconsistentFingerprints) || errors.IsCode(err, 110010)) != tt.wantErrIsInconsistentFP {
				t.Errorf("SearchContextPositions() error = %v, wantErrIsInconsistentFP %v", err, tt.wantErrIsInconsistentFP)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchContextPositions() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webServerConfigStore_UpdateConfig(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	webSrvOpts := metav1.WebServerOptions{
		GroupID:    1,
		HostID:     1,
		ServerName: "test-bifrost",
	}

	// Use fake data to get config and fingerprints
	fakeSvc := fake.WebServerConfigService{}
	_, fakeFP, _ := fakeSvc.Get("test-bifrost")
	ofp := fakeFP.Fingerprints()

	// Create diff fingerprint for error test case
	difffp := utilsV3.ConfigFingerprints{
		"C:\\config_test\\nginx.conf": "1111111111111111111111111111111111111111111111111111111111111111",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBifrostsManager := bifrosts.NewMockManager(ctrl)
	mockBifrostsManager.EXPECT().GetBifrostClient(webSrvOpts).AnyTimes().Return(&bifrostclinetv1.Client{Factory: new(fake.ServiceClient)}, nil)

	type fields struct {
		bm bifrosts.Manager
	}
	type args struct {
		in0            context.Context
		opts           metav1.WebServerOptions
		ofp            utilsV3.ConfigFingerprints
		configJsonData []byte
	}
	tests := []struct {
		name                    string
		fields                  fields
		args                    args
		wantErr                 bool
		wantErrIsInconsistentFP bool
	}{
		{
			name:   "normal test",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				in0:            nil,
				opts:           webSrvOpts,
				ofp:            ofp,
				configJsonData: []byte(`{"main-config":"C:\\config_test\\nginx.conf","configs":{}}`),
			},
			wantErr: false,
		},
		{
			name:   "invalid json data",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				in0:            nil,
				opts:           webSrvOpts,
				ofp:            ofp,
				configJsonData: []byte(`invalid json`),
			},
			wantErr: true,
		},
		{
			name:   "inconsistent fingerprints",
			fields: fields{bm: mockBifrostsManager},
			args: args{
				in0:            nil,
				opts:           webSrvOpts,
				ofp:            difffp,
				configJsonData: []byte(`{"main-config":"C:\\config_test\\nginx.conf","configs":{}}`),
			},
			wantErr:                 true,
			wantErrIsInconsistentFP: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigStore{
				bm: tt.fields.bm,
			}
			err := w.UpdateConfig(tt.args.in0, tt.args.opts, tt.args.ofp, tt.args.configJsonData)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateConfig() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && (errors.Is(err, metav1.ErrInconsistentFingerprints) || errors.IsCode(err, 110010)) != tt.wantErrIsInconsistentFP {
				t.Errorf("UpdateConfig() error = %v, wantErrIsInconsistentFP %v", err, tt.wantErrIsInconsistentFP)
			}
		})
	}
}
