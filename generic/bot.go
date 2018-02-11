package generic

import (
	"context"
	"sync"

	"github.com/andviro/middleware"
	"github.com/go-mixins/bot"
)

type Bot struct {
	pre, mw middleware.Middleware
	Handler middleware.Handler
	l       sync.RWMutex
}

var _ bot.Bot = (*Bot)(nil)

func (b *Bot) handle(ctx context.Context) (err error) {
	b.l.RLock()
	defer b.l.RUnlock()
	return b.pre.Use(b.mw).Then(b.Handler).Apply(ctx)
}

func (b *Bot) Run(driver bot.Driver) error {
	for driver.Next() {
		if err := b.handle(driver.Context()); err != nil {
			return err
		}
	}
	return nil
}

func (b *Bot) On(p middleware.Predicate, h middleware.Handler) {
	b.l.Lock()
	defer b.l.Unlock()
	b.mw = b.mw.On(p, h)
}

func (b *Bot) Use(mws ...middleware.Middleware) {
	b.l.Lock()
	defer b.l.Unlock()
	b.pre = b.pre.Use(mws...)
}
