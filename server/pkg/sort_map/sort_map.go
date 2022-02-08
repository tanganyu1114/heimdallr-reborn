package sort_map

type SortMap interface {
	Insert(keyer Keyer, v interface{}) error
	Get(key interface{}) (v interface{}, ok bool)
	Remove(key interface{}) error
	Range(func(keyer Keyer, v interface{}) bool)
}
