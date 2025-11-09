// для реализации интерфейса
package mapstorage

import (
	"context"
	"fmt"
	"orderService/models"
	"sync"

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

type mapStorage struct {
	m  map[string]models.Order
	mu sync.RWMutex
}

func NewMapStorage() *mapStorage {
	return &mapStorage{
		m:  make(map[string]models.Order),
		mu: sync.RWMutex{},
	}
}

func (ms *mapStorage) Create(ctx context.Context, order models.Order) (string, error) {
	id := uuid.New().String()
	order.ID = id
	ms.mu.Lock()
	ms.m[id] = order
	ms.mu.Unlock()
	return id, nil
}

func (ms *mapStorage) GetByID(ctx context.Context, id string) (models.Order, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	order, ok := ms.m[id]
	if !ok {
		return models.Order{}, fmt.Errorf("GetByID: order with id %s not found", id)
	}
	return order, nil
}

func (ms *mapStorage) Update(ctx context.Context, order models.Order) error {
	ms.mu.RLock()
	_, ok := ms.m[order.ID]
	ms.mu.RUnlock()
	if !ok {
		return fmt.Errorf("Update: order with id %s not found", order.ID)
	}
	ms.mu.Lock()
	ms.m[order.ID] = order
	ms.mu.Unlock()
	return nil
}

func (ms *mapStorage) Delete(ctx context.Context, id string) error {
	ms.mu.RLock()
	_, ok := ms.m[id]
	ms.mu.RUnlock()
	if !ok {
		return fmt.Errorf("Delete: order with id %s not found", id)
	}
	ms.mu.Lock()
	delete(ms.m, id)
	ms.mu.Unlock()
	return nil
}

func (ms *mapStorage) List(ctx context.Context) ([]models.Order, error) {
	orders := make([]models.Order, 0, len(ms.m))
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	for _, order := range ms.m {
		orders = append(orders, order)
	}
	return orders, nil
}
