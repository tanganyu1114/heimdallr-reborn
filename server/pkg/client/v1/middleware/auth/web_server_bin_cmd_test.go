package auth

import (
	"testing"

	metav1 "github.com/tanganyu1114/heimdallr-reborn/internal/pkg/meta/v1"
	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/pkg/client/v1/model"
	"github.com/tanganyu1114/heimdallr-reborn/pkg/client/v1/transport"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_webServerBinCMDMiddleware_Exec(t *testing.T) {
	type fields struct {
		md *authMiddleware
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "exec web server bin cmd middleware",
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
			mockWebServerBinCMDTransport := transport.NewMockWebServerBinCMDTransport(ctrl)
			mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.WebServerBinCMDExecRequest, modelclientv1.ResponseBody[*metav1.WebServerBinCMDExecResponse]](ctrl)

			mockFactory.EXPECT().WebServerBinCMDs().Return(mockWebServerBinCMDTransport)
			mockWebServerBinCMDTransport.EXPECT().Exec().Return(mockClientBuilder)
			mockClientBuilder.EXPECT().Use(gomock.Any()).Return(mockClientBuilder)
			mockClientBuilder.EXPECT().WithOptions(gomock.Any()).Return(mockClientBuilder)

			authMw := &authMiddleware{txp: mockFactory}
			w := &webServerBinCMDMiddleware{
				md: authMw,
			}
			got := w.Exec()
			assert.NotNil(t, got)
		})
	}
}
