package model

import (
	"fmt"
	"gin-vue-admin/global"
	"gin-vue-admin/pkg/sort_map"
	"github.com/ClessLi/bifrost/pkg/client/bifrost"
	"go.uber.org/zap"
)

type BifrostGroups struct {
	dataMap   map[uint]*BifrostGroup
	indexList []uint
}

func NewBifrostGroups() sort_map.SortMap {
	return &BifrostGroups{
		dataMap:   make(map[uint]*BifrostGroup),
		indexList: make([]uint, 0),
	}
}

func (bgs BifrostGroups) checkKeyType(k interface{}) uint {
	key, kOK := k.(uint)
	if !kOK {
		panic(fmt.Sprintf("key type(%T) is not uint", k))
	}
	return key
}

func (bgs BifrostGroups) checkValueType(v interface{}) *BifrostGroup {
	value, vOK := v.(*BifrostGroup)
	if !vOK {
		panic(fmt.Sprintf("value type(%T) is not *BifrostGroup", v))
	}
	return value
}

func (bgs *BifrostGroups) Insert(k, v interface{}) {
	key := bgs.checkKeyType(k)
	value := bgs.checkValueType(v)
	//index, idxKey := uintBSearchFirstGE(bgs.indexList, key)
	//if idxKey == key {
	//	bgs.indexList[index] = key
	//} else if idxKey < key {
	//	panic("uintBSearchFirstGE functions result error")
	//} else {
	//	uintInsert(&bgs.indexList, index, key)
	//}
	//uintInsert(&bgs.indexList, index, key)
	bgs.indexList = append(bgs.indexList, key)
	bgs.dataMap[key] = value
}

func (bgs BifrostGroups) Get(k interface{}) (v interface{}) {
	key := bgs.checkKeyType(k)
	if v, ok := bgs.dataMap[key]; ok {
		return v
	}
	return nil
}

func (bgs *BifrostGroups) Remove(k interface{}) {
	key := bgs.checkKeyType(k)
	if _, ok := bgs.dataMap[key]; ok {
		delete(bgs.dataMap, key)
	}
}

func (bgs *BifrostGroups) Range(operate func(k, v interface{}) bool) {
	for _, key := range bgs.indexList {
		value, ok := bgs.dataMap[key]
		if !ok {
			global.GVA_LOG.Error("BifrostGroups Range error", zap.String("err", fmt.Sprintf("index %d is not dataMap, and it will be removed", key)))
			bgs.Remove(key)
			continue
		}
		if !operate(key, value) {
			return
		}
	}
	//for key, value := range bgs.dataMap {
	//	if !operate(key, value) {
	//		return
	//	}
	//}
}

type BifrostHosts struct {
	dataMap   map[uint]*BifrostHost
	indexList []uint
}

func NewBifrostHosts() sort_map.SortMap {
	return &BifrostHosts{
		dataMap:   make(map[uint]*BifrostHost),
		indexList: make([]uint, 0),
	}
}

func (bhs BifrostHosts) checkKeyType(k interface{}) uint {
	key, kOK := k.(uint)
	if !kOK {
		panic(fmt.Sprintf("key type(%T) is not uint", k))
	}
	return key
}

func (bhs BifrostHosts) checkValueType(v interface{}) *BifrostHost {
	value, vOK := v.(*BifrostHost)
	if !vOK {
		panic(fmt.Sprintf("value type(%T) is not *BifrostHost", v))
	}
	return value
}

func (bhs *BifrostHosts) Insert(k, v interface{}) {
	key := bhs.checkKeyType(k)
	value := bhs.checkValueType(v)
	//index, idxKey := uintBSearchFirstGE(bhs.indexList, key)
	//if idxKey == key {
	//	bhs.indexList[index] = key
	//} else if idxKey < key {
	//	panic("uintBSearchFirstGE functions result error")
	//} else {
	//	uintInsert(&bhs.indexList, index, key)
	//}
	//uintInsert(&bhs.indexList, index, key)
	bhs.indexList = append(bhs.indexList, key)
	bhs.dataMap[key] = value
}

func (bhs BifrostHosts) Get(k interface{}) (v interface{}) {
	key := bhs.checkKeyType(k)
	if v, ok := bhs.dataMap[key]; ok {
		return v
	}
	return nil
}

func (bhs *BifrostHosts) Remove(k interface{}) {
	key := bhs.checkKeyType(k)
	if _, ok := bhs.dataMap[key]; ok {
		delete(bhs.dataMap, key)
	}
}

func (bhs *BifrostHosts) Range(operate func(k, v interface{}) bool) {
	for _, key := range bhs.indexList {
		value, ok := bhs.dataMap[key]
		if !ok {
			global.GVA_LOG.Error("BifrostHosts Range error", zap.String("err", fmt.Sprintf("index %d is not dataMap, and it will be removed", key)))
			bhs.Remove(key)
			continue
		}
		if !operate(key, value) {
			return
		}
	}
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
