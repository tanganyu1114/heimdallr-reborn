package bifrosts

import (
	"context"
	v1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	svcv1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/service/v1"
	storev1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/store/v1"
	metav1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/pkg/meta/v1"
)

type hostService struct {
	store storev1.Factory
}

var _ svcv1.HostSrv = (*hostService)(nil)

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

func newHosts(svc *service) svcv1.HostSrv {
	return &hostService{
		store: svc.store,
	}
}
