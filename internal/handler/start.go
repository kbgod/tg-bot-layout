package handler

import (
	"fmt"
	"github.com/kbgod/illuminate/router"
)

func (h *Handler) Start(ctx *router.Context) error {
	u := getUserFromContext(ctx.Context)
	return ctx.ReplyVoid(fmt.Sprintf(
		"Hello %s, I'm %s\n"+
			"Test FSM: /set_my_name\n",
		u.EscapedName(),
		ctx.Bot.User.Username))

}

func (h *Handler) OnMessage(ctx *router.Context) error {
	return ctx.ReplyVoid("undefined command")
}
