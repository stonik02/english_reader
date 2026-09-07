CREATE TABLE dictionary_pronunciations (
    id BIGSERIAL PRIMARY KEY,
    lemma_id BIGINT NOT NULL REFERENCES dictionary_lemmas(id) ON DELETE CASCADE,
    ipa TEXT NOT NULL DEFAULT '',
    accent TEXT NOT NULL DEFAULT '',
    audio_url TEXT NOT NULL DEFAULT '',
    source_url TEXT NOT NULL,
    attribution TEXT NOT NULL,
    license TEXT NOT NULL DEFAULT '',
    UNIQUE (lemma_id, ipa, accent, audio_url)
);

CREATE INDEX dictionary_pronunciations_lemma_idx ON dictionary_pronunciations (lemma_id, id);
