-- MovieMash Database Schema

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE,
    email VARCHAR(255) UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Movies table
CREATE TABLE IF NOT EXISTS movies (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    year INTEGER NOT NULL,
    director VARCHAR(255),
    poster_url TEXT,
    letterboxd_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Top 4 Sets table
CREATE TABLE IF NOT EXISTS top4_sets (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    movie_ids INTEGER[] NOT NULL, -- Array of 4 movie IDs
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CHECK (array_length(movie_ids, 1) = 4)
);

-- Comparisons table
CREATE TABLE IF NOT EXISTS comparisons (
    id SERIAL PRIMARY KEY,
    set_a_id INTEGER REFERENCES top4_sets(id),
    set_b_id INTEGER REFERENCES top4_sets(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Votes table
CREATE TABLE IF NOT EXISTS votes (
    id SERIAL PRIMARY KEY,
    comparison_id INTEGER REFERENCES comparisons(id),
    user_id INTEGER REFERENCES users(id),
    winner_set_id INTEGER REFERENCES top4_sets(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_votes_comparison_id ON votes(comparison_id);
CREATE INDEX IF NOT EXISTS idx_votes_user_id ON votes(user_id);
CREATE INDEX IF NOT EXISTS idx_votes_winner_set_id ON votes(winner_set_id);
CREATE INDEX IF NOT EXISTS idx_top4_sets_user_id ON top4_sets(user_id);
CREATE INDEX IF NOT EXISTS idx_comparisons_set_a_id ON comparisons(set_a_id);
CREATE INDEX IF NOT EXISTS idx_comparisons_set_b_id ON comparisons(set_b_id);

