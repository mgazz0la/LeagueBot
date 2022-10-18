package domain

import "fmt"

type (
	Team struct {
		Name      string
		ShortName string
		Bye       uint
	}

	Position int

	Player struct {
		Team

		FirstName string
		LastName  string
		Position  Position

		PointsByWeek []uint
	}

	Squad struct {
		Name      string
		OwnerName string
		Seed      uint

		Wins          uint
		Losses        uint
		PointsFor     float32
		PointsAgainst float32
		WaiverBudget  uint

		Starters []Player
		Bench    []Player
	}
)

func (s Squad) String() string {
	return fmt.Sprintf("#%d: (%d-%d) (%0.2f) %s", s.Seed, s.Wins, s.Losses, s.PointsFor, s.Name)
}

const (
	QB Position = iota
	RB
	WR
	TE
	K
	DEF
)
