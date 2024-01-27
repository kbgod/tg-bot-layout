package main

import (
	"context"
	"github.com/kbgod/illuminate"
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

	botClient, err := illuminate.NewBot(cfg.BotToken, nil)
	if err != nil {
		observer.Logger.Fatal().Err(err).Msg("create bot client")
	}
	observer.Logger.Info().Str("username", botClient.Username).Msg("bot authorized")

	h := handler.New(svc, botClient, botClient.User)

	if err := h.Run(ctx); err != nil {
		observer.Logger.Fatal().Err(err).Msg("run handler")
	}
}
