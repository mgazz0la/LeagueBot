package watcher

import (
	"log"
	"sort"
	"time"

	"golang.org/x/exp/slices"

	"github.com/mgazz0la/leaguebot/internal/controller/discord"
	"github.com/mgazz0la/leaguebot/internal/domain"
	"github.com/mgazz0la/leaguebot/internal/utils"
)

func fetchTransactions(bs *discord.BotState) map[domain.TransactionID]domain.Transaction {
	return bs.League.GetTransactions()
}

func doTransactionsDiffer(
	current map[domain.TransactionID]domain.Transaction,
	other map[domain.TransactionID]domain.Transaction,
) bool {
	for k := range other {
		if !slices.Contains(utils.Keys(current), k) {
			return true
		}
	}
	return false
}

func sendNewTransactionsToDiscord(
	bs *discord.BotState,
	old map[domain.TransactionID]domain.Transaction,
	current map[domain.TransactionID]domain.Transaction,
) {
	var newTxns []domain.Transaction
	for k, v := range current {
		if _, ok := old[k]; !ok {
			newTxns = append(newTxns, v)
		}
	}

	sort.Slice(newTxns, func(i, j int) bool {
		return newTxns[i].CompletedAt().Before(newTxns[j].CompletedAt())
	})

	for i := range newTxns {
		m, err := bs.League.TransactionToDiscordMessage(newTxns[i])
		if err != nil {
			log.Println(err.Error())
			continue
		}
		log.Println("new transactions!")
		bs.Session.ChannelMessageSendComplex(bs.GuildConfig.NotificationChannelID, m)
	}
}

func NewTransactionWatcher(
	bs *discord.BotState,
) Watcher[map[domain.TransactionID]domain.Transaction] {
	return NewWatcher(
		bs,
		5*time.Second,
		fetchTransactions,
		doTransactionsDiffer,
		sendNewTransactionsToDiscord,
	)
}
