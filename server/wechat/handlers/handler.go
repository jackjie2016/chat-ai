package handlers

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/eatmoreapple/openwechat"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	services "github.com/flipped-aurora/gin-vue-admin/server/logic"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	utils "github.com/flipped-aurora/gin-vue-admin/server/utils/discord"
	"github.com/flipped-aurora/gin-vue-admin/server/wechat/gtp"
	"github.com/google/uuid"
	"path"

	"go.mongodb.org/mongo-driver/bson"
	"io"
	"log"
	"net/http"
	url2 "net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// MessageHandlerInterface 消息处理接口
type MessageHandlerInterface interface {
	handle(*openwechat.Message) error
	ReplyText(*openwechat.Message) error
}

type HandlerType string

const (
	GroupHandler = "group"
	UserHandler  = "user"
)

// handlers 所有消息类型类型的处理器
var handlers map[HandlerType]MessageHandlerInterface

func init() {
	handlers = make(map[HandlerType]MessageHandlerInterface)
	handlers[GroupHandler] = NewGroupMessageHandler()
	handlers[UserHandler] = NewUserMessageHandler()
}

// Handler 全局处理入口
func Handler(msg *openwechat.Message) {
	log.Printf("hadler Received msg : %v", msg.Content)
	// 处理群消息
	if msg.IsSendByGroup() {
		handlers[GroupHandler].handle(msg)
		return
	}

	// 好友申请
	//if msg.IsFriendAdd() {
	//	fmt.Println("IsFriendAdd")
	//	if config.LoadConfig().AutoPass {
	//		_, err := msg.Agree("你好我是基于chatGPT引擎开发的微信机器人，你可以向我提问任何问题。")
	//		if err != nil {
	//			log.Fatalf("add friend agree error : %v", err)
	//			return
	//		}
	//	}
	//}

	// 私聊
	handlers[UserHandler].handle(msg)
}
func AtText(msg *openwechat.Message) (atText string) {
	atText = ""
	if msg.IsSendByGroup() {

		// 获取@我的用户
		groupSender, err := msg.SenderInGroup()
		if err != nil {
			log.Printf("get sender in group error :%v \n", err)
			return ""
		}
		// 回复@我的用户

		atText = "@" + groupSender.NickName

	}
	return atText
}
func Ai(msg *openwechat.Message) (err error) {
	var reply string
	atText := AtText(msg)

	// 向GPT发起请求
	requestText := strings.TrimSpace(msg.Content)
	requestText = strings.Trim(msg.Content, "\n")

	if strings.HasPrefix(msg.Content, "/gpt ") {
		fmt.Println("The string starts with /gpt ")

		reply, err = gtp.Completions(requestText)
		if err != nil {
			log.Printf("gtp request error: %v \n", err)
			msg.ReplyText("机器人神了，我一会发现了就去修。")
			return err
		}
		reply = atText + reply
		msg.ReplyText(reply)
		return err
	} else if strings.HasPrefix(msg.Content, "/mj ") {
		fmt.Println("The string starts with /mj ")
		Prompt := strings.TrimPrefix(msg.Content, "/mj ")
		if len(Prompt) == 0 {
			err = errors.New("invalid Prompt")
			return err
		}

		uuidStr := uuid.New().String()
		Prompt = fmt.Sprintf("[%s_%d_%s] %s", global.GVA_CONFIG.Discord.MSG_PRFIX, 1, uuidStr, Prompt)
		err = utils.GenerateImage(Prompt)
		if err != nil {
			msg.ReplyText("出问题了哦:" + err.Error())
			return err
		}
		msg.Content = strings.TrimPrefix(msg.Content, "/mj ")
		go DiscordImagele(uuidStr, msg)

	} else if strings.HasPrefix(msg.Content, "/up ") {

		Prompt := strings.TrimPrefix(msg.Content, "/up ")
		if len(Prompt) == 0 {
			err = errors.New("invalid Prompt")
			return err
		}

		PromptSlince := strings.Split(msg.Content, " ")

		// 获取切片中的三个字符串
		if len(PromptSlince) == 3 {

			msgid := PromptSlince[1]
			uStr := PromptSlince[2]

			index, err := strconv.Atoi(string(uStr[1]))
			if err != nil {
				fmt.Println("转换失败:", err)
				return err
			}
			fmt.Println(msgid)
			msgs := FindMango(msgid)

			fmt.Println(msgs)
			if msgs.ImgHash == "" {
				msg.ReplyText("源文件找不到了哦")
				return err
			}
			ImgHash := msgs.ImgHash
			for s, v := range []byte(uStr) {
				fmt.Println(s, v)
			}
			fmt.Println(uStr[0])
			if uStr[0] == 'U' {
				err = utils.ImageUpscale(int64(index), msgid, ImgHash)
				if err != nil {
					msg.ReplyText("出问题了哦:" + err.Error())
					return err
				}

				go DiscordUp(msgid, index, msg)

			} else if uStr[0] == 'V' {
				fmt.Println("ImageVariation")
				err = utils.ImageVariation(int64(index), msgid, ImgHash)
				if err != nil {
					msg.ReplyText("出问题了哦:" + err.Error())
					return err
				}

				go DiscordVariation(msgid, index, msg)
			}

		}

	}
	return err
}
func DiscordImagele(uuidStr string, msg *openwechat.Message) {

	reply := ""
	i := 0
	//第一步找
	isfound := false
	doOk := false
	for {

		data, err := services.GetMessages()
		if err != nil {
			msg.ReplyText("出问题了哦:" + err.Error())
			return
		}
		if !isfound {
			if i > 60 && !isfound {
				msg.ReplyText("很遗憾，未找到")
				break
			}

			for _, v := range data {

				//查找文本，根据第一级父节点
				if strings.Contains(v.Content, uuidStr) {
					reply = fmt.Sprintf(`%s
✅ 您的任务已提交
/imagine %s
正在快速处理中，请稍后`, AtText(msg), msg.Content)
					msg.ReplyText(reply)
					isfound = true
					break
				}

			}
		} else {
			if !doOk {
				for _, v := range data {
					//查找文本，根据第一级父节点
					if strings.Contains(v.Content, uuidStr) {
						if !strings.Contains(v.Content, "Waiting to start") {
							re := regexp.MustCompile(`\((\d+)%\)`)
							match := re.FindStringSubmatch(v.Content)
							if len(match) == 0 {
								reply = fmt.Sprintf(`%s
绘图成功，用时 48秒
Prompt:%s
📨 任务ID: %s
彩 放大 U1～U4 ，变换 V1～V4
✏️ 使用[/up 任务ID 操作]
/up %s U1`, AtText(msg), msg.Content, v.ID, v.ID)

								msg.ReplyText(reply)

								InsertMangodb(&v, 0, msg)

								doOk = true
								break
							}

						}
					}

				}
			} else {
				break
			}

		}

		time.Sleep(5 * time.Second)
		i++
	}
	//第二部查看有没有完成
}

