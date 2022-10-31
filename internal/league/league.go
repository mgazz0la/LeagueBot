package league

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/mgazz0la/leaguebot/internal/domain"
	"github.com/mgazz0la/leaguebot/internal/platform"
)

type (
	LeagueState struct {
		platform platform.Platform

		currentWeek   uint
		lastWeekFetch time.Time

		squadOwners map[domain.SquadID]string

		txmu         sync.Mutex
		transactions map[domain.TransactionID]domain.Transaction
	}
)

func NewLeague(platform platform.Platform, squadOwners map[domain.SquadID]string) *LeagueState {
	return &LeagueState{
		platform:    platform,
		squadOwners: squadOwners,
	}
}

func (ls *LeagueState) TransactionToDiscordMessage(
	txn domain.Transaction,
) (*discordgo.MessageSend, error) {
	m := new(discordgo.MessageSend)
	switch txn.Type() {
	case domain.TransactionTypeFreeAgent:
		fa, ok := txn.(*domain.FreeAgentTransaction)
		if !ok {
			return nil, errors.New("fa txn fail")
		}

		sq, ok := ls.GetSquadByID(fa.SquadID)
		if !ok {
			return nil, errors.New("squad fail")
		}

		s := ls.squadOwners[sq.SquadID]
		if fa.Add != nil {
			add, ok := ls.GetPlayerByID(*fa.Add)
			if !ok {
				return nil, errors.New("player fail")
			}
			s += " added " + add.String()
		}
		if fa.Add != nil && fa.Drop != nil {
			s += " and"
		}
		if fa.Drop != nil {
			drop, ok := ls.GetPlayerByID(*fa.Drop)
			if !ok {
				return nil, errors.New("player fail")
			}
			s += " dropped " + drop.String()
		}

		embed := &discordgo.MessageEmbed{
			Title:       "Roster Move",
			Description: s,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    sq.Name,
				IconURL: fmt.Sprintf("https://sleepercdn.com/avatars/thumbs/%s", sq.AvatarID),
			},
		}
		if fa.Add != nil {
			embed.Image = &discordgo.MessageEmbedImage{
				URL: fmt.Sprintf("https://sleepercdn.com/content/nfl/players/%s.jpg", *fa.Add),
			}
		}
		if fa.Drop != nil {
			embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
				URL: fmt.Sprintf("https://sleepercdn.com/content/nfl/players/%s.jpg", *fa.Drop),
			}
		}
		m.Embeds = []*discordgo.MessageEmbed{embed}
	case domain.TransactionTypeWaiver:
		w, ok := txn.(*domain.WaiverTransaction)
		if !ok {
			return nil, errors.New("waiver txn fail")
		}

		if !w.DidWin {
			return nil, errors.New("not supporting failed bids yet")
		}

		sq, ok := ls.GetSquadByID(w.SquadID)
		if !ok {
			return nil, errors.New("squad fail")
		}

		s := ls.squadOwners[sq.SquadID]

		add, ok := ls.GetPlayerByID(w.Add)
		if !ok {
			return nil, errors.New("player fail")
		}
		s += " added " + add.String() + " for $" + fmt.Sprint(w.Bid)

		if w.Drop != nil {
			drop, ok := ls.GetPlayerByID(*w.Drop)
			if !ok {
				return nil, errors.New("player fail")
			}
			s += " and dropped " + drop.String()
		}

		embed := &discordgo.MessageEmbed{
			Title:       "Waiver Move",
			Description: s,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    sq.Name,
				IconURL: fmt.Sprintf("https://sleepercdn.com/avatars/thumbs/%s", sq.AvatarID),
			},
			Image: &discordgo.MessageEmbedImage{
				URL: fmt.Sprintf("https://sleepercdn.com/content/nfl/players/%s.jpg", w.Add),
			},
		}
		if w.Drop != nil {
			embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
				URL: fmt.Sprintf("https://sleepercdn.com/content/nfl/players/%s.jpg", *w.Drop),
			}
		}
		m.Embeds = []*discordgo.MessageEmbed{embed}

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
	var err error
	ls.transactions, err = ls.platform.GetTransactions(ls.GetCurrentWeek())
	if err != nil {
		log.Println(err.Error())
	}
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

func (ls *LeagueState) GetCurrentWeek() uint {
	if time.Now().Sub(ls.lastWeekFetch) > 4*time.Hour {
		if week, err := ls.platform.GetCurrentWeek(); err == nil {
			ls.currentWeek = week
			ls.lastWeekFetch = time.Now()
		}
	}
	return ls.currentWeek
}
