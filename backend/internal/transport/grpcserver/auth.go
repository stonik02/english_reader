package grpcserver

import (
	"context"

	readerv1 "github.com/deniskrylov/english-reader/backend/gen/reader/v1"
	getme "github.com/deniskrylov/english-reader/backend/internal/handler/grpc/auth/getme"
	login "github.com/deniskrylov/english-reader/backend/internal/handler/grpc/auth/login"
	logout "github.com/deniskrylov/english-reader/backend/internal/handler/grpc/auth/logout"
	refresh "github.com/deniskrylov/english-reader/backend/internal/handler/grpc/auth/refresh"
	register "github.com/deniskrylov/english-reader/backend/internal/handler/grpc/auth/register"
)

// AuthService only assembles request-specific gRPC handlers into the generated
// service contract. It contains no endpoint business logic.
type AuthService struct {
	readerv1.UnimplementedAuthServiceServer
	*registerHandler
	*loginHandler
	*refreshHandler
	*logoutHandler
	*getMeHandler
}

func (s *AuthService) Register(ctx context.Context, request *readerv1.RegisterRequest) (*readerv1.AuthResponse, error) {
	return s.registerHandler.Register(ctx, request)
}

func (s *AuthService) Login(ctx context.Context, request *readerv1.LoginRequest) (*readerv1.AuthResponse, error) {
	return s.loginHandler.Login(ctx, request)
}

func (s *AuthService) Refresh(ctx context.Context, request *readerv1.RefreshRequest) (*readerv1.AuthResponse, error) {
	return s.refreshHandler.Refresh(ctx, request)
}

func (s *AuthService) Logout(ctx context.Context, request *readerv1.LogoutRequest) (*readerv1.Empty, error) {
	return s.logoutHandler.Logout(ctx, request)
}

func (s *AuthService) GetMe(ctx context.Context, request *readerv1.GetMeRequest) (*readerv1.UserResponse, error) {
	return s.getMeHandler.GetMe(ctx, request)
}

type registerHandler = register.Handler
type loginHandler = login.Handler
type refreshHandler = refresh.Handler
type logoutHandler = logout.Handler
type getMeHandler = getme.Handler

func NewAuthService(
	registerHandler *register.Handler,
	loginHandler *login.Handler,
	refreshHandler *refresh.Handler,
	logoutHandler *logout.Handler,
	getMeHandler *getme.Handler,
) *AuthService {
	return &AuthService{
		registerHandler: registerHandler,
		loginHandler:    loginHandler,
		refreshHandler:  refreshHandler,
		logoutHandler:   logoutHandler,
		getMeHandler:    getMeHandler,
	}
}
