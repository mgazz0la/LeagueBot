package platform

import "github.com/mgazz0la/leaguebot/internal/domain"

type (
	Platform interface {
		GetTransactions(week uint) (map[domain.TransactionID]domain.Transaction, error)
		GetSquads() (map[domain.SquadID]*domain.Squad, error)
		GetPlayers() (map[domain.PlayerID]*domain.Player, error)
		GetCurrentWeek() (uint, error)
	}
)
