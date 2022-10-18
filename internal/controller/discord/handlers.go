package discord

import (
	"fmt"
	"log"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/bwmarrin/discordgo"
	"github.com/mgazz0la/leaguebot/internal/domain"
	"github.com/mgazz0la/leaguebot/internal/utils"
)

type handler func(s *discordgo.Session, i *discordgo.InteractionCreate, bs *BotState)

var handlers = map[string]handler{
	"playoffs": playoffHandler,
}

func RegisterHandlers(d *discordgo.Session, botStates map[GuildID]*BotState) {
	d.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := handlers[i.ApplicationCommandData().Name]; ok {
			if bs, ok := botStates[GuildID(i.GuildID)]; ok {
				h(s, i, bs)
			}
		}
	})
}

func playoffHandler(s *discordgo.Session, i *discordgo.InteractionCreate, bs *BotState) {
	sqs, _ := bs.Platform.GetSquads()
	tBAMF := table.NewWriter()
	tBAMF.SetStyle(table.StyleRounded)
	tBAMF.Style().Options.SeparateColumns = false
	tBAMF.Style().Options.DrawBorder = false
	tBAMF.AppendRows(
		utils.Map(func(s *domain.Squad) table.Row {
			switch s.Seed {
			case 1, 2:
				return table.Row{fmt.Sprintf("%d*", s.Seed), s.Name, fmt.Sprintf("%d-%d", s.Wins, s.Losses), s.PointsFor}
			case 5, 6:
				return table.Row{fmt.Sprintf("%d†", s.Seed), s.Name, fmt.Sprintf("%d-%d", s.Wins, s.Losses), s.PointsFor}
			}
			return table.Row{s.Seed, s.Name, fmt.Sprintf("%d-%d", s.Wins, s.Losses), s.PointsFor}
		}, sqs[:6]),
	)

	tSacko := table.NewWriter()
	tSacko.SetStyle(table.StyleRounded)
	tSacko.Style().Options.SeparateColumns = false
	tSacko.Style().Options.DrawBorder = false
	tSacko.AppendRows(
		utils.Map(func(s *domain.Squad) table.Row {
			switch s.Seed {
			case 11, 12:
				return table.Row{fmt.Sprintf("%d*", s.Seed), s.Name, fmt.Sprintf("%d-%d", s.Wins, s.Losses), s.PointsFor}
			}
			return table.Row{s.Seed, s.Name, fmt.Sprintf("%d-%d", s.Wins, s.Losses), s.PointsFor}
		}, sqs[6:]),
	)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Type:  discordgo.EmbedTypeRich,
					Title: "Playoff Picture",
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "BAMF Bowl",
							Value: "```" + tBAMF.Render() + "```",
						},
						{
							Name:  "Sacko Bowl",
							Value: "```" + tSacko.Render() + "```",
						},
					},
					Footer: &discordgo.MessageEmbedFooter{
						Text: "* first-round bye, † wild card",
					},
				},
			},
		},
	})
	if err != nil {
		log.Println(err.Error())
	}
}
