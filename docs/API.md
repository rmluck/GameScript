# GameScript API Documentation

**Base URL:** `http://localhost:8080/api`

**Version:** 1.0.0

---

## Table of Contents

1. [Authentication](#authentication)
2. [Sports](#sports)
3. [Seasons](#seasons)
4. [Teams](#teams)
5. [Games](#games)
6. [Scenarios](#scenarios)
7. [Picks](#picks)
8. [Standings](#standings)
9. [Admin](#admin)
10. [Error Handling](#error-handling)

---

## Authentication

### Register User
**POST** `/auth/register`

Creates a new user account.

**Request Body:**
```json
{
  "email": "user@example.com",
  "username": "username",
  "password": "password123" // min 8 characters
}
```

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
- `400` - Missing required fields / Password too short / Email/username already exists

---

### Login User
**POST** `/auth/login`

Authenticates a user and returns a JWT token.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
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
- `401` - Invalid email or password

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

## Sports

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
  }
]
```

---

## Seasons

### Get Seasons for a Sport
**GET** `/sports/:sport_id/seasons`

Returns all seasons for a specific sport.

**Parameters:**
- `sport_id` (path) - Sport ID

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "sport_id": 1,
    "start_year": 2025,
    "end_year": null,
    "is_active": true,
    "created_at": "2025-01-01T00:00:00Z"
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
  "start_year": 2025,
  "end_year": null,
  "is_active": true,
  "created_at": "2025-01-01T00:00:00Z"
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
    "logo_url": "https://..."
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
  "abbreviation": "ARI",
  "city": "Arizona",
  "name": "Cardinals",
  "conference": "NFC",
  "division": "NFC West",
  "primary_color": "a40227",
  "logo_url": "https://..."
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
    "start_time": "2025-09-05T20:20:00-07:00",
    "day_of_week": "Thursday",
    "week": 1,
    "location": "Arrowhead Stadium, Kansas City, MO, USA",
    "primetime": "TNF",
    "network": "NBC",
    "home_score": null,
    "away_score": null,
    "status": "upcoming",
    "is_postseason": false,
    "home_team": {
      "id": 12,
      "abbreviation": "KC",
      "city": "Kansas City",
      "name": "Chiefs"
    },
    "away_team": {
      "id": 2,
      "abbreviation": "BAL",
      "city": "Baltimore",
      "name": "Ravens"
    }
  }
]
```

---

### Get Games by Week
**GET** `/seasons/:season_id/weeks/:week/games`

Returns all games for a specific week in a season.

**Parameters:**
- `season_id` (path) - Season ID
- `week` (path) - Week number

**Response:** Same as Get Games for a Season

---

### Get Games for a Team
**GET** `/teams/:team_id/games`

Returns all games for a specific team.

**Parameters:**
- `team_id` (path) - Team ID

**Response:** Same as Get Games for a Season

---

### Get Single Game
**GET** `/games/:game_id`

Returns details for a specific game.

**Parameters:**
- `game_id` (path) - Game ID

**Response (200 OK):**
```json
{
  "id": 1,
  "espn_id": "401671745",
  "start_time": "2025-09-05T20:20:00-07:00",
  "week": 1,
  "home_score": 27,
  "away_score": 24,
  "status": "final",
  "home_team": { ... },
  "away_team": { ... }
}
```

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
    "name": "My Scenario",
    "sport_id": 1,
    "season_id": 1,
    "season_start_year": 2025,
    "season_end_year": null,
    "is_public": true,
    "sport_short_name": "NFL",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
]
```

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
  "name": "My Scenario",
  "sport_id": 1,
  "season_id": 1,
  "is_public": true,
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
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
  "is_public": true,
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}
```

**Errors:**
- `400` - Missing required fields

---

### Update Scenario
**PUT** `/scenarios/:scenario_id`

Updates an existing scenario.

**Headers (Optional):**
```
Authorization: Bearer <token>
```

**Parameters:**
- `scenario_id` (path) - Scenario ID

**Request Body:**
```json
{
  "name": "Updated Name",
  "is_public": false
}
```

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "Updated Name",
  "is_public": false,
  "updated_at": "2025-01-02T00:00:00Z"
}
```

**Errors:**
- `400` - No fields to update
- `403` - Unauthorized (not owner)
- `404` - Scenario not found

---

### Delete Scenario
**DELETE** `/scenarios/:scenario_id`

Deletes a scenario.

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
    "predicted_away_score": 24,
    "status": "pending",
    "game": {
      "espn_id": "401671745",
      "start_time": "2025-09-05T20:20:00-07:00",
      "week": 1,
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
  "predicted_away_score": 24,
  "status": "pending"
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
  "predicted_away_score": 24
}
```

**Response (201 Created):**
```json
{
  "id": 1,
  "scenario_id": 1,
  "game_id": 1,
  "picked_team_id": 12,
  "predicted_home_score": 27,
  "predicted_away_score": 24,
  "status": "pending",
  "created_at": "2025-01-01T00:00:00Z"
}
```

**Errors:**
- `400` - Invalid request body
- `403` - Unauthorized (not owner)
- `500` - Database error (duplicate pick, etc.)

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

**Response (200 OK):**
```json
{
  "id": 1,
  "picked_team_id": 2,
  "predicted_home_score": 24,
  "predicted_away_score": 27,
  "updated_at": "2025-01-02T00:00:00Z"
}
```

**Errors:**
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

**Response (200 OK):**
```json
{
  "afc": {
    "divisions": {
      "AFC East": [
        {
          "rank": 1,
          "team_id": 2,
          "team_name": "Bills",
          "team_abbr": "BUF",
          "wins": 12,
          "losses": 5,
          "ties": 0,
          "win_pct": 0.706,
          "division_record": "5-1-0",
          "conference_record": "9-3-0",
          "points_for": 425,
          "points_against": 350,
          "point_diff": 75,
          "division_games_back": 0.0,
          "conference_games_back": 0.0
        }
      ]
    },
    "playoff_seeds": [
      {
        "seed": 1,
        "team_id": 2,
        "team_name": "Bills",
        "team_abbr": "BUF",
        "wins": 12,
        "losses": 5,
        "ties": 0,
        "is_division_winner": true
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
      "record": "2-15-0",
      "reason": "Non-playoff"
    }
  ]
}
```

**Errors:**
- `400` - Invalid scenario ID / Only NFL supported currently
- `404` - Scenario not found
- `500` - Error calculating standings

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

**Session Token:** 30 days (stored in HTTP-only cookie)

---

## Rate Limiting

Currently **no rate limiting** is implemented. This will be added in a future version.

---

## Notes

1. **Guest Scenarios**: Scenarios created without authentication are tied to a session token. They can be "claimed" after registering/logging in.

2. **Standings Tiebreakers**: Currently implements:
   - Win percentage
   - Head-to-head record
   - Division record (for division tiebreakers)
   - Conference record (for conference tiebreakers)
   - Point differential
   - Points scored/allowed
   
   *Not yet implemented:* Strength of victory/schedule, common games record, multi-team tiebreakers.

3. **Supported Sports**: Currently only NFL is fully supported. NBA and CFB coming soon.

4. **Time Zones**: All times are returned in Pacific Time (America/Los_Angeles).