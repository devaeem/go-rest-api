-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id         uuid PRIMARY KEY DEFAULT uuidv7(),
    name       varchar(100) NOT NULL,
    email      varchar(255) NOT NULL,
    password   text         NOT NULL,
    role       varchar(20)  NOT NULL DEFAULT 'user',
    created_at timestamptz,
    updated_at timestamptz,
    deleted_at timestamptz
);

-- databases created earlier by gorm AutoMigrate have no role column and a full unique index
ALTER TABLE users ADD COLUMN IF NOT EXISTS role varchar(20) NOT NULL DEFAULT 'user';
DROP INDEX IF EXISTS idx_users_email;

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_active ON users (email) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users (deleted_at);

-- +goose Down
DROP TABLE IF EXISTS users;
