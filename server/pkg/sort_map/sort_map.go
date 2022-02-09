package sort_map

type SortMap interface {
	Insert(keyer Keyer, v interface{}) error
	GetByKey(key interface{}) (v interface{}, ok bool)
	RemoveByKey(key interface{}) error
	Range(func(keyer Keyer, v interface{}) bool)
}
