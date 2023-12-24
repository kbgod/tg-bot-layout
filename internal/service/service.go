package service

import (
	observerpkg "github.com/kbgod/pigfish/internal/observer"
	"gorm.io/gorm"
)

type Service struct {
	db       *gorm.DB
	Observer *observerpkg.Observer
}

func New(db *gorm.DB, observer *observerpkg.Observer) *Service {
	return &Service{db: db, Observer: observer}
}
