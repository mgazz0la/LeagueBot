package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/mgazz0la/leaguebot/internal/platform/sleeper"
)

type (
	Guild struct {
		GuildName             string            `json:"guild_name"`
		GuildID               string            `json:"guild_id"`
		LeagueName            string            `json:"league_name"`
		SleeperLeagueID       sleeper.LeagueID  `json:"sleeper_league_id"`
		NotificationChannelID string            `json:"notification_channel_id"`
		SquadOwners           map[string]string `json:"squad_owners"`
	}

	Config struct {
		Token  string  `json:"token"`
		Guilds []Guild `json:"guilds"`
	}
)

const CONFIG_FILE = "config.json"

func LoadConfig() Config {
	bytes, err := ioutil.ReadFile(CONFIG_FILE)
	if err != nil {
		panic(err)
	}

	var c Config
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		panic(err)
	}

	return c
}
