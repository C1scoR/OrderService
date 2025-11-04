package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	pb "orderService/api"
	"orderService/internal/config"
	v1 "orderService/internal/v1"
	"orderService/pkg/logger"
	"orderService/pkg/logger/zaplogger"
	"orderService/pkg/repository"
	"orderService/pkg/repository/mapstorage"
	"os"
	"os/signal"
	"syscall"
	"time"

	//"orderService/internal/v1/gateway"
	"orderService/internal/v1/gateway"
	mws "orderService/internal/v1/middlewares"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Printf("Не удалось спарсить конфиг: %v", errors.Unwrap(err))
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	/*
		Далее просто передаём cl как currentLogger. Если захоти его поменять на другой, то просто в
		функции ниже нужно поменять NewLoggerAdapter от другого логгера и всё :)
	*/
	cl := logger.NewCurrentLogger(zaplogger.NewLoggerAdapter(cfg.Environment))
	/*
		type UnaryServerInterceptor func(ctx context.Context, req any, info *UnaryServerInfo, handler UnaryHandler) (resp any, err error)
	*/
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			mws.UnaryServerInterceptorLogger(*cl),
		),
	)
	repo := repository.NewOrderService(mapstorage.NewMapStorage())
	pb.RegisterOrderServiceServer(grpcServer, v1.NewServer(*repo))
	reflection.Register(grpcServer)
	httpPort := cfg.Port + 1
	httpEndpoint := fmt.Sprintf("0.0.0.0:%d", httpPort)
	grpcEndpoint := fmt.Sprintf("0.0.0.0:%d", cfg.Port)
	httpServer := gateway.ProvideHTTP(httpEndpoint, grpcEndpoint, grpcServer)
	go func() {
		log.Printf("gRPC server is running on localhost:%d", lis.Addr())
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("main/grpcServer.Serve не удалось запустить gRPC сервер: %v", err)
		}
	}()

	httpLis, err := net.Listen("tcp", httpEndpoint)
	if err != nil {
		log.Fatalf("failed to listen for HTTP gateway: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		time.Sleep(10 * time.Second)
		log.Println("Shutting down server...")
		httpServer.Shutdown(context.Background())
	}()

	if err := httpServer.Serve(httpLis); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server error: %v", err)
	}

	log.Println("Server exited properly")

	// // Keep the main function alive until an interrupt signal is received
	// select {}

}

// func newGateway(ctx context.Context, conn *grpc.ClientConn, opts []gwruntime.ServeMuxOption) (http.Handler, error) {
// 	mux := gwruntime.NewServeMux(opts...)

// 	if err := pb.RegisterOrderServiceHandler(ctx, mux, conn); err != nil {
// 		return nil, err
// 	}
// 	return mux, nil
// }

// type Endpoint struct {
// 	Network, Addr string
// }

// type Options struct {
// 	// Addr is the address to listen
// 	Addr string

// 	// GRPCServer defines an endpoint of a gRPC service
// 	GRPCServer Endpoint

// 	// OpenAPIDir is a path to a directory from which the server
// 	// serves OpenAPI specs.
// 	OpenAPIDir string

// 	// Mux is a list of options to be passed to the gRPC-Gateway multiplexer
// 	Mux []gwruntime.ServeMuxOption
// }

// func dial(network, addr string) (*grpc.ClientConn, error) {
// 	switch network {
// 	case "tcp":
// 		return dialTCP(addr)
// 	case "unix":
// 		return dialUnix(addr)
// 	default:
// 		return nil, fmt.Errorf("unsupported network type %q", network)
// 	}
// }

// func dialTCP(addr string) (*grpc.ClientConn, error) {
// 	return grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
// }

// func dialUnix(addr string) (*grpc.ClientConn, error) {
// 	d := func(ctx context.Context, addr string) (net.Conn, error) {
// 		return (&net.Dialer{}).DialContext(ctx, "unix", addr)
// 	}
// 	return grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(d))
// }

// func Run(ctx context.Context, opts Options) error {
// 	ctx, cancel := context.WithCancel(ctx)
// 	defer cancel()

// 	conn, err := dial(opts.GRPCServer.Network, opts.GRPCServer.Addr)
// 	if err != nil {
// 		return err
// 	}
// 	//graceful shutdown клиента
// 	go func() {
// 		<-ctx.Done()
// 		if err := conn.Close(); err != nil {
// 			grpclog.Errorf("Failed to close a client connection to the gRPC server: %v", err)
// 		}
// 	}()

// 	mux := http.NewServeMux()
// 	//mux.HandleFunc("/v1/create", )
// 	// mux.HandleFunc("/openapiv2/", openAPIServer(opts.OpenAPIDir))
// 	// mux.HandleFunc("/healthz", healthzServer(conn))

// 	gw, err := newGateway(ctx, conn, opts.Mux)
// 	if err != nil {
// 		return err
// 	}
// 	mux.Handle("/", gw)

// 	// Do not use logRequestBody for ExcessBodyServer because it will perform
// 	// io.ReadAll and mask the issue:
// 	// https://github.com/grpc-ecosystem/grpc-gateway/issues/5236
// 	// hmux := http.NewServeMux()
// 	// hmux.Handle("/rpc/excess-body/", allowCORS(mux))

// 	s := &http.Server{
// 		Addr:    opts.Addr,
// 		Handler: mux,
// 	}
// 	//graceful shutdown сервера
// 	go func() {
// 		<-ctx.Done()
// 		grpclog.Infof("Shutting down the http server")
// 		if err := s.Shutdown(context.Background()); err != nil {
// 			grpclog.Errorf("Failed to shutdown http server: %v", err)
// 		}
// 	}()

// 	grpclog.Infof("Starting listening at %s", opts.Addr)
// 	if err := s.ListenAndServe(); err != http.ErrServerClosed {
// 		grpclog.Errorf("Failed to listen and serve: %v", err)
// 		return err
// 	}
// 	return nil
// }

// var (
// 	endpoint   = flag.String("endpoint", "localhost:9090", "endpoint of the gRPC service")
// 	network    = flag.String("network", "tcp", `one of "tcp" or "unix". Must be consistent to -endpoint`)
// 	openAPIDir = flag.String("openapi_dir", "examples/internal/proto/examplepb", "path to the directory which contains OpenAPI definitions")
// )

// func main() {
// 	flag.Parse()

// 	ctx := context.Background()
// 	opts := Options{
// 		Addr: ":8080",
// 		GRPCServer: Endpoint{
// 			Network: *network,
// 			Addr:    *endpoint,
// 		},
// 		OpenAPIDir: *openAPIDir,
// 	}
// 	if err := Run(ctx, opts); err != nil {
// 		grpclog.Fatal(err)
// 	}
// }
