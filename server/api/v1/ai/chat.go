package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/wechat"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/wechat/handlers"
	"github.com/gin-gonic/gin"
	gogpt "github.com/sashabaranov/go-openai"
	"golang.org/x/net/proxy"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/service"
)

type ChatApi struct{}

var chatModels = []string{gogpt.GPT432K0314, gogpt.GPT4, gogpt.GPT40314, gogpt.GPT432K, gogpt.GPT3Dot5Turbo, gogpt.GPT3Dot5Turbo0301}

// CreateApi
// @Tags      SysApi
// @Summary   创建基础api
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysApi                  true  "api路径, api中文描述, api组, 方法"
// @Success   200   {object}  response.Response{msg=string}  "创建基础api"
// @Router    /chat/createApi [post]
func (s *ChatApi) CreateApi(c *gin.Context) {

	response.OkWithMessage("创建成功", c)
}

// Completion
// @Tags      chat
// @Summary   创建基础api
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      interface{}     true  "gpt上下文, chat分组, 方法"
// @Success   200   {object}  response.Response{msg=string}  "创建基础api"
// @Router    /chat/completion [post]
func (s *ChatApi) Completion(ctx *gin.Context) {
	var request gogpt.ChatCompletionRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)

		return
	}

	if len(request.Messages) == 0 {

		response.FailWithMessage("request messages required", ctx)
		return
	}

	cnf := global.GVA_CONFIG.Gpt
	gptConfig := gogpt.DefaultConfig(cnf.APIKey)

	if cnf.Proxy != "" {
		transport := &http.Transport{}

		if strings.HasPrefix(cnf.Proxy, "socks5h://") {
			// 创建一个 DialContext 对象，并设置代理服务器
			dialContext, err := newDialContext(cnf.Proxy[10:])
			if err != nil {
				panic(err)
			}
			transport.DialContext = dialContext
		} else {
			// 创建一个 HTTP Transport 对象，并设置代理服务器
			proxyUrl, err := url.Parse(cnf.Proxy)
			if err != nil {
				panic(err)
			}
			transport.Proxy = http.ProxyURL(proxyUrl)
		}
		// 创建一个 HTTP 客户端，并将 Transport 对象设置为其 Transport 字段
		gptConfig.HTTPClient = &http.Client{
			Transport: transport,
		}

	}

	// 自定义gptConfig.BaseURL
	if cnf.APIURL != "" {
		gptConfig.BaseURL = cnf.APIURL
	}

	client := gogpt.NewClientWithConfig(gptConfig)
	if request.Messages[0].Role != "system" {
		newMessage := append([]gogpt.ChatCompletionMessage{
			{Role: "system", Content: cnf.BotDesc},
		}, request.Messages...)
		request.Messages = newMessage

	}

	// cnf.Model 是否在 chatModels 中
	if utils.Contains(chatModels, cnf.Model) {
		request.Model = cnf.Model
		resp, err := client.CreateChatCompletion(ctx, request)
		if err != nil {
			response.FailWithMessage(err.Error(), ctx)
			return
		}
		response.OkWithDetailed(gin.H{
			"reply":    resp.Choices[0].Message.Content,
			"messages": append(request.Messages, resp.Choices[0].Message),
		}, "", ctx)

	} else {
		prompt := ""
		for _, item := range request.Messages {
			prompt += item.Content + "/n"
		}
		prompt = strings.Trim(prompt, "/n")

		req := gogpt.CompletionRequest{
			Model:            cnf.Model,
			MaxTokens:        cnf.MaxTokens,
			TopP:             cnf.TopP,
			FrequencyPenalty: cnf.FrequencyPenalty,
			PresencePenalty:  cnf.PresencePenalty,
			Prompt:           prompt,
		}

		resp, err := client.CreateCompletion(ctx, req)
		if err != nil {
			response.FailWithMessage(err.Error(), ctx)
			return
		}
		response.OkWithDetailed(gin.H{
			"reply": resp.Choices[0].Text,
			"messages": append(request.Messages, gogpt.ChatCompletionMessage{
				Role:    "assistant",
				Content: resp.Choices[0].Text,
			}),
		}, "", ctx)

	}
}

type requestJson struct {
	ID      string                        `json:"id"`
	Context []gogpt.ChatCompletionMessage `json:"context"`
}

func (s *ChatApi) Context(ctx *gin.Context) {
	var request requestJson

	err := ctx.BindJSON(&request)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)

		return
	}

	if len(request.Context) == 0 {
		response.FailWithMessage("request messages required", ctx)

		return
	}

	for _, msg := range request.Context {
		jsonStr, err := json.Marshal(msg)
		if err != nil {
			// 处理错误
		}

		err = global.GVA_REDIS.RPush(context.Background(), request.ID+":list", jsonStr).Err()
		if err != nil {
			// 处理错误
		}
	}

	// 设置一个键值对
	marshal, err := json.Marshal(request.Context)
	if err != nil {
		return
	}
	err = global.GVA_REDIS.Set(context.Background(), request.ID, marshal, 0).Err()
	if err != nil {
		panic(err)
	}
	response.OkWithDetailed(request.ID, "", ctx)

	return
}

