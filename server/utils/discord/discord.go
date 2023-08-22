package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	discord "github.com/bwmarrin/discordgo"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type Scene string

const (
	/**
	 * 首次触发生成
	 */
	FirstTrigger Scene = "FirstTrigger"
	/**
	 * 生成图片结束
	 */
	GenerateEnd Scene = "GenerateEnd"
	/**
	 * 发送的指令midjourney生成过程中发现错误
	 */
	GenerateEditError Scene = "GenerateEditError"
	/**
	 * 富文本
	 */
	RichText Scene = "RichText"
	/**
	 * 发送的指令midjourney直接报错或排队阻塞不在该项目中处理 在业务服务中处理
	 * 例如：首次触发生成多少秒后没有回调业务服务判定会指令错误或者排队阻塞
	 */
)

func DiscordMsgCreate(s *discord.Session, m *discord.MessageCreate) {
	// 过滤频道
	if m.ChannelID != global.GVA_CONFIG.Discord.DISCORD_CHANNEL_ID {
		return
	}

	// 过滤掉自己发送的消息
	if m.Author.ID == s.State.User.ID {
		return
	}

	/******** *********/
	if data, err := json.Marshal(m); err == nil {
		fmt.Println("discord message: ", string(data))
	}
	/******** *********/
	fmt.Println(m)
	/******** 提示词，首次触发 start ********/
	// 重新生成不发送
	// TODO 优化 使用 From
	if strings.Contains(m.Content, "(Waiting to start)") && !strings.Contains(m.Content, "Rerolling **") {
		trigger(m.Content, FirstTrigger)
		return
	}
	/******** end ********/

	/******** 图片生成回复 start ********/
	for _, attachment := range m.Attachments {
		if attachment.Width > 0 && attachment.Height > 0 {
			replay(m)
			return
		}
	}
	/******** end ********/
}

func DiscordMsgUpdate(s *discord.Session, m *discord.MessageUpdate) {
	// 过滤频道
	if m.ChannelID != global.GVA_CONFIG.Discord.DISCORD_CHANNEL_ID {
		return
	}

	if m.Author == nil {
		return
	}

	// 过滤掉自己发送的消息
	if m.Author.ID == s.State.User.ID {
		return
	}

	/******** *********/
	if data, err := json.Marshal(m); err == nil {
		fmt.Println("\ndiscord message update: ", string(data))
		//msg := &discord.MessageCreate{}
		//err := json.Unmarshal(data, msg)
		//if err == nil {
		//	oparetMongoDB(msg)
		//}

	}
	/******** *********/

	/******** 发送的指令midjourney生成发现错误 ********/
	if strings.Contains(m.Content, "(Stopped)") {
		trigger(m.Content, GenerateEditError)
		return
	}

	if len(m.Embeds) > 0 {
		send(m.Embeds)
		return
	}
}

type ReqCb struct {
	Embeds  []*discord.MessageEmbed `json:"embeds,omitempty"`
	Discord *discord.MessageCreate  `json:"discord,omitempty"`
	Content string                  `json:"content,omitempty"`
	Type    Scene                   `json:"type"`
}

