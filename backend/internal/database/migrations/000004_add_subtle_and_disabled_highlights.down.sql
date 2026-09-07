ALTER TABLE user_reader_settings
    DROP CONSTRAINT IF EXISTS user_reader_settings_highlight_color_check;

ALTER TABLE user_reader_settings
    ADD CONSTRAINT user_reader_settings_highlight_color_check
    CHECK (highlight_color IN ('yellow', 'blue', 'green', 'pink', 'orange', 'purple'));
