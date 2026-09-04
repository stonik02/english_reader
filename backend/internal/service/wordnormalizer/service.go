package wordnormalizer

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidWord = errors.New("invalid selected word")

type Service struct{ maxLength int }

func New(maxLength int) *Service { return &Service{maxLength: maxLength} }
func (s *Service) Normalize(value string) (string, error) {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.TrimFunc(value, func(r rune) bool { return unicode.IsPunct(r) || unicode.IsSpace(r) })
	if value == "" || len([]rune(value)) > s.maxLength || strings.ContainsAny(value, "<>\n\r\t ") {
		return "", ErrInvalidWord
	}
	for _, r := range value {
		if !unicode.IsLetter(r) && r != '\'' {
			return "", ErrInvalidWord
		}
	}
	return value, nil
}
