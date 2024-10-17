package database

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/kbgod/tg-bot-layout/internal/entity"
	"gorm.io/gorm"
)

func initMigration() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "init",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				&entity.Promo{},
				&entity.User{},
			)
		},
	}
}
