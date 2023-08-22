package wechat

import (
	"github.com/eatmoreapple/openwechat"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/wechat"
	wechatReq "github.com/flipped-aurora/gin-vue-admin/server/model/wechat/request"
)

type GroupService struct {
}

// CreateGroup 创建Group记录
// Author [piexlmax](https://github.com/piexlmax)
func (groupService *GroupService) CreateGroup(group *wechat.Group) (err error) {
	err = global.GVA_DB.Create(group).Error
	return err
}

// DeleteGroup 删除Group记录
// Author [piexlmax](https://github.com/piexlmax)
func (groupService *GroupService) DeleteGroup(group wechat.Group) (err error) {
	err = global.GVA_DB.Delete(&group).Error
	return err
}

// DeleteGroupByIds 批量删除Group记录
// Author [piexlmax](https://github.com/piexlmax)
func (groupService *GroupService) DeleteGroupByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]wechat.Group{}, "id in ?", ids.Ids).Error
	return err
}

// UpdateGroup 更新Group记录
// Author [piexlmax](https://github.com/piexlmax)
func (groupService *GroupService) UpdateGroup(group wechat.Group) (err error) {
	err = global.GVA_DB.Save(&group).Error
	return err
}

// GetGroup 根据id获取Group记录
// Author [piexlmax](https://github.com/piexlmax)
func (groupService *GroupService) GetGroup(id uint) (group wechat.Group, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&group).Error
	return
}

// GetGroupInfoList 分页获取Group记录
// Author [piexlmax](https://github.com/piexlmax)
func (groupService *GroupService) GetGroupInfoList(info wechatReq.GroupSearch) (list []wechat.Group, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&wechat.Group{})
	var groups []wechat.Group
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&groups).Error
	return groups, total, err
}
func (GroupService *GroupService) CreateGroups(Groups []*openwechat.Group, selfId string) (err error) {
	// 获取所有的好友
	var groups = make([]wechat.Group, 0)

	//global.GVA_DB.Where("date != ? ",startTime.Format("2006-01-02")).Unscoped().Delete(&wechat.Group{})
	tx := global.GVA_DB.Begin()
	//GroupUser := wechat.GroupUser{}
	//tx.Unscoped().Where("self_id = ?", selfId).Delete(&GroupUser)

	for _, v := range Groups {

		var group2 *wechat.Group
		if result := tx.Where("group_id=?", v.ID()).First(&group2); result.RowsAffected == 0 {
			group := wechat.Group{
				SelfId:   selfId,
				Username: v.UserName,
				GroupId:  v.ID(),
				NickName: v.NickName,
				Is
			}
			groups = append(groups, group)
		} else {
			group2.NickName = v.NickName
			group2.Username = v.UserName
			tx.Save(group2)
		}

	}
	if len(groups) > 0 {
		tx.Create(groups)
	}

	tx.Commit()

	return nil
}
