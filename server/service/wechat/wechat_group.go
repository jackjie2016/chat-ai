package wechat

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/wechat"
	wechatReq "github.com/flipped-aurora/gin-vue-admin/server/model/wechat/request"
	"time"
)

type WechatGroupService struct {
}

// CreateWechatGroup 创建WechatGroup记录
// Author [piexlmax](https://github.com/piexlmax)
func (wechatGroupService *WechatGroupService) CreateWechatGroup(wechatGroup *wechat.WechatGroup) (err error) {
	err = global.GVA_DB.Create(&wechatGroup).Error
	return err
}

// DeleteWechatGroup 删除WechatGroup记录
// Author [piexlmax](https://github.com/piexlmax)
func (wechatGroupService *WechatGroupService) DeleteWechatGroup(wechatGroup wechat.WechatGroup) (err error) {
	err = global.GVA_DB.Delete(&wechatGroup).Error
	return err
}

// DeleteWechatGroupByIds 批量删除WechatGroup记录
// Author [piexlmax](https://github.com/piexlmax)
func (wechatGroupService *WechatGroupService) DeleteWechatGroupByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]wechat.WechatGroup{}, "id in ?", ids.Ids).Error
	return err
}

// UpdateWechatGroup 更新WechatGroup记录
// Author [piexlmax](https://github.com/piexlmax)
func (wechatGroupService *WechatGroupService) UpdateWechatGroup(wechatGroup wechat.WechatGroup) (err error) {
	err = global.GVA_DB.Save(&wechatGroup).Error
	return err
}

// GetWechatGroup 根据id获取WechatGroup记录
// Author [piexlmax](https://github.com/piexlmax)
func (wechatGroupService *WechatGroupService) FindWechatGroup(wechatGroup *wechat.WechatGroup) (err error) {

	// 创建db
	db := global.GVA_DB.Model(&wechat.WechatGroup{})
	if wechatGroup.NickName != "" {
		db = db.Where("nick_name = ? ", wechatGroup.NickName)
	}

	//currentTime := time.Now()
	//
	//startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	//startTime.Format("2006/01/02 15:04:05")
	////如果有条件搜索 下方会自动创建搜索语句
	//db = db.Where("created_at elt ? ", startTime)

	err = db.First(&wechatGroup).Error
	return err
}

// GetWechatGroup 根据id获取WechatGroup记录
// Author [piexlmax](https://github.com/piexlmax)
func (wechatGroupService *WechatGroupService) GetWechatGroup(id uint) (wechatGroup wechat.WechatGroup, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&wechatGroup).Error
	return
}

// GetWechatGroupInfoList 分页获取WechatGroup记录
// Author [piexlmax](https://github.com/piexlmax)
func (wechatGroupService *WechatGroupService) GetWechatGroupInfoList(info wechatReq.WechatGroupSearch) (list []wechat.WechatGroup, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&wechat.WechatGroup{})
	var wechatGroups []wechat.WechatGroup
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("updated_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}

	// 如果有条件搜索 下方会自动创建搜索语句
	if info.NickName != "" {
		db = db.Where("nick_name= ?", info.NickName)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&wechatGroups).Error
	return wechatGroups, total, err
}

// GetWechatGroupInfoList 分页获取WechatGroup记录
// Author [piexlmax](https://github.com/piexlmax)
func (wechatGroupService *WechatGroupService) GetStatisticsWechatGroupList() (list []wechat.WechatGroup, total int64, err error) {
	currentTime := time.Now()
	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	// 创建db
	db2 := global.GVA_DB.Model(&wechat.Group{})
	var GroupName []string
	_ = db2.Select("nick_name").Where("need_statistics=?", true).Find(&GroupName).Error

	fmt.Println(GroupName)
	db := global.GVA_DB.Model(&wechat.WechatGroup{})
	var wechatGroups []wechat.WechatGroup
	// 如果有条件搜索 下方会自动创建搜索语句

	//db = db.Where("need_statistics= 1")
	db.Where("date", startTime.Format("2006-01-02"))
	db.Where("nick_name in ?", GroupName)
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Find(&wechatGroups).Error
	return wechatGroups, total, err
}

// GetWechatUserInfoList 分页获取WechatUser记录
// Author [piexlmax](https://github.com/piexlmax)
func (wechatGroupService *WechatGroupService) CreateWechatGroups(Groups []*openwechat.Group, selfId string) (err error) {
	// 获取所有的好友
	currentTime := time.Now()
	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	var groups = make([]wechat.WechatGroup, 0)

	//global.GVA_DB.Where("date != ? ",startTime.Format("2006-01-02")).Unscoped().Delete(&wechat.WechatGroup{})
	tx := global.GVA_DB.Begin()
	//WechatGroupUser := wechat.WechatGroupUser{}
	//tx.Unscoped().Where("self_id = ?", selfId).Delete(&WechatGroupUser)

	for _, v := range Groups {

		var group2 *wechat.WechatGroup
		if result := tx.Where("date", startTime.Format("2006-01-02")).First(&group2, "group_id = ?", v.ID()); result.RowsAffected == 0 {
			group := wechat.WechatGroup{
				SelfId:   selfId,
				Username: v.UserName,
				GroupId:  v.ID(),
				NickName: v.NickName,
				Date:     startTime.Format("2006-01-02"),
			}
			groups = append(groups, group)
		} else {
			//var wechatGroupUser wechat.WechatGroupUser
			//tx.Model(&wechatGroupUser).Where("group_id", group2.GroupId).Update("group_id", v.ID())
			//
			//group2.SelfId = selfId
			//group2.NickName = v.NickName
			//group2.GroupId = v.ID()
			//tx.Save(&group2)
		}
		WechatGroupUsers := make([]wechat.WechatGroupUser, 0)
		members, _ := v.Members()

		for _, vv := range members {
			UpdateWechatGroupUser := wechat.WechatGroupUser{}
			if result := tx.Where("self_id = ?", global.WechatSelf.ID()).Where("date", startTime.Format("2006-01-02")).Where("username = ?", vv.UserName).First(&UpdateWechatGroupUser, "group_id = ?", v.ID()); result.RowsAffected == 0 {
				if res := tx.Where("self_id = ?", global.WechatSelf.ID()).Where("date", startTime.Format("2006-01-02")).Where("nickname = ?", vv.NickName).First(&UpdateWechatGroupUser, "group_id = ?", v.ID()); res.RowsAffected == 0 {
					WechatGroupUser := wechat.WechatGroupUser{
						GVA_MODEL: global.GVA_MODEL{},
						SelfId:    selfId,
						GroupId:   v.ID(),
						Nickname:  vv.NickName,
						Username:  vv.UserName,
						Date:      startTime.Format("2006-01-02"),
					}
					WechatGroupUsers = append(WechatGroupUsers, WechatGroupUser)

				}

			} else {
				if UpdateWechatGroupUser.Nickname != vv.NickName {
					UpdateWechatGroupUser.Nickname = vv.NickName
					tx.Save(UpdateWechatGroupUser)
				}
			}
		}
		if len(WechatGroupUsers) > 0 {
			tx.Create(WechatGroupUsers)

		}

	}
	if len(groups) > 0 {
		tx.Create(groups)
	}

	tx.Commit()

	return nil
}
