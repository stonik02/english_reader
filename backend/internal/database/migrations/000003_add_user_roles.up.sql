CREATE TYPE user_role AS ENUM ('user', 'admin');

ALTER TABLE users
    ADD COLUMN role user_role NOT NULL DEFAULT 'user';

CREATE INDEX users_role_idx ON users (role);
