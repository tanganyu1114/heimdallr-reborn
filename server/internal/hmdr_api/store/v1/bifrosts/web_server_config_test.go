package bifrosts

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	"gin-vue-admin/internal/pkg/bifrosts"
	"gin-vue-admin/internal/pkg/bifrosts/fake"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	bifrostclinetv1 "github.com/ClessLi/bifrost/pkg/client/bifrost/v1"
	"go.uber.org/mock/gomock"
	"reflect"
	"testing"

	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
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
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:   "normal test",
			fields: fields{bm: mockBifrostsManager},
			args:   args{opts: webSrvOpts},
			want: []string{
				`# user  nobody;`,
				`worker_processes 1;`,
				`# error_log  logs/error.log;`,
				`# error_log  logs/error.log  notice;`,
				`# error_log  logs/error.log  info;`,
				`# pid        logs/nginx.pid;`,
				`events {`,
				`    worker_connections 1024;`,
				`}`,
				`stream {`,
				`    upstream vvvvv1 {`,
				`        server test.1.cn:22 weight=5;`,
				`        server test.2.cn:33 weight=5;`,
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
				`    # access_log  logs/access.log  main;`,
				`    sendfile on;`,
				`    # tcp_nopush     on;`,
				`    # keepalive_timeout  0;`,
				`    keepalive_timeout 65;`,
				`    # gzip  on;`,
				`    # include <== conf.d\test*.com.conf`,
				`    server {    # test inline comment for client.UpdateConfig with resolv.V2 at 2021-01-26 17:13:59.357884 +0800 CST m=+0.142618001  # test inline comment for client.UpdateConfig with resolv.V2 at 2021-01-29 16:08:16.7355344 +0800 CST m=+0.083795701`,
				`        listen 80;`,
				`        server_name test1.com;`,
				`        # include <== conf.d\location.conf`,
				`        location /test {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # include location.conf;`,
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
				`        # include <== conf.d\location.conf`,
				`        location /test {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # include location.conf;`,
				`        upstream upstream.test.1 {`,
				`            server 10.1.1.1:443;`,
				`            server 10.1.1.2:443;`,
				`        }`,
				`        location /test_proxy {`,
				`            proxy_pass upstream.test.1;`,
				`        }`,
				`    }`,
				`    # include <== conf.d\server_test*.conf`,
				`    server {`,
				`        listen 80;`,
				`        server_name test1.com;`,
				`        # include <== conf.d\location*.conf`,
				`        location /test {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # include location.conf;`,
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
				`            proxy_pass upstream.test.2;`,
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
				`        # include <== conf.d\location*.conf`,
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
				`            proxy_pass upstream.test.2;`,
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
				`        location /test {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # include location.conf;`,
				`        upstream upstream.test.1 {`,
				`            server 10.1.1.1:443;`,
				`            server 10.1.1.2:443;`,
				`        }`,
				`        location /test_proxy {`,
				`            proxy_pass upstream.test.1;`,
				`        }`,
				`    }`,
				`    server {`,
				`        listen 8080;`,
				`        server_name test2.com;`,
				`        # include <== conf.d\location*.conf`,
				`        location /test {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # include location.conf;`,
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
				`            proxy_pass upstream.test.2;`,
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
				`        # access_log  logs/host.access.log  main;`,
				`        location / {`,
				`            root html;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # include <== conf.d\location.conf`,
				`        location /test {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # include location.conf;`,
				`        upstream upstream.test.1 {`,
				`            server 10.1.1.1:443;`,
				`            server 10.1.1.2:443;`,
				`        }`,
				`        location /test_proxy {`,
				`            proxy_pass upstream.test.1;`,
				`        }`,
				`        # error_page  404              /404.html;`,
				`        # redirect server error pages to the static page /50x.html`,
				`        #`,
				`        error_page 500 502 503 504  /50x.html;`,
				`        location = /50x.html {`,
				`            root html;`,
				`        }`,
				`        # proxy the PHP scripts to Apache listening on 127.0.0.1:80`,
				`        #`,
				`        # location ~ \.php$ {`,
				`        # proxy_pass   http://127.0.0.1;`,
				`        # }`,
				`        # pass the PHP scripts to FastCGI server listening on 127.0.0.1:9000`,
				`        #`,
				`        # location ~ \.php$ {`,
				`        # root           html;`,
				`        # fastcgi_pass   127.0.0.1:9000;`,
				`        # fastcgi_index  index.php;`,
				`        # fastcgi_param  SCRIPT_FILENAME  /scripts$fastcgi_script_name;`,
				`        # include        fastcgi_params;`,
				`        # }`,
				`        # deny access to .htaccess files, if Apache's document root`,
				`        # concurs with nginx's one`,
				`        #`,
				`        # location ~ /\.ht {`,
				`        # deny  all;`,
				`        # }`,
				`    }`,
				`    server {`,
				`        listen 990 ssl;`,
				`        server_name localhost;`,
				`        # charset koi8-r;`,
				`        # access_log  logs/host.access.log  main;`,
				`        location / {`,
				`            root html;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # include <== conf.d\location.conf`,
				`        location /test {`,
				`            root html/test;`,
				`            index index.html index.htm;`,
				`        }`,
				`        # include location.conf;`,
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
				`    # listen       8000;`,
				`    # listen       somename:8080;`,
				`    # server_name  somename  alias  another.alias;`,
				`    # location / {`,
				`    # root   html;`,
				`    # index  index.html index.htm;`,
				`    # }`,
				`    # }`,
				`    # HTTPS server`,
				`    #`,
				`    # server {`,
				`    # listen       443 ssl;`,
				`    # server_name  localhost;`,
				`    # ssl_certificate      cert.pem;`,
				`    # ssl_certificate_key  cert.key;`,
				`    # ssl_session_cache    shared:SSL:1m;`,
				`    # ssl_session_timeout  5m;`,
				`    # ssl_ciphers  HIGH:!aNULL:!MD5;`,
				`    # ssl_prefer_server_ciphers  on;`,
				`    # location / {`,
				`    # root   html;`,
				`    # index  index.html index.htm;`,
				`    # }`,
				`    # }`,
				`}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigStore{
				bm: tt.fields.bm,
			}
			got, err := w.GetConfig(tt.args.ctx, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfig() got = %v, want %v", got, tt.want)
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
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigStore{
				bm: tt.fields.bm,
			}
			if err := w.InsertWithClone(tt.args.in0, tt.args.opts, tt.args.ctxmeta); (err != nil) != tt.wantErr {
				t.Errorf("InsertWithClone() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigStore_InsertWithNew(t *testing.T) {
	webSrvOpts := metav1.WebServerOptions{}
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigStore{
				bm: tt.fields.bm,
			}
			if err := w.InsertWithNew(tt.args.in0, tt.args.opts, tt.args.ctxmeta); (err != nil) != tt.wantErr {
				t.Errorf("InsertWithNew() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigStore_ModifyWithClone(t *testing.T) {
	webSrvOpts := metav1.WebServerOptions{}
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
			if err := w.ModifyWithClone(tt.args.in0, tt.args.opts, tt.args.ctxmeta); (err != nil) != tt.wantErr {
				t.Errorf("ModifyWithClone() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigStore_ModifyWithNew(t *testing.T) {
	webSrvOpts := metav1.WebServerOptions{}
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
			if err := w.ModifyWithNew(tt.args.in0, tt.args.opts, tt.args.ctxmeta); (err != nil) != tt.wantErr {
				t.Errorf("ModifyWithNew() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigStore_Remove(t *testing.T) {
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
			if err := w.Remove(tt.args.ctx, tt.args.opts, tt.args.pos); (err != nil) != tt.wantErr {
				t.Errorf("Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigStore_getConfig(t *testing.T) {
	type fields struct {
		bm bifrosts.Manager
	}
	type args struct {
		opts metav1.WebServerOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    configuration.NginxConfig
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigStore{
				bm: tt.fields.bm,
			}
			got, err := w.getConfig(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("getConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webServerConfigStore_parseContextPos(t *testing.T) {
	type fields struct {
		bm bifrosts.Manager
	}
	type args struct {
		nginxconfig configuration.NginxConfig
		pos         metav1.ConfigContextPos
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    nginx_context.Pos
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerConfigStore{
				bm: tt.fields.bm,
			}
			got, err := w.parseContextPos(tt.args.nginxconfig, tt.args.pos)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseContextPos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseContextPos() got = %v, want %v", got, tt.want)
			}
		})
	}
}
