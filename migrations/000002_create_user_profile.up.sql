CREATE TABLE IF NOT EXISTS profile (
    user_id bigint PRIMARY KEY REFERENCES users ON DELETE CASCADE,
    gender INTEGER NOT NULL CHECK(gender IN (0, 1, 2, 3)) DEFAULT 3,
    prefer_male BOOLEAN NOT NULL CHECK (prefer_male IN (0, 1)) DEFAULT 1,
    prefer_female BOOLEAN NOT NULL CHECK (prefer_male IN (0, 1)) DEFAULT 1,
    bio TEXT NOT NULL DEFAULT ''
);