package response

import (
	"errors"
	"time"

	readerv1 "github.com/deniskrylov/english-reader/backend/gen/reader/v1"
	domain "github.com/deniskrylov/english-reader/backend/internal/domain/vocabulary"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Entry(entry domain.Entry) *readerv1.VocabularyEntry {
	response := &readerv1.VocabularyEntry{
		Id:         entry.ID,
		LemmaId:    entry.LemmaID,
		Lemma:      entry.Lemma,
		SourceForm: entry.SourceForm,
		CreatedAt:  entry.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt:  entry.UpdatedAt.Format(time.RFC3339Nano),
	}
	if entry.ChosenSense != nil {
		response.ChosenSense = &readerv1.VocabularySense{
			Id:           entry.ChosenSense.ID,
			PartOfSpeech: entry.ChosenSense.PartOfSpeech,
			Translations: entry.ChosenSense.Translations,
			ExampleEn:    entry.ChosenSense.ExampleEN,
			ExampleRu:    entry.ChosenSense.ExampleRU,
		}
	}
	return response
}

func Error(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidInput), errors.Is(err, domain.ErrInvalidSense):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
