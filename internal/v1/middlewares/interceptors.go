package middewares

import (
	"context"
	"orderService/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// deprecated: UnaryServerInterceptorLogger - интерцептор для логирования входящих запросов, а также ответов сервера
func UnaryServerInterceptorLogger(cl logger.CurrentLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		//получаем reqID из контекста запроса клиента
		reqID, ok := logger.RequestID(ctx)
		//еслии нет, то создаём его и выдаём ошибку
		if !ok {
			reqID = uuid.NewString()
			ctx = logger.WithRequestID(ctx, reqID)
			cl.Error(ctx, "Missing metadata in context", zap.String("code", codes.InvalidArgument.String()))
			//По-хорошему нужно будет потом возвращать это
			//return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
		}
		//Логируем входящий запрос
		cl.Info(ctx, "Incoming RPC", zap.String("method", info.FullMethod), zap.Any("request", req), zap.String("request_id", reqID))
		m, err := handler(ctx, req)

		if err != nil {
			cl.Error(ctx, "RPC failed", zap.Error(err), zap.String("code", codes.Internal.String()))
		} else {
			cl.Info(ctx, "RPC succeeded", zap.Any("response", m))
		}
		return m, err
	}
}
