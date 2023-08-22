package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	aiReq "github.com/flipped-aurora/gin-vue-admin/server/model/ai/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type MdjService struct {
}

// CreateMdj 创建Mdj记录
// Author [piexlmax](https://github.com/piexlmax)
func (mdjService *MdjService) CreateMdj(mdj *ai.Mdj) (err error) {
	err = global.GVA_DB.Create(mdj).Error
	return err
}

// DeleteMdj 删除Mdj记录
// Author [piexlmax](https://github.com/piexlmax)
func (mdjService *MdjService) DeleteMdj(mdj ai.Mdj) (err error) {
	err = global.GVA_DB.Delete(&mdj).Error
	return err
}

// DeleteMdjByIds 批量删除Mdj记录
// Author [piexlmax](https://github.com/piexlmax)
func (mdjService *MdjService) DeleteMdjByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]ai.Mdj{}, "id in ?", ids.Ids).Error
	return err
}

// UpdateMdj 更新Mdj记录
// Author [piexlmax](https://github.com/piexlmax)
func (mdjService *MdjService) UpdateMdj(mdj ai.Mdj) (err error) {
	err = global.GVA_DB.Save(&mdj).Error
	return err
}

// GetMdj 根据id获取Mdj记录
// Author [piexlmax](https://github.com/piexlmax)
func (mdjService *MdjService) GetMdj(id uint) (mdj ai.Mdj, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&mdj).Error
	return
}

// GetMdjInfoList 分页获取Mdj记录
// Author [piexlmax](https://github.com/piexlmax)
func (mdjService *MdjService) GetMdjInfoList(info aiReq.MdjSearch) (list []ai.Mdj, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&ai.Mdj{})
	var mdjs []ai.Mdj
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&mdjs).Error
	return mdjs, total, err
}
