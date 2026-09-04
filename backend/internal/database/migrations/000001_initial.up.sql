CREATE EXTENSION IF NOT EXISTS citext;

CREATE TYPE book_status AS ENUM ('processing', 'ready', 'failed');
CREATE TYPE book_add_source AS ENUM ('button', 'first_read');

CREATE TABLE users (
    id UUID PRIMARY KEY,
    email CITEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE auth_sessions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token_hash BYTEA NOT NULL UNIQUE,
    device_label TEXT NOT NULL DEFAULT '',
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX auth_sessions_user_expires_idx ON auth_sessions (user_id, expires_at);

CREATE TABLE books (
    id UUID PRIMARY KEY,
    uploaded_by_user_id UUID NOT NULL REFERENCES users(id),
    title TEXT NOT NULL,
    author TEXT NOT NULL DEFAULT '',
    format TEXT NOT NULL CHECK (format = 'epub'),
    status book_status NOT NULL DEFAULT 'processing',
    content_sha256 BYTEA NOT NULL UNIQUE,
    source_file_path TEXT NOT NULL UNIQUE,
    cover_file_path TEXT,
    failure_code TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX books_status_created_idx ON books (status, created_at DESC);

CREATE TABLE ingestion_jobs (
    id UUID PRIMARY KEY,
    book_id UUID NOT NULL UNIQUE REFERENCES books(id) ON DELETE CASCADE,
    state TEXT NOT NULL CHECK (state IN ('pending', 'processing', 'completed', 'failed')),
    attempt INTEGER NOT NULL DEFAULT 0 CHECK (attempt >= 0),
    locked_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,
    error_detail TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX ingestion_jobs_active_idx ON ingestion_jobs (state, created_at)
    WHERE state IN ('pending', 'processing');

CREATE TABLE book_chapters (
    id UUID PRIMARY KEY,
    book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    sequence INTEGER NOT NULL CHECK (sequence >= 0),
    href TEXT NOT NULL,
    start_cfi TEXT NOT NULL,
    end_cfi TEXT NOT NULL,
    sanitized_html TEXT NOT NULL,
    plain_text TEXT NOT NULL,
    word_count INTEGER NOT NULL DEFAULT 0 CHECK (word_count >= 0),
    UNIQUE (book_id, sequence),
    UNIQUE (book_id, href)
);

CREATE TABLE user_books (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    added_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    added_via book_add_source NOT NULL,
    PRIMARY KEY (user_id, book_id)
);
CREATE INDEX user_books_user_added_idx ON user_books (user_id, added_at DESC);

CREATE TABLE reading_progress (
    user_id UUID NOT NULL,
    book_id UUID NOT NULL,
    chapter_id UUID REFERENCES book_chapters(id) ON DELETE SET NULL,
    epub_cfi TEXT NOT NULL,
    progress_percent NUMERIC(5,2) NOT NULL CHECK (progress_percent >= 0 AND progress_percent <= 100),
    revision BIGINT NOT NULL CHECK (revision >= 0),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, book_id),
    FOREIGN KEY (user_id, book_id) REFERENCES user_books(user_id, book_id) ON DELETE CASCADE
);

CREATE TABLE user_reader_settings (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    font_scale SMALLINT NOT NULL DEFAULT 100 CHECK (font_scale BETWEEN 80 AND 200),
    theme TEXT NOT NULL DEFAULT 'system' CHECK (theme IN ('system', 'light', 'dark')),
    line_height NUMERIC(3,2) NOT NULL DEFAULT 1.50 CHECK (line_height BETWEEN 1.00 AND 3.00),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE dictionary_lemmas (
    id BIGSERIAL PRIMARY KEY,
    language CHAR(2) NOT NULL,
    lemma CITEXT NOT NULL,
    source TEXT NOT NULL,
    source_version TEXT NOT NULL,
    UNIQUE (language, lemma, source, source_version)
);

CREATE TABLE dictionary_senses (
    id BIGSERIAL PRIMARY KEY,
    lemma_id BIGINT NOT NULL REFERENCES dictionary_lemmas(id) ON DELETE CASCADE,
    part_of_speech TEXT NOT NULL,
    translations JSONB NOT NULL DEFAULT '[]'::jsonb,
    example_en TEXT,
    example_ru TEXT,
    source_url TEXT NOT NULL,
    attribution TEXT NOT NULL,
    position INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX dictionary_senses_lemma_position_idx ON dictionary_senses (lemma_id, position);

CREATE TABLE vocabulary_entries (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    lemma_id BIGINT NOT NULL REFERENCES dictionary_lemmas(id),
    chosen_sense_id BIGINT REFERENCES dictionary_senses(id) ON DELETE SET NULL,
    source_form TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, lemma_id)
);
CREATE INDEX vocabulary_entries_user_created_idx ON vocabulary_entries (user_id, created_at DESC);

CREATE TABLE translation_cache (
    key_hash BYTEA PRIMARY KEY,
    model_version TEXT NOT NULL,
    source_lang CHAR(2) NOT NULL,
    target_lang CHAR(2) NOT NULL,
    normalized_text TEXT NOT NULL,
    translated_text TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX translation_cache_expires_idx ON translation_cache (expires_at);

CREATE TABLE dictionary_import_runs (
    id UUID PRIMARY KEY,
    source TEXT NOT NULL,
    source_version TEXT NOT NULL,
    checksum BYTEA NOT NULL,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    finished_at TIMESTAMPTZ,
    imported_lemmas INTEGER NOT NULL DEFAULT 0,
    error_detail TEXT
);
