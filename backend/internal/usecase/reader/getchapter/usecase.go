package getchapter

import (
	"context"
	domain "github.com/deniskrylov/english-reader/backend/internal/domain/reader"
)

type UseCase struct{ repository ReaderRepository }

func New(repository ReaderRepository) *UseCase { return &UseCase{repository: repository} }
func (u *UseCase) Execute(ctx context.Context, userID, bookID, chapterID, href string) (domain.Chapter, error) {
	if err := u.repository.EnsureReadyBook(ctx, userID, bookID); err != nil {
		return domain.Chapter{}, err
	}
	return u.repository.Chapter(ctx, bookID, chapterID, href)
}
