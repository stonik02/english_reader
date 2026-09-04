package lookupword

import (
	"context"
	readerv1 "github.com/deniskrylov/english-reader/backend/gen/reader/v1"
	uc "github.com/deniskrylov/english-reader/backend/internal/usecase/dictionary/lookupword"
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
func (h *Handler) LookupWord(ctx context.Context, q *readerv1.LookupWordRequest) (*readerv1.LookupWordResponse, error) {
	userID, err := h.tokens.Parse(q.GetAccessToken())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid access token")
	}
	value, err := h.usecase.Execute(ctx, uc.Request{UserID: userID, BookID: q.GetBookId(), ChapterID: q.GetChapterId(), SelectedText: q.GetSelectedText(), SentenceText: q.GetSentenceText()})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid lookup request")
	}
	response := &readerv1.LookupWordResponse{NormalizedLemma: value.NormalizedLemma, ContextVerified: value.ContextVerified, Source: &readerv1.SourceMetadata{Source: value.Source, SourceVersion: value.SourceVersion}, SentenceTranslation: &readerv1.SentenceTranslation{}, AlreadySaved: value.AlreadySaved, LemmaId: value.LemmaID}
	if value.ProviderError != "" {
		response.SentenceTranslation.Result = &readerv1.SentenceTranslation_ProviderError{ProviderError: value.ProviderError}
	} else {
		response.SentenceTranslation.Result = &readerv1.SentenceTranslation_TranslatedText{TranslatedText: value.SentenceTranslation}
	}
	for _, sense := range value.Senses {
		response.Senses = append(response.Senses, &readerv1.DictionarySense{Id: sense.ID, PartOfSpeech: sense.PartOfSpeech, Translations: sense.Translations, ExampleEn: sense.ExampleEN, ExampleRu: sense.ExampleRU, SourceUrl: sense.SourceURL, Attribution: sense.Attribution})
	}
	return response, nil
}
