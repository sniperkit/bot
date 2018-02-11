package main

import (
	"context"
	"os"

	"github.com/go-mixins/log/logrus"

	"github.com/go-mixins/bot"
	"github.com/go-mixins/bot/telegram"
)

func main() {
	logger := logrus.New()
	logger.Level = logrus.DebugLevel
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		logger.Fatalf("please set environment variable BOT_TOKEN")
	}
	b, err := telegram.New(token)
	if err != nil {
		logger.Fatalf("%+v", err)
	}
	defer b.Close()
	b.Handle(func(ctx context.Context) error {
		logger.Debug(telegram.DebugDump(ctx))
		return nil
	})
	b.On(b.Command("start"), func(ctx context.Context) error {
		return b.Reply(ctx, "Hello!")
	})
	b.On(b.Command("quit"), func(ctx context.Context) error {
		return b.Reply(ctx, "Bye!")
	})
	b.On(b.Hears("lol"), func(ctx context.Context) error {
		return b.Reply(ctx, "lol yourself", b.WithReply)
	})
	b.On(b.Message(bot.MsgText), func(ctx context.Context) error {
		return b.Reply(ctx, b.Msg(ctx).Text)
	})
	b.On(b.Message(bot.MsgSticker), func(ctx context.Context) error {
		return b.Reply(ctx, "sticker!")
	})
	b.On(b.Message(bot.MsgNewChatMembers), func(ctx context.Context) error {
		for _, u := range b.NewChatMembers(ctx) {
			if u.UserName == b.Me(ctx).UserName {
				continue
			}
			if err := b.Reply(ctx, "Hi, "+u.UserName+"!"); err != nil {
				return err
			}
		}
		return nil
	})
	if err = b.Run(); err != nil {
		logger.Fatalf("%+v", err)
	}
}
