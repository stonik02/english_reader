package listmybooks

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

func (h *Handler) ListMyBooks(ctx context.Context, request *readerv1.ListMyBooksRequest) (*readerv1.UserBookPage, error) {
	userID, err := h.tokens.Parse(request.GetAccessToken())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid access token")
	}
	page, err := h.usecase.Execute(ctx, userID, request.GetCursor(), int(request.GetLimit()))
	if err != nil {
		return nil, response.Error(err)
	}

	result := &readerv1.UserBookPage{NextCursor: page.NextCursor}
	for _, book := range page.Items {
		result.Books = append(result.Books, response.UserBook(book))
	}

	return result, nil
}
