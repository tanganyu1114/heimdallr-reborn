package web_server_config

import (
	"bytes"
	"encoding/json"
	"github.com/tanganyu1114/heimdallr-reborn/global"
	svcv1 "github.com/tanganyu1114/heimdallr-reborn/internal/hmdr_api/service/v1"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/internal/pkg/meta/v1"
	"github.com/tanganyu1114/heimdallr-reborn/model/response"
	"net/http/httptest"
	"testing"

	logV1 "github.com/ClessLi/component-base/pkg/log/v1"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

func TestWebServerConfigController_Remove(t *testing.T) {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWebServerConfigSrv.EXPECT().Remove(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

			w := &WebServerConfigController{
				svc: tt.fields.svc,
			}
			reqBodyBytes, err := json.Marshal(tt.args.requestMeta)
			if err != nil {
				t.Fatal(err)
			}
			resp := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(resp)
			c.Request = httptest.NewRequest("POST", "/api/conf/remove", bytes.NewBuffer(reqBodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")
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
