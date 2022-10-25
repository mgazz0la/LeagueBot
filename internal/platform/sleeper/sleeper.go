package sleeper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/mgazz0la/leaguebot/internal/domain"
	"github.com/mgazz0la/leaguebot/internal/platform"
	"github.com/mgazz0la/leaguebot/internal/utils"
)

type Sleeper struct {
	sc             *sleeperClient
	id             LeagueID
	sqmap          map[domain.SquadID]*domain.Squad
	lastSquadFetch time.Time
	pmap           map[domain.PlayerID]*domain.Player
}

const (
	playersJSON  = "internal/platform/sleeper/players.json"
	playerMapTTL = 24 * time.Hour
	squadMapTTL  = 15 * time.Minute
)

func NewSleeper(id LeagueID) (platform.Platform, error) {
	s := &Sleeper{
		sc: newSleeperClient(),
		id: id,
	}

	_, err := s.fetchPlayerFileIfNeeded()
	if err != nil {
		return nil, err
	}

	if err = s.loadPlayerMapFromFile(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Sleeper) GetTransactions(week uint) (map[domain.TransactionID]domain.Transaction, error) {
	txns, err := s.sc.GetTransactions(s.id, week)
	if err != nil {
		return nil, err
	}

	ts := make(map[domain.TransactionID]domain.Transaction)

	for _, txn := range txns {
		switch txn.Type {
		case transactionTypeFreeAgent:
			fa, err := s.handleFreeAgentTransaction(txn)
			if err != nil {
				log.Println(err)
				continue
			}
			ts[fa.ID()] = fa

		case transactionTypeWaiver:
			w, err := s.handleWaiverTransaction(txn)
			if err != nil {
				log.Println(err)
				continue
			}
			ts[w.ID()] = w

		case transactionTypeTrade:
			t, err := s.handleWaiverTransaction(txn)
			if err != nil {
				log.Println(err)
				continue
			}
			ts[t.ID()] = t
		}
	}

	return ts, nil
}

func (s *Sleeper) GetSquads() (map[domain.SquadID]*domain.Squad, error) {
	if time.Now().Sub(s.lastSquadFetch) <= squadMapTTL {
		return s.sqmap, nil
	}

	us, err := s.sc.GetUsers(s.id)
	if err != nil {
		return nil, err
	}
	// sort by User ID
	sort.Slice(us, func(i, j int) bool {
		return us[i].UserID < us[j].UserID
	})

	rs, err := s.sc.GetRosters(s.id)
	if err != nil {
		return nil, err
	}
	// sort by User ID
	sort.Slice(rs, func(i, j int) bool {
		return rs[i].UserID < rs[j].UserID
	})

	if len(us) != len(rs) {
		return nil, errors.New("user/roster mismatch")
	}

	sqs := make(map[domain.SquadID]*domain.Squad)

	for i := range us {
		// both us and rs should be sorted by user id, but let's double check here
		if us[i].UserID != rs[i].UserID {
			return nil, errors.New("user id mismatch")
		}

		bench := make(map[domain.PlayerID]bool)
		for _, p := range rs[i].Players {
			bench[domain.PlayerID(p)] = true
		}
		for _, p := range rs[i].Starters {
			delete(bench, domain.PlayerID(p))
		}
		for _, p := range rs[i].IR {
			delete(bench, domain.PlayerID(p))
		}

		sqid := domain.SquadID(fmt.Sprint(rs[i].RosterID))
		sqs[sqid] = &domain.Squad{
			Name:          us[i].Metadata.TeamName,
			OwnerName:     us[i].DisplayName,
			SquadID:       sqid,
			Wins:          rs[i].Settings.Wins,
			Losses:        rs[i].Settings.Losses,
			PointsFor:     rs[i].PointsFor(),
			PointsAgainst: rs[i].PointsAgainst(),
			WaiverBudget:  100 - rs[i].Settings.WaiverBudgetUsed,
			Starters: utils.Map(
				func(p playerID) domain.PlayerID { return domain.PlayerID(p) }, rs[i].Starters),
			IR: utils.Map(
				func(p playerID) domain.PlayerID { return domain.PlayerID(p) }, rs[i].IR),
			Bench: utils.Keys(bench),
		}
	}

	s.sqmap = sqs
	s.lastSquadFetch = time.Now()
	return sqs, nil
}

func (s *Sleeper) GetPlayers() (map[domain.PlayerID]*domain.Player, error) {
	updated, err := s.fetchPlayerFileIfNeeded()
	if err != nil {
		return nil, err
	}
	if updated {
		if err = s.loadPlayerMapFromFile(); err != nil {
			return nil, err
		}
	}

	return s.pmap, nil
}

func (s *Sleeper) loadPlayerMapFromFile() error {
	path, err := filepath.Abs(playersJSON)
	if err != nil {
		return err
	}

	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var players map[playerID]*player
	err = json.Unmarshal(fileData, &players)
	if err != nil {
		return err
	}

	s.pmap = make(map[domain.PlayerID]*domain.Player)
	for k, v := range players {
		s.pmap[domain.PlayerID(k)] = &domain.Player{
			FirstName: v.FirstName,
			LastName:  v.LastName,
			PlayerID:  domain.PlayerID(v.PlayerID),
		}
	}

	return nil
}

func (s *Sleeper) fetchPlayerFileIfNeeded() (bool, error) {
	path, err := filepath.Abs(playersJSON)
	if err != nil {
		return false, err
	}

	stat, err := os.Stat(path)
	if err == nil && time.Now().Sub(stat.ModTime()) < playerMapTTL {
		return false, nil
	}

	pmap, err := s.sc.GetPlayers()
	if err != nil {
		return false, err
	}

	bytes, err := json.Marshal(pmap)
	if err != nil {
		return false, err
	}

	// get file
	f, err := os.Create(playersJSON)
	if err != nil {
		return false, err
	}

	_, err = f.Write(bytes)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *Sleeper) handleFreeAgentTransaction(txn *transaction) (domain.Transaction, error) {
	fa := new(domain.FreeAgentTransaction)
	fa.TransactionID = domain.TransactionID(txn.TransactionID)
	fa.Timestamp = time.UnixMilli(int64(txn.TimestampMillis))

	if len(txn.InvolvedRosters) != 1 {
		return nil, errors.New("no roster tied to txn")
	}
	fa.SquadID = domain.SquadID(fmt.Sprint(txn.InvolvedRosters[0]))

	if len(txn.Adds)+len(txn.Drops) != 1 {
		return nil, errors.New("unexpected number of add/drops")
	}
	for add := range txn.Adds {
		addPlayerID := domain.PlayerID(add)
		fa.Add = &addPlayerID
	}
	for drop := range txn.Drops {
		dropPlayerID := domain.PlayerID(drop)
		fa.Drop = &dropPlayerID
	}

	return fa, nil
}

func (s *Sleeper) handleWaiverTransaction(txn *transaction) (domain.Transaction, error) {
	w := new(domain.WaiverTransaction)
	w.TransactionID = domain.TransactionID(txn.TransactionID)
	w.Timestamp = time.UnixMilli(int64(txn.TimestampMillis))

	if len(txn.InvolvedRosters) != 1 {
		return nil, errors.New("no roster tied to txn")
	}
	w.SquadID = domain.SquadID(fmt.Sprint(txn.InvolvedRosters[0]))
	w.DidWin = (txn.Status == transactionStatusSuccess)
	w.Bid = txn.Settings.WaiverBid

	if len(txn.Adds) != 1 {
		return nil, errors.New("unexpected number of add/drops")
	}
	w.Add = domain.PlayerID(utils.Keys(txn.Adds)[0])
	for drop := range txn.Drops {
		dropPlayerID := domain.PlayerID(drop)
		w.Drop = &dropPlayerID
	}

	return w, nil
}

func (s *Sleeper) handleTradeTransaction(txn *transaction) (domain.Transaction, error) {
	t := new(domain.TradeTransaction)
	t.TransactionID = domain.TransactionID(txn.TransactionID)
	t.Timestamp = time.UnixMilli(int64(txn.TimestampMillis))

	t.Transfers = make(map[domain.SquadID]*domain.TradeTransfer)
	for i := range txn.InvolvedRosters {
		s := domain.SquadID(fmt.Sprint(txn.InvolvedRosters[i]))
		t.Transfers[s] = new(domain.TradeTransfer)
	}
	for k, v := range txn.Adds {
		s := domain.SquadID(fmt.Sprint(v))
		t.Transfers[s].PlayersGained = append(t.Transfers[s].PlayersGained, domain.PlayerID(k))
	}
	for k, v := range txn.Drops {
		s := domain.SquadID(fmt.Sprint(v))
		t.Transfers[s].PlayersLost = append(t.Transfers[s].PlayersLost, domain.PlayerID(k))
	}
	for _, w := range txn.WaiverBudget {
		s := domain.SquadID(fmt.Sprint(w.Sender))
		t.Transfers[s].WaiverLost = append(t.Transfers[s].WaiverLost, w.Amount)
		r := domain.SquadID(fmt.Sprint(w.Receiver))
		t.Transfers[r].WaiverGained = append(t.Transfers[r].WaiverGained, w.Amount)
	}

	return t, nil
}
