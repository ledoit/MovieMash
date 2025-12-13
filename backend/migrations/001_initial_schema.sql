-- Movies table
CREATE TABLE IF NOT EXISTS movies (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    year INTEGER,
    director VARCHAR(255),
    genres TEXT[],
    letterboxd_id VARCHAR(100) UNIQUE,
    poster_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Top 4 sets table
CREATE TABLE IF NOT EXISTS top4_sets (
    id SERIAL PRIMARY KEY,
    user_letterboxd_id VARCHAR(100) NOT NULL,
    movie_ids INTEGER[] NOT NULL,
    scraped_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_top4_user_letterboxd ON top4_sets(user_letterboxd_id);

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Comparisons table
CREATE TABLE IF NOT EXISTS comparisons (
    id SERIAL PRIMARY KEY,
    set_a_id INTEGER NOT NULL REFERENCES top4_sets(id),
    set_b_id INTEGER NOT NULL REFERENCES top4_sets(id),
    votes_a INTEGER DEFAULT 0,
    votes_b INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP
);

CREATE INDEX idx_comparisons_set_a ON comparisons(set_a_id);
CREATE INDEX idx_comparisons_set_b ON comparisons(set_b_id);
CREATE INDEX idx_comparisons_expires ON comparisons(expires_at);

-- Votes table
CREATE TABLE IF NOT EXISTS votes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    comparison_id INTEGER NOT NULL REFERENCES comparisons(id),
    winner_set_id INTEGER NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_votes_user ON votes(user_id);
CREATE INDEX idx_votes_comparison ON votes(comparison_id);
CREATE INDEX idx_votes_timestamp ON votes(timestamp);

-- Recommendations table
CREATE TABLE IF NOT EXISTS recommendations (
    user_id INTEGER NOT NULL REFERENCES users(id),
    movie_id INTEGER NOT NULL REFERENCES movies(id),
    score FLOAT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, movie_id)
);

CREATE INDEX idx_recommendations_user ON recommendations(user_id);
CREATE INDEX idx_recommendations_score ON recommendations(score DESC);

-- Movie features for ML
CREATE TABLE IF NOT EXISTS movie_features (
    movie_id INTEGER PRIMARY KEY REFERENCES movies(id),
    features_json JSONB,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

