package bifrosts

import (
	"context"
	v1 "gin-vue-admin/api/heimdallr_api/v1"
	"gin-vue-admin/global"
	"gin-vue-admin/internal/pkg/bifrosts"
	metav1 "gin-vue-admin/internal/pkg/meta/v1"
	"github.com/marmotedu/errors"
	"sync"
)

var (
	hOnce           = new(sync.Once)
	singletonHStore *hostStore
)

type hostStore struct {
	bm bifrosts.Manager
}

func (b *hostStore) Create(ctx context.Context, host v1.Host) error {
	err := global.GVA_DB.Create(&host).Error
	if err != nil {
		return err
	}
	return b.bm.InsertHost(host)
}

func (b *hostStore) Delete(ctx context.Context, hostid uint) error {
	host := &v1.Host{
		GVA_MODEL: global.GVA_MODEL{ID: hostid},
	}
	err := global.GVA_DB.Delete(host).Error
	if err != nil {
		return err
	}
	return b.bm.RemoveHost(*host)
}

func (b *hostStore) DeleteCollection(ctx context.Context, ids metav1.IDsOptions) error {
	var hosts []v1.Host
	// find hosts from DB
	err := global.GVA_DB.Find(&hosts, "id in ?", ids.IDs).Error
	if err != nil {
		return err
	}
	// delete hosts from DB
	err = global.GVA_DB.Delete(&[]v1.Host{}, "id in ?", ids.IDs).Error
	if err != nil {
		return err
	}
	// delete hosts from bifrosts groups
	var errs []error
	for _, host := range hosts {
		errs = append(errs, b.bm.RemoveHost(host))
	}
	return errors.NewAggregate(errs)
}

func (b *hostStore) Get(ctx context.Context, hostid uint) (v1.Host, error) {
	host := v1.Host{}
	err := global.GVA_DB.Where("id = ?", hostid).First(&host).Error
	return host, err
}

func (b *hostStore) List(ctx context.Context, opts metav1.ListOptions) (v1.HostList, error) {
	limit := opts.PageSize
	offset := opts.PageSize * (opts.Page - 1)

	// 创建查询会话
	db := global.GVA_DB.Model(&v1.Host{})
	hlist := &v1.HostList{}
	var errs []error
	errs = append(errs, db.Count(&hlist.TotalCount).Error)
	errs = append(errs, db.Limit(limit).Offset(offset).Order("sequence").Find(&hlist.Items).Error)

	return *hlist, errors.NewAggregate(errs)
}

func (b *hostStore) Update(ctx context.Context, host v1.Host) error {
	oldHost := v1.Host{}
	err := global.GVA_DB.Find(&oldHost, host.ID).Error
	if err != nil {
		return err
	}
	err = global.GVA_DB.Save(&host).Error
	if err != nil {
		return err
	}

	return b.bm.UpdateHost(oldHost, host)
}

func newHostStore(store *bifrostsStore) *hostStore {
	hOnce.Do(func() {
		if singletonHStore == nil {
			singletonHStore = &hostStore{
				bm: store.bm,
			}
		}
	})
	if singletonHStore == nil {
		panic(errors.New("singleton host store is nil"))
	}
	return singletonHStore
}
