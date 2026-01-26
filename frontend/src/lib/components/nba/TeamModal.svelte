<script lang="ts">
    import { createEventDispatcher, onMount } from 'svelte';
    import { gamesAPI } from '$lib/api/games';
    import { picksAPI } from '$lib/api/picks';
    import { standingsAPI } from '$lib/api/standings';
    import { playoffsAPI } from '$lib/api/playoffs';
    import GamePickerRow from '../scenarios/GamePickerRow.svelte';
    import type { NBAPlayoffSeed, Game, Pick, PlayoffState } from '$types';

    // Props
    export let team: NBAPlayoffSeed;
    export let scenarioId: number;
    export let seasonId: number;

    // Event dispatcher
    const dispatch = createEventDispatcher();

    // State variables for games and picks
    let teamGames: Game[] = [];
    let picks: Map<number, Pick> = new Map();

    // State variables for team info
    let currentTeam: NBAPlayoffSeed = team;
    let teamDivision: string = '';

    let playoffState: PlayoffState | null = null;

    // Loading and error states
    let loading = true;
    let error = '';

    // Load data on mount
    onMount(async () => {
        await loadTeamGames();
        await loadPicks();
        await loadPlayoffState();
    });

    async function loadTeamGames() {
        try {
            loading = true;
            const allGames = await gamesAPI.getBySeason(seasonId);
            teamGames = allGames.filter(
                game => game.home_team_id === team.team_id || game.away_team_id === team.team_id
            );

            await loadTeamDivision();
        } catch (err: any) {
            error = 'Failed to load team games';
            console.error(err);
        } finally {
            loading = false;
        }
    }

    async function loadPicks() {
        try {
            const allPicks = await picksAPI.getByScenario(scenarioId);
            picks = new Map(allPicks.map(pick => [pick.game_id, pick]));
        } catch (err: any) {
            console.error('Error loading picks:', err);
        }
    }

    async function loadTeamDivision() {
        try {
            // Get team's division from the first game
            if (teamGames.length > 0) {
                const firstGame = teamGames[0];
                const teamData = firstGame.home_team_id === team.team_id 
                    ? firstGame.home_team 
                    : firstGame.away_team;
                teamDivision = teamData.division;
            }
        } catch (err: any) {
            console.error('Error loading team division:', err);
        }
    }

    async function loadPlayoffState() {
        try {
            const response = await playoffsAPI.getState(scenarioId);
            playoffState = response.playoff_state;
        } catch (err: any) {
            console.error('Failed to load playoff state:', err);
        }
    }

    // Reload team data to get updated stats after pick changes
    async function reloadTeamData() {
        try {
            const standings = await standingsAPI.getByNBAScenario(scenarioId);
            
            // Find updated team data in standings
            const allSeeds = [
                ...standings.eastern.playoff_seeds,
                ...standings.western.playoff_seeds
            ];
            
            const updatedTeam = allSeeds.find(seed => seed.team_id === team.team_id);
            if (updatedTeam) {
                currentTeam = updatedTeam;
            }
        } catch (err: any) {
            console.error('Error reloading team data:', err);
        }
    }

    function closeModal() {
        dispatch('close');
    }

    // Close modal when clicking outside content
    function handleClickOutside(event: MouseEvent) {
        if (event.target === event.currentTarget) {
            closeModal();
        }
    }

    async function handlePickChange(event: CustomEvent) {
        try {
            const { gameId, pickedTeamId, predictedHomeScore, predictedAwayScore, deletePick } = event.detail;
            const existingPick = picks.get(gameId);

            if (deletePick && existingPick) {
                // Delete the pick
                await picksAPI.delete(scenarioId, gameId);
                picks.delete(gameId);
            } else if (existingPick) {
                // Update existing pick
                const updated = await picksAPI.update(scenarioId, gameId, {
                    picked_team_id: pickedTeamId === undefined ? null : pickedTeamId,
                    predicted_home_score: predictedHomeScore,
                    predicted_away_score: predictedAwayScore
                });
                picks.set(gameId, updated);
            } else {
                // Create new pick
                const updated = await picksAPI.create(scenarioId, gameId, {
                    picked_team_id: pickedTeamId === undefined ? null : pickedTeamId,
                    predicted_home_score: predictedHomeScore,
                    predicted_away_score: predictedAwayScore
                });
                picks.set(gameId, updated);
            }
            
            picks = new Map(picks);
            
            // Reload team data to get updated stats
            await reloadTeamData();
            
            dispatch('pickUpdated');
        } catch (err: any) {
            console.error('Error saving pick:', err);
            alert('Failed to save pick. Please try again.');
        }
    }

    function formatRecord(wins: number, losses: number): string {
        return `${wins}-${losses}`;
    }

    function formatWinPct(winPct: number): string {
        if (winPct === -1.0) return '.000'; // Handle 0-0 case
        return winPct.toFixed(3);
    }

    function formatGamesBack(gb: number): string {
        if (gb === 0) return '-';
        return gb.toFixed(1);
    }

    function formatPoints(points: number, games_with_scores: number): string {
        const points_average = points / (games_with_scores);
        return points_average.toFixed(1);
    }

    function formatPointDiff(pointsFor: number, pointsAgainst: number, gamesWithScores: number): string {
        const diff_average = (pointsFor - pointsAgainst) / gamesWithScores;
        if (diff_average > 0) return `+${diff_average.toFixed(1)}`;
        return diff_average.toFixed(1);
    }

    function formatDate(dateString: string): string {
        const date = new Date(dateString);
        return date.toLocaleDateString('en-US', { 
            weekday: 'short', 
            month: 'short', 
            day: 'numeric' 
        });
    }

    function formatTime(dateString: string): string {
        const date = new Date(dateString);
        return date.toLocaleTimeString('en-US', { 
            hour: 'numeric', 
            minute: '2-digit',
        });
    }

    function isHomeGame(game: Game): boolean {
        return game.home_team_id === team.team_id;
    }

    function getGameResult(game: Game): 'win' | 'loss' | 'upcoming' {
        const pick = picks.get(game.id);
        
        // Check if there's a pick for this game
        if (pick?.picked_team_id !== undefined && pick.picked_team_id !== null) {
            return pick.picked_team_id === team.team_id ? 'win' : 'loss';
        }
        
        // If no pick, check actual game result if final
        if (game.status === 'final' && game.home_score !== null && game.away_score !== null) {
            const teamScore = isHomeGame(game) ? game.home_score : game.away_score;
            const oppScore = isHomeGame(game) ? game.away_score : game.home_score;
            return (teamScore ?? 0) > (oppScore ?? 0) ? 'win' : 'loss';
        }
        
        return 'upcoming';
    }
