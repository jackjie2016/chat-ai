package ai

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"net/http"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type MdjApi struct{}
type ImagineRequest struct {
	Param *Param `json:"param"`
	Type  string `json:"type"`
}
type Param struct {
	MsgID    string `json:"msg_Id"`
	Question string `json:"question"`
}

type RequesJson struct {
	Content string `json:"content"`
	Index   *int   `json:"index"`
	Id      string `json:"id"`
	Url     string `json:"url"`
}

const MsgID = "minong00001"

// CreateApi
// @Tags      SysApi
// @Summary   创建基础api
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysApi                  true  "api路径, api中文描述, api组, 方法"
// @Success   200   {object}  response.Response{msg=string}  "创建基础api"
// @Router    /api/createApi [post]
func (s *MdjApi) CreateApi(c *gin.Context) {

	response.OkWithMessage("创建成功", c)
}

func (s *MdjApi) RetrieveMessages(ctx *gin.Context) {

	var request RequesJson
	err := ctx.BindJSON(&request)
	if err != nil {

		response.FailWithMessage(err.Error(), ctx)

		return
	}

	Msg := &ai.MdjMsg{}
	err = global.GVA_MONGO.FindOne("mdj", "mdj_msg", bson.D{{"_id", request.Content}}, Msg)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	MdjMsg := SearchMdjMsg(Msg, "s")

	if MdjMsg == nil {
		response.FailWithMessage("未找到", ctx)
		return
	}

	response.OkWithDetailed(MdjMsg, "创建成功", ctx)

}

func (s *MdjApi) Imagine(ctx *gin.Context) {
	var request RequesJson
	err := ctx.BindJSON(&request)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		//c.ResponseJson(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	verify := utils.Rules{
		"content": {utils.NotEmpty()},
	}

	if err := utils.Verify(request, verify); err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	if request.Content == "" {
		response.FailWithMessage("请输入内容", ctx)
		return
	}

	uuid := uuid.New().String()
	ParamData := Param{
		MsgID:    fmt.Sprintf("%s-%s", MsgID, uuid),
		Question: request.Content,
	}

	requestJson := &ImagineRequest{
		Param: &ParamData,
		Type:  "imagine",
	}
	midjourney, err := Midjourney(requestJson)
	if err != nil {
		response.FailWithMessage("第三方请求失败", ctx)
		return
	}

	var resp bool
	err = json.Unmarshal(midjourney, &resp)
	if err != nil {
		return
	}

	fmt.Println(222)
	if !resp {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	MdjMsg := &ai.MdjMsg{
		ID:         uuid,
		RootHash:   uuid,
		UserId:     0,
		Status:     false,
		Qustion:    request.Content,
		CreateTime: time.Now().Format("2006-01-02 15:04:02"),
	}

	if err := global.GVA_MONGO.InsertMany("mdj", "mdj_msg", *MdjMsg); err != nil {
		fmt.Println("插入失败")
		response.FailWithMessage(err.Error(), ctx)
		return
	} else {
		fmt.Println("插入成功")
	}

	go FindMsgTask(MdjMsg, "")

	response.OkWithDetailed(uuid, "", ctx)

}

type ChangeResponse struct {
	ID       string `json:"id"`
	URI      string `json:"uri"`
	Hash     string `json:"hash"`
	Content  string `json:"content"`
	Progress string `json:"progress"`
}

