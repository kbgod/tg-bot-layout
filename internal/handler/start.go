package handler

import (
	"fmt"
	"github.com/kbgod/illuminate/router"
)

func (h *Handler) Start(ctx *router.Context) error {
	u := getUserFromContext(ctx.Context)
	return ctx.ReplyVoid(fmt.Sprintf(
		"Привіт %s, я ..., бот для чатів\n"+
			"Команди: /help\n",
		u.EscapedName()))

}
