package espn

import (
	"encoding/json"
	"fmt"

	"gamescript/internal/models"
)

const nbaTeamsURL = "https://site.api.espn.com/apis/site/v2/sports/basketball/nba/teams"

func (c *Client) FetchNBATeams() ([]models.Team, error) {
	body, err := c.Get(nbaTeamsURL)
	if err != nil {
		return nil, err
	}

	var apiResp models.ESPNTeamAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	var teams []models.Team
	for _, league := range apiResp.Sports[0].Leagues {
		for _, t := range league.Teams {
			team := t.Team
			conference := team.Conference.Name
			division := team.Division.Name
			var logoURL *string
			var alternateLogoURL *string
			if len(team.Logos) > 1 {
				logoURL = &team.Logos[0].Href
				alternateLogoURL = &team.Logos[1].Href
			}
			teams = append(teams, models.Team{
				SportID:		2,
				SeasonID:		2,
				ESPNID: 	  	team.ID,
				Abbreviation: 	team.Abbreviation,
				City:		 	team.Location,
				Name: 	   		team.Name,
				Conference: 	&conference,
				Division:  		&division,
				PrimaryColor:  	team.PrimaryColor,
				SecondaryColor:	team.SecondaryColor,
				LogoURL:		logoURL,
				AlternateLogoURL: alternateLogoURL,
			})
		}
	}
	
	return teams, nil
}