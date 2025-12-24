export interface User {
    id: number;
    email: string;
    username: string;
    is_admin: boolean;
    avatar_url?: string;
    failed_login_attempts: number;
    locked_until?: string;
    last_login?: string;
    password_changed_at: string;
    created_at: string;
    updated_at: string;
}

export interface AuthResponse {
    user: User;
    token: string;
}

export interface Sport {
    id: number;
    name: string;
    short_name: string;
    created_at: string;
}

export interface Season {
    id: number;
    sport_id: number;
    start_year: number;
    end_year?: number;
    is_active: boolean;
    created_at: string;
}

export interface Team {
    id: number;
    sport_id: number;
    season_id: number;
    espn_id: string;
    abbreviation: string;
    city: string;
    name: string;
    conference: string;
    division: string;
    primary_color: string;
    secondary_color?: string;
    logo_url?: string;
    alternate_logo_url?: string;
}

export interface Game {
    id: number;
    season_id: number;
    espn_id: string;
    home_team_id: number;
    away_team_id: number;
    start_time: string;
    day_of_week: string;
    week: number;
    location?: string;
    primetime?: string;
    network?: string;
    home_score?: number;
    away_score?: number;
    status: 'upcoming' | 'in_progress' | 'final';
    is_postseason: boolean;
    home_team: Team;
    away_team: Team;
}

export interface Scenario {
    id: number;
    name: string;
    sport_id: number;
    season_id: number;
    season_start_year?: number;
    season_end_year?: number;
    is_public: boolean;
    sport_short_name?: string;
    created_at: string;
    updated_at: string;
}

export interface Pick {
    id: number;
    scenario_id: number;
    game_id: number;
    picked_team_id: number | null;
    predicted_home_score?: number;
    predicted_away_score?: number;
    status: 'pending' | 'correct' | 'incorrect';
    game?: Game;
    created_at: string;
    updated_at: string;
}

export interface PlayoffSeed {
    seed: number;
    team_id: number;
    team_name: string;
    team_city: string;
    team_abbr: string;
    wins: number;
    losses: number;
    ties: number;
    win_pct: number;
    is_division_winner: boolean;
    logo_url: string;
    team_primary_color: string;
    team_secondary_color: string;
    conference_wins: number;
    conference_losses: number;
    conference_ties: number;
    division_wins: number;
    division_losses: number;
    division_ties: number;
    conference_games_back: number;
    division_games_back: number;
    points_for: number;
    points_against: number;
    point_diff: number;
    home_wins: number;
    home_losses: number;
    home_ties: number;
    away_wins: number;
    away_losses: number;
    away_ties: number;
    strength_of_schedule: number;
    strength_of_victory: number;
}

export interface ConferenceStandings {
    divisions: Record<string, PlayoffSeed[]>; 
    playoff_seeds: PlayoffSeed[];
}

export interface Standings {
    afc: ConferenceStandings;
    nfc: ConferenceStandings;
    draft_order: DraftPick[];
}

export interface DraftPick {
    pick: number;
    team_id: number;
    team_name: string;
    team_abbr: string;
    record: string;
    reason: string;
    logo_url: string;
    team_primary_color: string;
    team_secondary_color: string;
}

export interface PlayoffState {
    id: number;
    scenario_id: number;
    current_round: number;
    is_enabled: boolean;
    created_at: string;
    updated_at: string;
}

export interface PlayoffMatchup {
    id: number;
    playoff_state_id: number;
    round: number;
    matchup_order: number;
    conference?: string;
    higher_seed_team_id: number;
    lower_seed_team_id: number;
    higher_seed: number;
    lower_seed: number;
    picked_team_id?: number;
    predicted_higher_seed_score?: number;
    predicted_lower_seed_score?: number;
    status: 'pending' | 'completed';
    created_at: string;
    updated_at: string;
    higher_seed_team?: Team;
    lower_seed_team?: Team;
}

export const PLAYOFF_ROUNDS = {
    WILD_CARD: 1,
    DIVISIONAL: 2,
    CONFERENCE: 3,
    SUPER_BOWL: 4,
} as const;

export const PLAYOFF_ROUND_NAMES: Record<number, string> = {
    1: 'Wild Card',
    2: 'Divisional',
    3: 'Conference Championship',
    4: 'Super Bowl',
};