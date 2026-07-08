package transport

import (
	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/model"
	"reflect"
	"testing"

	v1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"go.uber.org/mock/gomock"
)

func Test_agentInfoTransport_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := httpclientv1.NewMockClientBuilder[httpclientv1.NilBody, modelclientv1.ResponseBody[[]v1.GroupInfo]](ctrl)

	type fields struct {
		getAgentInfoClient httpclientv1.ClientBuilder[httpclientv1.NilBody, modelclientv1.ResponseBody[[]v1.GroupInfo]]
	}
	tests := []struct {
		name   string
		fields fields
		want   httpclientv1.ClientBuilder[httpclientv1.NilBody, modelclientv1.ResponseBody[[]v1.GroupInfo]]
	}{
		{
			name: "returns get agent info client",
			fields: fields{
				getAgentInfoClient: mockClient,
			},
			want: mockClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &agentInfoTransport{
				getAgentInfoClient: tt.fields.getAgentInfoClient,
			}
			if got := a.Get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newAgentInfoTransport(t *testing.T) {
	type args struct {
		transport *transport
	}
	tests := []struct {
		name string
		args args
		want AgentInfoTransport
	}{
		{
			name: "creates agent info transport with valid base URL",
			args: args{
				transport: &transport{
					baseURL: "http://localhost:8080",
				},
			},
		},
		{
			name: "creates agent info transport with empty base URL",
			args: args{
				transport: &transport{
					baseURL: "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newAgentInfoTransport(tt.args.transport)
			if got == nil {
				t.Errorf("newAgentInfoTransport() = nil, want non-nil")
			}
		})
	}
}
