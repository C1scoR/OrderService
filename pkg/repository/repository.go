// Файл содержит интерфейс для встаривания в OrderServiceServer
package repository

import (
	"context"
	"orderService/models"
)

type OrderRepository interface {
	Create(ctx context.Context, order models.Order) (string, error)
	GetByID(ctx context.Context, id string) (models.Order, error)
	Update(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]models.Order, error)
}

type OrderService struct {
	Repository OrderRepository
}

func NewOrderService(repo OrderRepository) *OrderService {
	return &OrderService{
		Repository: repo,
	}
}

func (s *OrderService) Create(ctx context.Context, order models.Order) (string, error) {
	return s.Repository.Create(ctx, order)
}

func (s *OrderService) GetByID(ctx context.Context, id string) (models.Order, error) {
	return s.Repository.GetByID(ctx, id)
}

func (s *OrderService) Update(ctx context.Context, id string) error {
	return s.Repository.Update(ctx, id)
}

func (s *OrderService) Delete(ctx context.Context, id string) error {
	return s.Repository.Delete(ctx, id)
}

func (s *OrderService) List(ctx context.Context) ([]models.Order, error) {
	return s.Repository.List(ctx)
}
