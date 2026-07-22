package web_server_config

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	"github.com/tanganyu1114/heimdallr-reborn/server/global"
	svcv1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/service/v1"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/response"

	logV1 "github.com/ClessLi/component-base/pkg/log/v1"
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/errors"
	"go.uber.org/mock/gomock"
)

func TestWebServerConfigController_UpdateConfig(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := svcv1.NewMockFactory(ctrl)
	mockWebServerConfigSrv := svcv1.NewMockWebServerConfigSrv(ctrl)
	svc.EXPECT().WebServerConfigs().AnyTimes().Return(mockWebServerConfigSrv)

	nginxConfigJSON := `{
		"main-config": "E:\\config_test\\nginx.conf",
		"configs": {
			"E:\\config_test\\nginx.conf": {
				"enabled": true,
				"context-type": "config",
				"value": "E:\\config_test\\nginx.conf",
				"params": [
					{"context-type": "directive", "value": "worker_processes 1"},
					{"context-type": "events", "params": [{"context-type": "directive", "value": "worker_connections 1024"}]},
					{"context-type": "http", "params": [{"context-type": "directive", "value": "sendfile on"}]}
				]
			}
		}
	}`

	validMeta := v1.WebServerConfigUpdateOptions{
		WebServerOptions: v1.WebServerOptions{
			GroupID:    0,
			HostID:     0,
			ServerName: "test-bifrost",
		},
		Data: []byte(nginxConfigJSON),
		OriginalFingerprints: map[string]string{
			"E:\\config_test\\nginx.conf": "e2a36380e1591b13cca9d9eb5437bd1a2747901aa5f34caad39aa0960018d492",
		},
	}

	invalidMeta := v1.WebServerConfigUpdateOptions{
		WebServerOptions: v1.WebServerOptions{
			GroupID:    0,
			HostID:     0,
			ServerName: "",
		},
		Data:                 nil,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWebServerConfigSrv.EXPECT().UpdateConfig(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

			w := &WebServerConfigController{
				svc: tt.fields.svc,
			}
			reqBodyBytes, err := json.Marshal(tt.args.requestMeta)
			if err != nil {
				t.Fatal(err)
			}
			resp := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(resp)
			c.Request = httptest.NewRequest("POST", "/api/conf/update", bytes.NewBuffer(reqBodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")
			w.UpdateConfig(c)
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

func TestWebServerConfigController_ChangeContextEnabledState(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := svcv1.NewMockFactory(ctrl)
	mockWebServerConfigSrv := svcv1.NewMockWebServerConfigSrv(ctrl)
	svc.EXPECT().WebServerConfigs().AnyTimes().Return(mockWebServerConfigSrv)

	validMeta := v1.WebServerConfigContextUpdateOptions[v1.ConfigContextEnabledStateMeta]{
		WebServerOptions: v1.WebServerOptions{
			GroupID:    0,
			HostID:     0,
			ServerName: "test-bifrost",
		},
		TargetConfigContextOptions: v1.TargetConfigContextOptions[v1.ConfigContextEnabledStateMeta]{
			Position: v1.ConfigContextPos{
				Config:         "E:\\config_test\\nginx.conf",
				ContextPosPath: []int{0},
			},
			TargetContext: v1.ConfigContextEnabledStateMeta{
				Enabled: true,
			},
		},
		OriginalFingerprints: map[string]string{
			"E:\\config_test\\nginx.conf": "e2a36380e1591b13cca9d9eb5437bd1a2747901aa5f34caad39aa0960018d492",
		},
	}

	invalidMeta := v1.WebServerConfigContextUpdateOptions[v1.ConfigContextEnabledStateMeta]{
		WebServerOptions: v1.WebServerOptions{
			GroupID:    0,
			HostID:     0,
			ServerName: "",
		},
		TargetConfigContextOptions: v1.TargetConfigContextOptions[v1.ConfigContextEnabledStateMeta]{
			Position:      v1.ConfigContextPos{},
			TargetContext: v1.ConfigContextEnabledStateMeta{},
		},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWebServerConfigSrv.EXPECT().ChangeContextEnabledState(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

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
			c.Request.Header.Set("Content-Type", "application/json")
			w.ChangeContextEnabledState(c)
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

func TestWebServerConfigController_InsertWithClone(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := svcv1.NewMockFactory(ctrl)
	mockWebServerConfigSrv := svcv1.NewMockWebServerConfigSrv(ctrl)
	svc.EXPECT().WebServerConfigs().AnyTimes().Return(mockWebServerConfigSrv)

	validMeta := v1.WebServerConfigContextUpdateOptions[v1.CloneConfigContextMeta]{
		WebServerOptions: v1.WebServerOptions{
			GroupID:    0,
			HostID:     0,
			ServerName: "test-bifrost",
		},
		TargetConfigContextOptions: v1.TargetConfigContextOptions[v1.CloneConfigContextMeta]{
			Position: v1.ConfigContextPos{
				Config:         "E:\\config_test\\nginx.conf",
				ContextPosPath: []int{0},
			},
			TargetContext: v1.CloneConfigContextMeta{
				ConfigContextPos: v1.ConfigContextPos{
					Config:         "E:\\config_test\\clone.conf",
					ContextPosPath: []int{1},
				},
			},
		},
		DisableTheTarget: false,
		OriginalFingerprints: map[string]string{
			"E:\\config_test\\nginx.conf": "e2a36380e1591b13cca9d9eb5437bd1a2747901aa5f34caad39aa0960018d492",
		},
	}

	invalidMeta := v1.WebServerConfigContextUpdateOptions[v1.CloneConfigContextMeta]{
		WebServerOptions: v1.WebServerOptions{
			GroupID:    0,
			HostID:     0,
			ServerName: "",
		},
		TargetConfigContextOptions: v1.TargetConfigContextOptions[v1.CloneConfigContextMeta]{
			Position:      v1.ConfigContextPos{},
			TargetContext: v1.CloneConfigContextMeta{},
		},
		DisableTheTarget:     false,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWebServerConfigSrv.EXPECT().InsertWithClone(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

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
			c.Request.Header.Set("Content-Type", "application/json")
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
	global.GVA_LOG = logV1.ZapLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := svcv1.NewMockFactory(ctrl)
	mockWebServerConfigSrv := svcv1.NewMockWebServerConfigSrv(ctrl)
	svc.EXPECT().WebServerConfigs().AnyTimes().Return(mockWebServerConfigSrv)

	validMeta := v1.WebServerConfigContextUpdateOptions[v1.NewConfigContextMeta]{
		WebServerOptions: v1.WebServerOptions{
			GroupID:    0,
			HostID:     0,
			ServerName: "test-bifrost",
		},
		TargetConfigContextOptions: v1.TargetConfigContextOptions[v1.NewConfigContextMeta]{
			Position: v1.ConfigContextPos{
				Config:         "E:\\config_test\\nginx.conf",
				ContextPosPath: []int{0},
			},
			TargetContext: v1.NewConfigContextMeta{
				ContextType:  "directive",
				ContextValue: "worker_processes 2",
			},
		},
		DisableTheTarget: false,
		OriginalFingerprints: map[string]string{
			"E:\\config_test\\nginx.conf": "e2a36380e1591b13cca9d9eb5437bd1a2747901aa5f34caad39aa0960018d492",
		},
	}

	invalidMeta := v1.WebServerConfigContextUpdateOptions[v1.NewConfigContextMeta]{
		WebServerOptions: v1.WebServerOptions{
			GroupID:    0,
			HostID:     0,
			ServerName: "",
		},
		TargetConfigContextOptions: v1.TargetConfigContextOptions[v1.NewConfigContextMeta]{
			Position:      v1.ConfigContextPos{},
			TargetContext: v1.NewConfigContextMeta{},
		},
		DisableTheTarget:     false,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWebServerConfigSrv.EXPECT().InsertWithNew(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

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
			c.Request.Header.Set("Content-Type", "application/json")
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

func TestWebServerConfigController_ModifyContextValue(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := svcv1.NewMockFactory(ctrl)
	mockWebServerConfigSrv := svcv1.NewMockWebServerConfigSrv(ctrl)
	svc.EXPECT().WebServerConfigs().AnyTimes().Return(mockWebServerConfigSrv)

	validMeta := v1.WebServerConfigContextUpdateOptions[v1.NewConfigContextMeta]{
		WebServerOptions: v1.WebServerOptions{
			GroupID:    0,
			HostID:     0,
			ServerName: "test-bifrost",
		},
		TargetConfigContextOptions: v1.TargetConfigContextOptions[v1.NewConfigContextMeta]{
			Position: v1.ConfigContextPos{
				Config:         "E:\\config_test\\nginx.conf",
				ContextPosPath: []int{0},
			},
			TargetContext: v1.NewConfigContextMeta{
				ContextType:  "directive",
				ContextValue: "worker_processes 4",
			},
		},
		OriginalFingerprints: map[string]string{
			"E:\\config_test\\nginx.conf": "e2a36380e1591b13cca9d9eb5437bd1a2747901aa5f34caad39aa0960018d492",
		},
	}

	invalidMeta := v1.WebServerConfigContextUpdateOptions[v1.NewConfigContextMeta]{
		WebServerOptions: v1.WebServerOptions{
			GroupID:    0,
			HostID:     0,
			ServerName: "",
		},
		TargetConfigContextOptions: v1.TargetConfigContextOptions[v1.NewConfigContextMeta]{
			Position:      v1.ConfigContextPos{},
			TargetContext: v1.NewConfigContextMeta{},
		},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWebServerConfigSrv.EXPECT().ModifyContextValue(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

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
			c.Request.Header.Set("Content-Type", "application/json")
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
	global.GVA_LOG = logV1.ZapLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := svcv1.NewMockFactory(ctrl)
	mockWebServerConfigSrv := svcv1.NewMockWebServerConfigSrv(ctrl)
	svc.EXPECT().WebServerConfigs().AnyTimes().Return(mockWebServerConfigSrv)

	validMeta := v1.WebServerConfigContextUpdateOptions[v1.CloneConfigContextMeta]{
		WebServerOptions: v1.WebServerOptions{
			GroupID:    0,
			HostID:     0,
			ServerName: "test-bifrost",
		},
		TargetConfigContextOptions: v1.TargetConfigContextOptions[v1.CloneConfigContextMeta]{
			Position: v1.ConfigContextPos{
				Config:         "E:\\config_test\\nginx.conf",
				ContextPosPath: []int{0},
			},
			TargetContext: v1.CloneConfigContextMeta{
				ConfigContextPos: v1.ConfigContextPos{
					Config:         "E:\\config_test\\clone.conf",
					ContextPosPath: []int{1},
				},
			},
		},
		OriginalFingerprints: map[string]string{
			"E:\\config_test\\nginx.conf": "e2a36380e1591b13cca9d9eb5437bd1a2747901aa5f34caad39aa0960018d492",
		},
	}

	invalidMeta := v1.WebServerConfigContextUpdateOptions[v1.CloneConfigContextMeta]{
		WebServerOptions: v1.WebServerOptions{
			GroupID:    0,
			HostID:     0,
			ServerName: "",
		},
		TargetConfigContextOptions: v1.TargetConfigContextOptions[v1.CloneConfigContextMeta]{
			Position:      v1.ConfigContextPos{},
			TargetContext: v1.CloneConfigContextMeta{},
		},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWebServerConfigSrv.EXPECT().ModifyWithClone(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

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
			c.Request.Header.Set("Content-Type", "application/json")
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
	global.GVA_LOG = logV1.ZapLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := svcv1.NewMockFactory(ctrl)
	mockWebServerConfigSrv := svcv1.NewMockWebServerConfigSrv(ctrl)
	svc.EXPECT().WebServerConfigs().AnyTimes().Return(mockWebServerConfigSrv)

	validMeta := v1.WebServerConfigContextUpdateOptions[v1.NewConfigContextMeta]{
		WebServerOptions: v1.WebServerOptions{
			GroupID:    0,
			HostID:     0,
			ServerName: "test-bifrost",
		},
		TargetConfigContextOptions: v1.TargetConfigContextOptions[v1.NewConfigContextMeta]{
			Position: v1.ConfigContextPos{
				Config:         "E:\\config_test\\nginx.conf",
				ContextPosPath: []int{0},
			},
			TargetContext: v1.NewConfigContextMeta{
				ContextType:  "directive",
				ContextValue: "worker_processes 8",
			},
		},
		OriginalFingerprints: map[string]string{
			"E:\\config_test\\nginx.conf": "e2a36380e1591b13cca9d9eb5437bd1a2747901aa5f34caad39aa0960018d492",
		},
	}

	invalidMeta := v1.WebServerConfigContextUpdateOptions[v1.NewConfigContextMeta]{
		WebServerOptions: v1.WebServerOptions{
			GroupID:    0,
			HostID:     0,
			ServerName: "",
		},
		TargetConfigContextOptions: v1.TargetConfigContextOptions[v1.NewConfigContextMeta]{
			Position:      v1.ConfigContextPos{},
			TargetContext: v1.NewConfigContextMeta{},
		},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWebServerConfigSrv.EXPECT().ModifyWithNew(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

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
			c.Request.Header.Set("Content-Type", "application/json")
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

func TestWebServerConfigController_Move(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	svc := svcv1.NewMockFactory(ctrl)
	mockWebServerConfigSrv := svcv1.NewMockWebServerConfigSrv(ctrl)
	svc.EXPECT().WebServerConfigs().AnyTimes().Return(mockWebServerConfigSrv)

	validMeta := v1.WebServerConfigContextUpdateOptions[v1.CloneConfigContextMeta]{
		WebServerOptions: v1.WebServerOptions{
			GroupID:    0,
			HostID:     0,
			ServerName: "test-bifrost",
		},
		TargetConfigContextOptions: v1.TargetConfigContextOptions[v1.CloneConfigContextMeta]{
			Position: v1.ConfigContextPos{
				Config:         "E:\\config_test\\nginx.conf",
				ContextPosPath: []int{0},
			},
			TargetContext: v1.CloneConfigContextMeta{
				ConfigContextPos: v1.ConfigContextPos{
					Config:         "E:\\config_test\\nginx.conf",
					ContextPosPath: []int{2},
				},
			},
		},
		DisableTheTarget: false,
		OriginalFingerprints: map[string]string{
			"E:\\config_test\\nginx.conf": "e2a36380e1591b13cca9d9eb5437bd1a2747901aa5f34caad39aa0960018d492",
		},
	}

	invalidMeta := v1.WebServerConfigContextUpdateOptions[v1.CloneConfigContextMeta]{
		WebServerOptions: v1.WebServerOptions{
			GroupID:    0,
			HostID:     0,
			ServerName: "",
		},
		TargetConfigContextOptions: v1.TargetConfigContextOptions[v1.CloneConfigContextMeta]{
			Position:      v1.ConfigContextPos{},
			TargetContext: v1.CloneConfigContextMeta{},
		},
		DisableTheTarget:     false,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWebServerConfigSrv.EXPECT().Move(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

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
			c.Request.Header.Set("Content-Type", "application/json")
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

func Test_updateErrorHandle(t *testing.T) {
	global.GVA_LOG = logV1.ZapLogger()

	tests := []struct {
		name       string
		err        error
		okmsg      string
		failuremsg string
		wantCode   int
		wantMsg    string
	}{
		{
			name:       "no error",
			err:        nil,
			okmsg:      "success",
			failuremsg: "failure",
			wantCode:   200,
			wantMsg:    "success",
		},
		{
			name:       "fingerprint error - ErrInconsistentFingerprints",
			err:        v1.ErrInconsistentFingerprints,
			okmsg:      "success",
			failuremsg: "failure",
			wantCode:   200,
			wantMsg:    "指纹校验失败, 请重新查询, 刷新配置文件!",
		},
		{
			name:       "fingerprint error - error code 110010",
			err:        errors.WithCode(110010, "fingerprint mismatch"),
			okmsg:      "success",
			failuremsg: "failure",
			wantCode:   200,
			wantMsg:    "指纹校验失败, 请重新查询, 刷新配置文件!",
		},
		{
			name:       "generic error",
			err:        errors.New("some error"),
			okmsg:      "success",
			failuremsg: "operation failed",
			wantCode:   200,
			wantMsg:    "operation failed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(resp)
			c.Request = httptest.NewRequest("POST", "/test", nil)

			updateErrorHandle(c, tt.err, tt.okmsg, tt.failuremsg)

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
