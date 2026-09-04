package listentries

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/vocabulary"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

type Request struct {
	Cursor string
	Limit  int
	Query  string
}

type Response struct {
	Entries    []domain.Entry
	NextCursor string
}

type UseCase struct {
	repository Repository
	normalizer WordNormalizer
}

func New(repository Repository, normalizer WordNormalizer) *UseCase {
	return &UseCase{repository: repository, normalizer: normalizer}
}

func (u *UseCase) Execute(ctx context.Context, userID string, request Request) (Response, error) {
	if userID == "" {
		return Response{}, domain.ErrInvalidInput
	}
	limit := request.Limit
	if limit == 0 {
		limit = defaultLimit
	}
	if limit < 0 || limit > maxLimit {
		return Response{}, domain.ErrInvalidInput
	}

	cursor, err := decodeCursor(request.Cursor)
	if err != nil {
		return Response{}, domain.ErrInvalidInput
	}
	query, err := u.normalizeQuery(request.Query)
	if err != nil {
		return Response{}, domain.ErrInvalidInput
	}
	entries, err := u.repository.List(ctx, userID, domain.ListRequest{
		Cursor: cursor,
		Limit:  limit + 1,
		Query:  query,
	})
	if err != nil {
		return Response{}, err
	}

	response := Response{Entries: entries}
	if len(entries) <= limit {
		return response, nil
	}
	response.Entries = entries[:limit]
	response.NextCursor, err = encodeCursor(response.Entries[len(response.Entries)-1])
	if err != nil {
		return Response{}, err
	}

	return response, nil
}

func (u *UseCase) normalizeQuery(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", nil
	}
	return u.normalizer.Normalize(value)
}

func decodeCursor(value string) (*domain.Cursor, error) {
	if value == "" {
		return nil, nil
	}
	raw, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}
	var cursor domain.Cursor
	if err := json.Unmarshal(raw, &cursor); err != nil || cursor.ID == "" || cursor.CreatedAt.IsZero() {
		return nil, domain.ErrInvalidInput
	}
	return &cursor, nil
}

func encodeCursor(entry domain.Entry) (string, error) {
	raw, err := json.Marshal(domain.Cursor{CreatedAt: entry.CreatedAt, ID: entry.ID})
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}
