package watcher

import (
	"time"

	"github.com/mgazz0la/leaguebot/internal/controller/discord"
)

type (
	Watcher[T any] struct {
		bs           *discord.BotState
		current      T
		pollRate     time.Duration
		fetch        func(bs *discord.BotState) T
		shouldUpdate func(current T, other T) bool
		onUpdate     func(bs *discord.BotState, old T, current T)
	}
)

func NewWatcher[T any](
	bs *discord.BotState,
	pollRate time.Duration,
	fetch func(bs *discord.BotState) T,
	shouldUpdate func(current T, other T) bool,
	onUpdate func(bs *discord.BotState, old T, current T),
) Watcher[T] {
	return Watcher[T]{
		bs:           bs,
		current:      fetch(bs),
		pollRate:     pollRate,
		fetch:        fetch,
		shouldUpdate: shouldUpdate,
		onUpdate:     onUpdate,
	}
}

func (w Watcher[T]) Run() {
	t := time.NewTicker(w.pollRate)
	for range t.C {
		newT := w.fetch(w.bs)
		if w.shouldUpdate(w.current, newT) {
			oldT := w.current
			w.current = newT
			w.onUpdate(w.bs, oldT, w.current)
		}
	}
}
