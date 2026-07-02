package service

import (
	"context"
	"reflect"
	"testing"

	v1 "gin-vue-admin/api/heimdallr_api/v1"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	epclientv1 "gin-vue-admin/pkg/client/v1/endpoint"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"go.uber.org/mock/gomock"
)

func Test_hostService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockHostEndpoints(ctrl)
	ctx := context.Background()

	mockEndpoint := httpclientv1.NewEndpoint[metav1.IDOptions, *v1.Host](func(ctx context.Context, req interface{}) (interface{}, error) {
		return &v1.Host{Name: "test-host"}, nil
	})
	mockEndpoints.EXPECT().Get().Return(mockEndpoint)

	type fields struct {
		ctx context.Context
		eps epclientv1.HostEndpoints
	}
	type args struct {
		idOptions *metav1.IDOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *v1.Host
		wantErr bool
	}{
		{
			name: "successful get host",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			args: args{
				idOptions: &metav1.IDOptions{ID: 1},
			},
			want:    &v1.Host{Name: "test-host"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &hostService{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			got, err := s.Get(tt.args.idOptions)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil && tt.want != nil {
				t.Errorf("Get() got = nil, want %v", tt.want)
			}
		})
	}
}

func Test_hostService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockHostEndpoints(ctrl)
	ctx := context.Background()

	mockEndpoint := httpclientv1.NewEndpoint[metav1.ListOptions, *v1.HostList](func(ctx context.Context, req interface{}) (interface{}, error) {
		return &v1.HostList{Items: []*v1.Host{{Name: "test-host"}}}, nil
	})
	mockEndpoints.EXPECT().List().Return(mockEndpoint)

	type fields struct {
		ctx context.Context
		eps epclientv1.HostEndpoints
	}
	type args struct {
		listOptions *metav1.ListOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *v1.HostList
		wantErr bool
	}{
		{
			name: "successful list hosts",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			args: args{
				listOptions: &metav1.ListOptions{},
			},
			want:    &v1.HostList{Items: []*v1.Host{{Name: "test-host"}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &hostService{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			got, err := s.List(tt.args.listOptions)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil && tt.want != nil {
				t.Errorf("List() got = nil, want %v", tt.want)
			}
		})
	}
}

func Test_newHostService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockHostEndpoints(ctrl)
	mockEndpointFactory := epclientv1.NewMockFactory(ctrl)
	mockEndpointFactory.EXPECT().Hosts().Return(mockEndpoints).AnyTimes()

	ctx := context.Background()
	svcFactory := &factory{ctx: ctx, eps: mockEndpointFactory}

	type args struct {
		factory *factory
	}
	tests := []struct {
		name string
		args args
		want HostService
	}{
		{
			name: "successful creation",
			args: args{factory: svcFactory},
			want: &hostService{ctx: ctx, eps: mockEndpoints},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newHostService(tt.args.factory)
			if !reflect.DeepEqual(reflect.TypeOf(got), reflect.TypeOf(tt.want)) {
				t.Errorf("newHostService() type = %T, want %T", got, tt.want)
			}
		})
	}
}