func (s *MdjApi) Upscale(ctx *gin.Context) {

	var request RequesJson
	err := ctx.BindJSON(&request)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)

		return
	}
	msgs := ai.MdjMsg{}
	err = global.GVA_MONGO.FindOne("mdj", "mdj_msg", bson.D{bson.E{"mdjID", request.Content}}, &msgs)
	if !msgs.Status {
		response.FailWithMessage("当前的不能放大", ctx)
		return
	}
	UpscaleParam := ai.UpscaleParam{
		ID:       msgs.MdjID,
		Question: msgs.Qustion,
		Index:    *request.Index,
		URL:      msgs.ProxyURL,
	}
	upscaleRequest := &ai.UpscaleRequest{
		Type:  "upscale",
		Param: &UpscaleParam,
	}

	midjourney, err := Midjourney(upscaleRequest)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	fmt.Println(string(midjourney))
	var midMessages ChangeResponse
	err = json.Unmarshal(midjourney, &midMessages)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	MdjMsg := ai.MdjMsg{
		ID:         uuid.New().String(),
		PID:        msgs.ID,
		RootHash:   msgs.RootHash,
		ChangeType: "upscale",

		ParentID: msgs.MdjID,
		Index:    *request.Index,

		OrgURL: msgs.URL,

		UserId: 0,

		Status: true,

		Qustion: msgs.Qustion,

		MdjID:      midMessages.ID,
		URL:        midMessages.URI,
		Content:    fmt.Sprintf("%s Image # %d  （PID：%s） ", midMessages.Content, *request.Index, msgs.PID),
		CreateTime: time.Now().Format("2006-01-02 15:04:02"),
	}

	fmt.Printf("%+v", MdjMsg)
	if err := global.GVA_MONGO.InsertMany("mdj", "mdj_msg", MdjMsg); err != nil {
		fmt.Println("插入失败")
		response.FailWithMessage(err.Error(), ctx)
		return
	} else {
		fmt.Println("插入成功")
	}

	response.OkWithDetailed(midMessages, "", ctx)
}

