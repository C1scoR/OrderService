package gateway

import (
	"context"
	"log"
	"net/http"
	"strings"

	pb "orderService/api"
	v1 "orderService/internal/v1"
	"orderService/pkg/swagger"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ProvideHTTP(httpEndpoint, grpcEndpoint string, grpcServer *grpc.Server) *http.Server {
	ctx := context.Background()
	//creds, err := credentials.NewClientTLSFromFile("../tls/server.pem", "go-grpc-example")
	// if err != nil {
	// 	log.Fatalf("Failed to create TLS credentials %v", err)
	// }
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	gwmux := runtime.NewServeMux()
	err := pb.RegisterOrderServiceHandlerFromEndpoint(ctx, gwmux, grpcEndpoint, dopts)
	if err != nil {
		log.Fatalf("Register Endpoint err: %v", err)
	}
	mux := http.NewServeMux()
	mux.Handle("/", gwmux)
	mux.HandleFunc("/swagger.json", swagger.ServeSwaggerFile)
	mux.HandleFunc("/health", v1.HealthCheck)
	swagger.ServeSwaggerUI(mux)
	log.Println(httpEndpoint + " HTTP.Listing whth TLS and token...")
	return &http.Server{
		Addr:    httpEndpoint,
		Handler: grpcHandlerFunc(grpcServer, mux),
	}
}

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}
