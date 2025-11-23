-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    team_id INT REFERENCES teams(id) ON DELETE CASCADE,
    name TEXT UNIQUE NOT NULL,
    is_active BOOLEAN DEFAULT true
);

CREATE TABLE pull_requests (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    author_id INT REFERENCES users(id) ON DELETE SET NULL,
    status TEXT NOT NULL CHECK (status IN ('OPEN', 'MERGED')),
    reviewer_ids INTEGER[], -- Массив ID ревьюеров (специфично для Postgres)
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down
-- SQL in section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS pull_requests;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS teams;