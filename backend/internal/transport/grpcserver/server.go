package grpcserver

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	grpcHealth "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Server struct {
	GRPC   *grpc.Server
	health *grpcHealth.Server
}

func New(logger *slog.Logger) *Server {
	healthServer := grpcHealth.NewServer()
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(loggingInterceptor(logger), recoveryInterceptor(logger)),
	)
	healthpb.RegisterHealthServer(server, healthServer)
	reflection.Register(server)

	return &Server{GRPC: server, health: healthServer}
}

func (s *Server) SetServing(serving bool) {
	if serving {
		s.health.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
		return
	}
	s.health.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
}

func loggingInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		startedAt := time.Now()
		response, err := handler(ctx, req)
		logger.InfoContext(ctx, "grpc request",
			"method", info.FullMethod,
			"code", status.Code(err).String(),
			"duration", time.Since(startedAt).String(),
		)
		return response, err
	}
}

func recoveryInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ any, err error) {
		defer func() {
			if recovered := recover(); recovered != nil {
				logger.ErrorContext(ctx, "panic in grpc handler", "method", info.FullMethod, "panic", recovered)
				err = status.Error(codes.Internal, "internal server error")
			}
		}()
		return handler(ctx, req)
	}
}
