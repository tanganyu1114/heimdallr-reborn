package service

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	v1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	epclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/endpoint"
	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/model"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"go.uber.org/mock/gomock"
)

func Test_groupService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockGroupEndpoints(ctrl)
	ctx := context.Background()

	mockEndpoint := httpclientv1.NewEndpoint[v1.IDOptions, modelclientv1.ResponseBody[*v1.Group]](func(ctx context.Context, req interface{}) (interface{}, error) {
		data, _ := json.Marshal(&v1.Group{Name: "test-group"})
		return modelclientv1.ResponseBody[*v1.Group]{Data: data}, nil
	})
	mockEndpoints.EXPECT().Get().Return(mockEndpoint)

	type fields struct {
		ctx context.Context
		eps epclientv1.GroupEndpoints
	}
	type args struct {
		idOptions *v1.IDOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *v1.Group
		wantErr bool
	}{
		{
			name: "successful get group",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			args: args{
				idOptions: &v1.IDOptions{ID: 1},
			},
			want:    &v1.Group{Name: "test-group"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &groupService{
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

func Test_groupService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockGroupEndpoints(ctrl)
	ctx := context.Background()

	mockEndpoint := httpclientv1.NewEndpoint[v1.ListOptions, modelclientv1.ResponseBody[*v1.GroupList]](func(ctx context.Context, req interface{}) (interface{}, error) {
		data, _ := json.Marshal(&v1.GroupList{Items: []*v1.Group{{Name: "test-group"}}})
		return modelclientv1.ResponseBody[*v1.GroupList]{Data: data}, nil
	})
	mockEndpoints.EXPECT().List().Return(mockEndpoint)

	type fields struct {
		ctx context.Context
		eps epclientv1.GroupEndpoints
	}
	type args struct {
		listOptions *v1.ListOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *v1.GroupList
		wantErr bool
	}{
		{
			name: "successful list groups",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			args: args{
				listOptions: &v1.ListOptions{},
			},
			want:    &v1.GroupList{Items: []*v1.Group{{Name: "test-group"}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &groupService{
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

func Test_newGroupService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockGroupEndpoints(ctrl)
	mockEndpointFactory := epclientv1.NewMockFactory(ctrl)
	mockEndpointFactory.EXPECT().Groups().Return(mockEndpoints).AnyTimes()

	ctx := context.Background()
	svcFactory := &factory{ctx: ctx, eps: mockEndpointFactory}

	type args struct {
		factory *factory
	}
	tests := []struct {
		name string
		args args
		want GroupService
	}{
		{
			name: "successful creation",
			args: args{factory: svcFactory},
			want: &groupService{ctx: ctx, eps: mockEndpoints},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newGroupService(tt.args.factory)
			if !reflect.DeepEqual(reflect.TypeOf(got), reflect.TypeOf(tt.want)) {
				t.Errorf("newGroupService() type = %T, want %T", got, tt.want)
			}
		})
	}
}
