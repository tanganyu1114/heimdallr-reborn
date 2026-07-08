package endpoint

import (
	"testing"

	"github.com/tanganyu1114/heimdallr-reborn/api/heimdallr_api/v1"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/internal/pkg/meta/v1"
	"github.com/tanganyu1114/heimdallr-reborn/pkg/client/v1/transport"

	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/pkg/client/v1/model"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"go.uber.org/mock/gomock"
)

func Test_groupEndpoints_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockGroupTransport(ctrl)
	mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.IDOptions, modelclientv1.ResponseBody[*v1.Group]](ctrl)
	mockClient := httpclientv1.NewMockClient[metav1.IDOptions, modelclientv1.ResponseBody[*v1.Group]](ctrl)
	mockEndpoint := httpclientv1.NewEndpoint[metav1.IDOptions, modelclientv1.ResponseBody[*v1.Group]](nil)
	mockClientBuilder.EXPECT().Build().Return(mockClient).AnyTimes()
	mockClient.EXPECT().Endpoint().Return(mockEndpoint).AnyTimes()
	mockTransport.EXPECT().Get().Return(mockClientBuilder).AnyTimes()

	type fields struct {
		transport transport.GroupTransport
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.Endpoint[metav1.IDOptions, modelclientv1.ResponseBody[*v1.Group]]
	}{
		{
			name:   "returns get endpoint",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &groupEndpoints{
				transport: tt.fields.transport,
			}
			got := g.Get()
			if got == nil {
				t.Errorf("Get() = nil, want non-nil")
			}
		})
	}
}

func Test_groupEndpoints_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockGroupTransport(ctrl)
	mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.ListOptions, modelclientv1.ResponseBody[*v1.GroupList]](ctrl)
	mockClient := httpclientv1.NewMockClient[metav1.ListOptions, modelclientv1.ResponseBody[*v1.GroupList]](ctrl)
	mockEndpoint := httpclientv1.NewEndpoint[metav1.ListOptions, modelclientv1.ResponseBody[*v1.GroupList]](nil)
	mockClientBuilder.EXPECT().Build().Return(mockClient).AnyTimes()
	mockClient.EXPECT().Endpoint().Return(mockEndpoint).AnyTimes()
	mockTransport.EXPECT().List().Return(mockClientBuilder).AnyTimes()

	type fields struct {
		transport transport.GroupTransport
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.Endpoint[metav1.ListOptions, modelclientv1.ResponseBody[*v1.GroupList]]
	}{
		{
			name:   "returns list endpoint",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &groupEndpoints{
				transport: tt.fields.transport,
			}
			got := g.List()
			if got == nil {
				t.Errorf("List() = nil, want non-nil")
			}
		})
	}
}

func Test_newGroupEndpoints(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockFactory(ctrl)
	mockGroupTransport := transport.NewMockGroupTransport(ctrl)
	mockTransport.EXPECT().Groups().Return(mockGroupTransport).AnyTimes()

	type args struct {
		f *factory
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "creates group endpoints",
			args: args{
				f: &factory{
					transport: mockTransport,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newGroupEndpoints(tt.args.f)
			if got == nil {
				t.Errorf("newGroupEndpoints() = nil, want non-nil")
			}
		})
	}
}
