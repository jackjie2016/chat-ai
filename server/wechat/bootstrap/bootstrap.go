package bootstrap

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	//"github.com/eatmoreapple/openwechat"
	//"github.com/eatmoreapple/openwechat"
	"github.com/eatmoreapple/openwechat"
	"github.com/flipped-aurora/gin-vue-admin/server/wechat/handlers"
)

var WechatUserservice = service.ServiceGroupApp.WechatServiceGroup.WechatUserService
var WechatGroupservice = service.ServiceGroupApp.WechatServiceGroup.WechatGroupService

/*
*
* 只处理热登录
 */
func Run() {

	//bot := openwechat.DefaultBot()
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式，上面登录不上的可以尝试切换这种模式
	// 注册消息处理函数
	bot.MessageHandler = func(msg *openwechat.Message) {
		if msg.IsText() && msg.Content == "ping" {
			msg.ReplyText("pong")
		}
	}
	// 注册消息处理函数
	bot.MessageHandler = handlers.Handler
	// 注册登陆二维码回调
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl
	//var w http.ResponseWriter
	//fmt.Fprintln(w, "Hello, World!")

	// 创建热存储容器对象
	reloadStorage := openwechat.NewFileHotReloadStorage("storage.json")
	//执行热登录
	err := bot.HotLogin(reloadStorage)
	if err != nil {
		fmt.Println("未登录")
		return
		//if err = bot.Login(); err != nil {
		//	log.Printf("login error: %v \n", err)
		//	return
		//}
	}

	// 获取登陆的用户
	global.WechatSelf, err = bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取所有的好友
	friends, err := global.WechatSelf.Friends()
	err = WechatUserservice.CreateWechatUsers(friends, global.WechatSelf.ID())
	if err != nil {
		return
	}

	//var wUser *wechat.WechatUser
	//if result := global.GVA_DB.First(&wUser, "wechat_id = ?", "721830140"); result.RowsAffected != 0 {
	//	Uin, _ := strconv.Atoi(wUser.WechatId)
	//	friend := openwechat.Friend{
	//		User: &openwechat.User{
	//			Uin:      int64(Uin),
	//			NickName: wUser.Nickname,
	//			UserName: wUser.UserName,
	//		},
	//	}
	//
	//	toFriend, err := global.WechatSelf.SendTextToFriend(&friend, "sssss")
	//	if err != nil {
	//		return
	//	}
	//	fmt.Println(toFriend)
	//
	//}

	// 获取所有的群组
	groups, err := global.WechatSelf.Groups()
	err = WechatGroupservice.CreateWechatGroups(groups, global.WechatSelf.ID())
	if err != nil {
		return
	}

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}
