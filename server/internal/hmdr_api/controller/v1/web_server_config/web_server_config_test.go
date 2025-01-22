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

func TestWebServerConfigController_ChangeContextEnabledState(t *testing.T) {
	global.GVA_LOG = log.ZapLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := svcv1.NewMockFactory(ctrl)
	svc.EXPECT().WebServerConfigs().AnyTimes().Return(new(svcfake.WebServerConfigService))
	serverOpts := metav1.WebServerOptions{
		GroupID:    0,
		HostID:     0,
		ServerName: "test-bifrost",
	}
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
			name:   "enable context",
			fields: fields{svc: svc},
			args: args{
				requestMeta: metav1.WebServerConfigContextUpdateOptions[metav1.ConfigContextEnabledStateMeta]{
					WebServerOptions: serverOpts,
					TargetConfigContextOptions: metav1.TargetConfigContextOptions[metav1.ConfigContextEnabledStateMeta]{
						Position: metav1.ConfigContextPos{
							Config:         "C:\\config_test\\nginx.conf",
							ContextPosPath: []int{0},
						},
						TargetContext: metav1.ConfigContextEnabledStateMeta{Enabled: true},
					},
				},
			},
			wantRespBody: []byte(`{"code":0,"data":{},"msg":"修改启用状态成功"}`),
			wantErr:      false,
		},
		{
			name:   "disable context",
			fields: fields{svc: svc},
			args: args{
				requestMeta: metav1.WebServerConfigContextUpdateOptions[metav1.ConfigContextEnabledStateMeta]{
					WebServerOptions: serverOpts,
					TargetConfigContextOptions: metav1.TargetConfigContextOptions[metav1.ConfigContextEnabledStateMeta]{
						Position: metav1.ConfigContextPos{
							Config:         "conf.d\\location.conf",
							ContextPosPath: nil,
						},
						//TargetContext: metav1.ConfigContextEnabledStateMeta{Enabled: false},
					},
				},
			},
			wantRespBody: []byte(`{"code":0,"data":{},"msg":"修改启用状态成功"}`),
			wantErr:      false,
		},
		{
			name:   "wrong position",
			fields: fields{svc: svc},
			args: args{
				requestMeta: metav1.WebServerConfigContextUpdateOptions[metav1.ConfigContextEnabledStateMeta]{
					WebServerOptions: serverOpts,
					TargetConfigContextOptions: metav1.TargetConfigContextOptions[metav1.ConfigContextEnabledStateMeta]{
						Position: metav1.ConfigContextPos{
							Config:         "conf.d\\location.conf",
							ContextPosPath: []int{1, 2, 3},
						},
						TargetContext: metav1.ConfigContextEnabledStateMeta{Enabled: true},
					},
				},
			},
			wantRespBody: []byte(`{"code":7,"data":{},"msg":"修改启用状态失败"}`),
			wantErr:      true,
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
			c.Request = httptest.NewRequest("POST", "/api/conf/change-ctx-enabled-state", bytes.NewBuffer(reqBodyBytes))
			w.ChangeContextEnabledState(c)
			var respBody response.Response
			if err = json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
				t.Fatal(err)
			}
			if tt.wantErr != (resp.Code != 200 || respBody.Code != 0) {
				t.Errorf("Code: %d, Body: %s", resp.Code, resp.Body)
			}
			if got := resp.Body.Bytes(); !reflect.DeepEqual(got, tt.wantRespBody) {
				t.Errorf("ChangeContextEnabledState() got = %s, want %s", got, tt.wantRespBody)
			}
		})
	}
}

func TestWebServerConfigController_GetConfigTextLines(t *testing.T) {
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
			name:   "normal test",
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
			w.GetConfigTextLines(tt.args.c)
			t.Logf("Code: %d, Body: %s", tt.args.resp.Code, tt.args.resp.Body)
		})
	}
}

