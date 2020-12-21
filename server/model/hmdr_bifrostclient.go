package model

import (
	"fmt"
	"gin-vue-admin/global"
	"gin-vue-admin/pkg/sort_map"
	"github.com/ClessLi/bifrost/pkg/client/bifrost"
	"go.uber.org/zap"
)

type keyer struct {
	key   uint
	order uint
}

func (k keyer) GetOrder() uint {
	return k.order
}

func (k keyer) Key() interface{} {
	return k.key
}

func NewKeyer(key, order uint) sort_map.Keyer {
	return &keyer{
		key:   key,
		order: order,
	}
}

type BifrostGroups struct {
	dataMap   map[uint]*BifrostGroup
	indexList sort_map.MapIndexes
}

func NewBifrostGroups() sort_map.SortMap {
	return &BifrostGroups{
		dataMap:   make(map[uint]*BifrostGroup),
		indexList: sort_map.NewMapIndexes(),
	}
}

func (bgs BifrostGroups) checkValueType(v interface{}) *BifrostGroup {
	value, vOK := v.(*BifrostGroup)
	if !vOK {
		panic(fmt.Sprintf("value type(%T) is not *BifrostGroup", v))
	}
	return value
}

func (bgs *BifrostGroups) Insert(k, v interface{}) {
	key := checkKeyer(k)
	value := bgs.checkValueType(v)
	newKey := NewKeyer(key.Key().(uint), key.GetOrder())
	bgs.indexList.Insert(newKey)
	bgs.dataMap[newKey.Key().(uint)] = value
}

func (bgs BifrostGroups) Get(k interface{}) (v interface{}) {
	key := checkUINT(k)
	if v, ok := bgs.dataMap[key]; ok {
		return v
	}
	return nil
}

func (bgs *BifrostGroups) Remove(k interface{}) {
	key := checkUINT(k)
	if _, ok := bgs.dataMap[key]; ok {
		delete(bgs.dataMap, key)
	}
}

func (bgs *BifrostGroups) Range(operate func(k, v interface{}) bool) {
	idxRangeFunc := func(idx int, k sort_map.Keyer) bool {
		value, ok := bgs.dataMap[k.Key().(uint)]
		if !ok {
			global.GVA_LOG.Error("BifrostGroups Range error", zap.String("err", fmt.Sprintf("index %d is not dataMap, and it will be removed", k.Key())))
			bgs.Remove(k.Key().(uint))
			return true
		}
		return operate(k.Key(), value)
	}
	bgs.indexList.Range(idxRangeFunc)
	//for _, key := range bgs.indexList {
	//	value, ok := bgs.dataMap[key]
	//	if !ok {
	//		global.GVA_LOG.Error("BifrostGroups Range error", zap.String("err", fmt.Sprintf("index %d is not dataMap, and it will be removed", key)))
	//		bgs.Remove(key)
	//		continue
	//	}
	//	if !operate(key, value) {
	//		return
	//	}
	//}
	//for key, value := range bgs.dataMap {
	//	if !operate(key, value) {
	//		return
	//	}
	//}
}

type BifrostHosts struct {
	dataMap   map[uint]*BifrostHost
	indexList sort_map.MapIndexes
}

func NewBifrostHosts() sort_map.SortMap {
	return &BifrostHosts{
		dataMap:   make(map[uint]*BifrostHost),
		indexList: sort_map.NewMapIndexes(),
	}
}

func (bhs BifrostHosts) checkValueType(v interface{}) *BifrostHost {
	value, vOK := v.(*BifrostHost)
	if !vOK {
		panic(fmt.Sprintf("value type(%T) is not *BifrostHost", v))
	}
	return value
}

func (bhs *BifrostHosts) Insert(k, v interface{}) {
	key := checkKeyer(k)
	value := bhs.checkValueType(v)
	newKey := NewKeyer(key.Key().(uint), key.GetOrder())
	bhs.indexList.Insert(newKey)
	bhs.dataMap[newKey.Key().(uint)] = value
}

