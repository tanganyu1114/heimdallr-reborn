package web_server_config

import (
	"bytes"
	"encoding/json"
	"gin-vue-admin/global"
	svcv1 "gin-vue-admin/internal/hmdr_api/service/v1"
	svcfake "gin-vue-admin/internal/hmdr_api/service/v1/fake"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	"gin-vue-admin/model/response"
	log "github.com/ClessLi/component-base/pkg/log/v1"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewController(t *testing.T) {
	svc := new(svcfake.Service)
	type args struct {
		service svcv1.Factory
	}
	tests := []struct {
		name string
		args args
		want *WebServerConfigController
	}{
		{
			name: "normal test",
			args: args{service: svc},
			want: &WebServerConfigController{svc: svc},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewController(tt.args.service); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewController() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebServerConfigController_GetConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := svcv1.NewMockFactory(ctrl)
	svc.EXPECT().WebServerConfigs().AnyTimes().Return(new(svcfake.WebServerConfigService))
	getResponse := httptest.NewRecorder()
	serverOpts := metav1.WebServerOptions{
		GroupID:    0,
		HostID:     0,
		ServerName: "test-bifrost",
	}
	b, err := json.Marshal(serverOpts)
	if err != nil {
		t.Fatal(err)
	}
	getBody := bytes.NewBuffer(b)
	getRequest := httptest.NewRequest("GET", "/api/conf/getConfInfo", getBody)
	getRequestCtx, _ := gin.CreateTestContext(getResponse)
	getRequestCtx.Request = getRequest
	type fields struct {
		svc svcv1.Factory
	}
	type args struct {
		c    *gin.Context
		resp *httptest.ResponseRecorder
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "http method error",
			fields: fields{svc: svc},
			args: args{
				c:    getRequestCtx,
				resp: getResponse,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WebServerConfigController{
				svc: tt.fields.svc,
			}
			w.GetConfig(tt.args.c)
			t.Logf("Code: %d, Body: %s", tt.args.resp.Code, tt.args.resp.Body)
		})
	}
}

func TestWebServerConfigController_GetOptions(t *testing.T) {
	type fields struct {
		svc svcv1.Factory
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WebServerConfigController{
				svc: tt.fields.svc,
			}
			w.GetOptions(tt.args.c)
		})
	}
}

func TestWebServerConfigController_InsertWithClone(t *testing.T) {
	global.GVA_LOG = log.ZapLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := svcv1.NewMockFactory(ctrl)
	svc.EXPECT().WebServerConfigs().AnyTimes().Return(new(svcfake.WebServerConfigService))
	validMeta := metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta]{
		WebServerOptions: metav1.WebServerOptions{
			GroupID:    0,
			HostID:     0,
			ServerName: "test-bifrost",
		},
		TargetConfigContextOptions: metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]{
			Position: metav1.ConfigContextPos{
				Config:         "C:\\config_test\\nginx.conf",
				ContextPosPath: []int{8, 13, 4},
			},
			TargetContext: metav1.CloneConfigContextMeta{ConfigContextPos: metav1.ConfigContextPos{
				Config:         "C:\\config_test\\conf.d\\location2.conf",
				ContextPosPath: []int{4},
			}},
		},
	}
	invalidMeta := metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta]{
		WebServerOptions: metav1.WebServerOptions{
			GroupID:    0,
			HostID:     0,
			ServerName: "test-bifrost",
		},
		TargetConfigContextOptions: metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]{
			Position: metav1.ConfigContextPos{
				Config:         "C:\\config_test\\conf.d\\location2.conf",
				ContextPosPath: []int{2},
			},
			TargetContext: metav1.NewConfigContextMeta{
				ContextType:  "location",
				ContextValue: "~ /normal-test",
			},
		},
	}
	type fields struct {
		svc svcv1.Factory
	}
	type args struct {
		requestMeta any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "valid request body",
			fields: fields{svc: svc},
			args: args{
				requestMeta: validMeta,
			},
			wantErr: false,
		},
		{
			name:   "invalid request body",
			fields: fields{svc: svc},
			args: args{
				requestMeta: invalidMeta,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WebServerConfigController{
				svc: tt.fields.svc,
			}
			reqBodyBytes, err := json.Marshal(tt.args.requestMeta)
			if err != nil {
				t.Fatal(err)
			}
			resp := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(resp)
			c.Request = httptest.NewRequest("POST", "/api/conf/insert-clone-ctx", bytes.NewBuffer(reqBodyBytes))
			w.InsertWithClone(c)
			var respBody response.Response
			if err = json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
				t.Fatal(err)
			}
			if tt.wantErr != (resp.Code != 200 || respBody.Code != 0) {
				t.Errorf("Code: %d, Body: %s", resp.Code, resp.Body)
			}
		})
	}
}

