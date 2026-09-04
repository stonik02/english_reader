package gethighlights

import (
	"context"
	"errors"
	"math"

	readerv1 "github.com/deniskrylov/english-reader/backend/gen/reader/v1"
	domain "github.com/deniskrylov/english-reader/backend/internal/domain/vocabulary"
	uc "github.com/deniskrylov/english-reader/backend/internal/usecase/vocabulary/gethighlights"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	usecase UseCase
	tokens  TokenParser
}

func New(usecase UseCase, tokens TokenParser) *Handler {
	return &Handler{
		usecase: usecase,
		tokens:  tokens,
	}
}

func (h *Handler) GetHighlights(ctx context.Context, request *readerv1.GetHighlightsRequest) (*readerv1.GetHighlightsResponse, error) {
	userID, err := h.tokens.Parse(request.GetAccessToken())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid access token")
	}

	highlights, err := h.usecase.Execute(ctx, uc.Request{
		UserID:    userID,
		BookID:    request.GetBookId(),
		ChapterID: request.GetChapterId(),
	})
	if err != nil {
		return nil, mapError(err)
	}

	response := &readerv1.GetHighlightsResponse{
		Tokens: make([]*readerv1.HighlightToken, 0, len(highlights)),
	}
	for _, highlight := range highlights {
		positions := make([]int32, 0, len(highlight.Positions))
		for _, position := range highlight.Positions {
			if position > math.MaxInt32 {
				return nil, status.Error(codes.Internal, "highlight position is too large")
			}
			positions = append(positions, int32(position))
		}
		response.Tokens = append(response.Tokens, &readerv1.HighlightToken{
			LemmaId:   highlight.LemmaID,
			Lemma:     highlight.Lemma,
			Positions: positions,
			MatchKind: protobufMatchKind(highlight.MatchKind),
			Texts:     highlight.Texts,
		})
	}

	return response, nil
}

func protobufMatchKind(value domain.HighlightMatchKind) readerv1.HighlightToken_MatchKind {
	if value == domain.HighlightMatchKindLemma {
		return readerv1.HighlightToken_MATCH_KIND_LEMMA
	}
	return readerv1.HighlightToken_MATCH_KIND_EXACT_FALLBACK
}

func mapError(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, "invalid highlights request")
	case errors.Is(err, domain.ErrNotReady):
		return status.Error(codes.FailedPrecondition, "book is not ready")
	case errors.Is(err, domain.ErrNotFound):
		return status.Error(codes.NotFound, "book or chapter not found")
	default:
		return status.Error(codes.Internal, "get highlights")
	}
}
