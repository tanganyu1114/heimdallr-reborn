package web_server_bin_cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gin-vue-admin/global"
	svcv1 "gin-vue-admin/internal/hmdr_api/service/v1"
	"gin-vue-admin/internal/hmdr_api/service/v1/fake"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	"gin-vue-admin/model/response"
	log "github.com/ClessLi/component-base/pkg/log/v1"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestWebServerBinCMDController_Exec(t *testing.T) {
	webSrvOpts := metav1.WebServerOptions{
		ServerName: "test-bifrost",
	}
	global.GVA_LOG = log.ZapLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := svcv1.NewMockFactory(ctrl)
	svc.EXPECT().WebServerBinCMD().AnyTimes().Return(new(fake.Service).WebServerBinCMD())
	type fields struct {
		svc svcv1.Factory
	}
	type args struct {
		requestMeta any
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantRespBody []byte
		wantErr      bool
	}{
		{
			name:   "unknown request",
			fields: fields{svc: svc},
			args: args{requestMeta: metav1.WebServerBinCMDExecResponse{
				Successful: false,
				Stdout:     "errorrequest",
				Stderr:     "errorrequest",
			}},
			wantErr:      true,
			wantRespBody: []byte(`{"code":7,"data":{},"msg":"解析失败"}`),
		},
		{
			name:   "request with nil arg",
			fields: fields{svc: svc},
			args: args{requestMeta: metav1.WebServerBinCMDExecRequest{
				WebServerOptions: webSrvOpts,
				Args:             nil,
			}},
			wantRespBody: []byte(`{"code":0,"data":{"successful":true,"stdout":"nginx startup\n","stderr":""},"msg":"命令执行请求成功"}`),
			wantErr:      false,
		},
		{
			name:   "request with one arg",
			fields: fields{svc: svc},
			args: args{requestMeta: metav1.WebServerBinCMDExecRequest{
				WebServerOptions: webSrvOpts,
				Args:             []string{"-t"},
			}},
			wantRespBody: []byte(`{"code":0,"data":{"successful":false,"stdout":"","stderr":"config verify failure\n"},"msg":"命令执行请求成功"}`),
			wantErr:      false,
		},
		{
			name:   "request with two args",
			fields: fields{svc: svc},
			args: args{requestMeta: metav1.WebServerBinCMDExecRequest{
				WebServerOptions: webSrvOpts,
				Args:             []string{"-s", "reload"},
			}},
			wantRespBody: []byte(`{"code":0,"data":{"successful":true,"stdout":"nginx config reload successfully\n","stderr":""},"msg":"命令执行请求成功"}`),
		},
		{
			name:   "test unknown args",
			fields: fields{svc: svc},
			args: args{requestMeta: metav1.WebServerBinCMDExecRequest{
				WebServerOptions: webSrvOpts,
				Args:             []string{"-s", "unknown"},
			}},
			wantRespBody: []byte(fmt.Sprintf(`{"code":0,"data":{"successful":false,"stdout":"","stderr":"unknown argument:\n %v"},"msg":"命令执行请求成功"}`, []string{"-s", "unknown"})),
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WebServerBinCMDController{
				svc: tt.fields.svc,
			}
			reqBodyBytes, err := json.Marshal(tt.args.requestMeta)
			if err != nil {
				t.Fatal(err)
			}
			resp := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(resp)
			c.Request = httptest.NewRequest("POST", "/api/bin-cmd/exec", bytes.NewBuffer(reqBodyBytes))
			w.Exec(c)
			var respBody response.Response
			if err = json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
				t.Fatal(err)
			}
			if tt.wantErr != (resp.Code != 200 || respBody.Code != 0) {
				t.Errorf("Code: %d, Body: %s", resp.Code, resp.Body)
			}
			if got := resp.Body.Bytes(); !reflect.DeepEqual(got, tt.wantRespBody) {
				t.Errorf("Exec() got = %s, want %s", got, tt.wantRespBody)

			}
		})
	}
}
