package saveprogress

import (
	"context"
	domain "github.com/deniskrylov/english-reader/backend/internal/domain/reader"
	"strings"
)

type Request struct {
	ChapterID       string
	EPUBCFI         string
	ProgressPercent float64
	Revision        int64
}
type UseCase struct{ repository ReaderRepository }

func New(repository ReaderRepository) *UseCase { return &UseCase{repository: repository} }
func (u *UseCase) Execute(ctx context.Context, userID, bookID string, request Request) (domain.Progress, error) {
	if request.Revision < 0 || request.ProgressPercent < 0 || request.ProgressPercent > 100 || len(request.EPUBCFI) > 512 || !strings.HasPrefix(request.EPUBCFI, "epubcfi(") {
		return domain.Progress{}, domain.ErrInvalidInput
	}
	if err := u.repository.EnsureReadyBook(ctx, userID, bookID); err != nil {
		return domain.Progress{}, err
	}
	if _, err := u.repository.Chapter(ctx, bookID, request.ChapterID, ""); err != nil {
		return domain.Progress{}, err
	}
	return u.repository.SaveProgress(ctx, userID, bookID, domain.Progress{ChapterID: request.ChapterID, EPUBCFI: request.EPUBCFI, ProgressPercent: request.ProgressPercent, Revision: request.Revision})
}
