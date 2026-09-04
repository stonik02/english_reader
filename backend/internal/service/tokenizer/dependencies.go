package tokenizer

type WordNormalizer interface {
	Normalize(string) (string, error)
}

type Morphology interface {
	Lemma(string) (string, bool)
}
