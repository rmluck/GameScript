# GameScript

![Main Logo](/frontend/static/site-logo.png)

*By Rohan Mistry - Last updated January 29, 2026*

---

## 📖 Overview

GameScript is a comprehensive full-stack web application that enables users to create custom NFL and NBA playoff scenarios by simulating game outcomes and exploring their impact on standings, playoff seeding, and draft order. Built with a SvelteKit frontend and Go backend, the platform combines real-time sports data from ESPN's API with sophisticated standings calculation algorithms that account for complex tiebreaker rules specific to each league. Users can pick winners for individual games or entire playoff series, watch standings update in real-time, view team profiles, and save their scenarios to account profiles after registering with an email and password.

The platform demonstrates advanced full-stack development capabilities through its PostgreSQL database architecture, RESTful API design built with Fiber, and responsive user interface that handles complex state management across multiple sports leagues. Key technical achievements include implementing NFL and NBA playoff bracket generation systems, automated daily schedule updates via background schedulers, comprehensive tiebreaker logic following official league rules, and seamless deployment across cloud platforms (Vercel, Render, Supabase). The application processes 280+ NFL games and 1,230+ NBA games per season while managing intricate relational data structures and maintaining optimal performance through connection pooling and efficient query optimization.

**Target Users** are sports fans and analysts who want to explore hypothetical playoff scenarios and understand how game outcomes affect postseason positioning.

