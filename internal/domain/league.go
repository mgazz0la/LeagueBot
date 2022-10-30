package domain

import (
	"fmt"
	"time"
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

func (fa FreeAgentTransaction) CompletedAt() time.Time {
	return fa.Timestamp
}

func (w WaiverTransaction) Type() TransactionType {
	return TransactionTypeWaiver
}

func (w WaiverTransaction) ID() TransactionID {
	return w.TransactionID
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

func (tt TradeTransaction) CompletedAt() time.Time {
	return tt.Timestamp
}

func (p Player) String() string {
	return fmt.Sprintf("%s %s (%s - %s)", p.FirstName, p.LastName, p.Position, p.Team.Name)
}
