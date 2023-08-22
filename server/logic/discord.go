package logic

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	aiModel "github.com/flipped-aurora/gin-vue-admin/server/model/ai"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	url2 "net/url"
)

const (
	url        = "https://discord.com/api/v9/interactions"
	uploadUrl  = "https://discord.com/api/v9/channels/%s/attachments"
	messageUrl = "https://discord.com/api/v9/channels/%s/messages?limit=%d"
)

func GetNonce() string {
	max := int64(math.Pow10(10)) - 1
	randNum := rand.Int63n(max)
	randNum2 := rand.Int63n(int64(math.Pow10(8)))

	return fmt.Sprintf("%d%010d%08d", rand.Intn(8)+1, randNum, randNum2)

	//randNum := rand.Int63n(max)
	//return fmt.Sprintf("%019d", randNum)

}
func GenerateImage(prompt string) error {
	requestBody := aiModel.ReqTriggerDiscord{
		Type:          2,
		GuildID:       global.GVA_CONFIG.Discord.DISCORD_SERVER_ID,
		ChannelID:     global.GVA_CONFIG.Discord.DISCORD_CHANNEL_ID,
		ApplicationId: "936929561302675456",
		SessionId:     "cb06f61453064c0983f2adae2a88c223",
		Nonce:         GetNonce(),
		Data: aiModel.DSCommand{
			Version: global.GVA_CONFIG.Discord.VERSION,
			Id:      "938956540159881230",
			Name:    "imagine",
			Type:    1,
			Options: []aiModel.DSOption{{Type: 3, Name: "prompt", Value: prompt}},
			ApplicationCommand: aiModel.DSApplicationCommand{
				Id:                       "938956540159881230",
				ApplicationId:            "936929561302675456",
				Version:                  global.GVA_CONFIG.Discord.VERSION,
				DefaultPermission:        true,
				DefaultMemberPermissions: nil,
				Type:                     1,
				Nsfw:                     false,
				Name:                     "imagine",
				Description:              "Lucky you!",
				DmPermission:             true,
				Options:                  []aiModel.DSCommandOption{{Type: 3, Name: "prompt", Description: "The prompt to imagine", Required: true}},
			},
			Attachments: []aiModel.ReqCommandAttachments{},
		},
	}
	_, err := request(requestBody, url)
	return err
}

func Upscale(index int64, messageId string, messageHash string) error {
	requestBody := aiModel.ReqUpscaleDiscord{
		Type:          3,
		GuildId:       global.GVA_CONFIG.Discord.DISCORD_SERVER_ID,
		ChannelId:     global.GVA_CONFIG.Discord.DISCORD_CHANNEL_ID,
		MessageFlags:  0,
		MessageId:     messageId,
		ApplicationId: "936929561302675456",
		SessionId:     "cda4348d89fdb44f0f51b67801ca9a78",
		Nonce:         GetNonce(),
		Data: aiModel.UpscaleData{
			ComponentType: 2,
			CustomId:      fmt.Sprintf("MJ::JOB::upsample::%d::%s", index, messageHash),
		},
	}
	_, err := request(requestBody, url)
	return err
}

func Blend(prompt string) error {
	attachments := make([]aiModel.ReqCommandAttachments, 0)

	files := strings.Split(prompt, ";")
	if len(files) < 2 {
		return errors.New("垫图至少需要两张图片")
	}
	for index, file := range files {
		filename := path.Base(file)
		attachments = append(attachments, aiModel.ReqCommandAttachments{
			Id:             fmt.Sprintf("%d", index),
			Filename:       filename,
			UploadFilename: file,
		})
	}

	Options := make([]aiModel.DSOption, len(attachments))

	for i, _ := range attachments {
		Options[i] = aiModel.DSOption{
			Type:  11,
			Name:  fmt.Sprintf("image%d", i+1),
			Value: i,
		}
	}

	requestBody := aiModel.ReqBlendDiscord{
		Type:          2,
		GuildID:       global.GVA_CONFIG.Discord.DISCORD_SERVER_ID,
		ChannelID:     global.GVA_CONFIG.Discord.DISCORD_CHANNEL_ID,
		ApplicationID: "936929561302675456",
		SessionID:     "cb06f61453064c0983f2adae2a88c223",
		Nonce:         GetNonce(),
		Data: aiModel.BlendData{
			Version: "1118961510123847773",
			ID:      "1062880104792997970",
			Name:    "blend",
			Type:    1,
			Options: Options,
			ApplicationCommand: aiModel.ApplicationCommand{
				ID:                       "1062880104792997970",
				ApplicationID:            "936929561302675456",
				Version:                  "1118961510123847773",
				DmPermission:             true,
				DefaultMemberPermissions: nil,
				Type:                     1,
				Nsfw:                     false,
				Name:                     "blend",
				Description:              "Blend images together seamlessly!",
				Contexts:                 []int{0, 1, 2},
				Options: []aiModel.ApplicationOptions{
					{Type: 11, Name: "image1", Description: "First image to add to the blend", Required: true},
					{Type: 11, Name: "image2", Description: "Second image to add to the blend", Required: true},
					{Type: 3, Name: "dimensions", Description: "The dimensions of the image. If not specified, the image will be square.", Choices: []aiModel.Choices{
						{Name: "Portrait", Value: "--ar 2:3"},
						{Name: "Square", Value: "--ar 1:1"},
						{Name: "Landscape", Value: "--ar 3:2"},
					}},
					{Type: 11, Name: "image3", Description: "Third image to add to the blend (optional)"},
					{Type: 11, Name: "image4", Description: "Fourth image to add to the blend (optional)"},
					{Type: 11, Name: "image5", Description: "Fifth image to add to the blend (optional)"},
				},
			},
			Attachments: attachments,
		},
	}
	_, err := request(requestBody, url)
	return err
}

