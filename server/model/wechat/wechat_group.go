// 自动生成模板WechatGroup
package wechat

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// WechatGroup 结构体
type WechatGroup struct {
	global.GVA_MODEL
	Name           string `json:"name" form:"name" gorm:"column:name;comment:群名称;size:150;"`
	NickName       string `json:"nickName" form:"nickName" gorm:"column:nick_name;comment:用户昵称;size:150;"`
	Msg            string `json:"msg" form:"msg" gorm:"column:msg;comment:消息;size:1000;"`
	SelfId         string `json:"selfId" form:"selfId" gorm:"column:self_id;comment:;size:10;"`
	WechatId       string `json:"wechatId" form:"wechatId" gorm:"column:wechat_id;comment:;size:150;"`
	Nickname       string `json:"nickname" form:"nickname" gorm:"column:nickname;comment:;size:150;"`
	Username       string `json:"username" form:"username" gorm:"column:username;comment:;size:150;"`
	GroupId        string `json:"groupId" form:"groupId" gorm:"column:group_id;comment:;size:150;"`
	NeedStatistics bool   `json:"need_statistics" form:"need_statistics" gorm:"column:need_statistics;comment:是否需要统计;size:1;"`
	Date           string `json:"date" form:"date" gorm:"column:date;comment:打卡日期;size:50;"`
}

// TableName WechatGroup 表名
func (WechatGroup) TableName() string {
	return "wechat_group"
}
