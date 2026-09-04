package getreadingstate

import (
	"context"
	"errors"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/reader"
)

type UseCase struct{ repository ReaderRepository }

func New(repository ReaderRepository) *UseCase { return &UseCase{repository: repository} }
func (u *UseCase) Execute(ctx context.Context, userID, bookID string) (domain.State, error) {
	if err := u.repository.EnsureReadyBook(ctx, userID, bookID); err != nil {
		return domain.State{}, err
	}
	progress, err := u.repository.Progress(ctx, userID, bookID)
	if errors.Is(err, domain.ErrNotFound) {
		chapter, chapterErr := u.repository.FirstChapter(ctx, bookID)
		if chapterErr != nil {
			return domain.State{}, chapterErr
		}
		progress = domain.Progress{ChapterID: chapter.ID, EPUBCFI: chapter.StartCFI}
	} else if err != nil {
		return domain.State{}, err
	}
	chapter, err := u.repository.Chapter(ctx, bookID, progress.ChapterID, "")
	if errors.Is(err, domain.ErrNotFound) {
		chapter, err = u.repository.FirstChapter(ctx, bookID)
		if err != nil {
			return domain.State{}, err
		}
		progress.ChapterID = chapter.ID
		progress.EPUBCFI = chapter.StartCFI
	} else if err != nil {
		return domain.State{}, err
	}
	settings, err := u.repository.Settings(ctx, userID)
	if err != nil {
		return domain.State{}, err
	}
	return domain.State{Chapter: chapter, Progress: progress, Settings: settings}, nil
}