</script>

<!-- Modal Overlay -->
<div 
    class="fixed inset-0 bg-black/80 flex items-center justify-center z-50 p-4"
    on:click={handleClickOutside}
    on:keydown={(e) => { if (e.key === 'Escape') closeModal(); }}
    role="dialog"
    aria-modal="true"
    tabindex="-1"
>
    <section
        class="bg-neutral border-2 sm:border-4 rounded-lg max-w-5xl w-full max-h-[90vh] overflow-hidden flex flex-col"
        style="border-color: #{currentTeam.team_primary_color};"
        role="document"
        tabindex="-1"
        aria-label="Team details modal"
    >
        <!-- Team Info -->
        <div
            class="border-b-4 p-4 sm:p-6"
            style={`background-color: #${currentTeam.team_primary_color}70; border-color: #${currentTeam.team_primary_color};`}
        >
            <div class="flex items-center justify-between mb-3 sm:mb-4">
                <div class="flex items-center gap-3 sm:gap-4">
                    {#if currentTeam.logo_url}
                        <img 
                            src={currentTeam.logo_url} 
                            alt={currentTeam.team_abbr}
                            class="w-12 sm:w-16 h-12 sm:h-16 object-contain"
                        />
                    {/if}

                    <div>
                        <p class="text-md sm:text-lg font-sans text-black/70 -mb-1">
                            {currentTeam.team_city}
                        </p>
                        <h2 class="text-2xl sm:text-3xl font-heading font-bold text-black mb-1">
                            {currentTeam.team_name}
                        </h2>
                        {#if teamDivision}
                            <p class="text-sm font-sans text-black/60 uppercase">
                                {teamDivision}
                            </p>
                        {/if}
                    </div>
                </div>
                
                <button
                    on:click={closeModal}
                    class="p-2 rounded-lg cursor-pointer transition-colors bg-transparent"
                    style={`color: #${currentTeam.team_primary_color}E6;`}
                    on:mouseenter={(e) => e.currentTarget.style.backgroundColor = `#${currentTeam.team_primary_color}80`}
                    on:mouseleave={(e) => e.currentTarget.style.backgroundColor = `transparent`}
                    aria-label="Close modal"
                >
                    <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                </button>
            </div>

            <!-- Team Stats - 2 rows of 6 -->
            <div class="grid grid-cols-4 sm:grid-cols-6 gap-2 sm:gap-3">
                <div class="rounded-lg p-2 sm:p-3 bg-white/50">
                    <div class="text-xs font-sans font-semibold text-black/60 uppercase">Record</div>
                    <div class="text-sm sm:text-xl font-heading font-bold text-black">
                        {formatRecord(currentTeam.wins, currentTeam.losses)}
                    </div>
                </div>
                <div class="bg-white/50 rounded-lg p-2 sm:p-3">
                    <div class="text-xs font-sans font-semibold text-black/60 uppercase">Seed</div>
                    <div class="text-sm sm:text-xl font-heading font-bold text-black">#{currentTeam.seed}</div>
                </div>
                <div class="bg-white/50 rounded-lg p-2 sm:p-3">
                    <div class="text-xs font-sans font-semibold text-black/60 uppercase">Conf Record</div>
                    <div class="flex flex-row items-baseline gap-0.5 sm:gap-1">
                        <div class="text-sm sm:text-xl font-heading font-bold text-black">
                            {formatRecord(currentTeam.conference_wins, currentTeam.conference_losses)}
                        </div>
                        <div class="text-[10px] sm:text-xs font-sans text-black/60">
                            ({formatGamesBack(currentTeam.conference_games_back)} GB)
                        </div>
                    </div>
                </div>
                <div class="bg-white/50 rounded-lg p-2 sm:p-3">
                    <div class="text-xs font-sans font-semibold text-black/60 uppercase">Div Record</div>
                    <div class="flex flex-row items-baseline gap-0.5 sm:gap-1">
                        <div class="text-sm sm:text-xl font-heading font-bold text-black">
                            {formatRecord(currentTeam.division_wins, currentTeam.division_losses)}
                        </div>
                        <div class="text-[10px] sm:text-xs font-sans text-black/60">
                            ({formatGamesBack(currentTeam.division_games_back)} GB)
                        </div>
                    </div>
                </div>
                <div class="bg-white/50 rounded-lg p-2 sm:p-3">
                    <div class="text-xs font-sans font-semibold text-black/60 uppercase">Home Record</div>
                    <div class="text-sm sm:text-xl font-heading font-bold text-black">
                        {formatRecord(currentTeam.home_wins, currentTeam.home_losses)}
                    </div>
                </div>
                <div class="bg-white/50 rounded-lg p-2 sm:p-3">
                    <div class="text-xs font-sans font-semibold text-black/60 uppercase">Away Record</div>
                    <div class="text-sm sm:text-xl font-heading font-bold text-black">
                        {formatRecord(currentTeam.away_wins, currentTeam.away_losses)}
                    </div>
                </div>
                <div class="bg-white/50 rounded-lg p-2 sm:p-3">
                    <div class="text-xs font-sans font-semibold text-black/60 uppercase">Win %</div>
                    <div class="text-sm sm:text-xl font-heading font-bold text-black">
                        {formatWinPct(currentTeam.win_pct)}
                    </div>
                </div>
                <div class="bg-white/50 rounded-lg p-2 sm:p-3">
                    <div class="text-xs font-sans font-semibold text-black/60 uppercase">Point Diff</div>
                    <div class="text-sm sm:text-xl font-heading font-bold text-black">
                        {formatPointDiff(currentTeam.points_for, currentTeam.points_against, currentTeam.games_with_scores)}
                    </div>
                </div>
                <div class="bg-white/50 rounded-lg p-2 sm:p-3">
                    <div class="text-xs font-sans font-semibold text-black/60 uppercase">PPG</div>
                    <div class="text-sm sm:text-xl font-heading font-bold text-black">
                        {formatPoints(currentTeam.points_for, currentTeam.games_with_scores)}
                    </div>
                </div>
                <div class="bg-white/50 rounded-lg p-2 sm:p-3">
                    <div class="text-xs font-sans font-semibold text-black/60 uppercase">Opp PPG</div>
                    <div class="text-sm sm:text-xl font-heading font-bold text-black">
                        {formatPoints(currentTeam.points_against, currentTeam.games_with_scores)}
                    </div>
                </div>
                <div class="bg-white/50 rounded-lg p-2 sm:p-3">
                    <div class="text-xs font-sans font-semibold text-black/60 uppercase" title="Strength of Schedule">SOS</div>
                    <div class="text-sm sm:text-xl font-heading font-bold text-black">
                        {formatWinPct(currentTeam.strength_of_schedule)}
                    </div>
                </div>
                <div class="bg-white/50 rounded-lg p-2 sm:p-3">
                    <div class="text-xs font-sans font-semibold text-black/60 uppercase" title="Strength of Victory">SOV</div>
                    <div class="text-sm sm:text-xl font-heading font-bold text-black">
                        {formatWinPct(currentTeam.strength_of_victory)}
                    </div>
                </div>
            </div>
        </div>

        <!-- Team Schedule -->
        <div class="flex-1 overflow-y-auto p-6">
            <h3
                class="text-xl font-heading font-bold mb-4 uppercase tracking-wide"
                style={`color: #${currentTeam.team_primary_color};`}
            >
                Season Schedule
            </h3>

            {#if loading}
                <div class="flex items-center justify-center py-12">
                    <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-400"></div>
                </div>
            {:else if error}
                <div class="text-center py-8">
                    <p class="text-red-400 text-lg">{error}</p>
                </div>
            {:else if teamGames.length === 0}
                <div class="text-center py-8">
                    <p class="text-neutral/70 text-lg">No games found for this team</p>
                </div>
            {:else}
                <div class="space-y-2">
                    {#each teamGames as game (game.id)}
                        {@const result = getGameResult(game)}
                        {@const pick = picks.get(game.id)}
                        {@const isHome = isHomeGame(game)}
                        
                        <div
                            class="border-2 rounded-lg p-4 transition-colors"
                            style={`border-color: #${currentTeam.team_primary_color}60;`}
                            on:mouseenter={(e) => e.currentTarget.style.borderColor = `#${currentTeam.team_primary_color}90`}
                            on:mouseleave={(e) => e.currentTarget.style.borderColor = `#${currentTeam.team_primary_color}60`}
                            role="region"
                        >
                            <!-- Game Header -->
                            <div class="flex items-center justify-between mb-3">
                                <div class="flex items-center gap-2">
                                    <span class="text-sm font-sans font-semibold text-black">
                                        Week {game.week}
                                    </span>
                                    <span class="text-xs font-sans text-black/60">
                                        {formatDate(game.start_time)} • {formatTime(game.start_time)}
                                    </span>
                                </div>
                                
                                <div class="flex items-center gap-2">
                                    {#if result === 'win'}
                                        <span class="px-2 py-1 bg-green-500/20 text-green-700 text-xs font-sans font-bold rounded">W</span>
                                    {:else if result === 'loss'}
                                        <span class="px-2 py-1 bg-red-500/20 text-red-700 text-xs font-sans font-bold rounded">L</span>
                                    {/if}
                                    <span class="text-xs font-sans text-black/60">
                                        {isHome ? 'vs' : '@'}
                                    </span>
                                </div>
                            </div>

                            <!-- Game Picker Row (Reusable Component) -->
                            <GamePickerRow
                                {game}
                                {pick}
                                compact={true}
                                sportId={2}
                                hasPlayoffs={playoffState?.is_enabled || false}
                                on:pickChanged={handlePickChange}
                            />
                        </div>
                    {/each}
                </div>
            {/if}
        </div>
    </section>
</div>