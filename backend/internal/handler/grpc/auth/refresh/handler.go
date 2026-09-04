package refresh

import (
	"context"

	readerv1 "github.com/deniskrylov/english-reader/backend/gen/reader/v1"
	"github.com/deniskrylov/english-reader/backend/internal/handler/grpc/auth/response"
	uc "github.com/deniskrylov/english-reader/backend/internal/usecase/auth/refresh"
)

type Handler struct{ usecase UseCase }

func New(usecase UseCase) *Handler { return &Handler{usecase: usecase} }

func (h *Handler) Refresh(ctx context.Context, request *readerv1.RefreshRequest) (*readerv1.AuthResponse, error) {
	tokens, err := h.usecase.Execute(ctx, uc.Request{
		RefreshToken: request.GetRefreshToken(),
		DeviceLabel:  request.GetDeviceLabel(),
	})
	if err != nil {
		return nil, response.Error(err)
	}
	return response.Auth(tokens), nil
}
