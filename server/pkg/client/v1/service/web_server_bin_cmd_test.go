package service

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	metav1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	epclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/endpoint"
	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/model"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"go.uber.org/mock/gomock"
)

func Test_webServerBinCMDService_Exec(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockWebServerBinCMDEndpoints(ctrl)
	ctx := context.Background()

	mockEndpoint := httpclientv1.NewEndpoint[metav1.WebServerBinCMDExecRequest, modelclientv1.ResponseBody[*metav1.WebServerBinCMDExecResponse]](func(ctx context.Context, req interface{}) (interface{}, error) {
		data, _ := json.Marshal(&metav1.WebServerBinCMDExecResponse{Successful: true, Stdout: "command output"})
		return modelclientv1.ResponseBody[*metav1.WebServerBinCMDExecResponse]{Data: data}, nil
	})
	mockEndpoints.EXPECT().Exec().Return(mockEndpoint)

	type fields struct {
		ctx context.Context
		eps epclientv1.WebServerBinCMDEndpoints
	}
	type args struct {
		req *metav1.WebServerBinCMDExecRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *metav1.WebServerBinCMDExecResponse
		wantErr bool
	}{
		{
			name: "successful exec",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			args: args{
				req: &metav1.WebServerBinCMDExecRequest{},
			},
			want:    &metav1.WebServerBinCMDExecResponse{Successful: true, Stdout: "command output"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &webServerBinCMDService{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			got, err := s.Exec(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil && tt.want != nil {
				t.Errorf("Exec() got = nil, want %v", tt.want)
			}
		})
	}
}

func Test_newWebServerBinCMDService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockWebServerBinCMDEndpoints(ctrl)
	mockEndpointFactory := epclientv1.NewMockFactory(ctrl)
	mockEndpointFactory.EXPECT().WebServerBinCMDs().Return(mockEndpoints).AnyTimes()

	ctx := context.Background()
	svcFactory := &factory{ctx: ctx, eps: mockEndpointFactory}

	type args struct {
		factory *factory
	}
	tests := []struct {
		name string
		args args
		want WebServerBinCMDService
	}{
		{
			name: "successful creation",
			args: args{factory: svcFactory},
			want: &webServerBinCMDService{ctx: ctx, eps: mockEndpoints},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newWebServerBinCMDService(tt.args.factory)
			if !reflect.DeepEqual(reflect.TypeOf(got), reflect.TypeOf(tt.want)) {
				t.Errorf("newWebServerBinCMDService() type = %T, want %T", got, tt.want)
			}
		})
	}
}
