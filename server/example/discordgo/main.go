package main

import "github.com/bwmarrin/discordgo"

func main() {
	discord, err := discordgo.New("Bot " + "authentication token")
	if err != nil {
		panic(err)

	}
	discord.Close()
}
