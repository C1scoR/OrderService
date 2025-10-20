package v1

import (
	pb "orderService/api"
	"sync"
)

type routeGuideServer struct {
	pb.UnimplementedOrderServiceServer
	savedOrders []*pb.Order // read-only after initialized

	mu sync.Mutex // protect savedOrders
}
