package league

import (
	"sync"

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

func (ls *LeagueState) GetTransactionByID(txid domain.TransactionID) (domain.Transaction, bool) {
	ls.txmu.Lock()
	defer ls.txmu.Unlock()
	tx, ok := ls.transactions[txid]
	return tx, ok
}

func (ls *LeagueState) GetSquads() (map[domain.SquadID]*domain.Squad, error) {
	return ls.platform.GetSquads()
}
