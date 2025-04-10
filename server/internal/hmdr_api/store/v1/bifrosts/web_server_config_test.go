package bifrosts

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	"gin-vue-admin/global"
	"gin-vue-admin/internal/pkg/bifrosts"
	"gin-vue-admin/internal/pkg/bifrosts/fake"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	bifrostclinetv1 "github.com/ClessLi/bifrost/pkg/client/bifrost/v1"
	utilsV3 "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/utils"
	log "github.com/ClessLi/component-base/pkg/log/v1"
	"github.com/marmotedu/errors"
	"go.uber.org/mock/gomock"
	"reflect"
	"testing"

	nginx_context "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
)

func Test_newConfigContext(t *testing.T) {
	type args struct {
		meta metav1.NewConfigContextMeta
	}
	tests := []struct {
		name    string
		args    args
		wantCtx nginx_context.Context
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCtx := newConfigContext(tt.args.meta); !reflect.DeepEqual(gotCtx, tt.wantCtx) {
				t.Errorf("newConfigContext() = %v, want %v", gotCtx, tt.wantCtx)
			}
		})
	}
}

func Test_newWebServerConfigStore(t *testing.T) {
	type args struct {
		store *bifrostsStore
	}
	tests := []struct {
		name string
		args args
		want *webServerConfigStore
	}{
		// TODO: Add test cases.
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
	webSrvOpts := metav1.WebServerOptions{}
	ofp := utilsV3.ConfigFingerprints{
		"C:\\config_test\\conf.d\\location.conf":     "539813c0f45630e9feba9a10c6494b4d912f0733847a4f17c650492709299c75",
		"C:\\config_test\\conf.d\\location2.conf":    "fc2a0cf89b11602e6dbfe0aa2c98cb69220485e77ef0e32a623043eb125f2114",
		"C:\\config_test\\conf.d\\server_test1.conf": "151c5dd9a238cdd69f4fc35d0564ab448c0e11530862457b41671fd41ddb9a0b",
		"C:\\config_test\\conf.d\\server_test2.conf": "24480b9ef0c9c86cb90896d7871e10bab94cc25fda496050d48533d6bf542f53",
		"C:\\config_test\\conf.d\\test1.com.conf":    "775ed01e78add3b934de529cec247b2a970558d7d1832a3e776e29a30bfc131a",
		"C:\\config_test\\conf.d\\test2.com.conf":    "bd31d5d2604233bbac22fb73e9125375c4bc8fe4c612c4611363b8f93413b2ea",
		"C:\\config_test\\mime.types":                "3c6049a805154dc0122c7264153036205c8f27f69699dc8ba129f212afb66d5a",
		"C:\\config_test\\nginx.conf":                "e2a36380e1591b13cca9d9eb5437bd1a2747901aa5f34caad39aa0960018d492",
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
		name    string
		fields  fields
		args    args
		wantErr bool
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
			fields: fields{bm: mockBifrostsManager},
			args: args{
				opts: webSrvOpts,
				ofp:  ofp,
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
			w := &webServerConfigStore{
				bm: tt.fields.bm,
			}
			if err := w.ChangeContextEnabledState(tt.args.in0, tt.args.opts, tt.args.ofp, tt.args.ctxmeta); (err != nil) != tt.wantErr {
				t.Errorf("ChangeContextEnabledState() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigStore_GetConfig(t *testing.T) {
	webSrvOpts := metav1.WebServerOptions{}
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
			wantErr: false,
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
	global.GVA_LOG = log.ZapLogger()
	webSrvOpts := metav1.WebServerOptions{}
	ofp := utilsV3.ConfigFingerprints{
		"C:\\config_test\\conf.d\\location.conf":     "539813c0f45630e9feba9a10c6494b4d912f0733847a4f17c650492709299c75",
		"C:\\config_test\\conf.d\\location2.conf":    "fc2a0cf89b11602e6dbfe0aa2c98cb69220485e77ef0e32a623043eb125f2114",
		"C:\\config_test\\conf.d\\server_test1.conf": "151c5dd9a238cdd69f4fc35d0564ab448c0e11530862457b41671fd41ddb9a0b",
		"C:\\config_test\\conf.d\\server_test2.conf": "24480b9ef0c9c86cb90896d7871e10bab94cc25fda496050d48533d6bf542f53",
		"C:\\config_test\\conf.d\\test1.com.conf":    "775ed01e78add3b934de529cec247b2a970558d7d1832a3e776e29a30bfc131a",
		"C:\\config_test\\conf.d\\test2.com.conf":    "bd31d5d2604233bbac22fb73e9125375c4bc8fe4c612c4611363b8f93413b2ea",
		"C:\\config_test\\mime.types":                "3c6049a805154dc0122c7264153036205c8f27f69699dc8ba129f212afb66d5a",
		"C:\\config_test\\nginx.conf":                "e2a36380e1591b13cca9d9eb5437bd1a2747901aa5f34caad39aa0960018d492",
	}
	difffp := utilsV3.ConfigFingerprints{
		"C:\\config_test\\conf.d\\location.conf":     "539813c0f45630e9feba9a10c6494b4d912f0733847a4f17c650492709299c75",
		"C:\\config_test\\conf.d\\location2.conf":    "fc2a0cf89b11602e6dbfe0aa2c98cb69220485e77ef0e32a623043eb125f2114",
		"C:\\config_test\\conf.d\\server_test1.conf": "151c5dd9a238cdd69f4fc35d0564ab448c0e11530862457b41671fd41ddb9a0b",
		"C:\\config_test\\conf.d\\server_test2.conf": "24480b9ef0c9c86cb90896d7871e10bab94cc25fda496050d48533d6bf542f53",
		"C:\\config_test\\conf.d\\test1.com.conf":    "775ed01e78add3b934de529cec247b2a970558d7d1832a3e776e29a30bfc131a",
		"C:\\config_test\\conf.d\\test2.com.conf":    "bd31d5d2604233bbac22fb73e9125375c4bc8fe4c612c4611363b8f93413b2ea",
		"C:\\config_test\\mime.types":                "3c6049a805154dc0122c7264153036205c8f27f69699dc8ba129f212afb66d5a",
		"C:\\config_test\\nginx.conf":                "1111111111111111111111111111111111111111111111111111111111111111",
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
				"location /test {",
				"    root html/test;",
				"    index index.html index.htm;",
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
				"index index.html index.htm;",
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
	global.GVA_LOG = log.ZapLogger()
	webSrvOpts := metav1.WebServerOptions{}
	ofp := utilsV3.ConfigFingerprints{
		"C:\\config_test\\conf.d\\location.conf":     "539813c0f45630e9feba9a10c6494b4d912f0733847a4f17c650492709299c75",
		"C:\\config_test\\conf.d\\location2.conf":    "fc2a0cf89b11602e6dbfe0aa2c98cb69220485e77ef0e32a623043eb125f2114",
		"C:\\config_test\\conf.d\\server_test1.conf": "151c5dd9a238cdd69f4fc35d0564ab448c0e11530862457b41671fd41ddb9a0b",
		"C:\\config_test\\conf.d\\server_test2.conf": "24480b9ef0c9c86cb90896d7871e10bab94cc25fda496050d48533d6bf542f53",
		"C:\\config_test\\conf.d\\test1.com.conf":    "775ed01e78add3b934de529cec247b2a970558d7d1832a3e776e29a30bfc131a",
		"C:\\config_test\\conf.d\\test2.com.conf":    "bd31d5d2604233bbac22fb73e9125375c4bc8fe4c612c4611363b8f93413b2ea",
		"C:\\config_test\\mime.types":                "3c6049a805154dc0122c7264153036205c8f27f69699dc8ba129f212afb66d5a",
		"C:\\config_test\\nginx.conf":                "e2a36380e1591b13cca9d9eb5437bd1a2747901aa5f34caad39aa0960018d492",
	}
	difffp := utilsV3.ConfigFingerprints{
		"C:\\config_test\\conf.d\\location.conf":     "539813c0f45630e9feba9a10c6494b4d912f0733847a4f17c650492709299c75",
		"C:\\config_test\\conf.d\\location2.conf":    "fc2a0cf89b11602e6dbfe0aa2c98cb69220485e77ef0e32a623043eb125f2114",
		"C:\\config_test\\conf.d\\server_test1.conf": "151c5dd9a238cdd69f4fc35d0564ab448c0e11530862457b41671fd41ddb9a0b",
		"C:\\config_test\\conf.d\\server_test2.conf": "24480b9ef0c9c86cb90896d7871e10bab94cc25fda496050d48533d6bf542f53",
		"C:\\config_test\\conf.d\\test1.com.conf":    "775ed01e78add3b934de529cec247b2a970558d7d1832a3e776e29a30bfc131a",
		"C:\\config_test\\conf.d\\test2.com.conf":    "bd31d5d2604233bbac22fb73e9125375c4bc8fe4c612c4611363b8f93413b2ea",
		"C:\\config_test\\mime.types":                "3c6049a805154dc0122c7264153036205c8f27f69699dc8ba129f212afb66d5a",
		"C:\\config_test\\nginx.conf":                "1111111111111111111111111111111111111111111111111111111111111111",
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
					Config:         "C:\\config_test\\conf.d\\server_test1.conf",
					ContextPosPath: []int{0, 1},
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
					Config:         "C:\\config_test\\conf.d\\server_test1.conf",
					ContextPosPath: []int{0, 2},
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
		// TODO: Add test cases.
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
				t.Errorf("GetOptions() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webServerConfigStore_InsertWithClone(t *testing.T) {
	webSrvOpts := metav1.WebServerOptions{}
	ofp := utilsV3.ConfigFingerprints{
		"C:\\config_test\\conf.d\\location.conf":     "539813c0f45630e9feba9a10c6494b4d912f0733847a4f17c650492709299c75",
		"C:\\config_test\\conf.d\\location2.conf":    "fc2a0cf89b11602e6dbfe0aa2c98cb69220485e77ef0e32a623043eb125f2114",
		"C:\\config_test\\conf.d\\server_test1.conf": "151c5dd9a238cdd69f4fc35d0564ab448c0e11530862457b41671fd41ddb9a0b",
		"C:\\config_test\\conf.d\\server_test2.conf": "24480b9ef0c9c86cb90896d7871e10bab94cc25fda496050d48533d6bf542f53",
		"C:\\config_test\\conf.d\\test1.com.conf":    "775ed01e78add3b934de529cec247b2a970558d7d1832a3e776e29a30bfc131a",
		"C:\\config_test\\conf.d\\test2.com.conf":    "bd31d5d2604233bbac22fb73e9125375c4bc8fe4c612c4611363b8f93413b2ea",
		"C:\\config_test\\mime.types":                "3c6049a805154dc0122c7264153036205c8f27f69699dc8ba129f212afb66d5a",
		"C:\\config_test\\nginx.conf":                "e2a36380e1591b13cca9d9eb5437bd1a2747901aa5f34caad39aa0960018d492",
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
						ContextPosPath: []int{8, 13, 4},
					},
					TargetContext: metav1.CloneConfigContextMeta{ConfigContextPos: metav1.ConfigContextPos{
						Config:         "C:\\config_test\\conf.d\\location2.conf",
						ContextPosPath: []int{4},
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
	webSrvOpts := metav1.WebServerOptions{}
	ofp := utilsV3.ConfigFingerprints{
		"C:\\config_test\\conf.d\\location.conf":     "539813c0f45630e9feba9a10c6494b4d912f0733847a4f17c650492709299c75",
		"C:\\config_test\\conf.d\\location2.conf":    "fc2a0cf89b11602e6dbfe0aa2c98cb69220485e77ef0e32a623043eb125f2114",
		"C:\\config_test\\conf.d\\server_test1.conf": "151c5dd9a238cdd69f4fc35d0564ab448c0e11530862457b41671fd41ddb9a0b",
		"C:\\config_test\\conf.d\\server_test2.conf": "24480b9ef0c9c86cb90896d7871e10bab94cc25fda496050d48533d6bf542f53",
		"C:\\config_test\\conf.d\\test1.com.conf":    "775ed01e78add3b934de529cec247b2a970558d7d1832a3e776e29a30bfc131a",
		"C:\\config_test\\conf.d\\test2.com.conf":    "bd31d5d2604233bbac22fb73e9125375c4bc8fe4c612c4611363b8f93413b2ea",
		"C:\\config_test\\mime.types":                "3c6049a805154dc0122c7264153036205c8f27f69699dc8ba129f212afb66d5a",
		"C:\\config_test\\nginx.conf":                "e2a36380e1591b13cca9d9eb5437bd1a2747901aa5f34caad39aa0960018d492",
	}
	difffp := utilsV3.ConfigFingerprints{
		"C:\\config_test\\conf.d\\location.conf":     "539813c0f45630e9feba9a10c6494b4d912f0733847a4f17c650492709299c75",
		"C:\\config_test\\conf.d\\location2.conf":    "fc2a0cf89b11602e6dbfe0aa2c98cb69220485e77ef0e32a623043eb125f2114",
		"C:\\config_test\\conf.d\\server_test1.conf": "151c5dd9a238cdd69f4fc35d0564ab448c0e11530862457b41671fd41ddb9a0b",
		"C:\\config_test\\conf.d\\server_test2.conf": "24480b9ef0c9c86cb90896d7871e10bab94cc25fda496050d48533d6bf542f53",
		"C:\\config_test\\conf.d\\test1.com.conf":    "775ed01e78add3b934de529cec247b2a970558d7d1832a3e776e29a30bfc131a",
		"C:\\config_test\\conf.d\\test2.com.conf":    "bd31d5d2604233bbac22fb73e9125375c4bc8fe4c612c4611363b8f93413b2ea",
		"C:\\config_test\\mime.types":                "3c6049a805154dc0122c7264153036205c8f27f69699dc8ba129f212afb66d5a",
		"C:\\config_test\\nginx.conf":                "1111111111111111111111111111111111111111111111111111111111111111",
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
	webSrvOpts := metav1.WebServerOptions{}
	ofp := utilsV3.ConfigFingerprints{
		"C:\\config_test\\conf.d\\location.conf":     "539813c0f45630e9feba9a10c6494b4d912f0733847a4f17c650492709299c75",
		"C:\\config_test\\conf.d\\location2.conf":    "fc2a0cf89b11602e6dbfe0aa2c98cb69220485e77ef0e32a623043eb125f2114",
		"C:\\config_test\\conf.d\\server_test1.conf": "151c5dd9a238cdd69f4fc35d0564ab448c0e11530862457b41671fd41ddb9a0b",
		"C:\\config_test\\conf.d\\server_test2.conf": "24480b9ef0c9c86cb90896d7871e10bab94cc25fda496050d48533d6bf542f53",
		"C:\\config_test\\conf.d\\test1.com.conf":    "775ed01e78add3b934de529cec247b2a970558d7d1832a3e776e29a30bfc131a",
		"C:\\config_test\\conf.d\\test2.com.conf":    "bd31d5d2604233bbac22fb73e9125375c4bc8fe4c612c4611363b8f93413b2ea",
		"C:\\config_test\\mime.types":                "3c6049a805154dc0122c7264153036205c8f27f69699dc8ba129f212afb66d5a",
		"C:\\config_test\\nginx.conf":                "e2a36380e1591b13cca9d9eb5437bd1a2747901aa5f34caad39aa0960018d492",
	}
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
						Config:         "C:\\config_test\\conf.d\\location2.conf",
						ContextPosPath: []int{0},
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
	webSrvOpts := metav1.WebServerOptions{}
	ofp := utilsV3.ConfigFingerprints{
		"C:\\config_test\\conf.d\\location.conf":     "539813c0f45630e9feba9a10c6494b4d912f0733847a4f17c650492709299c75",
		"C:\\config_test\\conf.d\\location2.conf":    "fc2a0cf89b11602e6dbfe0aa2c98cb69220485e77ef0e32a623043eb125f2114",
		"C:\\config_test\\conf.d\\server_test1.conf": "151c5dd9a238cdd69f4fc35d0564ab448c0e11530862457b41671fd41ddb9a0b",
		"C:\\config_test\\conf.d\\server_test2.conf": "24480b9ef0c9c86cb90896d7871e10bab94cc25fda496050d48533d6bf542f53",
		"C:\\config_test\\conf.d\\test1.com.conf":    "775ed01e78add3b934de529cec247b2a970558d7d1832a3e776e29a30bfc131a",
		"C:\\config_test\\conf.d\\test2.com.conf":    "bd31d5d2604233bbac22fb73e9125375c4bc8fe4c612c4611363b8f93413b2ea",
		"C:\\config_test\\mime.types":                "3c6049a805154dc0122c7264153036205c8f27f69699dc8ba129f212afb66d5a",
		"C:\\config_test\\nginx.conf":                "e2a36380e1591b13cca9d9eb5437bd1a2747901aa5f34caad39aa0960018d492",
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
						Config:         "C:\\config_test\\conf.d\\location.conf",
						ContextPosPath: []int{0},
					},
					TargetContext: metav1.CloneConfigContextMeta{ConfigContextPos: metav1.ConfigContextPos{
						Config:         "C:\\config_test\\conf.d\\location2.conf",
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
	webSrvOpts := metav1.WebServerOptions{}
	ofp := utilsV3.ConfigFingerprints{
		"C:\\config_test\\conf.d\\location.conf":     "539813c0f45630e9feba9a10c6494b4d912f0733847a4f17c650492709299c75",
		"C:\\config_test\\conf.d\\location2.conf":    "fc2a0cf89b11602e6dbfe0aa2c98cb69220485e77ef0e32a623043eb125f2114",
		"C:\\config_test\\conf.d\\server_test1.conf": "151c5dd9a238cdd69f4fc35d0564ab448c0e11530862457b41671fd41ddb9a0b",
		"C:\\config_test\\conf.d\\server_test2.conf": "24480b9ef0c9c86cb90896d7871e10bab94cc25fda496050d48533d6bf542f53",
		"C:\\config_test\\conf.d\\test1.com.conf":    "775ed01e78add3b934de529cec247b2a970558d7d1832a3e776e29a30bfc131a",
		"C:\\config_test\\conf.d\\test2.com.conf":    "bd31d5d2604233bbac22fb73e9125375c4bc8fe4c612c4611363b8f93413b2ea",
		"C:\\config_test\\mime.types":                "3c6049a805154dc0122c7264153036205c8f27f69699dc8ba129f212afb66d5a",
		"C:\\config_test\\nginx.conf":                "e2a36380e1591b13cca9d9eb5437bd1a2747901aa5f34caad39aa0960018d492",
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
						Config:         "C:\\config_test\\conf.d\\location2.conf",
						ContextPosPath: []int{0},
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
	webSrvOpts := metav1.WebServerOptions{}
	ofp := utilsV3.ConfigFingerprints{
		"C:\\config_test\\conf.d\\location.conf":     "539813c0f45630e9feba9a10c6494b4d912f0733847a4f17c650492709299c75",
		"C:\\config_test\\conf.d\\location2.conf":    "fc2a0cf89b11602e6dbfe0aa2c98cb69220485e77ef0e32a623043eb125f2114",
		"C:\\config_test\\conf.d\\server_test1.conf": "151c5dd9a238cdd69f4fc35d0564ab448c0e11530862457b41671fd41ddb9a0b",
		"C:\\config_test\\conf.d\\server_test2.conf": "24480b9ef0c9c86cb90896d7871e10bab94cc25fda496050d48533d6bf542f53",
		"C:\\config_test\\conf.d\\test1.com.conf":    "775ed01e78add3b934de529cec247b2a970558d7d1832a3e776e29a30bfc131a",
		"C:\\config_test\\conf.d\\test2.com.conf":    "bd31d5d2604233bbac22fb73e9125375c4bc8fe4c612c4611363b8f93413b2ea",
		"C:\\config_test\\mime.types":                "3c6049a805154dc0122c7264153036205c8f27f69699dc8ba129f212afb66d5a",
		"C:\\config_test\\nginx.conf":                "e2a36380e1591b13cca9d9eb5437bd1a2747901aa5f34caad39aa0960018d492",
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
	webSrvOpts := metav1.WebServerOptions{}
	ofp := utilsV3.ConfigFingerprints{
		"C:\\config_test\\conf.d\\location.conf":     "539813c0f45630e9feba9a10c6494b4d912f0733847a4f17c650492709299c75",
		"C:\\config_test\\conf.d\\location2.conf":    "fc2a0cf89b11602e6dbfe0aa2c98cb69220485e77ef0e32a623043eb125f2114",
		"C:\\config_test\\conf.d\\server_test1.conf": "151c5dd9a238cdd69f4fc35d0564ab448c0e11530862457b41671fd41ddb9a0b",
		"C:\\config_test\\conf.d\\server_test2.conf": "24480b9ef0c9c86cb90896d7871e10bab94cc25fda496050d48533d6bf542f53",
		"C:\\config_test\\conf.d\\test1.com.conf":    "775ed01e78add3b934de529cec247b2a970558d7d1832a3e776e29a30bfc131a",
		"C:\\config_test\\conf.d\\test2.com.conf":    "bd31d5d2604233bbac22fb73e9125375c4bc8fe4c612c4611363b8f93413b2ea",
		"C:\\config_test\\mime.types":                "3c6049a805154dc0122c7264153036205c8f27f69699dc8ba129f212afb66d5a",
		"C:\\config_test\\nginx.conf":                "e2a36380e1591b13cca9d9eb5437bd1a2747901aa5f34caad39aa0960018d492",
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
						ContextPosPath: []int{8, 13, 4},
					},
					TargetContext: metav1.CloneConfigContextMeta{ConfigContextPos: metav1.ConfigContextPos{
						Config:         "C:\\config_test\\conf.d\\location2.conf",
						ContextPosPath: []int{4},
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
	global.GVA_LOG = log.ZapLogger()
	webSrvOpts := metav1.WebServerOptions{}
	ofp := utilsV3.ConfigFingerprints{
		"C:\\config_test\\conf.d\\location.conf":     "539813c0f45630e9feba9a10c6494b4d912f0733847a4f17c650492709299c75",
		"C:\\config_test\\conf.d\\location2.conf":    "fc2a0cf89b11602e6dbfe0aa2c98cb69220485e77ef0e32a623043eb125f2114",
		"C:\\config_test\\conf.d\\server_test1.conf": "151c5dd9a238cdd69f4fc35d0564ab448c0e11530862457b41671fd41ddb9a0b",
		"C:\\config_test\\conf.d\\server_test2.conf": "24480b9ef0c9c86cb90896d7871e10bab94cc25fda496050d48533d6bf542f53",
		"C:\\config_test\\conf.d\\test1.com.conf":    "775ed01e78add3b934de529cec247b2a970558d7d1832a3e776e29a30bfc131a",
		"C:\\config_test\\conf.d\\test2.com.conf":    "bd31d5d2604233bbac22fb73e9125375c4bc8fe4c612c4611363b8f93413b2ea",
		"C:\\config_test\\mime.types":                "3c6049a805154dc0122c7264153036205c8f27f69699dc8ba129f212afb66d5a",
		"C:\\config_test\\nginx.conf":                "e2a36380e1591b13cca9d9eb5437bd1a2747901aa5f34caad39aa0960018d492",
	}
	difffp := utilsV3.ConfigFingerprints{
		"C:\\config_test\\conf.d\\location.conf":     "539813c0f45630e9feba9a10c6494b4d912f0733847a4f17c650492709299c75",
		"C:\\config_test\\conf.d\\location2.conf":    "fc2a0cf89b11602e6dbfe0aa2c98cb69220485e77ef0e32a623043eb125f2114",
		"C:\\config_test\\conf.d\\server_test1.conf": "151c5dd9a238cdd69f4fc35d0564ab448c0e11530862457b41671fd41ddb9a0b",
		"C:\\config_test\\conf.d\\server_test2.conf": "24480b9ef0c9c86cb90896d7871e10bab94cc25fda496050d48533d6bf542f53",
		"C:\\config_test\\conf.d\\test1.com.conf":    "775ed01e78add3b934de529cec247b2a970558d7d1832a3e776e29a30bfc131a",
		"C:\\config_test\\conf.d\\test2.com.conf":    "bd31d5d2604233bbac22fb73e9125375c4bc8fe4c612c4611363b8f93413b2ea",
		"C:\\config_test\\mime.types":                "3c6049a805154dc0122c7264153036205c8f27f69699dc8ba129f212afb66d5a",
		"C:\\config_test\\nginx.conf":                "1111111111111111111111111111111111111111111111111111111111111111",
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
				{"C:\\config_test\\nginx.conf", []int{8, 13, 0}},
				{"C:\\config_test\\nginx.conf", []int{8, 17, 0}},
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
				{"C:\\config_test\\nginx.conf", []int{8, 13, 1}},
				{"C:\\config_test\\nginx.conf", []int{8, 14, 1}},
				{"C:\\config_test\\nginx.conf", []int{8, 20, 1}},
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
			fields: fields{bm: mockBifrostsManager},
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
