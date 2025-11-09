package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	pb "orderService/api"
	"orderService/internal/config"
	v1 "orderService/internal/v1"
	"orderService/internal/v1/gateway"
	"orderService/pkg/logger"
	"orderService/pkg/logger/zaplogger"
	"orderService/pkg/repository"
	"orderService/pkg/repository/postgres"
	"os"
	"os/signal"
	"syscall"
	"time"

	mws "orderService/internal/v1/middlewares"

	"orderService/pkg/repository/postgres/migrations"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	config     *config.Config
	logger     logger.Logger
	grpcServer *grpc.Server
	httpServer *http.Server
}

func New(cfg *config.Config) (*App, error) {
	// Инициализация логгера
	cl := logger.NewCurrentLogger(zaplogger.NewLoggerAdapter(cfg.Environment))
	// Инициализация БД и репозитория
	Store := postgres.New(cfg.PostgreSQL)
	repo := repository.NewOrderService(Store)
	//Запуск миграций
	migr := Store.GetGorm()
	migrations.Migrate(migr)
	// Создание gRPC сервера с middleware
	grpcServer := newGRPCServer(cl)
	// Регистрация сервисов
	pb.RegisterOrderServiceServer(grpcServer, v1.NewServer(*repo))
	reflection.Register(grpcServer)
	httpServer := newHTTPServer(cfg, grpcServer)

	return &App{
		config:     cfg,
		logger:     cl,
		grpcServer: grpcServer,
		httpServer: httpServer,
	}, nil
}

func (a *App) Run() error {
	// Запуск gRPC сервера
	grpcLis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", a.config.Port))
	if err != nil {
		return fmt.Errorf("failed to listen gRPC: %w", err)
	}

	// Запуск HTTP сервера
	httpPort := a.config.Port + 1
	httpLis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", httpPort))
	if err != nil {
		return fmt.Errorf("failed to listen HTTP: %w", err)
	}

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Запуск серверов в горутинах
	go a.runGRPCServer(grpcLis)
	go a.runHTTPServer(httpLis)

	// Ожидание сигнала завершения
	<-ctx.Done()
	return a.shutdown()
}

func (a *App) runGRPCServer(lis net.Listener) {
	a.logger.Info(context.Background(), "gRPC сервер запущен на порту: %d", a.config.Port)
	if err := a.grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC server error: %v", err)
	}
}

func (a *App) runHTTPServer(lis net.Listener) {
	if err := a.httpServer.Serve(lis); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server error: %v", err)
	}
}

func (a *App) shutdown() error {
	log.Println("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("HTTP server shutdown error: %w", err)
	}

	a.grpcServer.GracefulStop()
	log.Println("Server exited properly")
	return nil
}

func newGRPCServer(logger *logger.CurrentLogger) *grpc.Server {
	return grpc.NewServer(
		grpc.UnaryInterceptor(
			mws.UnaryServerInterceptorLogger(logger),
		),
	)
}

func newHTTPServer(cfg *config.Config, grpcServer *grpc.Server) *http.Server {
	httpPort := cfg.Port + 1
	httpEndpoint := fmt.Sprintf("0.0.0.0:%d", httpPort)
	grpcEndpoint := fmt.Sprintf("0.0.0.0:%d", cfg.Port)

	return gateway.ProvideHTTP(httpEndpoint, grpcEndpoint, grpcServer)
}
