package v1

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
)

type HostStore interface {
	Create(ctx context.Context, host v1.Host) error
	Delete(ctx context.Context, hostid uint) error
	DeleteCollection(ctx context.Context, ids metav1.IDsOptions) error
	Get(ctx context.Context, hostid uint) (v1.Host, error)
	List(ctx context.Context, opts metav1.ListOptions) (v1.HostList, error)
	Update(ctx context.Context, host v1.Host) error
}