func MaxUpscale(messageId string, messageHash string) error {
	requestBody := aiModel.ReqUpscaleDiscord{
		Type:          3,
		GuildId:       global.GVA_CONFIG.Discord.DISCORD_SERVER_ID,
		ChannelId:     global.GVA_CONFIG.Discord.DISCORD_CHANNEL_ID,
		MessageFlags:  0,
		MessageId:     messageId,
		ApplicationId: "936929561302675456",
		SessionId:     "cda4348d89fdb44f0f51b67801ca9a78",
		Nonce:         GetNonce(),
		Data: aiModel.UpscaleData{
			ComponentType: 2,
			CustomId:      fmt.Sprintf("MJ::JOB::variation::1::%s::SOLO", messageHash),
		},
	}

	data, _ := json.Marshal(requestBody)

	fmt.Println("max upscale request body: ", string(data))

	_, err := request(requestBody, url)
	return err
}

func Variate(index int64, messageId string, messageHash string) error {
	requestBody := aiModel.ReqVariationDiscord{
		Type:          3,
		GuildId:       global.GVA_CONFIG.Discord.DISCORD_SERVER_ID,
		ChannelId:     global.GVA_CONFIG.Discord.DISCORD_CHANNEL_ID,
		MessageFlags:  0,
		MessageId:     messageId,
		ApplicationId: "936929561302675456",
		SessionId:     "45bc04dd4da37141a5f73dfbfaf5bdcf",
		Nonce:         GetNonce(),
		Data: aiModel.UpscaleData{
			ComponentType: 2,
			CustomId:      fmt.Sprintf("MJ::JOB::variation::%d::%s", index, messageHash),
		},
	}
	_, err := request(requestBody, url)
	return err
}

func Reset(messageId string, messageHash string) error {
	requestBody := aiModel.ReqResetDiscord{
		Type:          3,
		GuildId:       global.GVA_CONFIG.Discord.DISCORD_SERVER_ID,
		ChannelId:     global.GVA_CONFIG.Discord.DISCORD_CHANNEL_ID,
		MessageFlags:  0,
		MessageId:     messageId,
		ApplicationId: "936929561302675456",
		SessionId:     "45bc04dd4da37141a5f73dfbfaf5bdcf",
		Nonce:         GetNonce(),
		Data: aiModel.UpscaleData{
			ComponentType: 2,
			CustomId:      fmt.Sprintf("MJ::JOB::reroll::0::%s::SOLO", messageHash),
		},
	}
	_, err := request(requestBody, url)
	return err
}

func Describe(uploadName string) error {
	requestBody := aiModel.ReqTriggerDiscord{
		Type:          2,
		GuildID:       global.GVA_CONFIG.Discord.DISCORD_SERVER_ID,
		ChannelID:     global.GVA_CONFIG.Discord.DISCORD_CHANNEL_ID,
		ApplicationId: "936929561302675456",
		SessionId:     "0033db636f7ce1a951e54cdac7044de3",
		Nonce:         GetNonce(),
		Data: aiModel.DSCommand{
			Version: "1118961510123847774",
			Id:      "1092492867185950852",
			Name:    "describe",
			Type:    1,
			Options: []aiModel.DSOption{{Type: 11, Name: "image", Value: 0}},
			ApplicationCommand: aiModel.DSApplicationCommand{
				Id:                       "1092492867185950852",
				ApplicationId:            "936929561302675456",
				Version:                  "1118961510123847774",
				DefaultPermission:        true,
				DefaultMemberPermissions: nil,
				Type:                     1,
				Nsfw:                     false,
				Name:                     "describe",
				Description:              "Writes a prompt based on your image.",
				DmPermission:             true,
				Options:                  []aiModel.DSCommandOption{{Type: 11, Name: "image", Description: "The image to describe", Required: true}},
			},
			Attachments: []aiModel.ReqCommandAttachments{{
				Id:             "0",
				Filename:       filepath.Base(uploadName),
				UploadFilename: uploadName,
			}},
		},
	}
	_, err := request(requestBody, url)
	return err
}

func GetMessages() ([]discordgo.Message, error) {
	replacedUrl := fmt.Sprintf(messageUrl, global.GVA_CONFIG.Discord.DISCORD_CHANNEL_ID, 10)
	body, err := GET(replacedUrl)
	var data []discordgo.Message
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, err
}

func Attachments(name string, size int64) (aiModel.ResAttachments, error) {
	requestBody := aiModel.ReqAttachments{
		Files: []aiModel.ReqFile{{
			Filename: name,
			FileSize: size,
			Id:       "1",
		}},
	}
	replacedUrl := fmt.Sprintf(uploadUrl, global.GVA_CONFIG.Discord.DISCORD_CHANNEL_ID)
	body, err := request(requestBody, replacedUrl)
	var data aiModel.ResAttachments
	json.Unmarshal(body, &data)
	return data, err
}

func request(params interface{}, url string) ([]byte, error) {
	requestData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(requestData))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", global.GVA_CONFIG.Discord.DISCORD_USER_TOKEN)
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

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	bod, respErr := ioutil.ReadAll(response.Body)
	fmt.Println("response: ", string(bod), respErr, response.Status)
	return bod, respErr
}
func GET(url string) ([]byte, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", global.GVA_CONFIG.Discord.DISCORD_USER_TOKEN)
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
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	bod, respErr := ioutil.ReadAll(response.Body)
	//fmt.Println("response: ", string(bod), respErr, response.Status)
	return bod, respErr
}
