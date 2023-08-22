package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/router/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/router/example"
	"github.com/flipped-aurora/gin-vue-admin/server/router/system"
	"github.com/flipped-aurora/gin-vue-admin/server/router/wechat"
)

type RouterGroup struct {
	System  system.RouterGroup
	Example example.RouterGroup
	Ai      ai.RouterGroup
	Wechat  wechat.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
