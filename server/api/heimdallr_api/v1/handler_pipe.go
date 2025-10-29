package v1

import "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/context"

type HandlerPipe[Item any, Set any] interface {
	Handle(target context.Context, in Item) (out Set, err error)

	PosMapHandle(target context.Context, in Item) (next context.PosSet, out Set, err error)
}
