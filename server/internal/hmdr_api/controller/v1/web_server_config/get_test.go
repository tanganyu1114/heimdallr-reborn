package web_server_config

import (
	"bytes"
	"encoding/json"
	"github.com/tanganyu1114/heimdallr-reborn/global"
	svcv1 "github.com/tanganyu1114/heimdallr-reborn/internal/hmdr_api/service/v1"
	"github.com/tanganyu1114/heimdallr-reborn/internal/pkg/bifrosts/fake"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/internal/pkg/meta/v1"
	"github.com/tanganyu1114/heimdallr-reborn/model/response"
	"net/http/httptest"
	"testing"

	nginx_context "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context/local"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context_type"
	logV1 "github.com/ClessLi/component-base/pkg/log/v1"
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/errors"
	"go.uber.org/mock/gomock"
)

func TestWebServerConfigController_GetConfig(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := svcv1.NewMockFactory(ctrl)
	mockWebServerConfigSrv := svcv1.NewMockWebServerConfigSrv(ctrl)
	svc.EXPECT().WebServerConfigs().AnyTimes().Return(mockWebServerConfigSrv)

	validMeta := metav1.WebServerOptions{
		GroupID:    0,
		HostID:     0,
		ServerName: "test-bifrost",
	}

	invalidMeta := metav1.WebServerOptions{
		GroupID:    0,
		HostID:     0,
		ServerName: "",
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
			name:   "get config returns error",
			fields: fields{svc: svc},
			args: args{
				requestMeta: validMeta,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "get config returns error" {
				mockWebServerConfigSrv.EXPECT().GetConfig(gomock.Any(), gomock.Any()).Return(metav1.WebServerConfig{}, errors.New("get config failed"))
			} else if tt.name == "valid request body" {
				mockWebServerConfigSrv.EXPECT().GetConfig(gomock.Any(), gomock.Any()).Return(metav1.WebServerConfig{}, nil)
			}

			w := &WebServerConfigController{
				svc: tt.fields.svc,
			}
			reqBodyBytes, err := json.Marshal(tt.args.requestMeta)
			if err != nil {
				t.Fatal(err)
			}
			resp := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(resp)
			c.Request = httptest.NewRequest("POST", "/api/conf/get", bytes.NewBuffer(reqBodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")
			w.GetConfig(c)
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

func TestWebServerConfigController_GetConfigTextLines(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := svcv1.NewMockFactory(ctrl)
	mockWebServerConfigSrv := svcv1.NewMockWebServerConfigSrv(ctrl)
	svc.EXPECT().WebServerConfigs().AnyTimes().Return(mockWebServerConfigSrv)

	validMeta := metav1.WebServerOptions{
		GroupID:    0,
		HostID:     0,
		ServerName: "test-bifrost",
	}

	invalidMeta := metav1.WebServerOptions{
		GroupID:    0,
		HostID:     0,
		ServerName: "",
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
			name:   "get config returns error",
			fields: fields{svc: svc},
			args: args{
				requestMeta: validMeta,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "get config returns error" {
				mockWebServerConfigSrv.EXPECT().GetConfig(gomock.Any(), gomock.Any()).Return(metav1.WebServerConfig{}, errors.New("get config failed"))
			} else if tt.name == "valid request body" {
				fakeSvc := fake.WebServerConfigService{}
				fakeConfig, _, _ := fakeSvc.Get("test-bifrost")
				mockWebServerConfigSrv.EXPECT().GetConfig(gomock.Any(), gomock.Any()).Return(metav1.WebServerConfig{
					Config: fakeConfig.Main(),
				}, nil)
			}

			w := &WebServerConfigController{
				svc: tt.fields.svc,
			}
			reqBodyBytes, err := json.Marshal(tt.args.requestMeta)
			if err != nil {
				t.Fatal(err)
			}
			resp := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(resp)
			c.Request = httptest.NewRequest("POST", "/api/conf/get-text-lines", bytes.NewBuffer(reqBodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")
			w.GetConfigTextLines(c)
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

func TestWebServerConfigController_GetContextTextLines(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := svcv1.NewMockFactory(ctrl)
	mockWebServerConfigSrv := svcv1.NewMockWebServerConfigSrv(ctrl)
	svc.EXPECT().WebServerConfigs().AnyTimes().Return(mockWebServerConfigSrv)

	validMeta := metav1.WebServerConfigTargetContextOptions{
		WebServerOptions: metav1.WebServerOptions{
			GroupID:    0,
			HostID:     0,
			ServerName: "test-bifrost",
		},
		ConfigContextPos: metav1.ConfigContextPos{
			Config:         "E:\\config_test\\nginx.conf",
			ContextPosPath: []int{0},
		},
		OriginalFingerprints: map[string]string{
			"E:\\config_test\\nginx.conf": "e2a36380e1591b13cca9d9eb5437bd1a2747901aa5f34caad39aa0960018d492",
		},
	}

	invalidMeta := metav1.WebServerConfigTargetContextOptions{
		WebServerOptions: metav1.WebServerOptions{
			GroupID:    0,
			HostID:     0,
			ServerName: "",
		},
		ConfigContextPos:     metav1.ConfigContextPos{},
		OriginalFingerprints: nil,
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
			name:   "get context returns error",
			fields: fields{svc: svc},
			args: args{
				requestMeta: validMeta,
			},
			wantErr: true,
		},
		{
			name:   "config lines returns error",
			fields: fields{svc: svc},
			args: args{
				requestMeta: validMeta,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "get context returns error" {
				mockWebServerConfigSrv.EXPECT().GetContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("get context failed"))
			} else if tt.name == "config lines returns error" {
				mockWebServerConfigSrv.EXPECT().GetContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nginx_context.ErrContext(errors.New("config lines failed")), nil)
			} else if tt.name == "valid request body" {
				mockWebServerConfigSrv.EXPECT().GetContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(local.NewContext(context_type.TypeDirHTTPProxyPass, "https://baidu.com"), nil)
			}

			w := &WebServerConfigController{
				svc: tt.fields.svc,
			}
			reqBodyBytes, err := json.Marshal(tt.args.requestMeta)
			if err != nil {
				t.Fatal(err)
			}
			resp := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(resp)
			c.Request = httptest.NewRequest("POST", "/api/conf/get-context-text-lines", bytes.NewBuffer(reqBodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")
			w.GetContextTextLines(c)
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

func TestWebServerConfigController_GetIncludedConfigs(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := svcv1.NewMockFactory(ctrl)
	mockWebServerConfigSrv := svcv1.NewMockWebServerConfigSrv(ctrl)
	svc.EXPECT().WebServerConfigs().AnyTimes().Return(mockWebServerConfigSrv)

	validMeta := metav1.WebServerConfigTargetContextOptions{
		WebServerOptions: metav1.WebServerOptions{
			GroupID:    0,
			HostID:     0,
			ServerName: "test-bifrost",
		},
		ConfigContextPos: metav1.ConfigContextPos{
			Config:         "E:\\config_test\\nginx.conf",
			ContextPosPath: []int{0},
		},
		OriginalFingerprints: map[string]string{
			"E:\\config_test\\nginx.conf": "e2a36380e1591b13cca9d9eb5437bd1a2747901aa5f34caad39aa0960018d492",
		},
	}

	invalidMeta := metav1.WebServerConfigTargetContextOptions{
		WebServerOptions: metav1.WebServerOptions{
			GroupID:    0,
			HostID:     0,
			ServerName: "",
		},
		ConfigContextPos:     metav1.ConfigContextPos{},
		OriginalFingerprints: nil,
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
			name:   "get included configs returns error",
			fields: fields{svc: svc},
			args: args{
				requestMeta: validMeta,
			},
			wantErr: true,
		},
		{
			name:   "fingerprint error",
			fields: fields{svc: svc},
			args: args{
				requestMeta: validMeta,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "get included configs returns error" {
				mockWebServerConfigSrv.EXPECT().GetIncludedConfigs(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("get included configs failed"))
			} else if tt.name == "fingerprint error" {
				mockWebServerConfigSrv.EXPECT().GetIncludedConfigs(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, metav1.ErrInconsistentFingerprints)
			} else if tt.name == "valid request body" {
				mockWebServerConfigSrv.EXPECT().GetIncludedConfigs(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
			}

			w := &WebServerConfigController{
				svc: tt.fields.svc,
			}
			reqBodyBytes, err := json.Marshal(tt.args.requestMeta)
			if err != nil {
				t.Fatal(err)
			}
			resp := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(resp)
			c.Request = httptest.NewRequest("POST", "/api/conf/get-included-configs", bytes.NewBuffer(reqBodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")
			w.GetIncludedConfigs(c)
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

func TestWebServerConfigController_GetOptions(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := svcv1.NewMockFactory(ctrl)
	mockWebServerConfigSrv := svcv1.NewMockWebServerConfigSrv(ctrl)
	svc.EXPECT().WebServerConfigs().AnyTimes().Return(mockWebServerConfigSrv)

	type fields struct {
		svc svcv1.Factory
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "success get options",
			fields:  fields{svc: svc},
			wantErr: false,
		},
		{
			name:    "get options returns error",
			fields:  fields{svc: svc},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				mockWebServerConfigSrv.EXPECT().GetOptions(gomock.Any()).Return(nil, errors.New("database error"))
			} else {
				mockWebServerConfigSrv.EXPECT().GetOptions(gomock.Any()).Return(nil, nil)
			}

			w := &WebServerConfigController{
				svc: tt.fields.svc,
			}
			resp := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(resp)
			c.Request = httptest.NewRequest("GET", "/api/conf/options", nil)
			w.GetOptions(c)
			var respBody response.Response
			if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
				t.Fatal(err)
			}
			if tt.wantErr != (resp.Code != 200 || respBody.Code != 0) {
				t.Errorf("Code: %d, Body: %s", resp.Code, resp.Body)
			}
		})
	}
}

func TestWebServerConfigController_SearchContextPositions(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := svcv1.NewMockFactory(ctrl)
	mockWebServerConfigSrv := svcv1.NewMockWebServerConfigSrv(ctrl)
	svc.EXPECT().WebServerConfigs().AnyTimes().Return(mockWebServerConfigSrv)

	validMeta := metav1.WebServerConfigContextPosSearchOptions{
		WebServerOptions: metav1.WebServerOptions{
			GroupID:    0,
			HostID:     0,
			ServerName: "test-bifrost",
		},
		SearchKeywordsMeta: metav1.SearchKeywordsMeta{
			StartingPositionList: []metav1.ConfigContextPos{
				{
					Config:         "C:\\config_test\\nginx.conf",
					ContextPosPath: []int{0, 0},
				},
			},
			Keywords: "worker_processes",
		},
		OriginalFingerprints: map[string]string{
			"E:\\config_test\\nginx.conf": "e2a36380e1591b13cca9d9eb5437bd1a2747901aa5f34caad39aa0960018d492",
		},
	}

	invalidMeta := metav1.WebServerConfigContextPosSearchOptions{
		WebServerOptions: metav1.WebServerOptions{
			GroupID:    0,
			HostID:     0,
			ServerName: "",
		},
		SearchKeywordsMeta:   metav1.SearchKeywordsMeta{},
		OriginalFingerprints: nil,
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
			name:   "search context positions returns error",
			fields: fields{svc: svc},
			args: args{
				requestMeta: validMeta,
			},
			wantErr: true,
		},
		{
			name:   "fingerprint error",
			fields: fields{svc: svc},
			args: args{
				requestMeta: validMeta,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "search context positions returns error" {
				mockWebServerConfigSrv.EXPECT().SearchContextPositions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("search failed"))
			} else if tt.name == "fingerprint error" {
				mockWebServerConfigSrv.EXPECT().SearchContextPositions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, metav1.ErrInconsistentFingerprints)
			} else if tt.name == "valid request body" {
				mockWebServerConfigSrv.EXPECT().SearchContextPositions(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]metav1.ConfigContextPos{
					{
						Config:         "C:\\config_test\\nginx.conf",
						ContextPosPath: []int{0, 0},
					},
				}, nil)
			}

			w := &WebServerConfigController{
				svc: tt.fields.svc,
			}
			reqBodyBytes, err := json.Marshal(tt.args.requestMeta)
			if err != nil {
				t.Fatal(err)
			}
			resp := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(resp)
			c.Request = httptest.NewRequest("POST", "/api/conf/search-context-positions", bytes.NewBuffer(reqBodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")
			w.SearchContextPositions(c)
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

func Test_getErrorHandle(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()

	tests := []struct {
		name        string
		err         error
		okdetailobj interface{}
		okmsg       string
		failuremsg  string
		wantCode    int
		wantMsg     string
	}{
		{
			name:        "no error",
			err:         nil,
			okdetailobj: "test data",
			okmsg:       "success",
			failuremsg:  "failure",
			wantCode:    200,
			wantMsg:     "success",
		},
		{
			name:        "fingerprint error - ErrInconsistentFingerprints",
			err:         metav1.ErrInconsistentFingerprints,
			okdetailobj: nil,
			okmsg:       "success",
			failuremsg:  "failure",
			wantCode:    200,
			wantMsg:     "指纹校验失败, 请重新查询, 刷新配置文件!",
		},
		{
			name:        "fingerprint error - error code 110010",
			err:         errors.WithCode(110010, "fingerprint mismatch"),
			okdetailobj: nil,
			okmsg:       "success",
			failuremsg:  "failure",
			wantCode:    200,
			wantMsg:     "指纹校验失败, 请重新查询, 刷新配置文件!",
		},
		{
			name:        "generic error",
			err:         errors.New("some error"),
			okdetailobj: nil,
			okmsg:       "success",
			failuremsg:  "operation failed",
			wantCode:    200,
			wantMsg:     "operation failed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(resp)
			c.Request = httptest.NewRequest("POST", "/test", nil)

			getErrorHandle(c, tt.err, tt.okdetailobj, tt.okmsg, tt.failuremsg)

			var respBody response.Response
			if err := json.Unmarshal(resp.Body.Bytes(), &respBody); err != nil {
				t.Fatal(err)
			}
			if resp.Code != tt.wantCode {
				t.Errorf("got code %d, want %d", resp.Code, tt.wantCode)
			}
			if respBody.Msg != tt.wantMsg {
				t.Errorf("got msg %q, want %q", respBody.Msg, tt.wantMsg)
			}
		})
	}
}
