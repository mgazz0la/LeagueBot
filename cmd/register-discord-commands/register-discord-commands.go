package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/mgazz0la/leaguebot/internal/config"
	"github.com/mgazz0la/leaguebot/internal/controller/discord"
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
	defer d.Close()

	for i := range cfg.Guilds {
		for _, c := range discord.GetCommandList() {
			_, err := d.ApplicationCommandCreate(d.State.User.ID, string(cfg.Guilds[i].GuildID), c)
			if err != nil {
				log.Panicf("Cannot create '%v' command: %v", c.Name, err)
			}
			log.Printf("Created command [%s] for guild [%s]", c.Name, cfg.Guilds[i].GuildName)
		}
	}
}
