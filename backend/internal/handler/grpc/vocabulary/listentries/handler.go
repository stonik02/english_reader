package listentries

import (
	"context"

	readerv1 "github.com/deniskrylov/english-reader/backend/gen/reader/v1"
	"github.com/deniskrylov/english-reader/backend/internal/handler/grpc/vocabulary/response"
	uc "github.com/deniskrylov/english-reader/backend/internal/usecase/vocabulary/listentries"
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

func (h *Handler) ListEntries(ctx context.Context, request *readerv1.ListEntriesRequest) (*readerv1.ListEntriesResponse, error) {
	userID, err := h.tokens.Parse(request.GetAccessToken())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid access token")
	}
	value, err := h.usecase.Execute(ctx, userID, uc.Request{
		Cursor: request.GetCursor(),
		Limit:  int(request.GetLimit()),
		Query:  request.GetQuery(),
	})
	if err != nil {
		return nil, response.Error(err)
	}
	result := &readerv1.ListEntriesResponse{NextCursor: value.NextCursor}
	for _, entry := range value.Entries {
		result.Entries = append(result.Entries, response.Entry(entry))
	}
	return result, nil
}
