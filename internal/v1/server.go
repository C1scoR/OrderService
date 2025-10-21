package v1

import (
	pb "orderService/api"
	"orderService/pkg/repository"
	"sync"
)

type OrderServiceServer struct {
	pb.UnimplementedOrderServiceServer
	OrdersRepository repository.OrderService // read-only after initialized
	mu               sync.Mutex              // protect savedOrders
}
