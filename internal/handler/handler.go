package handler

import (
	"context"
	"fmt"
	"github.com/kbgod/illuminate"
	"github.com/kbgod/illuminate/router"
	"github.com/kbgod/pigfish/internal/service"
	"time"
)

type Handler struct {
	svc     *service.Service
	bot     *illuminate.Bot
	botInfo illuminate.User
}

func New(svc *service.Service, bot *illuminate.Bot, botInfo illuminate.User) *Handler {
	return &Handler{svc: svc, bot: bot, botInfo: botInfo}
}

func (h *Handler) Run(ctx context.Context) error {
	if err := h.initCommands(); err != nil {
		return fmt.Errorf("init commands: %w", err)
	}
	r := h.initRoutes()
	updates := h.bot.GetUpdatesChan(&illuminate.GetUpdatesChanOpts{
		Buffer: 100,
		GetUpdatesOpts: &illuminate.GetUpdatesOpts{
			Timeout: 600,
			RequestOpts: &illuminate.RequestOpts{
				Timeout: 601 * time.Second,
			},
		},
		ErrorHandler: func(err error) {
			h.svc.Observer.Logger.Error().Err(err).Msg("get updates error")
		},
	})

	runWorkerPool(ctx, 100, r, updates)

	<-ctx.Done()

	h.svc.Observer.Logger.Info().Str("username", h.botInfo.Username).Msg("bot stopped")
	return nil
}

func runWorkerPool(ctx context.Context, size int, router *router.Router, updates <-chan illuminate.Update) {
	for i := 0; i < size; i++ {
		go func(id int) {
			for update := range updates {
				u := update
				_ = router.HandleUpdate(ctx, &u)
			}
		}(i)
	}
}
