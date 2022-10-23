package platform

import "github.com/mgazz0la/leaguebot/internal/domain"

type (
	Platform interface {
		GetTransactions(week uint) ([]*domain.Transaction, error)
		GetSquads() ([]*domain.Squad, error)
		GetPlayers() (map[domain.PlayerID]*domain.Player, error)
	}
)
