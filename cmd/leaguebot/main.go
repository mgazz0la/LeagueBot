package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/mgazz0la/leaguebot/internal/config"
	"github.com/mgazz0la/leaguebot/internal/controller/discord"
	"github.com/mgazz0la/leaguebot/internal/league"
	"github.com/mgazz0la/leaguebot/internal/platform/sleeper"
	"github.com/mgazz0la/leaguebot/internal/watcher"
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
		gid := cfg.Guilds[i].GuildID
		botStates[gid] = &discord.BotState{
			Platform: sleeper.NewSleeper(cfg.Guilds[i].SleeperLeagueID),
			League:   new(league.LeagueState),
		}
		if err = botStates[gid].League.Load(botStates[gid].Platform); err != nil {
			log.Fatalf("Failed to load league state: %v", err)
		}
		go watcher.TransactionWatcher(botStates[cfg.Guilds[i].GuildID])
	}

	discord.RegisterHandlers(d, botStates)

	log.Println("ready")

	defer d.Close()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
