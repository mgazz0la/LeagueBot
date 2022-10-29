package watcher

import (
	"log"
	"time"

	"golang.org/x/exp/slices"

	"github.com/mgazz0la/leaguebot/internal/controller/discord"
	"github.com/mgazz0la/leaguebot/internal/domain"
	"github.com/mgazz0la/leaguebot/internal/utils"
)

func NewTransactionWatcher(
	bs *discord.BotState,
) Watcher[map[domain.TransactionID]domain.Transaction] {
	return NewWatcher(
		bs,
		15*time.Second,
		func(bbs *discord.BotState) map[domain.TransactionID]domain.Transaction {
			return bbs.League.GetTransactions()
		},
		func(
			current map[domain.TransactionID]domain.Transaction,
			other map[domain.TransactionID]domain.Transaction,
		) bool {
			return slices.Equal(utils.Keys(current), utils.Keys(other))
		},
		func(
			bs *discord.BotState,
			old map[domain.TransactionID]domain.Transaction,
			current map[domain.TransactionID]domain.Transaction,
		) {
			var newTxns []domain.Transaction
			for k := range current {
				if _, ok := old[k]; !ok {
					newTxns = append(newTxns, current[k])
				}
			}
			for _, txn := range newTxns {
				m, err := bs.League.TransactionToDiscordMessage(txn)
				if err != nil {
					log.Println(err.Error())
				}
				bs.Session.ChannelMessageSendComplex(bs.GuildConfig.NotificationChannelID, m)
			}
		},
	)
}