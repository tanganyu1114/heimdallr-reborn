package cache

import (
	"context"
	"github.com/marmotedu/errors"
	v1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	"github.com/tanganyu1114/heimdallr-reborn/server/global"
	storev1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/store/v1"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/pkg/meta/v1"
	"go.uber.org/mock/gomock"
	"reflect"
	"testing"
	"time"
)

func Test_hostStore_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	nextStore := storev1.NewMockFactory(ctrl)
	nextHostStore := storev1.NewMockHostStore(ctrl)
	nextStore.EXPECT().Hosts().AnyTimes().Return(nextHostStore)
	h1 := v1.Host{
		GVA_MODEL:   global.GVA_MODEL{ID: 1},
		GroupId:     1,
		Name:        "10.1.1.1",
		Description: "test 1 host",
		Status:      true,
		Ipaddr:      "10.1.1.1",
		Port:        "12321",
		Token:       "",
		Sequence:    100,
	}
	h2 := v1.Host{
		GVA_MODEL:   global.GVA_MODEL{ID: 2},
		GroupId:     1,
		Name:        "10.1.1.2",
		Description: "test 2 host",
		Status:      true,
		Ipaddr:      "10.1.1.2",
		Port:        "12321",
		Token:       "",
		Sequence:    100,
	}
	he := v1.Host{
		GVA_MODEL:   global.GVA_MODEL{ID: 3},
		GroupId:     2,
		Name:        "localhost",
		Description: "invalid host",
		Status:      false,
		Ipaddr:      "localhost",
		Port:        "",
		Token:       "",
		Sequence:    0,
	}
	nextHostStore.EXPECT().Create(nil, h1).AnyTimes().Return(nil)
	nextHostStore.EXPECT().Create(nil, h2).AnyTimes().Return(nil)
	nextHostStore.EXPECT().Create(nil, he).AnyTimes().Return(errors.New("invalid host"))
	cs := GetCacheStore(nextStore, time.Second)
	err := cs.Hosts().Create(nil, h1)
	if err != nil {
		t.Fatal(err)
	}
	type fields struct {
		cacheStore *cacheStore
	}
	type args struct {
		ctx  context.Context
		host v1.Host
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "normal test",
			fields: fields{cacheStore: cs.(*cacheStore)},
			args:   args{host: h2},
		},
		{
			name:    "error test",
			fields:  fields{cacheStore: cs.(*cacheStore)},
			args:    args{host: he},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &hostStore{
				cacheStore: tt.fields.cacheStore,
			}
			if err := h.Create(tt.args.ctx, tt.args.host); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_hostStore_Delete(t *testing.T) {
	type fields struct {
		cacheStore *cacheStore
	}
	type args struct {
		ctx    context.Context
		hostid uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &hostStore{
				cacheStore: tt.fields.cacheStore,
			}
			if err := h.Delete(tt.args.ctx, tt.args.hostid); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_hostStore_DeleteCollection(t *testing.T) {
	type fields struct {
		cacheStore *cacheStore
	}
	type args struct {
		ctx context.Context
		ids metav1.IDsOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &hostStore{
				cacheStore: tt.fields.cacheStore,
			}
			if err := h.DeleteCollection(tt.args.ctx, tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("DeleteCollection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_hostStore_Get(t *testing.T) {
	type fields struct {
		cacheStore *cacheStore
	}
	type args struct {
		ctx    context.Context
		hostid uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    v1.Host
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &hostStore{
				cacheStore: tt.fields.cacheStore,
			}
			got, err := h.Get(tt.args.ctx, tt.args.hostid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hostStore_List(t *testing.T) {
	type fields struct {
		cacheStore *cacheStore
	}
	type args struct {
		ctx  context.Context
		opts metav1.ListOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    v1.HostList
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &hostStore{
				cacheStore: tt.fields.cacheStore,
			}
			got, err := h.List(tt.args.ctx, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hostStore_Update(t *testing.T) {
	type fields struct {
		cacheStore *cacheStore
	}
	type args struct {
		ctx  context.Context
		host v1.Host
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &hostStore{
				cacheStore: tt.fields.cacheStore,
			}
			if err := h.Update(tt.args.ctx, tt.args.host); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_newHostStore(t *testing.T) {
	type args struct {
		cacheStore *cacheStore
	}
	tests := []struct {
		name string
		args args
		want *hostStore
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newHostStore(tt.args.cacheStore); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newHostStore() = %v, want %v", got, tt.want)
			}
		})
	}
}
