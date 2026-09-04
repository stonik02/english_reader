package refreshtoken

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

type Service struct{}

func New() *Service { return &Service{} }

func (s *Service) New() (string, []byte, error) {
	raw := make([]byte, 32)

	if _, e := rand.Read(raw); e != nil {
		return "", nil, fmt.Errorf("generate refresh token: %w", e)
	}

	v := base64.RawURLEncoding.EncodeToString(raw)
	h := sha256.Sum256([]byte(v))

	return v, h[:], nil
}
func (s *Service) Hash(v string) []byte {
	h := sha256.Sum256([]byte(v))
	return h[:]
}
