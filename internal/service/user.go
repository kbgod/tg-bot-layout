package service

import (
	"errors"
	"fmt"
	"github.com/kbgod/illuminate"
	"github.com/kbgod/pigfish/internal/entity"
	"gorm.io/gorm"
)

func (s *Service) GetUser(tgUser *illuminate.User, isPrivate bool, promo *string) (*entity.User, error) {
	var user entity.User
	if err := s.db.Take(&user, tgUser.ID).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("get user: %w", err)
	} else if err == nil {
		mustUpdate := make(map[string]any)
		if isPrivate && user.StoppedAt != nil {
			mustUpdate["stopped_at"] = nil
		}
		if tgUser.FirstName != user.FirstName {
			mustUpdate["first_name"] = tgUser.FirstName
		}
		if string(tgUser.Username) != user.Username {
			mustUpdate["username"] = string(tgUser.Username)
		}
		if len(mustUpdate) > 0 {
			if err := s.db.Model(&entity.User{}).Where("id", user.ID).Updates(mustUpdate).Error; err != nil {
				return nil, fmt.Errorf("update user: %w", err)
			}
		}

		return &user, nil
	}

	user.ID = int64(tgUser.ID)
	user.FirstName = tgUser.FirstName
	user.Username = string(tgUser.Username)

	if promo != nil {
		p, err := entity.GetPromoByName(s.db, *promo)
		if err != nil {
			return nil, fmt.Errorf("get promo by name: %w", err)
		}
		user.PromoID = &p.ID
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &user, nil
}
