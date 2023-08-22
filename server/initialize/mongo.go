package initialize

import (
	"fmt"
	"gitee.com/phper95/pkg/nosql"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

func InitMongoClient() {
	mongoCfg := &global.GVA_CONFIG.MongoDB
	fmt.Println(mongoCfg)
	err := nosql.InitMongoClient(nosql.DefaultMongoClient, mongoCfg.User,
		mongoCfg.Password, mongoCfg.Hosts, 200)
	if err != nil {

		panic(err)
	}
	global.GVA_MONGO = nosql.GetMongoClient(nosql.DefaultMongoClient)
}
