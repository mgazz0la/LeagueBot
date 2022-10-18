package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/mgazz0la/leaguebot/internal/controller/discord"
	"github.com/mgazz0la/leaguebot/internal/platform/sleeper"
)

type Config struct {
	Token  string `json:"token"`
	Guilds []struct {
		GuildName       string           `json:"guild_name"`
		GuildID         discord.GuildID  `json:"guild_id"`
		LeagueName      string           `json:"league_name"`
		SleeperLeagueID sleeper.LeagueID `json:"sleeper_league_id"`
	} `json:"guilds"`
}

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
