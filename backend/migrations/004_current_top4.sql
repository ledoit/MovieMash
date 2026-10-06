ALTER TABLE users
    ADD COLUMN IF NOT EXISTS display_name VARCHAR(255);

ALTER TABLE top4_sets
    ADD COLUMN IF NOT EXISTS is_current BOOLEAN NOT NULL DEFAULT TRUE;

CREATE INDEX IF NOT EXISTS idx_top4_sets_current ON top4_sets (user_id) WHERE is_current;
