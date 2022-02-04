package v1

import "context"

type AgentInfoSrv interface {
	Get(ctx context.Context)
}
