CREATE TABLE nicknames
(
    id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name       TEXT UNIQUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);