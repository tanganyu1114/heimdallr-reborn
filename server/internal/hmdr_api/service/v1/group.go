package v1

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
)

type GroupSrv interface {
	Create(ctx context.Context, group v1.Group) error
	Delete(ctx context.Context, groupid uint) error
	DeleteCollections(ctx context.Context, ids metav1.IDsOptions) error
	Get(ctx context.Context, groupid uint) (v1.Group, error)
	List(ctx context.Context, opts metav1.ListOptions) (v1.GroupList, error)
	Update(ctx context.Context, group v1.Group) error
}
