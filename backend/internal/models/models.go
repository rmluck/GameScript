package models

import "time"

type Sport struct {
	ID				int       	`json:"id"`
	Name			string    	`json:"name"`
	ShortName 		string    	`json:"short_name"`
	CreatedAt		time.Time 	`json:"created_at"`
}

type Season struct {
	ID       		int       	`json:"id"`
	SportID  		int       	`json:"sport_id"`
	StartYear		int       	`json:"start_year"`
	EndYear  		*int      	`json:"end_year"`
	IsActive		bool     	`json:"is_active"`
	CreatedAt		time.Time 	`json:"created_at"`
}

type Team struct {
	ID       		int       	`json:"id"`
	SportID  		int       	`json:"sport_id"`
	SeasonID 		int       	`json:"season_id"`
	ESPNID			string    	`json:"espn_id"`
	Abbreviation 	string    	`json:"abbreviation"`
	City     		string    	`json:"city"`
	Name	 		string    	`json:"name"`
	Conference 		*string   	`json:"conference"`
	Division		*string   	`json:"division"`
	PrimaryColor	string     	`json:"primary_color"`
	SecondaryColor	string     	`json:"secondary_color"`
	LogoURL			*string   	`json:"logo_url"`
	AlternateLogoURL *string   	`json:"alternate_logo_url"`
	CreatedAt		time.Time 	`json:"created_at"`
}

type Game struct {
	ID       		int       	`json:"id"`
	SeasonID 		int       	`json:"season_id"`
	ESPNID			string    	`json:"espn_id"`
	HomeTeamID		int       	`json:"home_team_id"`
	AwayTeamID		int       	`json:"away_team_id"`
	StartTime		time.Time 	`json:"start_time"`
	DayOfWeek		*string   	`json:"day_of_week"`
	Week      		*int      	`json:"week"`
	Location  		*string   	`json:"location"`
	Primetime		*string   	`json:"primetime"`
	Network			*string   	`json:"network"`
	HomeScore		*int      	`json:"home_score"`
	AwayScore		*int      	`json:"away_score"`
	Status   		*string   	`json:"status"`
	CreatedAt		time.Time 	`json:"created_at"`

	// Temporary fields for ESPN integration (not stored in DB)
	HomeTeamESPNID  *string   	`json:"home_team_espn_id"`
	AwayTeamESPNID  *string  	`json:"away_team_espn_id"`
}

type Scenario struct {
	ID	   			int       	`json:"id"`
	UserID 			*int      	`json:"user_id"`
	SessionToken 	*string		`json:"session_token"`
	Name			string    	`json:"name"`
	SportID 		int       	`json:"sport_id"`
	SeasonID 		int       	`json:"season_id"`
	IsPublic		bool     	`json:"is_public"`
	CreatedAt		time.Time 	`json:"created_at"`
	UpdatedAt		time.Time 	`json:"updated_at"`
}

type Pick struct {
	ID	   			int       	`json:"id"`
	ScenarioID 		int       	`json:"scenario_id"`
	GameID			int      	`json:"game_id"`
	PickedTeamID 	int      	`json:"picked_team_id"`
	PredictedHomeScore *int     `json:"predicted_home_score"`
	PredictedAwayScore *int     `json:"predicted_away_score"`
	Status			*string   	`json:"status"`
	CreatedAt		time.Time 	`json:"created_at"`
	UpdatedAt		time.Time 	`json:"updated_at"`
}

type PlayoffState struct {
	ID	   			int       	`json:"id"`
	ScenarioID 		int       	`json:"scenario_id"`
	CurrentRound	int       	`json:"current_round"`
	IsEnabled		bool     	`json:"is_enabled"`
	CreatedAt		time.Time 	`json:"created_at"`
	UpdatedAt		time.Time 	`json:"updated_at"`
}

type PlayoffSeries struct {
	ID	   			int       	`json:"id"`
	PlayoffStateID 	int       	`json:"playoff_state_id"`
	Round			int       	`json:"round"`
	SeriesOrder		int       	`json:"series_order"`
	Conference		*string   	`json:"conference"`
	HigherSeedTeamID *int     	`json:"higher_seed_team_id"`
	LowerSeedTeamID *int      	`json:"lower_seed_team_id"`
	HigherSeed		*int      	`json:"higher_seed"`
	LowerSeed		*int      	`json:"lower_seed"`
	PickedTeamID 	*int      	`json:"picked_team_id"`
	PredictedHigherSeedWins *int     `json:"predicted_higher_seed_wins"`
	PredictedLowerSeedWins *int      `json:"predicted_lower_seed_wins"`
	BestOf			*int      	`json:"best_of"`
	Status 			*string   	`json:"status"`
	CreatedAt		time.Time 	`json:"created_at"`
	UpdatedAt		time.Time 	`json:"updated_at"`
}

type PlayoffMatchup struct {
	ID	   			int       	`json:"id"`
	PlayoffStateID 	int       	`json:"playoff_state_id"`
	PlayoffSeriesID	int       	`json:"playoff_series_id"`
	Round			int       	`json:"round"`
	MatchupOrder	int       	`json:"matchup_order"`
	GameNumber		*int      	`json:"game_number"`
	Conference		*string   	`json:"conference"`
	HigherSeedTeamID *int     	`json:"higher_seed_team_id"`
	LowerSeedTeamID *int      	`json:"lower_seed_team_id"`
	HigherSeed		*int      	`json:"higher_seed"`
	LowerSeed		*int      	`json:"lower_seed"`
	PickedTeamID 	*int      	`json:"picked_team_id"`
	PredictedHigherSeedScore *int     `json:"predicted_higher_seed_score"`
	PredictedLowerSeedScore *int      `json:"predicted_lower_seed_score"`
	Status 			*string   	`json:"status"`
	CreatedAt		time.Time 	`json:"created_at"`
	UpdatedAt		time.Time 	`json:"updated_at"`
}

type User struct {
	ID 	 			int       	`json:"id"`
	Email 			string    	`json:"email"`
	Username 		string    	`json:"username"`
	PasswordHash 	string    	`json:"-"`
	IsAdmin			bool     	`json:"is_admin"`
	AvatarURL 		*string   	`json:"avatar_url"`
	FailedLoginAttempts int 	 `json:"failed_login_attempts"`
	LockedUntil		*time.Time	`json:"locked_until"`
	LastLogin		*time.Time	`json:"last_login"`
	PasswordChangedAt time.Time  `json:"password_changed_at"`
	CreatedAt		time.Time 	`json:"created_at"`
	UpdatedAt		time.Time 	`json:"updated_at"`
}