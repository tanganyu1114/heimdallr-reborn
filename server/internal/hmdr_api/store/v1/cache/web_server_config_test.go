package cache

import (
	"context"
	v1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	"github.com/tanganyu1114/heimdallr-reborn/server/global"
	storev1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/store/v1"
	"github.com/tanganyu1114/heimdallr-reborn/server/internal/pkg/bifrosts/fake"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/pkg/meta/v1"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	utilsV3 "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/utils"
	logV1 "github.com/ClessLi/component-base/pkg/log/v1"
	"github.com/marmotedu/errors"
	"go.uber.org/mock/gomock"
)

func Test_newWebServerConfigStore(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFactory := storev1.NewMockFactory(ctrl)
	cache := &cacheStore{
		next:           mockFactory,
		expireDuration: time.Minute,
	}

	type args struct {
		cacheStore *cacheStore
	}
	tests := []struct {
		name string
		args args
		want *webServerConfigStore
	}{
		{
			name: "normal_test",
			args: args{cacheStore: cache},
			want: &webServerConfigStore{cacheStore: cache},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newWebServerConfigStore(tt.args.cacheStore); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newWebServerConfigStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webServerConfigStore_GetOptions(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFactory := storev1.NewMockFactory(ctrl)
	mockWebServerConfigStore := storev1.NewMockWebServerConfigStore(ctrl)

	cache := &cacheStore{
		next:           mockFactory,
		expireDuration: time.Minute,
	}

	mockFactory.EXPECT().WebServerConfigs().Return(mockWebServerConfigStore).AnyTimes()

	expectedResult := []v1.BifrostGroupMeta{
		{
			UintObjectMeta: metav1.UintObjectMeta{Label: "test-group", Value: 1},
			Children: []*v1.BifrostMeta{
				{
					UintObjectMeta: metav1.UintObjectMeta{Label: "test-host", Value: 1},
					Children:       []*v1.WebSrvMeta{},
				},
			},
		},
	}

	type fields struct {
		cacheStore *cacheStore
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
			name:    "normal_test",
			fields:  fields{cacheStore: cache},
			args:    args{ctx: context.Background()},
			want:    expectedResult,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWebServerConfigStore.EXPECT().GetOptions(tt.args.ctx).Return(expectedResult, nil)

			w := &webServerConfigStore{
				cacheStore: tt.fields.cacheStore,
			}
			got, err := w.GetOptions(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOptions() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webServerConfigStore_GetConfig(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	opts := metav1.WebServerOptions{GroupID: 1, HostID: 1}

	fakeSvc := fake.WebServerConfigService{}
	expectedConfig, expectedFp, _ := fakeSvc.Get("test-server")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFactory := storev1.NewMockFactory(ctrl)
	mockWebServerConfigStore := storev1.NewMockWebServerConfigStore(ctrl)

	mockFactory.EXPECT().WebServerConfigs().Return(mockWebServerConfigStore).AnyTimes()
	mockWebServerConfigStore.EXPECT().GetConfig(context.Background(), opts).Return(expectedConfig, expectedFp, nil).AnyTimes()

	cache := &cacheStore{
		next: mockFactory,
		cache: &groupsCache{
			cache:  make(map[uint]*hostsCache),
			locker: &sync.RWMutex{},
		},
		expireDuration: time.Minute,
	}

	type fields struct {
		cacheStore *cacheStore
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
			name:   "normal_test",
			fields: fields{cacheStore: cache},
			args:   args{ctx: context.Background(), opts: opts},
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
			wantOriginalFingerprints: expectedFp.Fingerprints(),
			wantErr:                  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigStore{
				cacheStore: tt.fields.cacheStore,
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
	opts := metav1.WebServerOptions{GroupID: 1, HostID: 1}

	fakeSvc := fake.WebServerConfigService{}
	_, fakeFP, _ := fakeSvc.Get("test-server")
	ofp := fakeFP.Fingerprints()

	difffp := utilsV3.ConfigFingerprints{
		"C:\\config_test\\nginx.conf": "1111111111111111111111111111111111111111111111111111111111111111",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFactory := storev1.NewMockFactory(ctrl)
	mockWebServerConfigStore := storev1.NewMockWebServerConfigStore(ctrl)

	mockFactory.EXPECT().WebServerConfigs().Return(mockWebServerConfigStore).AnyTimes()
	mockWebServerConfigStore.EXPECT().GetConfig(gomock.Any(), opts).AnyTimes().DoAndReturn(func(ctx context.Context, opts metav1.WebServerOptions) (configuration.NginxConfig, utilsV3.ConfigFingerprinter, error) {
		return fakeSvc.Get("test-server")
	})

	cache := &cacheStore{
		next: mockFactory,
		cache: &groupsCache{
			cache:  make(map[uint]*hostsCache),
			locker: &sync.RWMutex{},
		},
		expireDuration: time.Minute,
	}

	type fields struct {
		cacheStore *cacheStore
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
			fields: fields{cacheStore: cache},
			args: args{
				opts: opts,
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
			fields: fields{cacheStore: cache},
			args: args{
				opts: opts,
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
			fields: fields{cacheStore: cache},
			args: args{
				opts: opts,
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
			wantErr: false,
		},
		{
			name:   "wrong position",
			fields: fields{cacheStore: cache},
			args: args{
				opts: opts,
				fp:   ofp,
				pos: metav1.ConfigContextPos{
					Config:         "C:\\config_test\\nginx.conf",
					ContextPosPath: []int{1, 2, 3, 4, 5, 6, 7},
				},
			},
			wantConfigLines: nil,
			wantErr:         true,
		},
		{
			name:   "inconsistent fingerprints",
			fields: fields{cacheStore: cache},
			args: args{
				opts: opts,
				fp:   difffp,
				pos: metav1.ConfigContextPos{
					Config:         "C:\\config_test\\nginx.conf",
					ContextPosPath: []int{0},
				},
			},
			wantConfigLines:         nil,
			wantErr:                 true,
			wantErrIsInconsistentFP: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigStore{
				cacheStore: tt.fields.cacheStore,
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
	opts := metav1.WebServerOptions{GroupID: 1, HostID: 1}

	fakeSvc := fake.WebServerConfigService{}
	_, fakeFP, _ := fakeSvc.Get("test-server")
	ofp := fakeFP.Fingerprints()

	difffp := utilsV3.ConfigFingerprints{
		"C:\\config_test\\nginx.conf": "1111111111111111111111111111111111111111111111111111111111111111",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFactory := storev1.NewMockFactory(ctrl)
	mockWebServerConfigStore := storev1.NewMockWebServerConfigStore(ctrl)

	mockFactory.EXPECT().WebServerConfigs().Return(mockWebServerConfigStore).AnyTimes()
	mockWebServerConfigStore.EXPECT().GetConfig(context.Background(), opts).AnyTimes().DoAndReturn(func(ctx context.Context, opts metav1.WebServerOptions) (configuration.NginxConfig, utilsV3.ConfigFingerprinter, error) {
		return fakeSvc.Get("test-server")
	})

	cache := &cacheStore{
		next: mockFactory,
		cache: &groupsCache{
			cache:  make(map[uint]*hostsCache),
			locker: &sync.RWMutex{},
		},
		expireDuration: time.Minute,
	}

	type fields struct {
		cacheStore *cacheStore
	}
	type args struct {
		ctx  context.Context
		opts metav1.WebServerOptions
		ofp  utilsV3.ConfigFingerprints
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
			name:   "normal_test",
			fields: fields{cacheStore: cache},
			args: args{
				ctx:  context.Background(),
				opts: opts,
				ofp:  ofp,
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
			fields: fields{cacheStore: cache},
			args: args{
				ctx:  context.Background(),
				opts: opts,
				ofp:  ofp,
				pos: metav1.ConfigContextPos{
					Config:         "C:\\config_test\\conf.d\\server_test1.con",
					ContextPosPath: nil,
				},
			},
			wantErr: true,
		},
		{
			name:   "the target is not an `include` context",
			fields: fields{cacheStore: cache},
			args: args{
				ctx:  context.Background(),
				opts: opts,
				ofp:  ofp,
				pos: metav1.ConfigContextPos{
					Config:         "C:\\config_test\\nginx.conf",
					ContextPosPath: []int{8, 1},
				},
			},
			wantErr: true,
		},
		{
			name:   "inconsistent fingerprints",
			fields: fields{cacheStore: cache},
			args: args{
				ctx:  context.Background(),
				opts: opts,
				ofp:  difffp,
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
				cacheStore: tt.fields.cacheStore,
			}
			got, err := w.GetIncludedConfigs(tt.args.ctx, tt.args.opts, tt.args.ofp, tt.args.pos)
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

func Test_webServerConfigStore_SearchContextPositions(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	opts := metav1.WebServerOptions{GroupID: 1, HostID: 1}

	fakeSvc := fake.WebServerConfigService{}
	_, fakeFP, _ := fakeSvc.Get("test-server")
	ofp := fakeFP.Fingerprints()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFactory := storev1.NewMockFactory(ctrl)
	mockWebServerConfigStore := storev1.NewMockWebServerConfigStore(ctrl)

	mockFactory.EXPECT().WebServerConfigs().Return(mockWebServerConfigStore).AnyTimes()
	mockWebServerConfigStore.EXPECT().GetConfig(context.Background(), opts).AnyTimes().DoAndReturn(func(ctx context.Context, opts metav1.WebServerOptions) (configuration.NginxConfig, utilsV3.ConfigFingerprinter, error) {
		return fakeSvc.Get("test-server")
	})

	cache := &cacheStore{
		next: mockFactory,
		cache: &groupsCache{
			cache:  make(map[uint]*hostsCache),
			locker: &sync.RWMutex{},
		},
		expireDuration: time.Minute,
	}

	type fields struct {
		cacheStore *cacheStore
	}
	type args struct {
		ctx    context.Context
		opts   metav1.WebServerOptions
		fp     utilsV3.ConfigFingerprints
		kwmeta metav1.SearchKeywordsMeta
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "normal_test",
			fields: fields{cacheStore: cache},
			args: args{
				ctx:  context.Background(),
				opts: opts,
				fp:   ofp,
				kwmeta: metav1.SearchKeywordsMeta{
					Keywords:             "proxy_pass",
					IsOnlyInCurrent:      false,
					IsRegexpRule:         false,
					StartingPositionList: []metav1.ConfigContextPos{},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigStore{
				cacheStore: tt.fields.cacheStore,
			}
			_, err := w.SearchContextPositions(tt.args.ctx, tt.args.opts, tt.args.fp, tt.args.kwmeta)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchContextPositions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigStore_InsertWithClone(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	opts := metav1.WebServerOptions{GroupID: 1, HostID: 1}
	ofp := utilsV3.ConfigFingerprints{}
	ctxmeta := metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]{
		Position: metav1.ConfigContextPos{
			Config:         "C:\\config_test\\nginx.conf",
			ContextPosPath: []int{0},
		},
		TargetContext: metav1.CloneConfigContextMeta{
			ConfigContextPos: metav1.ConfigContextPos{
				Config:         "C:\\config_test\\nginx.conf",
				ContextPosPath: []int{1},
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFactory := storev1.NewMockFactory(ctrl)
	mockWebServerConfigStore := storev1.NewMockWebServerConfigStore(ctrl)

	mockFactory.EXPECT().WebServerConfigs().Return(mockWebServerConfigStore).AnyTimes()
	mockWebServerConfigStore.EXPECT().InsertWithClone(context.Background(), opts, ofp, ctxmeta, false).Return(nil)

	cache := &cacheStore{
		next: mockFactory,
		cache: &groupsCache{
			cache:  make(map[uint]*hostsCache),
			locker: &sync.RWMutex{},
		},
		expireDuration: time.Minute,
	}

	type fields struct {
		cacheStore *cacheStore
	}
	type args struct {
		ctx            context.Context
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
			name:    "normal_test",
			fields:  fields{cacheStore: cache},
			args:    args{ctx: context.Background(), opts: opts, ofp: ofp, ctxmeta: ctxmeta, disabledTarget: false},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigStore{
				cacheStore: tt.fields.cacheStore,
			}
			if err := w.InsertWithClone(tt.args.ctx, tt.args.opts, tt.args.ofp, tt.args.ctxmeta, tt.args.disabledTarget); (err != nil) != tt.wantErr {
				t.Errorf("InsertWithClone() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigStore_InsertWithNew(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	opts := metav1.WebServerOptions{GroupID: 1, HostID: 1}
	ofp := utilsV3.ConfigFingerprints{}
	ctxmeta := metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]{
		Position: metav1.ConfigContextPos{
			Config:         "C:\\config_test\\nginx.conf",
			ContextPosPath: []int{0},
		},
		TargetContext: metav1.NewConfigContextMeta{
			Enabled:      true,
			ContextType:  context_type.TypeServer,
			ContextValue: "server {}",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFactory := storev1.NewMockFactory(ctrl)
	mockWebServerConfigStore := storev1.NewMockWebServerConfigStore(ctrl)

	mockFactory.EXPECT().WebServerConfigs().Return(mockWebServerConfigStore).AnyTimes()
	mockWebServerConfigStore.EXPECT().InsertWithNew(context.Background(), opts, ofp, ctxmeta, false).Return(nil)

	cache := &cacheStore{
		next: mockFactory,
		cache: &groupsCache{
			cache:  make(map[uint]*hostsCache),
			locker: &sync.RWMutex{},
		},
		expireDuration: time.Minute,
	}

	type fields struct {
		cacheStore *cacheStore
	}
	type args struct {
		ctx            context.Context
		opts           metav1.WebServerOptions
		ofp            utilsV3.ConfigFingerprints
		ctxmeta        metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]
		disabledTarget bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "normal_test",
			fields:  fields{cacheStore: cache},
			args:    args{ctx: context.Background(), opts: opts, ofp: ofp, ctxmeta: ctxmeta, disabledTarget: false},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigStore{
				cacheStore: tt.fields.cacheStore,
			}
			if err := w.InsertWithNew(tt.args.ctx, tt.args.opts, tt.args.ofp, tt.args.ctxmeta, tt.args.disabledTarget); (err != nil) != tt.wantErr {
				t.Errorf("InsertWithNew() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigStore_Remove(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	opts := metav1.WebServerOptions{GroupID: 1, HostID: 1}
	ofp := utilsV3.ConfigFingerprints{}
	pos := metav1.ConfigContextPos{
		Config:         "C:\\config_test\\nginx.conf",
		ContextPosPath: []int{0},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFactory := storev1.NewMockFactory(ctrl)
	mockWebServerConfigStore := storev1.NewMockWebServerConfigStore(ctrl)

	mockFactory.EXPECT().WebServerConfigs().Return(mockWebServerConfigStore).AnyTimes()
	mockWebServerConfigStore.EXPECT().Remove(context.Background(), opts, ofp, pos).Return(nil)

	cache := &cacheStore{
		next: mockFactory,
		cache: &groupsCache{
			cache:  make(map[uint]*hostsCache),
			locker: &sync.RWMutex{},
		},
		expireDuration: time.Minute,
	}

	type fields struct {
		cacheStore *cacheStore
	}
	type args struct {
		ctx  context.Context
		opts metav1.WebServerOptions
		ofp  utilsV3.ConfigFingerprints
		pos  metav1.ConfigContextPos
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "normal_test",
			fields:  fields{cacheStore: cache},
			args:    args{ctx: context.Background(), opts: opts, ofp: ofp, pos: pos},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigStore{
				cacheStore: tt.fields.cacheStore,
			}
			if err := w.Remove(tt.args.ctx, tt.args.opts, tt.args.ofp, tt.args.pos); (err != nil) != tt.wantErr {
				t.Errorf("Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigStore_UpdateConfig(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	opts := metav1.WebServerOptions{GroupID: 1, HostID: 1}
	ofp := utilsV3.ConfigFingerprints{}
	configJsonData := []byte(`{"main-config":"C:\\config_test\\nginx.conf"}`)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFactory := storev1.NewMockFactory(ctrl)
	mockWebServerConfigStore := storev1.NewMockWebServerConfigStore(ctrl)

	mockFactory.EXPECT().WebServerConfigs().Return(mockWebServerConfigStore).AnyTimes()
	mockWebServerConfigStore.EXPECT().UpdateConfig(context.Background(), opts, ofp, configJsonData).Return(nil)

	cache := &cacheStore{
		next: mockFactory,
		cache: &groupsCache{
			cache:  make(map[uint]*hostsCache),
			locker: &sync.RWMutex{},
		},
		expireDuration: time.Minute,
	}

	type fields struct {
		cacheStore *cacheStore
	}
	type args struct {
		ctx            context.Context
		opts           metav1.WebServerOptions
		ofp            utilsV3.ConfigFingerprints
		configJsonData []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "normal_test",
			fields:  fields{cacheStore: cache},
			args:    args{ctx: context.Background(), opts: opts, ofp: ofp, configJsonData: configJsonData},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigStore{
				cacheStore: tt.fields.cacheStore,
			}
			if err := w.UpdateConfig(tt.args.ctx, tt.args.opts, tt.args.ofp, tt.args.configJsonData); (err != nil) != tt.wantErr {
				t.Errorf("UpdateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigStore_ModifyWithClone(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	opts := metav1.WebServerOptions{GroupID: 1, HostID: 1}
	ofp := utilsV3.ConfigFingerprints{}
	ctxmeta := metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]{
		Position: metav1.ConfigContextPos{
			Config:         "C:\\config_test\\nginx.conf",
			ContextPosPath: []int{0},
		},
		TargetContext: metav1.CloneConfigContextMeta{
			ConfigContextPos: metav1.ConfigContextPos{
				Config:         "C:\\config_test\\nginx.conf",
				ContextPosPath: []int{1},
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFactory := storev1.NewMockFactory(ctrl)
	mockWebServerConfigStore := storev1.NewMockWebServerConfigStore(ctrl)

	mockFactory.EXPECT().WebServerConfigs().Return(mockWebServerConfigStore).AnyTimes()
	mockWebServerConfigStore.EXPECT().ModifyWithClone(context.Background(), opts, ofp, ctxmeta).Return(nil)

	cache := &cacheStore{
		next: mockFactory,
		cache: &groupsCache{
			cache:  make(map[uint]*hostsCache),
			locker: &sync.RWMutex{},
		},
		expireDuration: time.Minute,
	}

	type fields struct {
		cacheStore *cacheStore
	}
	type args struct {
		ctx     context.Context
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
			name:    "normal_test",
			fields:  fields{cacheStore: cache},
			args:    args{ctx: context.Background(), opts: opts, ofp: ofp, ctxmeta: ctxmeta},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigStore{
				cacheStore: tt.fields.cacheStore,
			}
			if err := w.ModifyWithClone(tt.args.ctx, tt.args.opts, tt.args.ofp, tt.args.ctxmeta); (err != nil) != tt.wantErr {
				t.Errorf("ModifyWithClone() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigStore_ModifyWithNew(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	opts := metav1.WebServerOptions{GroupID: 1, HostID: 1}
	ofp := utilsV3.ConfigFingerprints{}
	ctxmeta := metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]{
		Position: metav1.ConfigContextPos{
			Config:         "C:\\config_test\\nginx.conf",
			ContextPosPath: []int{0},
		},
		TargetContext: metav1.NewConfigContextMeta{
			Enabled:      true,
			ContextType:  context_type.TypeServer,
			ContextValue: "server {}",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFactory := storev1.NewMockFactory(ctrl)
	mockWebServerConfigStore := storev1.NewMockWebServerConfigStore(ctrl)

	mockFactory.EXPECT().WebServerConfigs().Return(mockWebServerConfigStore).AnyTimes()
	mockWebServerConfigStore.EXPECT().ModifyWithNew(context.Background(), opts, ofp, ctxmeta).Return(nil)

	cache := &cacheStore{
		next: mockFactory,
		cache: &groupsCache{
			cache:  make(map[uint]*hostsCache),
			locker: &sync.RWMutex{},
		},
		expireDuration: time.Minute,
	}

	type fields struct {
		cacheStore *cacheStore
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
			name:    "normal_test",
			fields:  fields{cacheStore: cache},
			args:    args{ctx: context.Background(), opts: opts, ofp: ofp, ctxmeta: ctxmeta},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigStore{
				cacheStore: tt.fields.cacheStore,
			}
			if err := w.ModifyWithNew(tt.args.ctx, tt.args.opts, tt.args.ofp, tt.args.ctxmeta); (err != nil) != tt.wantErr {
				t.Errorf("ModifyWithNew() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigStore_ChangeContextEnabledState(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	opts := metav1.WebServerOptions{GroupID: 1, HostID: 1}
	ofp := utilsV3.ConfigFingerprints{}
	ctxmeta := metav1.TargetConfigContextOptions[metav1.ConfigContextEnabledStateMeta]{
		Position: metav1.ConfigContextPos{
			Config:         "C:\\config_test\\nginx.conf",
			ContextPosPath: []int{0},
		},
		TargetContext: metav1.ConfigContextEnabledStateMeta{
			Enabled: true,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFactory := storev1.NewMockFactory(ctrl)
	mockWebServerConfigStore := storev1.NewMockWebServerConfigStore(ctrl)

	mockFactory.EXPECT().WebServerConfigs().Return(mockWebServerConfigStore).AnyTimes()
	mockWebServerConfigStore.EXPECT().ChangeContextEnabledState(context.Background(), opts, ofp, ctxmeta).Return(nil)

	cache := &cacheStore{
		next: mockFactory,
		cache: &groupsCache{
			cache:  make(map[uint]*hostsCache),
			locker: &sync.RWMutex{},
		},
		expireDuration: time.Minute,
	}

	type fields struct {
		cacheStore *cacheStore
	}
	type args struct {
		ctx     context.Context
		opts    metav1.WebServerOptions
		ofp     utilsV3.ConfigFingerprints
		ctxmeta metav1.TargetConfigContextOptions[metav1.ConfigContextEnabledStateMeta]
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "normal_test",
			fields:  fields{cacheStore: cache},
			args:    args{ctx: context.Background(), opts: opts, ofp: ofp, ctxmeta: ctxmeta},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigStore{
				cacheStore: tt.fields.cacheStore,
			}
			if err := w.ChangeContextEnabledState(tt.args.ctx, tt.args.opts, tt.args.ofp, tt.args.ctxmeta); (err != nil) != tt.wantErr {
				t.Errorf("ChangeContextEnabledState() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigStore_ModifyContextValue(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	opts := metav1.WebServerOptions{GroupID: 1, HostID: 1}
	ofp := utilsV3.ConfigFingerprints{}
	ctxmeta := metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]{
		Position: metav1.ConfigContextPos{
			Config:         "C:\\config_test\\nginx.conf",
			ContextPosPath: []int{0},
		},
		TargetContext: metav1.NewConfigContextMeta{
			Enabled:      true,
			ContextType:  context_type.TypeDirective,
			ContextValue: "worker_processes 2",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFactory := storev1.NewMockFactory(ctrl)
	mockWebServerConfigStore := storev1.NewMockWebServerConfigStore(ctrl)

	mockFactory.EXPECT().WebServerConfigs().Return(mockWebServerConfigStore).AnyTimes()
	mockWebServerConfigStore.EXPECT().ModifyContextValue(context.Background(), opts, ofp, ctxmeta).Return(nil)

	cache := &cacheStore{
		next: mockFactory,
		cache: &groupsCache{
			cache:  make(map[uint]*hostsCache),
			locker: &sync.RWMutex{},
		},
		expireDuration: time.Minute,
	}

	type fields struct {
		cacheStore *cacheStore
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
			name:    "normal_test",
			fields:  fields{cacheStore: cache},
			args:    args{ctx: context.Background(), opts: opts, ofp: ofp, ctxmeta: ctxmeta},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigStore{
				cacheStore: tt.fields.cacheStore,
			}
			if err := w.ModifyContextValue(tt.args.ctx, tt.args.opts, tt.args.ofp, tt.args.ctxmeta); (err != nil) != tt.wantErr {
				t.Errorf("ModifyContextValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigStore_Move(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	opts := metav1.WebServerOptions{GroupID: 1, HostID: 1}
	ofp := utilsV3.ConfigFingerprints{}
	ctxmeta := metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]{
		Position: metav1.ConfigContextPos{
			Config:         "C:\\config_test\\nginx.conf",
			ContextPosPath: []int{0},
		},
		TargetContext: metav1.CloneConfigContextMeta{
			ConfigContextPos: metav1.ConfigContextPos{
				Config:         "C:\\config_test\\nginx.conf",
				ContextPosPath: []int{1},
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFactory := storev1.NewMockFactory(ctrl)
	mockWebServerConfigStore := storev1.NewMockWebServerConfigStore(ctrl)

	mockFactory.EXPECT().WebServerConfigs().Return(mockWebServerConfigStore).AnyTimes()
	mockWebServerConfigStore.EXPECT().Move(context.Background(), opts, ofp, ctxmeta, false).Return(nil)

	cache := &cacheStore{
		next: mockFactory,
		cache: &groupsCache{
			cache:  make(map[uint]*hostsCache),
			locker: &sync.RWMutex{},
		},
		expireDuration: time.Minute,
	}

	type fields struct {
		cacheStore *cacheStore
	}
	type args struct {
		ctx            context.Context
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
			name:    "normal_test",
			fields:  fields{cacheStore: cache},
			args:    args{ctx: context.Background(), opts: opts, ofp: ofp, ctxmeta: ctxmeta, disabledTarget: false},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigStore{
				cacheStore: tt.fields.cacheStore,
			}
			if err := w.Move(tt.args.ctx, tt.args.opts, tt.args.ofp, tt.args.ctxmeta, tt.args.disabledTarget); (err != nil) != tt.wantErr {
				t.Errorf("Move() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigStore_getConfigAndVerifyOFP(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	opts := metav1.WebServerOptions{GroupID: 1, HostID: 1}

	fakeSvc := fake.WebServerConfigService{}
	expectedConfig, expectedFp, _ := fakeSvc.Get("test-server")
	fp := expectedFp.Fingerprints()

	difffp := utilsV3.ConfigFingerprints{
		"C:\\config_test\\nginx.conf": "1111111111111111111111111111111111111111111111111111111111111111",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFactory := storev1.NewMockFactory(ctrl)
	mockWebServerConfigStore := storev1.NewMockWebServerConfigStore(ctrl)

	mockFactory.EXPECT().WebServerConfigs().Return(mockWebServerConfigStore).AnyTimes()
	mockWebServerConfigStore.EXPECT().GetConfig(context.Background(), opts).AnyTimes().DoAndReturn(func(ctx context.Context, opts metav1.WebServerOptions) (configuration.NginxConfig, utilsV3.ConfigFingerprinter, error) {
		return fakeSvc.Get("test-server")
	})

	cache := &cacheStore{
		next: mockFactory,
		cache: &groupsCache{
			cache:  make(map[uint]*hostsCache),
			locker: &sync.RWMutex{},
		},
		expireDuration: time.Minute,
	}

	type fields struct {
		cacheStore *cacheStore
	}
	type args struct {
		ctx  context.Context
		opts metav1.WebServerOptions
		fp   utilsV3.ConfigFingerprints
	}
	tests := []struct {
		name                    string
		fields                  fields
		args                    args
		wantErr                 bool
		wantErrIsInconsistentFP bool
	}{
		{
			name:    "matching fingerprints",
			fields:  fields{cacheStore: cache},
			args:    args{ctx: context.Background(), opts: opts, fp: fp},
			wantErr: false,
		},
		{
			name:                    "mismatching fingerprints",
			fields:                  fields{cacheStore: cache},
			args:                    args{ctx: context.Background(), opts: opts, fp: difffp},
			wantErr:                 true,
			wantErrIsInconsistentFP: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigStore{
				cacheStore: tt.fields.cacheStore,
			}
			got, err := w.getConfigAndVerifyOFP(tt.args.ctx, tt.args.opts, tt.args.fp)
			if (err != nil) != tt.wantErr {
				t.Errorf("getConfigAndVerifyOFP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && (errors.Is(err, metav1.ErrInconsistentFingerprints)) != tt.wantErrIsInconsistentFP {
				t.Errorf("getConfigAndVerifyOFP() error = %v, wantErrIsInconsistentFP %v", err, tt.wantErrIsInconsistentFP)
			}
			if !tt.wantErr && !reflect.DeepEqual(got.TextLines(), expectedConfig.TextLines()) {
				t.Errorf("getConfigAndVerifyOFP() got.TextLines() = %v, want %v", got.TextLines(), expectedConfig.TextLines())
			}
		})
	}
}
