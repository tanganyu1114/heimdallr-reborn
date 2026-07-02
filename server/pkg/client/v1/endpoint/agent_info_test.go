package endpoint

import (
	"testing"

	"gin-vue-admin/api/heimdallr_api/v1"
	"gin-vue-admin/pkg/client/v1/transport"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"go.uber.org/mock/gomock"
)

func Test_agentInfoEndpoints_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockAgentInfoTransport(ctrl)
	mockClientBuilder := httpclientv1.NewMockClientBuilder[httpclientv1.NilBody, []v1.GroupInfo](ctrl)
	mockClient := httpclientv1.NewMockClient[httpclientv1.NilBody, []v1.GroupInfo](ctrl)
	mockEndpoint := httpclientv1.NewEndpoint[httpclientv1.NilBody, []v1.GroupInfo](nil)
	mockClientBuilder.EXPECT().Build().Return(mockClient).AnyTimes()
	mockClient.EXPECT().Endpoint().Return(mockEndpoint).AnyTimes()
	mockTransport.EXPECT().Get().Return(mockClientBuilder).AnyTimes()

	type fields struct {
		transport transport.AgentInfoTransport
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.Endpoint[httpclientv1.NilBody, []v1.GroupInfo]
	}{
		{
			name:   "returns get endpoint",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &agentInfoEndpoints{
				transport: tt.fields.transport,
			}
			got := a.Get()
			if got == nil {
				t.Errorf("Get() = nil, want non-nil")
			}
		})
	}
}

func Test_newAgentInfoEndpoints(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockFactory(ctrl)
	mockAgentTransport := transport.NewMockAgentInfoTransport(ctrl)
	mockTransport.EXPECT().AgentInfos().Return(mockAgentTransport).AnyTimes()

	type args struct {
		f *factory
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "creates agent info endpoints",
			args: args{
				f: &factory{
					transport: mockTransport,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newAgentInfoEndpoints(tt.args.f)
			if got == nil {
				t.Errorf("newAgentInfoEndpoints() = nil, want non-nil")
			}
		})
	}
}
