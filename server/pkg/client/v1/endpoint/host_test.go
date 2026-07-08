package endpoint

import (
	"testing"

	"gin-vue-admin/api/heimdallr_api/v1"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	"gin-vue-admin/pkg/client/v1/transport"

	modelclientv1 "gin-vue-admin/pkg/client/v1/model"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"go.uber.org/mock/gomock"
)

func Test_hostEndpoints_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockHostTransport(ctrl)
	mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.IDOptions, modelclientv1.ResponseBody[*v1.Host]](ctrl)
	mockClient := httpclientv1.NewMockClient[metav1.IDOptions, modelclientv1.ResponseBody[*v1.Host]](ctrl)
	mockEndpoint := httpclientv1.NewEndpoint[metav1.IDOptions, modelclientv1.ResponseBody[*v1.Host]](nil)
	mockClientBuilder.EXPECT().Build().Return(mockClient).AnyTimes()
	mockClient.EXPECT().Endpoint().Return(mockEndpoint).AnyTimes()
	mockTransport.EXPECT().Get().Return(mockClientBuilder).AnyTimes()

	type fields struct {
		transport transport.HostTransport
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.Endpoint[metav1.IDOptions, modelclientv1.ResponseBody[*v1.Host]]
	}{
		{
			name:   "returns get endpoint",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &hostEndpoints{
				transport: tt.fields.transport,
			}
			got := h.Get()
			if got == nil {
				t.Errorf("Get() = nil, want non-nil")
			}
		})
	}
}

func Test_hostEndpoints_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockHostTransport(ctrl)
	mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.ListOptions, modelclientv1.ResponseBody[*v1.HostList]](ctrl)
	mockClient := httpclientv1.NewMockClient[metav1.ListOptions, modelclientv1.ResponseBody[*v1.HostList]](ctrl)
	mockEndpoint := httpclientv1.NewEndpoint[metav1.ListOptions, modelclientv1.ResponseBody[*v1.HostList]](nil)
	mockClientBuilder.EXPECT().Build().Return(mockClient).AnyTimes()
	mockClient.EXPECT().Endpoint().Return(mockEndpoint).AnyTimes()
	mockTransport.EXPECT().List().Return(mockClientBuilder).AnyTimes()

	type fields struct {
		transport transport.HostTransport
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.Endpoint[metav1.ListOptions, modelclientv1.ResponseBody[*v1.HostList]]
	}{
		{
			name:   "returns list endpoint",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &hostEndpoints{
				transport: tt.fields.transport,
			}
			got := h.List()
			if got == nil {
				t.Errorf("List() = nil, want non-nil")
			}
		})
	}
}

func Test_newHostEndpoints(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockFactory(ctrl)
	mockHostTransport := transport.NewMockHostTransport(ctrl)
	mockTransport.EXPECT().Hosts().Return(mockHostTransport).AnyTimes()

	type args struct {
		f *factory
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "creates host endpoints",
			args: args{
				f: &factory{
					transport: mockTransport,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newHostEndpoints(tt.args.f)
			if got == nil {
				t.Errorf("newHostEndpoints() = nil, want non-nil")
			}
		})
	}
}
