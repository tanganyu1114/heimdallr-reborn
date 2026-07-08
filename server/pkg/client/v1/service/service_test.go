package service

import (
	"context"
	"testing"

	epclientv1 "github.com/tanganyu1114/heimdallr-reborn/pkg/client/v1/endpoint"

	"go.uber.org/mock/gomock"
)

func TestNewService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockFactory(ctrl)

	type args struct {
		ctx context.Context
		eps epclientv1.Factory
	}
	tests := []struct {
		name string
		args args
		want Factory
	}{
		{
			name: "creates service with valid context and endpoints",
			args: args{
				ctx: context.Background(),
				eps: mockEndpoints,
			},
		},
		{
			name: "creates service with nil endpoints",
			args: args{
				ctx: context.Background(),
				eps: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewService(tt.args.ctx, tt.args.eps)
			if tt.args.eps == nil {
				if err == nil {
					t.Errorf("NewService() expected error for nil endpoints, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("NewService() unexpected error: %v", err)
			}
			if got == nil {
				t.Errorf("NewService() = nil, want non-nil")
			}
		})
	}
}

func Test_factory_AgentInfos(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockFactory(ctrl)
	mockAgentInfoEndpoints := epclientv1.NewMockAgentInfoEndpoints(ctrl)
	mockEndpoints.EXPECT().AgentInfos().Return(mockAgentInfoEndpoints).AnyTimes()

	type fields struct {
		ctx context.Context
		eps epclientv1.Factory
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "returns agent info service",
			fields: fields{
				ctx: context.Background(),
				eps: mockEndpoints,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &factory{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
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

	mockEndpoints := epclientv1.NewMockFactory(ctrl)
	mockGroupEndpoints := epclientv1.NewMockGroupEndpoints(ctrl)
	mockEndpoints.EXPECT().Groups().Return(mockGroupEndpoints).AnyTimes()

	type fields struct {
		ctx context.Context
		eps epclientv1.Factory
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "returns group service",
			fields: fields{
				ctx: context.Background(),
				eps: mockEndpoints,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &factory{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
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

	mockEndpoints := epclientv1.NewMockFactory(ctrl)
	mockHostEndpoints := epclientv1.NewMockHostEndpoints(ctrl)
	mockEndpoints.EXPECT().Hosts().Return(mockHostEndpoints).AnyTimes()

	type fields struct {
		ctx context.Context
		eps epclientv1.Factory
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "returns host service",
			fields: fields{
				ctx: context.Background(),
				eps: mockEndpoints,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &factory{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			got := f.Hosts()
			if got == nil {
				t.Errorf("Hosts() = nil, want non-nil")
			}
		})
	}
}

func Test_factory_WebServerConfigs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockFactory(ctrl)
	mockWebServerConfigEndpoints := epclientv1.NewMockWebServerConfigEndpoints(ctrl)
	mockEndpoints.EXPECT().WebServerConfigs().Return(mockWebServerConfigEndpoints).AnyTimes()

	type fields struct {
		ctx context.Context
		eps epclientv1.Factory
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "returns web server config service",
			fields: fields{
				ctx: context.Background(),
				eps: mockEndpoints,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &factory{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			got := f.WebServerConfigs()
			if got == nil {
				t.Errorf("WebServerConfigs() = nil, want non-nil")
			}
		})
	}
}

func Test_factory_WebServerBinCMDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockFactory(ctrl)
	mockWebServerBinCMDEndpoints := epclientv1.NewMockWebServerBinCMDEndpoints(ctrl)
	mockEndpoints.EXPECT().WebServerBinCMDs().Return(mockWebServerBinCMDEndpoints).AnyTimes()

	type fields struct {
		ctx context.Context
		eps epclientv1.Factory
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "returns web server bin cmd service",
			fields: fields{
				ctx: context.Background(),
				eps: mockEndpoints,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &factory{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			got := f.WebServerBinCMDs()
			if got == nil {
				t.Errorf("WebServerBinCMDs() = nil, want non-nil")
			}
		})
	}
}

func Test_factory_WebServerStatistics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockFactory(ctrl)
	mockWebServerStatisticsEndpoints := epclientv1.NewMockWebServerStatisticsEndpoints(ctrl)
	mockEndpoints.EXPECT().WebServerStatistics().Return(mockWebServerStatisticsEndpoints).AnyTimes()

	type fields struct {
		ctx context.Context
		eps epclientv1.Factory
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "returns web server statistics service",
			fields: fields{
				ctx: context.Background(),
				eps: mockEndpoints,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &factory{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			got := f.WebServerStatistics()
			if got == nil {
				t.Errorf("WebServerStatistics() = nil, want non-nil")
			}
		})
	}
}