🔗 **Try it live**: [https://gamescript.live](https://gamescript.live)

<br>

**Home Page**

![Home Page](/static/img/home_page.png)

**NFL Scenario Page**

![NFL Scenario Page](/static/img/nfl_scenario_page.png)

**NBA Scenario Page**

![NBA Scenario Page](/static/img/nba_scenario_page.png)

**Profile Page**

![Profile Page](/static/img/profile_page.png)

---

## 📁 Contents

```bash
|-- .github/
|   └── workflows/
|       └── keepalive.yml               # GitHub Actions workflow for backend health checks
|-- backend
|   |-- cmd/
|   |   └── server/
|   |       └── main.go                 # Application entry point
|   |-- database/
|   |   |-- nba/
|   |   |   |-- schedules/              # NBA schedule JSON data
|   |   |   └── teams/                  # NBA teams JSON data
|   |   |-- nfl/
|   |   |   |-- schedules/              # NFL schedule JSON data
|   |   |   └── teams/                  # NFL teams JSON data
|   |   └── schema.sql                  # PostgreSQL database schema
|   |-- internal/
|   |   |-- database/
|   |   |   └── db.go                   # Database connection management
|   |   |-- handlers/
|   |   |   |-- auth.go                 # Authentication endpoints
|   |   |   |-- games.go                # Games API handlers
|   |   |   |-- handlers.go             # Route setup
|   |   |   |-- picks.go                # User picks handlers
|   |   |   |-- playoffs.go             # Playoff bracket handlers
|   |   |   |-- scenarios.go            # Scenario CRUD handlers
|   |   |   |-- standings.go            # Standings calculation handlers
|   |   |   └── teams.go                # Teams API handlers
|   |   |-- middleware/
|   |   |   |-- auth.go                 # JWT authentication middleware
|   |   |   └── rate_limit.go           # Rate limiting middleware
|   |   |-- models/
|   |   |   |-- espn.go                 # ESPN API response models
|   |   |   └── models.go               # Core data models
|   |   |-- playoffs/
|   |   |   |-- nba_playoffs.go         # NBA playoff bracket generation
|   |   |   └── nfl_playoffs.go         # NFL playoff bracket generation
|   |   |-- scheduler/
|   |   |   |-- scheduler.go            # Background job scheduler
|   |   |   |-- nba_scheduler.go        # NBA daily updates
|   |   |   └── nfl_scheduler.go        # NFL daily updates
|   |   |-- services/
|   |   |   └── espn/
|   |   |       |-- client.go           # ESPN API client
|   |   |       |-- nba_schedule.go     # NBA schedule fetcher
|   |   |       |-- nba_teams.go        # NBA teams fetcher
|   |   |       |-- nfl_schedule.go     # NFL schedule fetcher
|   |   |       └── nfl_teams.go        # NFL teams fetcher
|   |   └── standings/
|   |       |-- nba_standings.go        # NBA standings & tiebreaker logic
|   |       └── nfl_standings.go        # NFL standings & tiebreaker logic
|   └── scripts/
|       |-- fetch_data/
|       |   |-- fetch_nba_schedule.go   # Fetch NBA data from ESPN
|       |   |-- fetch_nba_teams.go
|       |   |-- fetch_nfl_schedule.go   # Fetch NFL data from ESPN
|       |   └── fetch_nfl_teams.go
|       └── import_data/
|           |-- import_nba_schedule.go  # Import NBA data to database
|           |-- import_nba_teams.go
|           |-- import_nfl_schedule.go  # Import NFL data to database
|           └── import_nfl_teams.go
|-- docs/
|   |-- API.md                          # API documentation
|   └── Standings Rules.md              # Sport-specific tiebreaker rules
|-- frontend/
|   |-- src/
|   |   |-- lib/
|   |   |   |-- api/
|   |   |   |   |-- auth.ts             # Authentication API calls
|   |   |   |   |-- client.ts           # Axios client configuration
|   |   |   |   |-- games.ts            # Games API calls
|   |   |   |   |-- picks.ts            # Picks API calls
|   |   |   |   |-- playoffs.ts         # Playoffs API calls
|   |   |   |   |-- scenarios.ts        # Scenarios API calls
|   |   |   |   |-- standings.ts        # Standings API calls
|   |   |   |   └── teams.ts            # Teams API calls
|   |   |   |-- components/
|   |   |   |   |-- nba/                # NBA-specific components
|   |   |   |   |   |-- PlayoffGameCard.svelte
|   |   |   |   |   |-- PlayoffPicksBox.svelte
|   |   |   |   |   |-- StandingsBox.svelte
|   |   |   |   |   |-- StandingsBoxExpanded.svelte
|   |   |   |   |   └── TeamModal.svelte
|   |   |   |   |-- nfl/                # NFL-specific components
|   |   |   |   |   |-- PlayoffGameCard.svelte
|   |   |   |   |   |-- PlayoffPicksBox.svelte
|   |   |   |   |   |-- StandingsBox.svelte
|   |   |   |   |   |-- StandingsBoxExpanded.svelte
|   |   |   |   |   └── TeamModal.svelte
|   |   |   |   └── scenarios/          # Shared scenario components
|   |   |   |       |-- ByeTeams.svelte
|   |   |   |       |-- ComingSoonModal.svelte
|   |   |   |       |-- ConfirmationModal.svelte
|   |   |   |       |-- CreateScenarioModal.svelte
|   |   |   |       |-- DraftOrderBox.svelte
|   |   |   |       |-- GameCard.svelte
|   |   |   |       |-- GamePickerRow.svelte
|   |   |   |       |-- PicksBox.svelte
|   |   |   |       |-- ScenarioHeader.svelte
|   |   |   |       |-- ScenarioInfo.svelte
|   |   |   |       |-- ScenarioSettings.svelte
|   |   |   |       |-- WeekNavigator.svelte
|   |   |   |-- stores/
|   |   |   |   └── auth.ts             # Svelte auth store
|   |   |   |-- types/
|   |   |   |   └── index.ts            # TypeScript type definitions
|   |   |   └── utils/
|   |   |       |-- nba/
|   |   |       |   └── dates.ts        # NBA week calculations
|   |   |       |-- nfl/
|   |   |       |   └── dates.ts        # NFL week calculations
|   |   |       └── validation.ts       # Form validation utilities
|   |   └── routes/
|   |       |-- auth/
|   |       |   |-- login/              # Login page
|   |       |   └── register/           # Registration page
|   |       |-- nfl/                    # NFL quick-create page
|   |       |-- nba/                    # NBA quick-create page
|   |       |-- profile/                # User profile page
|   |       |-- scenarios/
|   |       |   |-- nfl/[id]/           # NFL scenario page
|   |       |   └── nba/[id]/           # NBA scenario page
|   |       |-- +layout.svelte          # Root layout
|   |       └── +page.svelte            # Home page
|   |-- static/                         # Static assets (favicon, images)
|   └── .env.production                 # Production environment variables
|-- LICENSE                             # MIT License
└── README.md                           # Project documentation
```

---

## 🌟 Features

### Core Functionality

* **Multi-Sport Support**: Create scenarios for NFL and NBA with sport-specific rules.
* **Real-Time Standings**: Instantly see how each pick affects conference/division standings.
* **Playoff Seeding**: Automatic playoff bracket generation following official league rules.
* **Game Picking**: Select winners for regular season games with optional score predictions.
* **Playoff Brackets**: Simulate entire playoff tournaments with series and/or single-game predictions.
* **Tiebreakers**: Comprehensive implementation of sport-specific tiebreaker procedures.

### User Experience

* **Week Navigation**: Jump to any week of the season with keyboard shortcuts.
* **Team Profiles**: Click any team to view their complete schedule and stats.
* **Scenario Management**: Name, edit, and delete your scenarios.
* **Share Links**: Copy unique URLs to share scenarios with friends.
* **Responsive Design**: Optimized for desktop and mobile viewing.
* **Dark Theme**: Modern, eye-friendly dark color scheme.

### Advanced Features

* **User Accounts**: Save unlimited scenarios with authentication.
* **Guest Mode**: Create scenarios without signing up (session-based).
* **Undo/Reset**: Modify picks and reset playoff rounds.
* **Game Info**: Hover infobox for each game with its date, time, location, and network.
* **Bye Weeks**: NFL bye team tracking.
* **Play-In Tournament**: NBA play-in bracket generation.
* **Standings Views**: Toggle between conference and division view.
* **Team Colors**: Dynamic theming based on team primary colors.
* **Custom Logos**: High-quality team logos with alternate versions.

### Technical Features

* **Auto-Updates**: Daily schedule and score updates at midnight PST.
* **API Access**: RESTful API with comprehensive documentation.
* **Session Management**: JWT authentication with 7-day expiration.
* **Rate Limiting**: Protection against brute force attacks.
* **CORS Security**: Secure cross-origin resource sharing.
* **Database Migrations**: Version-controlled schema changes.
* **Background Jobs**: Automated data synchronization.
* **Keepalive System**: GitHub Actions workflow to prevent cold starts.

---

## 🛠️ Installation Instructions

### Prerequisites
* Node.js 18+ and npm
* Go 1.21+
* PostgreSQL 14+

### Backend Setup

1. Clone the repository:
```bash
git clone https://github.com/rmluck/GameScript.git
cd GameScript/backend
```

2. Create __.env__ file for backend:
```bash
# Database (use Supabase connection pooling URL)
DATABASE_URL=postgresql://user:password@host:6543/database?pgbouncer=true

# JWT Secret
JWT_SECRET=your_jwt_secret_key

# Server
PORT=8080

# CORS
ALLOWED_ORIGINS=http://localhost:5173,https://gamescript.live,https://www.gamescript.live

# Security
MAX_LOGIN_ATTEMPTS=5
LOCKOUT_DURATION_MINUTES=15
```

3. Set up database:
```bash
# Create database and run schema
psql -d your_database < database/schema.sql

# Import sports and seasons data
psql -d your_database < database/sports.sql
psql -d your_database < database/seasons.sql
```

4. Import team and schedule data:
```bash
# Import NFL teams
go run scripts/import_data/import_nfl_teams.go

# Import NBA teams
go run scripts/import_data/import_nba_teams.go

# Import NFL schedule
go run scripts/import_data/import_nfl_schedule.go

# Import NBA schedule
go run scripts/import_data/import_nba_schedule.go`
```

5. Run the server
```bash
go run cmd/server/main.go
```

Backend will be available at __http://localhost:8080

### Frontend Setup

1. Naviage to frontend
```bash
cd ../frontend
```

2. Install dependencies
```bash
npm install
```

3. Create __.env__ file for frontend:
```bash
PUBLIC_API_URL=http://localhost:8080/api
```

4. Run development server:
```bash
npm run dev
```

Frontend will be available at __http://localhost:5173__

---

## 💡 Usage

**Step 1: Create a Scenario**

Choose your sport and season. You can either:
* Click "Create Scenario" to manually configure settings
* Use quick-create buttons in the navigation bar that automatically generate scenarios

Give your scenario a custom name and choose whether to make it public (shareable) or private.

![Create Scenario Modal](/static/img/create_scenario_modal.png)

**Step 2: Make Picks**

Regular Season:
* Navigate between weeks using arrow keys or the week selector
* Click team logos for each game to select winners
* Optionally enter score predictions for tiebreaker purposes
* See game details including date, time, location, and network

Standings:
* Watch standings update in real-time as you make picks
* Can toggle between conference and division view
* In expanded standings view, scroll through stats for each team (overall/home/away records, division/conference records, point differential, games back calculations, etc.)
* Click any team in the standings to view their full schedule

Team Modal:
* View complete win-loss record and statistics breakdown
* See all games for the team's season with results
* Make picks directly from the team schedule view

![Team Modal](/static/img/team_modal.png)

**Step 3: Enable Playoffs**

Once all regular season games are complete (picked or final), you can enable playoffs:
* Banner appears prompting you to enable playoffs
Playoff bracket auto-generates based on final seeding
* Navigate through playoff rounds sequentially

NFL Playoffs:
* Wild Card Round (6 games)
* Divisional Round (4 games)
* Conference Championships (2 games)
* Super Bowl (1 game)

NBA Playoffs:
* Play-In Tournament Round A (4 games: 7v8, 9v10)
* Play-In Tournament Round B (2 games: Winner 9v10 vs Loser 7v8)
* Conference Quarterfinals (8 series, best-of-7)
* Conference Semifinals (4 series, best-of-7)
* Conference Finals (2 series, best-of-7)
* NBA Finals (1 series, best-of-7)

Playoff Predictions:
* For single-elimination games, pick winner and optionally input scores
* For series, pick series winner and optionally input games won
* Warning when changing earlier round picks (resets later rounds)

---

## 🚧 Future Improvements

### Near-Term Enhancements
* College Football support
* View public scenarios (cannot edit)
* Draft order view
* Draft order lottery simulator for NBA

### User Experience
* Reset all picks
* Undo individual picks (not just last pick)
* Bulk pick selection (sim-to-end of season)
* Dark/light mode toggle
* Mobile app (iOS/Android)

### Social Features
* Public scenarios list
* Social media integration
* Comments and discussions

### Advanced Analytics
* Win probability models
* Team-by-team playoff percentage calculators

### Technical Improvements
* GraphQL API option
* WebSocket real-time updats
* Redis caching layer
* Database query optimization
* Advanced error handling
* Comprehensive test coverage
* API rate limiting per user
* Backup and restore functionality

---

## 🧰 Tech Stack

### Frontend Technologies
* **TypeScript**: Type-safe JavaScript for robust development
* **HTML 5**: Semantic markup structure
* **CSS3**: Modern styling with custom properties
* **Svelte 4**: Reactive UI framework with minimal runtime
* **SvelteKit**: Full-stack framework for Svelte apps
* **Vite**: Lightning-fast build tool and dev server
* **Axios**: Promise-based HTTP client for API calls
* **Tailwind CSS**: Utility-first CSS framework

### Backend & API
* **Go 1.21**: High-performance backend language
* **Fiber v2**: Express-inspired web framework for Go
* **SQL**: Database queries and migrations
* **ESPN API**: Real-time sports data (schedules, scores, teams)
* **JWT**: JSON Web Tokens for authentication
* **bcrypt**: Secure password hashing

### Database & Data Management
* **PostgreSQL**: Robust relational database
* **Supabase**: PostgreSQL hosting and management
* **JSON**: Data interchange format

### Real-Time & State Management
* **Svelte Stores**: Reactive state management
* **Custom Hooks**: Encapsulated business logic
* **WebSockets**: (Planned) Real-time updates

### Development Tools
* **ESLint**: Code linting
* **Prettier**: Code formatting
* **Git**: Version control
* **GitHub Actions**: CI/CD pipelines

### Build & Development
* **Node.js**: JavaScript runtime
* **npm**: Package management

### Security & Performance
* **CORS**: Cross-origin security
* **Rate Limiting**: API protection
* **Connection Pooling**: Database optimization
* **Gzip Compression**: Response optimization
* **SSL/TLS**: Encrypted connections

### Deployment
* **Vercel**: Frontend hosting with CDN
* **Render**: Backend hosting with auto-scaling
* **Squarespace**: Custom domain name registration

## 📊 Project Statistics
* **Total Lines of Code**: ~15,000+
* **Backend**: ~8,000 lines (Go)
* **Frontend**: ~7,000 lines (TypeScript/Svelte)
* **API Endpoints**: 30+
* **Database Tables**: 10
* **NFL Games Per Season**: 285
* **NBA Games Per Season**: 1,230
* **Supported Leagues**: 2 (NFL, NBA)

## 🙏 Contributions/Acknowledgements

This project was built independently as a portfolio project. Inspired by playoff scenario calculators such as [Playoff Predictors](https://playoffpredictors.com) and [Pro Football Network's playoff predictors](https://www.profootballnetwork.com/nfl-playoff-predictor).

## 🪪 License

This project is licensed under the [MIT License](/LICENSE).