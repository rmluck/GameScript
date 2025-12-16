CREATE TABLE playoff_states (
    id SERIAL PRIMARY KEY,
    scenario_id INTEGER NOT NULL REFERENCES scenarios(id) ON DELETE CASCADE,
    current_round INTEGER DEFAULT 0,
    is_enabled BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(scenario_id)
);

CREATE TABLE playoff_matchups (
    id SERIAL PRIMARY KEY,
    playoff_state_id INTEGER NOT NULL REFERENCES playoff_states(id) ON DELETE CASCADE,
    round INTEGER NOT NULL,
    matchup_order INTEGER NOT NULL,
    conference VARCHAR(100),
    higher_seed_team_id INTEGER NOT NULL REFERENCES teams(id),
    lower_seed_team_id INTEGER NOT NULL REFERENCES teams(id),
    higher_seed INTEGER NOT NULL,
    lower_seed INTEGER NOT NULL,
    picked_team_id INTEGER REFERENCES teams(id),
    predicted_higher_seed_score INTEGER,
    predicted_lower_seed_score INTEGER,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(playoff_state_id, round, matchup_order, conference)
);

CREATE INDEX idx_playoff_states_scenario ON playoff_states(scenario_id);
CREATE INDEX idx_playoff_matchups_state ON playoff_matchups(playoff_state_id);
CREATE INDEX idx_playoff_matchups_round ON playoff_matchups(playoff_state_id, round);