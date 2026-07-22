package service

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	v1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"

	epclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/endpoint"
	modelclientv1 "github.com/tanganyu1114/heimdallr-reborn/server/pkg/client/v1/model"

	utilsV3 "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/utils"
	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
	"github.com/marmotedu/errors"
	"go.uber.org/mock/gomock"
)

func Test_newWebServerConfigService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockWebServerConfigEndpoints(ctrl)
	mockEndpointFactory := epclientv1.NewMockFactory(ctrl)
	mockEndpointFactory.EXPECT().WebServerConfigs().Return(mockEndpoints).AnyTimes()

	ctx := context.Background()
	svcFactory := &factory{ctx: ctx, eps: mockEndpointFactory}

	type args struct {
		factory *factory
	}
	tests := []struct {
		name string
		args args
		want WebServerConfigService
	}{
		{
			name: "successful creation",
			args: args{factory: svcFactory},
			want: &webServerConfigService{ctx: ctx, eps: mockEndpoints},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newWebServerConfigService(tt.args.factory)
			if !reflect.DeepEqual(reflect.TypeOf(got), reflect.TypeOf(tt.want)) {
				t.Errorf("newWebServerConfigService() type = %T, want %T", got, tt.want)
			}
		})
	}
}

func Test_webServerConfigService_ChangeContextEnabledState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockWebServerConfigEndpoints(ctrl)
	ctx := context.Background()

	mockEndpoint := httpclientv1.NewEndpoint[v1.WebServerConfigContextUpdateOptions[v1.ConfigContextEnabledStateMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]](func(ctx context.Context, req interface{}) (interface{}, error) {
		return modelclientv1.ResponseBody[httpclientv1.NilBody]{}, nil
	})
	mockEndpoints.EXPECT().ChangeContextEnabledState().Return(mockEndpoint).AnyTimes()

	type fields struct {
		ctx context.Context
		eps epclientv1.WebServerConfigEndpoints
	}
	type args struct {
		opts *v1.WebServerConfigContextUpdateOptions[v1.ConfigContextEnabledStateMeta]
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "successful change context enabled state",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			args: args{
				opts: &v1.WebServerConfigContextUpdateOptions[v1.ConfigContextEnabledStateMeta]{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &webServerConfigService{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			if err := s.ChangeContextEnabledState(tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("ChangeContextEnabledState() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigService_GetConfig(t *testing.T) {
	ctx := context.Background()

	type fields struct {
		ctx context.Context
		eps epclientv1.WebServerConfigEndpoints
	}
	type args struct {
		opts *v1.WebServerOptions
	}
	tests := []struct {
		name             string
		fields           fields
		args             args
		setupMock        func(ctrl *gomock.Controller) epclientv1.WebServerConfigEndpoints
		wantConfig       *v1.WebServerConfig
		wantFingerprints utilsV3.ConfigFingerprints
		wantErr          bool
	}{
		{
			name: "opts is nil",
			fields: fields{
				ctx: ctx,
			},
			args: args{
				opts: nil,
			},
			setupMock: func(ctrl *gomock.Controller) epclientv1.WebServerConfigEndpoints {
				return epclientv1.NewMockWebServerConfigEndpoints(ctrl)
			},
			wantConfig:       nil,
			wantFingerprints: utilsV3.ConfigFingerprints{},
			wantErr:          true,
		},
		{
			name: "endpoint call error",
			fields: fields{
				ctx: ctx,
			},
			args: args{
				opts: &v1.WebServerOptions{},
			},
			setupMock: func(ctrl *gomock.Controller) epclientv1.WebServerConfigEndpoints {
				mockEndpoints := epclientv1.NewMockWebServerConfigEndpoints(ctrl)
				mockErrorEndpoint := httpclientv1.NewEndpoint[v1.WebServerOptions, modelclientv1.ResponseBody[*modelclientv1.WebServerConfig]](func(ctx context.Context, req interface{}) (interface{}, error) {
					return nil, errors.New("endpoint call error")
				})
				mockEndpoints.EXPECT().GetConfig().Return(mockErrorEndpoint)
				return mockEndpoints
			},
			wantConfig:       nil,
			wantFingerprints: utilsV3.ConfigFingerprints{},
			wantErr:          true,
		},
		{
			name: "invalid response data",
			fields: fields{
				ctx: ctx,
			},
			args: args{
				opts: &v1.WebServerOptions{},
			},
			setupMock: func(ctrl *gomock.Controller) epclientv1.WebServerConfigEndpoints {
				mockEndpoints := epclientv1.NewMockWebServerConfigEndpoints(ctrl)
				mockInvalidResponseEndpoint := httpclientv1.NewEndpoint[v1.WebServerOptions, modelclientv1.ResponseBody[*modelclientv1.WebServerConfig]](func(ctx context.Context, req interface{}) (interface{}, error) {
					return modelclientv1.ResponseBody[*modelclientv1.WebServerConfig]{Data: json.RawMessage(`invalid json`)}, nil
				})
				mockEndpoints.EXPECT().GetConfig().Return(mockInvalidResponseEndpoint)
				return mockEndpoints
			},
			wantConfig:       nil,
			wantFingerprints: utilsV3.ConfigFingerprints{},
			wantErr:          true,
		},
		{
			name: "successful get config with fingerprints",
			fields: fields{
				ctx: ctx,
			},
			args: args{
				opts: &v1.WebServerOptions{},
			},
			setupMock: func(ctrl *gomock.Controller) epclientv1.WebServerConfigEndpoints {
				mockEndpoints := epclientv1.NewMockWebServerConfigEndpoints(ctrl)
				mockEndpoint := httpclientv1.NewEndpoint[v1.WebServerOptions, modelclientv1.ResponseBody[*modelclientv1.WebServerConfig]](func(ctx context.Context, req interface{}) (interface{}, error) {
					nginxConfigJSON := `{
						"main-config": "E:\\config_test\\nginx.conf",
						"configs": {
							"E:\\config_test\\nginx.conf": {
								"enabled": true,
								"context-type": "config",
								"value": "E:\\config_test\\nginx.conf",
								"params": [
									{"context-type": "directive", "value": "worker_processes 1"},
									{"context-type": "events", "params": [{"context-type": "directive", "value": "worker_connections 1024"}]},
									{"context-type": "http", "params": [{"context-type": "directive", "value": "sendfile on"}]}
								]
							}
						}
					}`
					config := modelclientv1.WebServerConfig{
						Config: json.RawMessage(nginxConfigJSON),
						OriginalFingerprints: utilsV3.ConfigFingerprints{
							"E:\\config_test\\nginx.conf": "main-config-hash-123",
						},
					}
					respData, _ := json.Marshal(&config)
					return modelclientv1.ResponseBody[*modelclientv1.WebServerConfig]{Data: respData}, nil
				})
				mockEndpoints.EXPECT().GetConfig().Return(mockEndpoint)
				return mockEndpoints
			},
			wantConfig: &v1.WebServerConfig{},
			wantFingerprints: utilsV3.ConfigFingerprints{
				"E:\\config_test\\nginx.conf": "main-config-hash-123",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tt.fields.eps = tt.setupMock(ctrl)
			s := &webServerConfigService{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			gotConfig, gotFingerprints, err := s.GetConfig(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotConfig == nil && tt.wantConfig != nil {
				t.Errorf("GetConfig() got = nil, want %v", tt.wantConfig)
			}
			if !reflect.DeepEqual(gotFingerprints, tt.wantFingerprints) {
				t.Errorf("GetConfig() fingerprints = %v, want %v", gotFingerprints, tt.wantFingerprints)
			}
		})
	}
}

func Test_webServerConfigService_GetConfigTextLines(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockWebServerConfigEndpoints(ctrl)
	ctx := context.Background()

	mockEndpoint := httpclientv1.NewEndpoint[v1.WebServerOptions, modelclientv1.ResponseBody[string]](func(ctx context.Context, req interface{}) (interface{}, error) {
		data, _ := json.Marshal("config text lines")
		return modelclientv1.ResponseBody[string]{Data: data}, nil
	})
	mockEndpoints.EXPECT().GetConfigTextLines().Return(mockEndpoint).AnyTimes()

	type fields struct {
		ctx context.Context
		eps epclientv1.WebServerConfigEndpoints
	}
	type args struct {
		opts *v1.WebServerOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "successful get config text lines",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			args: args{
				opts: &v1.WebServerOptions{},
			},
			want:    "config text lines",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &webServerConfigService{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			got, err := s.GetConfigTextLines(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfigTextLines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetConfigTextLines() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webServerConfigService_GetContextTextLines(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockWebServerConfigEndpoints(ctrl)
	ctx := context.Background()

	mockEndpoint := httpclientv1.NewEndpoint[v1.WebServerConfigTargetContextOptions, modelclientv1.ResponseBody[string]](func(ctx context.Context, req interface{}) (interface{}, error) {
		data, _ := json.Marshal("context text lines")
		return modelclientv1.ResponseBody[string]{Data: data}, nil
	})
	mockEndpoints.EXPECT().GetContextTextLines().Return(mockEndpoint).AnyTimes()

	type fields struct {
		ctx context.Context
		eps epclientv1.WebServerConfigEndpoints
	}
	type args struct {
		opts *v1.WebServerConfigTargetContextOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "successful get context text lines",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			args: args{
				opts: &v1.WebServerConfigTargetContextOptions{},
			},
			want:    "context text lines",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &webServerConfigService{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			got, err := s.GetContextTextLines(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContextTextLines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetContextTextLines() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webServerConfigService_GetIncludedConfigs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockWebServerConfigEndpoints(ctrl)
	ctx := context.Background()

	mockEndpoint := httpclientv1.NewEndpoint[v1.WebServerConfigTargetContextOptions, modelclientv1.ResponseBody[[]string]](func(ctx context.Context, req interface{}) (interface{}, error) {
		data, _ := json.Marshal([]string{"config1.conf", "config2.conf"})
		return modelclientv1.ResponseBody[[]string]{Data: data}, nil
	})
	mockEndpoints.EXPECT().GetIncludedConfigs().Return(mockEndpoint).AnyTimes()

	type fields struct {
		ctx context.Context
		eps epclientv1.WebServerConfigEndpoints
	}
	type args struct {
		opts *v1.WebServerConfigTargetContextOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "successful get included configs",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			args: args{
				opts: &v1.WebServerConfigTargetContextOptions{},
			},
			want:    []string{"config1.conf", "config2.conf"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &webServerConfigService{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			got, err := s.GetIncludedConfigs(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIncludedConfigs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("GetIncludedConfigs() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webServerConfigService_GetOptions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockWebServerConfigEndpoints(ctrl)
	ctx := context.Background()

	mockEndpoint := httpclientv1.NewEndpoint[httpclientv1.NilBody, modelclientv1.ResponseBody[[]v1.BifrostGroupMeta]](func(ctx context.Context, req interface{}) (interface{}, error) {
		data, _ := json.Marshal([]v1.BifrostGroupMeta{{UintObjectMeta: v1.UintObjectMeta{Label: "test-group", Value: 1}}})
		return modelclientv1.ResponseBody[[]v1.BifrostGroupMeta]{Data: data}, nil
	})
	mockEndpoints.EXPECT().GetOptions().Return(mockEndpoint).AnyTimes()

	type fields struct {
		ctx context.Context
		eps epclientv1.WebServerConfigEndpoints
	}
	tests := []struct {
		name    string
		fields  fields
		want    []v1.BifrostGroupMeta
		wantErr bool
	}{
		{
			name: "successful get options",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			want:    []v1.BifrostGroupMeta{{UintObjectMeta: v1.UintObjectMeta{Label: "test-group", Value: 1}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &webServerConfigService{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			got, err := s.GetOptions()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil && tt.want != nil {
				t.Errorf("GetOptions() got = nil, want %v", tt.want)
			}
		})
	}
}

func Test_webServerConfigService_InsertWithClone(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockWebServerConfigEndpoints(ctrl)
	ctx := context.Background()

	mockEndpoint := httpclientv1.NewEndpoint[v1.WebServerConfigContextUpdateOptions[v1.CloneConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]](func(ctx context.Context, req interface{}) (interface{}, error) {
		return modelclientv1.ResponseBody[httpclientv1.NilBody]{}, nil
	})
	mockEndpoints.EXPECT().InsertWithClone().Return(mockEndpoint).AnyTimes()

	type fields struct {
		ctx context.Context
		eps epclientv1.WebServerConfigEndpoints
	}
	type args struct {
		opts *v1.WebServerConfigContextUpdateOptions[v1.CloneConfigContextMeta]
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "successful insert with clone",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			args: args{
				opts: &v1.WebServerConfigContextUpdateOptions[v1.CloneConfigContextMeta]{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &webServerConfigService{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			if err := s.InsertWithClone(tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("InsertWithClone() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigService_InsertWithNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockWebServerConfigEndpoints(ctrl)
	ctx := context.Background()

	mockEndpoint := httpclientv1.NewEndpoint[v1.WebServerConfigContextUpdateOptions[v1.NewConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]](func(ctx context.Context, req interface{}) (interface{}, error) {
		return modelclientv1.ResponseBody[httpclientv1.NilBody]{}, nil
	})
	mockEndpoints.EXPECT().InsertWithNew().Return(mockEndpoint)

	type fields struct {
		ctx context.Context
		eps epclientv1.WebServerConfigEndpoints
	}
	type args struct {
		opts *v1.WebServerConfigContextUpdateOptions[v1.NewConfigContextMeta]
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "successful insert with new",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			args: args{
				opts: &v1.WebServerConfigContextUpdateOptions[v1.NewConfigContextMeta]{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &webServerConfigService{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			if err := s.InsertWithNew(tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("InsertWithNew() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigService_ModifyContextValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockWebServerConfigEndpoints(ctrl)
	ctx := context.Background()

	mockEndpoint := httpclientv1.NewEndpoint[v1.WebServerConfigContextUpdateOptions[v1.NewConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]](func(ctx context.Context, req interface{}) (interface{}, error) {
		return modelclientv1.ResponseBody[httpclientv1.NilBody]{}, nil
	})
	mockEndpoints.EXPECT().ModifyContextValue().Return(mockEndpoint)

	type fields struct {
		ctx context.Context
		eps epclientv1.WebServerConfigEndpoints
	}
	type args struct {
		opts *v1.WebServerConfigContextUpdateOptions[v1.NewConfigContextMeta]
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "successful modify context value",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			args: args{
				opts: &v1.WebServerConfigContextUpdateOptions[v1.NewConfigContextMeta]{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &webServerConfigService{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			if err := s.ModifyContextValue(tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("ModifyContextValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigService_ModifyWithClone(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockWebServerConfigEndpoints(ctrl)
	ctx := context.Background()

	mockEndpoint := httpclientv1.NewEndpoint[v1.WebServerConfigContextUpdateOptions[v1.CloneConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]](func(ctx context.Context, req interface{}) (interface{}, error) {
		return modelclientv1.ResponseBody[httpclientv1.NilBody]{}, nil
	})
	mockEndpoints.EXPECT().ModifyWithClone().Return(mockEndpoint)

	type fields struct {
		ctx context.Context
		eps epclientv1.WebServerConfigEndpoints
	}
	type args struct {
		opts *v1.WebServerConfigContextUpdateOptions[v1.CloneConfigContextMeta]
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "successful modify with clone",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			args: args{
				opts: &v1.WebServerConfigContextUpdateOptions[v1.CloneConfigContextMeta]{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &webServerConfigService{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			if err := s.ModifyWithClone(tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("ModifyWithClone() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigService_ModifyWithNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockWebServerConfigEndpoints(ctrl)
	ctx := context.Background()

	mockEndpoint := httpclientv1.NewEndpoint[v1.WebServerConfigContextUpdateOptions[v1.NewConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]](func(ctx context.Context, req interface{}) (interface{}, error) {
		return modelclientv1.ResponseBody[httpclientv1.NilBody]{}, nil
	})
	mockEndpoints.EXPECT().ModifyWithNew().Return(mockEndpoint)

	type fields struct {
		ctx context.Context
		eps epclientv1.WebServerConfigEndpoints
	}
	type args struct {
		opts *v1.WebServerConfigContextUpdateOptions[v1.NewConfigContextMeta]
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "successful modify with new",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			args: args{
				opts: &v1.WebServerConfigContextUpdateOptions[v1.NewConfigContextMeta]{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &webServerConfigService{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			if err := s.ModifyWithNew(tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("ModifyWithNew() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigService_Move(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockWebServerConfigEndpoints(ctrl)
	ctx := context.Background()

	mockEndpoint := httpclientv1.NewEndpoint[v1.WebServerConfigContextUpdateOptions[v1.CloneConfigContextMeta], modelclientv1.ResponseBody[httpclientv1.NilBody]](func(ctx context.Context, req interface{}) (interface{}, error) {
		return modelclientv1.ResponseBody[httpclientv1.NilBody]{}, nil
	})
	mockEndpoints.EXPECT().Move().Return(mockEndpoint)

	type fields struct {
		ctx context.Context
		eps epclientv1.WebServerConfigEndpoints
	}
	type args struct {
		opts *v1.WebServerConfigContextUpdateOptions[v1.CloneConfigContextMeta]
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "successful move config",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			args: args{
				opts: &v1.WebServerConfigContextUpdateOptions[v1.CloneConfigContextMeta]{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &webServerConfigService{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			if err := s.Move(tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("Move() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigService_Remove(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockWebServerConfigEndpoints(ctrl)
	ctx := context.Background()

	mockEndpoint := httpclientv1.NewEndpoint[v1.WebServerConfigTargetContextOptions, modelclientv1.ResponseBody[httpclientv1.NilBody]](func(ctx context.Context, req interface{}) (interface{}, error) {
		return modelclientv1.ResponseBody[httpclientv1.NilBody]{}, nil
	})
	mockEndpoints.EXPECT().Remove().Return(mockEndpoint)

	type fields struct {
		ctx context.Context
		eps epclientv1.WebServerConfigEndpoints
	}
	type args struct {
		opts *v1.WebServerConfigTargetContextOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "successful remove config",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			args: args{
				opts: &v1.WebServerConfigTargetContextOptions{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &webServerConfigService{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			if err := s.Remove(tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webServerConfigService_SearchContextPositions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockWebServerConfigEndpoints(ctrl)
	ctx := context.Background()

	mockEndpoint := httpclientv1.NewEndpoint[v1.WebServerConfigContextPosSearchOptions, modelclientv1.ResponseBody[[]v1.ConfigContextPos]](func(ctx context.Context, req interface{}) (interface{}, error) {
		data, _ := json.Marshal([]v1.ConfigContextPos{{Config: "test.conf", ContextPosPath: []int{0, 1}}})
		return modelclientv1.ResponseBody[[]v1.ConfigContextPos]{Data: data}, nil
	})
	mockEndpoints.EXPECT().SearchContextPositions().Return(mockEndpoint).AnyTimes()

	type fields struct {
		ctx context.Context
		eps epclientv1.WebServerConfigEndpoints
	}
	type args struct {
		opts *v1.WebServerConfigContextPosSearchOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []v1.ConfigContextPos
		wantErr bool
	}{
		{
			name: "successful search context positions",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			args: args{
				opts: &v1.WebServerConfigContextPosSearchOptions{},
			},
			want:    []v1.ConfigContextPos{{Config: "test.conf", ContextPosPath: []int{0, 1}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &webServerConfigService{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			got, err := s.SearchContextPositions(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchContextPositions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("SearchContextPositions() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webServerConfigService_UpdateConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockWebServerConfigEndpoints(ctrl)
	ctx := context.Background()

	mockEndpoint := httpclientv1.NewEndpoint[*v1.WebServerConfigUpdateOptions, modelclientv1.ResponseBody[httpclientv1.NilBody]](func(ctx context.Context, req interface{}) (interface{}, error) {
		return modelclientv1.ResponseBody[httpclientv1.NilBody]{}, nil
	})
	mockEndpoints.EXPECT().UpdateConfig().Return(mockEndpoint).AnyTimes()

	type fields struct {
		ctx context.Context
		eps epclientv1.WebServerConfigEndpoints
	}
	type args struct {
		opts *v1.WebServerConfigUpdateOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "successful update config",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			args: args{
				opts: &v1.WebServerConfigUpdateOptions{},
			},
			wantErr: false,
		},
		{
			name: "nil opts returns error",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			args: args{
				opts: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &webServerConfigService{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			if err := s.UpdateConfig(tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("UpdateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
