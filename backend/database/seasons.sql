-- Insert new active seasons for NFL, NBA, and CFB for the current year

INSERT INTO seasons (sport_id, start_year, end_year, is_active) VALUES (
    (SELECT id FROM sports WHERE short_name = 'NFL'), 2025, 2026, TRUE
);

INSERT INTO seasons (sport_id, start_year, end_year, is_active) VALUES (
    (SELECT id FROM sports WHERE short_name = 'NBA'), 2025, 2026, TRUE
);

INSERT INTO seasons (sport_id, start_year, end_year, is_active) VALUES (
    (SELECT id FROM sports WHERE short_name = 'CFB'), 2025, 2026, TRUE
);