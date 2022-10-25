package discord

import (
	"github.com/mgazz0la/leaguebot/internal/league"
)

type (
	ChannelID string
	GuildID   string
	BotState  struct {
		GuildID GuildID
		League  *league.LeagueState
	}
)
