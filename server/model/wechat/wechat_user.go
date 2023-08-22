// 自动生成模板WechatUser
package wechat

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	
)

// WechatUser 结构体
type WechatUser struct {
      global.GVA_MODEL
      SelfId  string `json:"selfId" form:"selfId" gorm:"column:self_id;comment:;size:10;"`
      GroupId  *int `json:"groupId" form:"groupId" gorm:"column:group_id;comment:;"`
      Nickname  string `json:"nickname" form:"nickname" gorm:"column:nickname;comment:;size:150;"`
      WechatId  string `json:"wechatId" form:"wechatId" gorm:"column:wechat_id;comment:;size:150;"`
      Username  string `json:"username" form:"username" gorm:"column:username;comment:;size:150;"`
      Mobile  string `json:"mobile" form:"mobile" gorm:"column:mobile;comment:;size:150;"`
}


// TableName WechatUser 表名
func (WechatUser) TableName() string {
  return "wechat_user"
}

