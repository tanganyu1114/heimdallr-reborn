package logging

import (
	"context"

	v1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
	svcv1 "github.com/tanganyu1114/heimdallr-reborn/server/internal/hmdr_api/service/v1"
	"go.uber.org/zap/zapcore"
)

type groupService struct {
	svc svcv1.Factory
}

var _ svcv1.GroupSrv = (*groupService)(nil)

func (g *groupService) Create(ctx context.Context, group v1.Group) (err error) {
	defer func() {
		level := zapcore.DebugLevel
		if err != nil {
			level = zapcore.ErrorLevel
		}
		log(level, "创建组", group, nil, err)
	}()
	return g.svc.Groups().Create(ctx, group)
}

func (g *groupService) Delete(ctx context.Context, groupid uint) (err error) {
	defer func() {
		level := zapcore.DebugLevel
		if err != nil {
			level = zapcore.ErrorLevel
		}
		log(level, "删除组", groupid, nil, err)
	}()
	return g.svc.Groups().Delete(ctx, groupid)
}

func (g *groupService) DeleteCollections(ctx context.Context, ids v1.IDsOptions) (err error) {
	defer func() {
		level := zapcore.DebugLevel
		if err != nil {
			level = zapcore.ErrorLevel
		}
		log(level, "批量删除组", ids, nil, err)
	}()
	return g.svc.Groups().DeleteCollections(ctx, ids)
}

func (g *groupService) Get(ctx context.Context, groupid uint) (group v1.Group, err error) {
	defer func() {
		level := zapcore.DebugLevel
		if err != nil {
			level = zapcore.ErrorLevel
		}
		log(level, "查询组", groupid, group, err)
	}()
	return g.svc.Groups().Get(ctx, groupid)
}

func (g *groupService) List(ctx context.Context, opts v1.ListOptions) (groups v1.GroupList, err error) {
	defer func() {
		level := zapcore.DebugLevel
		if err != nil {
			level = zapcore.ErrorLevel
		}
		log(level, "批量查询组", opts, groups, err)
	}()
	return g.svc.Groups().List(ctx, opts)
}

func (g *groupService) Update(ctx context.Context, group v1.Group) (err error) {
	defer func() {
		level := zapcore.DebugLevel
		if err != nil {
			level = zapcore.ErrorLevel
		}
		log(level, "更新组", group, nil, err)
	}()
	return g.svc.Groups().Update(ctx, group)
}

func newGroups(svc svcv1.Factory) *groupService {
	return &groupService{svc: svc}
}
