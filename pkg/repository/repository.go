package repository

/*
type Order struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Item     string `protobuf:"bytes,2,opt,name=item,proto3" json:"item,omitempty"`
	Quantity int32  `protobuf:"varint,3,opt,name=quantity,proto3" json:"quantity,omitempty"`
}

type OrderRepository interface {
	Create(ctx context.Context, OrderRequest api.CreateOrderRequest) string
	GetByID(ctx context.Context, _) api.Order
	Update(ctx context.Context, _) error
	Delete(ctx context.Context, _) error
	List(ctx context.Context, _)
}

type OrderServiceServer interface {
	CreateOrder(context.Context, *CreateOrderRequest) (*CreateOrderResponse, error)
	GetOrder(context.Context, *GetOrderRequest) (*GetOrderResponse, error)
	UpdateOrder(context.Context, *UpdateOrderRequest) (*UpdateOrderResponse, error)
	DeleteOrder(context.Context, *DeleteOrderRequest) (*DeleteOrderResponse, error)
	ListOrders(context.Context, *ListOrdersRequest) (*ListOrdersResponse, error)
	mustEmbedUnimplementedOrderServiceServer()
}
*/
