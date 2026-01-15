package standings

import (
	"testing"
)

func TestCalculateWinPct(t *testing.T) {
	tests := []struct {
		name     string
		wins     int
		losses   int
		ties     int
		expected float64
	}{
		{"Perfect season", 17, 0, 0, 1.0},
		{"Winless season", 0, 17, 0, 0.0},
		{"With ties", 10, 6, 1, 0.617647}, // (10 + 0.5) / 17
		{"All ties", 0, 0, 17, 0.5},
		{"No games", 0, 0, 0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateNFLWinPct(tt.wins, tt.losses, tt.ties)
			if result != tt.expected {
				t.Errorf("calculateWinPct(%d, %d, %d) = %f; want %f",
					tt.wins, tt.losses, tt.ties, result, tt.expected)
			}
		})
	}
}

func TestCalculateGamesBack(t *testing.T) {
	tests := []struct {
		name     string
		leader   NFLTeamRecord
		team     NFLTeamRecord
		expected float64
	}{
		{
			name:     "Same team",
			leader:   NFLTeamRecord{TeamID: 1, Wins: 10, Losses: 5, Ties: 0},
			team:     NFLTeamRecord{TeamID: 1, Wins: 10, Losses: 5, Ties: 0},
			expected: 0.0,
		},
		{
			name:     "1 game back",
			leader:   NFLTeamRecord{TeamID: 1, Wins: 10, Losses: 5, Ties: 0},
			team:     NFLTeamRecord{TeamID: 2, Wins: 9, Losses: 6, Ties: 0},
			expected: 1.0,
		},
		{
			name:     "2.5 games back",
			leader:   NFLTeamRecord{TeamID: 1, Wins: 10, Losses: 5, Ties: 0},
			team:     NFLTeamRecord{TeamID: 2, Wins: 8, Losses: 7, Ties: 0},
			expected: 2.0,
		},
		{
			name:     "With ties",
			leader:   NFLTeamRecord{TeamID: 1, Wins: 10, Losses: 5, Ties: 2},
			team:     NFLTeamRecord{TeamID: 2, Wins: 9, Losses: 6, Ties: 2},
			expected: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateNFLGamesBack(tt.leader, tt.team)
			if result != tt.expected {
				t.Errorf("calculateGamesBack() = %f; want %f", result, tt.expected)
			}
		})
	}
}

func TestGetHeadToHeadWinner(t *testing.T) {
	teamA := NFLTeamRecord{TeamID: 1}
	teamB := NFLTeamRecord{TeamID: 2}

	tests := []struct {
		name     string
		games    []NFLGameResult
		expected []NFLTeamRecord
	}{
		{
			name: "Team A wins head-to-head",
			games: []NFLGameResult{
				{HomeTeamID: 1, AwayTeamID: 2, HomeScore: 24, AwayScore: 17},
			},
			expected: []NFLTeamRecord{teamA},
		},
		{
			name: "Team B wins head-to-head",
			games: []NFLGameResult{
				{HomeTeamID: 2, AwayTeamID: 1, HomeScore: 30, AwayScore: 20},
			},
			expected: []NFLTeamRecord{teamB},
		},
		{
			name: "Split series (tie)",
			games: []NFLGameResult{
				{HomeTeamID: 1, AwayTeamID: 2, HomeScore: 24, AwayScore: 17},
				{HomeTeamID: 2, AwayTeamID: 1, HomeScore: 30, AwayScore: 20},
			},
			expected: []NFLTeamRecord{},
		},
		{
			name:     "No games played",
			games:    []NFLGameResult{},
			expected: []NFLTeamRecord{},
		},
		{
			name: "Tie game doesn't count",
			games: []NFLGameResult{
				{HomeTeamID: 1, AwayTeamID: 2, HomeScore: 24, AwayScore: 24},
			},
			expected: []NFLTeamRecord{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			teams := []NFLTeamRecord{teamA, teamB}
			result := compareNFLHeadToHead(teams, tt.games)
			if !equalTeamRecordSlices(result, tt.expected) {
				t.Errorf("compareHeadToHead() = %v; want %v", result, tt.expected)
			}
		})
	}
}

