package sleeper

import (
	"errors"
	"fmt"

	"github.com/mgazz0la/leaguebot/internal/domain"
	"github.com/mgazz0la/leaguebot/internal/platform"
	"github.com/mgazz0la/leaguebot/internal/utils"
)

type Sleeper struct {
	sc *sleeperClient
	id LeagueID
}

func NewSleeper(id LeagueID) platform.Platform {
	return &Sleeper{
		sc: newSleeperClient(),
		id: id,
	}
}

func (s *Sleeper) GetTransactions(week uint) ([]*domain.Transaction, error) {
	txns, err := s.sc.GetTransactions(s.id, week)
	return utils.Map(func(t transaction) *domain.Transaction {
		return &domain.Transaction{
			Type: t.Type,
		}
	}, txns), err
}

func (s *Sleeper) GetSquads() ([]*domain.Squad, error) {
	us, err := s.sc.GetUsers(s.id)
	if err != nil {
		return nil, err
	}

	rs, err := s.sc.GetRosters(s.id)
	if err != nil {
		return nil, err
	}

	if len(us) != len(rs) {
		return nil, errors.New("user/roster mismatch")
	}

	sqs := make([]*domain.Squad, len(us))

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

		sqs[i] = &domain.Squad{
			Name:          us[i].Metadata.TeamName,
			OwnerName:     us[i].DisplayName,
			SquadID:       domain.SquadID(fmt.Sprint(rs[i].RosterID)),
			Wins:          rs[i].Settings.Wins,
			Losses:        rs[i].Settings.Losses,
			PointsFor:     rs[i].PointsFor(),
			PointsAgainst: rs[i].PointsAgainst(),
			WaiverBudget:  100 - rs[i].Settings.WaiverBudgetUsed,
			Starters:      utils.Map(func(p playerID) domain.PlayerID { return domain.PlayerID(p) }, rs[i].Starters),
			Bench:         utils.MapKeys(bench),
			IR:            utils.Map(func(p playerID) domain.PlayerID { return domain.PlayerID(p) }, rs[i].IR),
		}
	}

	return sqs, nil
}

func (s *Sleeper) GetPlayers() (map[domain.PlayerID]*domain.Player, error) {
	players, err := s.sc.GetPlayers()
	if err != nil {
		return nil, err
	}

	m := make(map[domain.PlayerID]*domain.Player)
	for k, v := range players {
		m[domain.PlayerID(k)] = &domain.Player{
			FirstName: v.FirstName,
			LastName:  v.LastName,
			PlayerID:  domain.PlayerID(v.PlayerID),
		}
	}

	return m, nil
}
