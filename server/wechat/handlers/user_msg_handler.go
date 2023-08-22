package handlers

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/wechat"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"go.uber.org/zap"
	"log"
	"regexp"
	"strings"
)

var _ MessageHandlerInterface = (*UserMessageHandler)(nil)

// UserMessageHandler 私聊消息处理
type UserMessageHandler struct {
}

// handle 处理消息
func (g *UserMessageHandler) handle(msg *openwechat.Message) error {
	if msg.IsText() {
		return g.ReplyText(msg)
	}
	return nil
}

// NewUserMessageHandler 创建私聊处理器
func NewUserMessageHandler() MessageHandlerInterface {
	return &UserMessageHandler{}
}

// ReplyText 发送文本消息到群
func (g *UserMessageHandler) ReplyText(msg *openwechat.Message) error {
	// 接收私聊消息
	sender, err := msg.Sender()
	log.Printf("Received User %v Text Msg : %v", sender.NickName, msg.Content)
	//var reply string

	if strings.HasPrefix(msg.Content, "/bind ") {
		fmt.Println("The string starts with /gpt ")

		return err
	}

	if strings.Contains(msg.Content, "绑定") {
		// 进入绑定逻辑
		phone := strings.TrimSpace(strings.Replace(msg.Content, "绑定", "", -1))
		fmt.Println(phone)
		if isValidPhone(phone) {
			fmt.Println("手机号合法")
			var wUser *wechat.WechatUser
			if result := global.GVA_DB.First(&wUser, "wechat_id = ?", sender.ID()); result.RowsAffected != 0 {
				//Uin, _ := strconv.Atoi(wUser.WechatId)
				//friend := openwechat.Friend{
				//	User: &openwechat.User{
				//		Uin:      int64(Uin),
				//		NickName: wUser.Nickname,
				//		UserName: wUser.UserName,
				//	},
				//}
			}
			var isbind bool
			var oldmobile string
			if result := global.GVA_DB.First(&wUser, "wechat_id = ?", sender.ID()); result.RowsAffected == 0 {
				user := wechat.WechatUser{
					SelfId:   global.WechatSelf.ID(),
					Username: sender.UserName,
					WechatId: sender.ID(),
					Nickname: sender.NickName,
					Mobile:   phone,
				}
				isbind = true
				global.GVA_DB.Create(&user)
			} else {
				if wUser.Mobile == phone {
					_, err = msg.ReplyText(fmt.Sprintf("修改失败：原绑定就是【%s】", oldmobile))
					if err != nil {
						log.Printf("response user error: %v \n", err)
					}
					return nil
				}
				if wUser.Mobile == "" {
					isbind = true
				} else {
					isbind = false
					oldmobile = wUser.Mobile
				}
				wUser.SelfId = global.WechatSelf.ID()
				wUser.Nickname = sender.NickName
				wUser.Mobile = phone

				global.GVA_DB.Save(&wUser)
			}
			if isbind {
				_, err = msg.ReplyText("绑定成功")
			} else {
				_, err = msg.ReplyText(fmt.Sprintf("修改成功：原绑定【%s】调整为【%s】", oldmobile, phone))
			}
			if err != nil {
				log.Printf("response user error: %v \n", err)
			}
			return err
			// TODO: 执行绑定逻辑
		} else {
			_, err = msg.ReplyText("手机号不合法")
			if err != nil {
				log.Printf("response user error: %v \n", err)
			}
			return err
		}
	}

	if msg.Content == "打卡汇总" {

		daka(msg)
	}
	//err = Ai(msg)
	//if err != nil {
	//	return err
	//}

	//if reply == "" {
	//	return nil
	//}

	// 回复用户
	//reply = strings.TrimSpace(reply)
	//reply = strings.Trim(reply, "\n")
	//_, err = msg.ReplyText(reply)
	//if err != nil {
	//	log.Printf("response user error: %v \n", err)
	//}
	return err
}

var WechatGroupService = service.ServiceGroupApp.WechatServiceGroup.WechatGroupService
var GroupService = service.ServiceGroupApp.WechatServiceGroup.GroupService
var WechatGroupUserService = service.ServiceGroupApp.WechatServiceGroup.WechatGroupUserService

func daka(msg *openwechat.Message) (err error) {

	var wechatGroups []wechat.WechatGroup
	//var total int64
	var reply string

	// 获取所有的群组和群成员
	groups, err := global.WechatSelf.Groups()

	_ = GroupService.CreateGroups(groups, global.WechatSelf.ID())

	err = WechatGroupService.CreateWechatGroups(groups, global.WechatSelf.ID())

	//wechatGroup.Name = "成语接龙测试群"
	wechatGroups, _, err = WechatGroupService.GetStatisticsWechatGroupList()

	for _, DbGroup := range wechatGroups {
		fmt.Println("群id：", DbGroup.GroupId)
		if wechatGroupUsers, total, err := WechatGroupUserService.GetWechatGroupSleep(DbGroup.GroupId, false); err != nil {
			global.GVA_LOG.Error("更新失败!", zap.Error(err))
			return err
		} else {
			//已打卡的
			reply = fmt.Sprintf("今日群《%s》参与打卡%d人：\n", DbGroup.NickName, total)

			for index, user := range wechatGroupUsers {
				if index == 0 {
					reply = fmt.Sprintf("%s%s", reply, user.Nickname)
				} else {
					reply = fmt.Sprintf("%s\n%s", reply, user.Nickname)
				}
			}

			reply = strings.TrimSpace(reply)
			reply = strings.Trim(reply, "\n")

			_, err = msg.ReplyText(reply)
			if err != nil {
				log.Printf("response user error: %v \n", err)
			}
		}

		//未打卡的
		if wechatGroupUsers, total, err := WechatGroupUserService.GetWechatGroupSleep(DbGroup.GroupId, true); err != nil {
			global.GVA_LOG.Error("更新失败!", zap.Error(err))
			return err
		} else {
			reply = fmt.Sprintf("今日群《%s》未打卡%d人：\n", DbGroup.NickName, total)

			for index, user := range wechatGroupUsers {

				if index == 0 {
					reply = fmt.Sprintf("%s%s", reply, user.Nickname)
				} else {
					reply = fmt.Sprintf("%s\n%s", reply, user.Nickname)
				}
			}

			reply = strings.TrimSpace(reply)
			reply = strings.Trim(reply, "\n")

			_, err = msg.ReplyText(reply)
			if err != nil {
				log.Printf("response user error: %v \n", err)
			}
		}
	}

	return nil

}
func isValidPhone(phone string) bool {
	pattern := `^1[3456789]\d{9}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(phone)
}
