package handler

import "github.com/kbgod/illuminate/router"

func (h *Handler) initRoutes() *router.Router {
	botRouter := router.New(h.bot)
	botRouter.Use(h.Recovery)
	botRouter.Use(h.ErrorHandler)
	botRouter.Use(h.CallbackQueryAutoAnswer)
	botRouter.Use(h.UserMiddleware)

	// global events (accessible from any state)
	botRouter.OnStart(h.Start)
	botRouter.OnCommand("set_my_name", h.SetUserName)

	enterNameScene := botRouter.UseState("enter_name")

	// this handler will be called only if user is in "enter_name" state
	enterNameScene.OnCommand("inside", h.InsideEnterName)
	enterNameScene.OnCommand("exit", h.ExitFromEnterName)
	enterNameScene.OnMessage(h.OnName)

	// this handler will be called only if routes with states and type OnMessage are not defined
	botRouter.OnMessage(h.OnMessage)

	return botRouter
}
