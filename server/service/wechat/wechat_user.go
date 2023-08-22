package wechat

import (
	"github.com/eatmoreapple/openwechat"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/wechat"
	wechatReq "github.com/flipped-aurora/gin-vue-admin/server/model/wechat/request"
)

type WechatUserService struct {
}

// CreateWechatUser 创建WechatUser记录
// Author [piexlmax](https://github.com/piexlmax)
func (wechatUserService *WechatUserService) CreateWechatUser(wechatUser *wechat.WechatUser) (err error) {
	err = global.GVA_DB.Create(&wechatUser).Error
	return err
}

// DeleteWechatUser 删除WechatUser记录
// Author [piexlmax](https://github.com/piexlmax)
func (wechatUserService *WechatUserService) DeleteWechatUser(wechatUser wechat.WechatUser) (err error) {
	err = global.GVA_DB.Delete(&wechatUser).Error
	return err
}

// DeleteWechatUserByIds 批量删除WechatUser记录
// Author [piexlmax](https://github.com/piexlmax)
func (wechatUserService *WechatUserService) DeleteWechatUserByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]wechat.WechatUser{}, "id in ?", ids.Ids).Error
	return err
}

// UpdateWechatUser 更新WechatUser记录
// Author [piexlmax](https://github.com/piexlmax)
func (wechatUserService *WechatUserService) UpdateWechatUser(wechatUser wechat.WechatUser) (err error) {
	err = global.GVA_DB.Save(&wechatUser).Error
	return err
}

// GetWechatUser 根据id获取WechatUser记录
// Author [piexlmax](https://github.com/piexlmax)
func (wechatUserService *WechatUserService) GetWechatUser(id uint) (wechatUser wechat.WechatUser, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&wechatUser).Error
	return
}

// GetWechatUserInfoList 分页获取WechatUser记录
// Author [piexlmax](https://github.com/piexlmax)
func (wechatUserService *WechatUserService) GetWechatUserInfoList(info wechatReq.WechatUserSearch) (list []wechat.WechatUser, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&wechat.WechatUser{})
	var wechatUsers []wechat.WechatUser
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&wechatUsers).Error
	return wechatUsers, total, err
}

// GetWechatUserInfoList 分页获取WechatUser记录
// Author [piexlmax](https://github.com/piexlmax)
func (wechatUserService *WechatUserService) CreateWechatUsers(friends []*openwechat.Friend, selfId string) (err error) {
	// 获取所有的好友

	var users = make([]wechat.WechatUser, 0)
	tx := global.GVA_DB.Begin()

	for _, v := range friends {
		var user2 *wechat.WechatUser
		if result := tx.First(&user2, "wechat_id = ?", v.ID()); result.RowsAffected == 0 {
			user := wechat.WechatUser{
				SelfId:   selfId,
				Username: v.UserName,
				WechatId: v.ID(),
				Nickname: v.NickName,
			}
			users = append(users, user)
		} else {
			if v.ID() == user2.WechatId {
				user2.SelfId = selfId
				user2.Nickname = v.NickName
				user2.Username = v.UserName
				tx.Save(&user2)
			}

		}

	}
	if len(users) > 0 {
		tx.Create(users)
	}

	tx.Commit()

	return nil
}
