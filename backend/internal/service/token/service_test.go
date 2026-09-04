package token

import (
	"testing"
	"time"
)

func TestServiceIssueAndParse(t *testing.T) {
	service := New("test-secret-with-at-least-thirty-two-characters", time.Hour)
	token, _, err := service.Issue("user-1")
	if err != nil {
		t.Fatalf("Issue() error = %v", err)
	}
	subject, err := service.Parse(token)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if subject != "user-1" {
		t.Fatalf("Parse() subject = %q, want user-1", subject)
	}
	if _, err := New("another-test-secret-with-at-least-32-characters", time.Hour).Parse(token); err == nil {
		t.Fatal("Parse() accepted token signed with another secret")
	}
}
