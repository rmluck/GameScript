package models

import "time"

type Sport struct {
	ID			int       `json:"id"`
	Name		string    `json:"name"`
	ShortName 	string    `json:"short_name"`
	CreatedAt	time.Time `json:"created_at"`
}

type Season struct {
	ID       	int       `json:"id"`
	SportID  	int       `json:"sport_id"`
	StartYear	int       `json:"start_year"`
	EndYear  	*int       `json:"end_year"`
	IsActive	bool      `json:"is_active"`
	CreatedAt	time.Time `json:"created_at"`
}

type Team struct {
	ID       	int       `json:"id"`
	SportID  	int       `json:"sport_id"`
	SeasonID 	int       `json:"season_id"`
	Abbreviation string    `json:"abbreviation"`
	City     	string    `json:"city"`
	Name	 	string    `json:"name"`
	Conference 	*string   `json:"conference"`
	Division	*string   `json:"division"`
	PrimaryColor	*string   `json:"primary_color"`
	SecondaryColor	*string   `json:"secondary_color"`
	LogoURL		*string   `json:"logo_url"`
	CreatedAt	time.Time `json:"created_at"`
}

type Game struct {
	ID       	int       `json:"id"`
	SeasonID 	int       `json:"season_id"`
	HomeTeamID	int       `json:"home_team_id"`
	AwayTeamID	int       `json:"away_team_id"`
	StartTime	time.Time `json:"start_time"`
	DayOfWeek	*string   `json:"day_of_week"`
	Week      	*int      `json:"week"`
	Location  	*string   `json:"location"`
	Primetime	*string    `json:"primetime"`
	HomeScore	*int      `json:"home_score"`
	AwayScore	*int      `json:"away_score"`
	Status   	*string   `json:"status"`
	IsPostseason	bool      `json:"is_postseason"`
	CreatedAt	time.Time `json:"created_at"`
}