package main

import (
	"context"
	"github.com/kbgod/illuminate"
	zerologadapter "github.com/kbgod/illuminate/log/adapter/zerolog"
	"github.com/kbgod/pigfish/config"
	"github.com/kbgod/pigfish/internal/database"
	"github.com/kbgod/pigfish/internal/handler"
	observerpkg "github.com/kbgod/pigfish/internal/observer"
	"github.com/kbgod/pigfish/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	observer := observerpkg.New(cfg.LogLevel, cfg.Debug)

	db, err := gorm.Open(postgres.Open(cfg.PostgresDSN()), &gorm.Config{})
	if err != nil {
		observer.Logger.Fatal().Err(err).Msg("connect to database")
	}
	if cfg.DBDebug {
		db = db.Debug()
	}

	migrator := database.NewMigrator(db, observer)

	if len(os.Args) > 1 {
		if os.Args[1] == "fresh" && !cfg.FreshAllowed {
			observer.Logger.Fatal().Msg("fresh command not allowed")
		}
		migrator.RunCommand(os.Args[1])
	}

	svc := service.New(db, observer)

	ctx, cancel := context.WithCancel(context.Background())
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-exit
		cancel()
	}()

	botClient := illuminate.NewBot(
		illuminate.WithToken(cfg.BotToken),
		illuminate.WithLogger(zerologadapter.NewAdapter(observer.Logger)),
	)
	me, err := botClient.GetMe(ctx, nil)
	if err != nil {
		observer.Logger.Fatal().Err(err).Msg("get bot info")
	}
	observer.Logger.Info().Str("username", me.Username.PeerID()).Msg("bot authorized")

	h := handler.New(svc, botClient, me)

	if err := h.Run(ctx); err != nil {
		observer.Logger.Fatal().Err(err).Msg("run handler")
	}
}
