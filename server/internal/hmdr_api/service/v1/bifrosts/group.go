package bifrosts

import (
	"context"

	v1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	svcv1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/service/v1"
	storev1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/store/v1"
)

type groupService struct {
	store storev1.Factory
}

var _ svcv1.GroupSrv = (*groupService)(nil)

func (g *groupService) Create(ctx context.Context, group v1.Group) error {
	return g.store.Groups().Create(ctx, group)
}

func (g *groupService) Delete(ctx context.Context, groupid uint) error {
	return g.store.Groups().Delete(ctx, groupid)
}

func (g *groupService) DeleteCollections(ctx context.Context, ids v1.IDsOptions) error {
	return g.store.Groups().DeleteCollections(ctx, ids)
}

func (g *groupService) Get(ctx context.Context, groupid uint) (v1.Group, error) {
	return g.store.Groups().Get(ctx, groupid)
}

func (g *groupService) List(ctx context.Context, opts v1.ListOptions) (v1.GroupList, error) {
	return g.store.Groups().List(ctx, opts)
}

func (g *groupService) Update(ctx context.Context, group v1.Group) error {
	return g.store.Groups().Update(ctx, group)
}

func newGroups(svc *service) svcv1.GroupSrv {
	return &groupService{
		store: svc.store,
	}
}
