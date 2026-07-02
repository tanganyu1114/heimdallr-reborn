package service

import (
	"context"
	"reflect"
	"testing"

	v1 "gin-vue-admin/api/heimdallr_api/v1"
	epclientv1 "gin-vue-admin/pkg/client/v1/endpoint"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"go.uber.org/mock/gomock"
)

func Test_agentInfoService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockAgentInfoEndpoints(ctrl)
	ctx := context.Background()

	mockEndpoint := httpclientv1.NewEndpoint[httpclientv1.NilBody, []v1.GroupInfo](func(ctx context.Context, req interface{}) (interface{}, error) {
		return []v1.GroupInfo{{Name: "test-group"}}, nil
	})
	mockEndpoints.EXPECT().Get().Return(mockEndpoint)

	type fields struct {
		ctx context.Context
		eps epclientv1.AgentInfoEndpoints
	}
	tests := []struct {
		name    string
		fields  fields
		want    []v1.GroupInfo
		wantErr bool
	}{
		{
			name: "successful get agent info",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			want:    []v1.GroupInfo{{Name: "test-group"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &agentInfoService{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			got, err := s.Get()
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newAgentInfoService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockAgentInfoEndpoints(ctrl)
	mockEndpointFactory := epclientv1.NewMockFactory(ctrl)
	mockEndpointFactory.EXPECT().AgentInfos().Return(mockEndpoints).AnyTimes()

	ctx := context.Background()
	svcFactory := &factory{ctx: ctx, eps: mockEndpointFactory}

	type args struct {
		factory *factory
	}
	tests := []struct {
		name string
		args args
		want AgentInfoService
	}{
		{
			name: "successful creation",
			args: args{factory: svcFactory},
			want: &agentInfoService{ctx: ctx, eps: mockEndpoints},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newAgentInfoService(tt.args.factory)
			if !reflect.DeepEqual(reflect.TypeOf(got), reflect.TypeOf(tt.want)) {
				t.Errorf("newAgentInfoService() type = %T, want %T", got, tt.want)
			}
		})
	}
}
