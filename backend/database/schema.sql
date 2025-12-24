-- SPORTS
CREATE TABLE sports (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    short_name VARCHAR(10) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- SEASONS
CREATE TABLE seasons (
    id SERIAL PRIMARY KEY,
    sport_id INTEGER NOT NULL REFERENCES sports(id) ON DELETE CASCADE,
    start_year INTEGER NOT NULL,
    end_year INTEGER,
    is_active BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- TEAMS
CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    sport_id INTEGER NOT NULL REFERENCES sports(id) ON DELETE CASCADE,
    season_id INTEGER NOT NULL REFERENCES seasons(id) ON DELETE CASCADE,
    espn_id VARCHAR(16),
    abbreviation VARCHAR(10) NOT NULL,
    city VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    conference VARCHAR(50),
    division VARCHAR(50),
    primary_color VARCHAR(20),
    secondary_color VARCHAR(20),
    logo_url VARCHAR(255),
    alternate_logo_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(season_id, espn_id)
);

-- GAMES
CREATE TABLE games (
    id SERIAL PRIMARY KEY,
    season_id INTEGER NOT NULL REFERENCES seasons(id) ON DELETE CASCADE,
    espn_id VARCHAR(32),
    home_team_id INTEGER NOT NULL REFERENCES teams(id),
    away_team_id INTEGER NOT NULL REFERENCES teams(id),
    start_time TIMESTAMP NOT NULL,
    day_of_week VARCHAR(20),
    week INTEGER,
    location VARCHAR(100),
    primetime VARCHAR(100),
    network VARCHAR(100),
    home_score INTEGER,
    away_score INTEGER,
    status VARCHAR(50) DEFAULT 'upcoming',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(season_id, espn_id)
);

-- USERS
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    is_admin BOOLEAN DEFAULT FALSE,
    avatar_url VARCHAR(255),
    -- Need to update database to add these fields
    failed_login_attempts INTEGER DEFAULT 0,
    locked_until TIMESTAMPTZ,
    last_login TIMESTAMPTZ,
    password_changed_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- SCENARIOS
CREATE TABLE scenarios (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    sport_id INTEGER NOT NULL REFERENCES sports(id) ON DELETE CASCADE,
    season_id INTEGER NOT NULL REFERENCES seasons(id) ON DELETE CASCADE,
    is_public BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    session_token VARCHAR(255)
);

-- PICKS
CREATE TABLE picks (
    id SERIAL PRIMARY KEY,
    scenario_id INTEGER NOT NULL REFERENCES scenarios(id) ON DELETE CASCADE,
    game_id INTEGER NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    picked_team_id INTEGER,
    predicted_home_score INTEGER,
    predicted_away_score INTEGER,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(scenario_id, game_id)
);

-- PLAYOFF STATES
CREATE TABLE playoff_states (
    id SERIAL PRIMARY KEY,
    scenario_id INTEGER NOT NULL REFERENCES scenarios(id) ON DELETE CASCADE,
    current_round INTEGER DEFAULT 0,
    is_enabled BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(scenario_id)
);

-- PLAYOFF MATCHUPS
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

-- Indexes for performance optimization
CREATE INDEX idx_games_season ON games(season_id);
CREATE INDEX idx_teams_sport ON teams(sport_id);
CREATE INDEX idx_scenarios_user_id ON scenarios(user_id);
CREATE INDEX idx_picks_scenario ON picks(scenario_id);
CREATE INDEX idx_picks_game ON picks(game_id);
CREATE INDEX idx_scenarios_session_token ON scenarios(session_token);
CREATE INDEX idx_playoff_states_scenario ON playoff_states(scenario_id);
CREATE INDEX idx_playoff_matchups_state ON playoff_matchups(playoff_state_id);
CREATE INDEX idx_playoff_matchups_round ON playoff_matchups(playoff_state_id, round);