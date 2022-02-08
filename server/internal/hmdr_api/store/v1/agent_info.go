package v1

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
)

type AgentInfoStore interface {
	Get(ctx context.Context) ([]v1.GroupInfo, error)
	SyncAgentInfos()
}
