package postgres

import (
	"context"
	"orderService/models"
	"orderService/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func (ms *userRepo) Create(ctx context.Context, user *models.User) (int, error) {
	cl := logger.LoggerFromContext(ctx)
	result := gorm.WithResult()
	err := gorm.G[models.User](ms.db, result).Create(ctx, user)
	if err != nil {
		cl.Error(ctx, "Ошибка при создании заказа в БД ", zap.Error(err))
		return 0, err
	}
	return user.ID, nil
}

func (ms *userRepo) GetByID(ctx context.Context, id int) (*models.User, error) {
	cl := logger.LoggerFromContext(ctx)

	// Using numeric primary key
	user, err := gorm.G[models.User](ms.db).Where("id = ?", id).First(ctx)
	if err != nil {
		cl.Error(ctx, "Ошибка при создании заказа в БД ", zap.Error(err))
		return &models.User{}, err
	}
	return &user, nil
}
