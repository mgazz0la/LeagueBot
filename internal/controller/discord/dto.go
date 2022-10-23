package discord

import (
	"github.com/mgazz0la/leaguebot/internal/league"
	"github.com/mgazz0la/leaguebot/internal/platform"
)

type (
	ChannelID string
	GuildID   string
	BotState  struct {
		GuildID  GuildID
		Platform platform.Platform
		League   *league.LeagueState
	}
)
