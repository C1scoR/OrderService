// Файл содержит интерфейс для встаривания в OrderServiceServer
package repository

import (
	"context"
	"orderService/models"
)

type OrderRepository interface {
	Create(ctx context.Context, order models.Order) (string, error)
	GetByID(ctx context.Context, id string) (models.Order, error)
	Update(ctx context.Context, order models.Order) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]models.Order, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (int, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
}

type Repository interface {
	Order() OrderRepository
	User() UserRepository
}

// Главный интерфейс-контейнер. Он предоставляет доступ к репозиториям сущностей.
// type Service struct {
// 	Repository Repository
// }

// func NewService(repo Repository) *Service {
// 	return &Service{
// 		Repository: repo,
// 	}
// }

// func (s *Service) CreateOrder(ctx context.Context, order models.Order) (string, error) {
// 	return s.Repository.Order().Create(ctx, order)
// 	//return s.Repository.Create(ctx, order)
// }

// func (s *Service) GetOrderByID(ctx context.Context, id string) (models.Order, error) {
// 	return s.Repository.Order().GetByID(ctx, id)
// 	//return s.Repository.GetByID(ctx, id)
// }

// func (s *Service) UpdateOrder(ctx context.Context, order models.Order) error {
// 	return s.Repository.Order().Update(ctx, order)
// 	//return s.Repository.Update(ctx, order)
// }

// func (s *Service) DeleteOrder(ctx context.Context, id string) error {
// 	return s.Repository.Order().Delete(ctx, id)
// 	//return s.Repository.Delete(ctx, id)
// }

// func (s *Service) ListOrder(ctx context.Context) ([]models.Order, error) {
// 	return s.Repository.Order().List(ctx)
// 	//return s.Repository.List(ctx)
// }

// func (s *Service) CreateUser(ctx context.Context, user *models.User) (int, error) {
// 	return s.Repository.User().Create(ctx, user)

// }

// func (s *Service) GetUserByID(ctx context.Context, id int) (*models.User, error) {
// 	return s.Repository.User().GetByID(ctx, id)
// }
