package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	aiReq "github.com/flipped-aurora/gin-vue-admin/server/model/ai/request"
)

type MdjTaskService struct {
}

// CreateMdjTask 创建MdjTask记录
// Author [piexlmax](https://github.com/piexlmax)
func (mdjTaskService *MdjTaskService) CreateMdjTask(mdjTask *ai.MdjTask) (err error) {
	err = global.GVA_DB.Create(mdjTask).Error
	return err
}

// DeleteMdjTask 删除MdjTask记录
// Author [piexlmax](https://github.com/piexlmax)
func (mdjTaskService *MdjTaskService) DeleteMdjTask(mdjTask ai.MdjTask) (err error) {
	err = global.GVA_DB.Delete(&mdjTask).Error
	return err
}

// DeleteMdjTaskByIds 批量删除MdjTask记录
// Author [piexlmax](https://github.com/piexlmax)
func (mdjTaskService *MdjTaskService) DeleteMdjTaskByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]ai.MdjTask{}, "id in ?", ids.Ids).Error
	return err
}

// UpdateMdjTask 更新MdjTask记录
// Author [piexlmax](https://github.com/piexlmax)
func (mdjTaskService *MdjTaskService) UpdateMdjTask(mdjTask ai.MdjTask) (err error) {
	err = global.GVA_DB.Save(&mdjTask).Error
	return err
}

// GetMdjTask 根据id获取MdjTask记录
// Author [piexlmax](https://github.com/piexlmax)
func (mdjTaskService *MdjTaskService) GetMdjTask(id uint) (mdjTask ai.MdjTask, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&mdjTask).Error
	return
}

// GetMdjTaskInfoList 分页获取MdjTask记录
// Author [piexlmax](https://github.com/piexlmax)
func (mdjTaskService *MdjTaskService) GetMdjTaskInfoList(info aiReq.MdjTaskSearch) (list []ai.MdjTask, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&ai.MdjTask{})
	var mdjTasks []ai.MdjTask
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&mdjTasks).Error
	return mdjTasks, total, err
}