func TestWebServerConfigController_InsertWithNew(t *testing.T) {
	global.GVA_LOG = log.ZapLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := svcv1.NewMockFactory(ctrl)
	svc.EXPECT().WebServerConfigs().AnyTimes().Return(new(svcfake.WebServerConfigService))
	invalidMeta := metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta]{
		WebServerOptions: metav1.WebServerOptions{
			GroupID:    0,
			HostID:     0,
			ServerName: "test-bifrost",
		},
		TargetConfigContextOptions: metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]{
			Position: metav1.ConfigContextPos{
				Config:         "C:\\config_test\\nginx.conf",
				ContextPosPath: []int{8, 13, 4},
			},
			TargetContext: metav1.CloneConfigContextMeta{ConfigContextPos: metav1.ConfigContextPos{
				Config:         "C:\\config_test\\conf.d\\location2.conf",
				ContextPosPath: []int{4},
			}},
		},
	}
	validMeta := metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta]{
		WebServerOptions: metav1.WebServerOptions{
			GroupID:    0,
			HostID:     0,
			ServerName: "test-bifrost",
		},
		TargetConfigContextOptions: metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]{
			Position: metav1.ConfigContextPos{
				Config:         "C:\\config_test\\conf.d\\location2.conf",
				ContextPosPath: []int{2},
			},
			TargetContext: metav1.NewConfigContextMeta{
				ContextType:  "location",
				ContextValue: "~ /normal-test",
			},
		},
	}
	type fields struct {
		svc svcv1.Factory
	}
	type args struct {
		requestMeta any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "valid request body",
			fields: fields{svc: svc},
			args: args{
				requestMeta: validMeta,
			},
			wantErr: false,
		},
		{
			name:   "invalid request body",
			fields: fields{svc: svc},
			args: args{
				requestMeta: invalidMeta,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WebServerConfigController{
				svc: tt.fields.svc,
			}
			reqBodyBytes, err := json.Marshal(tt.args.requestMeta)
			if err != nil {
				t.Fatal(err)
			}
			resp := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(resp)
			c.Request = httptest.NewRequest("POST", "/api/conf/insert-new-ctx", bytes.NewBuffer(reqBodyBytes))
			w.InsertWithNew(c)
			var respBody response.Response
			if err = json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
				t.Fatal(err)
			}
			if tt.wantErr != (resp.Code != 200 || respBody.Code != 0) {
				t.Errorf("Code: %d, Body: %s", resp.Code, resp.Body)
			}
		})
	}
}

func TestWebServerConfigController_ModifyWithClone(t *testing.T) {
	global.GVA_LOG = log.ZapLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := svcv1.NewMockFactory(ctrl)
	svc.EXPECT().WebServerConfigs().AnyTimes().Return(new(svcfake.WebServerConfigService))
	validMeta := metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta]{
		WebServerOptions: metav1.WebServerOptions{
			GroupID:    0,
			HostID:     0,
			ServerName: "test-bifrost",
		},
		TargetConfigContextOptions: metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]{
			Position: metav1.ConfigContextPos{
				Config:         "C:\\config_test\\nginx.conf",
				ContextPosPath: []int{8, 13, 4},
			},
			TargetContext: metav1.CloneConfigContextMeta{ConfigContextPos: metav1.ConfigContextPos{
				Config:         "C:\\config_test\\conf.d\\location2.conf",
				ContextPosPath: []int{4},
			}},
		},
	}
	invalidMeta := metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta]{
		WebServerOptions: metav1.WebServerOptions{
			GroupID:    0,
			HostID:     0,
			ServerName: "test-bifrost",
		},
		TargetConfigContextOptions: metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]{
			Position: metav1.ConfigContextPos{
				Config:         "C:\\config_test\\conf.d\\location2.conf",
				ContextPosPath: []int{2},
			},
			TargetContext: metav1.NewConfigContextMeta{
				ContextType:  "location",
				ContextValue: "~ /normal-test",
			},
		},
	}
	type fields struct {
		svc svcv1.Factory
	}
	type args struct {
		requestMeta any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "valid request body",
			fields: fields{svc: svc},
			args: args{
				requestMeta: validMeta,
			},
			wantErr: false,
		},
		{
			name:   "invalid request body",
			fields: fields{svc: svc},
			args: args{
				requestMeta: invalidMeta,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WebServerConfigController{
				svc: tt.fields.svc,
			}
			reqBodyBytes, err := json.Marshal(tt.args.requestMeta)
			if err != nil {
				t.Fatal(err)
			}
			resp := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(resp)
			c.Request = httptest.NewRequest("POST", "/api/conf/modify-clone-ctx", bytes.NewBuffer(reqBodyBytes))
			w.ModifyWithClone(c)
			var respBody response.Response
			if err = json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
				t.Fatal(err)
			}
			if tt.wantErr != (resp.Code != 200 || respBody.Code != 0) {
				t.Errorf("Code: %d, Body: %s", resp.Code, resp.Body)
			}
		})
	}
}