func DiscordUp(MsgID string, index int, msg *openwechat.Message) {

	reply := ""
	i := 0
	//第一步找
	isfound := false
	doOk := false
	for {

		data, err := services.GetMessages()
		if err != nil {
			msg.ReplyText("出问题了哦:" + err.Error())
			return
		}
		if !isfound {
			if i > 60 && !isfound {
				msg.ReplyText("很遗憾，未找到")
				break
			}

			for _, v := range data {

				//查找文本，根据第一级父节点
				if v.ReferencedMessage != nil && v.ReferencedMessage.ID == MsgID {

					if strings.Contains(v.Content, "Image") {

						re := regexp.MustCompile(`Image #(\d+)`)
						match := re.FindStringSubmatch(v.Content)

						if len(match) > 1 {
							number := match[1]
							num, _ := strconv.Atoi(number)
							if num == index {
								reply = fmt.Sprintf(`%s
 
✅ 您的任务已提交
 %s
正在快速处理中，请稍后`, AtText(msg), msg.Content)
								msg.ReplyText(reply)
								isfound = true
								break
							}
						}

					}

				}

			}
		} else {
			if !doOk {
				for _, v := range data {
					//查找文本，根据第一级父节点
					if v.ReferencedMessage != nil && v.ReferencedMessage.ID == MsgID {
						if strings.Contains(v.Content, "Image") {
							//re := regexp.MustCompile(`\(Image #%d\)`)
							//match := re.FindStringSubmatch(v.Content)
							//fmt.Println("match:", match)

							re := regexp.MustCompile(`Image #(\d+)`)
							match := re.FindStringSubmatch(v.Content)

							if len(match) > 1 {
								number := match[1]
								num, _ := strconv.Atoi(number)
								if num == index {
									if !strings.Contains(v.Content, "Waiting to start") {
										re := regexp.MustCompile(`\((\d+)%\)`)
										match := re.FindStringSubmatch(v.Content)

										if len(match) == 0 {
											reply = fmt.Sprintf(`%s
 图片放大，用时: 4秒
 %s`, AtText(msg), msg.Content)
											msg.ReplyText(reply)

											InsertMangodb(&v, index, msg)
											doOk = true
											break
										}

									}
								}

							} else {
								fmt.Println("No match found")
							}

						}

					}

				}
			} else {
				break
			}

		}

		time.Sleep(5 * time.Second)
		i++
	}
	//第二部查看有没有完成
}

