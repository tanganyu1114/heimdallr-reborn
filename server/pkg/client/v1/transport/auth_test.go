package transport

import (
	"reflect"
	"testing"

	"github.com/tanganyu1114/heimdallr-reborn/server/model/request"
	"github.com/tanganyu1114/heimdallr-reborn/server/model/response"
	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/model"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"go.uber.org/mock/gomock"
)

func Test_newAuthTransport(t *testing.T) {
	type args struct {
		transport *transport
	}
	tests := []struct {
		name string
		args args
		want AuthTransport
	}{
		{
			name: "creates auth transport with valid base URL",
			args: args{
				transport: &transport{
					baseURL: "http://localhost:8080",
				},
			},
		},
		{
			name: "creates auth transport with empty base URL",
			args: args{
				transport: &transport{
					baseURL: "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newAuthTransport(tt.args.transport)
			if got == nil {
				t.Errorf("newAuthTransport() = nil, want non-nil")
			}
		})
	}
}

func Test_authTransport_ValidateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[*request.ValidateTokenRequest, modelclientv1.ResponseBody[*response.ValidateTokenResponse]](ctrl)

	type fields struct {
		validateTokenClient   httpclientv1.ClientBuilder[*request.ValidateTokenRequest, modelclientv1.ResponseBody[*response.ValidateTokenResponse]]
		refreshTokenClient    httpclientv1.ClientBuilder[*request.RefreshTokenRequest, modelclientv1.ResponseBody[*response.RefreshTokenResponse]]
		checkPermissionClient httpclientv1.ClientBuilder[*request.CheckPermissionRequest, modelclientv1.ResponseBody[*response.CheckPermissionResponse]]
		logOperationClient    httpclientv1.ClientBuilder[*request.LogOperationRequest, modelclientv1.ResponseBody[httpclientv1.NilBody]]
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[*request.ValidateTokenRequest, modelclientv1.ResponseBody[*response.ValidateTokenResponse]]
	}{
		{
			name: "returns validate token client",
			fields: fields{
				validateTokenClient: mockClient,
			},
			want: mockClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &authTransport{
				validateTokenClient:   tt.fields.validateTokenClient,
				refreshTokenClient:    tt.fields.refreshTokenClient,
				checkPermissionClient: tt.fields.checkPermissionClient,
				logOperationClient:    tt.fields.logOperationClient,
			}
			if got := a.ValidateToken(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authTransport_RefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[*request.RefreshTokenRequest, modelclientv1.ResponseBody[*response.RefreshTokenResponse]](ctrl)

	type fields struct {
		validateTokenClient   httpclientv1.ClientBuilder[*request.ValidateTokenRequest, modelclientv1.ResponseBody[*response.ValidateTokenResponse]]
		refreshTokenClient    httpclientv1.ClientBuilder[*request.RefreshTokenRequest, modelclientv1.ResponseBody[*response.RefreshTokenResponse]]
		checkPermissionClient httpclientv1.ClientBuilder[*request.CheckPermissionRequest, modelclientv1.ResponseBody[*response.CheckPermissionResponse]]
		logOperationClient    httpclientv1.ClientBuilder[*request.LogOperationRequest, modelclientv1.ResponseBody[httpclientv1.NilBody]]
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[*request.RefreshTokenRequest, modelclientv1.ResponseBody[*response.RefreshTokenResponse]]
	}{
		{
			name: "returns refresh token client",
			fields: fields{
				refreshTokenClient: mockClient,
			},
			want: mockClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &authTransport{
				validateTokenClient:   tt.fields.validateTokenClient,
				refreshTokenClient:    tt.fields.refreshTokenClient,
				checkPermissionClient: tt.fields.checkPermissionClient,
				logOperationClient:    tt.fields.logOperationClient,
			}
			if got := a.RefreshToken(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RefreshToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authTransport_CheckPermission(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[*request.CheckPermissionRequest, modelclientv1.ResponseBody[*response.CheckPermissionResponse]](ctrl)

	type fields struct {
		validateTokenClient   httpclientv1.ClientBuilder[*request.ValidateTokenRequest, modelclientv1.ResponseBody[*response.ValidateTokenResponse]]
		refreshTokenClient    httpclientv1.ClientBuilder[*request.RefreshTokenRequest, modelclientv1.ResponseBody[*response.RefreshTokenResponse]]
		checkPermissionClient httpclientv1.ClientBuilder[*request.CheckPermissionRequest, modelclientv1.ResponseBody[*response.CheckPermissionResponse]]
		logOperationClient    httpclientv1.ClientBuilder[*request.LogOperationRequest, modelclientv1.ResponseBody[httpclientv1.NilBody]]
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[*request.CheckPermissionRequest, modelclientv1.ResponseBody[*response.CheckPermissionResponse]]
	}{
		{
			name: "returns check permission client",
			fields: fields{
				checkPermissionClient: mockClient,
			},
			want: mockClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &authTransport{
				validateTokenClient:   tt.fields.validateTokenClient,
				refreshTokenClient:    tt.fields.refreshTokenClient,
				checkPermissionClient: tt.fields.checkPermissionClient,
				logOperationClient:    tt.fields.logOperationClient,
			}
			if got := a.CheckPermission(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckPermission() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authTransport_LogOperation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[*request.LogOperationRequest, modelclientv1.ResponseBody[httpclientv1.NilBody]](ctrl)

	type fields struct {
		validateTokenClient   httpclientv1.ClientBuilder[*request.ValidateTokenRequest, modelclientv1.ResponseBody[*response.ValidateTokenResponse]]
		refreshTokenClient    httpclientv1.ClientBuilder[*request.RefreshTokenRequest, modelclientv1.ResponseBody[*response.RefreshTokenResponse]]
		checkPermissionClient httpclientv1.ClientBuilder[*request.CheckPermissionRequest, modelclientv1.ResponseBody[*response.CheckPermissionResponse]]
		logOperationClient    httpclientv1.ClientBuilder[*request.LogOperationRequest, modelclientv1.ResponseBody[httpclientv1.NilBody]]
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[*request.LogOperationRequest, modelclientv1.ResponseBody[httpclientv1.NilBody]]
	}{
		{
			name: "returns log operation client",
			fields: fields{
				logOperationClient: mockClient,
			},
			want: mockClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &authTransport{
				validateTokenClient:   tt.fields.validateTokenClient,
				refreshTokenClient:    tt.fields.refreshTokenClient,
				checkPermissionClient: tt.fields.checkPermissionClient,
				logOperationClient:    tt.fields.logOperationClient,
			}
			if got := a.LogOperation(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LogOperation() = %v, want %v", got, tt.want)
			}
		})
	}
}
