package removefrommylibrary

import (
	"context"

	readerv1 "github.com/deniskrylov/english-reader/backend/gen/reader/v1"
	"github.com/deniskrylov/english-reader/backend/internal/handler/grpc/library/response"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	usecase UseCase
	tokens  TokenParser
}

func New(usecase UseCase, tokens TokenParser) *Handler {
	return &Handler{usecase: usecase, tokens: tokens}
}

func (h *Handler) RemoveFromMyLibrary(ctx context.Context, request *readerv1.RemoveFromMyLibraryRequest) (*readerv1.Empty, error) {
	userID, err := h.tokens.Parse(request.GetAccessToken())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid access token")
	}
	if err := h.usecase.Execute(ctx, userID, request.GetBookId()); err != nil {
		return nil, response.Error(err)
	}

	return &readerv1.Empty{}, nil
}
