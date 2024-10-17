package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/kbgod/illuminate"
	"github.com/kbgod/illuminate/router"
	"github.com/kbgod/tg-bot-layout/internal/service"
)

type Handler struct {
	svc *service.Service
	bot *illuminate.Bot
}

func New(svc *service.Service, bot *illuminate.Bot) *Handler {
	return &Handler{svc: svc, bot: bot}
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

	h.svc.Observer.Logger.Info().Str("username", h.bot.Username).Msg("bot stopped")
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