func (s *ChatApi) Sse(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")

	var request requestForm

	err := ctx.BindQuery(&request)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)

		return
	}
	fmt.Println(request)
	//msgs := make([]Message, 0)

	// 获取一个键值对
	sliceJson, err := global.GVA_REDIS.Get(context.Background(), request.ID).Bytes()
	if err != nil {
		panic(err)
	}

	var Messages []gogpt.ChatCompletionMessage
	err = json.Unmarshal(sliceJson, &Messages)
	if err != nil {
		panic(err)
	}
	//fmt.Println(Messages)

	gptConfig := gogpt.DefaultConfig(global.GVA_CONFIG.Gpt.APIKey)

	if global.GVA_CONFIG.Gpt.Proxy != "" {
		transport := &http.Transport{}

		if strings.HasPrefix(global.GVA_CONFIG.Gpt.Proxy, "socks5h://") {
			// 创建一个 DialContext 对象，并设置代理服务器
			dialContext, err := newDialContext(global.GVA_CONFIG.Gpt.Proxy[10:])
			if err != nil {
				panic(err)
			}
			transport.DialContext = dialContext
		} else {
			// 创建一个 HTTP Transport 对象，并设置代理服务器
			proxyUrl, err := url.Parse(global.GVA_CONFIG.Gpt.Proxy)
			if err != nil {
				panic(err)
			}
			transport.Proxy = http.ProxyURL(proxyUrl)
		}
		// 创建一个 HTTP 客户端，并将 Transport 对象设置为其 Transport 字段
		gptConfig.HTTPClient = &http.Client{
			Transport: transport,
		}

	}

	// 自定义gptConfig.BaseURL
	if global.GVA_CONFIG.Gpt.APIURL != "" {
		gptConfig.BaseURL = global.GVA_CONFIG.Gpt.APIURL
	}

	//Messages = append(Messages, request.Context...)

	gogptRequest := &gogpt.ChatCompletionRequest{
		Model:            "gpt-3.5-turbo",
		Messages:         Messages,
		MaxTokens:        0,
		Temperature:      0.5,
		TopP:             0,
		N:                0,
		Stream:           false,
		Stop:             nil,
		PresencePenalty:  0,
		FrequencyPenalty: 0,
		LogitBias:        nil,
		User:             "",
	}

	client := gogpt.NewClientWithConfig(gptConfig)
	if gogptRequest.Messages[0].Role != "system" {
		newMessage := append([]gogpt.ChatCompletionMessage{
			{Role: "system", Content: global.GVA_CONFIG.Gpt.BotDesc},
		}, gogptRequest.Messages...)
		gogptRequest.Messages = newMessage

	}

	// cnf.Model 是否在 chatModels 中
	if utils.Contains(chatModels, global.GVA_CONFIG.Gpt.Model) {
		gogptRequest.Model = global.GVA_CONFIG.Gpt.Model
		resp, err := client.CreateChatCompletion(ctx, *gogptRequest)
		if err != nil {
			response.FailWithMessage(err.Error(), ctx)

			return
		}
		//c.ResponseJson(ctx, http.StatusOK, "", gin.H{
		//	"reply":    resp.Choices[0].Message.Content,
		//	"messages": append(gogptRequest.Messages, resp.Choices[0].Message),
		//})
		var runeStr = []rune(resp.Choices[0].Message.Content)

		for _, c := range runeStr {
			//fmt.Fprintf(ctx.Writer, "data: %c\n\n", c)
			ctx.Writer.Flush()

			message := Message{
				Text:      string(c),
				Timestamp: time.Now(),
			}

			data, err := json.Marshal(message)
			if err != nil {
				ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			ctx.SSEvent("message", string(data))
			time.Sleep(100 * time.Millisecond)
		}

	} else {
		prompt := ""
		for _, item := range gogptRequest.Messages {
			prompt += item.Content + "/n"
		}
		prompt = strings.Trim(prompt, "/n")

		req := gogpt.CompletionRequest{
			Model:            global.GVA_CONFIG.Gpt.Model,
			MaxTokens:        global.GVA_CONFIG.Gpt.MaxTokens,
			TopP:             global.GVA_CONFIG.Gpt.TopP,
			FrequencyPenalty: global.GVA_CONFIG.Gpt.FrequencyPenalty,
			PresencePenalty:  global.GVA_CONFIG.Gpt.PresencePenalty,
			Prompt:           prompt,
		}

		resp, err := client.CreateCompletion(ctx, req)
		if err != nil {
			response.FailWithMessage(err.Error(), ctx)

			return
		}

		gogptRequest.Messages = append(gogptRequest.Messages, gogpt.ChatCompletionMessage{
			Role:    "assistant",
			Content: resp.Choices[0].Text,
		})
		marshal, err := json.Marshal(gogptRequest.Messages)
		if err != nil {
			return
		}

		err = global.GVA_REDIS.Set(context.Background(), request.ID, marshal, 0).Err()
		if err != nil {
			panic(err)
		}

		message := Message{
			Text:      resp.Choices[0].Text,
			Timestamp: time.Now(),
		}

		data, err := json.Marshal(message)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.SSEvent("message", string(data))

	}
}

