package v1

import (
	"context"
	v1 "github.com/tanganyu1114/heimdallr-reborn/server/api/heimdallr_api/v1"
)

type AgentInfoSrv interface {
	Get(ctx context.Context) ([]v1.GroupInfo, error)
}
