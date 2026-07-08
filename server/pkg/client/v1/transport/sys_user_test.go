package transport

import (
	modelclientv1 "gin-vue-admin/pkg/client/v1/model"
	"reflect"
	"testing"

	"gin-vue-admin/model/request"
	"gin-vue-admin/model/response"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"go.uber.org/mock/gomock"
)

func Test_sysUserTransport_SDKLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[*request.SDKLogin, modelclientv1.ResponseBody[*response.LoginResponse]](ctrl)

	type fields struct {
		sdkLoginClient httpclientv1.ClientBuilder[*request.SDKLogin, modelclientv1.ResponseBody[*response.LoginResponse]]
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[*request.SDKLogin, modelclientv1.ResponseBody[*response.LoginResponse]]
	}{
		{
			name: "returns SDK login client",
			fields: fields{
				sdkLoginClient: mockClient,
			},
			want: mockClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sysUserTransport{
				sdkLoginClient: tt.fields.sdkLoginClient,
			}
			if got := s.SDKLogin(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SDKLogin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newSysUserTransport(t *testing.T) {
	type args struct {
		transport *transport
	}
	tests := []struct {
		name string
		args args
		want SysUserTransport
	}{
		{
			name: "creates sys user transport with valid base URL",
			args: args{
				transport: &transport{
					baseURL: "http://localhost:8080",
				},
			},
		},
		{
			name: "creates sys user transport with empty base URL",
			args: args{
				transport: &transport{
					baseURL: "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newSysUserTransport(tt.args.transport)
			if got == nil {
				t.Errorf("newSysUserTransport() = nil, want non-nil")
			}
		})
	}
}
