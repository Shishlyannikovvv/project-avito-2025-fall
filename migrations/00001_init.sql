-- +goose Up
CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    team_id INTEGER REFERENCES teams(id) ON DELETE CASCADE,
    is_active BOOLEAN DEFAULT true NOT NULL
);

CREATE TABLE pull_requests (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    author_id INTEGER REFERENCES users(id),
    status TEXT CHECK (status IN ('OPEN', 'MERGED')) DEFAULT 'OPEN',
    reviewer1_id INTEGER REFERENCES users(id),
    reviewer2_id INTEGER REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    merged_at TIMESTAMPTZ,
    UNIQUE(author_id, id)
);

-- +goose Down
DROP TABLE pull_requests;
DROP TABLE users;
DROP TABLE teams;
