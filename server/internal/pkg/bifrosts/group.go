package bifrosts

import (
	"fmt"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	"gin-vue-admin/global"
	"gin-vue-admin/pkg/sort_map"
	"github.com/marmotedu/errors"
	"go.uber.org/zap"
	"sync"
)

type Group struct {
	Meta     v1.Group
	Bifrosts sort_map.SortMap
}

func NewGroup(meta v1.Group) *Group {
	return &Group{
		Meta:     meta,
		Bifrosts: NewBifrosts(),
	}
}

type Groups struct {
	mu        *sync.RWMutex
	dataMap   map[uint]*Group
	indexList sort_map.MapIndexes
}

func newGroups() sort_map.SortMap {
	return &Groups{
		mu:        new(sync.RWMutex),
		dataMap:   make(map[uint]*Group),
		indexList: sort_map.NewMapIndexes(),
	}
}

func (gs Groups) checkValueType(v interface{}) (*Group, error) {
	value, ok := v.(*Group)
	if !ok {
		return nil, errors.Errorf("value type(%T) is not *Group", v)
	}
	return value, nil
}

func (gs *Groups) Insert(keyer sort_map.Keyer, v interface{}) error {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	ukeyer, err := checkUINTKeyer(keyer)
	if err != nil {
		return errors.Wrap(err, "params error")
	}
	value, err := gs.checkValueType(v)
	if err != nil {
		return errors.Wrap(err, "params error")
	}
	gs.indexList.Insert(ukeyer)
	gs.dataMap[ukeyer.Key().(uint)] = value
	return nil
}

func (gs Groups) Get(key interface{}) (v interface{}, ok bool) {
	gs.mu.RLock()
	defer gs.mu.RUnlock()

	ukey, err := uintKey(key)
	if err != nil {
		global.GVA_LOG.Warn("call Groups.Get() params error", zap.String("err", err.Error()))
		return nil, false
	}
	if v, ok := gs.dataMap[ukey]; ok {
		return v, true
	}
	return nil, false
}

func (gs *Groups) Remove(key interface{}) error {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	ukey, err := uintKey(key)
	if err != nil {
		return errors.Wrap(err, "params error")
	}
	gs.indexList.Remove(NewUINTKeyer(ukey, 0))
	var errs []error
	if g, ok := gs.dataMap[ukey]; ok {
		g.Bifrosts.Range(func(keyer sort_map.Keyer, v interface{}) bool {
			errs = append(errs, g.Bifrosts.Remove(keyer))
			return true
		})
		delete(gs.dataMap, ukey)
	}
	return errors.NewAggregate(errs)
}

func (gs *Groups) Range(operate func(keyer sort_map.Keyer, v interface{}) bool) {
	idxRangeFunc := func(idx int, keyer sort_map.Keyer) bool {
		value, ok := gs.Get(keyer.Key())
		if !ok {
			global.GVA_LOG.Error("Groups Range error", zap.String("err", fmt.Sprintf("faied to get value with index `%v` in dataMap, and it will be removed", keyer.Key())))
			err := gs.Remove(keyer.Key())
			if err != nil {
				global.GVA_LOG.Error("Groups.Remove() error", zap.String("err", err.Error()))
			}
			return true
		}
		return operate(keyer, value)
	}
	gs.indexList.Range(idxRangeFunc)
}
