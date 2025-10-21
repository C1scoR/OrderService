package main

import (
	"fmt"
	"log"
	"net"
	pb "orderService/api"
	v1 "orderService/internal/v1"
	"orderService/pkg/repository"
	"orderService/pkg/repository/mapstorage"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 50051))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	repo := repository.NewOrderService(mapstorage.NewMapStorage())
	pb.RegisterOrderServiceServer(grpcServer, v1.NewServer(*repo))
	//log.Printf("gRPC server is running on localhost:%d", lis.Addr())
	reflection.Register(grpcServer)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("main/grpcServer.Serve не удалось запустить сервер: %v", err)
	}
}
