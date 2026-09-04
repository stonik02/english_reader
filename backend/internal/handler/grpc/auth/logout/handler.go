package logout

import (
	"context"

	readerv1 "github.com/deniskrylov/english-reader/backend/gen/reader/v1"
	"github.com/deniskrylov/english-reader/backend/internal/handler/grpc/auth/response"
)

type Handler struct{ usecase UseCase }

func New(usecase UseCase) *Handler { return &Handler{usecase: usecase} }

func (h *Handler) Logout(ctx context.Context, request *readerv1.LogoutRequest) (*readerv1.Empty, error) {
	if err := h.usecase.Execute(ctx, request.GetRefreshToken()); err != nil {
		return nil, response.Error(err)
	}
	return &readerv1.Empty{}, nil
}
