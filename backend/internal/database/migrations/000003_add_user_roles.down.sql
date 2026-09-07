DROP INDEX users_role_idx;

ALTER TABLE users
    DROP COLUMN role;

DROP TYPE user_role;
