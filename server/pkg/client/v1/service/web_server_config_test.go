package service

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	"reflect"
	"testing"

	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	epclientv1 "gin-vue-admin/pkg/client/v1/endpoint"

	httpclientv1 "github.com/ClessLi/component-base/pkg/client-sdk/http/v1"
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

	mockEndpoint := httpclientv1.NewEndpoint[metav1.WebServerConfigContextUpdateOptions[metav1.ConfigContextEnabledStateMeta], httpclientv1.NilBody](func(ctx context.Context, req interface{}) (interface{}, error) {
		return httpclientv1.NilBody{}, nil
	})
	mockEndpoints.EXPECT().ChangeContextEnabledState().Return(mockEndpoint)

	type fields struct {
		ctx context.Context
		eps epclientv1.WebServerConfigEndpoints
	}
	type args struct {
		opts *metav1.WebServerConfigContextUpdateOptions[metav1.ConfigContextEnabledStateMeta]
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
				opts: &metav1.WebServerConfigContextUpdateOptions[metav1.ConfigContextEnabledStateMeta]{},
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockWebServerConfigEndpoints(ctrl)
	ctx := context.Background()

	mockEndpoint := httpclientv1.NewEndpoint[metav1.WebServerOptions, *metav1.WebServerConfig](func(ctx context.Context, req interface{}) (interface{}, error) {
		return &metav1.WebServerConfig{}, nil
	})
	mockEndpoints.EXPECT().GetConfig().Return(mockEndpoint)

	type fields struct {
		ctx context.Context
		eps epclientv1.WebServerConfigEndpoints
	}
	type args struct {
		opts *metav1.WebServerOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *metav1.WebServerConfig
		wantErr bool
	}{
		{
			name: "successful get config",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			args: args{
				opts: &metav1.WebServerOptions{},
			},
			want:    &metav1.WebServerConfig{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &webServerConfigService{
				ctx: tt.fields.ctx,
				eps: tt.fields.eps,
			}
			got, err := s.GetConfig(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil && tt.want != nil {
				t.Errorf("GetConfig() got = nil, want %v", tt.want)
			}
		})
	}
}

func Test_webServerConfigService_GetConfigTextLines(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEndpoints := epclientv1.NewMockWebServerConfigEndpoints(ctrl)
	ctx := context.Background()

	mockEndpoint := httpclientv1.NewEndpoint[metav1.WebServerOptions, string](func(ctx context.Context, req interface{}) (interface{}, error) {
		return "config text lines", nil
	})
	mockEndpoints.EXPECT().GetConfigTextLines().Return(mockEndpoint)

	type fields struct {
		ctx context.Context
		eps epclientv1.WebServerConfigEndpoints
	}
	type args struct {
		opts *metav1.WebServerOptions
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
				opts: &metav1.WebServerOptions{},
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

	mockEndpoint := httpclientv1.NewEndpoint[metav1.WebServerConfigTargetContextOptions, string](func(ctx context.Context, req interface{}) (interface{}, error) {
		return "context text lines", nil
	})
	mockEndpoints.EXPECT().GetContextTextLines().Return(mockEndpoint)

	type fields struct {
		ctx context.Context
		eps epclientv1.WebServerConfigEndpoints
	}
	type args struct {
		opts *metav1.WebServerConfigTargetContextOptions
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
				opts: &metav1.WebServerConfigTargetContextOptions{},
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

	mockEndpoint := httpclientv1.NewEndpoint[metav1.WebServerConfigTargetContextOptions, []string](func(ctx context.Context, req interface{}) (interface{}, error) {
		return []string{"config1.conf", "config2.conf"}, nil
	})
	mockEndpoints.EXPECT().GetIncludedConfigs().Return(mockEndpoint)

	type fields struct {
		ctx context.Context
		eps epclientv1.WebServerConfigEndpoints
	}
	type args struct {
		opts *metav1.WebServerConfigTargetContextOptions
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
				opts: &metav1.WebServerConfigTargetContextOptions{},
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

	mockEndpoint := httpclientv1.NewEndpoint[httpclientv1.NilBody, []v1.BifrostGroupMeta](func(ctx context.Context, req interface{}) (interface{}, error) {
		return []v1.BifrostGroupMeta{{UintObjectMeta: metav1.UintObjectMeta{Label: "test-group", Value: 1}}}, nil
	})
	mockEndpoints.EXPECT().GetOptions().Return(mockEndpoint)

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
			want:    []v1.BifrostGroupMeta{{UintObjectMeta: metav1.UintObjectMeta{Label: "test-group", Value: 1}}},
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

	mockEndpoint := httpclientv1.NewEndpoint[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], httpclientv1.NilBody](func(ctx context.Context, req interface{}) (interface{}, error) {
		return httpclientv1.NilBody{}, nil
	})
	mockEndpoints.EXPECT().InsertWithClone().Return(mockEndpoint)

	type fields struct {
		ctx context.Context
		eps epclientv1.WebServerConfigEndpoints
	}
	type args struct {
		opts *metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta]
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
				opts: &metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta]{},
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

	mockEndpoint := httpclientv1.NewEndpoint[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], httpclientv1.NilBody](func(ctx context.Context, req interface{}) (interface{}, error) {
		return httpclientv1.NilBody{}, nil
	})
	mockEndpoints.EXPECT().InsertWithNew().Return(mockEndpoint)

	type fields struct {
		ctx context.Context
		eps epclientv1.WebServerConfigEndpoints
	}
	type args struct {
		opts *metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta]
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
				opts: &metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta]{},
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

	mockEndpoint := httpclientv1.NewEndpoint[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], httpclientv1.NilBody](func(ctx context.Context, req interface{}) (interface{}, error) {
		return httpclientv1.NilBody{}, nil
	})
	mockEndpoints.EXPECT().ModifyContextValue().Return(mockEndpoint)

	type fields struct {
		ctx context.Context
		eps epclientv1.WebServerConfigEndpoints
	}
	type args struct {
		opts *metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta]
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
				opts: &metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta]{},
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

	mockEndpoint := httpclientv1.NewEndpoint[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], httpclientv1.NilBody](func(ctx context.Context, req interface{}) (interface{}, error) {
		return httpclientv1.NilBody{}, nil
	})
	mockEndpoints.EXPECT().ModifyWithClone().Return(mockEndpoint)

	type fields struct {
		ctx context.Context
		eps epclientv1.WebServerConfigEndpoints
	}
	type args struct {
		opts *metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta]
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
				opts: &metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta]{},
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

	mockEndpoint := httpclientv1.NewEndpoint[metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta], httpclientv1.NilBody](func(ctx context.Context, req interface{}) (interface{}, error) {
		return httpclientv1.NilBody{}, nil
	})
	mockEndpoints.EXPECT().ModifyWithNew().Return(mockEndpoint)

	type fields struct {
		ctx context.Context
		eps epclientv1.WebServerConfigEndpoints
	}
	type args struct {
		opts *metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta]
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
				opts: &metav1.WebServerConfigContextUpdateOptions[metav1.NewConfigContextMeta]{},
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

	mockEndpoint := httpclientv1.NewEndpoint[metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta], httpclientv1.NilBody](func(ctx context.Context, req interface{}) (interface{}, error) {
		return httpclientv1.NilBody{}, nil
	})
	mockEndpoints.EXPECT().Move().Return(mockEndpoint)

	type fields struct {
		ctx context.Context
		eps epclientv1.WebServerConfigEndpoints
	}
	type args struct {
		opts *metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta]
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
				opts: &metav1.WebServerConfigContextUpdateOptions[metav1.CloneConfigContextMeta]{},
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

	mockEndpoint := httpclientv1.NewEndpoint[metav1.WebServerConfigTargetContextOptions, httpclientv1.NilBody](func(ctx context.Context, req interface{}) (interface{}, error) {
		return httpclientv1.NilBody{}, nil
	})
	mockEndpoints.EXPECT().Remove().Return(mockEndpoint)

	type fields struct {
		ctx context.Context
		eps epclientv1.WebServerConfigEndpoints
	}
	type args struct {
		opts *metav1.WebServerConfigTargetContextOptions
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
				opts: &metav1.WebServerConfigTargetContextOptions{},
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

	mockEndpoint := httpclientv1.NewEndpoint[metav1.WebServerConfigContextPosSearchOptions, []metav1.ConfigContextPos](func(ctx context.Context, req interface{}) (interface{}, error) {
		return []metav1.ConfigContextPos{{Config: "test.conf", ContextPosPath: []int{0, 1}}}, nil
	})
	mockEndpoints.EXPECT().SearchContextPositions().Return(mockEndpoint)

	type fields struct {
		ctx context.Context
		eps epclientv1.WebServerConfigEndpoints
	}
	type args struct {
		opts *metav1.WebServerConfigContextPosSearchOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []metav1.ConfigContextPos
		wantErr bool
	}{
		{
			name: "successful search context positions",
			fields: fields{
				ctx: ctx,
				eps: mockEndpoints,
			},
			args: args{
				opts: &metav1.WebServerConfigContextPosSearchOptions{},
			},
			want:    []metav1.ConfigContextPos{{Config: "test.conf", ContextPosPath: []int{0, 1}}},
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
