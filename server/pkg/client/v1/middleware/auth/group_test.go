package auth

import (
	"testing"

	v1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/model"
	"github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/transport"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_groupMiddleware_Get(t *testing.T) {
	type fields struct {
		md *authMiddleware
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "get group middleware",
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
			mockGroupTransport := transport.NewMockGroupTransport(ctrl)
			mockClientBuilder := httpclientv1.NewMockClientBuilder[v1.IDOptions, modelclientv1.ResponseBody[*v1.Group]](ctrl)

			mockFactory.EXPECT().Groups().Return(mockGroupTransport)
			mockGroupTransport.EXPECT().Get().Return(mockClientBuilder)
			mockClientBuilder.EXPECT().Use(gomock.Any()).Return(mockClientBuilder)
			mockClientBuilder.EXPECT().WithOptions(gomock.Any()).Return(mockClientBuilder)

			authMw := &authMiddleware{txp: mockFactory}
			g := &groupMiddleware{
				md: authMw,
			}
			got := g.Get()
			assert.NotNil(t, got)
		})
	}
}

func Test_groupMiddleware_List(t *testing.T) {
	type fields struct {
		md *authMiddleware
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "list group middleware",
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
			mockGroupTransport := transport.NewMockGroupTransport(ctrl)
			mockClientBuilder := httpclientv1.NewMockClientBuilder[v1.ListOptions, modelclientv1.ResponseBody[*v1.GroupList]](ctrl)

			mockFactory.EXPECT().Groups().Return(mockGroupTransport)
			mockGroupTransport.EXPECT().List().Return(mockClientBuilder)
			mockClientBuilder.EXPECT().Use(gomock.Any()).Return(mockClientBuilder)
			mockClientBuilder.EXPECT().WithOptions(gomock.Any()).Return(mockClientBuilder)

			authMw := &authMiddleware{txp: mockFactory}
			g := &groupMiddleware{
				md: authMw,
			}
			got := g.List()
			assert.NotNil(t, got)
		})
	}
}
