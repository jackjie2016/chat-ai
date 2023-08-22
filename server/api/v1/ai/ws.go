package ai

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"sync"
)

type WsApi struct{}

var (
	upgrader = websocket.Upgrader{}
	counter  int
	mutex    sync.Mutex
)

func generateID() string {
	mutex.Lock()
	counter++
	id := counter
	mutex.Unlock()
	return strconv.Itoa(id)
}

func handleWebSocket(conn *websocket.Conn, userID string) {
	// 处理WebSocket连接逻辑
	// 在这里可以使用userID标识特定的用户连接
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read message from WebSocket:", err)
			break
		}

		log.Println("Received message:", string(msg))

		err = conn.WriteMessage(websocket.TextMessage, []byte("Received your message"))
		if err != nil {
			log.Println("Failed to send message to WebSocket:", err)
			break
		}
	}
}
func (a WsApi) Talk(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}

	// 分配唯一的用户ID
	userID := generateID()

	// 处理WebSocket连接
	handleWebSocket(conn, userID)

	response.OkWithDetailed("", "获取成功", ctx)
}
