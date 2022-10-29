package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/mgazz0la/leaguebot/internal/utils"
)

type (
	AvatarID        string
	SquadID         string
	PlayerID        string
	TransactionID   string
	TransactionType int

	Team struct {
		Name string
		Bye  uint
	}

	Position string

	Player struct {
		Team Team

		PlayerID PlayerID

		FirstName string
		LastName  string
		Position  Position

		PointsByWeek []uint
	}

	Squad struct {
		Name      string
		OwnerName string
		Seed      uint
		SquadID   SquadID
		AvatarID  AvatarID

		Wins          uint
		Losses        uint
		PointsFor     float32
		PointsAgainst float32
		WaiverBudget  uint

		Starters []PlayerID
		Bench    []PlayerID
		IR       []PlayerID
	}

	Transaction interface {
		Type() TransactionType
		ID() TransactionID
		String() string
		CompletedAt() time.Time
	}

	FreeAgentTransaction struct {
		TransactionID TransactionID
		Timestamp     time.Time
		SquadID       SquadID
		Add           *PlayerID
		Drop          *PlayerID
	}

	WaiverTransaction struct {
		TransactionID TransactionID
		Timestamp     time.Time
		SquadID       SquadID
		Bid           uint
		DidWin        bool
		Add           PlayerID
		Drop          *PlayerID
	}

	TradeTransfer struct {
		PlayersGained []PlayerID
		WaiverGained  []uint
		PlayersLost   []PlayerID
		WaiverLost    []uint
	}

	TradeTransaction struct {
		TransactionID TransactionID
		Timestamp     time.Time
		Transfers     map[SquadID]*TradeTransfer
	}

	TransactionMap map[TransactionID]Transaction
)

const (
	TransactionTypeFreeAgent TransactionType = iota
	TransactionTypeWaiver
	TransactionTypeTrade
)

func (fa FreeAgentTransaction) Type() TransactionType {
	return TransactionTypeFreeAgent
}

func (fa FreeAgentTransaction) ID() TransactionID {
	return fa.TransactionID
}

func (fa FreeAgentTransaction) String() string {
	s := fmt.Sprintf("Squad %s", fa.SquadID)
	if fa.Add != nil {
		s += fmt.Sprintf(" picked up player %s", *fa.Add)
	}
	if fa.Add != nil && fa.Drop != nil {
		s += fmt.Sprintf(" and")
	}
	if fa.Drop != nil {
		s += fmt.Sprintf(" dropped player %s", *fa.Drop)
	}

	return s
}

func (fa FreeAgentTransaction) CompletedAt() time.Time {
	return fa.Timestamp
}

func (w WaiverTransaction) Type() TransactionType {
	return TransactionTypeWaiver
}

func (w WaiverTransaction) ID() TransactionID {
	return w.TransactionID
}

func (w WaiverTransaction) String() string {
	s := fmt.Sprintf("Squad %s", w.SquadID)
	if w.DidWin {
		s += fmt.Sprintf(" picked up player %s for $%d", w.Add, w.Bid)
		if w.Drop != nil {
			s += fmt.Sprintf(" and dropped player %s", *w.Drop)
		}
	} else {
		s += fmt.Sprintf(" failed to pick up player %s for $%d", w.Add, w.Bid)
	}

	return s
}

func (w WaiverTransaction) CompletedAt() time.Time {
	return w.Timestamp
}

func (tt TradeTransaction) Type() TransactionType {
	return TransactionTypeTrade
}

func (tt TradeTransaction) ID() TransactionID {
	return tt.TransactionID
}

func (tt TradeTransaction) String() string {
	var s string
	for sq, t := range tt.Transfers {
		s += fmt.Sprintf("Squad %s receives", sq)
		if len(t.PlayersGained) > 0 {
			s += " Players [" + strings.Join(utils.Map(func(p PlayerID) string {
				return string(p)
			}, t.PlayersGained), ",") + "]"
		}
		if len(t.WaiverGained) > 0 {
			s += " FAAB [" + strings.Join(utils.Map(func(f uint) string {
				return fmt.Sprintf("$%d", f)
			}, t.WaiverGained), ",") + "]"

		}
		s += fmt.Sprintf(" and gives away")
		if len(t.PlayersLost) > 0 {
			s += " Players [" + strings.Join(utils.Map(func(p PlayerID) string {
				return string(p)
			}, t.PlayersLost), ",") + "]"
		}
		if len(t.WaiverLost) > 0 {
			s += " FAAB [" + strings.Join(utils.Map(func(f uint) string {
				return fmt.Sprintf("$%d", f)
			}, t.WaiverLost), ",") + "]"

		}
		s += "\n"
	}

	return strings.TrimSuffix(s, "\n")
}

func (tt TradeTransaction) CompletedAt() time.Time {
	return tt.Timestamp
}

func (p Player) String() string {
	return fmt.Sprintf("%s %s (%s - %s)", p.FirstName, p.LastName, p.Position, p.Team.Name)
}
