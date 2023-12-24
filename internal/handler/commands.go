package handler

import (
	"context"
	"fmt"
	"github.com/kbgod/illuminate"
)

func (h *Handler) initCommands(ctx context.Context) error {
	ok, err := h.bot.SetMyCommands(ctx, []illuminate.BotCommand{
		{
			Command:     "start",
			Description: "привітальне меню",
		},
		{
			Command:     "help",
			Description: "допомога",
		},
	}, nil)
	if err != nil {
		return fmt.Errorf("set commands: %w", err)
	}
	if !ok {
		return fmt.Errorf("set commands: not ok")
	}

	return nil
}
