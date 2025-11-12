package v1

import (
	"context"
	"fmt"
	"net/http"
	pb "orderService/api"
	"orderService/models"
	"orderService/pkg/repository"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/*
	type OrderServiceServer interface {
		CreateOrder(context.Context, *CreateOrderRequest) (*CreateOrderResponse, error)
		GetOrder(context.Context, *GetOrderRequest) (*GetOrderResponse, error)
		UpdateOrder(context.Context, *UpdateOrderRequest) (*UpdateOrderResponse, error)
		DeleteOrder(context.Context, *DeleteOrderRequest) (*DeleteOrderResponse, error)
		ListOrders(context.Context, *ListOrdersRequest) (*ListOrdersResponse, error)
		mustEmbedUnimplementedOrderServiceServer()
	}
*/
//OrderServiceServer - это структура, которая имеет все методы для реализации интерфейса OrderServiceServer interface;
//Она позволяет добавить свою реализацию методов для gPRC сервера. Для этого создав её экземпляр через func NewServer()
//мы регистрируем обработчики (т.е. методы которые содержит эта структура) через функцию .proto файла: RegisterOrderServiceServer()
type OrderServiceServer struct {
	pb.UnimplementedOrderServiceServer                       //inherited example with stubs for OrderServiceServer interface
	Repository                         repository.Repository // read-only after initialized
	mu                                 sync.Mutex            // protect savedOrders
}

// NewServer это функция, которая создаёт экземпляр структуры OrderServiceServer.
// Она принимает экземпляр хранилища repository.OrderService, которое мы будем использоват для хранения заказов.
func NewServer(Repository repository.Repository) *OrderServiceServer {
	s := &OrderServiceServer{Repository: Repository}
	return s
}

/*
Вот этими месседжами мы оперируем в функции создания заказа

	message CreateOrderRequest {
	  string item = 1;
	  int32 quantity = 2;
	}

	message CreateOrderResponse {
	  string id = 1;
	}
*/
func (s *OrderServiceServer) CreateOrder(ctx context.Context, orderRequest *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	order := models.Order{
		Item:     orderRequest.GetItem(),
		Quantity: orderRequest.GetQuantity(),
	}
	id, err := s.Repository.Order().Create(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("CreateOrder/Create error creating an order: %v", err)
	}
	return &pb.CreateOrderResponse{Id: id}, nil
}

/*
message GetOrderRequest {
  	string id = 1;
}
message GetOrderResponse {
	Order order = 1;
}
*/

func (s *OrderServiceServer) GetOrder(ctx context.Context, orderRequest *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	order, err := s.Repository.Order().GetByID(ctx, orderRequest.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "GetOrder/GetByID error getting an order: %v", err)
	}
	responseOrder := pb.Order{
		Id:       order.ID,
		Item:     order.Item,
		Quantity: order.Quantity,
	}
	return &pb.GetOrderResponse{Order: &responseOrder}, nil
}

/*
	message UpdateOrderRequest {
	  string id = 1;
	  string item = 2;
	  int32 quantity = 3;
	}
*/

/*
	message UpdateOrderResponse {
	  Order order = 1;
	}
*/
func (s *OrderServiceServer) UpdateOrder(ctx context.Context, uor *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error) {
	order := models.Order{
		ID:       uor.GetId(),
		Item:     uor.GetItem(),
		Quantity: uor.GetQuantity(),
	}
	err := s.Repository.Order().Update(ctx, order)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "UpdateOrder/Update error updating an order: %v", err)
	}

	return &pb.UpdateOrderResponse{
		Order: &pb.Order{
			Id:       order.ID,
			Item:     order.Item,
			Quantity: order.Quantity,
		},
	}, nil
}

/*
	message DeleteOrderRequest {
	  string id = 1;
	}

	message DeleteOrderResponse {
	  bool success = 1;
	}
*/
func (s *OrderServiceServer) DeleteOrder(ctx context.Context, dor *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	err := s.Repository.Order().Delete(ctx, dor.GetId())
	if err != nil {
		return &pb.DeleteOrderResponse{Success: false}, fmt.Errorf("DeleteOrder/Delete error deleting an order: %v", err)
	}
	return &pb.DeleteOrderResponse{Success: true}, nil
}

func (s *OrderServiceServer) ListOrders(ctx context.Context, lor *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	orders, err := s.Repository.Order().List(ctx)
	if err != nil {
		return nil, fmt.Errorf("ListOrders/List error getting all orders: %v", err)
	}
	var responseOrders []*pb.Order
	for _, order := range orders {
		responseOrders = append(responseOrders, &pb.Order{
			Id:       order.ID,
			Item:     order.Item,
			Quantity: order.Quantity,
		})
	}
	return &pb.ListOrdersResponse{Orders: responseOrders}, nil
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok!", http.StatusOK)
}
