ALTER TABLE user_reader_settings
    ADD COLUMN highlight_color TEXT NOT NULL DEFAULT 'yellow'
    CHECK (highlight_color IN ('yellow', 'blue', 'green', 'pink', 'orange', 'purple'));
