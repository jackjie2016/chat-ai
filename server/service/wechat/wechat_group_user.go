package wechat

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/wechat"
	wechatReq "github.com/flipped-aurora/gin-vue-admin/server/model/wechat/request"
	"time"
)

type WechatGroupUserService struct {
}

// CreateWechatGroupUser 创建WechatGroupUser记录
// Author [piexlmax](https://github.com/piexlmax)
func (wechatGroupUserService *WechatGroupUserService) CreateWechatGroupUser(wechatGroupUser *wechat.WechatGroupUser) (err error) {
	err = global.GVA_DB.Create(wechatGroupUser).Error
	return err
}

// DeleteWechatGroupUser 删除WechatGroupUser记录
// Author [piexlmax](https://github.com/piexlmax)
func (wechatGroupUserService *WechatGroupUserService) DeleteWechatGroupUser(wechatGroupUser wechat.WechatGroupUser) (err error) {
	err = global.GVA_DB.Delete(&wechatGroupUser).Error
	return err
}

// DeleteWechatGroupUserByIds 批量删除WechatGroupUser记录
// Author [piexlmax](https://github.com/piexlmax)
func (wechatGroupUserService *WechatGroupUserService) DeleteWechatGroupUserByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]wechat.WechatGroupUser{}, "id in ?", ids.Ids).Error
	return err
}

// UpdateWechatGroupUser 更新WechatGroupUser记录
// Author [piexlmax](https://github.com/piexlmax)
func (wechatGroupUserService *WechatGroupUserService) UpdateWechatGroupUser(wechatGroupUser wechat.WechatGroupUser) (err error) {
	err = global.GVA_DB.Save(&wechatGroupUser).Error
	return err
}

// GetWechatGroupUser 根据id获取WechatGroupUser记录
// Author [piexlmax](https://github.com/piexlmax)
func (wechatGroupUserService *WechatGroupUserService) GetWechatGroupUser(id uint) (wechatGroupUser wechat.WechatGroupUser, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&wechatGroupUser).Error
	return
}

// GetWechatGroupUserInfoList 分页获取WechatGroupUser记录
// Author [piexlmax](https://github.com/piexlmax)
func (wechatGroupUserService *WechatGroupUserService) GetWechatGroupUserInfoList(info wechatReq.WechatGroupUserSearch) (list []wechat.WechatGroupUser, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&wechat.WechatGroupUser{})
	var wechatGroupUsers []wechat.WechatGroupUser
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&wechatGroupUsers).Error
	return wechatGroupUsers, total, err
}

// GetWechatGroupSleep 获取群里未打卡的人 flag false 当日打卡的人，true 未打卡的人
// Author [piexlmax](https://github.com/piexlmax)
func (wechatGroupUserService *WechatGroupUserService) GetWechatGroupSleep(groupid string, flag bool) (list []wechat.WechatGroupUser, total int64, err error) {
	var wechatGroupUsers []wechat.WechatGroupUser
	// 创建db
	db := global.GVA_DB.Model(&wechat.WechatGroupUser{})
	//if wechatGroup.NickName != "" {
	//	db = db.Where("nick_name = ? ", wechatGroup.NickName)
	//}

	currentTime := time.Now()

	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	startTime.Format("2006/01/02 15:04:05")
	//如果有条件搜索 下方会自动创建搜索语句

	//now := time.Now() // 获取当前时间

	// 设置时间为 9 点
	//targetTime := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, now.Location())

	// 格式化时间为 "2006/01/02 15:04:05" 格式
	//lastTime := targetTime.Format("2006/01/02 15:04:05")
	//
	////currentTime := time.Now()
	////startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), global.GVA_CONFIG.Daka.EndHour, 0, 0, 0, currentTime.Location())
	////dakaBegin := startTime.Format("2006-01-02 15:04:05")

	db = db.Where("date = ? ", startTime.Format("2006-01-02"))
	if flag {
		db = db.Where("sign = ? ", false)
	} else {
		//db = db.Where("updated_at BETWEEN ? AND ?", now.Format("2006/01/02 15:04:05"), lastTime)

		db = db.Where("sign = ? ", true)
	}
	db = db.Where("group_id = ? ", groupid)
	err = db.Find(&wechatGroupUsers).Error

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	return wechatGroupUsers, total, err
}
