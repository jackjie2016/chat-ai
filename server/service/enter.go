package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/service/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/service/example"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/wechat"
)

type WechatServiceGroup struct {
	SystemServiceGroup  system.ServiceGroup
	ExampleServiceGroup example.ServiceGroup
	AiServiceGroup      ai.ServiceGroup
	WechatServiceGroup  wechat.ServiceGroup
}

var ServiceGroupApp = new(WechatServiceGroup)