func (s *MdjApi) Variation(ctx *gin.Context) {

	var request RequesJson
	err := ctx.BindJSON(&request)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	msgs := ai.MdjMsg{}
	err = global.GVA_MONGO.FindOne("mdj", "mdj_msg", bson.D{bson.E{"mdjID", request.Content}}, &msgs)
	if !msgs.Status {
		response.FailWithMessage("当前的不能放大", ctx)
		return
	}
	VariationParam := ai.UpscaleParam{
		ID:       msgs.MdjID,
		Question: msgs.Qustion,
		Index:    *request.Index,
		URL:      msgs.ProxyURL,
	}
	upscaleRequest := &ai.UpscaleRequest{
		Type:  "variation",
		Param: &VariationParam,
	}

	midjourney, err := Midjourney(upscaleRequest)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	var midMessages ChangeResponse
	err = json.Unmarshal(midjourney, &midMessages)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	MdjMsg := ai.MdjMsg{
		ID:         uuid.New().String(),
		PID:        msgs.ID,
		RootHash:   msgs.RootHash,
		ChangeType: "variation",

		ParentID: msgs.MdjID,
		Index:    *request.Index,

		OrgURL: msgs.URL,

		UserId: 0,

		Status: true,

		Qustion: msgs.Qustion,

		MdjID:      midMessages.ID,
		URL:        midMessages.URI,
		Content:    fmt.Sprintf("%s Variations  # %d （PID：%s） ", midMessages.Content, *request.Index, msgs.PID),
		CreateTime: time.Now().Format("2006-01-02 15:04:02"),
	}

	if err := global.GVA_MONGO.InsertMany("mdj", "mdj_msg", MdjMsg); err != nil {
		fmt.Println("插入失败")
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	response.OkWithDetailed(midMessages, "", ctx)
}

func (s *MdjApi) Msglist(c *gin.Context) {

	//ShouldBindQuery和ShouldBindJSON要搞清楚
	var reqInfo request.MdjSearch

	err := c.ShouldBindJSON(&reqInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	opts := &options.FindOptions{}
	filter := bson.D{}

	if reqInfo.Keyword != "" {

		pattern := fmt.Sprintf(".*%s.*", reqInfo.Keyword)

		regex := primitive.Regex{Pattern: pattern, Options: "i"}

		filter = append(filter, bson.E{"content", regex})

	}
	//if reqInfo.Status {
	//	filter = append(filter, bson.E{"status", reqInfo.Status})
	//}

	total, err := global.GVA_MONGO.QueryCount("mdj", "mdj_msg", filter, 0)
	fmt.Printf("%+v\n", filter)
	//liveInfo.settlementAmountThirtyRecent

	msgs := []*ai.MdjMsg{}

	_, err = global.GVA_MONGO.FindWithOpts("mdj", "mdj_msg", int64(reqInfo.Page-1), int64(reqInfo.PageSize), filter, opts.SetSort(bson.M{"create_time": -1}), &msgs)

	response.OkWithDetailed(response.PageResult{
		List:     msgs,
		Total:    total,
		Page:     int(reqInfo.Page),
		PageSize: int(reqInfo.PageSize),
	}, "获取成功", c)
}

func Midjourney(request interface{}) (responseData []byte, err error) {

	url := "https://rpcqmo.laf.dev/mj-service"
	method := "POST"

	var payload *strings.Reader
	client := &http.Client{}
	req := &http.Request{}
	var goodsJson []byte
	goodsJson, err = json.Marshal(request)
	fmt.Println(string(goodsJson))
	payload = strings.NewReader(string(goodsJson))
	req, err = http.NewRequest(method, url, payload)
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")

	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	return body, err
}

func UploadToOss(Id, url string) (ossUrl string) {
	fmt.Println("UploadToOss 1", url)
	// OSS配置信息
	endpoint := "oss-ap-southeast-1.aliyuncs.com"
	accessKeyID := "LTAI5tN9YEpdtL4N94mWRDYk"
	accessKeySecret := "5mISqOyU0g6EX5wK5gXPZwmfww9jtS"
	bucketName := "minong-mj"

	objectName := global.GVA_CONFIG.AliyunOSS.BasePath + "/" + "uploads" + "/" + time.Now().Format("2006-01-02") + "/"

	fileName := path.Base(url)
	objectName = objectName + fileName
	// 创建一个自定义的Transport，用于支持HTTPS请求
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// 创建一个Client，指定Transport
	httpclient := &http.Client{Transport: tr}

	// 发起HTTPS GET请求
	resp, err := httpclient.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	fmt.Println("读取成功了吗")
	// 将图片内容读取到内存
	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error while reading", err)
		return
	}

	// 创建OSS Client对象
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		fmt.Println("Error while creating OSS client", err)
		return
	}

	// 获取Bucket对象
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		fmt.Println("Error while getting bucket", err)
		return
	}

	// 上传图片到OSS
	err = bucket.PutObject(objectName, bytes.NewReader(imageBytes))
	if err != nil {
		fmt.Println("Error while uploading", err)
		return
	}
	fmt.Println("UploadToOss ok to mongodb")
	msg := make(map[string]string)
	msg["ossUrl"] = global.GVA_CONFIG.AliyunOSS.BucketUrl + "/" + objectName
	if err := global.GVA_MONGO.UpdateOne("mdj", "mdj_msg", bson.D{{"_id", Id}}, msg); err != nil {
		fmt.Println("更新失败")
		panic(err)
	}
	fmt.Println("all ok")
	return global.GVA_CONFIG.AliyunOSS.BucketUrl + "/" + objectName
}
func FindMsgTask(resMsg *ai.MdjMsg, Msgtype string) {
	i := 1
	for {
		if i == 11 {
			fmt.Println("找不到")
			break
		}
		i = i + 1
		fmt.Println(fmt.Sprintf("%d", i))
		MdjMsg := SearchMdjMsg(resMsg, Msgtype)
		if MdjMsg != nil && MdjMsg.Status {
			break
		}

		time.Sleep(60 * time.Second)

	}
	fmt.Println("结束了")

}
func SearchMdjMsg(resMsg *ai.MdjMsg, Msgtype string) (MdjMsg *ai.MdjMsg) {
	fmt.Println(resMsg, Msgtype)
	requestJson := &ImagineRequest{
		Type: "RetrieveMessages",
	}

	midjourney, err := Midjourney(requestJson)
	if err != nil {
		return
	}
	var midMessages []ai.MidMessages
	err = json.Unmarshal(midjourney, &midMessages)
	if err != nil {
		return
	}

	for _, v := range midMessages {
		//fmt.Println(fmt.Sprintf("%s-%s", MsgID, uuid))

		//查找文本，根据第一级父节点
		if strings.Contains(v.Content, fmt.Sprintf("%s-%s", MsgID, resMsg.RootHash)) {
			switch Msgtype {
			case "u":
				re := regexp.MustCompile(`\(Image #%d\)`)
				match := re.FindStringSubmatch(v.Content)
				fmt.Println(match)
			case "v":
				re := regexp.MustCompile(`\(Image #%d\)`)
				match := re.FindStringSubmatch(v.Content)
				fmt.Println(match)
			}
			if len(v.Attachments) > 0 {
				re := regexp.MustCompile(`\((\d+)%\)`)
				match := re.FindStringSubmatch(v.Content)

				Status := false
				OssURL := ""

				if len(match) == 2 {
					Status = false
				}

				if len(match) == 0 && v.Attachments[0].ProxyURL != "" {
					Status = true
					fmt.Println(" OssURL 上传开始")
					//go UploadToOss(v.ID, v.Attachments[0].ProxyURL)
					go UploadToOss(v.ID, v.Attachments[0].URL)
					fmt.Println(" OssURL 上传", OssURL)
				}

				MdjMsg = &ai.MdjMsg{
					ID:          resMsg.ID,
					MdjID:       v.ID,
					RootHash:    resMsg.RootHash,
					ParentID:    v.MessageReference.MessageId,
					UserId:      0,
					AttachID:    v.Attachments[0].ID,
					Type:        v.Type,
					Content:     v.Content,
					Filename:    v.Attachments[0].Filename,
					Size:        v.Attachments[0].Size,
					URL:         v.Attachments[0].URL,
					ProxyURL:    v.Attachments[0].ProxyURL,
					OssURL:      OssURL,
					Width:       v.Attachments[0].Width,
					Height:      v.Attachments[0].Height,
					ContentType: v.Attachments[0].ContentType,
					Status:      Status,
					CreateTime:  time.Now().Format("2006-01-02 15:04:02"),
					UpdateTime:  time.Now().Format("2006-01-02 15:04:02"),
				}
				//fmt.Println("搜索resMsg：", resMsg.ID)
				//fmt.Printf("搜索到这边：%+v", MdjMsg)

			} else {
				MdjMsg = &ai.MdjMsg{
					ID:         resMsg.ID,
					RootHash:   resMsg.RootHash,
					MdjID:      v.ID,
					UserId:     0,
					Type:       v.Type,
					Content:    v.Content,
					CreateTime: time.Now().Format("2006-01-02 15:04:02"),
					UpdateTime: time.Now().Format("2006-01-02 15:04:02"),
				}
				//fmt.Println("找到了，但是图片还没有出来")

			}
			MdjMsg2 := &ai.MdjMsg{}
			err := global.GVA_MONGO.FindOne("mdj", "mdj_msg", bson.D{{"_id", MdjMsg.ID}}, MdjMsg2)
			if err != nil {
				return nil
			}
			fmt.Println(MdjMsg.ID)
			fmt.Printf("%+v", MdjMsg2)
			if MdjMsg2 != nil {
				//判断是不是 放大或者 重绘的，放大重绘的绑定父节点,要更新到指定的节点

				MdjMsg.Qustion = MdjMsg2.Qustion
				MdjMsg.CreateTime = MdjMsg2.CreateTime
				if err := global.GVA_MONGO.UpdateOne("mdj", "mdj_msg", bson.D{{"_id", MdjMsg.ID}}, MdjMsg); err != nil {
					fmt.Println("更新失败")
					panic(err)
				}
			} else {
				if err := global.GVA_MONGO.InsertMany("mdj", "mdj_msg", MdjMsg); err != nil {
					fmt.Println("插入失败")
					panic(err)
				}
			}

			break

		}
	}
	return MdjMsg
}
