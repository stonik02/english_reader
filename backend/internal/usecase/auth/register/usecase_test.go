package register

import (
	"context"
	"errors"
	"testing"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/auth"
)

var errHashing = errors.New("hashing failed")

type failingHasher struct{}

func (failingHasher) Hash(string) (string, error) { return "", errHashing }

func TestUseCaseRejectsPasswordsShorterThanEightCharacters(t *testing.T) {
	useCase := New(nil, nil, nil, nil, 0)

	_, err := useCase.Execute(context.Background(), Request{
		Email:    "reader@example.com",
		Password: "1234567",
	})
	if err != domain.ErrInvalidInput {
		t.Fatalf("Execute() error = %v, want %v", err, domain.ErrInvalidInput)
	}
}

func TestUseCaseAcceptsEightCharacterPasswordForFurtherProcessing(t *testing.T) {
	useCase := New(nil, failingHasher{}, nil, nil, 0)

	_, err := useCase.Execute(context.Background(), Request{
		Email:    "reader@example.com",
		Password: "12345678",
	})
	if !errors.Is(err, errHashing) {
		t.Fatalf("Execute() error = %v, want password-hasher error", err)
	}
}
