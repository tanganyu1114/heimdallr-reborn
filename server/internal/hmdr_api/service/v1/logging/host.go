package logging

import (
	"context"

	v1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	svcv1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/service/v1"
	"go.uber.org/zap/zapcore"
)

type hostService struct {
	svc svcv1.Factory
}

var _ svcv1.HostSrv = (*hostService)(nil)

func (h *hostService) Create(ctx context.Context, host v1.Host) (err error) {
	defer func() {
		level := zapcore.DebugLevel
		if err != nil {
			level = zapcore.ErrorLevel
		}
		log(level, "创建主机", host, nil, err)
	}()
	return h.svc.Hosts().Create(ctx, host)
}

func (h *hostService) Delete(ctx context.Context, hostid uint) (err error) {
	defer func() {
		level := zapcore.DebugLevel
		if err != nil {
			level = zapcore.ErrorLevel
		}
		log(level, "删除主机", hostid, nil, err)
	}()
	return h.svc.Hosts().Delete(ctx, hostid)
}

func (h *hostService) DeleteCollection(ctx context.Context, ids v1.IDsOptions) (err error) {
	defer func() {
		level := zapcore.DebugLevel
		if err != nil {
			level = zapcore.ErrorLevel
		}
		log(level, "批量删除主机", ids, nil, err)
	}()
	return h.svc.Hosts().DeleteCollection(ctx, ids)
}

func (h *hostService) Get(ctx context.Context, hostid uint) (host v1.Host, err error) {
	defer func() {
		level := zapcore.DebugLevel
		if err != nil {
			level = zapcore.ErrorLevel
		}
		log(level, "查询主机", hostid, host, err)
	}()
	return h.svc.Hosts().Get(ctx, hostid)
}

func (h *hostService) List(ctx context.Context, opts v1.ListOptions) (hosts v1.HostList, err error) {
	defer func() {
		level := zapcore.DebugLevel
		if err != nil {
			level = zapcore.ErrorLevel
		}
		log(level, "批量查询主机", opts, hosts, err)
	}()
	return h.svc.Hosts().List(ctx, opts)
}

func (h *hostService) Update(ctx context.Context, host v1.Host) (err error) {
	defer func() {
		level := zapcore.DebugLevel
		if err != nil {
			level = zapcore.ErrorLevel
		}
		log(level, "更新主机", host, nil, err)
	}()
	return h.svc.Hosts().Update(ctx, host)
}

func newHosts(svc svcv1.Factory) *hostService {
	return &hostService{svc: svc}
}
