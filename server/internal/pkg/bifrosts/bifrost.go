package bifrosts

import (
	"fmt"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	"gin-vue-admin/global"
	"gin-vue-admin/pkg/sort_map"
	bifrost "github.com/ClessLi/bifrost/pkg/client/bifrost/v1"
	"github.com/marmotedu/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"sync"
)

type Bifrost struct {
	Meta   v1.Host
	Client *bifrost.Client
}

func NewBifrost(host v1.Host) (*Bifrost, error) {
	client, err := bifrost.New(host.Ipaddr+":"+host.Port, grpc.WithInsecure()) // 暂时使用非ssl加密协议
	if err != nil {
		return nil, err
	}
	return &Bifrost{
		Meta:   host,
		Client: client,
	}, nil
}

type Bifrosts struct {
	mu        *sync.RWMutex
	dataMap   map[uint]*Bifrost
	indexList sort_map.MapIndexes
}

func NewBifrosts() sort_map.SortMap {
	return &Bifrosts{
		mu:        new(sync.RWMutex),
		dataMap:   make(map[uint]*Bifrost),
		indexList: sort_map.NewMapIndexes(),
	}
}

func (bs Bifrosts) checkValueType(v interface{}) (*Bifrost, error) {
	value, ok := v.(*Bifrost)
	if !ok {
		return nil, errors.Errorf("value type(%T) is not *Bifrost", v)
	}
	return value, nil
}

func (bs *Bifrosts) Insert(keyer sort_map.Keyer, v interface{}) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	ukeyer, err := checkUINTKeyer(keyer)
	if err != nil {
		return errors.Wrap(err, "params error")
	}
	value, err := bs.checkValueType(v)
	if err != nil {
		return errors.Wrap(err, "params error")
	}
	bs.indexList.Insert(ukeyer)
	bs.dataMap[ukeyer.Key().(uint)] = value
	return nil
}

func (bs Bifrosts) Get(key interface{}) (v interface{}, ok bool) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	ukey, err := uintKey(key)
	if err != nil {
		global.GVA_LOG.Warn("call Bifrots.Get() params error", zap.String("err", err.Error()))
		return nil, false
	}
	if v, ok = bs.dataMap[ukey]; ok {
		return v, true
	}
	return nil, false
}

func (bs *Bifrosts) Remove(key interface{}) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	ukey, err := uintKey(key)
	if err != nil {
		return errors.Wrap(err, "params error")
	}
	bs.indexList.Remove(NewUINTKeyer(ukey, 0))
	if b, ok := bs.dataMap[ukey]; ok {
		err = b.Client.Close()
		delete(bs.dataMap, ukey)
	}
	return err
}

func (bs *Bifrosts) Range(operate func(keyer sort_map.Keyer, v interface{}) bool) {
	idxRangeFunc := func(idx int, keyer sort_map.Keyer) bool {
		value, ok := bs.Get(keyer.Key())
		if !ok {
			global.GVA_LOG.Error("Bifrosts Range error", zap.String("err", fmt.Sprintf("faied to get value with index `%v` in dataMap, and it will be removed", keyer.Key())))
			err := bs.Remove(keyer.Key())
			if err != nil {
				global.GVA_LOG.Warn("Bifrosts.Remove() error", zap.String("err", err.Error()))
			}
			return true
		}
		return operate(keyer, value)
	}
	bs.indexList.Range(idxRangeFunc)
}
