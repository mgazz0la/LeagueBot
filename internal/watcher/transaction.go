package watcher

import (
	"log"
	"time"

	"github.com/mgazz0la/leaguebot/internal/controller/discord"
	"github.com/mgazz0la/leaguebot/internal/utils"
)

func TransactionWatcher(bs *discord.BotState) {
	t := time.NewTicker(15 * time.Second)
	txnCount := 0
	for _ = range t.C {
		txns, err := bs.Platform.GetTransactions(utils.GetCurrentWeek())
		if err != nil {
			continue
		}

		if txnCount == 0 || len(txns) == 0 {
			txnCount = len(txns)
			continue
		}

		if txnCount != len(txns) {
			txnCount = len(txns)
			log.Println("TRANSACTION!")
		}
	}
}
