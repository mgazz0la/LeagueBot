package league

import (
	"errors"
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/mgazz0la/leaguebot/internal/domain"
	"github.com/mgazz0la/leaguebot/internal/platform"
)

type (
	LeagueState struct {
		platform platform.Platform

		txmu         sync.Mutex
		transactions map[domain.TransactionID]domain.Transaction
	}
)

func NewLeague(platform platform.Platform) *LeagueState {
	return &LeagueState{
		platform: platform,
	}
}

func (ls *LeagueState) TransactionToDiscordMessage(
	txn domain.Transaction,
) (*discordgo.MessageSend, error) {
	m := new(discordgo.MessageSend)
	switch txn.Type() {
	case domain.TransactionTypeFreeAgent:
		fa, ok := txn.(domain.FreeAgentTransaction)
		if !ok {
			return nil, errors.New("fa txn fail")
		}

		sq, ok := ls.GetSquadByID(fa.SquadID)
		if !ok {
			return nil, errors.New("squad fail")
		}

		var thumbnailID domain.PlayerID
		var embedFields []*discordgo.MessageEmbedField
		if fa.Add != nil {
			add, ok := ls.GetPlayerByID(*fa.Add)
			if !ok {
				return nil, errors.New("player fail")
			}
			embedFields = append(embedFields, &discordgo.MessageEmbedField{
				Name:   "Add",
				Value:  add.String(),
				Inline: true,
			})
			thumbnailID = *fa.Add
		}
		if fa.Drop != nil {
			drop, ok := ls.GetPlayerByID(*fa.Drop)
			if !ok {
				return nil, errors.New("player fail")
			}
			embedFields = append(embedFields, &discordgo.MessageEmbedField{
				Name:   "Drop",
				Value:  drop.String(),
				Inline: true,
			})
			if thumbnailID == "" {
				thumbnailID = *fa.Drop
			}
		}

		m.Embeds = []*discordgo.MessageEmbed{
			{
				Title: "Roster Move",
				Author: &discordgo.MessageEmbedAuthor{
					Name:    sq.Name,
					IconURL: fmt.Sprintf("https://sleepercdn.com/avatars/thumbs/%s", sq.AvatarID),
				},
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: fmt.Sprintf("https://sleepercdn.com/content/nfl/players/%s.jpg", thumbnailID),
				},
				Fields: embedFields,
			},
		}
	}
	return m, nil
}

func (ls *LeagueState) GetPlayerByID(pid domain.PlayerID) (domain.Player, bool) {
	pmap, err := ls.platform.GetPlayers()
	if err != nil {
		return domain.Player{}, false
	}

	p, ok := pmap[pid]
	if ok {
		return *p, true
	}
	return domain.Player{}, false
}

func (ls *LeagueState) GetSquadByID(sqid domain.SquadID) (domain.Squad, bool) {
	sqs, err := ls.GetSquads()
	if err != nil {
		return domain.Squad{}, false
	}

	sq, ok := sqs[sqid]
	if ok {
		return *sq, true
	}
	return domain.Squad{}, false
}

func (ls *LeagueState) GetTransactions() map[domain.TransactionID]domain.Transaction {
	ls.txmu.Lock()
	defer ls.txmu.Unlock()
	ls.transactions, _ = ls.platform.GetTransactions(GetCurrentWeek())
	return ls.transactions
}

func (ls *LeagueState) GetTransactionByID(txid domain.TransactionID) (domain.Transaction, bool) {
	ls.txmu.Lock()
	defer ls.txmu.Unlock()
	tx, ok := ls.transactions[txid]
	return tx, ok
}

func (ls *LeagueState) GetSquads() (map[domain.SquadID]*domain.Squad, error) {
	return ls.platform.GetSquads()
}
