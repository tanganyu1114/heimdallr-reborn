package sort_map

//go:generate mockgen -self_package=gin-vue-admin/pkg/sort_map -destination=mock_sort_map.go -package=sort_map gin-vue-admin/pkg/sort_map SortMap,MapIndexes,Keyer
type SortMap interface {
	Insert(keyer Keyer, v interface{}) error
	GetByKey(key interface{}) (v interface{}, ok bool)
	RemoveByKey(key interface{}) error
	Range(func(keyer Keyer, v interface{}) bool)
}
