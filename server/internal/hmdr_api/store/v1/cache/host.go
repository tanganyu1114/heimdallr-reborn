package cache

import (
	"context"
	v1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	"github.com/tanganyu1114/heimdallr-reborn/server/global"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/pkg/meta/v1"
)

type hostStore struct {
	cacheStore *cacheStore
}

func (h *hostStore) Create(ctx context.Context, host v1.Host) error {
	h.cacheStore.cache.GetGroup(host.GroupId).ReleaseHost(host.ID)
	return h.cacheStore.next.Hosts().Create(ctx, host)
}

func (h *hostStore) Delete(ctx context.Context, hostid uint) error {
	host := &v1.Host{
		GVA_MODEL: global.GVA_MODEL{ID: hostid},
	}
	err := global.GVA_DB.Find(host).Error
	if err != nil {
		return err
	}
	h.cacheStore.cache.GetGroup(host.GroupId).ReleaseHost(host.ID)
	return h.cacheStore.next.Hosts().Delete(ctx, hostid)
}

func (h *hostStore) DeleteCollection(ctx context.Context, ids metav1.IDsOptions) error {
	var hosts []v1.Host
	// find hosts from DB
	err := global.GVA_DB.Find(&hosts, "id in ?", ids.IDs).Error
	if err != nil {
		return err
	}
	for _, host := range hosts {
		h.cacheStore.cache.GetGroup(host.GroupId).ReleaseHost(host.ID)
	}
	return h.cacheStore.next.Hosts().DeleteCollection(ctx, ids)
}

func (h *hostStore) Get(ctx context.Context, hostid uint) (v1.Host, error) {
	return h.cacheStore.next.Hosts().Get(ctx, hostid)
}

func (h *hostStore) List(ctx context.Context, opts metav1.ListOptions) (v1.HostList, error) {
	return h.cacheStore.next.Hosts().List(ctx, opts)
}

func (h *hostStore) Update(ctx context.Context, host v1.Host) error {
	h.cacheStore.cache.GetGroup(host.GroupId).ReleaseHost(host.ID)
	return h.cacheStore.next.Hosts().Update(ctx, host)
}

func newHostStore(cacheStore *cacheStore) *hostStore {
	return &hostStore{cacheStore: cacheStore}
}
