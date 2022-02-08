package v1

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	storev1 "gin-vue-admin/internal/hmdr_api/store/v1"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
)

type HostSrv interface {
	Create(ctx context.Context, host v1.Host) error
	Delete(ctx context.Context, hostid uint) error
	DeleteCollection(ctx context.Context, ids metav1.IDsOptions) error
	Get(ctx context.Context, hostid uint) (v1.Host, error)
	List(ctx context.Context, opts metav1.ListOptions) (v1.HostList, error)
	Update(ctx context.Context, host v1.Host) error
}

type hostService struct {
	store storev1.Factory
}

var _ HostSrv = (*hostService)(nil)

func (h *hostService) Create(ctx context.Context, host v1.Host) error {
	return h.store.Hosts().Create(ctx, host)
}

func (h *hostService) Delete(ctx context.Context, hostid uint) error {
	return h.store.Hosts().Delete(ctx, hostid)
}

func (h *hostService) DeleteCollection(ctx context.Context, ids metav1.IDsOptions) error {
	return h.store.Hosts().DeleteCollection(ctx, ids)
}

func (h *hostService) Get(ctx context.Context, hostid uint) (v1.Host, error) {
	return h.store.Hosts().Get(ctx, hostid)
}

func (h *hostService) List(ctx context.Context, opts metav1.ListOptions) (v1.HostList, error) {
	return h.store.Hosts().List(ctx, opts)
}

func (h *hostService) Update(ctx context.Context, host v1.Host) error {
	return h.store.Hosts().Update(ctx, host)
}

func newHosts(svc *service) HostSrv {
	return &hostService{
		store: svc.store,
	}
}
