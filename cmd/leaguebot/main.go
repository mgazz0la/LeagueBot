package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/mgazz0la/leaguebot/internal/config"
	"github.com/mgazz0la/leaguebot/internal/controller/discord"
	"github.com/mgazz0la/leaguebot/internal/platform/sleeper"
)

func main() {
	cfg := config.LoadConfig()
	d, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	err = d.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	botStates := make(map[discord.GuildID]*discord.BotState)
	for i := range cfg.Guilds {
		botStates[cfg.Guilds[i].GuildID] = &discord.BotState{
			Platform: sleeper.NewSleeper(cfg.Guilds[i].SleeperLeagueID),
		}
	}

	discord.RegisterHandlers(d, botStates)

	log.Println("ready")

	defer d.Close()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
