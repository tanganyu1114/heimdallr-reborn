package v1

import (
	"context"
	v1 "github.com/tanganyu1114/heimdallr-reborn/api/heimdallr_api/v1"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/internal/pkg/meta/v1"
)

type HostSrv interface {
	Create(ctx context.Context, host v1.Host) error
	Delete(ctx context.Context, hostid uint) error
	DeleteCollection(ctx context.Context, ids metav1.IDsOptions) error
	Get(ctx context.Context, hostid uint) (v1.Host, error)
	List(ctx context.Context, opts metav1.ListOptions) (v1.HostList, error)
	Update(ctx context.Context, host v1.Host) error
}
