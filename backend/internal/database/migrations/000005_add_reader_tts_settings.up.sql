ALTER TABLE user_reader_settings
    ADD COLUMN tts_voice_uri TEXT NOT NULL DEFAULT '',
    ADD COLUMN tts_rate NUMERIC(3,2) NOT NULL DEFAULT 0.90
        CHECK (tts_rate BETWEEN 0.60 AND 1.20);
