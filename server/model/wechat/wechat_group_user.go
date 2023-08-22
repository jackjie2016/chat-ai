// 自动生成模板WechatGroupUser
package wechat

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// WechatGroupUser 结构体
type WechatGroupUser struct {
	global.GVA_MODEL
	SelfId   string `json:"selfId" form:"selfId" gorm:"column:self_id;comment:;size:10; index:daka_tag;"`
	GroupId  string `json:"groupId" form:"groupId" gorm:"column:group_id;comment:;size:150;index:daka_tag;"`
	Nickname string `json:"nickname" form:"nickname" gorm:"column:nickname;comment:;size:150;"`
	Username string `json:"username" form:"username" gorm:"column:username;comment:;size:150;index:daka_tag;"`
	Date     string `json:"date" form:"date" gorm:"column:date;comment:打卡日期;size:150;index:daka_tag;"`
	Sign     bool   `json:"sign" form:"sign" gorm:"column:sign;comment:签到;size:150;"`
}

// TableName WechatGroupUser 表名
func (WechatGroupUser) TableName() string {
	return "wechat_group_user"
}
