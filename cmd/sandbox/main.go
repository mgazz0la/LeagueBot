package main

import (
	"fmt"
	"strings"

	"github.com/mgazz0la/leaguebot/internal/domain"
	"github.com/mgazz0la/leaguebot/internal/platform/sleeper"
	"github.com/mgazz0la/leaguebot/internal/utils"
)

func main() {
	s := sleeper.NewSleeper("787735231225532416")
	fmt.Println(strings.Join(utils.Map(func(t *domain.Transaction) string {
		return t.Type
	}, utils.First(s.GetTransactions(utils.GetCurrentWeek()))), "\n"))
}
