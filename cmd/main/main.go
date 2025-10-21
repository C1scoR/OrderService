package main

import (
	"fmt"
	"log"
	"net"
	pb "orderService/api"
	"orderService/internal/config"
	v1 "orderService/internal/v1"
	"orderService/pkg/repository"
	"orderService/pkg/repository/mapstorage"

	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.ParseConfig("internal/config/.env")
	if err != nil {
		log.Fatalf("main/ParseConfig: Ошибка загрузки config: %v", err)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	repo := repository.NewOrderService(mapstorage.NewMapStorage())
	pb.RegisterOrderServiceServer(grpcServer, v1.NewServer(*repo))
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("main/grpcServer.Serve не удалось запустить сервер: %v", err)
	}
}
