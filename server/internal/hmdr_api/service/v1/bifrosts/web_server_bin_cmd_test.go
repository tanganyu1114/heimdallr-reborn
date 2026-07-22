package bifrosts

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	metav1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	storev1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/store/v1"
	v1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/store/v1"
	"github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/store/v1/fake"
	"go.uber.org/mock/gomock"
)

func Test_webServerBinCMDService_Exec(t *testing.T) {
	webSrvOpts := metav1.WebServerOptions{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := storev1.NewMockFactory(ctrl)
	store.EXPECT().WebServerBinCMD().AnyTimes().Return(new(fake.Store).WebServerBinCMD())
	type fields struct {
		store v1.Factory
	}
	type args struct {
		ctx  context.Context
		opts metav1.WebServerOptions
		arg  []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    metav1.WebServerBinCMDExecResponse
		wantErr bool
	}{
		{
			name:   "request with nil arg",
			fields: fields{store: store},
			args:   args{opts: webSrvOpts},
			want: metav1.WebServerBinCMDExecResponse{
				Successful: true,
				Stdout:     "nginx startup\n",
				Stderr:     "",
			},
			wantErr: false,
		},
		{
			name:   "request with one arg",
			fields: fields{store: store},
			args: args{
				opts: webSrvOpts,
				arg:  []string{"-t"},
			},
			want: metav1.WebServerBinCMDExecResponse{
				Successful: false,
				Stdout:     "",
				Stderr:     "config verify failure\n",
			},
			wantErr: false,
		},
		{
			name:   "request with two args",
			fields: fields{store: store},
			args: args{
				opts: webSrvOpts,
				arg:  []string{"-s", "reload"},
			},
			want: metav1.WebServerBinCMDExecResponse{
				Successful: true,
				Stdout:     "nginx config reload successfully\n",
				Stderr:     "",
			},
			wantErr: false,
		},
		{
			name:   "test unknown args",
			fields: fields{store: store},
			args: args{
				opts: webSrvOpts,
				arg:  []string{"-s", "unknown"},
			},
			want: metav1.WebServerBinCMDExecResponse{
				Successful: false,
				Stdout:     "",
				Stderr:     fmt.Sprintf("unknown argument:\n %v", []string{"-s", "unknown"}),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &webServerBinCMDService{
				store: tt.fields.store,
			}
			got, err := w.Exec(tt.args.ctx, tt.args.opts, tt.args.arg...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Exec() got = %v, want %v", got, tt.want)
			}
		})
	}
}
