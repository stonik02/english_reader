package libretranslate

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

var ErrUnavailable = errors.New("translation provider unavailable")

type Service struct {
	url       string
	client    *http.Client
	maxLength int
	failures  int
	openUntil time.Time
	mutex     sync.Mutex
}

func New(url string, timeout time.Duration, maxLength int) *Service {
	return &Service{url: strings.TrimRight(url, "/"), client: &http.Client{Timeout: timeout}, maxLength: maxLength}
}
func (s *Service) Translate(ctx context.Context, text string) (string, error) {
	if len([]rune(text)) == 0 || len([]rune(text)) > s.maxLength {
		return "", fmt.Errorf("invalid translation text")
	}
	if s.isOpen() {
		return "", ErrUnavailable
	}
	for attempt := 0; attempt < 1; attempt++ {
		translated, temporary, err := s.request(ctx, text)
		if err == nil {
			s.succeed()
			return translated, nil
		}
		if !temporary {
			s.fail()
			return "", ErrUnavailable
		}
	}
	s.fail()
	return "", ErrUnavailable
}
func (s *Service) request(ctx context.Context, text string) (string, bool, error) {
	payload, _ := json.Marshal(map[string]string{"q": text, "source": "en", "target": "ru", "format": "text"})

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, s.url+"/translate", bytes.NewReader(payload))
	if err != nil {
		return "", false, err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := s.client.Do(request)
	if err != nil {
		return "", true, err
	}

	defer response.Body.Close()

	if response.StatusCode >= 500 {
		return "", true, ErrUnavailable
	}

	if response.StatusCode != http.StatusOK {
		return "", false, ErrUnavailable
	}

	var body struct {
		TranslatedText string `json:"translatedText"`
	}

	if err := json.NewDecoder(response.Body).Decode(&body); err != nil || body.TranslatedText == "" {
		return "", false, ErrUnavailable
	}

	return body.TranslatedText, false, nil
}

func (s *Service) isOpen() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return time.Now().Before(s.openUntil)
}

func (s *Service) succeed() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.failures = 0
	s.openUntil = time.Time{}
}

func (s *Service) fail() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.failures++
	if s.failures >= 3 {
		s.openUntil = time.Now().Add(15 * time.Second)
	}
}
