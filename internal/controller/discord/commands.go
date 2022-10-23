package discord

import "github.com/bwmarrin/discordgo"

func GetCommandList() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "playoffs",
			Description: "Shows the current playoff picture",
			Type:        discordgo.ChatApplicationCommand,
		},
		{
			Name:        "roster",
			Description: "Shows the roster for the specified team",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "team",
					Description:  "whom?",
					Type:         discordgo.ApplicationCommandOptionString,
					Required:     true,
					Autocomplete: true,
				},
			},
		},
	}
}
