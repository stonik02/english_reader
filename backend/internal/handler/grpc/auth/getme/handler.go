package getme

import (
	"context"

	readerv1 "github.com/deniskrylov/english-reader/backend/gen/reader/v1"
	"github.com/deniskrylov/english-reader/backend/internal/handler/grpc/auth/response"
)

type Handler struct{ usecase UseCase }

func New(usecase UseCase) *Handler { return &Handler{usecase: usecase} }

func (h *Handler) GetMe(ctx context.Context, request *readerv1.GetMeRequest) (*readerv1.UserResponse, error) {
	user, err := h.usecase.Execute(ctx, request.GetAccessToken())
	if err != nil {
		return nil, response.Error(err)
	}
	return &readerv1.UserResponse{User: response.User(user)}, nil
}
