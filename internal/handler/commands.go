package handler

import (
	"fmt"
	"github.com/kbgod/illuminate"
)

func (h *Handler) initCommands() error {
	ok, err := h.bot.SetMyCommands([]illuminate.BotCommand{
		{
			Command:     "start",
			Description: "start page",
		},
		{
			Command:     "set_my_name",
			Description: "test FSM",
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
