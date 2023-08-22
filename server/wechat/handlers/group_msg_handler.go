package handlers

import (
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/wechat"
	"github.com/flipped-aurora/gin-vue-admin/server/wechat/gtp"
	"log"
	"strings"
	"time"
)

var _ MessageHandlerInterface = (*GroupMessageHandler)(nil)

// GroupMessageHandler 群消息处理
type GroupMessageHandler struct {
}

// handle 处理消息
func (g *GroupMessageHandler) handle(msg *openwechat.Message) error {
	if msg.IsText() {
		return g.ReplyText(msg)
	} else if msg.IsRenameGroup() {
		return g.Rename(msg)
	} else if msg.IsJoinGroup() {
		return g.Join(msg)
	}
	return nil
}

// NewGroupMessageHandler 创建群消息处理器
func NewGroupMessageHandler() MessageHandlerInterface {
	return &GroupMessageHandler{}
}

// ReplyText 发送文本消息到群
func (g *GroupMessageHandler) Rename(msg *openwechat.Message) error {
	sender, _ := msg.Sender()
	group := openwechat.Group{sender}
	log.Printf("Received Group %+v sender Info: %v", group.ID(), sender)

	fmt.Println(group.ID(), group.NickName, group.UserName)
	//		group := wechat.WechatGroup{
	//			SelfId:   global.WechatSelf.ID(),
	//			Username: v.UserName,
	//			GroupId:  v.ID(),
	//			NickName: v.NickName,
	//		}
	return nil
}

// ReplyText 发送文本消息到群
func (g *GroupMessageHandler) Join(msg *openwechat.Message) error {
	sender, _ := msg.Sender()
	group := openwechat.Group{sender}
	log.Printf("Received Group %+v sender Info: %v", group.ID(), sender)

	fmt.Println(group.ID(), group.NickName, group.UserName)
	//		group := wechat.WechatGroup{
	//			SelfId:   global.WechatSelf.ID(),
	//			Username: v.UserName,
	//			GroupId:  v.ID(),
	//			NickName: v.NickName,
	//		}
	return nil
}

// ReplyText 发送文本消息到群
func (g *GroupMessageHandler) ReplyText(msg *openwechat.Message) error {
	// 接收群消息
	fmt.Println("群消息")
	sender, err := msg.Sender()
	group := openwechat.Group{sender}

	// 获取@我的用户
	groupSender, err := msg.SenderInGroup()
	if err != nil {
		log.Printf("get sender in group error :%v \n", err)
		return err
	}
	//采集每天7点之前的打卡数据，每个人只收集一条
	t := time.Now() // 获取当前时间
	var wechatGroupUser wechat.WechatGroupUser

	currentTime := time.Now()
	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), global.GVA_CONFIG.Daka.EndHour, 0, 0, 0, currentTime.Location())
	dakaBegin := startTime.Format("2006-01-02")
	//if err := global.GVA_DB.Unscoped().Where("updated_at < ?", dakaBegin).Where("self_id = ?", global.WechatSelf.ID()).Delete(&wechatGroupUser).Error; err != nil {
	//	log.Printf("delete  :%v \n", err)
	//	return err
	//}

	if t.Hour() < global.GVA_CONFIG.Daka.EndHour {

		if res := global.GVA_DB.Where("date = ?", dakaBegin).Where("group_id = ?", group.ID()).Where("username = ?", groupSender.UserName).Where("self_id = ?", global.WechatSelf.ID()).First(&wechatGroupUser); res.RowsAffected != 0 {

			if !wechatGroupUser.Sign {
				wechatGroupUser.Sign = true
				global.GVA_DB.Save(&wechatGroupUser)
			} else {
				log.Println("已打卡")
			}
		} else {

			//if sender.NickName != sender.Self().NickName {\
			GroupUser := wechat.WechatGroupUser{
				Nickname: groupSender.NickName,
				Username: groupSender.UserName,
				SelfId:   global.WechatSelf.ID(),
				GroupId:  group.ID(),
				Date:     startTime.Format("2006-01-02"),
				Sign:     true,
			}

			global.GVA_DB.Create(&GroupUser)

		}

		//} else {
		//	fmt.Println("更新时间：", wechatGroup.UpdatedAt.Format("2006-01-02 15:04:05"))
		//	currentTime := time.Now()
		//	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
		//	lingdian := startTime.Format("2006-01-02 15:04:05")
		//	if wechatGroup.UpdatedAt.Format("2006-01-02 15:04:05") > lingdian {
		//		fmt.Println("更新时间：", wechatGroup.UpdatedAt.Format("2006-01-02 15:04:05"))
		//	} else {
		//		if err = WechatGroupService.UpdateWechatGroup(wechatGroup); err != nil {
		//			global.GVA_LOG.Error("更新失败!", zap.Error(err))
		//			return err
		//		} else {
		//			fmt.Println("今日更新成功")
		//		}
		//		//fmt.Println("今日未更新")
		//	}
		//	//if wechatGroup.UpdatedAt>
		//}

	}
	//msg.ReplyText("机器人神了，我一会发现了就去修。")

	// 不是@的不处理
	if !msg.IsAt() {
		return nil
	}

	// 替换掉@文本，然后向GPT发起请求
	replaceText := "@ " + sender.Self().NickName
	requestText := strings.TrimSpace(strings.ReplaceAll(msg.Content, replaceText, ""))
	reply := ""

	//gpt的文本查询
	reply, err = gtp.Completions(requestText)
	if err != nil {
		log.Printf("gtp request error: %v \n", err)
		msg.ReplyText("机器人神了，我一会发现了就去修。")
		return err
	}
	if reply == "" {
		return nil
	}
	err = Ai(msg)
	if err != nil {
		return err
	}

	// 回复@我的用户
	reply = strings.TrimSpace(reply)
	reply = strings.Trim(reply, "\n")
	atText := "@" + groupSender.NickName
	replyText := atText + reply

	_, err = msg.ReplyText(replyText)
	if err != nil {
		log.Printf("response group error: %v \n", err)
	}
	return err
}
