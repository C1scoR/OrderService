// для реализации интерфейса
package mapstorage

import (
	"context"
	"fmt"
	"orderService/models"

	"github.com/google/uuid"
)

/*
type OrderRepository interface {
	Create(ctx context.Context, order models.Order) (string, error)
	GetByID(ctx context.Context, id string) (models.Order, error)
	Update(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]models.Order, error)
}
*/

type mapStorage map[string]models.Order

func (ms *mapStorage) Create(ctx context.Context, order models.Order) (string, error) {
	id := uuid.New().String()
	(*ms)[id] = order
	return id, nil
}

func (ms *mapStorage) GetByID(ctx context.Context, id string) (models.Order, error) {
	order, ok := (*ms)[id]
	if !ok {
		return models.Order{}, fmt.Errorf("GetByID: order with id %s not found", id)
	}
	return order, nil
}

func NewMapStorage() *mapStorage {
	ms := make(mapStorage)
	return &ms
}

func (ms *mapStorage) Update(ctx context.Context, id string, order models.Order) error {
	_, ok := (*ms)[id]
	if !ok {
		return fmt.Errorf("Update: order with id %s not found", id)
	}
	(*ms)[id] = order
	return nil
}

func (ms *mapStorage) Delete(ctx context.Context, id string) error {
	_, ok := (*ms)[id]
	if !ok {
		return fmt.Errorf("Delete: order with id %s not found", id)
	}
	delete(*ms, id)
	return nil
}

func (ms *mapStorage) List(ctx context.Context) ([]models.Order, error) {
	orders := make([]models.Order, 0, len(*ms))
	for _, order := range *ms {
		orders = append(orders, order)
	}
	return orders, nil
}
