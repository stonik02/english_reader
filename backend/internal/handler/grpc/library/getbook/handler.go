package getbook

import (
	"context"

	readerv1 "github.com/deniskrylov/english-reader/backend/gen/reader/v1"
	"github.com/deniskrylov/english-reader/backend/internal/handler/grpc/library/response"
)

type Handler struct {
	usecase UseCase
}

func New(usecase UseCase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) GetBook(ctx context.Context, request *readerv1.GetBookRequest) (*readerv1.Book, error) {
	book, err := h.usecase.Execute(ctx, request.GetBookId())
	if err != nil {
		return nil, response.Error(err)
	}

	return response.Book(book), nil
}
