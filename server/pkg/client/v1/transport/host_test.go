package transport

import (
	"reflect"
	"testing"

	v1 "github.com/tanganyu1114/heimdallr-reborn/api/heimdallr_api/v1"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/internal/pkg/meta/v1"
	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/pkg/client/v1/model"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"go.uber.org/mock/gomock"
)

func Test_hostTransport_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[metav1.IDOptions, modelclientv1.ResponseBody[*v1.Host]](ctrl)

	type fields struct {
		getHostClient   httpclientv1.ClientBuilder[metav1.IDOptions, modelclientv1.ResponseBody[*v1.Host]]
		listHostsClient httpclientv1.ClientBuilder[metav1.ListOptions, modelclientv1.ResponseBody[*v1.HostList]]
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[metav1.IDOptions, modelclientv1.ResponseBody[*v1.Host]]
	}{
		{
			name: "returns get host client",
			fields: fields{
				getHostClient: mockClient,
			},
			want: mockClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &hostTransport{
				getHostClient:   tt.fields.getHostClient,
				listHostsClient: tt.fields.listHostsClient,
			}
			if got := h.Get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hostTransport_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[metav1.ListOptions, modelclientv1.ResponseBody[*v1.HostList]](ctrl)

	type fields struct {
		getHostClient   httpclientv1.ClientBuilder[metav1.IDOptions, modelclientv1.ResponseBody[*v1.Host]]
		listHostsClient httpclientv1.ClientBuilder[metav1.ListOptions, modelclientv1.ResponseBody[*v1.HostList]]
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[metav1.ListOptions, modelclientv1.ResponseBody[*v1.HostList]]
	}{
		{
			name: "returns list hosts client",
			fields: fields{
				listHostsClient: mockClient,
			},
			want: mockClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &hostTransport{
				getHostClient:   tt.fields.getHostClient,
				listHostsClient: tt.fields.listHostsClient,
			}
			if got := h.List(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newHostTransport(t *testing.T) {
	type args struct {
		transport *transport
	}
	tests := []struct {
		name string
		args args
		want HostTransport
	}{
		{
			name: "creates host transport with valid base URL",
			args: args{
				transport: &transport{
					baseURL: "http://localhost:8080",
				},
			},
		},
		{
			name: "creates host transport with empty base URL",
			args: args{
				transport: &transport{
					baseURL: "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newHostTransport(tt.args.transport)
			if got == nil {
				t.Errorf("newHostTransport() = nil, want non-nil")
			}
		})
	}
}
