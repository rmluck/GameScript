#!/bin/bash

# GameScript API Testing Script
# Run this after starting server with: go run cmd/server/main.go

BASE_URL="http://localhost:8080/api"
TOKEN=""
SCENARIO_ID=""
GAME_ID=""

echo "=========================================="
echo "GameScript API Manual Testing"
echo "=========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Test function
test_endpoint() {
    local name=$1
    local method=$2
    local endpoint=$3
    local data=$4
    local auth=$5
    
    echo "Testing: $name"
    echo "Method: $method $endpoint"
    
    if [ -n "$auth" ]; then
        response=$(curl -s -X $method "$BASE_URL$endpoint" \
            -H "Content-Type: application/json" \
            -H "Authorization: Bearer $TOKEN" \
            -d "$data")
    else
        response=$(curl -s -X $method "$BASE_URL$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data")
    fi
    
    echo "Response: $response"
    echo ""
    echo "$response"
}

# 1. Health Check
echo "=========================================="
echo "1. HEALTH CHECK"
echo "=========================================="
test_endpoint "Health Check" "GET" "/health"

# 2. Public Endpoints (No Auth)
echo "=========================================="
echo "2. PUBLIC ENDPOINTS"
echo "=========================================="

test_endpoint "Get Sports" "GET" "/sports"
test_endpoint "Get Seasons for NFL" "GET" "/sports/1/seasons"
test_endpoint "Get Season Details" "GET" "/seasons/1"
test_endpoint "Get Teams for Season" "GET" "/seasons/1/teams"
test_endpoint "Get Single Team" "GET" "/teams/1"
test_endpoint "Get Games for Season" "GET" "/seasons/1/games"
test_endpoint "Get Games for Week 1" "GET" "/seasons/1/weeks/1/games"
test_endpoint "Get Games for Team" "GET" "/teams/1/games"

# 3. Authentication
echo "=========================================="
echo "3. AUTHENTICATION"
echo "=========================================="

# Register new user
echo "Registering new user..."
register_response=$(curl -s -X POST "$BASE_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "test@example.com",
        "username": "testuser",
        "password": "password123"
    }')
echo "Register Response: $register_response"
TOKEN=$(echo $register_response | jq -r '.token')
echo "Token: $TOKEN"
echo ""

# Login
echo "Logging in..."
login_response=$(curl -s -X POST "$BASE_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "test@example.com",
        "password": "password123"
    }')
echo "Login Response: $login_response"
TOKEN=$(echo $login_response | jq -r '.token')
echo "Token: $TOKEN"
echo ""

# Get Current User
test_endpoint "Get Current User" "GET" "/auth/me" "" "auth"

# 4. Scenarios (Protected)
echo "=========================================="
echo "4. SCENARIOS (PROTECTED)"
echo "=========================================="

# Create Scenario
echo "Creating scenario..."
create_scenario_response=$(curl -s -X POST "$BASE_URL/scenarios" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{
        "name": "My Test Scenario",
        "sport_id": 1,
        "season_id": 1,
        "is_public": true
    }')
echo "Create Scenario Response: $create_scenario_response"
SCENARIO_ID=$(echo $create_scenario_response | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
echo "Scenario ID: $SCENARIO_ID"
echo ""

test_endpoint "Get All Scenarios" "GET" "/scenarios" "" "auth"
test_endpoint "Get Single Scenario" "GET" "/scenarios/$SCENARIO_ID" "" "auth"

# Update Scenario
test_endpoint "Update Scenario" "PUT" "/scenarios/$SCENARIO_ID" '{
    "name": "Updated Scenario Name",
    "is_public": false
}' "auth"

# 5. Picks (Protected)
echo "=========================================="
echo "5. PICKS (PROTECTED)"
echo "=========================================="

# Get a game ID first
echo "Getting first game..."
games_response=$(curl -s "$BASE_URL/seasons/1/weeks/1/games")
GAME_ID=$(echo $games_response | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
echo "Game ID: $GAME_ID"
echo ""

# Create Pick
test_endpoint "Create Pick" "POST" "/picks/scenarios/$SCENARIO_ID/games/$GAME_ID" '{
    "picked_team_id": 1,
    "predicted_home_score": 24,
    "predicted_away_score": 17
}' "auth"

test_endpoint "Get All Picks for Scenario" "GET" "/picks/scenarios/$SCENARIO_ID" "" "auth"
test_endpoint "Get Single Pick" "GET" "/picks/scenarios/$SCENARIO_ID/games/$GAME_ID" "" "auth"

# Update Pick
test_endpoint "Update Pick" "PUT" "/picks/scenarios/$SCENARIO_ID/games/$GAME_ID" '{
    "picked_team_id": 2,
    "predicted_home_score": 27,
    "predicted_away_score": 20
}' "auth"

# 6. Standings
echo "=========================================="
echo "6. STANDINGS"
echo "=========================================="

test_endpoint "Get Standings" "GET" "/scenarios/$SCENARIO_ID/standings" "" "auth"

# 7. Edge Cases
echo "=========================================="
echo "7. EDGE CASES"
echo "=========================================="

test_endpoint "Invalid Scenario ID" "GET" "/scenarios/99999" "" "auth"
test_endpoint "Unauthorized Access (No Token)" "GET" "/scenarios"
test_endpoint "Invalid Login" "POST" "/auth/login" '{
    "email": "wrong@example.com",
    "password": "wrongpassword"
}'
test_endpoint "Missing Required Fields" "POST" "/scenarios" '{
    "name": "Incomplete Scenario"
}' "auth"
test_endpoint "Invalid Pick Data" "POST" "/picks/scenarios/$SCENARIO_ID/games/$GAME_ID" '{
    "picked_team_id": "not_a_number"
}' "auth"

# 8. Delete Operations
echo "=========================================="
echo "8. DELETE OPERATIONS"
echo "=========================================="

test_endpoint "Delete Pick" "DELETE" "/picks/scenarios/$SCENARIO_ID/games/$GAME_ID" "" "auth"
test_endpoint "Delete Scenario" "DELETE" "/scenarios/$SCENARIO_ID" "" "auth"

echo "=========================================="
echo "Testing Complete!"
echo "=========================================="