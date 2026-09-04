package listcatalog

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

func (h *Handler) ListCatalog(ctx context.Context, request *readerv1.ListCatalogRequest) (*readerv1.BookPage, error) {
	page, err := h.usecase.Execute(ctx, request.GetCursor(), int(request.GetLimit()))
	if err != nil {
		return nil, response.Error(err)
	}

	result := &readerv1.BookPage{NextCursor: page.NextCursor}
	for _, book := range page.Items {
		result.Books = append(result.Books, response.Book(book))
	}

	return result, nil
}
