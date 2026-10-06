-- Profiles, Elo, watches, and claim accounts.
-- Existing users are system rows (the anonymous seed owner). Letterboxd ingest sets control = simulator.

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS letterboxd_username VARCHAR(255),
    ADD COLUMN IF NOT EXISTS control VARCHAR(16) NOT NULL DEFAULT 'system',
    ADD COLUMN IF NOT EXISTS elo DOUBLE PRECISION NOT NULL DEFAULT 1500,
    ADD COLUMN IF NOT EXISTS match_count INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS wins INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS claimed_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS avatar_url TEXT,
    ADD COLUMN IF NOT EXISTS bio TEXT;

ALTER TABLE movies
    ADD COLUMN IF NOT EXISTS elo DOUBLE PRECISION NOT NULL DEFAULT 1500,
    ADD COLUMN IF NOT EXISTS match_count INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS wins INTEGER NOT NULL DEFAULT 0;

ALTER TABLE top4_sets
    ADD COLUMN IF NOT EXISTS elo DOUBLE PRECISION NOT NULL DEFAULT 1500,
    ADD COLUMN IF NOT EXISTS match_count INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS wins INTEGER NOT NULL DEFAULT 0;

ALTER TABLE votes
    ADD COLUMN IF NOT EXISTS actor VARCHAR(16) NOT NULL DEFAULT 'guest',
    ADD COLUMN IF NOT EXISTS voter_key TEXT,
    ADD COLUMN IF NOT EXISTS loser_set_id INTEGER REFERENCES top4_sets(id);

ALTER TABLE votes ALTER COLUMN user_id DROP NOT NULL;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'users_control_check'
    ) THEN
        ALTER TABLE users
            ADD CONSTRAINT users_control_check
            CHECK (control IN ('simulator', 'human', 'system'));
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'votes_actor_check'
    ) THEN
        ALTER TABLE votes
            ADD CONSTRAINT votes_actor_check
            CHECK (actor IN ('guest', 'account', 'simulator'));
    END IF;
END $$;

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_letterboxd_username
    ON users (letterboxd_username)
    WHERE letterboxd_username IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_movies_letterboxd_id
    ON movies (letterboxd_id)
    WHERE letterboxd_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_users_elo ON users (elo DESC);
CREATE INDEX IF NOT EXISTS idx_movies_elo ON movies (elo DESC);
CREATE INDEX IF NOT EXISTS idx_top4_sets_elo ON top4_sets (elo DESC);
CREATE INDEX IF NOT EXISTS idx_votes_voter_key ON votes (voter_key);

CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    user_id INTEGER REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS watches (
    user_id INTEGER NOT NULL REFERENCES users(id),
    movie_id INTEGER NOT NULL REFERENCES movies(id),
    watched_at TIMESTAMPTZ,
    letterboxd_rating DOUBLE PRECISION,
    PRIMARY KEY (user_id, movie_id)
);

CREATE INDEX IF NOT EXISTS idx_watches_movie_id ON watches (movie_id);
CREATE INDEX IF NOT EXISTS idx_accounts_user_id ON accounts (user_id);
