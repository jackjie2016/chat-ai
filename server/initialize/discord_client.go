package initialize

import (
	"fmt"
	discord "github.com/bwmarrin/discordgo"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

var discordClient *discord.Session

func LoadDiscordClient(create func(s *discord.Session, m *discord.MessageCreate), update func(s *discord.Session, m *discord.MessageUpdate)) {
	var err error
	discordClient, err = discord.New("Bot " + global.GVA_CONFIG.Discord.DISCORD_BOT_TOKEN)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	err = discordClient.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	discordClient.AddHandler(create)
	discordClient.AddHandler(update)
	//go func() {
	//	for {
	//		_, _ = utils.GetMessages()
	//		time.Sleep(30 * time.Second)
	//	}
	//}()
}

func GetDiscordClient() *discord.Session {
	return discordClient
}
