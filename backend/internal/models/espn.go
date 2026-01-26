// Models for ESPN API responses

package models

type ESPNTeamAPIResponse struct {
	Sports []struct {
		Leagues []struct {
			Teams []struct {
				Team struct {
					ID          	string `json:"id"`
					Abbreviation 	string `json:"abbreviation"`
					DisplayName  	string `json:"displayName"`
					Location     	string `json:"location"`
					Name         	string `json:"name"`
					Conference struct {
						Name 		string `json:"name"`
					} `json:"conference"`
					Division struct {
						Name 		string `json:"name"`
					} `json:"division"`
					PrimaryColor  	string `json:"color"`
					SecondaryColor 	string `json:"alternateColor"`
					Logos []struct {
						Href 		string `json:"href"`
						Rel 		[]string `json:"rel"`
					} `json:"logos"`
				} `json:"team"`
			} `json:"teams"`
		} `json:"leagues"`
	} `json:"sports"`
}

type ESPNScheduleAPIResponse struct {
	Week struct {
		Number int `json:"number"`
	}
	Events []struct {
		Competitions []struct {
			ID string `json:"id"`
			Date string `json:"date"`
			Type struct {
				ID string `json:"id"`
			} `json:"type"`
			Venue struct {
				FullName string `json:"fullName"`
				Address struct {
					City string `json:"city"`
					State string `json:"state"`
					Country string `json:"country"`
				} `json:"address"`
			} `json:"venue"`
			Competitors []struct {
				HomeAway string `json:"homeAway"`
				Winner bool `json:"winner"`
				Team struct {
					ID string `json:"id"`
				} `json:"team"`
				Score string `json:"score"`
			} `json:"competitors"`
			Status struct {
				Type struct {
					Name string `json:"name"`
				} `json:"type"`
			} `json:"status"`
			Broadcasts []struct {
				Names  []string `json:"names"`
			} `json:"broadcasts"`
		} `json:"competitions"`
	} `json:"events"`
}