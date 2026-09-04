package saveentry

import (
	"context"

	readerv1 "github.com/deniskrylov/english-reader/backend/gen/reader/v1"
	"github.com/deniskrylov/english-reader/backend/internal/handler/grpc/vocabulary/response"
	uc "github.com/deniskrylov/english-reader/backend/internal/usecase/vocabulary/saveentry"
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

func (h *Handler) SaveEntry(ctx context.Context, request *readerv1.SaveEntryRequest) (*readerv1.SaveEntryResponse, error) {
	userID, err := h.tokens.Parse(request.GetAccessToken())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid access token")
	}
	entry, alreadySaved, err := h.usecase.Execute(ctx, userID, uc.Request{
		LemmaID:       request.GetLemmaId(),
		ChosenSenseID: request.ChosenSenseId,
		SourceForm:    request.GetSourceForm(),
	})
	if err != nil {
		return nil, response.Error(err)
	}
	return &readerv1.SaveEntryResponse{
		Entry:        response.Entry(entry),
		AlreadySaved: alreadySaved,
	}, nil
}
