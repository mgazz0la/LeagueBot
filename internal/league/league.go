package league

import (
	"fmt"
	"strings"
	"sync"

	"github.com/mgazz0la/leaguebot/internal/domain"
	"github.com/mgazz0la/leaguebot/internal/platform"
	"github.com/mgazz0la/leaguebot/internal/utils"
)

type (
	LeagueState struct {
		mu sync.Mutex

		players map[domain.PlayerID]*domain.Player
		squads  map[domain.SquadID]*domain.Squad
	}
)

func SquadToString(s *domain.Squad, pmap map[domain.PlayerID]*domain.Player) string {
	if s == nil {
		return ""
	}
	return fmt.Sprintf("%s\n%s\n%s", s.Name, strings.Join(utils.Map(func(p domain.PlayerID) string {
		return fmt.Sprintf("%s\n", pmap[p].String())
	}, s.Starters), "\n"),
		strings.Join(utils.Map(func(p domain.PlayerID) string {
			return fmt.Sprintf("%s\n", pmap[p].String())
		}, s.Bench), "\n"))
}

func (ls *LeagueState) GetSquads() map[domain.SquadID]*domain.Squad {
	return ls.squads
}

func (ls *LeagueState) GetPlayerMap() map[domain.PlayerID]*domain.Player {
	return ls.players
}

func (ls *LeagueState) Load(p platform.Platform) error {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	ls.players = make(map[domain.PlayerID]*domain.Player)
	ls.squads = make(map[domain.SquadID]*domain.Squad)

	sqs, err := p.GetSquads()
	if err != nil {
		return err
	}

	for _, sq := range sqs {
		ls.squads[sq.SquadID] = sq
	}

	players, err := p.GetPlayers()
	if err != nil {
		return err
	}

	for _, p := range players {
		ls.players[p.PlayerID] = p
	}

	return nil
}
