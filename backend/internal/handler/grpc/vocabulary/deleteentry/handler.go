package deleteentry

import (
	"context"

	readerv1 "github.com/deniskrylov/english-reader/backend/gen/reader/v1"
	"github.com/deniskrylov/english-reader/backend/internal/handler/grpc/vocabulary/response"
	uc "github.com/deniskrylov/english-reader/backend/internal/usecase/vocabulary/deleteentry"
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

func (h *Handler) DeleteEntry(ctx context.Context, request *readerv1.DeleteEntryRequest) (*readerv1.DeleteEntryResponse, error) {
	userID, err := h.tokens.Parse(request.GetAccessToken())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid access token")
	}
	deleteRequest := uc.Request{}
	switch target := request.GetTarget().(type) {
	case *readerv1.DeleteEntryRequest_EntryId:
		deleteRequest.EntryID = target.EntryId
	case *readerv1.DeleteEntryRequest_LemmaId:
		deleteRequest.LemmaID = target.LemmaId
	}
	if err := h.usecase.Execute(ctx, userID, deleteRequest); err != nil {
		return nil, response.Error(err)
	}
	return &readerv1.DeleteEntryResponse{}, nil
}
