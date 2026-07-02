package endpoint

import (
	"testing"

	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	"gin-vue-admin/pkg/client/v1/transport"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"go.uber.org/mock/gomock"
)

func Test_newWebServerBinCMDEndpoints(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockFactory(ctrl)
	mockWebServerBinCMDTransport := transport.NewMockWebServerBinCMDTransport(ctrl)
	mockTransport.EXPECT().WebServerBinCMDs().Return(mockWebServerBinCMDTransport).AnyTimes()

	type args struct {
		f *factory
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "creates web server bin cmd endpoints",
			args: args{
				f: &factory{
					transport: mockTransport,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newWebServerBinCMDEndpoints(tt.args.f)
			if got == nil {
				t.Errorf("newWebServerBinCMDEndpoints() = nil, want non-nil")
			}
		})
	}
}

func Test_webServerBinCMDEndpoints_Exec(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockWebServerBinCMDTransport(ctrl)
	mockClientBuilder := httpclientv1.NewMockClientBuilder[metav1.WebServerBinCMDExecRequest, *metav1.WebServerBinCMDExecResponse](ctrl)
	mockClient := httpclientv1.NewMockClient[metav1.WebServerBinCMDExecRequest, *metav1.WebServerBinCMDExecResponse](ctrl)
	mockEndpoint := httpclientv1.NewEndpoint[metav1.WebServerBinCMDExecRequest, *metav1.WebServerBinCMDExecResponse](nil)
	mockClientBuilder.EXPECT().Build().Return(mockClient).AnyTimes()
	mockClient.EXPECT().Endpoint().Return(mockEndpoint).AnyTimes()
	mockTransport.EXPECT().Exec().Return(mockClientBuilder).AnyTimes()

	type fields struct {
		transport transport.WebServerBinCMDTransport
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.Endpoint[metav1.WebServerBinCMDExecRequest, *metav1.WebServerBinCMDExecResponse]
	}{
		{
			name:   "returns exec endpoint",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerBinCMDEndpoints{
				transport: tt.fields.transport,
			}
			got := w.Exec()
			if got == nil {
				t.Errorf("Exec() = nil, want non-nil")
			}
		})
	}
}