func TestWebServerConfigController_GetContextTextLines(t *testing.T) {
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
	targetCtxOpts := metav1.WebServerConfigTargetContextOptions{
		WebServerOptions: serverOpts,
		ConfigContextPos: metav1.ConfigContextPos{
			Config:         "C:\\config_test\\conf.d\\location.conf",
			ContextPosPath: []int{0},
		},
	}
	b, err := json.Marshal(targetCtxOpts)
	if err != nil {
		t.Fatal(err)
	}
	getBody := bytes.NewBuffer(b)
	getRequest := httptest.NewRequest("GET", "/api/conf/get-context-text", getBody)
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
			name:   "normal test",
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
			w.GetContextTextLines(tt.args.c)
			t.Logf("Code: %d, Body: %s", tt.args.resp.Code, tt.args.resp.Body)
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
	getRequest := httptest.NewRequest("GET", "/api/conf/get-conf-struct", getBody)
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
			name:   "normal test",
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

func TestWebServerConfigController_GetIncludedConfigs(t *testing.T) {
	global.GVA_LOG = log.ZapLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := svcv1.NewMockFactory(ctrl)
	svc.EXPECT().WebServerConfigs().AnyTimes().Return(new(svcfake.WebServerConfigService))
	serverOpts := metav1.WebServerOptions{
		GroupID:    0,
		HostID:     0,
		ServerName: "test-bifrost",
	}
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
			name:   "normal test",
			fields: fields{svc: svc},
			args: args{
				requestMeta: metav1.WebServerConfigTargetContextOptions{
					WebServerOptions: serverOpts,
					ConfigContextPos: metav1.ConfigContextPos{
						Config:         "C:\\config_test\\conf.d\\server_test1.conf",
						ContextPosPath: []int{0, 2},
					},
				},
			},
			wantRespBody: []byte(`{"code":0,"data":["C:\\config_test\\conf.d\\location.conf","C:\\config_test\\conf.d\\location2.conf"],"msg":"获取包含的配置文件成功"}`),
			wantErr:      false,
		},
		{
			name:   "wrong position",
			fields: fields{svc: svc},
			args: args{
				requestMeta: metav1.WebServerConfigTargetContextOptions{
					WebServerOptions: serverOpts,
					ConfigContextPos: metav1.ConfigContextPos{
						Config:         "C:\\config_test\\conf.d\\server_test1.con",
						ContextPosPath: nil,
					},
				},
			},
			wantRespBody: []byte(`{"code":7,"data":{},"msg":"获取包含的配置文件失败"}`),
			wantErr:      true,
		},
		{
			name:   "the target is not an `include` context",
			fields: fields{svc: svc},
			args: args{
				requestMeta: metav1.WebServerConfigTargetContextOptions{
					WebServerOptions: serverOpts,
					ConfigContextPos: metav1.ConfigContextPos{
						Config:         "C:\\config_test\\conf.d\\server_test1.conf",
						ContextPosPath: []int{0, 1},
					},
				},
			},
			wantRespBody: []byte(`{"code":7,"data":{},"msg":"获取包含的配置文件失败"}`),
			wantErr:      true,
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
			c.Request = httptest.NewRequest("POST", "/api/conf/get-includes", bytes.NewBuffer(reqBodyBytes))
			w.GetIncludedConfigs(c)
			var respBody response.Response
			if err = json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
				t.Fatal(err)
			}
			if tt.wantErr != (resp.Code != 200 || respBody.Code != 0) {
				t.Errorf("Code: %d, Body: %s", resp.Code, resp.Body)
			}
			if got := resp.Body.Bytes(); !reflect.DeepEqual(got, tt.wantRespBody) {
				t.Errorf("GetIncludedConfigs() got = %s, want %s", got, tt.wantRespBody)
			}
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
		name         string
		fields       fields
		args         args
		wantRespBody []byte
		wantErr      bool
	}{
		{
			name:   "valid request body",
			fields: fields{svc: svc},
			args: args{
				requestMeta: validMeta,
			},
			wantRespBody: []byte(`{"code":0,"data":{},"msg":"新增成功"}`),
			wantErr:      false,
		},
		{
			name:   "invalid request body",
			fields: fields{svc: svc},
			args: args{
				requestMeta: invalidMeta,
			},
			wantRespBody: []byte(`{"code":7,"data":{},"msg":"解析失败"}`),
			wantErr:      true,
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
			if got := resp.Body.Bytes(); !reflect.DeepEqual(got, tt.wantRespBody) {
				t.Errorf("InsertWithClone() got = %s, want %s", got, tt.wantRespBody)
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
			name:   "valid request body",
			fields: fields{svc: svc},
			args: args{
				requestMeta: validMeta,
			},
			wantRespBody: []byte(`{"code":0,"data":{},"msg":"新增成功"}`),
			wantErr:      false,
		},
		{
			name:   "invalid request body",
			fields: fields{svc: svc},
			args: args{
				requestMeta: invalidMeta,
			},
			wantRespBody: []byte(`{"code":7,"data":{},"msg":"解析失败"}`),
			wantErr:      true,
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
			if got := resp.Body.Bytes(); !reflect.DeepEqual(got, tt.wantRespBody) {
				t.Errorf("InsertWithNew() got = %s, want %s", got, tt.wantRespBody)
			}
		})
	}
}

func TestWebServerConfigController_ModifyContextValue(t *testing.T) {
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
				ContextPosPath: []int{0},
			},
			TargetContext: metav1.NewConfigContextMeta{
				ContextType:  "location",
				ContextValue: "~ /normal-test",
			},
		},
	}
	unmatchedTypeMeta := metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta]{
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
		{
			name:   "unmatched context type",
			fields: fields{svc: svc},
			args: args{
				requestMeta: unmatchedTypeMeta,
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
			c.Request = httptest.NewRequest("POST", "/api/conf/modify-ctx-value", bytes.NewBuffer(reqBodyBytes))
			w.ModifyContextValue(c)
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
	validMeta := metav1.WebServerConfigTargetContextOptions{
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
	invalidMeta := metav1.WebServerConfigTargetContextOptions{}
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

func TestWebServerConfigController_Move(t *testing.T) {
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
			c.Request = httptest.NewRequest("POST", "/api/conf/move-ctx", bytes.NewBuffer(reqBodyBytes))
			w.Move(c)
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
