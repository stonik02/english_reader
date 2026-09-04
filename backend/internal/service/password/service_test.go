package password

import "testing"

func TestServiceHashAndVerify(t *testing.T) {
	service := New()

	hash, err := service.Hash("correct horse battery staple")
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}
	if !service.Verify(hash, "correct horse battery staple") {
		t.Fatal("Verify() returned false for the matching password")
	}
	if service.Verify(hash, "wrong password") {
		t.Fatal("Verify() returned true for another password")
	}
}
