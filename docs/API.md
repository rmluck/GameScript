# GameScript API Documentation

**Base URL:** `http://localhost:8080/api`

**Version:** 1.0.0

---

## Table of Contents

1. [Authentication](#authentication)
2. [Sports & Seasons](#sports--seasons)
3. [Teams](#teams)
4. [Games](#games)
5. [Scenarios](#scenarios)
6. [Picks](#picks)
7. [Standings](#standings)
8. [Playoffs](#playoffs)
9. [Admin](#admin)
10. [Error Handling](#error-handling)

---

## Authentication

### Register User
**POST** `/auth/register`

Creates a new user account with password validation.

**Request Body:**
```json
{
  "email": "user@example.com",
  "username": "username",
  "password": "SecurePass123!"
}
```

**Password Requirements:**
- Minimum 8 characters
- At least one uppercase letter
- At least one lowercase letter
- At least one number
- At least one special character

**Username Requirements:**
- 3-50 characters
- Letters, numbers, hyphens, and underscores only

**Response (201 Created):**
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "username": "username",
    "is_admin": false,
    "created_at": "2025-01-01T00:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Errors:**
- `400` - Validation errors (weak password, invalid email, etc.)
- `400` - Email or username already exists

---

### Login User
**POST** `/auth/login`

Authenticates a user and returns a JWT token.

**Security Features:**
- Rate limited (5 attempts per 15 minutes)
- Account lockout after 5 failed attempts (15 minutes)
- Failed attempt tracking

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

**Response (200 OK):**
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "username": "username",
    "is_admin": false,
    "created_at": "2025-01-01T00:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Errors:**
- `400` - Missing email or password
- `401` - Invalid credentials
- `423` - Account locked (too many failed attempts)

---

### Get Current User
**GET** `/auth/me`

Returns the currently authenticated user's information.

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "id": 1,
  "email": "user@example.com",
  "username": "username",
  "is_admin": false,
  "avatar_url": null,
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

**Errors:**
- `401` - Missing or invalid token
- `404` - User not found

---

### Update Profile
**PUT** `/auth/profile`

Updates user profile information (username, email, password).

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "username": "new_username",
  "email": "newemail@example.com",
  "current_password": "OldPass123!",
  "new_password": "NewSecurePass456!"
}
```

**Notes:**
- All fields are optional
- `current_password` required if changing password
- New password must meet password requirements

**Response (200 OK):**
```json
{
  "id": 1,
  "email": "newemail@example.com",
  "username": "new_username",
  "is_admin": false,
  "avatar_url": null,
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-02T00:00:00Z"
}
```

**Errors:**
- `400` - Validation errors
- `401` - Current password incorrect
- `500` - Update failed

---

## Sports & Seasons

### Get All Sports
**GET** `/sports`

Returns a list of all available sports.

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "name": "National Football League",
    "short_name": "NFL",
    "created_at": "2025-01-01T00:00:00Z"
  },
  {
    "id": 2,
    "name": "National Basketball Association",
    "short_name": "NBA",
    "created_at": "2025-01-01T00:00:00Z"
  }
]
```

---

### Get Seasons for a Sport
**GET** `/sports/:sport_id/seasons`

Returns all seasons for a specific sport.

**Parameters:**
- `sport_id` (path) - Sport ID (1=NFL, 2=NBA)

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "sport_id": 1,
    "start_year": 2024,
    "end_year": 2025,
    "is_active": true,
    "created_at": "2024-09-01T00:00:00Z"
  }
]
```

---

### Get Season Details
**GET** `/seasons/:season_id`

Returns details for a specific season.

**Parameters:**
- `season_id` (path) - Season ID

**Response (200 OK):**
```json
{
  "id": 1,
  "sport_id": 1,
  "start_year": 2024,
  "end_year": 2025,
  "is_active": true,
  "created_at": "2024-09-01T00:00:00Z"
}
```

**Errors:**
- `404` - Season not found

---

## Teams

### Get Teams for a Season
**GET** `/seasons/:season_id/teams`

Returns all teams for a specific season.

**Parameters:**
- `season_id` (path) - Season ID

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "sport_id": 1,
    "season_id": 1,
    "espn_id": "22",
    "abbreviation": "ARI",
    "city": "Arizona",
    "name": "Cardinals",
    "conference": "NFC",
    "division": "NFC West",
    "primary_color": "a40227",
    "secondary_color": "face07",
    "logo_url": "https://a.espncdn.com/...",
    "alternate_logo_url": "https://a.espncdn.com/..."
  }
]
```

---

### Get Single Team
**GET** `/teams/:team_id`

Returns details for a specific team.

**Parameters:**
- `team_id` (path) - Team ID

**Response (200 OK):**
```json
{
  "id": 1,
  "sport_id": 1,
  "season_id": 1,
  "espn_id": "22",
  "abbreviation": "ARI",
  "city": "Arizona",
  "name": "Cardinals",
  "conference": "NFC",
  "division": "NFC West",
  "primary_color": "a40227",
  "secondary_color": "face07",
  "logo_url": "https://...",
  "alternate_logo_url": "https://..."
}
```

**Errors:**
- `404` - Team not found

---

## Games

### Get Games for a Season
**GET** `/seasons/:season_id/games`

Returns all games for a specific season.

**Parameters:**
- `season_id` (path) - Season ID

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "season_id": 1,
    "espn_id": "401671745",
    "home_team_id": 12,
    "away_team_id": 2,
    "start_time": "2024-09-05T20:20:00Z",
    "day_of_week": "Thursday",
    "week": 1,
    "location": "Arrowhead Stadium, Kansas City, MO, USA",
    "primetime": "TNF",
    "network": "NBC",
    "home_score": 27,
    "away_score": 20,
    "status": "final",
    "home_team": {
      "id": 12,
      "abbreviation": "KC",
      "city": "Kansas City",
      "name": "Chiefs",
      "conference": "AFC",
      "division": "AFC West",
      "primary_color": "e31837",
      "secondary_color": "ffb612",
      "logo_url": "https://...",
      "alternate_logo_url": "https://..."
    },
    "away_team": {
      "id": 2,
      "abbreviation": "BAL",
      "city": "Baltimore",
      "name": "Ravens",
      "conference": "AFC",
      "division": "AFC North",
      "primary_color": "241773",
      "secondary_color": "000000",
      "logo_url": "https://...",
      "alternate_logo_url": "https://..."
    }
  }
]
```

**Game Status Values:**
- `"upcoming"` - Game hasn't started
- `"final"` - Game completed

---

### Get Games by Week
**GET** `/seasons/:season_id/weeks/:week/games`

Returns all games for a specific week in a season.

**Parameters:**
- `season_id` (path) - Season ID
- `week` (path) - Week number

**Response:** Same format as Get Games for a Season

---

### Get Games for a Team
**GET** `/teams/:team_id/games`

Returns all games for a specific team.

**Parameters:**
- `team_id` (path) - Team ID

**Response:** Same format as Get Games for a Season

---

### Get Single Game
**GET** `/games/:game_id`

Returns details for a specific game.

**Parameters:**
- `game_id` (path) - Game ID

**Response:** Same format as individual game in Get Games for a Season

**Errors:**
- `404` - Game not found

---

## Scenarios

### Get All Scenarios
**GET** `/scenarios`

Returns all scenarios for the current user (authenticated) or session (guest).

**Headers (Optional):**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "name": "My NFL Playoff Scenario",
    "sport_id": 1,
    "season_id": 1,
    "season_start_year": 2024,
    "season_end_year": 2025,
    "is_public": true,
    "sport_short_name": "NFL",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-15T12:30:00Z"
  }
]
```

**Notes:**
- Guest users get scenarios tied to their session token
- Authenticated users get their owned scenarios
- Scenarios sorted by `updated_at` DESC

---

### Get Single Scenario
**GET** `/scenarios/:scenario_id`

Returns details for a specific scenario.

**Headers (Optional):**
```
Authorization: Bearer <token>
```

**Parameters:**
- `scenario_id` (path) - Scenario ID

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "My NFL Playoff Scenario",
  "sport_id": 1,
  "season_id": 1,
  "season_start_year": 2024,
  "season_end_year": 2025,
  "is_public": true,
  "sport_short_name": "NFL",
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-15T12:30:00Z"
}
```

**Errors:**
- `403` - Unauthorized (not owner)
- `404` - Scenario not found

---

### Create Scenario
**POST** `/scenarios`

Creates a new scenario.

**Headers (Optional):**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "name": "My New Scenario",
  "sport_id": 1,
  "season_id": 1,
  "is_public": true
}
```

**Response (201 Created):**
```json
{
  "id": 1,
  "name": "My New Scenario",
  "sport_id": 1,
  "season_id": 1,
  "season_start_year": 2024,
  "season_end_year": 2025,
  "is_public": true,
  "sport_short_name": "NFL",
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

**Notes:**
- Guest users automatically get a session token cookie
- Session tokens valid for 7 days

**Errors:**
- `400` - Missing required fields

---

### Update Scenario
**PUT** `/scenarios/:scenario_id`

Updates an existing scenario (name and/or public status).

**Headers (Optional):**
```
Authorization: Bearer <token>
```

**Parameters:**
- `scenario_id` (path) - Scenario ID

**Request Body:**
```json
{
  "name": "Updated Scenario Name",
  "is_public": false
}
```

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "Updated Scenario Name",
  "sport_id": 1,
  "season_id": 1,
  "is_public": false,
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-02T10:00:00Z"
}
```

**Errors:**
- `400` - No fields to update
- `403` - Unauthorized (not owner)
- `404` - Scenario not found

---

### Delete Scenario
**DELETE** `/scenarios/:scenario_id`

Deletes a scenario and all associated picks.

**Headers (Optional):**
```
Authorization: Bearer <token>
```

**Parameters:**
- `scenario_id` (path) - Scenario ID

**Response (200 OK):**
```json
{
  "message": "Scenario deleted successfully"
}
```

**Errors:**
- `403` - Unauthorized (not owner)
- `404` - Scenario not found

---

### Claim Guest Scenario
**POST** `/scenarios/:scenario_id/claim`

Claims a guest scenario after registering/logging in.

**Headers:**
```
Authorization: Bearer <token>
```

**Parameters:**
- `scenario_id` (path) - Scenario ID

**Response (200 OK):**
```json
{
  "message": "Scenario claimed successfully",
  "id": 1
}
```

**Errors:**
- `400` - No session token found
- `404` - Scenario not found or already claimed

---

## Picks

### Get All Picks for Scenario
**GET** `/picks/scenarios/:scenario_id`

Returns all picks for a specific scenario.

**Headers (Optional):**
```
Authorization: Bearer <token>
```

**Parameters:**
- `scenario_id` (path) - Scenario ID

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "scenario_id": 1,
    "game_id": 1,
    "picked_team_id": 12,
    "predicted_home_score": 27,
    "predicted_away_score": 20,
    "status": "pending",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z",
    "game": {
      "espn_id": "401671745",
      "start_time": "2024-09-05T20:20:00Z",
      "week": 1,
      "home_score": null,
      "away_score": null,
      "status": "upcoming",
      "home_team": { ... },
      "away_team": { ... }
    }
  }
]
```

**Errors:**
- `403` - Unauthorized (not owner)

---

### Get Single Pick
**GET** `/picks/scenarios/:scenario_id/games/:game_id`

Returns a specific pick.

**Headers (Optional):**
```
Authorization: Bearer <token>
```

**Parameters:**
- `scenario_id` (path) - Scenario ID
- `game_id` (path) - Game ID

**Response (200 OK):**
```json
{
  "id": 1,
  "scenario_id": 1,
  "game_id": 1,
  "picked_team_id": 12,
  "predicted_home_score": 27,
  "predicted_away_score": 20,
  "status": "pending",
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

**Errors:**
- `403` - Unauthorized (not owner)
- `404` - Pick not found

---

### Create Pick
**POST** `/picks/scenarios/:scenario_id/games/:game_id`

Creates a new pick for a game.

**Headers (Optional):**
```
Authorization: Bearer <token>
```

**Parameters:**
- `scenario_id` (path) - Scenario ID
- `game_id` (path) - Game ID

**Request Body:**
```json
{
  "picked_team_id": 12,
  "predicted_home_score": 27,
  "predicted_away_score": 20
}
```

**Notes:**
- If both scores provided, `picked_team_id` automatically set to winner
- Scores are optional (can pick winner without scores)
- Updates scenario's `updated_at` timestamp

**Response (201 Created):**
```json
{
  "id": 1,
  "scenario_id": 1,
  "game_id": 1,
  "picked_team_id": 12,
  "predicted_home_score": 27,
  "predicted_away_score": 20,
  "status": "pending",
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

**Errors:**
- `400` - Invalid request body
- `403` - Unauthorized (not owner)
- `500` - Database error (duplicate pick)

---

### Update Pick
**PUT** `/picks/scenarios/:scenario_id/games/:game_id`

Updates an existing pick.

**Headers (Optional):**
```
Authorization: Bearer <token>
```

**Parameters:**
- `scenario_id` (path) - Scenario ID
- `game_id` (path) - Game ID

**Request Body:**
```json
{
  "picked_team_id": 2,
  "predicted_home_score": 24,
  "predicted_away_score": 27
}
```

**Notes:**
- Deletes any playoff brackets if they exist
- Updates scenario's `updated_at` timestamp

**Response (200 OK):**
```json
{
  "id": 1,
  "scenario_id": 1,
  "game_id": 1,
  "picked_team_id": 2,
  "predicted_home_score": 24,
  "predicted_away_score": 27,
  "status": "pending",
  "updated_at": "2025-01-02T00:00:00Z"
}
```

**Errors:**
- `400` - Invalid request body
- `403` - Unauthorized (not owner)
- `404` - Pick not found
- `500` - Database error

---

### Delete Pick
**DELETE** `/picks/scenarios/:scenario_id/games/:game_id`

Deletes a pick.

**Headers (Optional):**
```
Authorization: Bearer <token>
```

**Parameters:**
- `scenario_id` (path) - Scenario ID
- `game_id` (path) - Game ID

**Response (200 OK):**
```json
{
  "message": "Pick deleted successfully"
}
```

**Errors:**
- `403` - Unauthorized (not owner)
- `500` - Database error

---

## Standings

### Get Standings for Scenario
**GET** `/scenarios/:scenario_id/standings`

Calculates and returns standings based on actual results + user picks.

**Parameters:**
- `scenario_id` (path) - Scenario ID

**NFL Response (200 OK):**
```json
{
  "afc": {
    "divisions": {
      "AFC East": [
        {
          "seed": 2,
          "team_id": 2,
          "team_name": "Bills",
          "team_city": "Buffalo",
          "team_abbr": "BUF",
          "wins": 13,
          "losses": 4,
          "ties": 0,
          "win_pct": 0.765,
          "home_wins": 7,
          "home_losses": 1,
          "home_ties": 0,
          "away_wins": 6,
          "away_losses": 3,
          "away_ties": 0,
          "division_wins": 5,
          "division_losses": 1,
          "division_ties": 0,
          "conference_wins": 10,
          "conference_losses": 2,
          "conference_ties": 0,
          "division_games_back": 0.0,
          "conference_games_back": 1.0,
          "points_for": 482,
          "points_against": 345,
          "point_diff": 137,
          "strength_of_schedule": 0.524,
          "strength_of_victory": 0.612,
          "is_division_winner": true,
          "logo_url": "https://...",
          "team_primary_color": "00338d",
          "team_secondary_color": "c60c30"
        }
      ]
    },
    "playoff_seeds": [
      {
        "seed": 1,
        "team_id": 12,
        "team_name": "Chiefs",
        "team_city": "Kansas City",
        "team_abbr": "KC",
        "wins": 15,
        "losses": 2,
        "ties": 0,
        "win_pct": 0.882,
        "home_wins": 8,
        "home_losses": 0,
        "home_ties": 0,
        "away_wins": 7,
        "away_losses": 2,
        "away_ties": 0,
        "division_wins": 6,
        "division_losses": 0,
        "division_ties": 0,
        "conference_wins": 11,
        "conference_losses": 1,
        "conference_ties": 0,
        "division_games_back": 0.0,
        "conference_games_back": 0.0,
        "points_for": 456,
        "points_against": 312,
        "point_diff": 144,
        "strength_of_schedule": 0.498,
        "strength_of_victory": 0.589,
        "is_division_winner": true,
        "logo_url": "https://...",
        "team_primary_color": "e31837",
        "team_secondary_color": "ffb612"
      }
    ]
  },
  "nfc": { ... },
  "draft_order": [
    {
      "pick": 1,
      "team_id": 5,
      "team_name": "Panthers",
      "team_abbr": "CAR",
      "record": "3-14-0",
      "logo_url": "https://...",
      "team_primary_color": "0085ca",
      "team_secondary_color": "101820"
    }
  ]
}
```

**NBA Response (200 OK):**
```json
{
  "eastern": {
    "divisions": {
      "Atlantic": [
        {
          "seed": 1,
          "team_id": 3,
          "team_name": "Celtics",
          "team_city": "Boston",
          "team_abbr": "BOS",
          "wins": 35,
          "losses": 15,
          "win_pct": 0.700,
          "home_wins": 20,
          "home_losses": 5,
          "away_wins": 15,
          "away_losses": 10,
          "division_wins": 10,
          "division_losses": 2,
          "conference_wins": 24,
          "conference_losses": 8,
          "division_games_back": 0.0,
          "conference_games_back": 0.0,
          "points_for": 5650,
          "points_against": 5200,
          "games_with_scores": 50,
          "strength_of_schedule": 0.512,
          "strength_of_victory": 0.598,
          "is_division_winner": true,
          "logo_url": "https://...",
          "team_primary_color": "007a33",
          "team_secondary_color": "ba9653"
        }
      ]
    },
    "playoff_seeds": [ ... ]
  },
  "western": { ... },
  "draft_order": [ ... ]
}
```

**NFL Tiebreaker Rules (in order):**
1. Win percentage
2. Head-to-head record
3. Division record (division tiebreakers)
4. Conference record
5. Common games (minimum 4)
6. Strength of victory
7. Strength of schedule
8. Point differential
9. Points scored
10. Points allowed

**NBA Tiebreaker Rules (in order):**
1. Win percentage
2. Head-to-head record
3. Division winner (if applicable)
4. Division win percentage (same division)
5. Conference win percentage
6. Point differential

**Errors:**
- `400` - Invalid scenario ID
- `404` - Scenario not found
- `500` - Error calculating standings

---

## Playoffs

### Get Playoff State
**GET** `/playoffs/scenarios/:scenario_id/state`

Returns playoff state for a scenario.

**Headers (Optional):**
```
Authorization: Bearer <token>
```

**Parameters:**
- `scenario_id` (path) - Scenario ID

**Response (200 OK):**
```json
{
  "playoff_state": {
    "id": 1,
    "scenario_id": 1,
    "current_round": 1,
    "is_enabled": true,
    "created_at": "2025-01-15T00:00:00Z",
    "updated_at": "2025-01-15T00:00:00Z"
  },
  "can_enable": true
}
```

**Notes:**
- `can_enable` is `true` when all regular season games are complete/picked
- `playoff_state` is `null` if playoffs not yet enabled

**Playoff Round Numbers:**
- **NFL**: 1=Wild Card, 2=Divisional, 3=Conference Championships, 4=Super Bowl
- **NBA**: 1=Play-In A, 2=Play-In B, 3=Conference Quarterfinals, 4=Conference Semifinals, 5=Conference Finals, 6=NBA Finals

**Errors:**
- `403` - Unauthorized (not owner)
- `404` - Scenario not found

---

### Enable Playoffs
**POST** `/playoffs/scenarios/:scenario_id/enable`

Enables playoffs for a scenario (generates first round).

**Headers (Optional):**
```
Authorization: Bearer <token>
```

**Parameters:**
- `scenario_id` (path) - Scenario ID

**Response (200 OK):**
```json
{
  "message": "NFL playoffs enabled successfully"
}
```

**Notes:**
- **NFL**: Generates Wild Card round (6 games per conference)
- **NBA**: Generates Play-In Round A (7v8, 9v10 per conference)
- Requires all regular season games complete/picked
- Seeds determined by standings

**Errors:**
- `400` - Not all regular season games complete
- `403` - Unauthorized (not owner)
- `404` - Scenario not found

---

### Get Playoff Matchups/Series
**GET** `/playoffs/scenarios/:scenario_id/rounds/:round`

Returns playoff matchups or series for a specific round.

**Headers (Optional):**
```
Authorization: Bearer <token>
```

**Parameters:**
- `scenario_id` (path) - Scenario ID
- `round` (path) - Round number

**NFL/NBA Play-In Response (200 OK):**
```json
[
  {
    "id": 1,
    "round": 1,
    "matchup_order": 1,
    "game_number": null,
    "conference": "AFC",
    "higher_seed": 2,
    "lower_seed": 7,
    "higher_seed_team_id": 5,
    "lower_seed_team_id": 18,
    "picked_team_id": 5,
    "predicted_higher_seed_score": 24,
    "predicted_lower_seed_score": 17,
    "status": "pending",
    "created_at": "2025-01-15T00:00:00Z",
    "updated_at": "2025-01-15T00:00:00Z",
    "higher_seed_team": {
      "id": 5,
      "abbreviation": "BUF",
      "city": "Buffalo",
      "name": "Bills",
      "logo_url": "https://...",
      "alternate_logo_url": "https://...",
      "primary_color": "00338d",
      "secondary_color": "c60c30"
    },
    "lower_seed_team": { ... }
  }
]
```

**NBA Series Response (200 OK):**
```json
[
  {
    "id": 1,
    "round": 3,
    "series_order": 1,
    "conference": "Eastern",
    "higher_seed": 1,
    "lower_seed": 8,
    "higher_seed_team_id": 3,
    "lower_seed_team_id": 15,
    "picked_team_id": 3,
    "predicted_higher_seed_wins": 4,
    "predicted_lower_seed_wins": 2,
    "best_of": 7,
    "status": "pending",
    "created_at": "2025-01-20T00:00:00Z",
    "updated_at": "2025-01-20T00:00:00Z",
    "higher_seed_team": { ... },
    "lower_seed_team": { ... }
  }
]
```

**Errors:**
- `403` - Unauthorized (not owner)
- `404` - Playoffs not enabled / Round not generated

---

### Update Playoff Pick
**PUT** `/playoffs/scenarios/:scenario_id/matchups/:matchup_id`

Updates a playoff matchup or series pick.

**Headers (Optional):**
```
Authorization: Bearer <token>
```

**Parameters:**
- `scenario_id` (path) - Scenario ID
- `matchup_id` (path) - Matchup or Series ID

**Request Body (Single Game):**
```json
{
  "picked_team_id": 5,
  "predicted_higher_seed_score": 24,
  "predicted_lower_seed_score": 17
}
```

**Request Body (Series):**
```json
{
  "picked_team_id": 3,
  "predicted_higher_seed_wins": 4,
  "predicted_lower_seed_wins": 2
}
```

**Notes:**
- If both scores/wins provided, `picked_team_id` auto-set to winner
- Deletes all subsequent playoff rounds
- Series must have one team reach 4 wins
- Updates scenario's `updated_at` timestamp

**Response (200 OK):**
```json
{
  "id": 1,
  "picked_team_id": 5,
  "predicted_higher_seed_score": 24,
  "predicted_lower_seed_score": 17
}
```

**Errors:**
- `400` - Invalid scores/wins
- `403` - Unauthorized (not owner)
- `404` - Matchup/series not found
- `500` - Database error

---

### Generate Next Playoff Round
**POST** `/playoffs/scenarios/:scenario_id/generate`

Generates the next playoff round based on current round picks.

**Headers (Optional):**
```
Authorization: Bearer <token>
```

**Parameters:**
- `scenario_id` (path) - Scenario ID

**Response (200 OK):**
```json
{
  "message": "Next round generated successfully"
}
```

**Notes:**
- Requires all picks in current round to be complete
- Automatically determines matchups based on winners
- **NFL**: Divisional round reseeds based on original seeds
- **NBA Play-In**: Round B matches winner 9v10 vs loser 7v8

**Errors:**
- `400` - Current round not complete
- `403` - Unauthorized (not owner)
- `404` - Playoffs not enabled
- `500` - Generation failed

---

### Delete Playoff Pick
**DELETE** `/playoffs/scenarios/:scenario_id/matchups/:matchup_id`

Deletes a playoff pick.

**Headers (Optional):**
```
Authorization: Bearer <token>
```

**Parameters:**
- `scenario_id` (path) - Scenario ID
- `matchup_id` (path) - Matchup or Series ID

**Response (200 OK):**
```json
{
  "message": "Playoff pick deleted successfully"
}
```

**Notes:**
- Sets pick fields to NULL rather than deleting the matchup/series
- Updates scenario's `updated_at` timestamp

**Errors:**
- `403` - Unauthorized (not owner)
- `404` - Matchup/series not found

---

## Admin

### Trigger NFL Schedule Update
**POST** `/admin/update-schedule/nfl`

Manually triggers an NFL schedule update from ESPN API.

**Response (200 OK):**
```json
{
  "status": "ok",
  "message": "NFL schedule update triggered"
}
```

**Notes:**
- Updates game scores, start times, and status
- Only updates final scores for completed games
- Runs automatically daily at midnight PST

---

### Trigger NBA Schedule Update
**POST** `/admin/update-schedule/nba`

Manually triggers an NBA schedule update from ESPN API.

**Response (200 OK):**
```json
{
  "status": "ok",
  "message": "NBA schedule update triggered"
}
```

**Notes:**
- Updates game scores, start times, and status
- Only updates final scores for completed games
- Runs automatically daily at midnight PST

---

## Error Handling

All error responses follow this format:

```json
{
  "error": "Error message description"
}
```

### Common HTTP Status Codes:
- `200 OK` - Request succeeded
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request (missing fields, validation errors)
- `401 Unauthorized` - Authentication required or invalid token
- `403 Forbidden` - Authenticated but not authorized (not owner)
- `404 Not Found` - Resource doesn't exist
- `423 Locked` - Account temporarily locked
- `500 Internal Server Error` - Server error (database, unexpected errors)

---

## Authentication

Most endpoints support **optional authentication**:
- **Guest users**: Get a session token automatically (stored in cookies)
- **Registered users**: Use JWT token in `Authorization` header

**JWT Token Format:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**Token Expiration:** 7 days

**Session Token:** 7 days (stored in HTTP-only cookie)

---

## Rate Limiting

**Auth endpoints** are rate limited:
- 5 attempts per 15 minutes for `/auth/register` and `/auth/login`
- Account locks for 15 minutes after 5 failed login attempts

---

## Data Sources

- **Game Data**: Fetched from ESPN API
- **Schedule Updates**: Automatic daily at midnight PST
- **Supported Sports**: NFL (fully supported), NBA (fully supported)
- **Time Zones**: All times in UTC, converted to PST for display

---

## Notes

1. **Guest Scenarios**: 
   - Created without authentication, tied to session token
   - Can be "claimed" after registering/logging in
   - Session expires after 7 days

2. **Standings Calculation**:
   - Uses actual game results + user picks
   - Games without picks use actual results if completed
   - Picks without scores use 1-0 dummy scores for W/L only
   - Points calculations exclude dummy scores

3. **Playoff Bracket**:
   - NFL: Single elimination, reseeds in Divisional round
   - NBA: Best-of-7 series, Play-In Tournament for seeds 7-10
   - Changing regular season picks deletes playoff brackets

4. **Password Security**:
   - Bcrypt hashing with cost factor 12
   - Strict validation (8+ chars, uppercase, lowercase, number, special char)
   - Account lockout after failed attempts

5. **Database Constraints**:
   - Unique constraint on (season_id, espn_id) for games
   - Unique constraint on (scenario_id, game_id) for picks
   - Cascading deletes for scenarios → picks → playoff data