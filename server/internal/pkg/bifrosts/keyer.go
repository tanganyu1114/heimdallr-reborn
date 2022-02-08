package bifrosts

import "gin-vue-admin/pkg/sort_map"

type uintKeyer struct {
	key   uint
	order uint
}

func (k uintKeyer) GetOrder() uint {
	return k.order
}

func (k uintKeyer) Key() interface{} {
	return k.key
}

func NewUINTKeyer(key, order uint) sort_map.Keyer {
	return &uintKeyer{
		key:   key,
		order: order,
	}
}
