package endpoint

import (
	"testing"

	"gin-vue-admin/pkg/client/v1/transport"

	"go.uber.org/mock/gomock"
)

func TestNewEndpoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockFactory(ctrl)

	type args struct {
		transport transport.Factory
	}
	tests := []struct {
		name string
		args args
		want Factory
	}{
		{
			name: "creates endpoint with valid transport",
			args: args{
				transport: mockTransport,
			},
		},
		{
			name: "creates endpoint with nil transport",
			args: args{
				transport: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEndpoint(tt.args.transport)
			if tt.args.transport == nil {
				if err == nil {
					t.Errorf("NewEndpoint() expected error for nil transport, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("NewEndpoint() unexpected error: %v", err)
			}
			if got == nil {
				t.Errorf("NewEndpoint() = nil, want non-nil")
			}
		})
	}
}

func Test_factory_AgentInfos(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockFactory(ctrl)
	mockAgentTransport := transport.NewMockAgentInfoTransport(ctrl)
	mockTransport.EXPECT().AgentInfos().Return(mockAgentTransport).AnyTimes()

	type fields struct {
		transport transport.Factory
	}
	tests := []struct {
		name   string
		fields fields
		want   AgentInfoEndpoints
	}{
		{
			name:   "returns agent info endpoints",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &factory{
				transport: tt.fields.transport,
			}
			got := f.AgentInfos()
			if got == nil {
				t.Errorf("AgentInfos() = nil, want non-nil")
			}
		})
	}
}

func Test_factory_Groups(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockFactory(ctrl)
	mockGroupTransport := transport.NewMockGroupTransport(ctrl)
	mockTransport.EXPECT().Groups().Return(mockGroupTransport).AnyTimes()

	type fields struct {
		transport transport.Factory
	}
	tests := []struct {
		name   string
		fields fields
		want   GroupEndpoints
	}{
		{
			name:   "returns group endpoints",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &factory{
				transport: tt.fields.transport,
			}
			got := f.Groups()
			if got == nil {
				t.Errorf("Groups() = nil, want non-nil")
			}
		})
	}
}

func Test_factory_Hosts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockFactory(ctrl)
	mockHostTransport := transport.NewMockHostTransport(ctrl)
	mockTransport.EXPECT().Hosts().Return(mockHostTransport).AnyTimes()

	type fields struct {
		transport transport.Factory
	}
	tests := []struct {
		name   string
		fields fields
		want   HostEndpoints
	}{
		{
			name:   "returns host endpoints",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &factory{
				transport: tt.fields.transport,
			}
			got := f.Hosts()
			if got == nil {
				t.Errorf("Hosts() = nil, want non-nil")
			}
		})
	}
}

func Test_factory_WebServerBinCMDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockFactory(ctrl)
	mockWebServerBinCMDTransport := transport.NewMockWebServerBinCMDTransport(ctrl)
	mockTransport.EXPECT().WebServerBinCMDs().Return(mockWebServerBinCMDTransport).AnyTimes()

	type fields struct {
		transport transport.Factory
	}
	tests := []struct {
		name   string
		fields fields
		want   WebServerBinCMDEndpoints
	}{
		{
			name:   "returns web server bin cmd endpoints",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &factory{
				transport: tt.fields.transport,
			}
			got := f.WebServerBinCMDs()
			if got == nil {
				t.Errorf("WebServerBinCMDs() = nil, want non-nil")
			}
		})
	}
}

func Test_factory_WebServerConfigs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockFactory(ctrl)
	mockWebServerConfigTransport := transport.NewMockWebServerConfigTransport(ctrl)
	mockTransport.EXPECT().WebServerConfigs().Return(mockWebServerConfigTransport).AnyTimes()

	type fields struct {
		transport transport.Factory
	}
	tests := []struct {
		name   string
		fields fields
		want   WebServerConfigEndpoints
	}{
		{
			name:   "returns web server config endpoints",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &factory{
				transport: tt.fields.transport,
			}
			got := f.WebServerConfigs()
			if got == nil {
				t.Errorf("WebServerConfigs() = nil, want non-nil")
			}
		})
	}
}

func Test_factory_WebServerStatistics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransport := transport.NewMockFactory(ctrl)
	mockWebServerStatisticsTransport := transport.NewMockWebServerStatisticsTransport(ctrl)
	mockTransport.EXPECT().WebServerStatistics().Return(mockWebServerStatisticsTransport).AnyTimes()

	type fields struct {
		transport transport.Factory
	}
	tests := []struct {
		name   string
		fields fields
		want   WebServerStatisticsEndpoints
	}{
		{
			name:   "returns web server statistics endpoints",
			fields: fields{transport: mockTransport},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &factory{
				transport: tt.fields.transport,
			}
			got := f.WebServerStatistics()
			if got == nil {
				t.Errorf("WebServerStatistics() = nil, want non-nil")
			}
		})
	}
}
