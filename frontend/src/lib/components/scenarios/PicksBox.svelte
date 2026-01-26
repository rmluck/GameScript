<script lang="ts">
    import { onMount, createEventDispatcher } from 'svelte';
    import { gamesAPI } from '$lib/api/games';
    import { picksAPI } from '$lib/api/picks';
    import { teamsAPI } from '$lib/api/teams';
    import { scenariosAPI } from '$lib/api/scenarios';
    import type { Game, Pick, Team, PlayoffState } from '$types';
    import WeekNavigator from './WeekNavigator.svelte';
    import GameCard from './GameCard.svelte';
    import ByeTeams from './ByeTeams.svelte';

    // Props
    export let scenarioId: number;
    export let playoffState: PlayoffState | null = null;
    export let currentWeek: number;
    export let canEnablePlayoffs: boolean = false;

    // Event dispatcher
    const dispatch = createEventDispatcher();

    // State variables for scenario
    let seasonId: number | null = null;
    let sportId: number | null = null;

    // State variable for games
    let games: Game[] = [];
    let allGames: Game[] = [];

    // State variable for picks
    let picks: Map<number, Pick> = new Map();

    // State variable for teams
    let allTeams: Team[] = [];
    let byeTeams: Team[] = [];

    // Loading and error states
    let loading = true;
    let error = '';

    // Group games by day
    $: gamesByDay = groupGamesByDay(games);

    // Filter games when week changes
    $: if (currentWeek && allGames.length > 0) {
        games = allGames.filter(game => game.week === currentWeek);
        if (sportId === 1) {
            calculateByeTeams();
        }
    }

    // Update box height based on sport
    $: boxHeight = sportId === 1 ? '119.5vh' : sportId === 2 ? '114vh' : '100vh';

    // Load data on mount
    onMount(async () => {
        await loadData();
    });

    // Load all necessary data
    async function loadData() {
        try {
            loading = true;
            error = '';
            
            const scenario = await scenariosAPI.getById(scenarioId);
            seasonId = scenario.season_id;
            sportId = scenario.sport_id;

            // Load all teams
            allTeams = await teamsAPI.getBySeason(seasonId);

            // Load all games for the season once
            allGames = await gamesAPI.getBySeason(seasonId);
            
            // Filter to current week
            games = allGames.filter(game => game.week === currentWeek);
            
            // Calculate bye teams for NFL
            if (sportId === 1) {
                calculateByeTeams();
            }

            await loadPicks();
        } catch (err: any) {
            error = err.response?.data?.error || 'Failed to load games';
            console.error('Error loading data:', err);
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

    // Export function to reload picks from parent
    export async function reloadPicks() {
        await loadPicks();
    }

    function calculateByeTeams() {
        const playingTeamIds = new Set<number>();
        games.forEach(game => {
            playingTeamIds.add(game.home_team_id);
            playingTeamIds.add(game.away_team_id);
        });

        byeTeams = allTeams.filter(team => !playingTeamIds.has(team.id));
    }

    function groupGamesByDay(games: Game[]): Map<string, Game[]> {
        const grouped = new Map<string, Game[]>();
        
        // Group games by day of week
        games.forEach(game => {
            const day = game.day_of_week || 'Unknown';
            if (!grouped.has(day)) {
                grouped.set(day, []);
            }
            grouped.get(day)!.push(game);
        });

        // Sort days according to sport-specific order
        let dayOrder: string[] = [];
        if (sportId === 1) {
            dayOrder = ['Thursday', 'Friday', 'Saturday', 'Sunday', 'Monday', 'Tuesday', 'Wednesday'];
        } else if (sportId === 2) {
            dayOrder = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'];
        }

        const sortedMap = new Map<string, Game[]>();
        dayOrder.forEach((day: string) => {
            if (grouped.has(day)) {
                sortedMap.set(day, grouped.get(day)!);
            }
        });

        return sortedMap;
    }

    // Handle pick changes from child components
    async function handlePickChange(event: CustomEvent<{
        gameId: number;
        pickedTeamId?: number | null;
        predictedHomeScore?: number;
        predictedAwayScore?: number;
        deletePick?: boolean;
    }>) {
        try {
            const { gameId, pickedTeamId, predictedHomeScore, predictedAwayScore, deletePick } = event.detail;
            
            const existingPick = picks.get(gameId);

            if (deletePick && existingPick) {
                // Delete existing pick
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
            dispatch('pickUpdated');
        } catch (err: any) {
            console.error('Error saving pick:', err);
            alert('Failed to save pick. Please try again.');
        }
    }

    function handleWeekChange(event: CustomEvent<{ week: number }>) {
        currentWeek = event.detail.week;
        // Dispatch up to parent
        dispatch('weekChanged', { week: currentWeek });
    }
</script>

<div class="bg-neutral border-2 border-primary-700 rounded-lg p-3 sm:p-4 md:p-6 w-full max-w-full flex flex-col overflow-hidden" style="height: {boxHeight};">
    <!-- Week Navigator -->
    <WeekNavigator 
        {currentWeek}
        {allGames}
        {playoffState}
        {canEnablePlayoffs}
        sportId={sportId}
        isCurrentRoundComplete={false}
        on:weekChanged={handleWeekChange}
    />

    {#if loading}
        <div class="flex items-center justify-center py-12">
            <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-400"></div>
        </div>
    {:else if error}
        <div class="text-center py-8">
            <p class="text-red-400 text-lg">{error}</p>
        </div>
    {:else if games.length === 0}
        <div class="text-center py-8">
            <p class="text-neutral/70 text-lg">No games scheduled for Week {currentWeek}</p>
        </div>
    {:else}
        <!-- Games -->
        <div class="mt-4 md:mt-6 space-y-4 overflow-y-auto flex-1 pr-2">
            {#each [...gamesByDay.entries()] as [day, dayGames]}
                <div>
                    <h3 class="text-lg sm:text-xl font-heading font-bold text-primary-700 mb-2 md:mb-3 uppercase tracking-wide">
                        {day}
                    </h3>
                    
                    <div class="grid grid-cols-1 sm:grid-cols-2 space-y-2 sm:space-y-3">
                        {#each dayGames as game (game.id)}
                            <GameCard 
                                {game}
                                sportId={sportId}
                                pick={picks.get(game.id)}
                                hasPlayoffs={playoffState?.is_enabled || false}
                                on:pickChanged={handlePickChange}
                            />
                        {/each}
                    </div>
                </div>
            {/each}
        </div>

        <!-- Bye Teams (for NFL) -->
        {#if byeTeams.length > 0}
            <ByeTeams teams={byeTeams} />
        {/if}
    {/if}
</div>