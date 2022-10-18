package discord

import "github.com/mgazz0la/leaguebot/internal/platform"

type (
	GuildID  string
	BotState struct {
		GuildID  GuildID
		Platform platform.Platform
	}
)
