-- Convert PST times to UTC (add 8 hours during PST, 7 during PDT)
-- This assumes all times are currently stored as PST
UPDATE games 
SET start_time = start_time AT TIME ZONE 'America/Los_Angeles' AT TIME ZONE 'UTC'
WHERE start_time < '2025-03-09 02:00:00'  -- Before DST starts
   OR start_time >= '2025-11-02 02:00:00'; -- After DST ends

UPDATE games 
SET start_time = start_time AT TIME ZONE 'America/Los_Angeles' AT TIME ZONE 'UTC'
WHERE start_time >= '2025-03-09 02:00:00' 
  AND start_time < '2025-11-02 02:00:00';

ALTER TABLE teams ADD COLUMN alternate_logo_url VARCHAR(255);