package v1

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
)

type AgentInfoSrv interface {
	Get(ctx context.Context) ([]*v1.GroupInfo, error)
}
