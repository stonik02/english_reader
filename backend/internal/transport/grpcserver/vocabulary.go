package grpcserver

import (
	"context"

	readerv1 "github.com/deniskrylov/english-reader/backend/gen/reader/v1"
)

type VocabularyService struct {
	readerv1.UnimplementedVocabularyServiceServer
	saveEntry   SaveEntryHandler
	listEntries ListEntriesHandler
	deleteEntry DeleteEntryHandler
	highlights  GetHighlightsHandler
}

type SaveEntryHandler interface {
	SaveEntry(context.Context, *readerv1.SaveEntryRequest) (*readerv1.SaveEntryResponse, error)
}

type ListEntriesHandler interface {
	ListEntries(context.Context, *readerv1.ListEntriesRequest) (*readerv1.ListEntriesResponse, error)
}

type DeleteEntryHandler interface {
	DeleteEntry(context.Context, *readerv1.DeleteEntryRequest) (*readerv1.DeleteEntryResponse, error)
}

type GetHighlightsHandler interface {
	GetHighlights(context.Context, *readerv1.GetHighlightsRequest) (*readerv1.GetHighlightsResponse, error)
}

func NewVocabularyService(
	saveEntry SaveEntryHandler,
	listEntries ListEntriesHandler,
	deleteEntry DeleteEntryHandler,
	highlights GetHighlightsHandler,
) *VocabularyService {
	return &VocabularyService{
		saveEntry:   saveEntry,
		listEntries: listEntries,
		deleteEntry: deleteEntry,
		highlights:  highlights,
	}
}

func (s *VocabularyService) SaveEntry(ctx context.Context, request *readerv1.SaveEntryRequest) (*readerv1.SaveEntryResponse, error) {
	return s.saveEntry.SaveEntry(ctx, request)
}

func (s *VocabularyService) ListEntries(ctx context.Context, request *readerv1.ListEntriesRequest) (*readerv1.ListEntriesResponse, error) {
	return s.listEntries.ListEntries(ctx, request)
}

func (s *VocabularyService) DeleteEntry(ctx context.Context, request *readerv1.DeleteEntryRequest) (*readerv1.DeleteEntryResponse, error) {
	return s.deleteEntry.DeleteEntry(ctx, request)
}

func (s *VocabularyService) GetHighlights(ctx context.Context, request *readerv1.GetHighlightsRequest) (*readerv1.GetHighlightsResponse, error) {
	return s.highlights.GetHighlights(ctx, request)
}
