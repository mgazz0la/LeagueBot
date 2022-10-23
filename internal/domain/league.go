package domain

import (
	"fmt"
)

type (
	SquadID  string
	PlayerID string

	Team struct {
		Name string
		Bye  uint
	}

	Position int

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

		Wins          uint
		Losses        uint
		PointsFor     float32
		PointsAgainst float32
		WaiverBudget  uint

		Starters []PlayerID
		Bench    []PlayerID
		IR       []PlayerID
	}

	Transaction struct {
		Type string
	}
)

func (p Player) String() string {
	return fmt.Sprintf("%s %s", p.FirstName, p.LastName)
}

const (
	QB Position = iota
	RB
	WR
	TE
	K
	DEF
)
