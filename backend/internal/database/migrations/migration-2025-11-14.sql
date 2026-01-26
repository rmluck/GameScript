-- Migration: Add ESPN IDs to games and teams

ALTER TABLE games ADD COLUMN espn_id VARCHAR(32);
CREATE UNIQUE INDEX IF NOT EXISTS games_season_espn_id_idx ON games (season_id, espn_id);

ALTER TABLE teams ADD COLUMN espn_id VARCHAR(16);
CREATE UNIQUE INDEX IF NOT EXISTS teams_sport_espn_id_idx ON teams(sport_id, espn_id);