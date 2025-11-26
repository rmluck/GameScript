# Error Handling Checklist

## Authentication Errors
- [ ] Missing email/username/password returns `400`
- [ ] Invalid email format returns `400`
- [ ] Password < 8 characters returns `400`
- [ ] Duplicate email returns `400` with clear message
- [ ] Duplicate username returns `400` with clear message
- [ ] Invalid login credentials return `401`
- [ ] Missing Authorization header returns `401` (for protected routes)
- [ ] Invalid JWT token returns `401`
- [ ] Expired JWT token returns `401`

## Scenarios Errors
- [ ] Missing required fields (name, sport_id, season_id) return `400`
- [ ] Invalid scenario ID returns `404`
- [ ] Accessing another user's scenario returns `403`
- [ ] Deleting non-existent scenario returns `404`
- [ ] Updating with no fields returns `400`

## Picks Errors
- [ ] Invalid picked_team_id returns `400` or `500`
- [ ] Duplicate pick (same scenario + game) returns `500`
- [ ] Invalid game ID returns `404`
- [ ] Accessing another user's pick returns `403`
- [ ] Missing predicted scores returns error

## Games Errors
- [ ] Invalid season ID returns empty array or `404`
- [ ] Invalid week number returns empty array
- [ ] Invalid team ID returns `404`
- [ ] Invalid game ID returns `404`

## Standings Errors
- [ ] Invalid scenario ID returns `404`
- [ ] Non-NFL sport returns `400` with "Only NFL supported" message
- [ ] Database error returns `500`

## General Errors
- [ ] Malformed JSON returns `400`
- [ ] Database connection failure returns `500`
- [ ] Panic recovery middleware catches crashes
- [ ] All errors return JSON (not HTML)
- [ ] Error messages are user-friendly (no stack traces)