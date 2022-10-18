package discord

import "github.com/bwmarrin/discordgo"

func GetCommandList() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "playoffs",
			Description: "Shows the current playoff picture",
		},
	}
}
