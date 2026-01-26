-- Migration: Remove the foreign key constraint on picked_team_id
ALTER TABLE picks DROP CONSTRAINT IF EXISTS picks_picked_team_id_fkey;

-- picked_team_id can now be:
-- - NULL (no pick made)
-- - 0 (tie picked)
-- - Any valid team_id (team picked to win)