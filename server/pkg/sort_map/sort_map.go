package sort_map

type SortMap interface {
	Insert(k, v interface{})
	Get(k interface{}) (v interface{})
	Remove(k interface{})
	Range(func(k, v interface{}) bool)
}
