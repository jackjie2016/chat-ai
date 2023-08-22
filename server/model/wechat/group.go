// 自动生成模板Group
package wechat

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	
)

// Group 结构体
type Group struct {
      global.GVA_MODEL
      Name  string `json:"name" form:"name" gorm:"column:name;comment:群名称;size:150;"`
      NickName  string `json:"nickName" form:"nickName" gorm:"column:nick_name;comment:用户昵称;size:150;"`
      Msg  string `json:"msg" form:"msg" gorm:"column:msg;comment:消息;size:1000;"`
      SelfId  string `json:"selfId" form:"selfId" gorm:"column:self_id;comment:;size:10;"`
      WechatId  string `json:"wechatId" form:"wechatId" gorm:"column:wechat_id;comment:;size:150;"`
      Username  string `json:"username" form:"username" gorm:"column:username;comment:;size:150;"`
      GroupId  string `json:"groupId" form:"groupId" gorm:"column:group_id;comment:;size:150;"`
      NeedStatistics  *bool `json:"needStatistics" form:"needStatistics" gorm:"column:need_statistics;comment:是否需要统计;"`
      IsCollocet  *bool `json:"isCollocet" form:"isCollocet" gorm:"column:is_collocet;comment:是否采集;"`
}


// TableName Group 表名
func (Group) TableName() string {
  return "group"
}

