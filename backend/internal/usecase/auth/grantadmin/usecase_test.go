package grantadmin

import (
	"context"
	"testing"

	domain "github.com/deniskrylov/english-reader/backend/internal/domain/auth"
)

type users struct{ user domain.User }

func (u users) FindUserByEmail(context.Context, string) (domain.User, string, error) {
	return u.user, "", nil
}

type roles struct {
	id   string
	role string
}

func (r *roles) SetRole(_ context.Context, id, role string) error {
	r.id = id
	r.role = role
	return nil
}

func TestUseCaseGrantsAdminRoleToExistingUser(t *testing.T) {
	updater := &roles{}
	result, err := New(users{user: domain.User{ID: "user-1", Email: "reader@example.com", Role: "user"}}, updater).Execute(context.Background(), " READER@example.com ")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if updater.id != "user-1" || updater.role != "admin" {
		t.Fatalf("updated = (%q, %q), want user-1/admin", updater.id, updater.role)
	}
	if result.Role != "admin" {
		t.Fatalf("result role = %q, want admin", result.Role)
	}
}
