package bifrosts

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	svcv1 "gin-vue-admin/internal/hmdr_api/service/v1"
	storev1 "gin-vue-admin/internal/hmdr_api/store/v1"
	"gin-vue-admin/internal/hmdr_api/store/v1/cache"
	storefake "gin-vue-admin/internal/hmdr_api/store/v1/fake"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	utilsV3 "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/utils"
	"go.uber.org/mock/gomock"
	"reflect"
	"strings"
	"testing"
	"time"
)

func Test_newWebServerConfigs(t *testing.T) {
	store := new(storefake.Store)
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
		ServerName: "test",
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := storev1.NewMockFactory(ctrl)
	cacheStore := cache.GetCacheStore(store, time.Second)
	wscstore := storev1.NewMockWebServerConfigStore(ctrl)
	store.EXPECT().WebServerConfigs().AnyTimes().Return(wscstore)
	wscstore.EXPECT().GetConfig(nil, webSrvOpts).AnyTimes().Return(new(storefake.WebServerConfigStore).GetConfig(nil, webSrvOpts))
	_, ofp, err := cacheStore.WebServerConfigs().GetConfig(nil, webSrvOpts)
	if err != nil {
		t.Fatal(err)
	}
	wscstore.EXPECT().ChangeContextEnabledState(nil, webSrvOpts, ofp.Fingerprints(), metav1.TargetConfigContextOptions[metav1.ConfigContextEnabledStateMeta]{
		Position: metav1.ConfigContextPos{
			Config:         "C:\\config_test\\nginx.conf",
			ContextPosPath: []int{0},
		},
		TargetContext: metav1.ConfigContextEnabledStateMeta{Enabled: true},
	}).AnyTimes().Return(new(storefake.WebServerConfigStore).ChangeContextEnabledState(nil, webSrvOpts, ofp.Fingerprints(), metav1.TargetConfigContextOptions[metav1.ConfigContextEnabledStateMeta]{
		Position: metav1.ConfigContextPos{
			Config:         "C:\\config_test\\nginx.conf",
			ContextPosPath: []int{0},
		},
		TargetContext: metav1.ConfigContextEnabledStateMeta{Enabled: true},
	}))
	wscstore.EXPECT().ChangeContextEnabledState(nil, webSrvOpts, ofp.Fingerprints(), metav1.TargetConfigContextOptions[metav1.ConfigContextEnabledStateMeta]{
		Position: metav1.ConfigContextPos{
			Config:         "conf.d\\location.conf",
			ContextPosPath: nil,
		},
		//TargetContext: metav1.ConfigContextEnabledStateMeta{Enabled: false},
	}).AnyTimes().Return(new(storefake.WebServerConfigStore).ChangeContextEnabledState(nil, webSrvOpts, ofp.Fingerprints(), metav1.TargetConfigContextOptions[metav1.ConfigContextEnabledStateMeta]{
		Position: metav1.ConfigContextPos{
			Config:         "conf.d\\location.conf",
			ContextPosPath: nil,
		},
		//TargetContext: metav1.ConfigContextEnabledStateMeta{Enabled: false},
	}))
	wscstore.EXPECT().ChangeContextEnabledState(nil, webSrvOpts, ofp.Fingerprints(), metav1.TargetConfigContextOptions[metav1.ConfigContextEnabledStateMeta]{
		Position: metav1.ConfigContextPos{
			Config:         "conf.d\\location.conf",
			ContextPosPath: []int{1, 2, 3},
		},
		TargetContext: metav1.ConfigContextEnabledStateMeta{Enabled: true},
	}).AnyTimes().Return(new(storefake.WebServerConfigStore).ChangeContextEnabledState(nil, webSrvOpts, ofp.Fingerprints(), metav1.TargetConfigContextOptions[metav1.ConfigContextEnabledStateMeta]{
		Position: metav1.ConfigContextPos{
			Config:         "conf.d\\location.conf",
			ContextPosPath: []int{1, 2, 3},
		},
		TargetContext: metav1.ConfigContextEnabledStateMeta{Enabled: true},
	}))
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
	webSrvOpts := metav1.WebServerOptions{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := storev1.NewMockFactory(ctrl)
	cacheStore := cache.GetCacheStore(store, time.Second)
	wscstore := storev1.NewMockWebServerConfigStore(ctrl)
	store.EXPECT().WebServerConfigs().AnyTimes().Return(wscstore)
	wscstore.EXPECT().GetConfig(nil, webSrvOpts).AnyTimes().Return(new(storefake.WebServerConfigStore).GetConfig(nil, webSrvOpts))
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
			fields: fields{store: cacheStore},
			args:   args{opts: webSrvOpts},
			wantConfigTextLines: []string{
				`# user nobody;`,
				`worker_processes 1;`,
				`# error_log logs/error.log;`,
				`# error_log logs/error.log  notice;`,
				`# error_log logs/error.log  info;`,
				`# pid logs/nginx.pid;`,
				`events {`,
				`    worker_connections 1024;`,
				`}`,
				`stream {`,
				`    upstream vvvvv1 {`,
				`        server test.1.cn:22 weight=5;`,
				`        # server test.2.cn:33 weight=5;`,
				`    }`,
				`    upstream vvvvv2 {`,
				`        server vvvvv1:55;`,
				`        server test.vv2.3.cn:44;`,
				`    }`,
				`    server {`,
				`        listen 5510;`,
				`        # test`,
				`        proxy_connect_timeout 10s;`,
				`        proxy_timeout 300s;`,
				`        proxy_pass www.baidu.com:443;`,
				`    }`,
				`    server {`,
				`        listen 5520;`,
				`        proxy_pass vvvvv1;`,
				`    }`,
				`    server {`,
				`        listen 5530;`,
				`        proxy_pass vvvvv2;`,
				`    }`,
				`}`,
				`http {`,
				`    # include <== mime.types`,
				`    types {`,
				`        text/html html htm shtml;`,
				`        text/css css;`,
				`        text/xml xml;`,
				`        image/gif gif;`,
				`        image/jpeg jpeg jpg;`,
				`        application/javascript js;`,
				`        application/atom+xml atom;`,
				`        application/rss+xml rss;`,
				`        text/mathml mml;`,
				`        text/plain txt;`,
				`        text/vnd.sun.j2me.app-descriptor jad;`,
				`        text/vnd.wap.wml wml;`,
				`        text/x-component htc;`,
				`        image/png png;`,
				`        image/svg+xml svg svgz;`,
				`        image/tiff tif tiff;`,
				`        image/vnd.wap.wbmp wbmp;`,
				`        image/webp webp;`,
				`        image/x-icon ico;`,
				`        image/x-jng jng;`,
				`        image/x-ms-bmp bmp;`,
				`        application/font-woff woff;`,
				`        application/java-archive jar war ear;`,
				`        application/json json;`,
				`        application/mac-binhex40 hqx;`,
				`        application/msword doc;`,
				`        application/pdf pdf;`,
				`        application/postscript ps eps ai;`,
				`        application/rtf rtf;`,
				`        application/vnd.apple.mpegurl m3u8;`,
				`        application/vnd.google-earth.kml+xml kml;`,
				`        application/vnd.google-earth.kmz kmz;`,
				`        application/vnd.ms-excel xls;`,
				`        application/vnd.ms-fontobject eot;`,
				`        application/vnd.ms-powerpoint ppt;`,
				`        application/vnd.oasis.opendocument.graphics odg;`,
				`        application/vnd.oasis.opendocument.presentation odp;`,
				`        application/vnd.oasis.opendocument.spreadsheet ods;`,
				`        application/vnd.oasis.opendocument.text odt;`,
				`        application/vnd.openxmlformats-officedocument.presentationml.presentation pptx;`,
				`        application/vnd.openxmlformats-officedocument.spreadsheetml.sheet xlsx;`,
				`        application/vnd.openxmlformats-officedocument.wordprocessingml.document docx;`,
				`        application/vnd.wap.wmlc wmlc;`,
				`        application/x-7z-compressed 7z;`,
				`        application/x-cocoa cco;`,
				`        application/x-java-archive-diff jardiff;`,
				`        application/x-java-jnlp-file jnlp;`,
				`        application/x-makeself run;`,
				`        application/x-perl pl pm;`,
				`        application/x-pilot prc pdb;`,
				`        application/x-rar-compressed rar;`,
				`        application/x-redhat-package-manager rpm;`,
				`        application/x-sea sea;`,
				`        application/x-shockwave-flash swf;`,
				`        application/x-stuffit sit;`,
				`        application/x-tcl tcl tk;`,
				`        application/x-x509-ca-cert der pem crt;`,
				`        application/x-xpinstall xpi;`,
				`        application/xhtml+xml xhtml;`,
				`        application/xspf+xml xspf;`,
				`        application/zip zip;`,
				`        application/octet-stream bin exe dll;`,
				`        application/octet-stream deb;`,
				`        application/octet-stream dmg;`,
				`        application/octet-stream iso img;`,
				`        application/octet-stream msi msp msm;`,
				`        audio/midi mid midi kar;`,
				`        audio/mpeg mp3;`,
				`        audio/ogg ogg;`,
				`        audio/x-m4a m4a;`,
				`        audio/x-realaudio ra;`,
				`        video/3gpp 3gpp 3gp;`,
				`        video/mp2t ts;`,
				`        video/mp4 mp4;`,
				`        video/mpeg mpeg mpg;`,
				`        video/quicktime mov;`,
				`        video/webm webm;`,
				`        video/x-flv flv;`,
				`        video/x-m4v m4v;`,
				`        video/x-mng mng;`,
				`        video/x-ms-asf asx asf;`,
				`        video/x-ms-wmv wmv;`,
				`        video/x-msvideo avi;`,
				`    }`,
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
				`    server {    # test inline comment for client.UpdateConfig with resolv.V2 at 2021-01-26 17:13:59.357884 +0800 CST m=+0.142618001  # test inline comment for client.UpdateConfig with resolv.V2 at 2021-01-29 16:08:16.7355344 +0800 CST m=+0.083795701`,
				`        listen 80;`,
				`        server_name test1.com;`,
				`        # include <== ./conf.d/location.conf`,
				`        location /test {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # # include <== location.conf`,
				`        upstream upstream.test.1 {`,
				`            server 10.1.1.1:443;`,
				`            server 10.1.1.2:443;`,
				`        }`,
				`        location /test_proxy {`,
				`            proxy_pass upstream.test.1;`,
				`        }`,
				`    }`,
				`    server {`,
				`        listen 80;`,
				`        server_name test2.com;`,
				`        # include <== ./conf.d/location.conf`,
				`        location /test {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # # include <== location.conf`,
				`        upstream upstream.test.1 {`,
				`            server 10.1.1.1:443;`,
				`            server 10.1.1.2:443;`,
				`        }`,
				`        location /test_proxy {`,
				`            proxy_pass upstream.test.1;`,
				`        }`,
				`    }`,
				`    # include <== ./conf.d/server_test*.conf`,
				`    server {`,
				`        listen 80;`,
				`        server_name test1.com;`,
				`        # include <== ./conf.d/location*.conf`,
				`        location /test {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # # include <== location.conf`,
				`        upstream upstream.test.1 {`,
				`            server 10.1.1.1:443;`,
				`            server 10.1.1.2:443;`,
				`        }`,
				`        location /test_proxy {`,
				`            proxy_pass upstream.test.1;`,
				`        }`,
				`        location /test2 {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        upstream upstream.test.3 {`,
				`            server 10.1.3.1:443;`,
				`            server 10.1.3.2:443;`,
				`        }`,
				`        upstream upstream.test.2 {`,
				`            server 10.1.2.1:443;`,
				`            server 10.1.2.2:443;`,
				`            server upstream.test.3;`,
				`        }`,
				`        location /test_proxy {`,
				`            # proxy_pass upstream.test.2;`,
				`        }`,
				`        location /test_proxy2 {`,
				`            proxy_pass https://baidu.com;`,
				`        }`,
				`        location /test_proxy3 {`,
				`            proxy_pass http://baidu2.com/test1;`,
				`        }`,
				`        location /test_proxy4 {`,
				`            proxy_pass http://10.1.1.2:333/test1;`,
				`        }`,
				`    }`,
				`    server {`,
				`        listen 8080;`,
				`        server_name test1.com;`,
				`        # include <== ./conf.d/location*.conf`,
				`        location /test {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # # include <== location.conf`,
				`        upstream upstream.test.1 {`,
				`            server 10.1.1.1:443;`,
				`            server 10.1.1.2:443;`,
				`        }`,
				`        location /test_proxy {`,
				`            proxy_pass upstream.test.1;`,
				`        }`,
				`        location /test2 {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        upstream upstream.test.3 {`,
				`            server 10.1.3.1:443;`,
				`            server 10.1.3.2:443;`,
				`        }`,
				`        upstream upstream.test.2 {`,
				`            server 10.1.2.1:443;`,
				`            server 10.1.2.2:443;`,
				`            server upstream.test.3;`,
				`        }`,
				`        location /test_proxy {`,
				`            # proxy_pass upstream.test.2;`,
				`        }`,
				`        location /test_proxy2 {`,
				`            proxy_pass https://baidu.com;`,
				`        }`,
				`        location /test_proxy3 {`,
				`            proxy_pass http://baidu2.com/test1;`,
				`        }`,
				`        location /test_proxy4 {`,
				`            proxy_pass http://10.1.1.2:333/test1;`,
				`        }`,
				`    }`,
				`    server {`,
				`        listen 8080;`,
				`        server_name test2.com;`,
				`        # include <== ./conf.d/location*.conf`,
				`        location /test {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # # include <== location.conf`,
				`        upstream upstream.test.1 {`,
				`            server 10.1.1.1:443;`,
				`            server 10.1.1.2:443;`,
				`        }`,
				`        location /test_proxy {`,
				`            proxy_pass upstream.test.1;`,
				`        }`,
				`        location /test2 {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        upstream upstream.test.3 {`,
				`            server 10.1.3.1:443;`,
				`            server 10.1.3.2:443;`,
				`        }`,
				`        upstream upstream.test.2 {`,
				`            server 10.1.2.1:443;`,
				`            server 10.1.2.2:443;`,
				`            server upstream.test.3;`,
				`        }`,
				`        location /test_proxy {`,
				`            # proxy_pass upstream.test.2;`,
				`        }`,
				`        location /test_proxy2 {`,
				`            proxy_pass https://baidu.com;`,
				`        }`,
				`        location /test_proxy3 {`,
				`            proxy_pass http://baidu2.com/test1;`,
				`        }`,
				`        location /test_proxy4 {`,
				`            proxy_pass http://10.1.1.2:333/test1;`,
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
				`        # include <== ./conf.d/location.conf`,
				`        location /test {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # # include <== location.conf`,
				`        upstream upstream.test.1 {`,
				`            server 10.1.1.1:443;`,
				`            server 10.1.1.2:443;`,
				`        }`,
				`        location /test_proxy {`,
				`            proxy_pass upstream.test.1;`,
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
				`    server {`,
				`        listen 990 ssl;`,
				`        server_name localhost;`,
				`        # charset koi8-r;`,
				`        # access_log logs/host.access.log  main;`,
				`        location / {`,
				`            root html;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # include <== ./conf.d/location.conf`,
				`        location /test {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # # include <== location.conf`,
				`        upstream upstream.test.1 {`,
				`            server 10.1.1.1:443;`,
				`            server 10.1.1.2:443;`,
				`        }`,
				`        location /test_proxy {`,
				`            proxy_pass upstream.test.1;`,
				`        }`,
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
				"C:\\config_test\\conf.d\\location.conf":     "539813c0f45630e9feba9a10c6494b4d912f0733847a4f17c650492709299c75",
				"C:\\config_test\\conf.d\\location2.conf":    "fc2a0cf89b11602e6dbfe0aa2c98cb69220485e77ef0e32a623043eb125f2114",
				"C:\\config_test\\conf.d\\server_test1.conf": "151c5dd9a238cdd69f4fc35d0564ab448c0e11530862457b41671fd41ddb9a0b",
				"C:\\config_test\\conf.d\\server_test2.conf": "24480b9ef0c9c86cb90896d7871e10bab94cc25fda496050d48533d6bf542f53",
				"C:\\config_test\\conf.d\\test1.com.conf":    "775ed01e78add3b934de529cec247b2a970558d7d1832a3e776e29a30bfc131a",
				"C:\\config_test\\conf.d\\test2.com.conf":    "bd31d5d2604233bbac22fb73e9125375c4bc8fe4c612c4611363b8f93413b2ea",
				"C:\\config_test\\mime.types":                "3c6049a805154dc0122c7264153036205c8f27f69699dc8ba129f212afb66d5a",
				"C:\\config_test\\nginx.conf":                "e2a36380e1591b13cca9d9eb5437bd1a2747901aa5f34caad39aa0960018d492",
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
	webSrvOpts := metav1.WebServerOptions{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := storev1.NewMockFactory(ctrl)
	cacheStore := cache.GetCacheStore(store, time.Second)
	wscstore := storev1.NewMockWebServerConfigStore(ctrl)
	store.EXPECT().WebServerConfigs().AnyTimes().Return(wscstore)
	wscstore.EXPECT().GetConfig(nil, webSrvOpts).AnyTimes().Return(new(storefake.WebServerConfigStore).GetConfig(nil, webSrvOpts))
	//_, _, err := cacheStore.WebServerConfigs().GetConfig(nil, webSrvOpts)
	//if err != nil {
	//	t.Fatal(err)
	//}
	wscstore.EXPECT().GetContext(nil, webSrvOpts, metav1.ConfigContextPos{
		Config:         "C:\\config_test\\conf.d\\location.conf",
		ContextPosPath: []int{0},
	}).AnyTimes().Return(new(storefake.WebServerConfigStore).GetContext(nil, webSrvOpts, metav1.ConfigContextPos{
		Config:         "C:\\config_test\\conf.d\\location.conf",
		ContextPosPath: []int{0},
	}))
	wscstore.EXPECT().GetContext(nil, webSrvOpts, metav1.ConfigContextPos{
		Config:         "C:\\config_test\\conf.d\\location.conf",
		ContextPosPath: []int{0, 1},
	}).AnyTimes().Return(new(storefake.WebServerConfigStore).GetContext(nil, webSrvOpts, metav1.ConfigContextPos{
		Config:         "C:\\config_test\\conf.d\\location.conf",
		ContextPosPath: []int{0, 1},
	}))
	wscstore.EXPECT().GetContext(nil, webSrvOpts, metav1.ConfigContextPos{
		Config:         "C:\\config_test\\nginx.conf",
		ContextPosPath: []int{},
	}).AnyTimes().Return(new(storefake.WebServerConfigStore).GetContext(nil, webSrvOpts, metav1.ConfigContextPos{
		Config:         "C:\\config_test\\nginx.conf",
		ContextPosPath: []int{},
	}))
	wscstore.EXPECT().GetContext(nil, webSrvOpts, metav1.ConfigContextPos{
		Config:         "C:\\config_test\\nginx.conf",
		ContextPosPath: nil,
	}).AnyTimes().Return(new(storefake.WebServerConfigStore).GetContext(nil, webSrvOpts, metav1.ConfigContextPos{
		Config:         "C:\\config_test\\nginx.conf",
		ContextPosPath: nil,
	}))
	wscstore.EXPECT().GetContext(nil, webSrvOpts, metav1.ConfigContextPos{
		Config:         "C:\\config_test\\nginx.conf",
		ContextPosPath: []int{1, 2, 3, 4, 5, 6, 7},
	}).AnyTimes().Return(new(storefake.WebServerConfigStore).GetContext(nil, webSrvOpts, metav1.ConfigContextPos{
		Config:         "C:\\config_test\\nginx.conf",
		ContextPosPath: []int{1, 2, 3, 4, 5, 6, 7},
	}))
	type fields struct {
		store storev1.Factory
	}
	type args struct {
		ctx  context.Context
		opts metav1.WebServerOptions
		pos  metav1.ConfigContextPos
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantConfigLines []string
		wantErr         bool
	}{
		{
			name:   "one level pos path",
			fields: fields{store: cacheStore},
			args: args{
				opts: webSrvOpts,
				pos: metav1.ConfigContextPos{
					Config:         "C:\\config_test\\conf.d\\location.conf",
					ContextPosPath: []int{0},
				},
			},
			wantConfigLines: []string{
				"location /test {",
				"    root html/test;",
				"    index index.html index.htm;",
				"}",
			},
			wantErr: false,
		},
		{
			name:   "two level pos path",
			fields: fields{store: cacheStore},
			args: args{
				opts: webSrvOpts,
				pos: metav1.ConfigContextPos{
					Config:         "C:\\config_test\\conf.d\\location.conf",
					ContextPosPath: []int{0, 1},
				},
			},
			wantConfigLines: []string{
				"index index.html index.htm;",
			},
			wantErr: false,
		},
		{
			name:   "null pos path",
			fields: fields{store: cacheStore},
			args: args{
				opts: webSrvOpts,
				pos: metav1.ConfigContextPos{
					Config:         "C:\\config_test\\nginx.conf",
					ContextPosPath: []int{},
				},
			},
			wantConfigLines: []string{
				`# user nobody;`,
				`worker_processes 1;`,
				`# error_log logs/error.log;`,
				`# error_log logs/error.log  notice;`,
				`# error_log logs/error.log  info;`,
				`# pid logs/nginx.pid;`,
				`events {`,
				`    worker_connections 1024;`,
				`}`,
				`stream {`,
				`    upstream vvvvv1 {`,
				`        server test.1.cn:22 weight=5;`,
				`        # server test.2.cn:33 weight=5;`,
				`    }`,
				`    upstream vvvvv2 {`,
				`        server vvvvv1:55;`,
				`        server test.vv2.3.cn:44;`,
				`    }`,
				`    server {`,
				`        listen 5510;`,
				`        # test`,
				`        proxy_connect_timeout 10s;`,
				`        proxy_timeout 300s;`,
				`        proxy_pass www.baidu.com:443;`,
				`    }`,
				`    server {`,
				`        listen 5520;`,
				`        proxy_pass vvvvv1;`,
				`    }`,
				`    server {`,
				`        listen 5530;`,
				`        proxy_pass vvvvv2;`,
				`    }`,
				`}`,
				`http {`,
				`    # include <== mime.types`,
				`    types {`,
				`        text/html html htm shtml;`,
				`        text/css css;`,
				`        text/xml xml;`,
				`        image/gif gif;`,
				`        image/jpeg jpeg jpg;`,
				`        application/javascript js;`,
				`        application/atom+xml atom;`,
				`        application/rss+xml rss;`,
				`        text/mathml mml;`,
				`        text/plain txt;`,
				`        text/vnd.sun.j2me.app-descriptor jad;`,
				`        text/vnd.wap.wml wml;`,
				`        text/x-component htc;`,
				`        image/png png;`,
				`        image/svg+xml svg svgz;`,
				`        image/tiff tif tiff;`,
				`        image/vnd.wap.wbmp wbmp;`,
				`        image/webp webp;`,
				`        image/x-icon ico;`,
				`        image/x-jng jng;`,
				`        image/x-ms-bmp bmp;`,
				`        application/font-woff woff;`,
				`        application/java-archive jar war ear;`,
				`        application/json json;`,
				`        application/mac-binhex40 hqx;`,
				`        application/msword doc;`,
				`        application/pdf pdf;`,
				`        application/postscript ps eps ai;`,
				`        application/rtf rtf;`,
				`        application/vnd.apple.mpegurl m3u8;`,
				`        application/vnd.google-earth.kml+xml kml;`,
				`        application/vnd.google-earth.kmz kmz;`,
				`        application/vnd.ms-excel xls;`,
				`        application/vnd.ms-fontobject eot;`,
				`        application/vnd.ms-powerpoint ppt;`,
				`        application/vnd.oasis.opendocument.graphics odg;`,
				`        application/vnd.oasis.opendocument.presentation odp;`,
				`        application/vnd.oasis.opendocument.spreadsheet ods;`,
				`        application/vnd.oasis.opendocument.text odt;`,
				`        application/vnd.openxmlformats-officedocument.presentationml.presentation pptx;`,
				`        application/vnd.openxmlformats-officedocument.spreadsheetml.sheet xlsx;`,
				`        application/vnd.openxmlformats-officedocument.wordprocessingml.document docx;`,
				`        application/vnd.wap.wmlc wmlc;`,
				`        application/x-7z-compressed 7z;`,
				`        application/x-cocoa cco;`,
				`        application/x-java-archive-diff jardiff;`,
				`        application/x-java-jnlp-file jnlp;`,
				`        application/x-makeself run;`,
				`        application/x-perl pl pm;`,
				`        application/x-pilot prc pdb;`,
				`        application/x-rar-compressed rar;`,
				`        application/x-redhat-package-manager rpm;`,
				`        application/x-sea sea;`,
				`        application/x-shockwave-flash swf;`,
				`        application/x-stuffit sit;`,
				`        application/x-tcl tcl tk;`,
				`        application/x-x509-ca-cert der pem crt;`,
				`        application/x-xpinstall xpi;`,
				`        application/xhtml+xml xhtml;`,
				`        application/xspf+xml xspf;`,
				`        application/zip zip;`,
				`        application/octet-stream bin exe dll;`,
				`        application/octet-stream deb;`,
				`        application/octet-stream dmg;`,
				`        application/octet-stream iso img;`,
				`        application/octet-stream msi msp msm;`,
				`        audio/midi mid midi kar;`,
				`        audio/mpeg mp3;`,
				`        audio/ogg ogg;`,
				`        audio/x-m4a m4a;`,
				`        audio/x-realaudio ra;`,
				`        video/3gpp 3gpp 3gp;`,
				`        video/mp2t ts;`,
				`        video/mp4 mp4;`,
				`        video/mpeg mpeg mpg;`,
				`        video/quicktime mov;`,
				`        video/webm webm;`,
				`        video/x-flv flv;`,
				`        video/x-m4v m4v;`,
				`        video/x-mng mng;`,
				`        video/x-ms-asf asx asf;`,
				`        video/x-ms-wmv wmv;`,
				`        video/x-msvideo avi;`,
				`    }`,
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
				`    server {    # test inline comment for client.UpdateConfig with resolv.V2 at 2021-01-26 17:13:59.357884 +0800 CST m=+0.142618001  # test inline comment for client.UpdateConfig with resolv.V2 at 2021-01-29 16:08:16.7355344 +0800 CST m=+0.083795701`,
				`        listen 80;`,
				`        server_name test1.com;`,
				`        # include <== ./conf.d/location.conf`,
				`        location /test {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # # include <== location.conf`,
				`        upstream upstream.test.1 {`,
				`            server 10.1.1.1:443;`,
				`            server 10.1.1.2:443;`,
				`        }`,
				`        location /test_proxy {`,
				`            proxy_pass upstream.test.1;`,
				`        }`,
				`    }`,
				`    server {`,
				`        listen 80;`,
				`        server_name test2.com;`,
				`        # include <== ./conf.d/location.conf`,
				`        location /test {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # # include <== location.conf`,
				`        upstream upstream.test.1 {`,
				`            server 10.1.1.1:443;`,
				`            server 10.1.1.2:443;`,
				`        }`,
				`        location /test_proxy {`,
				`            proxy_pass upstream.test.1;`,
				`        }`,
				`    }`,
				`    # include <== ./conf.d/server_test*.conf`,
				`    server {`,
				`        listen 80;`,
				`        server_name test1.com;`,
				`        # include <== ./conf.d/location*.conf`,
				`        location /test {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # # include <== location.conf`,
				`        upstream upstream.test.1 {`,
				`            server 10.1.1.1:443;`,
				`            server 10.1.1.2:443;`,
				`        }`,
				`        location /test_proxy {`,
				`            proxy_pass upstream.test.1;`,
				`        }`,
				`        location /test2 {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        upstream upstream.test.3 {`,
				`            server 10.1.3.1:443;`,
				`            server 10.1.3.2:443;`,
				`        }`,
				`        upstream upstream.test.2 {`,
				`            server 10.1.2.1:443;`,
				`            server 10.1.2.2:443;`,
				`            server upstream.test.3;`,
				`        }`,
				`        location /test_proxy {`,
				`            # proxy_pass upstream.test.2;`,
				`        }`,
				`        location /test_proxy2 {`,
				`            proxy_pass https://baidu.com;`,
				`        }`,
				`        location /test_proxy3 {`,
				`            proxy_pass http://baidu2.com/test1;`,
				`        }`,
				`        location /test_proxy4 {`,
				`            proxy_pass http://10.1.1.2:333/test1;`,
				`        }`,
				`    }`,
				`    server {`,
				`        listen 8080;`,
				`        server_name test1.com;`,
				`        # include <== ./conf.d/location*.conf`,
				`        location /test {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # # include <== location.conf`,
				`        upstream upstream.test.1 {`,
				`            server 10.1.1.1:443;`,
				`            server 10.1.1.2:443;`,
				`        }`,
				`        location /test_proxy {`,
				`            proxy_pass upstream.test.1;`,
				`        }`,
				`        location /test2 {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        upstream upstream.test.3 {`,
				`            server 10.1.3.1:443;`,
				`            server 10.1.3.2:443;`,
				`        }`,
				`        upstream upstream.test.2 {`,
				`            server 10.1.2.1:443;`,
				`            server 10.1.2.2:443;`,
				`            server upstream.test.3;`,
				`        }`,
				`        location /test_proxy {`,
				`            # proxy_pass upstream.test.2;`,
				`        }`,
				`        location /test_proxy2 {`,
				`            proxy_pass https://baidu.com;`,
				`        }`,
				`        location /test_proxy3 {`,
				`            proxy_pass http://baidu2.com/test1;`,
				`        }`,
				`        location /test_proxy4 {`,
				`            proxy_pass http://10.1.1.2:333/test1;`,
				`        }`,
				`    }`,
				`    server {`,
				`        listen 8080;`,
				`        server_name test2.com;`,
				`        # include <== ./conf.d/location*.conf`,
				`        location /test {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # # include <== location.conf`,
				`        upstream upstream.test.1 {`,
				`            server 10.1.1.1:443;`,
				`            server 10.1.1.2:443;`,
				`        }`,
				`        location /test_proxy {`,
				`            proxy_pass upstream.test.1;`,
				`        }`,
				`        location /test2 {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        upstream upstream.test.3 {`,
				`            server 10.1.3.1:443;`,
				`            server 10.1.3.2:443;`,
				`        }`,
				`        upstream upstream.test.2 {`,
				`            server 10.1.2.1:443;`,
				`            server 10.1.2.2:443;`,
				`            server upstream.test.3;`,
				`        }`,
				`        location /test_proxy {`,
				`            # proxy_pass upstream.test.2;`,
				`        }`,
				`        location /test_proxy2 {`,
				`            proxy_pass https://baidu.com;`,
				`        }`,
				`        location /test_proxy3 {`,
				`            proxy_pass http://baidu2.com/test1;`,
				`        }`,
				`        location /test_proxy4 {`,
				`            proxy_pass http://10.1.1.2:333/test1;`,
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
				`        # include <== ./conf.d/location.conf`,
				`        location /test {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # # include <== location.conf`,
				`        upstream upstream.test.1 {`,
				`            server 10.1.1.1:443;`,
				`            server 10.1.1.2:443;`,
				`        }`,
				`        location /test_proxy {`,
				`            proxy_pass upstream.test.1;`,
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
				`    server {`,
				`        listen 990 ssl;`,
				`        server_name localhost;`,
				`        # charset koi8-r;`,
				`        # access_log logs/host.access.log  main;`,
				`        location / {`,
				`            root html;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # include <== ./conf.d/location.conf`,
				`        location /test {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # # include <== location.conf`,
				`        upstream upstream.test.1 {`,
				`            server 10.1.1.1:443;`,
				`            server 10.1.1.2:443;`,
				`        }`,
				`        location /test_proxy {`,
				`            proxy_pass upstream.test.1;`,
				`        }`,
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
			fields: fields{store: cacheStore},
			args: args{
				opts: webSrvOpts,
				pos: metav1.ConfigContextPos{
					Config:         "C:\\config_test\\nginx.conf",
					ContextPosPath: nil,
				},
			},
			wantConfigLines: []string{
				`# user nobody;`,
				`worker_processes 1;`,
				`# error_log logs/error.log;`,
				`# error_log logs/error.log  notice;`,
				`# error_log logs/error.log  info;`,
				`# pid logs/nginx.pid;`,
				`events {`,
				`    worker_connections 1024;`,
				`}`,
				`stream {`,
				`    upstream vvvvv1 {`,
				`        server test.1.cn:22 weight=5;`,
				`        # server test.2.cn:33 weight=5;`,
				`    }`,
				`    upstream vvvvv2 {`,
				`        server vvvvv1:55;`,
				`        server test.vv2.3.cn:44;`,
				`    }`,
				`    server {`,
				`        listen 5510;`,
				`        # test`,
				`        proxy_connect_timeout 10s;`,
				`        proxy_timeout 300s;`,
				`        proxy_pass www.baidu.com:443;`,
				`    }`,
				`    server {`,
				`        listen 5520;`,
				`        proxy_pass vvvvv1;`,
				`    }`,
				`    server {`,
				`        listen 5530;`,
				`        proxy_pass vvvvv2;`,
				`    }`,
				`}`,
				`http {`,
				`    # include <== mime.types`,
				`    types {`,
				`        text/html html htm shtml;`,
				`        text/css css;`,
				`        text/xml xml;`,
				`        image/gif gif;`,
				`        image/jpeg jpeg jpg;`,
				`        application/javascript js;`,
				`        application/atom+xml atom;`,
				`        application/rss+xml rss;`,
				`        text/mathml mml;`,
				`        text/plain txt;`,
				`        text/vnd.sun.j2me.app-descriptor jad;`,
				`        text/vnd.wap.wml wml;`,
				`        text/x-component htc;`,
				`        image/png png;`,
				`        image/svg+xml svg svgz;`,
				`        image/tiff tif tiff;`,
				`        image/vnd.wap.wbmp wbmp;`,
				`        image/webp webp;`,
				`        image/x-icon ico;`,
				`        image/x-jng jng;`,
				`        image/x-ms-bmp bmp;`,
				`        application/font-woff woff;`,
				`        application/java-archive jar war ear;`,
				`        application/json json;`,
				`        application/mac-binhex40 hqx;`,
				`        application/msword doc;`,
				`        application/pdf pdf;`,
				`        application/postscript ps eps ai;`,
				`        application/rtf rtf;`,
				`        application/vnd.apple.mpegurl m3u8;`,
				`        application/vnd.google-earth.kml+xml kml;`,
				`        application/vnd.google-earth.kmz kmz;`,
				`        application/vnd.ms-excel xls;`,
				`        application/vnd.ms-fontobject eot;`,
				`        application/vnd.ms-powerpoint ppt;`,
				`        application/vnd.oasis.opendocument.graphics odg;`,
				`        application/vnd.oasis.opendocument.presentation odp;`,
				`        application/vnd.oasis.opendocument.spreadsheet ods;`,
				`        application/vnd.oasis.opendocument.text odt;`,
				`        application/vnd.openxmlformats-officedocument.presentationml.presentation pptx;`,
				`        application/vnd.openxmlformats-officedocument.spreadsheetml.sheet xlsx;`,
				`        application/vnd.openxmlformats-officedocument.wordprocessingml.document docx;`,
				`        application/vnd.wap.wmlc wmlc;`,
				`        application/x-7z-compressed 7z;`,
				`        application/x-cocoa cco;`,
				`        application/x-java-archive-diff jardiff;`,
				`        application/x-java-jnlp-file jnlp;`,
				`        application/x-makeself run;`,
				`        application/x-perl pl pm;`,
				`        application/x-pilot prc pdb;`,
				`        application/x-rar-compressed rar;`,
				`        application/x-redhat-package-manager rpm;`,
				`        application/x-sea sea;`,
				`        application/x-shockwave-flash swf;`,
				`        application/x-stuffit sit;`,
				`        application/x-tcl tcl tk;`,
				`        application/x-x509-ca-cert der pem crt;`,
				`        application/x-xpinstall xpi;`,
				`        application/xhtml+xml xhtml;`,
				`        application/xspf+xml xspf;`,
				`        application/zip zip;`,
				`        application/octet-stream bin exe dll;`,
				`        application/octet-stream deb;`,
				`        application/octet-stream dmg;`,
				`        application/octet-stream iso img;`,
				`        application/octet-stream msi msp msm;`,
				`        audio/midi mid midi kar;`,
				`        audio/mpeg mp3;`,
				`        audio/ogg ogg;`,
				`        audio/x-m4a m4a;`,
				`        audio/x-realaudio ra;`,
				`        video/3gpp 3gpp 3gp;`,
				`        video/mp2t ts;`,
				`        video/mp4 mp4;`,
				`        video/mpeg mpeg mpg;`,
				`        video/quicktime mov;`,
				`        video/webm webm;`,
				`        video/x-flv flv;`,
				`        video/x-m4v m4v;`,
				`        video/x-mng mng;`,
				`        video/x-ms-asf asx asf;`,
				`        video/x-ms-wmv wmv;`,
				`        video/x-msvideo avi;`,
				`    }`,
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
				`    server {    # test inline comment for client.UpdateConfig with resolv.V2 at 2021-01-26 17:13:59.357884 +0800 CST m=+0.142618001  # test inline comment for client.UpdateConfig with resolv.V2 at 2021-01-29 16:08:16.7355344 +0800 CST m=+0.083795701`,
				`        listen 80;`,
				`        server_name test1.com;`,
				`        # include <== ./conf.d/location.conf`,
				`        location /test {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # # include <== location.conf`,
				`        upstream upstream.test.1 {`,
				`            server 10.1.1.1:443;`,
				`            server 10.1.1.2:443;`,
				`        }`,
				`        location /test_proxy {`,
				`            proxy_pass upstream.test.1;`,
				`        }`,
				`    }`,
				`    server {`,
				`        listen 80;`,
				`        server_name test2.com;`,
				`        # include <== ./conf.d/location.conf`,
				`        location /test {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # # include <== location.conf`,
				`        upstream upstream.test.1 {`,
				`            server 10.1.1.1:443;`,
				`            server 10.1.1.2:443;`,
				`        }`,
				`        location /test_proxy {`,
				`            proxy_pass upstream.test.1;`,
				`        }`,
				`    }`,
				`    # include <== ./conf.d/server_test*.conf`,
				`    server {`,
				`        listen 80;`,
				`        server_name test1.com;`,
				`        # include <== ./conf.d/location*.conf`,
				`        location /test {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # # include <== location.conf`,
				`        upstream upstream.test.1 {`,
				`            server 10.1.1.1:443;`,
				`            server 10.1.1.2:443;`,
				`        }`,
				`        location /test_proxy {`,
				`            proxy_pass upstream.test.1;`,
				`        }`,
				`        location /test2 {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        upstream upstream.test.3 {`,
				`            server 10.1.3.1:443;`,
				`            server 10.1.3.2:443;`,
				`        }`,
				`        upstream upstream.test.2 {`,
				`            server 10.1.2.1:443;`,
				`            server 10.1.2.2:443;`,
				`            server upstream.test.3;`,
				`        }`,
				`        location /test_proxy {`,
				`            # proxy_pass upstream.test.2;`,
				`        }`,
				`        location /test_proxy2 {`,
				`            proxy_pass https://baidu.com;`,
				`        }`,
				`        location /test_proxy3 {`,
				`            proxy_pass http://baidu2.com/test1;`,
				`        }`,
				`        location /test_proxy4 {`,
				`            proxy_pass http://10.1.1.2:333/test1;`,
				`        }`,
				`    }`,
				`    server {`,
				`        listen 8080;`,
				`        server_name test1.com;`,
				`        # include <== ./conf.d/location*.conf`,
				`        location /test {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # # include <== location.conf`,
				`        upstream upstream.test.1 {`,
				`            server 10.1.1.1:443;`,
				`            server 10.1.1.2:443;`,
				`        }`,
				`        location /test_proxy {`,
				`            proxy_pass upstream.test.1;`,
				`        }`,
				`        location /test2 {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        upstream upstream.test.3 {`,
				`            server 10.1.3.1:443;`,
				`            server 10.1.3.2:443;`,
				`        }`,
				`        upstream upstream.test.2 {`,
				`            server 10.1.2.1:443;`,
				`            server 10.1.2.2:443;`,
				`            server upstream.test.3;`,
				`        }`,
				`        location /test_proxy {`,
				`            # proxy_pass upstream.test.2;`,
				`        }`,
				`        location /test_proxy2 {`,
				`            proxy_pass https://baidu.com;`,
				`        }`,
				`        location /test_proxy3 {`,
				`            proxy_pass http://baidu2.com/test1;`,
				`        }`,
				`        location /test_proxy4 {`,
				`            proxy_pass http://10.1.1.2:333/test1;`,
				`        }`,
				`    }`,
				`    server {`,
				`        listen 8080;`,
				`        server_name test2.com;`,
				`        # include <== ./conf.d/location*.conf`,
				`        location /test {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # # include <== location.conf`,
				`        upstream upstream.test.1 {`,
				`            server 10.1.1.1:443;`,
				`            server 10.1.1.2:443;`,
				`        }`,
				`        location /test_proxy {`,
				`            proxy_pass upstream.test.1;`,
				`        }`,
				`        location /test2 {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        upstream upstream.test.3 {`,
				`            server 10.1.3.1:443;`,
				`            server 10.1.3.2:443;`,
				`        }`,
				`        upstream upstream.test.2 {`,
				`            server 10.1.2.1:443;`,
				`            server 10.1.2.2:443;`,
				`            server upstream.test.3;`,
				`        }`,
				`        location /test_proxy {`,
				`            # proxy_pass upstream.test.2;`,
				`        }`,
				`        location /test_proxy2 {`,
				`            proxy_pass https://baidu.com;`,
				`        }`,
				`        location /test_proxy3 {`,
				`            proxy_pass http://baidu2.com/test1;`,
				`        }`,
				`        location /test_proxy4 {`,
				`            proxy_pass http://10.1.1.2:333/test1;`,
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
				`        # include <== ./conf.d/location.conf`,
				`        location /test {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # # include <== location.conf`,
				`        upstream upstream.test.1 {`,
				`            server 10.1.1.1:443;`,
				`            server 10.1.1.2:443;`,
				`        }`,
				`        location /test_proxy {`,
				`            proxy_pass upstream.test.1;`,
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
				`    server {`,
				`        listen 990 ssl;`,
				`        server_name localhost;`,
				`        # charset koi8-r;`,
				`        # access_log logs/host.access.log  main;`,
				`        location / {`,
				`            root html;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # include <== ./conf.d/location.conf`,
				`        location /test {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # # include <== location.conf`,
				`        upstream upstream.test.1 {`,
				`            server 10.1.1.1:443;`,
				`            server 10.1.1.2:443;`,
				`        }`,
				`        location /test_proxy {`,
				`            proxy_pass upstream.test.1;`,
				`        }`,
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
			fields: fields{store: cacheStore},
			args: args{
				opts: webSrvOpts,
				pos: metav1.ConfigContextPos{
					Config:         "C:\\config_test\\nginx.conf",
					ContextPosPath: []int{1, 2, 3, 4, 5, 6, 7},
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
			got, err := w.GetContext(tt.args.ctx, tt.args.opts, tt.args.pos)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotConfigLines, _ := got.ConfigLines(false)
			if !reflect.DeepEqual(gotConfigLines, tt.wantConfigLines) {
				//for _, line := range gotConfigLines {
				//	fmt.Printf("`%s`,\n", line)
				//}
				t.Errorf("GetContext() got.ConfigLines( false ) = %v, want %v", gotConfigLines, tt.wantConfigLines)
			}
		})
	}
}

func Test_webServerConfigService_GetIncludedConfigs(t *testing.T) {
	webSrvOpts := metav1.WebServerOptions{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := storev1.NewMockFactory(ctrl)
	cacheStore := cache.GetCacheStore(store, time.Second)
	wscstore := storev1.NewMockWebServerConfigStore(ctrl)
	store.EXPECT().WebServerConfigs().AnyTimes().Return(wscstore)
	wscstore.EXPECT().GetConfig(nil, webSrvOpts).AnyTimes().Return(new(storefake.WebServerConfigStore).GetConfig(nil, webSrvOpts))
	//_, _, err := cacheStore.WebServerConfigs().GetConfig(nil, webSrvOpts)
	//if err != nil {
	//	t.Fatal(err)
	//}
	wscstore.EXPECT().GetIncludedConfigs(nil, webSrvOpts, metav1.ConfigContextPos{
		Config:         "C:\\config_test\\conf.d\\server_test1.conf",
		ContextPosPath: []int{0, 2},
	}).AnyTimes().Return(new(storefake.WebServerConfigStore).GetIncludedConfigs(nil, webSrvOpts, metav1.ConfigContextPos{
		Config:         "C:\\config_test\\conf.d\\server_test1.conf",
		ContextPosPath: []int{0, 2},
	}))
	wscstore.EXPECT().GetIncludedConfigs(nil, webSrvOpts, metav1.ConfigContextPos{
		Config:         "C:\\config_test\\conf.d\\server_test1.con",
		ContextPosPath: nil,
	}).AnyTimes().Return(new(storefake.WebServerConfigStore).GetIncludedConfigs(nil, webSrvOpts, metav1.ConfigContextPos{
		Config:         "C:\\config_test\\conf.d\\server_test1.con",
		ContextPosPath: nil,
	}))
	wscstore.EXPECT().GetIncludedConfigs(nil, webSrvOpts, metav1.ConfigContextPos{
		Config:         "C:\\config_test\\conf.d\\server_test1.conf",
		ContextPosPath: []int{0, 1},
	}).AnyTimes().Return(new(storefake.WebServerConfigStore).GetIncludedConfigs(nil, webSrvOpts, metav1.ConfigContextPos{
		Config:         "C:\\config_test\\conf.d\\server_test1.conf",
		ContextPosPath: []int{0, 1},
	}))
	type fields struct {
		store storev1.Factory
	}
	type args struct {
		ctx  context.Context
		opts metav1.WebServerOptions
		pos  metav1.ConfigContextPos
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:   "normal test",
			fields: fields{store: cacheStore},
			args: args{
				opts: webSrvOpts,
				pos: metav1.ConfigContextPos{
					Config:         "C:\\config_test\\conf.d\\server_test1.conf",
					ContextPosPath: []int{0, 2},
				},
			},
			want: []string{
				"C:\\config_test\\conf.d\\location.conf",
				"C:\\config_test\\conf.d\\location2.conf",
			},
			wantErr: false,
		},
		{
			name:   "wrong position",
			fields: fields{store: cacheStore},
			args: args{
				opts: webSrvOpts,
				pos: metav1.ConfigContextPos{
					Config:         "C:\\config_test\\conf.d\\server_test1.con",
					ContextPosPath: nil,
				},
			},
			wantErr: true,
		},
		{
			name:   "the target is not an `include` context",
			fields: fields{store: cacheStore},
			args: args{
				opts: webSrvOpts,
				pos: metav1.ConfigContextPos{
					Config:         "C:\\config_test\\conf.d\\server_test1.conf",
					ContextPosPath: []int{0, 1},
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
			got, err := w.GetIncludedConfigs(tt.args.ctx, tt.args.opts, tt.args.pos)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIncludedConfigs() error = %v, wantErr %v", err, tt.wantErr)
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
	webSrvOpts := metav1.WebServerOptions{}
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
	wscstore.EXPECT().GetConfig(nil, webSrvOpts).AnyTimes().Return(new(storefake.WebServerConfigStore).GetConfig(nil, webSrvOpts))
	_, ofp, err := cacheStore.WebServerConfigs().GetConfig(nil, webSrvOpts)
	if err != nil {
		t.Fatal(err)
	}
	wscstore.EXPECT().InsertWithClone(nil, webSrvOpts, ofp.Fingerprints(), ctxmeta).AnyTimes().Return(new(storefake.WebServerConfigStore).InsertWithClone(nil, webSrvOpts, ofp.Fingerprints(), ctxmeta))
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
			if err := w.InsertWithClone(tt.args.ctx, tt.args.opts, tt.args.ofp, tt.args.ctxmeta); (err != nil) != tt.wantErr {
				t.Errorf("InsertWithClone() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigService_InsertWithNew(t *testing.T) {
	webSrvOpts := metav1.WebServerOptions{}
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
	wscstore.EXPECT().GetConfig(nil, webSrvOpts).AnyTimes().Return(new(storefake.WebServerConfigStore).GetConfig(nil, webSrvOpts))
	_, ofp, err := cacheStore.WebServerConfigs().GetConfig(nil, webSrvOpts)
	if err != nil {
		t.Fatal(err)
	}
	wscstore.EXPECT().InsertWithNew(nil, webSrvOpts, ofp.Fingerprints(), ctxmeta).AnyTimes().Return(new(storefake.WebServerConfigStore).InsertWithNew(nil, webSrvOpts, ofp.Fingerprints(), ctxmeta))
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
			if err := w.InsertWithNew(tt.args.ctx, tt.args.opts, tt.args.ofp, tt.args.ctxmeta); (err != nil) != tt.wantErr {
				t.Errorf("InsertWithNew() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigService_ModifyContextValue(t *testing.T) {
	webSrvOpts := metav1.WebServerOptions{}
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
	wscstore.EXPECT().GetConfig(nil, webSrvOpts).AnyTimes().Return(new(storefake.WebServerConfigStore).GetConfig(nil, webSrvOpts))
	_, ofp, err := cacheStore.WebServerConfigs().GetConfig(nil, webSrvOpts)
	if err != nil {
		t.Fatal(err)
	}
	wscstore.EXPECT().ModifyContextValue(nil, webSrvOpts, ofp.Fingerprints(), ctxmeta).AnyTimes().Return(new(storefake.WebServerConfigStore).ModifyContextValue(nil, webSrvOpts, ofp.Fingerprints(), ctxmeta))
	wscstore.EXPECT().ModifyContextValue(nil, webSrvOpts, ofp.Fingerprints(), unmatchedTypeCtxmeta).AnyTimes().Return(new(storefake.WebServerConfigStore).ModifyContextValue(nil, webSrvOpts, ofp.Fingerprints(), unmatchedTypeCtxmeta))
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
	webSrvOpts := metav1.WebServerOptions{}
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
	wscstore.EXPECT().GetConfig(nil, webSrvOpts).AnyTimes().Return(new(storefake.WebServerConfigStore).GetConfig(nil, webSrvOpts))
	_, ofp, err := cacheStore.WebServerConfigs().GetConfig(nil, webSrvOpts)
	if err != nil {
		t.Fatal(err)
	}
	wscstore.EXPECT().ModifyWithClone(nil, webSrvOpts, ofp.Fingerprints(), ctxmeta).AnyTimes().Return(new(storefake.WebServerConfigStore).ModifyWithClone(nil, webSrvOpts, ofp.Fingerprints(), ctxmeta))
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
	webSrvOpts := metav1.WebServerOptions{}
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
	wscstore.EXPECT().GetConfig(nil, webSrvOpts).AnyTimes().Return(new(storefake.WebServerConfigStore).GetConfig(nil, webSrvOpts))
	_, ofp, err := cacheStore.WebServerConfigs().GetConfig(nil, webSrvOpts)
	if err != nil {
		t.Fatal(err)
	}
	wscstore.EXPECT().ModifyWithNew(nil, webSrvOpts, ofp.Fingerprints(), ctxmeta).AnyTimes().Return(new(storefake.WebServerConfigStore).ModifyWithNew(nil, webSrvOpts, ofp.Fingerprints(), ctxmeta))
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
	webSrvOpts := metav1.WebServerOptions{}
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
	wscstore.EXPECT().GetConfig(nil, webSrvOpts).AnyTimes().Return(new(storefake.WebServerConfigStore).GetConfig(nil, webSrvOpts))
	_, ofp, err := cacheStore.WebServerConfigs().GetConfig(nil, webSrvOpts)
	if err != nil {
		t.Fatal(err)
	}

	wscstore.EXPECT().Remove(nil, webSrvOpts, ofp.Fingerprints(), pos).AnyTimes().Return(new(storefake.WebServerConfigStore).Remove(nil, webSrvOpts, ofp.Fingerprints(), pos))
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
	webSrvOpts := metav1.WebServerOptions{}
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
	wscstore.EXPECT().GetConfig(nil, webSrvOpts).AnyTimes().Return(new(storefake.WebServerConfigStore).GetConfig(nil, webSrvOpts))
	_, ofp, err := cacheStore.WebServerConfigs().GetConfig(nil, webSrvOpts)
	if err != nil {
		t.Fatal(err)
	}
	wscstore.EXPECT().Move(nil, webSrvOpts, ofp.Fingerprints(), ctxmeta).AnyTimes().Return(new(storefake.WebServerConfigStore).Move(nil, webSrvOpts, ofp.Fingerprints(), ctxmeta))
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
			if err := w.Move(tt.args.ctx, tt.args.opts, tt.args.ofp, tt.args.ctxmeta); (err != nil) != tt.wantErr {
				t.Errorf("Move() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