func DiscordVariation(MsgID string, index int, msg *openwechat.Message) {

	reply := ""
	i := 0
	//第一步找
	isfound := false
	doOk := false
	for {

		data, err := services.GetMessages()
		if err != nil {
			msg.ReplyText("出问题了哦:" + err.Error())
			return
		}
		if !isfound {
			if i > 60 && !isfound {
				msg.ReplyText("很遗憾，未找到")
				break
			}

			for _, v := range data {

				//查找文本，根据第一级父节点
				if v.ReferencedMessage != nil && v.ReferencedMessage.ID == MsgID {

					if strings.Contains(v.Content, "Variations") {

						reply = fmt.Sprintf(`%s
 
✅ 您的任务已提交
 %s
正在快速处理中，请稍后`, AtText(msg), msg.Content)
						msg.ReplyText(reply)
						isfound = true
						break
					}

				}

			}
		} else {
			if !doOk {
				for _, v := range data {
					//查找文本，根据第一级父节点
					if v.ReferencedMessage != nil && v.ReferencedMessage.ID == MsgID {
						if strings.Contains(v.Content, "Variations") {

							if !strings.Contains(v.Content, "Waiting to start") {
								re := regexp.MustCompile(`\((\d+)%\)`)
								match := re.FindStringSubmatch(v.Content)

								if len(match) == 0 {
									reply = fmt.Sprintf(`%s
 图片放大，用时: 4秒
 %s`, AtText(msg), msg.Content)
									msg.ReplyText(reply)

									InsertMangodb(&v, index, msg)
									doOk = true
									break
								}

							}

						}

					}

				}
			} else {
				break
			}

		}

		time.Sleep(5 * time.Second)
		i++
	}
	//第二部查看有没有完成
}

func SendImage(url string, msg *openwechat.Message) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("发生了异常:", r)
			SendImage(url, msg)
		}
	}()
	// 创建代理 URL
	client := &http.Client{}

	if DEBUG := os.Getenv("DEBUG"); DEBUG == "true" {
		proxyUrl := &url2.URL{}
		// 创建代理 URL
		proxyUrl, _ = url2.Parse("http://127.0.0.1:7890")

		// 创建 Transport 对象
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}

		// 创建 Client 对象
		client = &http.Client{
			Transport: transport,
		}

	}

	resp, err := client.Get(url)
	if err != nil {

		panic(err)
	}
	defer resp.Body.Close()

	fileName := path.Base(url)
	fileName = strings.Split(fileName, "?")[0]

	//filename := uuid.New().String() + ".png"
	file, err := os.Create("./temp/" + fileName)
	if err != nil {
		fmt.Println("Error while creating file:", err)
		return
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("Error while copying image:", err)
		return
	}

	defer func() {
		file.Close()
		//err = os.Remove("image.jpg")
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}
	}()
	//file.Seek(0, 0)
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		log.Fatal(err)
	}
	_, err = msg.ReplyImage(file)
	if err != nil {
		return
	}
	/*	_, err = friends[2].SendImage(file)
		if err != nil {
			fmt.Println("send image:", err)
			return
		}*/
}

func FindMango(msgid string) (MdjMsg2 ai.MdjMsg2) {
	msgs := ai.MdjMsg2{}
	err := global.GVA_MONGO.FindOne("mdj", "mdj_msg3", bson.D{bson.E{"_id", msgid}}, &msgs)
	if err != nil {
		panic(err)
	}
	return msgs
}
func InsertMangodb(m *discordgo.Message, index int, msg *openwechat.Message) error {
	ImgHash := ""
	if len(m.Attachments) > 0 && m.Attachments[0].URL != "" {
		ImgHash = getImgHash(m.Attachments[0].URL, 0)

		SendImage(m.Attachments[0].URL, msg)
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
		uid, hashId, _ := getUIDAndHash(m.Content)
		MdjMsg := &ai.MdjMsg2{
			ID:          m.ID,
			MsgId:       m.ID,
			ChannelId:   m.ChannelID,
			GuildId:     m.GuildID,
			AuthorId:    m.Author.ID,
			ParentID:    ParentID,
			ChangeType:  ChangeType,
			Index:       index,
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

		if err := global.GVA_MONGO.InsertMany("mdj", "mdj_msg3", MdjMsg); err != nil {
			fmt.Println("插入失败", err.Error())
			return err
		}
	}
	return nil
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
