package platform

import "github.com/mgazz0la/leaguebot/internal/domain"

type (
	Platform interface {
		GetSquads() ([]*domain.Squad, error)
	}
)
