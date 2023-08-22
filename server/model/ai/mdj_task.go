// 自动生成模板MdjTask
package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// MdjTask 结构体
type MdjTask struct {
	global.GVA_MODEL
	Uuid    string `json:"uuid" form:"uuid" gorm:"column:uuid;comment:;"`
	Content string `json:"content" form:"content" gorm:"column:content;comment:;"`
	Prompt  string `json:"prompt" form:"prompt" gorm:"column:prompt;comment:;"`
	Status  bool   `json:"status" form:"status" gorm:"column:status;comment:;default:0;"`
}

// TableName MdjTask 表名
func (MdjTask) TableName() string {
	return "mdj_task"
}