type ReplyBody struct {
	Mobile  string `json:"mobile"`
	Content string `json:"content"`
}

func (s *ChatApi) Reply(ctx *gin.Context) {

	var request ReplyBody

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	err = utils.Verify(request, utils.WechatReplyVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	var wUser *wechat.WechatUser

	if global.WechatSelf == nil {
		response.FailWithMessage("先登录", ctx)

		return
	}
	fmt.Println(global.WechatSelf.ID())
	if result := global.GVA_DB.Where("self_id", global.WechatSelf.ID()).First(&wUser, "mobile = ?", request.Mobile); result.RowsAffected != 0 {
		fmt.Println(wUser)
		Uin, _ := strconv.Atoi(wUser.WechatId)
		friend := openwechat.Friend{
			User: &openwechat.User{
				Uin:      int64(Uin),
				NickName: wUser.Nickname,
				UserName: wUser.Username,
			},
		}
		fmt.Println(global.WechatSelf)
		toFriend, err := global.WechatSelf.Self().SendTextToFriend(&friend, request.Content)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(toFriend)

	} else {
		response.FailWithMessage("找不到用户", ctx)

		return
	}

	response.OkWithDetailed("调用成功", "", ctx)

	return
}

var WechatUserservice = service.ServiceGroupApp.WechatServiceGroup.WechatUserService
var WechatGroupservice = service.ServiceGroupApp.WechatServiceGroup.WechatGroupService

func PrintlnQrcodeUrl(uuid string) {
	println("访问下面网址扫描二维码登录")
	qrcodeUrl := openwechat.GetQrcodeUrl(uuid)
	println(qrcodeUrl)
	qrcodeCh <- qrcodeUrl
}

var qrcodeCh = make(chan string, 1)

func (s *ChatApi) Login(ctx *gin.Context) {
	if global.WechatSelf == nil || !global.WechatSelf.Bot().Alive() {

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
		bot.UUIDCallback = PrintlnQrcodeUrl

		//var w http.ResponseWriter
		//fmt.Fprintln(w, openwechat.PrintlnQrcodeUrl)
		go func() {
			// 创建热存储容器对象
			reloadStorage := openwechat.NewFileHotReloadStorage("storage.json")
			//执行热登录
			err := bot.HotLogin(reloadStorage)

			if err != nil {
				if err = bot.Login(); err != nil {
					log.Printf("login error: %v \n", err)
					return
				}
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

			// 获取所有的群组
			groups, err := global.WechatSelf.Groups()
			err = WechatGroupservice.CreateWechatGroups(groups, global.WechatSelf.ID())
			if err != nil {
				return
			}

			// 阻塞主goroutine, 直到发生异常或者用户主动退出
			bot.Block()
		}()
		var qrcode string
		select {
		case qrcode = <-qrcodeCh:
			fmt.Println("qrcode；", qrcode)
		}
		fmt.Println("返回结束")
		response.OkWithDetailed(map[string]string{"qrcode": qrcode}, "请扫码登录", ctx)
		return
	} else {
		fmt.Println(global.WechatSelf.NickName)
		response.OkWithDetailed(map[string]string{
			"self_id":  global.WechatSelf.ID(),
			"username": global.WechatSelf.NickName,
		}, "当前是登录状态", ctx)

		return
	}

}

func (s *ChatApi) WechatWs(ctx *gin.Context) {

}

type requestForm struct {
	ID                  string  `form:"id"`
	Model               string  `form:"model"`
	PresencePenalty     int     `form:"presencePenalty"`
	MaxTokens           int     `form:"maxTokens"`
	Temperature         float32 `form:"temperature"`
	HistoryMessageCount int     `form:"historyMessageCount"`
}
type Message struct {
	Text      string    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}
type dialContextFunc func(ctx context.Context, network, address string) (net.Conn, error)

func newDialContext(socks5 string) (dialContextFunc, error) {
	baseDialer := &net.Dialer{
		Timeout:   60 * time.Second,
		KeepAlive: 60 * time.Second,
	}

	if socks5 != "" {
		// split socks5 proxy string [username:password@]host:port
		var auth *proxy.Auth = nil

		if strings.Contains(socks5, "@") {
			proxyInfo := strings.SplitN(socks5, "@", 2)
			proxyUser := strings.Split(proxyInfo[0], ":")
			if len(proxyUser) == 2 {
				auth = &proxy.Auth{
					User:     proxyUser[0],
					Password: proxyUser[1],
				}
			}
			socks5 = proxyInfo[1]
		}

		dialSocksProxy, err := proxy.SOCKS5("tcp", socks5, auth, baseDialer)
		if err != nil {
			return nil, err
		}

		contextDialer, ok := dialSocksProxy.(proxy.ContextDialer)
		if !ok {
			return nil, err
		}

		return contextDialer.DialContext, nil
	} else {
		return baseDialer.DialContext, nil
	}
}
