package refreshtoken

import (
	"bytes"
	"testing"
)

func TestServiceNewAndHash(t *testing.T) {
	service := New()
	token, hash, err := service.New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if token == "" || len(hash) == 0 {
		t.Fatal("New() returned an empty token or hash")
	}
	if !bytes.Equal(hash, service.Hash(token)) {
		t.Fatal("Hash() does not reproduce the generated token hash")
	}
}
