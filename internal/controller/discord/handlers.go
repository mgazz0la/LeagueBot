package discord

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lithammer/fuzzysearch/fuzzy"

	"github.com/bwmarrin/discordgo"
	"github.com/mgazz0la/leaguebot/internal/domain"
	"github.com/mgazz0la/leaguebot/internal/league"
	"github.com/mgazz0la/leaguebot/internal/utils"
)

type handler func(s *discordgo.Session, i *discordgo.InteractionCreate, bs *BotState)

var handlers = map[string]handler{
	"playoffs": playoffHandler,
	"roster":   rosterHandler,
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

func completedRosterHandler(
	s *discordgo.Session, i *discordgo.InteractionCreate, bs *BotState,
) error {
	sqid := domain.SquadID(i.ApplicationCommandData().Options[0].StringValue())
	sq, ok := bs.League.GetSquadByID(sqid)
	if !ok {
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("who the heck is %s???", sqid),
			},
		})
	}

	makePlayerNameList := func(pid domain.PlayerID) string {
		p, ok := bs.League.GetPlayerByID(pid)
		if !ok {
			log.Printf("could not find player [%v]", pid)
			p.FirstName = "D'Pez"
			p.LastName = "Poopsie"
		}

		return p.FirstName + " " + p.LastName
	}

	embedFields := []*discordgo.MessageEmbedField{
		{
			Name:  "Starters",
			Value: strings.Join(utils.Map(makePlayerNameList, sq.Starters), "\n"),
		},
		{
			Name:  "Bench",
			Value: strings.Join(utils.Map(makePlayerNameList, sq.Bench), "\n"),
		},
	}
	if len(sq.IR) > 0 {
		embedFields = append(embedFields, &discordgo.MessageEmbedField{
			Name:  "IR",
			Value: strings.Join(utils.Map(makePlayerNameList, sq.IR), "\n"),
		})
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Type:   discordgo.EmbedTypeRich,
					Title:  sq.Name,
					Fields: embedFields,
				},
			},
		},
	})
}

func autocompleteRosterHandler(
	s *discordgo.Session, i *discordgo.InteractionCreate, bs *BotState,
) error {
	txt := i.ApplicationCommandData().Options[0].StringValue()
	sqs, err := bs.League.GetSquads()
	if err != nil {
		log.Println(err.Error())
		return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "could not get roster list--go complain to Commish",
			},
		})
	}

	choices := utils.Map(func(s *domain.Squad) *discordgo.ApplicationCommandOptionChoice {
		return &discordgo.ApplicationCommandOptionChoice{
			Name:  s.Name,
			Value: s.SquadID,
		}
	}, utils.Values(sqs))

	if txt != "" {
		sort.Slice(choices, func(i, j int) bool {
			return fuzzy.RankMatchFold(txt, choices[i].Name) >
				fuzzy.RankMatchFold(txt, choices[j].Name)
		})
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}

func rosterHandler(s *discordgo.Session, i *discordgo.InteractionCreate, bs *BotState) {
	var err error

	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		err = completedRosterHandler(s, i, bs)
	case discordgo.InteractionApplicationCommandAutocomplete:
		err = autocompleteRosterHandler(s, i, bs)
	}

	if err != nil {
		log.Println(err.Error())
	}
}

func playoffHandler(s *discordgo.Session, i *discordgo.InteractionCreate, bs *BotState) {
	sqmap, err := bs.League.GetSquads()
	if err != nil {
		log.Println(err.Error())
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "could not get roster list--go complain to Commish",
			},
		})
		return
	}

	sqs := utils.Values(sqmap)
	league.ApplySeeds(sqs)

	tBAMF := table.NewWriter()
	tBAMF.SetStyle(table.StyleRounded)
	tBAMF.Style().Options.SeparateColumns = false
	tBAMF.Style().Options.DrawBorder = false
	tBAMF.AppendRows(
		utils.Map(func(s *domain.Squad) table.Row {
			winLoss := fmt.Sprintf("%d-%d", s.Wins, s.Losses)
			seed := fmt.Sprint(s.Seed)
			switch s.Seed {
			case 1, 2:
				seed += "*"
			case 5, 6:
				seed += "†"
			}
			return table.Row{seed, s.Name, winLoss, fmt.Sprintf("%.2f", s.PointsFor)}
		}, sqs[:6]),
	)

	tSacko := table.NewWriter()
	tSacko.SetStyle(table.StyleRounded)
	tSacko.Style().Options.SeparateColumns = false
	tSacko.Style().Options.DrawBorder = false
	tSacko.AppendRows(
		utils.Map(func(s *domain.Squad) table.Row {
			winLoss := fmt.Sprintf("%d-%d", s.Wins, s.Losses)
			seed := fmt.Sprint(s.Seed)
			if s.Seed == 11 || s.Seed == 12 {
				seed += "*"
			}
			return table.Row{seed, s.Name, winLoss, fmt.Sprintf("%.2f", s.PointsFor)}
		}, sqs[6:]),
	)

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
