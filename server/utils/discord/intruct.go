package utils

import (
	"github.com/bwmarrin/discordgo"
	services "github.com/flipped-aurora/gin-vue-admin/server/logic"
)

func GenerateImage(prompt string) error {
	err := services.GenerateImage(prompt)
	return err
}

func ImageUpscale(index int64, discordMsgId string, msgHash string) error {
	err := services.Upscale(index, discordMsgId, msgHash)
	return err
}

func ImageVariation(index int64, discordMsgId string, msgHash string) error {
	err := services.Variate(index, discordMsgId, msgHash)
	return err
}

func GetMessages() ([]discordgo.Message, error) {
	data, err := services.GetMessages()
	for _, v := range data {
		d := v
		//fmt.Printf("%p,%+v\n", d, d.ID)
		go InsertMongoDB(&d)
	}
	return data, err
}
func ImageMaxUpscale(discordMsgId string, msgHash string) error {
	err := services.MaxUpscale(discordMsgId, msgHash)
	return err
}
func ImageReset(discordMsgId string, msgHash string) error {
	err := services.Reset(discordMsgId, msgHash)
	return err
}

func ImageDescribe(uploadName string) error {
	err := services.Describe(uploadName)
	return err
}
func ImageBlend(uploadNames string) error {
	err := services.Blend(uploadNames)
	return err
}
