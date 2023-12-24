package handler

import "github.com/kbgod/illuminate/router"

func (h *Handler) initRoutes() *router.Router {
	botRouter := router.New(h.bot)
	botRouter.Use(h.Recovery)
	botRouter.Use(h.ErrorHandler)
	botRouter.Use(h.CallbackQueryAutoAnswer)
	botRouter.Use(h.UserMiddleware)

	botRouter.OnStart(h.Start)

	return botRouter
}
