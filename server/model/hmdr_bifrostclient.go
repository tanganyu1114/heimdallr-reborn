package model

import "github.com/ClessLi/bifrost/pkg/client/bifrost"

type BifrostGroup struct {
	HmdrGroup HmdrGroup
	Hosts     map[uint]*BifrostHost
}

type BifrostHost struct {
	HmdrHost HmdrHost
	Client   *bifrost.Client
}
