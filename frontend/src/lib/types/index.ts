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

// export interface TeamRecord {
//     rank: number;
//     team_id: number;
//     team_name: string;
//     team_abbr: string;
//     wins: number;
//     losses: number;
//     ties: number;
//     win_pct: number;
//     division_record: string;
//     conference_record: string;
//     points_for: number;
//     points_against: number;
//     point_diff: number;
//     division_games_back?: number;
//     conference_games_back?: number;
//     logo_url: string;
//     team_primary_color: string;
//     team_secondary_color: string;
// }

export interface PlayoffSeed {
    seed: number;
    team_id: number;
    team_name: string;
    team_abbr: string;
    wins: number;
    losses: number;
    ties: number;
    win_pct: number;
    is_division_winner: boolean;
    logo_url: string;
    team_primary_color: string;
    team_secondary_color: string;
    conference_record: string;
    division_record: string;
    conference_games_back: number;
    division_games_back: number;
    points_for: number;
    points_against: number;
    point_diff: number;
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