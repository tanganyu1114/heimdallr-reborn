package bifrosts

import (
	"fmt"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	"gin-vue-admin/global"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	"gin-vue-admin/pkg/sort_map"
	bifrost "github.com/ClessLi/bifrost/pkg/client/bifrost/v1"
	"github.com/marmotedu/errors"
	"go.uber.org/zap"
	"sync"
)

type Manager interface {
	SyncServersStatus()

	InsertGroup(group v1.Group) error
	InsertHost(host v1.Host) error

	GetGroup(groupid uint) (*Group, error)
	GetBifrostClient(opts metav1.WebServerOptions) (*bifrost.Client, error)

	RemoveAll() error
	RemoveGroupByID(groupid uint) error
	RemoveHost(host v1.Host) error

	UpdateGroup(old, new v1.Group) error
	UpdateHost(old, new v1.Host) error

	Range(func(keyer sort_map.Keyer, v interface{}) bool)

	GetServersStatus() []v1.GroupInfo
}

type manager struct {
	mu *sync.RWMutex
	//bgs         *Groups
	bgs         sort_map.SortMap
	statusCache []v1.GroupInfo
}

var _ Manager = &manager{}

func (m *manager) SyncServersStatus() {

	tmpGroupInfos := make([]v1.GroupInfo, 0)

	syncSrvStatusFromGroups := func(keyer sort_map.Keyer, v interface{}) bool {
		g := v.(*Group)
		groupInfo := v1.GroupInfo{
			Name:  g.Meta.Name,
			Hosts: make([]v1.HostInfo, 0),
		}

		syncSrvStatusFromBifrosts := func(keyer sort_map.Keyer, v interface{}) bool {
			b := v.(*Bifrost)
			var statusList []v1.WebServerStatus
			var status = true
			metrics, err := b.Client.WebServerStatus().Get()
			if err != nil {
				global.GVA_LOG.Error("access the client status failed, hostIp:"+b.Meta.Ipaddr, zap.String("err", err.Error()))
				status = false
				//return false  // 注意：如果报错，散列表后续元素将不被加载
				return true
			}

			for _, info := range metrics.StatusList {
				statusList = append(statusList, v1.WebServerStatus{
					Name:    info.Name,
					Status:  int(info.Status),
					Version: info.Version,
				})
			}

			hostInfo := v1.HostInfo{
				Name:    b.Meta.Name,
				Ipaddr:  b.Meta.Ipaddr,
				Descrip: b.Meta.Description,
				Status:  status,
				AgentInfo: v1.AgentInfo{
					OS:             metrics.OS,
					Time:           metrics.Time,
					Cpu:            metrics.Cpu,
					Mem:            metrics.Mem,
					Disk:           metrics.Disk,
					StatusList:     statusList,
					BifrostVersion: metrics.BifrostVersion,
				},
			}
			groupInfo.Hosts = append(groupInfo.Hosts, hostInfo)
			return true
		}

		g.Bifrosts.Range(syncSrvStatusFromBifrosts)
		tmpGroupInfos = append(tmpGroupInfos, groupInfo)
		return true
	}

	// 遍历并同步获取状态信息
	m.bgs.Range(syncSrvStatusFromGroups)

	// 同步状态信息到缓存
	m.mu.Lock()
	defer m.mu.Unlock()
	m.statusCache = m.statusCache[:0]
	m.statusCache = append(m.statusCache, tmpGroupInfos...)
}

func (m *manager) InsertGroup(group v1.Group) error {
	return m.bgs.Insert(group, NewGroup(group))
}

func (m *manager) InsertHost(host v1.Host) error {
	g, err := m.GetGroup(host.GroupId)
	if err != nil {
		return err
	}

	b, err := NewBifrost(host)
	if err != nil {
		return err
	}

	err = g.Bifrosts.Insert(host, b)
	if err != nil {
		return errors.NewAggregate([]error{err, b.Client.Close()})
	}
	return nil
}

func (m manager) GetGroup(groupid uint) (*Group, error) {
	g, ok := m.bgs.GetByKey(groupid)
	if !ok {
		return nil, errors.Errorf("there is no group %d", groupid)
	}
	return g.(*Group), nil
}

func (m manager) GetBifrostClient(opts metav1.WebServerOptions) (*bifrost.Client, error) {
	g, err := m.GetGroup(opts.GroupID)
	if err != nil {
		return nil, err
	}
	b, ok := g.Bifrosts.GetByKey(opts.HostID)
	if !ok {
		return nil, errors.Errorf("there is no bifrost %d in group %d", opts.HostID, opts.GroupID)
	}
	return b.(*Bifrost).Client, nil
}

func (m *manager) RemoveAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.statusCache = m.statusCache[:0]

	var errs []error
	m.bgs.Range(func(keyer sort_map.Keyer, v interface{}) bool {
		errs = append(errs, m.bgs.RemoveByKey(keyer.Key()))
		return true
	})
	return errors.NewAggregate(errs)
}

func (m *manager) RemoveGroupByID(groupid uint) error {
	return m.bgs.RemoveByKey(groupid)
}

func (m *manager) RemoveHost(host v1.Host) error {
	g, err := m.GetGroup(host.GroupId)
	if err != nil {
		return err
	}

	return g.Bifrosts.RemoveByKey(host.Key())
}

func (m *manager) UpdateHost(old, new v1.Host) error {
	err := m.RemoveHost(old)
	if err != nil {
		return err
	}
	return m.InsertHost(new)
}

func (m *manager) UpdateGroup(old, new v1.Group) error {
	oldG, err := m.GetGroup(old.ID)
	if err != nil {
		return err
	}

	if old.ID == new.ID {
		oldG.Meta = new
		return nil
	}

	err = m.InsertGroup(new)
	if err != nil {
		return err
	}

	newG, err := m.GetGroup(new.ID)
	if err != nil {
		return err
	}

	var errs []error
	// change bifrosts to new group from old group
	oldG.Bifrosts.Range(func(keyer sort_map.Keyer, v interface{}) bool {
		if b, ok := v.(*Bifrost); ok {
			global.GVA_LOG.Info(fmt.Sprintf("bifrost's(id: %d) groupid(%d) change to %d", b.Meta.ID, b.Meta.GroupId, new.ID))
			b.Meta.GroupId = new.ID
			err := newG.Bifrosts.Insert(keyer, b)
			if err != nil {
				global.GVA_LOG.Error("new group insert bifrost failed", zap.String("err", err.Error()))
				errs = append(errs, err)
			}
		}
		return true
	})
	// remove old group index
	m.bgs.(*Groups).indexList.Remove(old)
	delete(m.bgs.(*Groups).dataMap, old.ID)

	return errors.NewAggregate(errs)
}

func (m *manager) Range(f func(keyer sort_map.Keyer, v interface{}) bool) {
	m.bgs.Range(f)
}

func (m *manager) GetServersStatus() []v1.GroupInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.statusCache
}

func New() Manager {
	return &manager{
		mu:          new(sync.RWMutex),
		bgs:         newGroups(),
		statusCache: make([]v1.GroupInfo, 0),
	}
}