func TestWebServerConfigController_ModifyWithNew(t *testing.T) {
	global.GVA_LOG = log.ZapLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := svcv1.NewMockFactory(ctrl)
	svc.EXPECT().WebServerConfigs().AnyTimes().Return(new(svcfake.WebServerConfigService))
	invalidMeta := metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta]{
		WebServerOptions: metav1.WebServerOptions{
			GroupID:    0,
			HostID:     0,
			ServerName: "test-bifrost",
		},
		TargetConfigContextOptions: metav1.TargetConfigContextOptions[metav1.CloneConfigContextMeta]{
			Position: metav1.ConfigContextPos{
				Config:         "C:\\config_test\\nginx.conf",
				ContextPosPath: []int{8, 13, 4},
			},
			TargetContext: metav1.CloneConfigContextMeta{ConfigContextPos: metav1.ConfigContextPos{
				Config:         "C:\\config_test\\conf.d\\location2.conf",
				ContextPosPath: []int{4},
			}},
		},
	}
	validMeta := metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta]{
		WebServerOptions: metav1.WebServerOptions{
			GroupID:    0,
			HostID:     0,
			ServerName: "test-bifrost",
		},
		TargetConfigContextOptions: metav1.TargetConfigContextOptions[metav1.NewConfigContextMeta]{
			Position: metav1.ConfigContextPos{
				Config:         "C:\\config_test\\conf.d\\location2.conf",
				ContextPosPath: []int{2},
			},
			TargetContext: metav1.NewConfigContextMeta{
				ContextType:  "location",
				ContextValue: "~ /normal-test",
			},
		},
	}
	type fields struct {
		svc svcv1.Factory
	}
	type args struct {
		requestMeta any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "valid request body",
			fields: fields{svc: svc},
			args: args{
				requestMeta: validMeta,
			},
			wantErr: false,
		},
		{
			name:   "invalid request body",
			fields: fields{svc: svc},
			args: args{
				requestMeta: invalidMeta,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WebServerConfigController{
				svc: tt.fields.svc,
			}
			reqBodyBytes, err := json.Marshal(tt.args.requestMeta)
			if err != nil {
				t.Fatal(err)
			}
			resp := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(resp)
			c.Request = httptest.NewRequest("POST", "/api/conf/modify-new-ctx", bytes.NewBuffer(reqBodyBytes))
			w.ModifyWithNew(c)
			var respBody response.Response
			if err = json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
				t.Fatal(err)
			}
			if tt.wantErr != (resp.Code != 200 || respBody.Code != 0) {
				t.Errorf("Code: %d, Body: %s", resp.Code, resp.Body)
			}
		})
	}
}

func TestWebServerConfigController_Remove(t *testing.T) {
	global.GVA_LOG = log.ZapLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := svcv1.NewMockFactory(ctrl)
	svc.EXPECT().WebServerConfigs().AnyTimes().Return(new(svcfake.WebServerConfigService))
	validMeta := metav1.WebServerConfigContextDeleteOptions{
		WebServerOptions: metav1.WebServerOptions{
			GroupID:    0,
			HostID:     0,
			ServerName: "test-bifrost",
		},
		ConfigContextPos: metav1.ConfigContextPos{
			Config:         "C:\\config_test\\nginx.conf",
			ContextPosPath: []int{8, 13, 4},
		},
	}
	invalidMeta := metav1.WebServerConfigContextDeleteOptions{}
	type fields struct {
		svc svcv1.Factory
	}
	type args struct {
		requestMeta any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "valid request body",
			fields: fields{svc: svc},
			args: args{
				requestMeta: validMeta,
			},
			wantErr: false,
		},
		{
			name:   "invalid request body",
			fields: fields{svc: svc},
			args: args{
				requestMeta: invalidMeta,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WebServerConfigController{
				svc: tt.fields.svc,
			}
			reqBodyBytes, err := json.Marshal(tt.args.requestMeta)
			if err != nil {
				t.Fatal(err)
			}
			resp := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(resp)
			c.Request = httptest.NewRequest("DELETE", "/api/conf/remove-ctx", bytes.NewBuffer(reqBodyBytes))
			w.Remove(c)
			var respBody response.Response
			if err = json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
				t.Fatal(err)
			}
			if tt.wantErr != (resp.Code != 200 || respBody.Code != 0) {
				t.Errorf("Code: %d, Body: %s", resp.Code, resp.Body)
			}
		})
	}
}
