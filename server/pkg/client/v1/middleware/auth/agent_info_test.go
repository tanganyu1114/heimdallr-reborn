package auth

import (
	"testing"

	v1 "gin-vue-admin/api/heimdallr_api/v1"
	"gin-vue-admin/pkg/client/v1/transport"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_agentInfoMiddleware_Get(t *testing.T) {
	type fields struct {
		md *authMiddleware
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "get agent info middleware",
			fields: fields{
				md: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFactory := transport.NewMockFactory(ctrl)
			mockAgentInfoTransport := transport.NewMockAgentInfoTransport(ctrl)
			mockClientBuilder := httpclientv1.NewMockClientBuilder[httpclientv1.NilBody, []v1.GroupInfo](ctrl)

			mockFactory.EXPECT().AgentInfos().Return(mockAgentInfoTransport)
			mockAgentInfoTransport.EXPECT().Get().Return(mockClientBuilder)
			mockClientBuilder.EXPECT().Use(gomock.Any()).Return(mockClientBuilder)
			mockClientBuilder.EXPECT().WithOptions(gomock.Any()).Return(mockClientBuilder)

			authMw := &authMiddleware{txp: mockFactory}
			a := &agentInfoMiddleware{
				md: authMw,
			}
			got := a.Get()
			assert.NotNil(t, got)
		})
	}
}
