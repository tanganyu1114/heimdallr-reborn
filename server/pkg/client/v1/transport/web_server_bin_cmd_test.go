package transport

import (
	"reflect"
	"testing"

	metav1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/model"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"go.uber.org/mock/gomock"
)

func Test_newWebServerBinCMDTransport(t *testing.T) {
	type args struct {
		transport *transport
	}
	tests := []struct {
		name string
		args args
		want WebServerBinCMDTransport
	}{
		{
			name: "creates web server bin cmd transport with valid base URL",
			args: args{
				transport: &transport{
					baseURL: "http://localhost:8080",
				},
			},
		},
		{
			name: "creates web server bin cmd transport with empty base URL",
			args: args{
				transport: &transport{
					baseURL: "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newWebServerBinCMDTransport(tt.args.transport)
			if got == nil {
				t.Errorf("newWebServerBinCMDTransport() = nil, want non-nil")
			}
		})
	}
}

func Test_webServerBinCMDTransport_Exec(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[metav1.WebServerBinCMDExecRequest, modelclientv1.ResponseBody[*metav1.WebServerBinCMDExecResponse]](ctrl)

	type fields struct {
		execClient httpclientv1.ClientBuilder[metav1.WebServerBinCMDExecRequest, modelclientv1.ResponseBody[*metav1.WebServerBinCMDExecResponse]]
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[metav1.WebServerBinCMDExecRequest, modelclientv1.ResponseBody[*metav1.WebServerBinCMDExecResponse]]
	}{
		{
			name: "returns exec client",
			fields: fields{
				execClient: mockClient,
			},
			want: mockClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerBinCMDTransport{
				execClient: tt.fields.execClient,
			}
			if got := w.Exec(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Exec() = %v, want %v", got, tt.want)
			}
		})
	}
}
