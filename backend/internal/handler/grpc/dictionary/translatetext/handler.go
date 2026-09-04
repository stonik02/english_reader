package translatetext

import (
	"context"

	readerv1 "github.com/deniskrylov/english-reader/backend/gen/reader/v1"
	uc "github.com/deniskrylov/english-reader/backend/internal/usecase/dictionary/translatetext"
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

func (h *Handler) TranslateText(ctx context.Context, q *readerv1.TranslateTextRequest) (*readerv1.TranslateTextResponse, error) {
	if _, err := h.tokens.Parse(q.GetAccessToken()); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid access token")
	}
	value, err := h.usecase.Execute(ctx, uc.Request{BookID: q.GetBookId(), ChapterID: q.GetChapterId(), Text: q.GetText()})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid translation request")
	}
	response := &readerv1.TranslateTextResponse{ContextVerified: value.ContextVerified, SentenceTranslation: &readerv1.SentenceTranslation{}}
	if value.ProviderError != "" {
		response.SentenceTranslation.Result = &readerv1.SentenceTranslation_ProviderError{ProviderError: value.ProviderError}
	} else {
		response.SentenceTranslation.Result = &readerv1.SentenceTranslation_TranslatedText{TranslatedText: value.Text}
	}
	return response, nil
}
