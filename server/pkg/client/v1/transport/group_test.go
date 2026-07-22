package transport

import (
	"reflect"
	"testing"

	v1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/model"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"go.uber.org/mock/gomock"
)

func Test_groupTransport_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[v1.IDOptions, modelclientv1.ResponseBody[*v1.Group]](ctrl)

	type fields struct {
		getGroupClient   httpclientv1.ClientBuilder[v1.IDOptions, modelclientv1.ResponseBody[*v1.Group]]
		listGroupsClient httpclientv1.ClientBuilder[v1.ListOptions, modelclientv1.ResponseBody[*v1.GroupList]]
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[v1.IDOptions, modelclientv1.ResponseBody[*v1.Group]]
	}{
		{
			name: "returns get group client",
			fields: fields{
				getGroupClient: mockClient,
			},
			want: mockClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &groupTransport{
				getGroupClient:   tt.fields.getGroupClient,
				listGroupsClient: tt.fields.listGroupsClient,
			}
			if got := g.Get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_groupTransport_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[v1.ListOptions, modelclientv1.ResponseBody[*v1.GroupList]](ctrl)

	type fields struct {
		getGroupClient   httpclientv1.ClientBuilder[v1.IDOptions, modelclientv1.ResponseBody[*v1.Group]]
		listGroupsClient httpclientv1.ClientBuilder[v1.ListOptions, modelclientv1.ResponseBody[*v1.GroupList]]
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[v1.ListOptions, modelclientv1.ResponseBody[*v1.GroupList]]
	}{
		{
			name: "returns list groups client",
			fields: fields{
				listGroupsClient: mockClient,
			},
			want: mockClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &groupTransport{
				getGroupClient:   tt.fields.getGroupClient,
				listGroupsClient: tt.fields.listGroupsClient,
			}
			if got := g.List(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newGroupTransport(t *testing.T) {
	type args struct {
		transport *transport
	}
	tests := []struct {
		name string
		args args
		want GroupTransport
	}{
		{
			name: "creates group transport with valid base URL",
			args: args{
				transport: &transport{
					baseURL: "http://localhost:8080",
				},
			},
		},
		{
			name: "creates group transport with empty base URL",
			args: args{
				transport: &transport{
					baseURL: "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newGroupTransport(tt.args.transport)
			if got == nil {
				t.Errorf("newGroupTransport() = nil, want non-nil")
			}
		})
	}
}
