package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/mgazz0la/leaguebot/internal/config"
	"github.com/mgazz0la/leaguebot/internal/controller/discord"
	"github.com/mgazz0la/leaguebot/internal/domain"
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
		gid := discord.GuildID(cfg.Guilds[i].GuildID)

		sleeper, err := sleeper.NewSleeper(cfg.Guilds[i].SleeperLeagueID)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		sqOwners := make(map[domain.SquadID]string)
		for k, v := range cfg.Guilds[i].SquadOwners {
			sqOwners[domain.SquadID(k)] = v
		}
		botStates[gid] = &discord.BotState{
			Session:     d,
			GuildConfig: cfg.Guilds[i],
			League:      league.NewLeague(sleeper, sqOwners),
		}

		go watcher.NewTransactionWatcher(botStates[gid]).Run()
	}

	discord.RegisterHandlers(d, botStates)

	log.Println("ready")
	/*
		gid := discord.GuildID(cfg.Guilds[0].GuildID)
		add := domain.PlayerID("4666")
		drop := domain.PlayerID("1920")
		m, err := botStates[gid].League.TransactionToDiscordMessage(domain.WaiverTransaction{
			SquadID: domain.SquadID("1"),
			DidWin:  true,
			Bid:     99,
			Add:     add,
			Drop:    &drop,
		})
		if err != nil {
			log.Println(err.Error())
		} else {
			botStates[gid].Session.ChannelMessageSendComplex(botStates[gid].GuildConfig.NotificationChannelID, m)
		}
	*/

	defer d.Close()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
