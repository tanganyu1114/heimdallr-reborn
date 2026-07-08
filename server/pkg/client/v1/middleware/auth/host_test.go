package auth

import (
	"testing"

	v1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/pkg/meta/v1"
	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/model"
	"github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/transport"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_hostMiddleware_Get(t *testing.T) {
	type fields struct {
		md *authMiddleware
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "get host middleware",
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
			mockHostTransport := transport.NewMockHostTransport(ctrl)
			mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.IDOptions, modelclientv1.ResponseBody[*v1.Host]](ctrl)

			mockFactory.EXPECT().Hosts().Return(mockHostTransport)
			mockHostTransport.EXPECT().Get().Return(mockClientBuilder)
			mockClientBuilder.EXPECT().Use(gomock.Any()).Return(mockClientBuilder)
			mockClientBuilder.EXPECT().WithOptions(gomock.Any()).Return(mockClientBuilder)

			authMw := &authMiddleware{txp: mockFactory}
			h := &hostMiddleware{
				md: authMw,
			}
			got := h.Get()
			assert.NotNil(t, got)
		})
	}
}

func Test_hostMiddleware_List(t *testing.T) {
	type fields struct {
		md *authMiddleware
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "list host middleware",
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
			mockHostTransport := transport.NewMockHostTransport(ctrl)
			mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.ListOptions, modelclientv1.ResponseBody[*v1.HostList]](ctrl)

			mockFactory.EXPECT().Hosts().Return(mockHostTransport)
			mockHostTransport.EXPECT().List().Return(mockClientBuilder)
			mockClientBuilder.EXPECT().Use(gomock.Any()).Return(mockClientBuilder)
			mockClientBuilder.EXPECT().WithOptions(gomock.Any()).Return(mockClientBuilder)

			authMw := &authMiddleware{txp: mockFactory}
			h := &hostMiddleware{
				md: authMw,
			}
			got := h.List()
			assert.NotNil(t, got)
		})
	}
}
