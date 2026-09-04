package password

import (
	"crypto/rand"
	"encoding/base64"
	"strings"

	"golang.org/x/crypto/argon2"
)

type Service struct{}

func New() *Service { return &Service{} }

func (s *Service) Hash(value string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(value), salt, 1, 64*1024, 4, 32)

	return "$argon2id$v=19$m=65536,t=1,p=4$" + base64.RawStdEncoding.EncodeToString(salt) + "$" + base64.RawStdEncoding.EncodeToString(hash), nil
}
func (s *Service) Verify(encoded, value string) bool {
	p := strings.Split(encoded, "$")
	if len(p) != 6 || p[1] != "argon2id" {
		return false
	}

	salt, e := base64.RawStdEncoding.DecodeString(p[4])
	if e != nil {
		return false
	}

	expected, e := base64.RawStdEncoding.DecodeString(p[5])
	if e != nil {
		return false
	}

	actual := argon2.IDKey([]byte(value), salt, 1, 64*1024, 4, uint32(len(expected)))
	if len(actual) != len(expected) {
		return false
	}

	var d byte
	for i := range actual {
		d |= actual[i] ^ expected[i]
	}

	return d == 0
}
