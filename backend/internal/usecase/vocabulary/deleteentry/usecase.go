package deleteentry

import (
	"context"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/vocabulary"
)

type Request struct {
	EntryID string
	LemmaID int64
}

type UseCase struct {
	repository Repository
}

func New(repository Repository) *UseCase {
	return &UseCase{repository: repository}
}

func (u *UseCase) Execute(ctx context.Context, userID string, request Request) error {
	if userID == "" || (request.EntryID == "" && request.LemmaID <= 0) || (request.EntryID != "" && request.LemmaID > 0) {
		return domain.ErrInvalidInput
	}
	if request.EntryID != "" {
		return u.repository.DeleteByID(ctx, userID, request.EntryID)
	}
	return u.repository.DeleteByLemmaID(ctx, userID, request.LemmaID)
}