func (bhs BifrostHosts) Get(k interface{}) (v interface{}) {
	key := checkUINT(k)
	if v, ok := bhs.dataMap[key]; ok {
		return v
	}
	return nil
}

func (bhs *BifrostHosts) Remove(k interface{}) {
	key := checkUINT(k)
	if _, ok := bhs.dataMap[key]; ok {
		delete(bhs.dataMap, key)
	}
}

func (bhs *BifrostHosts) Range(operate func(k, v interface{}) bool) {
	idxRangeFunc := func(idx int, k sort_map.Keyer) bool {
		value, ok := bhs.dataMap[k.Key().(uint)]
		if !ok {
			global.GVA_LOG.Error("BifrostHosts Range error", zap.String("err", fmt.Sprintf("index %d is not dataMap, and it will be removed", k.Key())))
			bhs.Remove(k.Key().(uint))
			return true
		}
		return operate(k.Key(), value)
	}
	bhs.indexList.Range(idxRangeFunc)
	//for _, key := range bhs.indexList {
	//	value, ok := bhs.dataMap[key]
	//	if !ok {
	//		global.GVA_LOG.Error("BifrostHosts Range error", zap.String("err", fmt.Sprintf("index %d is not dataMap, and it will be removed", key)))
	//		bhs.Remove(key)
	//		continue
	//	}
	//	if !operate(key, value) {
	//		return
	//	}
	//}
	//for key, value := range bhs.dataMap {
	//	if !operate(key, value) {
	//		return
	//	}
	//}
}

type BifrostGroup struct {
	HmdrGroup HmdrGroup
	Hosts     sort_map.SortMap
}

func NewBifrostGroup(group HmdrGroup) *BifrostGroup {
	return &BifrostGroup{
		HmdrGroup: group,
		Hosts:     NewBifrostHosts(),
	}
}

type BifrostHost struct {
	HmdrHost HmdrHost
	Client   *bifrost.Client
}

func checkUINT(k interface{}) uint {
	key, ok := k.(uint)
	if !ok {
		panic(fmt.Sprintf("k type(%T) is not uint", k))
	}
	return key
}

func checkKeyer(k interface{}) sort_map.Keyer {
	key, kOK := k.(sort_map.Keyer)
	if !kOK {
		panic(fmt.Sprintf("k type(%T) is not Keyer", k))
	}
	return key
}

//func uintInsert(slice *[]uint, index int, key uint) {
//	n := len(*slice)
//	*slice = append(*slice, key)
//	if index == -1 {
//		return
//	}
//	for i := n; i < index; i-- {
//		(*slice)[i] = (*slice)[i-1]
//	}
//	(*slice)[index] = key
//}

//func uintBSearchFirstGE(ints []uint, val uint) (int, uint) {
//	return uintBSearchFirstGEInternally(ints, 0, len(ints)-1, val)
//}
//
//func uintBSearchFirstGEInternally(ints []uint, low int, high int, val uint) (int, uint) {
//	if low > high {
//		return -1, 0
//	}
//
//	if ints[low] >= val {
//		return low, ints[low]
//	}
//	mid := low + ((high - low) >> 1)
//	if ints[mid] < val {
//		return uintBSearchFirstGEInternally(ints, mid+1, high, val)
//	} else {
//		return uintBSearchFirstGEInternally(ints, low, mid, val)
//	}
//}

//func uintBSearchFirstGT(ints []uint, val uint) (int, uint) {
//	return uintBSearchFirstGTInternally(ints, 0, len(ints)-1, val)
//}
//
//func uintBSearchFirstGTInternally(ints []uint, low int, high int, val uint) (int, uint) {
//	if low > high {
//		return -1, 0
//	}
//
//	if ints[low] > val {
//		return low, ints[low]
//	}
//	mid := low + ((high - low) >> 1)
//	if ints[mid] > val {
//		return uintBSearchFirstGTInternally(ints, low, mid, val)
//	} else {
//		return uintBSearchFirstGTInternally(ints, mid+1, high, val)
//	}
//}
