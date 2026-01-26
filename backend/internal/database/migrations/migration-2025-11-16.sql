-- Migration: Modify teams unique constraint and add session token to scenarios

ALTER TABLE teams DROP CONSTRAINT IF EXISTS teams_sport_id_espn_id_key;
ALTER TABLE teams ADD CONSTRAINT teams_season_espn_id_unique UNIQUE (season_id, espn_id);

ALTER TABLE scenarios ADD COLUMN session_token VARCHAR(255);
ALTER TABLE scenarios ALTER COLUMN user_id DROP NOT NULL;
CREATE INDEX idx_scenarios_session_token ON scenarios(session_token);