func replay(m *discord.MessageCreate) {
	body := ReqCb{
		Discord: m,
		Type:    GenerateEnd,
	}
	//OparetMongoDB(m)
	request(body)
}
func OparetMongoDB(m *discord.MessageCreate) {
	MdjMsg := &ai.MdjMsg2{}
	err := global.GVA_MONGO.FindOne("mdj", "mdj_msg2", bson.D{{"_id", m.ID}}, MdjMsg)
	if err != nil {
		fmt.Println("查询失败：", err.Error())
	} else {
		fmt.Printf("查询结果：%+v", MdjMsg)
	}

	ImgHash := ""
	if len(m.Attachments) > 0 && m.Attachments[0].URL != "" {
		ImgHash = getImgHash(m.Attachments[0].URL, 0)
	}

	if len(m.Embeds) > 0 {
		fmt.Println("当前是图生文", m.Embeds[0].Description)

	}

	if MdjMsg.ID != "" {
		fmt.Println("更新", MdjMsg)
		if MdjMsg.ProxyURL == "" {
			MdjMsg.MsgId = m.ID
			MdjMsg.Content = m.Content
			MdjMsg.ProxyURL = m.Attachments[0].ProxyURL
			MdjMsg.URL = m.Attachments[0].URL
			MdjMsg.URL = m.Attachments[0].URL
			MdjMsg.Size = m.Attachments[0].Size
			MdjMsg.Filename = m.Attachments[0].Filename
			MdjMsg.Filename = m.Attachments[0].Filename
			MdjMsg.Width = m.Attachments[0].Width
			MdjMsg.Height = m.Attachments[0].Height
			MdjMsg.UpdateTime = time.Now().Format("2006-01-02 15:04:02")
			if err := global.GVA_MONGO.UpdateOne("mdj", "mdj_msg2", bson.D{{"_id", MdjMsg.ID}}, MdjMsg); err != nil {
				fmt.Println("更新失败")
				panic(err)
			}
		}
	} else {
		fmt.Println("执行插入：", MdjMsg)

		uid, hashId, err := getUIDAndHash(m.Content)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if len(m.Attachments) > 0 && m.Attachments[0].ProxyURL != "" {

			MdjMsg = &ai.MdjMsg2{
				ID:        m.ID,
				MsgId:     m.ID,
				ChannelId: m.ChannelID,
				GuildId:   m.GuildID,
				AuthorId:  m.Author.ID,
				ParentID:  "",
				Index:     0,
				UserId:    int(uid),
				HashID:    hashId,
				AttachID:  m.Attachments[0].ID,
				Type:      0,
				Content:   m.Content,
				Filename:  m.Attachments[0].Filename,
				Size:      m.Attachments[0].Size,
				OrgURL:    "",
				URL:       m.Attachments[0].URL,
				ImgHash:   ImgHash,
				ProxyURL:  m.Attachments[0].ProxyURL,

				Width:       m.Attachments[0].Width,
				Height:      m.Attachments[0].Height,
				ContentType: "",
				CreateTime:  time.Now().Format("2006-01-02 15:04:02"),
				UpdateTime:  time.Now().Format("2006-01-02 15:04:02"),
			}
		} else {
			MdjMsg = &ai.MdjMsg2{
				ID:          m.ID,
				MsgId:       m.ID,
				ChannelId:   m.Content,
				GuildId:     m.GuildID,
				AuthorId:    m.Author.ID,
				ParentID:    "",
				Index:       0,
				UserId:      int(uid),
				HashID:      hashId,
				AttachID:    "",
				Type:        0,
				Content:     m.Content,
				Filename:    "",
				Size:        0,
				OrgURL:      "",
				URL:         "",
				ProxyURL:    "",
				OssURL:      "",
				Width:       0,
				Height:      0,
				ContentType: "",
				CreateTime:  time.Now().Format("2006-01-02 15:04:02"),
				UpdateTime:  time.Now().Format("2006-01-02 15:04:02"),
			}

		}
		fmt.Println("执行插入：", MdjMsg)
		if err := global.GVA_MONGO.InsertMany("mdj", "mdj_msg2", MdjMsg); err != nil {
			fmt.Println("插入失败", err.Error())

			return
		}
	}

	if MdjMsg.URL != "" && MdjMsg.ImgHash != "" {
		go UploadToOss(MdjMsg)
	}
}

