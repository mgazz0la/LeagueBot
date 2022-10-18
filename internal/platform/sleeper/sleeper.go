package sleeper

import (
	"errors"

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

		sqs[i] = &domain.Squad{
			Name:          us[i].Metadata.TeamName,
			OwnerName:     us[i].DisplayName,
			Wins:          rs[i].Settings.Wins,
			Losses:        rs[i].Settings.Losses,
			PointsFor:     rs[i].PointsFor(),
			PointsAgainst: rs[i].PointsAgainst(),
			WaiverBudget:  100 - rs[i].Settings.WaiverBudgetUsed,
		}
	}

	utils.ApplySeeds(sqs)

	return sqs, nil
}
