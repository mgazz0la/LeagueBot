package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/mgazz0la/leaguebot/internal/config"
	"github.com/mgazz0la/leaguebot/internal/league"
)

type (
	ChannelID string
	GuildID   string
	BotState  struct {
		GuildConfig config.Guild
		Session     *discordgo.Session
		League      *league.LeagueState
	}
)
