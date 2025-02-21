package cache

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
)

type groupStore struct {
	cacheStore *cacheStore
}

func (g *groupStore) Create(ctx context.Context, group v1.Group) error {
	g.cacheStore.cache.ReleaseGroup(group.ID)
	return g.cacheStore.store.Groups().Create(ctx, group)
}

func (g *groupStore) Delete(ctx context.Context, groupid uint) error {
	g.cacheStore.cache.ReleaseGroup(groupid)
	return g.cacheStore.store.Groups().Delete(ctx, groupid)
}

func (g *groupStore) DeleteCollections(ctx context.Context, ids metav1.IDsOptions) error {
	for _, id := range ids.IDs {
		g.cacheStore.cache.ReleaseGroup(uint(id))
	}
	return g.cacheStore.store.Groups().DeleteCollections(ctx, ids)
}

func (g *groupStore) Get(ctx context.Context, groupid uint) (v1.Group, error) {
	return g.cacheStore.store.Groups().Get(ctx, groupid)
}

func (g *groupStore) List(ctx context.Context, opts metav1.ListOptions) (v1.GroupList, error) {
	return g.cacheStore.store.Groups().List(ctx, opts)
}

func (g *groupStore) Update(ctx context.Context, group v1.Group) error {
	g.cacheStore.cache.ReleaseGroup(group.ID)
	return g.cacheStore.store.Groups().Update(ctx, group)
}

func newGroupStore(cacheStore *cacheStore) *groupStore {
	return &groupStore{cacheStore: cacheStore}
}
