package bifrosts

import (
	"context"
	"fmt"
	bifrostclinetv1 "github.com/ClessLi/bifrost/pkg/client/bifrost/v1"
	"github.com/tanganyu1114/heimdallr-reborn/internal/pkg/bifrosts"
	"github.com/tanganyu1114/heimdallr-reborn/internal/pkg/bifrosts/fake"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/internal/pkg/meta/v1"
	v1 "github.com/tanganyu1114/heimdallr-reborn/internal/pkg/meta/v1"
	"go.uber.org/mock/gomock"
	"reflect"
	"testing"
)

func Test_webServerBinCMDStore_Exec(t *testing.T) {
	webSrvOpts := metav1.WebServerOptions{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBifrostsManager := bifrosts.NewMockManager(ctrl)
	mockBifrostsManager.EXPECT().GetBifrostClient(webSrvOpts).AnyTimes().Return(&bifrostclinetv1.Client{Factory: new(fake.ServiceClient)}, nil)
	type fields struct {
		bm bifrosts.Manager
	}
	type args struct {
		ctx  context.Context
		opts v1.WebServerOptions
		arg  []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    v1.WebServerBinCMDExecResponse
		wantErr bool
	}{
		{
			name:   "request with nil arg",
			fields: fields{bm: mockBifrostsManager},
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
			fields: fields{bm: mockBifrostsManager},
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
			fields: fields{bm: mockBifrostsManager},
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
			fields: fields{bm: mockBifrostsManager},
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
			w := &webServerBinCMDStore{
				bm: tt.fields.bm,
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