func InsertMongoDB(m *discord.Message) {
	//fmt.Printf("InsertMongoDB:%p,%+v\n", m, m.ID)
	//fmt.Println(m.ID, m.Content)
	if m.ID == "" {
		return
	}
	ImgHash := ""
	if len(m.Attachments) > 0 && m.Attachments[0].URL != "" {
		ImgHash = getImgHash(m.Attachments[0].URL, 0)
	}

	if len(m.Embeds) > 0 {
		fmt.Println("当前是图生文", m.Embeds[0].Description)

	}
	//fmt.Println(m)
	if !CheckTask(m.Content) {
		fmt.Println("任务未完成", m.Content)
		return
	}
	MdjMsg := &ai.MdjMsg2{}
	err := global.GVA_MONGO.FindOne("mdj", "mdj_msg2", bson.D{{"_id", m.ID}}, MdjMsg)
	if err != nil {
		return
	}

	if MdjMsg.ID != "" {
		ParentID := ""
		if m.ReferencedMessage != nil {
			ParentID = m.ReferencedMessage.ID
		}
		ChangeType := "generate"
		if strings.Contains(m.Content, "Variations") {
			ChangeType = "variation"
		}
		if strings.Contains(m.Content, "Image #") {
			ChangeType = "upscale"
		}

		//UpCheck := regexp.MustCompile(`\(Image #%d\)`)
		//UpCheckmatch := UpCheck.FindStringSubmatch(m.Content)
		//fmt.Println(UpCheckmatch)
		//fmt.Println("更新", MdjMsg)
		if len(m.Attachments) > 0 && MdjMsg.ProxyURL != "" {
			MdjMsg.MsgId = m.ID
			MdjMsg.ChangeType = ChangeType
			MdjMsg.Content = m.Content
			MdjMsg.ProxyURL = m.Attachments[0].ProxyURL
			MdjMsg.ParentID = ParentID

			MdjMsg.URL = m.Attachments[0].URL
			MdjMsg.Size = m.Attachments[0].Size
			MdjMsg.Filename = m.Attachments[0].Filename
			MdjMsg.Filename = m.Attachments[0].Filename
			MdjMsg.Width = m.Attachments[0].Width
			MdjMsg.Height = m.Attachments[0].Height
			MdjMsg.Embeds = m.Embeds
			MdjMsg.ImgHash = ImgHash

			MdjMsg.UpdateTime = time.Now().Format("2006-01-02 15:04:02")
			if err := global.GVA_MONGO.UpdateOne("mdj", "mdj_msg2", bson.D{{"_id", MdjMsg.ID}}, MdjMsg); err != nil {
				fmt.Println("更新失败")
				panic(err)
			} else {
				//fmt.Println("更新成功")
			}
		}
	} else {
		//fmt.Println(" 插入start", MdjMsg, m)
		uid, hashId, _ := getUIDAndHash(m.Content)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if len(m.Attachments) > 0 && m.Attachments[0].ProxyURL != "" {
			ParentID := ""
			if m.ReferencedMessage != nil {
				ParentID = m.ReferencedMessage.ID
			}
			ChangeType := "generate"
			if strings.Contains(m.Content, "Variations") {
				ChangeType = "variation"
			}
			if strings.Contains(m.Content, "Image #") {
				ChangeType = "upscale"
			}

			MdjMsg = &ai.MdjMsg2{
				ID:          m.ID,
				MsgId:       m.ID,
				ChannelId:   m.ChannelID,
				GuildId:     m.GuildID,
				AuthorId:    m.Author.ID,
				ParentID:    ParentID,
				ChangeType:  ChangeType,
				Index:       0,
				UserId:      int(uid),
				HashID:      hashId,
				AttachID:    m.Attachments[0].ID,
				Type:        int(m.Type),
				Content:     m.Content,
				Filename:    m.Attachments[0].Filename,
				Size:        m.Attachments[0].Size,
				OrgURL:      "",
				URL:         m.Attachments[0].URL,
				ImgHash:     ImgHash,
				ProxyURL:    m.Attachments[0].ProxyURL,
				Embeds:      m.Embeds,
				Width:       m.Attachments[0].Width,
				Height:      m.Attachments[0].Height,
				ContentType: "",
				CreateTime:  time.Now().Format("2006-01-02 15:04:02"),
				UpdateTime:  time.Now().Format("2006-01-02 15:04:02"),
				IsRead:      false,
			}

			if err := global.GVA_MONGO.InsertMany("mdj", "mdj_msg2", MdjMsg); err != nil {
				fmt.Println("插入失败", err.Error())
				return
			}
		}

	}

	if MdjMsg.URL != "" && MdjMsg.ImgHash != "" && MdjMsg.OssURL == "" {
		go UploadToOss(MdjMsg)
	}
}
func UploadToOss(msg *ai.MdjMsg2) {
	fmt.Println("UploadToOss 1", msg.ID, msg.URL)
	// OSS配置信息
	endpoint := "oss-ap-southeast-1.aliyuncs.com"
	accessKeyID := "LTAI5tN9YEpdtL4N94mWRDYk"
	accessKeySecret := "5mISqOyU0g6EX5wK5gXPZwmfww9jtS"
	bucketName := "minong-mj"

	objectName := global.GVA_CONFIG.AliyunOSS.BasePath + "/" + "uploads" + "/" + time.Now().Format("2006-01-02") + "/"

	fileName := path.Base(msg.URL)
	objectName = objectName + fileName
	// 创建一个自定义的Transport，用于支持HTTPS请求
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// 创建一个Client，指定Transport
	httpclient := &http.Client{Transport: tr}

	// 发起HTTPS GET请求
	resp, err := httpclient.Get(msg.URL)
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

	// 设置更新内容
	//update := bson.M{"$set": bson.M{"ossUrl": global.GVA_CONFIG.AliyunOSS.BucketUrl + "/" + objectName, "UpdateTime": time.Now().Format("2006-01-02 15:04:02")}}
	msg.OssURL = global.GVA_CONFIG.AliyunOSS.BucketUrl + "/" + objectName
	msg.UpdateTime = time.Now().Format("2006-01-02 15:04:02")
	if err := global.GVA_MONGO.UpdateOne("mdj", "mdj_msg2", bson.D{{"_id", msg.ID}}, msg); err != nil {
		fmt.Println("更新失败")
		panic(err)
	}
	fmt.Println("all ok")
	//return global.GVA_CONFIG.AliyunOSS.BucketUrl + "/" + objectName
}
func CheckTask(str string) bool {
	re := regexp.MustCompile(`\((\d+)%\)`) // 匹配括号中的数字
	match := re.FindStringSubmatch(str)
	if len(match) > 1 {
		return false
	}

	if strings.Contains(str, "(Waiting to start)") {
		return false
	}
	return true

}
func getImgHash(url string, ctype int) (Imghash string) {

	rule0 := `__([a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}).`
	rule1 := `_([a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}).`

	re := regexp.MustCompile(rule0)
	match := re.FindStringSubmatch(url)

	if len(match) > 1 {
		return match[1]
	} else {
		re := regexp.MustCompile(rule1)
		match := re.FindStringSubmatch(url)
		if len(match) > 1 {
			return match[1]
		}
	}
	return ""
}
func getUIDAndHash(s string) (uint64, string, error) {
	r := regexp.MustCompile(fmt.Sprintf(`\[%s_(\d+)_(.*)\]`, global.GVA_CONFIG.Discord.MSG_PRFIX))
	match := r.FindStringSubmatch(s)
	fmt.Println(match)
	if len(match) != 3 {
		return 0, "", fmt.Errorf("no uid or hash found")
	}

	uidStr := match[1]
	hashStr := match[2]

	uid, err := strconv.ParseUint(uidStr, 10, 64)
	if err != nil {
		return 0, "", err
	}

	return uid, hashStr, nil
}
func getUid(content string) (int, error) {
	re := regexp.MustCompile(fmt.Sprintf(`\[%s_(\d+)_(.*)\]`, global.GVA_CONFIG.Discord.MSG_PRFIX))
	match := re.FindStringSubmatch(content)
	if len(match) > 1 {
		return strconv.Atoi(match[1])
	}
	return 0, nil
}

func send(embeds []*discord.MessageEmbed) {
	body := ReqCb{
		Embeds: embeds,
		Type:   RichText,
	}
	request(body)
}

func trigger(content string, t Scene) {
	body := ReqCb{
		Content: content,
		Type:    t,
	}
	request(body)
}

func request(params interface{}) {
	data, err := json.Marshal(params)
	if err != nil {
		fmt.Println("json marshal error: ", err)
		return
	}
	req, err := http.NewRequest("POST", global.GVA_CONFIG.Discord.CB_URL, strings.NewReader(string(data)))
	if err != nil {
		fmt.Println("http request error: ", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("http request error: ", err)
		return
	}
	defer resp.Body.Close()
}
