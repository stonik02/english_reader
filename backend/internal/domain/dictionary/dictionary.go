package dictionary

import "time"

type Sense struct {
	ID           int64
	PartOfSpeech string
	Translations []string
	ExampleEN    string
	ExampleRU    string
	SourceURL    string
	Attribution  string
}

type LookupResult struct {
	LemmaID       int64
	Lemma         string
	Source        string
	SourceVersion string
	Senses        []Sense
}

type CachedTranslation struct {
	Text      string
	ExpiresAt time.Time
}

type LookupResponse struct {
	LemmaID             int64
	NormalizedLemma     string
	Senses              []Sense
	SentenceTranslation string
	ProviderError       string
	ContextVerified     bool
	Source              string
	SourceVersion       string
	AlreadySaved        bool
}

type TextTranslationResponse struct {
	Text            string
	ProviderError   string
	ContextVerified bool
}
