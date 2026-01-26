-- Migration: Reset scores for upcoming games with zero scores

UPDATE games 
SET home_score = NULL, away_score = NULL 
WHERE status = 'upcoming' AND home_score = 0 AND away_score = 0;