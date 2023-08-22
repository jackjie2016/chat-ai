package ai

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	services "github.com/flipped-aurora/gin-vue-admin/server/logic"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	utils "github.com/flipped-aurora/gin-vue-admin/server/utils/discord"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"net/http"
	"sort"
	"time"
)

type MidBotApi struct{}
type RequestTrigger struct {
	Type         string `json:"type"`
	DiscordMsgId string `json:"discordMsgId,omitempty"`
	MsgHash      string `json:"msgHash,omitempty"`
	Prompt       string `json:"prompt,omitempty"`
	Index        int64  `json:"index,omitempty"`
}
type CustomMsg struct {
	ai.ChatMsg
	Mdj_Msg *ai.MdjMsg2 `json:"mdj_Msg"`
}

func (a MidBotApi) MidjourneyList(c *gin.Context) {

	//ShouldBindQuery和ShouldBindJSON要搞清楚
	var reqInfo request.MdjSearch

	err := c.ShouldBindJSON(&reqInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	opts := &options.FindOptions{}
	filter := bson.D{}

	total, err := global.GVA_MONGO.QueryCount("mdj", "mdj_chat", filter, 0)
	fmt.Printf("%+v\n", filter)
	//liveInfo.settlementAmountThirtyRecent

	Chatmsgs := []ai.ChatMsg{}

	customMsgs := make([]*CustomMsg, 0)

	_, err = global.GVA_MONGO.FindWithOpts("mdj", "mdj_chat", int64(reqInfo.Page-1), int64(reqInfo.PageSize), filter, opts.SetSort(bson.M{"create_time": -1}), &Chatmsgs)

	for _, chatmsg := range Chatmsgs {
		msg := &ai.MdjMsg2{}
		if chatmsg.IsSender == false {
			//, {"userId", chatmsg.UserId}
			err = global.GVA_MONGO.FindOne("mdj", "mdj_msg2", bson.D{{"hashid", chatmsg.HashId}}, msg)
		}
		customMsg := &CustomMsg{
			ChatMsg: chatmsg,
			Mdj_Msg: msg,
		}
		customMsgs = append(customMsgs, customMsg)
	}

	sort.Slice(customMsgs, func(i, j int) bool {
		return customMsgs[i].CreateTime < customMsgs[j].CreateTime
	})

	response.OkWithDetailed(response.PageResult{
		List:     customMsgs,
		Total:    total,
		Page:     int(reqInfo.Page),
		PageSize: int(reqInfo.PageSize),
	}, "获取成功", c)
}

func (a MidBotApi) MidjourneyBot(c *gin.Context) {
	var body RequestTrigger
	if err := c.ShouldBindJSON(&body); err != nil {

		response.FailWithMessage(err.Error(), c)
		return
	}

	var err error
	uid := 10
	uuidStr := uuid.New().String()

	Msg := ai.ChatMsg{
		HashId:     uuidStr,
		UserId:     uid,
		IsSender:   true,
		Content:    body.Prompt,
		Type:       body.Type,
		CreateTime: time.Now().Format("2006-01-02 15:04:02"),
		UpdateTime: time.Now().Format("2006-01-02 15:04:02"),
	}
	switch body.Type {
	case "generate":
		if len(body.Prompt) == 0 {
			err = errors.New("invalid Prompt")

			break
		}
		Msg.Content = body.Prompt
		body.Prompt = fmt.Sprintf("[%s_%d_%s] %s ~~", global.GVA_CONFIG.Discord.MSG_PRFIX, uid, uuidStr, body.Prompt)

		err = utils.GenerateImage(body.Prompt)

	case "upscale":
		err = utils.ImageUpscale(body.Index, body.DiscordMsgId, body.MsgHash)
		Msg.Content = fmt.Sprintf("对任务：%s，第：%d,放大处理", body.DiscordMsgId, body.Index)
	case "variation":
		err = utils.ImageVariation(body.Index, body.DiscordMsgId, body.MsgHash)
		Msg.Content = fmt.Sprintf("以任务：%s，第：%d为准重绘处理", body.DiscordMsgId, body.Index)
	case "maxUpscale":
		err = utils.ImageMaxUpscale(body.DiscordMsgId, body.MsgHash)
		Msg.Content = fmt.Sprintf("以任务：%s为准重绘处理", body.DiscordMsgId)
	case "reset":
		err = utils.ImageReset(body.DiscordMsgId, body.MsgHash)
		Msg.Content = fmt.Sprintf("以任务：%s 为准重绘处理", body.DiscordMsgId)
	case "describe":
		//if len(body.Prompt) == 0 {
		//	err = errors.New("invalid Prompt")
		//	break
		//}
		//body.Prompt = fmt.Sprintf("[%s_%d_%s] %s ~~", global.GVA_CONFIG.Discord.MSG_PRFIX, uid, uuidStr, body.Prompt)
		//
		err = utils.ImageDescribe(body.Prompt)

	case "blend":
		//if len(body.Prompt) == 0 {
		//	err = errors.New("invalid Prompt")
		//	break
		//}
		//body.Prompt = fmt.Sprintf("[%s_%d_%s] %s ~~", global.GVA_CONFIG.Discord.MSG_PRFIX, uid, uuidStr, body.Prompt)
		//
		err = utils.ImageBlend(body.Prompt)

	default:
		err = errors.New("invalid type")
	}
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	Msgs := make([]interface{}, 0)
	Msgs = append(Msgs, Msg)
	Msg.IsSender = false
	Msgs = append(Msgs, Msg)

	if err := global.GVA_MONGO.InsertMany("mdj", "mdj_chat", Msgs...); err != nil {
		//fmt.Println("插入失败", err.Error())
		response.FailWithMessage(err.Error(), c)
		return
	} else {

	}

	response.OkWithMessage("success", c)

}
func GetImgHash(msgId string) (string, error) {

	filter := bson.D{}
	filter = append(filter, bson.E{"_id", msgId})

	Message := &ai.MdjMsg2{}

	_ = global.GVA_MONGO.FindOne("mdj", "mdj_msg2", filter, &Message)
	return "", nil
}

func (a MidBotApi) Health(c *gin.Context) {
	response.OkWithMessage("success", c)
}
func (a MidBotApi) Ok(c *gin.Context) {
	response.OkWithMessage("success", c)
}

type ReqUploadFile struct {
	ImgData []byte `form:"imgData"`
	Name    string `form:"name"`
	Size    int64  `form:"size"`
}

func (a MidBotApi) Upload(c *gin.Context) {

	//ctx, cancel := context.WithTimeout(c, time.Second*10) // 设置10秒的超时时间
	//defer cancel()                                        // 在处理函数返回之前取消上下文
	//
	//// 在这里编写你的处理逻辑
	//
	//if ctx.Err() == context.DeadlineExceeded {
	//	// 超时处理逻辑
	//	c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{"error": "请求超时"})
	//	return
	//}
	//var body ReqUploadFile
	//if err := c.ShouldBind(&body); err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	//需要把上传的名称保存起来，根据文件名来确认describe的数据是谁的，文件名需要uuid

	file, err := c.FormFile("file")
	if err != nil {
		// 处理错误
	}

	// 读取文件内容
	f, err := file.Open()
	if err != nil {
		// 处理错误
		return
	}
	defer f.Close()

	// 将文件内容转换为二进制数据
	data2 := make([]byte, file.Size)
	if _, err := io.ReadFull(f, data2); err != nil {
		// 处理错误
		return
	}
	//var body ReqUploadFile
	//if err := c.ShouldBindWith(&body, binding.Form); err != nil {
	//	// 处理错误
	//	return
	//}

	data, err := services.Attachments(file.Filename, file.Size)
	if err != nil {
		//c.String(http.StatusInternalServerError, "Error saving image: %s", err)
		response.FailWithMessage(err.Error(), c)
		return
	}
	if len(data.Attachments) == 0 {
		//c.String(http.StatusInternalServerError, "上传图片失败: %s", err)
		response.FailWithMessage(err.Error(), c)
		return
	}
	payload := bytes.NewReader(data2)
	client := &http.Client{}
	fmt.Printf("Attachments:%+v\n", data.Attachments)
	req, err := http.NewRequest("PUT", data.Attachments[0].UploadUrl, payload)
	fmt.Println("req:", req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.Header.Add("Content-Type", "image/png")

	res, err := client.Do(req)

	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	defer res.Body.Close()

	response.OkWithDetailed(map[string]interface{}{
		"filename": data.Attachments[0].UploadFilename,
	}, "success", c)
}

func (a MidBotApi) Messages(c *gin.Context) {
	//Messages := make([]discordgo.Message, 0)
	Messages, err := utils.GetMessages()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
	}

	response.OkWithDetailed(Messages, "获取成功", c)
}

func (a MidBotApi) List(c *gin.Context) {

	//ShouldBindQuery和ShouldBindJSON要搞清楚
	var reqInfo request.MdjSearch

	err := c.ShouldBindJSON(&reqInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	opts := &options.FindOptions{}
	filter := bson.D{}

	total, err := global.GVA_MONGO.QueryCount("mdj", "mdj_msg2", filter, 0)

	//liveInfo.settlementAmountThirtyRecent

	//Message := []discordgo.Message{}

	Messages := make([]ai.MdjMsg2, 0)

	_, err = global.GVA_MONGO.FindWithOpts("mdj", "mdj_msg2", int64(reqInfo.Page-1), int64(reqInfo.PageSize), filter, opts.SetSort(bson.M{"create_time": -1}), &Messages)

	sort.Slice(Messages, func(i, j int) bool {
		return Messages[i].CreateTime < Messages[j].CreateTime
	})

	response.OkWithDetailed(response.PageResult{
		List:     Messages,
		Total:    total,
		Page:     int(reqInfo.Page),
		PageSize: int(reqInfo.PageSize),
	}, "获取成功", c)
}

type FindForm struct {
	CreateTime string `json:"create_time" bson:"create_time"`
}

func (a MidBotApi) Find(c *gin.Context) {
	//ShouldBindQuery和ShouldBindJSON要搞清楚
	var reqInfo FindForm

	err := c.ShouldBindJSON(&reqInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	filter := bson.D{}
	filter = append(filter, bson.E{"create_time", bson.D{
		{"$gt", reqInfo.CreateTime},
	}})

	//Message := &ai.MdjMsg2{}
	//
	//_ = global.GVA_MONGO.FindOne("mdj", "mdj_msg2", filter, &Message)
	Messages := make([]ai.MdjMsg2, 0)

	_, err = global.GVA_MONGO.FindWithOrder("mdj", "mdj_msg2", filter, map[string]int{
		"create_time": -1,
	}, &Messages)

	response.OkWithDetailed(Messages, "获取成功", c)
}