func equalTeamRecordSlices(a, b []NFLTeamRecord) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestFilterByConference(t *testing.T) {
	teams := []NFLTeamRecord{
		{TeamID: 1, Conference: "AFC"},
		{TeamID: 2, Conference: "NFC"},
		{TeamID: 3, Conference: "AFC"},
		{TeamID: 4, Conference: "NFC"},
	}

	afcTeams := filterByNFLConference(teams, "AFC")
	if len(afcTeams) != 2 {
		t.Errorf("Expected 2 AFC teams, got %d", len(afcTeams))
	}

	nfcTeams := filterByNFLConference(teams, "NFC")
	if len(nfcTeams) != 2 {
		t.Errorf("Expected 2 NFC teams, got %d", len(nfcTeams))
	}
}

func TestCalculateTeamRecords(t *testing.T) {
	teams := []NFLTeamRecord{
		{TeamID: 1, Conference: "AFC", Division: "AFC East"},
		{TeamID: 2, Conference: "AFC", Division: "AFC East"},
	}

	games := []NFLGameResult{
		{HomeTeamID: 1, AwayTeamID: 2, HomeScore: 24, AwayScore: 17},
		{HomeTeamID: 2, AwayTeamID: 1, HomeScore: 30, AwayScore: 27},
	}

	records := calculateNFLTeamRecords(teams, games)

	// Team 1: 1-1, scored 51, allowed 47
	if records[0].Wins != 1 || records[0].Losses != 1 {
		t.Errorf("Team 1: Expected 1-1, got %d-%d", records[0].Wins, records[0].Losses)
	}
	if records[0].PointsFor != 51 {
		t.Errorf("Team 1: Expected 51 points for, got %d", records[0].PointsFor)
	}
	if records[0].PointsAgainst != 47 {
		t.Errorf("Team 1: Expected 47 points against, got %d", records[0].PointsAgainst)
	}

	// Team 2: 1-1, scored 47, allowed 51
	if records[1].Wins != 1 || records[1].Losses != 1 {
		t.Errorf("Team 2: Expected 1-1, got %d-%d", records[1].Wins, records[1].Losses)
	}

	// Both are division games
	if records[0].DivisionWins != 1 || records[0].DivisionLosses != 1 {
		t.Errorf("Team 1: Expected 1-1 division record, got %d-%d",
			records[0].DivisionWins, records[0].DivisionLosses)
	}
}

func TestApplyTiebreakers(t *testing.T) {
	// Test that higher win percentage ranks first
	teams := []NFLTeamRecord{
		{TeamID: 1, Wins: 8, Losses: 9, WinPct: 0.470},
		{TeamID: 2, Wins: 10, Losses: 7, WinPct: 0.588},
		{TeamID: 3, Wins: 12, Losses: 5, WinPct: 0.706},
	}

	sorted := applyNFLDivisionTiebreakers(teams, []NFLGameResult{})

	if sorted[0].TeamID != 3 || sorted[1].TeamID != 2 || sorted[2].TeamID != 1 {
		t.Errorf("Teams not sorted correctly by win percentage")
	}
}

func TestApplyTiebreakersWithPointDifferential(t *testing.T) {
	// Test point differential tiebreaker when records are identical
	teams := []NFLTeamRecord{
		{TeamID: 1, Wins: 10, Losses: 7, WinPct: 0.588, PointsFor: 350, PointsAgainst: 320},
		{TeamID: 2, Wins: 10, Losses: 7, WinPct: 0.588, PointsFor: 380, PointsAgainst: 340},
	}

	sorted := applyNFLDivisionTiebreakers(teams, []NFLGameResult{})

	// Team 2 should rank higher (better point differential: +40 vs +30)
	if sorted[0].TeamID != 2 {
		t.Errorf("Expected Team 2 to rank first due to point differential")
	}
}
