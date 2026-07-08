package bifrosts

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	"gin-vue-admin/global"
	svcv1 "gin-vue-admin/internal/hmdr_api/service/v1"
	storev1 "gin-vue-admin/internal/hmdr_api/store/v1"
	"gin-vue-admin/internal/hmdr_api/store/v1/cache"
	"gin-vue-admin/internal/pkg/bifrosts/fake"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	"reflect"
	"strings"
	"testing"
	"time"

	nginx_context "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	utilsV3 "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/utils"
	log "github.com/ClessLi/component-base/pkg/log/v1"
	"github.com/marmotedu/errors"
	"go.uber.org/mock/gomock"
)

func Test_newWebServerConfigs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := storev1.NewMockFactory(ctrl)
	svc := &service{store: store}
	type args struct {
		svc *service
	}
	tests := []struct {
		name string
		args args
		want svcv1.WebServerConfigSrv
	}{
		{
			name: "normal test",
			args: args{svc: svc},
			want: &webServerConfigService{svc.store},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newWebServerConfigs(tt.args.svc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newWebServerConfigs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webServerConfigService_ChangeContextEnabledState(t *testing.T) {
	webSrvOpts := metav1.WebServerOptions{
		GroupID:    1,
		HostID:     1,
		ServerName: "test-bifrost",
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := storev1.NewMockFactory(ctrl)
	cacheStore := cache.GetCacheStore(store, time.Second)
	wscstore := storev1.NewMockWebServerConfigStore(ctrl)
	store.EXPECT().WebServerConfigs().AnyTimes().Return(wscstore)

	// Use fake data to get config and fingerprints
	fakeSvc := fake.WebServerConfigService{}
	fakeConfig, fakeFP, _ := fakeSvc.Get("test-bifrost")
	wscstore.EXPECT().GetConfig(nil, webSrvOpts).AnyTimes().Return(fakeConfig, fakeFP, nil)

	_, ofp, err := cacheStore.WebServerConfigs().GetConfig(nil, webSrvOpts)
	if err != nil {
		t.Fatal(err)
	}

	// Mock ChangeContextEnabledState to return nil (success)
	wscstore.EXPECT().ChangeContextEnabledState(nil, webSrvOpts, ofp.Fingerprints(), metav1.TargetConfigContextOptions[metav1.ConfigContextEnabledStateMeta]{
		Position: metav1.ConfigContextPos{
			Config:         "C:\\config_test\\nginx.conf",
			ContextPosPath: []int{0},
		},
		TargetContext: metav1.ConfigContextEnabledStateMeta{Enabled: true},
	}).AnyTimes().Return(nil)

	wscstore.EXPECT().ChangeContextEnabledState(nil, webSrvOpts, ofp.Fingerprints(), metav1.TargetConfigContextOptions[metav1.ConfigContextEnabledStateMeta]{
		Position: metav1.ConfigContextPos{
			Config:         "conf.d\\location.conf",
			ContextPosPath: nil,
		},
	}).AnyTimes().Return(nil)

	wscstore.EXPECT().ChangeContextEnabledState(nil, webSrvOpts, ofp.Fingerprints(), metav1.TargetConfigContextOptions[metav1.ConfigContextEnabledStateMeta]{
		Position: metav1.ConfigContextPos{
			Config:         "conf.d\\location.conf",
			ContextPosPath: []int{1, 2, 3},
		},
		TargetContext: metav1.ConfigContextEnabledStateMeta{Enabled: true},
	}).AnyTimes().Return(errors.New("invalid context position"))
	type fields struct {
		store storev1.Factory
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
			name:   "enable context",
			fields: fields{store: cacheStore},
			args: args{
				opts: webSrvOpts,
				ofp:  ofp.Fingerprints(),
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
			fields: fields{store: cacheStore},
			args: args{
				opts: webSrvOpts,
				ofp:  ofp.Fingerprints(),
				ctxmeta: metav1.TargetConfigContextOptions[metav1.ConfigContextEnabledStateMeta]{
					Position: metav1.ConfigContextPos{
						Config:         "conf.d\\location.conf",
						ContextPosPath: nil,
					},
					//TargetContext: metav1.ConfigContextEnabledStateMeta{Enabled: false},
				},
			},
			wantErr: false,
		},
		{
			name:   "wrong position",
			fields: fields{store: cacheStore},
			args: args{
				opts: webSrvOpts,
				ofp:  ofp.Fingerprints(),
				ctxmeta: metav1.TargetConfigContextOptions[metav1.ConfigContextEnabledStateMeta]{
					Position: metav1.ConfigContextPos{
						Config:         "conf.d\\location.conf",
						ContextPosPath: []int{1, 2, 3},
					},
					TargetContext: metav1.ConfigContextEnabledStateMeta{Enabled: true},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigService{
				store: tt.fields.store,
			}
			if err := w.ChangeContextEnabledState(tt.args.ctx, tt.args.opts, tt.args.ofp, tt.args.ctxmeta); (err != nil) != tt.wantErr {
				t.Errorf("ChangeContextEnabledState() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigService_GetConfig(t *testing.T) {
	webSrvOpts := metav1.WebServerOptions{
		GroupID:    1,
		HostID:     1,
		ServerName: "test-bifrost",
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := storev1.NewMockFactory(ctrl)
	wscstore := storev1.NewMockWebServerConfigStore(ctrl)
	store.EXPECT().WebServerConfigs().AnyTimes().Return(wscstore)

	// Use fake data to get config and fingerprints
	fakeSvc := fake.WebServerConfigService{}
	fakeConfig, fakeFP, _ := fakeSvc.Get("test-bifrost")
	wscstore.EXPECT().GetConfig(nil, webSrvOpts).AnyTimes().Return(fakeConfig, fakeFP, nil)
	type fields struct {
		store storev1.Factory
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
			fields: fields{store: store},
			args:   args{opts: webSrvOpts},
			wantConfigTextLines: []string{
				`# testConfigManager`,
				`# testConfigManager`,
				`# testConfigManager_update...`,
				`# testConfigManager_update...`,
				`# user nobody;`,
				`worker_processes 1;    # test inline comments`,
				`# error_log logs/error.log;`,
				`# error_log logs/error.log  notice;`,
				`# error_log logs/error.log  info;`,
				`# pid logs/nginx.pid;`,
				`events {`,
				`    worker_connections 1024;`,
				`}`,
				`stream {`,
				`    server {`,
				`        listen 8888;`,
				`        proxy_pass baidu.com:22;`,
				`    }`,
				`}`,
				`http {`,
				`    # include <== mime.types`,
				`    default_type application/octet-stream;`,
				`    # log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '`,
				`    # '$status $body_bytes_sent "$http_referer" '`,
				`    # '"$http_user_agent" "$http_x_forwarded_for"';`,
				`    # access_log logs/access.log  main;`,
				`    sendfile on;`,
				`    # tcp_nopush on;`,
				`    # keepalive_timeout 0;`,
				`    keepalive_timeout 65;`,
				`    # gzip on;`,
				`    # include <== ./conf.d/test*.com.conf`,
				`    server {`,
				`        listen 80;`,
				`        server_name test1.com;`,
				`        # include <== ./conf.d/location.conf`,
				`        location /test1-location {`,
				`            if ($http_api_name != '') {`,
				`                proxy_pass http://wrong_proxy;`,
				`                break;`,
				`            }`,
				`            proxy_pass http://right_proxy;` + "    # test inline comments",
				`        }`,
				`        location /test-to-baidu {`,
				`            proxy_pass https://www.baidu.com;`,
				`        }`,
				`    }`,
				`    server {`,
				`        listen 8080;`,
				`        server_name test1.com;`,
				`        # include <== ./conf.d/location*.conf`,
				`        location /test1-location {`,
				`            if ($http_api_name != '') {`,
				`                proxy_pass http://wrong_proxy;`,
				`                break;`,
				`            }`,
				`            proxy_pass http://right_proxy;` + "    # test inline comments",
				`        }`,
				`        location /test-to-baidu {`,
				`            proxy_pass https://www.baidu.com;`,
				`        }`,
				`        location /test2 {`,
				`            proxy_pass http://test2.test.com;`,
				`        }`,
				`        # # include <== ./conf.d/location.conf`,
				`    }`,
				`    server {`,
				`        listen 8080;`,
				`        server_name test2.com;`,
				`        # include <== ./conf.d/location*.conf`,
				`        location /test1-location {`,
				`            if ($http_api_name != '') {`,
				`                proxy_pass http://wrong_proxy;`,
				`                break;`,
				`            }`,
				`            proxy_pass http://right_proxy;` + "    # test inline comments",
				`        }`,
				`        location /test-to-baidu {`,
				`            proxy_pass https://www.baidu.com;`,
				`        }`,
				`        location /test2 {`,
				`            proxy_pass http://test2.test.com;`,
				`        }`,
				`        # # include <== ./conf.d/location.conf`,
				`    }`,
				`    server {`,
				`        listen 80;`,
				`        server_name localhost;`,
				`        # charset koi8-r;`,
				`        # access_log logs/host.access.log  main;`,
				`        location / {`,
				`            root html;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # error_page 404              /404.html;`,
				`        # redirect server error pages to the static page /50x.html`,
				`        #`,
				`        error_page 500 502 503 504  /50x.html;`,
				`        location = /50x.html {`,
				`            root html;`,
				`        }`,
				`        # proxy the PHP scripts to Apache listening on 127.0.0.1:80`,
				`        #`,
				`        # location ~ \.php$ {`,
				`        #     proxy_pass http://127.0.0.1;`,
				`        # }`,
				`        # pass the PHP scripts to FastCGI server listening on 127.0.0.1:9000`,
				`        #`,
				`        # location ~ \.php$ {`,
				`        #     root html;`,
				`        #     fastcgi_pass 127.0.0.1:9000;`,
				`        #     fastcgi_index index.php;`,
				`        #     fastcgi_param SCRIPT_FILENAME  /scripts$fastcgi_script_name;`,
				`        #     # include <== fastcgi_params`,
				`        # }`,
				`        # deny access to .htaccess files, if Apache's document root`,
				`        # concurs with nginx's one`,
				`        #`,
				`        # location ~ /\.ht {`,
				`        #     deny all;`,
				`        # }`,
				`    }`,
				`    # another virtual host using mix of IP-, name-, and port-based configuration`,
				`    #`,
				`    # server {`,
				`    #     listen 8000;`,
				`    #     listen somename:8080;`,
				`    #     server_name somename  alias  another.alias;`,
				`    #     location / {`,
				`    #         root html;`,
				`    #         index index.html index.htm;`,
				`    #     }`,
				`    # }`,
				`    # HTTPS server`,
				`    #`,
				`    # server {`,
				`    #     listen 443 ssl;`,
				`    #     server_name localhost;`,
				`    #     ssl_certificate cert.pem;`,
				`    #     ssl_certificate_key cert.key;`,
				`    #     ssl_session_cache shared:SSL:1m;`,
				`    #     ssl_session_timeout 5m;`,
				`    #     ssl_ciphers HIGH:!aNULL:!MD5;`,
				`    #     ssl_prefer_server_ciphers on;`,
				`    #     location / {`,
				`    #         root html;`,
				`    #         index index.html index.htm;`,
				`    #     }`,
				`    # }`,
				`}`,
			},
			wantOriginalFingerprints: utilsV3.ConfigFingerprints{
				"C:\\config_test\\conf.d\\location.conf":  "190dbe607c04b9050b88097b1f4df3ecb4081a45dc7ccf115382a5d471fa8fe7",
				"C:\\config_test\\conf.d\\location1.conf": "bc018a68698d57e136fbc54002341d2f173d4803bea9974f89acc376e155ba86",
				"C:\\config_test\\conf.d\\test1.com.conf": "bcb0c82680cceef957db7165d5138f733f2864738557cfc286c635c091457ca3",
				"C:\\config_test\\conf.d\\test2.com.conf": "24480b9ef0c9c86cb90896d7871e10bab94cc25fda496050d48533d6bf542f53",
				"C:\\config_test\\nginx.conf":             "06b19f6a6452402f21b2a28e90205626d18e97b33a743c2846d86431ca236919",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigService{
				store: tt.fields.store,
			}
			got, err := w.GetConfig(tt.args.ctx, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotLines, _ := got.Config.ConfigLines(false); !reflect.DeepEqual(gotLines, tt.wantConfigTextLines) {
				t.Errorf("GetConfig() got.NginxConfig.TextLines = %v, want %v", strings.Join(gotLines, "\n"), strings.Join(tt.wantConfigTextLines, "\n"))
			}
			if gotFP := got.OriginalFingerprints; !reflect.DeepEqual(gotFP, tt.wantOriginalFingerprints) {
				t.Errorf("GetConfig() got.OriginalFingerprints = %v, want %v", gotFP, tt.wantOriginalFingerprints)
			}
		})
	}
}

func Test_webServerConfigService_GetContext(t *testing.T) {
	global.GVA_LOG = log.ZapLogger()
	webSrvOpts := metav1.WebServerOptions{
		GroupID:    1,
		HostID:     1,
		ServerName: "test-bifrost",
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := storev1.NewMockFactory(ctrl)
	wscstore := storev1.NewMockWebServerConfigStore(ctrl)
	store.EXPECT().WebServerConfigs().AnyTimes().Return(wscstore)

	// Use fake data to get config and fingerprints
	fakeSvc := fake.WebServerConfigService{}
	fakeConfig, fakeFP, _ := fakeSvc.Get("test-bifrost")
	ofp := fakeFP.Fingerprints()
	wscstore.EXPECT().GetConfig(nil, webSrvOpts).AnyTimes().Return(fakeConfig, fakeFP, nil)

	// Create diff fingerprint for error test case
	difffp := utilsV3.ConfigFingerprints{
		"C:\\config_test\\nginx.conf": "1111111111111111111111111111111111111111111111111111111111111111",
	}

	// Mock GetContext calls
	wscstore.EXPECT().GetContext(nil, webSrvOpts, ofp, metav1.ConfigContextPos{
		Config:         "C:\\config_test\\conf.d\\location.conf",
		ContextPosPath: []int{0},
	}).AnyTimes().DoAndReturn(func(ctx context.Context, opts metav1.WebServerOptions, fp utilsV3.ConfigFingerprints, pos metav1.ConfigContextPos) (nginx_context.Context, error) {
		target, _ := fakeConfig.Main().GetConfig("C:\\config_test\\conf.d\\location.conf")
		return target.Child(0), nil
	})
	wscstore.EXPECT().GetContext(nil, webSrvOpts, ofp, metav1.ConfigContextPos{
		Config:         "C:\\config_test\\conf.d\\location.conf",
		ContextPosPath: []int{0, 1},
	}).AnyTimes().DoAndReturn(func(ctx context.Context, opts metav1.WebServerOptions, fp utilsV3.ConfigFingerprints, pos metav1.ConfigContextPos) (nginx_context.Context, error) {
		target, _ := fakeConfig.Main().GetConfig("C:\\config_test\\conf.d\\location.conf")
		return target.Child(0).Child(1), nil
	})
	wscstore.EXPECT().GetContext(nil, webSrvOpts, ofp, metav1.ConfigContextPos{
		Config:         "C:\\config_test\\nginx.conf",
		ContextPosPath: []int{},
	}).AnyTimes().DoAndReturn(func(ctx context.Context, opts metav1.WebServerOptions, fp utilsV3.ConfigFingerprints, pos metav1.ConfigContextPos) (nginx_context.Context, error) {
		return fakeConfig.Main(), nil
	})
	wscstore.EXPECT().GetContext(nil, webSrvOpts, ofp, metav1.ConfigContextPos{
		Config:         "C:\\config_test\\nginx.conf",
		ContextPosPath: nil,
	}).AnyTimes().DoAndReturn(func(ctx context.Context, opts metav1.WebServerOptions, fp utilsV3.ConfigFingerprints, pos metav1.ConfigContextPos) (nginx_context.Context, error) {
		return fakeConfig.Main(), nil
	})
	wscstore.EXPECT().GetContext(nil, webSrvOpts, ofp, metav1.ConfigContextPos{
		Config:         "C:\\config_test\\nginx.conf",
		ContextPosPath: []int{1, 2, 3, 4, 5, 6, 7},
	}).AnyTimes().Return(nginx_context.NullContext(), errors.New("invalid context position"))
	wscstore.EXPECT().GetContext(nil, webSrvOpts, difffp, metav1.ConfigContextPos{
		Config:         "C:\\config_test\\conf.d\\location.conf",
		ContextPosPath: []int{0, 1},
	}).AnyTimes().Return(nginx_context.NullContext(), errors.Wrapf(metav1.ErrInconsistentFingerprints, "fingerprint mismatch"))
	type fields struct {
		store storev1.Factory
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
		wantConfigLines         []string
		wantErr                 bool
		wantErrIsInconsistentFP bool
	}{
		{
			name:   "one level pos path",
			fields: fields{store: store},
			args: args{
				opts: webSrvOpts,
				ofp:  ofp,
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
			fields: fields{store: store},
			args: args{
				opts: webSrvOpts,
				ofp:  ofp,
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
			fields: fields{store: store},
			args: args{
				opts: webSrvOpts,
				ofp:  ofp,
				pos: metav1.ConfigContextPos{
					Config:         "C:\\config_test\\nginx.conf",
					ContextPosPath: []int{},
				},
			},
			wantConfigLines: []string{
				`# testConfigManager`,
				`# testConfigManager`,
				`# testConfigManager_update...`,
				`# testConfigManager_update...`,
				`# user nobody;`,
				`worker_processes 1;    # test inline comments`,
				`# error_log logs/error.log;`,
				`# error_log logs/error.log  notice;`,
				`# error_log logs/error.log  info;`,
				`# pid logs/nginx.pid;`,
				`events {`,
				`    worker_connections 1024;`,
				`}`,
				`stream {`,
				`    server {`,
				`        listen 8888;`,
				`        proxy_pass baidu.com:22;`,
				`    }`,
				`}`,
				`http {`,
				`    # include <== mime.types`,
				`    default_type application/octet-stream;`,
				`    # log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '`,
				`    # '$status $body_bytes_sent "$http_referer" '`,
				`    # '"$http_user_agent" "$http_x_forwarded_for"';`,
				`    # access_log logs/access.log  main;`,
				`    sendfile on;`,
				`    # tcp_nopush on;`,
				`    # keepalive_timeout 0;`,
				`    keepalive_timeout 65;`,
				`    # gzip on;`,
				`    # include <== ./conf.d/test*.com.conf`,
				`    server {`,
				`        listen 8080;`,
				`        server_name test1.com;`,
				`        # include <== ./conf.d/location*.conf`,
				`        location /test1-location {`,
				`            if ($http_api_name != '') {`,
				`                proxy_pass http://wrong_proxy;`,
				`                break;`,
				`            }`,
				`            proxy_pass http://right_proxy;    # test inline comments`,
				`        }`,
				`        location /test-to-baidu {`,
				`            proxy_pass https://www.baidu.com;`,
				`        }`,
				`        location /test2 {`,
				`            proxy_pass http://test2.test.com;`,
				`        }`,
				`        # # include <== ./conf.d/location.conf`,
				`    }`,
				`    server {`,
				`        listen 8080;`,
				`        server_name test2.com;`,
				`        # include <== ./conf.d/location*.conf`,
				`        location /test1-location {`,
				`            if ($http_api_name != '') {`,
				`                proxy_pass http://wrong_proxy;`,
				`                break;`,
				`            }`,
				`            proxy_pass http://right_proxy;    # test inline comments`,
				`        }`,
				`        location /test-to-baidu {`,
				`            proxy_pass https://www.baidu.com;`,
				`        }`,
				`        location /test2 {`,
				`            proxy_pass http://test2.test.com;`,
				`        }`,
				`        # # include <== ./conf.d/location.conf`,
				`    }`,
				`    server {`,
				`        listen 80;`,
				`        server_name test1.com;`,
				`        # include <== ./conf.d/location.conf`,
				`        location /test1-location {`,
				`            if ($http_api_name != '') {`,
				`                proxy_pass http://wrong_proxy;`,
				`                break;`,
				`            }`,
				`            proxy_pass http://right_proxy;    # test inline comments`,
				`        }`,
				`        location /test-to-baidu {`,
				`            proxy_pass https://www.baidu.com;`,
				`        }`,
				`    }`,
				`    server {`,
				`        listen 80;`,
				`        server_name localhost;`,
				`        # charset koi8-r;`,
				`        # access_log logs/host.access.log  main;`,
				`        location / {`,
				`            root html;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # error_page 404              /404.html;`,
				`        # redirect server error pages to the static page /50x.html`,
				`        #`,
				`        error_page 500 502 503 504  /50x.html;`,
				`        location = /50x.html {`,
				`            root html;`,
				`        }`,
				`        # proxy the PHP scripts to Apache listening on 127.0.0.1:80`,
				`        #`,
				`        # location ~ \.php$ {`,
				`        #     proxy_pass http://127.0.0.1;`,
				`        # }`,
				`        # pass the PHP scripts to FastCGI server listening on 127.0.0.1:9000`,
				`        #`,
				`        # location ~ \.php$ {`,
				`        #     root html;`,
				`        #     fastcgi_pass 127.0.0.1:9000;`,
				`        #     fastcgi_index index.php;`,
				`        #     fastcgi_param SCRIPT_FILENAME  /scripts$fastcgi_script_name;`,
				`        #     # include <== fastcgi_params`,
				`        # }`,
				`        # deny access to .htaccess files, if Apache's document root`,
				`        # concurs with nginx's one`,
				`        #`,
				`        # location ~ /\.ht {`,
				`        #     deny all;`,
				`        # }`,
				`    }`,
				`    # another virtual host using mix of IP-, name-, and port-based configuration`,
				`    #`,
				`    # server {`,
				`    #     listen 8000;`,
				`    #     listen somename:8080;`,
				`    #     server_name somename  alias  another.alias;`,
				`    #     location / {`,
				`    #         root html;`,
				`    #         index index.html index.htm;`,
				`    #     }`,
				`    # }`,
				`    # HTTPS server`,
				`    #`,
				`    # server {`,
				`    #     listen 443 ssl;`,
				`    #     server_name localhost;`,
				`    #     ssl_certificate cert.pem;`,
				`    #     ssl_certificate_key cert.key;`,
				`    #     ssl_session_cache shared:SSL:1m;`,
				`    #     ssl_session_timeout 5m;`,
				`    #     ssl_ciphers HIGH:!aNULL:!MD5;`,
				`    #     ssl_prefer_server_ciphers on;`,
				`    #     location / {`,
				`    #         root html;`,
				`    #         index index.html index.htm;`,
				`    #     }`,
				`    # }`,
				`}`,
			},
		},
		{
			name:   "nil pos path",
			fields: fields{store: store},
			args: args{
				opts: webSrvOpts,
				ofp:  ofp,
				pos: metav1.ConfigContextPos{
					Config:         "C:\\config_test\\nginx.conf",
					ContextPosPath: nil,
				},
			},
			wantConfigLines: []string{
				`# testConfigManager`,
				`# testConfigManager`,
				`# testConfigManager_update...`,
				`# testConfigManager_update...`,
				`# user nobody;`,
				`worker_processes 1;    # test inline comments`,
				`# error_log logs/error.log;`,
				`# error_log logs/error.log  notice;`,
				`# error_log logs/error.log  info;`,
				`# pid logs/nginx.pid;`,
				`events {`,
				`    worker_connections 1024;`,
				`}`,
				`stream {`,
				`    server {`,
				`        listen 8888;`,
				`        proxy_pass baidu.com:22;`,
				`    }`,
				`}`,
				`http {`,
				`    # include <== mime.types`,
				`    default_type application/octet-stream;`,
				`    # log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '`,
				`    # '$status $body_bytes_sent "$http_referer" '`,
				`    # '"$http_user_agent" "$http_x_forwarded_for"';`,
				`    # access_log logs/access.log  main;`,
				`    sendfile on;`,
				`    # tcp_nopush on;`,
				`    # keepalive_timeout 0;`,
				`    keepalive_timeout 65;`,
				`    # gzip on;`,
				`    # include <== ./conf.d/test*.com.conf`,
				`    server {`,
				`        listen 8080;`,
				`        server_name test1.com;`,
				`        # include <== ./conf.d/location*.conf`,
				`        location /test1-location {`,
				`            if ($http_api_name != '') {`,
				`                proxy_pass http://wrong_proxy;`,
				`                break;`,
				`            }`,
				`            proxy_pass http://right_proxy;    # test inline comments`,
				`        }`,
				`        location /test-to-baidu {`,
				`            proxy_pass https://www.baidu.com;`,
				`        }`,
				`        location /test2 {`,
				`            proxy_pass http://test2.test.com;`,
				`        }`,
				`        # # include <== ./conf.d/location.conf`,
				`    }`,
				`    server {`,
				`        listen 8080;`,
				`        server_name test2.com;`,
				`        # include <== ./conf.d/location*.conf`,
				`        location /test1-location {`,
				`            if ($http_api_name != '') {`,
				`                proxy_pass http://wrong_proxy;`,
				`                break;`,
				`            }`,
				`            proxy_pass http://right_proxy;    # test inline comments`,
				`        }`,
				`        location /test-to-baidu {`,
				`            proxy_pass https://www.baidu.com;`,
				`        }`,
				`        location /test2 {`,
				`            proxy_pass http://test2.test.com;`,
				`        }`,
				`        # # include <== ./conf.d/location.conf`,
				`    }`,
				`    server {`,
				`        listen 80;`,
				`        server_name test1.com;`,
				`        # include <== ./conf.d/location.conf`,
				`        location /test1-location {`,
				`            if ($http_api_name != '') {`,
				`                proxy_pass http://wrong_proxy;`,
				`                break;`,
				`            }`,
				`            proxy_pass http://right_proxy;    # test inline comments`,
				`        }`,
				`        location /test-to-baidu {`,
				`            proxy_pass https://www.baidu.com;`,
				`        }`,
				`    }`,
				`    server {`,
				`        listen 80;`,
				`        server_name localhost;`,
				`        # charset koi8-r;`,
				`        # access_log logs/host.access.log  main;`,
				`        location / {`,
				`            root html;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # error_page 404              /404.html;`,
				`        # redirect server error pages to the static page /50x.html`,
				`        #`,
				`        error_page 500 502 503 504  /50x.html;`,
				`        location = /50x.html {`,
				`            root html;`,
				`        }`,
				`        # proxy the PHP scripts to Apache listening on 127.0.0.1:80`,
				`        #`,
				`        # location ~ \.php$ {`,
				`        #     proxy_pass http://127.0.0.1;`,
				`        # }`,
				`        # pass the PHP scripts to FastCGI server listening on 127.0.0.1:9000`,
				`        #`,
				`        # location ~ \.php$ {`,
				`        #     root html;`,
				`        #     fastcgi_pass 127.0.0.1:9000;`,
				`        #     fastcgi_index index.php;`,
				`        #     fastcgi_param SCRIPT_FILENAME  /scripts$fastcgi_script_name;`,
				`        #     # include <== fastcgi_params`,
				`        # }`,
				`        # deny access to .htaccess files, if Apache's document root`,
				`        # concurs with nginx's one`,
				`        #`,
				`        # location ~ /\.ht {`,
				`        #     deny all;`,
				`        # }`,
				`    }`,
				`    # another virtual host using mix of IP-, name-, and port-based configuration`,
				`    #`,
				`    # server {`,
				`    #     listen 8000;`,
				`    #     listen somename:8080;`,
				`    #     server_name somename  alias  another.alias;`,
				`    #     location / {`,
				`    #         root html;`,
				`    #         index index.html index.htm;`,
				`    #     }`,
				`    # }`,
				`    # HTTPS server`,
				`    #`,
				`    # server {`,
				`    #     listen 443 ssl;`,
				`    #     server_name localhost;`,
				`    #     ssl_certificate cert.pem;`,
				`    #     ssl_certificate_key cert.key;`,
				`    #     ssl_session_cache shared:SSL:1m;`,
				`    #     ssl_session_timeout 5m;`,
				`    #     ssl_ciphers HIGH:!aNULL:!MD5;`,
				`    #     ssl_prefer_server_ciphers on;`,
				`    #     location / {`,
				`    #         root html;`,
				`    #         index index.html index.htm;`,
				`    #     }`,
				`    # }`,
				`}`,
			},
		},
		{
			name:   "wrong pos path",
			fields: fields{store: store},
			args: args{
				opts: webSrvOpts,
				ofp:  ofp,
				pos: metav1.ConfigContextPos{
					Config:         "C:\\config_test\\nginx.conf",
					ContextPosPath: []int{1, 2, 3, 4, 5, 6, 7},
				},
			},
			wantErr: true,
		},
		{
			name:   "inconsistent fingerprints",
			fields: fields{store: store},
			args: args{
				opts: webSrvOpts,
				ofp:  difffp,
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
			w := &webServerConfigService{
				store: tt.fields.store,
			}
			got, err := w.GetContext(tt.args.ctx, tt.args.opts, tt.args.ofp, tt.args.pos)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err != nil && (errors.Is(err, metav1.ErrInconsistentFingerprints) || errors.IsCode(err, 110010)) != tt.wantErrIsInconsistentFP {
				t.Errorf("GetContext() error = %v, wantErrIsInconsistentFP %v", err, tt.wantErrIsInconsistentFP)
				return
			}
			gotConfigLines, _ := got.ConfigLines(false)
			if !reflect.DeepEqual(gotConfigLines, tt.wantConfigLines) {
				//for _, line := range gotConfigLines {
				//	fmt.Printf("`%s`,\n", line)
				//}
				// Find and print first difference
				minLen := len(gotConfigLines)
				if len(tt.wantConfigLines) < minLen {
					minLen = len(tt.wantConfigLines)
				}
				for i := 0; i < minLen; i++ {
					if gotConfigLines[i] != tt.wantConfigLines[i] {
						t.Errorf("GetContext() first diff at line %d:\n  got: %q\n want: %q", i+1, gotConfigLines[i], tt.wantConfigLines[i])
						break
					}
				}
				if len(gotConfigLines) != len(tt.wantConfigLines) {
					t.Errorf("GetContext() line count mismatch: got %d, want %d", len(gotConfigLines), len(tt.wantConfigLines))
				}
			}
		})
	}
}

func Test_webServerConfigService_GetIncludedConfigs(t *testing.T) {
	global.GVA_LOG = log.ZapLogger()
	webSrvOpts := metav1.WebServerOptions{
		GroupID:    1,
		HostID:     1,
		ServerName: "test-bifrost",
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := storev1.NewMockFactory(ctrl)
	wscstore := storev1.NewMockWebServerConfigStore(ctrl)
	store.EXPECT().WebServerConfigs().AnyTimes().Return(wscstore)

	// Use fake data to get config and fingerprints
	fakeSvc := fake.WebServerConfigService{}
	fakeConfig, fakeFP, _ := fakeSvc.Get("test-bifrost")
	ofp := fakeFP.Fingerprints()
	wscstore.EXPECT().GetConfig(nil, webSrvOpts).AnyTimes().Return(fakeConfig, fakeFP, nil)

	// Create diff fingerprint for error test case
	difffp := utilsV3.ConfigFingerprints{
		"C:\\config_test\\nginx.conf": "1111111111111111111111111111111111111111111111111111111111111111",
	}

	// Mock GetContext calls for GetIncludedConfigs tests
	normalPos := metav1.ConfigContextPos{
		Config:         "C:\\config_test\\nginx.conf",
		ContextPosPath: []int{1, 0},
	}
	wrongPos := metav1.ConfigContextPos{
		Config:         "C:\\config_test\\nginx.con",
		ContextPosPath: nil,
	}
	notIncludePos := metav1.ConfigContextPos{
		Config:         "C:\\config_test\\nginx.conf",
		ContextPosPath: []int{0, 1},
	}

	// Create a mock Include context for normal test
	mockIncludeCtx := nginx_context.NullContext()
	wscstore.EXPECT().GetContext(nil, webSrvOpts, ofp, normalPos).AnyTimes().Return(mockIncludeCtx, nil)
	wscstore.EXPECT().GetContext(nil, webSrvOpts, ofp, wrongPos).AnyTimes().Return(nginx_context.NullContext(), errors.New("wrong position"))
	wscstore.EXPECT().GetContext(nil, webSrvOpts, ofp, notIncludePos).AnyTimes().Return(nginx_context.NullContext(), nil)
	wscstore.EXPECT().GetContext(nil, webSrvOpts, difffp, normalPos).AnyTimes().Return(nginx_context.NullContext(), errors.Wrapf(metav1.ErrInconsistentFingerprints, "fingerprint mismatch"))

	// Mock GetIncludedConfigs calls
	wscstore.EXPECT().GetIncludedConfigs(nil, webSrvOpts, ofp, normalPos).AnyTimes().Return([]string{
		"C:\\config_test\\conf.d\\location.conf",
		"C:\\config_test\\conf.d\\location2.conf",
	}, nil)
	wscstore.EXPECT().GetIncludedConfigs(nil, webSrvOpts, ofp, wrongPos).AnyTimes().Return(nil, errors.New("wrong position"))
	wscstore.EXPECT().GetIncludedConfigs(nil, webSrvOpts, ofp, notIncludePos).AnyTimes().Return(nil, errors.New("the target is not an `include` context"))
	wscstore.EXPECT().GetIncludedConfigs(nil, webSrvOpts, difffp, normalPos).AnyTimes().Return(nil, errors.Wrapf(metav1.ErrInconsistentFingerprints, "fingerprint mismatch"))
	type fields struct {
		store storev1.Factory
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
			name:   "normal test",
			fields: fields{store: store},
			args: args{
				opts: webSrvOpts,
				ofp:  ofp,
				pos:  normalPos,
			},
			want: []string{
				"C:\\config_test\\conf.d\\location.conf",
				"C:\\config_test\\conf.d\\location2.conf",
			},
			wantErr: false,
		},
		{
			name:   "wrong position",
			fields: fields{store: store},
			args: args{
				opts: webSrvOpts,
				ofp:  ofp,
				pos:  wrongPos,
			},
			wantErr: true,
		},
		{
			name:   "the target is not an `include` context",
			fields: fields{store: store},
			args: args{
				opts: webSrvOpts,
				ofp:  ofp,
				pos:  notIncludePos,
			},
			wantErr: true,
		},
		{
			name:   "inconsistent fingerprints",
			fields: fields{store: store},
			args: args{
				opts: webSrvOpts,
				ofp:  difffp,
				pos:  normalPos,
			},
			wantErr:                 true,
			wantErrIsInconsistentFP: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigService{
				store: tt.fields.store,
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

func Test_webServerConfigService_GetOptions(t *testing.T) {
	type fields struct {
		store storev1.Factory
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigService{
				store: tt.fields.store,
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

func Test_webServerConfigService_InsertWithClone(t *testing.T) {
	webSrvOpts := metav1.WebServerOptions{
		GroupID:    1,
		HostID:     1,
		ServerName: "test-bifrost",
	}
	ctxmeta := metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]{
		Position: metav1.ConfigContextPos{
			Config:         "C:\\config_test\\nginx.conf",
			ContextPosPath: []int{8, 13, 4},
		},
		TargetContext: metav1.CloneConfigContextMeta{ConfigContextPos: metav1.ConfigContextPos{
			Config:         "C:\\config_test\\conf.d\\location2.conf",
			ContextPosPath: []int{4},
		}},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := storev1.NewMockFactory(ctrl)
	cacheStore := cache.GetCacheStore(store, time.Second)
	wscstore := storev1.NewMockWebServerConfigStore(ctrl)
	store.EXPECT().WebServerConfigs().AnyTimes().Return(wscstore)

	// Use fake data to get config and fingerprints
	fakeSvc := fake.WebServerConfigService{}
	fakeConfig, fakeFP, _ := fakeSvc.Get("test-bifrost")
	wscstore.EXPECT().GetConfig(nil, webSrvOpts).AnyTimes().Return(fakeConfig, fakeFP, nil)

	_, ofp, err := cacheStore.WebServerConfigs().GetConfig(nil, webSrvOpts)
	if err != nil {
		t.Fatal(err)
	}
	disableTheTarget := true
	wscstore.EXPECT().InsertWithClone(nil, webSrvOpts, ofp.Fingerprints(), ctxmeta, disableTheTarget).AnyTimes().Return(nil)
	type fields struct {
		store storev1.Factory
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
			name:   "normal test",
			fields: fields{store: cacheStore},
			args: args{
				opts:           webSrvOpts,
				ofp:            ofp.Fingerprints(),
				ctxmeta:        ctxmeta,
				disabledTarget: disableTheTarget,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigService{
				store: tt.fields.store,
			}
			if err := w.InsertWithClone(tt.args.ctx, tt.args.opts, tt.args.ofp, tt.args.ctxmeta, tt.args.disabledTarget); (err != nil) != tt.wantErr {
				t.Errorf("InsertWithClone() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigService_InsertWithNew(t *testing.T) {
	webSrvOpts := metav1.WebServerOptions{
		GroupID:    1,
		HostID:     1,
		ServerName: "test-bifrost",
	}
	ctxmeta := metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]{
		Position: metav1.ConfigContextPos{
			Config:         "C:\\config_test\\conf.d\\location2.conf",
			ContextPosPath: []int{2},
		},
		TargetContext: metav1.NewConfigContextMeta{
			ContextType:  "location",
			ContextValue: "~ /normal-test",
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := storev1.NewMockFactory(ctrl)
	cacheStore := cache.GetCacheStore(store, time.Second)
	wscstore := storev1.NewMockWebServerConfigStore(ctrl)
	store.EXPECT().WebServerConfigs().AnyTimes().Return(wscstore)

	// Use fake data to get config and fingerprints
	fakeSvc := fake.WebServerConfigService{}
	fakeConfig, fakeFP, _ := fakeSvc.Get("test-bifrost")
	wscstore.EXPECT().GetConfig(nil, webSrvOpts).AnyTimes().Return(fakeConfig, fakeFP, nil)

	_, ofp, err := cacheStore.WebServerConfigs().GetConfig(nil, webSrvOpts)
	if err != nil {
		t.Fatal(err)
	}
	disableTheTarget := true
	wscstore.EXPECT().InsertWithNew(nil, webSrvOpts, ofp.Fingerprints(), ctxmeta, disableTheTarget).AnyTimes().Return(nil)
	type fields struct {
		store storev1.Factory
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
			name:   "normal test",
			fields: fields{store: cacheStore},
			args: args{
				opts:           webSrvOpts,
				ofp:            ofp.Fingerprints(),
				ctxmeta:        ctxmeta,
				disabledTarget: disableTheTarget,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigService{
				store: tt.fields.store,
			}
			if err := w.InsertWithNew(tt.args.ctx, tt.args.opts, tt.args.ofp, tt.args.ctxmeta, tt.args.disabledTarget); (err != nil) != tt.wantErr {
				t.Errorf("InsertWithNew() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigService_ModifyContextValue(t *testing.T) {
	webSrvOpts := metav1.WebServerOptions{
		GroupID:    1,
		HostID:     1,
		ServerName: "test-bifrost",
	}
	ctxmeta := metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]{
		Position: metav1.ConfigContextPos{
			Config:         "C:\\config_test\\conf.d\\location2.conf",
			ContextPosPath: []int{0},
		},
		TargetContext: metav1.NewConfigContextMeta{
			ContextType:  "location",
			ContextValue: "~ /normal-test",
		},
	}
	unmatchedTypeCtxmeta := metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]{
		Position: metav1.ConfigContextPos{
			Config:         "C:\\config_test\\conf.d\\location2.conf",
			ContextPosPath: []int{2},
		},
		TargetContext: metav1.NewConfigContextMeta{
			ContextType:  "location",
			ContextValue: "~ /normal-test",
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := storev1.NewMockFactory(ctrl)
	cacheStore := cache.GetCacheStore(store, time.Second)
	wscstore := storev1.NewMockWebServerConfigStore(ctrl)
	store.EXPECT().WebServerConfigs().AnyTimes().Return(wscstore)

	// Use fake data to get config and fingerprints
	fakeSvc := fake.WebServerConfigService{}
	fakeConfig, fakeFP, _ := fakeSvc.Get("test-bifrost")
	wscstore.EXPECT().GetConfig(nil, webSrvOpts).AnyTimes().Return(fakeConfig, fakeFP, nil)

	_, ofp, err := cacheStore.WebServerConfigs().GetConfig(nil, webSrvOpts)
	if err != nil {
		t.Fatal(err)
	}
	wscstore.EXPECT().ModifyContextValue(nil, webSrvOpts, ofp.Fingerprints(), ctxmeta).AnyTimes().Return(nil)
	wscstore.EXPECT().ModifyContextValue(nil, webSrvOpts, ofp.Fingerprints(), unmatchedTypeCtxmeta).AnyTimes().Return(errors.New("unmatched context type"))
	type fields struct {
		store storev1.Factory
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
			fields: fields{store: cacheStore},
			args: args{
				opts:    webSrvOpts,
				ofp:     ofp.Fingerprints(),
				ctxmeta: ctxmeta,
			},
			wantErr: false,
		},
		{
			name:   "unmatched context type",
			fields: fields{store: cacheStore},
			args: args{
				opts:    webSrvOpts,
				ofp:     ofp.Fingerprints(),
				ctxmeta: unmatchedTypeCtxmeta,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigService{
				store: tt.fields.store,
			}
			if err := w.ModifyContextValue(tt.args.ctx, tt.args.opts, tt.args.ofp, tt.args.ctxmeta); (err != nil) != tt.wantErr {
				t.Errorf("ModifyContextValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigService_ModifyWithClone(t *testing.T) {
	webSrvOpts := metav1.WebServerOptions{
		GroupID:    1,
		HostID:     1,
		ServerName: "test-bifrost",
	}
	ctxmeta := metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]{
		Position: metav1.ConfigContextPos{
			Config:         "C:\\config_test\\nginx.conf",
			ContextPosPath: []int{8, 13, 4},
		},
		TargetContext: metav1.CloneConfigContextMeta{ConfigContextPos: metav1.ConfigContextPos{
			Config:         "C:\\config_test\\conf.d\\location2.conf",
			ContextPosPath: []int{4},
		}},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := storev1.NewMockFactory(ctrl)
	cacheStore := cache.GetCacheStore(store, time.Second)
	wscstore := storev1.NewMockWebServerConfigStore(ctrl)
	store.EXPECT().WebServerConfigs().AnyTimes().Return(wscstore)

	// Use fake data to get config and fingerprints
	fakeSvc := fake.WebServerConfigService{}
	fakeConfig, fakeFP, _ := fakeSvc.Get("test-bifrost")
	wscstore.EXPECT().GetConfig(nil, webSrvOpts).AnyTimes().Return(fakeConfig, fakeFP, nil)

	_, ofp, err := cacheStore.WebServerConfigs().GetConfig(nil, webSrvOpts)
	if err != nil {
		t.Fatal(err)
	}
	wscstore.EXPECT().ModifyWithClone(nil, webSrvOpts, ofp.Fingerprints(), ctxmeta).AnyTimes().Return(nil)
	type fields struct {
		store storev1.Factory
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
			name:   "normal test",
			fields: fields{store: cacheStore},
			args: args{
				opts:    webSrvOpts,
				ofp:     ofp.Fingerprints(),
				ctxmeta: ctxmeta,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigService{
				store: tt.fields.store,
			}
			if err := w.ModifyWithClone(tt.args.ctx, tt.args.opts, tt.args.ofp, tt.args.ctxmeta); (err != nil) != tt.wantErr {
				t.Errorf("ModifyWithClone() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigService_ModifyWithNew(t *testing.T) {
	webSrvOpts := metav1.WebServerOptions{
		GroupID:    1,
		HostID:     1,
		ServerName: "test-bifrost",
	}
	ctxmeta := metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]{
		Position: metav1.ConfigContextPos{
			Config:         "C:\\config_test\\conf.d\\location2.conf",
			ContextPosPath: []int{2},
		},
		TargetContext: metav1.NewConfigContextMeta{
			ContextType:  "location",
			ContextValue: "~ /normal-test",
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := storev1.NewMockFactory(ctrl)
	cacheStore := cache.GetCacheStore(store, time.Second)
	wscstore := storev1.NewMockWebServerConfigStore(ctrl)
	store.EXPECT().WebServerConfigs().AnyTimes().Return(wscstore)

	// Use fake data to get config and fingerprints
	fakeSvc := fake.WebServerConfigService{}
	fakeConfig, fakeFP, _ := fakeSvc.Get("test-bifrost")
	wscstore.EXPECT().GetConfig(nil, webSrvOpts).AnyTimes().Return(fakeConfig, fakeFP, nil)

	_, ofp, err := cacheStore.WebServerConfigs().GetConfig(nil, webSrvOpts)
	if err != nil {
		t.Fatal(err)
	}
	wscstore.EXPECT().ModifyWithNew(nil, webSrvOpts, ofp.Fingerprints(), ctxmeta).AnyTimes().Return(nil)
	type fields struct {
		store storev1.Factory
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
			fields: fields{store: cacheStore},
			args: args{
				opts:    webSrvOpts,
				ofp:     ofp.Fingerprints(),
				ctxmeta: ctxmeta,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigService{
				store: tt.fields.store,
			}
			if err := w.ModifyWithNew(tt.args.ctx, tt.args.opts, tt.args.ofp, tt.args.ctxmeta); (err != nil) != tt.wantErr {
				t.Errorf("ModifyWithNew() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigService_Remove(t *testing.T) {
	webSrvOpts := metav1.WebServerOptions{
		GroupID:    1,
		HostID:     1,
		ServerName: "test-bifrost",
	}
	pos := metav1.ConfigContextPos{
		Config:         "C:\\config_test\\nginx.conf",
		ContextPosPath: []int{8, 13, 4},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := storev1.NewMockFactory(ctrl)
	cacheStore := cache.GetCacheStore(store, time.Second)
	wscstore := storev1.NewMockWebServerConfigStore(ctrl)
	store.EXPECT().WebServerConfigs().AnyTimes().Return(wscstore)

	// Use fake data to get config and fingerprints
	fakeSvc := fake.WebServerConfigService{}
	fakeConfig, fakeFP, _ := fakeSvc.Get("test-bifrost")
	wscstore.EXPECT().GetConfig(nil, webSrvOpts).AnyTimes().Return(fakeConfig, fakeFP, nil)

	_, ofp, err := cacheStore.WebServerConfigs().GetConfig(nil, webSrvOpts)
	if err != nil {
		t.Fatal(err)
	}

	wscstore.EXPECT().Remove(nil, webSrvOpts, ofp.Fingerprints(), pos).AnyTimes().Return(nil)
	type fields struct {
		store storev1.Factory
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
			name:   "normal test",
			fields: fields{store: cacheStore},
			args: args{
				opts: webSrvOpts,
				ofp:  ofp.Fingerprints(),
				pos:  pos,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigService{
				store: tt.fields.store,
			}
			if err := w.Remove(tt.args.ctx, tt.args.opts, tt.args.ofp, tt.args.pos); (err != nil) != tt.wantErr {
				t.Errorf("Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigService_Move(t *testing.T) {
	webSrvOpts := metav1.WebServerOptions{
		GroupID:    1,
		HostID:     1,
		ServerName: "test-bifrost",
	}
	ctxmeta := metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]{
		Position: metav1.ConfigContextPos{
			Config:         "C:\\config_test\\nginx.conf",
			ContextPosPath: []int{8, 13, 4},
		},
		TargetContext: metav1.CloneConfigContextMeta{ConfigContextPos: metav1.ConfigContextPos{
			Config:         "C:\\config_test\\conf.d\\location2.conf",
			ContextPosPath: []int{4},
		}},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := storev1.NewMockFactory(ctrl)
	cacheStore := cache.GetCacheStore(store, time.Second)
	wscstore := storev1.NewMockWebServerConfigStore(ctrl)
	store.EXPECT().WebServerConfigs().AnyTimes().Return(wscstore)

	// Use fake data to get config and fingerprints
	fakeSvc := fake.WebServerConfigService{}
	fakeConfig, fakeFP, _ := fakeSvc.Get("test-bifrost")
	wscstore.EXPECT().GetConfig(nil, webSrvOpts).AnyTimes().Return(fakeConfig, fakeFP, nil)

	_, ofp, err := cacheStore.WebServerConfigs().GetConfig(nil, webSrvOpts)
	if err != nil {
		t.Fatal(err)
	}
	disableTheTarget := true
	wscstore.EXPECT().Move(nil, webSrvOpts, ofp.Fingerprints(), ctxmeta, disableTheTarget).AnyTimes().Return(nil)
	type fields struct {
		store storev1.Factory
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
			name:   "normal test",
			fields: fields{store: cacheStore},
			args: args{
				opts:           webSrvOpts,
				ofp:            ofp.Fingerprints(),
				ctxmeta:        ctxmeta,
				disabledTarget: disableTheTarget,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigService{
				store: tt.fields.store,
			}
			if err := w.Move(tt.args.ctx, tt.args.opts, tt.args.ofp, tt.args.ctxmeta, tt.args.disabledTarget); (err != nil) != tt.wantErr {
				t.Errorf("Move() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigService_SearchContextPositions(t *testing.T) {
	global.GVA_LOG = log.ZapLogger()
	webSrvOpts := metav1.WebServerOptions{
		GroupID:    1,
		HostID:     1,
		ServerName: "test-bifrost",
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := storev1.NewMockFactory(ctrl)
	wscstore := storev1.NewMockWebServerConfigStore(ctrl)
	store.EXPECT().WebServerConfigs().AnyTimes().Return(wscstore)

	// Use fake data to get config and fingerprints
	fakeSvc := fake.WebServerConfigService{}
	_, fakeFP, _ := fakeSvc.Get("test-bifrost")
	ofp := fakeFP.Fingerprints()

	// Create diff fingerprint for error test case
	difffp := utilsV3.ConfigFingerprints{
		"C:\\config_test\\nginx.conf": "1111111111111111111111111111111111111111111111111111111111111111",
	}
	type fields struct {
		store storev1.Factory
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
			fields: fields{store: store},
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
				{"C:\\config_test\\nginx.conf", []int{8, 13, 0}},
				{"C:\\config_test\\nginx.conf", []int{8, 17, 0}},
			},
		},
		{
			name:   "regexp match rule, only in current config",
			fields: fields{store: store},
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
				{"C:\\config_test\\nginx.conf", []int{8, 13, 1}},
				{"C:\\config_test\\nginx.conf", []int{8, 14, 1}},
				{"C:\\config_test\\nginx.conf", []int{8, 20, 1}},
			},
		},
		{
			name:   "string match rule, not only in current config",
			fields: fields{store: store},
			args: args{
				ctx:  nil,
				opts: webSrvOpts,
				fp:   ofp,
				kwmeta: metav1.SearchKeywordsMeta{
					StartingPositionList: []metav1.ConfigContextPos{{Config: "C:\\config_test\\nginx.conf", ContextPosPath: []int{8}}},
					Keywords:             "listen 80",
					IsRegexpRule:         false,
					IsOnlyInCurrent:      false,
				},
			},
			want: []metav1.ConfigContextPos{
				{"C:\\config_test\\nginx.conf", []int{8, 13, 0}},
				{"C:\\config_test\\nginx.conf", []int{8, 17, 0}},
				{"conf.d\\server_test1.conf", []int{0, 0}},
				{"conf.d\\server_test2.conf", []int{0, 0}},
				{"conf.d\\server_test2.conf", []int{1, 0}},
				{"conf.d\\test1.com.conf", []int{0, 1}},
				{"conf.d\\test2.com.conf", []int{0, 0}},
			},
		},
		{
			name:   "regexp match rule, not only in current config",
			fields: fields{store: store},
			args: args{
				ctx:  nil,
				opts: webSrvOpts,
				fp:   ofp,
				kwmeta: metav1.SearchKeywordsMeta{
					StartingPositionList: []metav1.ConfigContextPos{{Config: "C:\\config_test\\nginx.conf", ContextPosPath: []int{8}}},
					Keywords:             `^server_name\s+test.*$`,
					IsRegexpRule:         true,
					IsOnlyInCurrent:      false,
				},
			},
			want: []metav1.ConfigContextPos{
				{"conf.d\\server_test1.conf", []int{0, 1}},
				{"conf.d\\server_test2.conf", []int{0, 1}},
				{"conf.d\\server_test2.conf", []int{1, 1}},
				{"conf.d\\test1.com.conf", []int{0, 2}},
				{"conf.d\\test2.com.conf", []int{0, 1}},
			},
		},
		{
			name:   "context not found",
			fields: fields{store: store},
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
			fields: fields{store: store},
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
			fields: fields{store: store},
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
			if tt.name == "context not found" {
				wscstore.EXPECT().SearchContextPositions(tt.args.ctx, tt.args.opts, tt.args.fp, tt.args.kwmeta).Return([]metav1.ConfigContextPos{}, nil)
			} else if tt.name == "config not found" {
				wscstore.EXPECT().SearchContextPositions(tt.args.ctx, tt.args.opts, tt.args.fp, tt.args.kwmeta).Return(nil, errors.New("config not found"))
			} else if tt.name == "inconsistent fingerprints" {
				wscstore.EXPECT().SearchContextPositions(tt.args.ctx, tt.args.opts, tt.args.fp, tt.args.kwmeta).Return(nil, errors.Wrapf(metav1.ErrInconsistentFingerprints, "fingerprint mismatch"))
			} else {
				wscstore.EXPECT().SearchContextPositions(tt.args.ctx, tt.args.opts, tt.args.fp, tt.args.kwmeta).Return(tt.want, nil)
			}

			w := &webServerConfigService{
				store: tt.fields.store,
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

func Test_webServerConfigService_UpdateConfig(t *testing.T) {
	webSrvOpts := metav1.WebServerOptions{
		GroupID:    1,
		HostID:     1,
		ServerName: "test-bifrost",
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := storev1.NewMockFactory(ctrl)
	wscstore := storev1.NewMockWebServerConfigStore(ctrl)
	store.EXPECT().WebServerConfigs().AnyTimes().Return(wscstore)

	// Use fake data to get config and fingerprints
	fakeSvc := fake.WebServerConfigService{}
	fakeConfig, fakeFP, _ := fakeSvc.Get("test-bifrost")
	wscstore.EXPECT().GetConfig(nil, webSrvOpts).AnyTimes().Return(fakeConfig, fakeFP, nil)

	cacheStore := cache.GetCacheStore(store, time.Second)
	_, ofp, err := cacheStore.WebServerConfigs().GetConfig(nil, webSrvOpts)
	if err != nil {
		t.Fatal(err)
	}

	validJsonData := []byte(`{"main-config":"C:\\config_test\\nginx.conf","configs":{"C:\\config_test\\nginx.conf":{"enabled":true,"context-type":"config","value":"C:\\config_test\\nginx.conf","params":[{"context-type":"directive","value":"user nobody"},{"enabled":true,"context-type":"directive","value":"worker_processes 1"}]}}}`)
	invalidJsonData := []byte(`invalid json data`)

	type fields struct {
		store storev1.Factory
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
			name:   "valid config json data",
			fields: fields{store: cacheStore},
			args: args{
				ctx:            nil,
				opts:           webSrvOpts,
				ofp:            ofp.Fingerprints(),
				configJsonData: validJsonData,
			},
			wantErr: false,
		},
		{
			name:   "invalid config json data",
			fields: fields{store: cacheStore},
			args: args{
				ctx:            nil,
				opts:           webSrvOpts,
				ofp:            ofp.Fingerprints(),
				configJsonData: invalidJsonData,
			},
			wantErr: true,
		},
		{
			name:   "inconsistent fingerprints",
			fields: fields{store: cacheStore},
			args: args{
				ctx:            nil,
				opts:           webSrvOpts,
				ofp:            utilsV3.ConfigFingerprints{"C:\\config_test\\nginx.conf": "invalid_fingerprint"},
				configJsonData: validJsonData,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid config json data" {
				wscstore.EXPECT().UpdateConfig(tt.args.ctx, tt.args.opts, tt.args.ofp, tt.args.configJsonData).Return(nil)
			} else if tt.name == "invalid config json data" {
				wscstore.EXPECT().UpdateConfig(tt.args.ctx, tt.args.opts, tt.args.ofp, tt.args.configJsonData).Return(
					errors.New("invalid config json data"),
				)
			} else if tt.name == "inconsistent fingerprints" {
				wscstore.EXPECT().UpdateConfig(tt.args.ctx, tt.args.opts, tt.args.ofp, tt.args.configJsonData).Return(
					errors.Wrapf(metav1.ErrInconsistentFingerprints, "fingerprint mismatch"),
				)
			}

			w := &webServerConfigService{
				store: tt.fields.store,
			}
			if err := w.UpdateConfig(tt.args.ctx, tt.args.opts, tt.args.ofp, tt.args.configJsonData); (err != nil) != tt.wantErr {
				t.Errorf("UpdateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
