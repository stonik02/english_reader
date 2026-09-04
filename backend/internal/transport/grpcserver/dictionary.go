package grpcserver

import (
	"context"

	readerv1 "github.com/deniskrylov/english-reader/backend/gen/reader/v1"
)

type DictionaryService struct {
	readerv1.UnimplementedDictionaryServiceServer
	lookup    LookupWordHandler
	translate TranslateTextHandler
}

type LookupWordHandler interface {
	LookupWord(context.Context, *readerv1.LookupWordRequest) (*readerv1.LookupWordResponse, error)
}
type TranslateTextHandler interface {
	TranslateText(context.Context, *readerv1.TranslateTextRequest) (*readerv1.TranslateTextResponse, error)
}

func NewDictionaryService(lookup LookupWordHandler, translate TranslateTextHandler) *DictionaryService {
	return &DictionaryService{lookup: lookup, translate: translate}
}

func (s *DictionaryService) LookupWord(ctx context.Context, request *readerv1.LookupWordRequest) (*readerv1.LookupWordResponse, error) {
	return s.lookup.LookupWord(ctx, request)
}

func (s *DictionaryService) TranslateText(ctx context.Context, request *readerv1.TranslateTextRequest) (*readerv1.TranslateTextResponse, error) {
	return s.translate.TranslateText(ctx, request)
}
