package transport

import (
	"reflect"
	"testing"

	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/model"

	"github.com/tanganyu1114/heimdallr-reborn/server/model/request"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/response"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"go.uber.org/mock/gomock"
)

func Test_sysUserTransport_GetSDKChallenge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[*request.SDKChallengeRequest, modelclientv1.ResponseBody[*response.SDKChallengeResponse]](ctrl)

	type fields struct {
		sdkChallengeClient httpclientv1.ClientBuilder[*request.SDKChallengeRequest, modelclientv1.ResponseBody[*response.SDKChallengeResponse]]
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[*request.SDKChallengeRequest, modelclientv1.ResponseBody[*response.SDKChallengeResponse]]
	}{
		{
			name: "returns SDK challenge client",
			fields: fields{
				sdkChallengeClient: mockClient,
			},
			want: mockClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sysUserTransport{
				sdkChallengeClient: tt.fields.sdkChallengeClient,
			}
			if got := s.GetSDKChallenge(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSDKChallenge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sysUserTransport_SDKLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[*request.EncryptedLoginRequest, modelclientv1.ResponseBody[*response.LoginResponse]](ctrl)

	type fields struct {
		sdkLoginClient httpclientv1.ClientBuilder[*request.EncryptedLoginRequest, modelclientv1.ResponseBody[*response.LoginResponse]]
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[*request.EncryptedLoginRequest, modelclientv1.ResponseBody[*response.LoginResponse]]
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
