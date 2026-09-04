package saveentry

import (
	"context"
	"strings"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/vocabulary"
)

const maxSourceFormLength = 128

type Request struct {
	LemmaID       int64
	ChosenSenseID *int64
	SourceForm    string
}

type UseCase struct {
	repository Repository
}

func New(repository Repository) *UseCase {
	return &UseCase{repository: repository}
}

func (u *UseCase) Execute(ctx context.Context, userID string, request Request) (domain.Entry, bool, error) {
	request.SourceForm = strings.TrimSpace(request.SourceForm)
	if userID == "" || request.LemmaID <= 0 || request.SourceForm == "" || len(request.SourceForm) > maxSourceFormLength {
		return domain.Entry{}, false, domain.ErrInvalidInput
	}
	if request.ChosenSenseID != nil && *request.ChosenSenseID <= 0 {
		return domain.Entry{}, false, domain.ErrInvalidInput
	}

	return u.repository.Save(ctx, userID, domain.SaveRequest{
		LemmaID:       request.LemmaID,
		ChosenSenseID: request.ChosenSenseID,
		SourceForm:    request.SourceForm,
	})
}
