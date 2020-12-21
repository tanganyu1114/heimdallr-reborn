package sort_map

type Keyer interface {
	GetOrder() uint
	Key() interface{}
}
