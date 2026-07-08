package web_server_config

import (
	v1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/service/v1"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestNewController(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type args struct {
		service v1.Factory
	}
	tests := []struct {
		name string
		args args
		want *WebServerConfigController
	}{
		{
			name: "create controller with mock service",
			args: args{
				service: v1.NewMockFactory(ctrl),
			},
			want: &WebServerConfigController{
				svc: v1.NewMockFactory(ctrl),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewController(tt.args.service)
			if got == nil {
				t.Errorf("NewController() returned nil")
				return
			}
			if !reflect.DeepEqual(got.svc != nil, tt.want.svc != nil) {
				t.Errorf("NewController() svc = %v, want %v", got.svc, tt.want.svc)
			}
		})
	}
}
