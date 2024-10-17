package entity

import (
	"time"

	"github.com/kbgod/tg-bot-layout/pkg/tgutil"
)

type UserRole string

type User struct {
	ID              int64 `gorm:"primaryKey;autoIncrement:false"`
	FirstName       string
	Username        string
	Role            UserRole
	BanReason       *string
	PromoID         *int64
	BotState        *string
	BotStateContext *string
	StoppedAt       *time.Time
	BannedAt        *time.Time
	CreatedAt       time.Time
}

func (u *User) EscapedName() string {
	return tgutil.Escape(u.FirstName)
}
