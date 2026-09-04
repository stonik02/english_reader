package grpcserver

import (
	"context"

	readerv1 "github.com/deniskrylov/english-reader/backend/gen/reader/v1"
)

// DictionaryService is registered before the LookupWord application layer is
// wired. It keeps the public contract discoverable without pretending that the
// endpoint is implemented.
type DictionaryService struct {
	readerv1.UnimplementedDictionaryServiceServer
	handler LookupWordHandler
}

type LookupWordHandler interface {
	LookupWord(context.Context, *readerv1.LookupWordRequest) (*readerv1.LookupWordResponse, error)
}

func NewDictionaryService(handler LookupWordHandler) *DictionaryService {
	return &DictionaryService{handler: handler}
}

func (s *DictionaryService) LookupWord(ctx context.Context, request *readerv1.LookupWordRequest) (*readerv1.LookupWordResponse, error) {
	return s.handler.LookupWord(ctx, request)
}
