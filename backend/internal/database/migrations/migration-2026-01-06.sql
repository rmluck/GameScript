-- Migration: Create playoff_series table and modify playoff_matchups

CREATE TABLE IF NOT EXISTS playoff_series (
    id SERIAL PRIMARY KEY,
    playoff_state_id INTEGER NOT NULL REFERENCES playoff_states(id) ON DELETE CASCADE,
    round INTEGER NOT NULL,
    series_order INTEGER NOT NULL,
    conference VARCHAR(100),
    higher_seed_team_id INTEGER NOT NULL REFERENCES teams(id),
    lower_seed_team_id INTEGER NOT NULL REFERENCES teams(id),
    higher_seed INTEGER NOT NULL,
    lower_seed INTEGER NOT NULL,
    picked_team_id INTEGER REFERENCES teams(id),
    predicted_higher_seed_wins INTEGER DEFAULT 0,
    predicted_lower_seed_wins INTEGER DEFAULT 0,
    best_of INTEGER DEFAULT 7,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(playoff_state_id, round, series_order, conference)
);

ALTER TABLE playoff_matchups ADD COLUMN IF NOT EXISTS playoff_series_id INTEGER REFERENCES playoff_series(id) ON DELETE CASCADE;
ALTER TABLE playoff_matchups ADD COLUMN IF NOT EXISTS game_number INTEGER;
ALTER TABLE playoff_matchups DROP CONSTRAINT IF EXISTS playoff_matchups_playoff_state_id_round_matchup_order_conference_key;
ALTER TABLE playoff_matchups ADD CONSTRAINT playoff_matchups_unique_key UNIQUE(playoff_state_id, round, matchup_order, conference, game_number);

CREATE INDEX IF NOT EXISTS idx_playoff_series_state ON playoff_series(playoff_state_id);
CREATE INDEX IF NOT EXISTS idx_playoff_series_round ON playoff_series(playoff_state_id, round);
CREATE INDEX IF NOT EXISTS idx_playoff_matchups_series ON playoff_matchups(playoff_series_id);
