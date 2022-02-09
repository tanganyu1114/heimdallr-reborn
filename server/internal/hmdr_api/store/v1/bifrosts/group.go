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
	gOnce           = new(sync.Once)
	singletonGStore *groupStore
)

type groupStore struct {
	bm bifrosts.Manager
}

func (g *groupStore) Create(ctx context.Context, group v1.Group) error {
	err := global.GVA_DB.Create(&group).Error
	if err != nil {
		return err
	}
	return g.bm.InsertGroup(group)
}

func (g *groupStore) Delete(ctx context.Context, groupid uint) error {
	err := global.GVA_DB.Delete(&[]v1.Host{}, "group_id = ?", groupid).Error
	if err != nil {
		return err
	}

	group := &v1.Group{
		GVA_MODEL: global.GVA_MODEL{ID: groupid},
	}
	err = global.GVA_DB.Delete(group).Error
	if err != nil {
		return err
	}

	return g.bm.RemoveGroupByID(groupid)
}

func (g *groupStore) DeleteCollections(ctx context.Context, ids metav1.IDsOptions) error {
	err := global.GVA_DB.Delete(&[]v1.Group{}, "id in ?", ids.IDs).Error
	if err != nil {
		return err
	}
	var errs []error
	for _, id := range ids.IDs {
		errs = append(errs, g.bm.RemoveGroupByID(uint(id)))
	}
	return errors.NewAggregate(errs)
}

func (g *groupStore) Get(ctx context.Context, groupid uint) (v1.Group, error) {
	group := v1.Group{}
	err := global.GVA_DB.Where("id = ?", groupid).First(&group).Error
	return group, err
}

func (g *groupStore) List(ctx context.Context, opts metav1.ListOptions) (v1.GroupList, error) {
	limit := opts.PageSize
	offset := opts.PageSize * (opts.Page - 1)

	// 创建查询会话
	db := global.GVA_DB.Model(&v1.Group{})
	glist := &v1.GroupList{}
	var errs []error
	errs = append(errs, db.Count(&glist.TotalCount).Error)
	errs = append(errs, db.Limit(limit).Offset(offset).Order("sequence").Find(&glist.Items).Error)

	return *glist, errors.NewAggregate(errs)
}

func (g *groupStore) Update(ctx context.Context, group v1.Group) error {
	oldgroup := v1.Group{}
	err := global.GVA_DB.Find(&oldgroup, group.ID).Error
	if err != nil {
		return err
	}
	err = global.GVA_DB.Save(&group).Error
	if err != nil {
		return err
	}

	return g.bm.UpdateGroup(oldgroup, group)
}

func newGroupStore(store *bifrostsStore) *groupStore {
	gOnce.Do(func() {
		if singletonGStore == nil {
			singletonGStore = &groupStore{
				bm: store.bm,
			}
		}
	})
	if singletonGStore == nil {
		panic(errors.New("singleton group store is nil"))
	}
	return singletonGStore
}
