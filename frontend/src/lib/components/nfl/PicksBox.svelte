<script lang="ts">
    import { onMount, createEventDispatcher } from 'svelte';
    import { gamesAPI } from '$lib/api/games';
    import { picksAPI } from '$lib/api/picks';
    import { teamsAPI } from '$lib/api/teams';
    import { scenariosAPI } from '$lib/api/scenarios';
    import type { Game, Pick, Team } from '$types';
    
    import WeekNavigator from './WeekNavigator.svelte';
    import GameCard from './GameCard.svelte';
    import ByeTeams from './ByeTeams.svelte';

    export let scenarioId: number;
    export let currentWeek: number;

    const dispatch = createEventDispatcher();

    let games: Game[] = [];
    let picks: Map<number, Pick> = new Map();
    let allTeams: Team[] = [];
    let byeTeams: Team[] = [];
    let seasonId: number | null = null;
    let loading = true;
    let error = '';

    // Group games by day
    $: gamesByDay = groupGamesByDay(games);

    onMount(async () => {
        await loadData();
    });

    // Reload when week changes (but only if we have a seasonId)
    $: if (currentWeek && seasonId) {
        loadGamesForWeek();
    }

    async function loadData() {
        try {
            loading = true;
            error = '';
            
            const scenario = await scenariosAPI.getById(scenarioId);
            seasonId = scenario.season_id;

            allTeams = await teamsAPI.getBySeason(seasonId);

            await loadGamesForWeek();
            await loadPicks();
        } catch (err: any) {
            error = err.response?.data?.error || 'Failed to load games';
            console.error('Error loading data:', err);
        } finally {
            loading = false;
        }
    }

    async function loadGamesForWeek() {
        if (!seasonId) return;
        
        try {
            games = await gamesAPI.getBySeasonAndWeek(seasonId, currentWeek);
            console.log('Loaded games for week', currentWeek, ':', games);
            calculateByeTeams();
        } catch (err: any) {
            console.error('Error loading games:', err);
            error = 'Failed to load games for this week';
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
        
        games.forEach(game => {
            const day = game.day_of_week || 'Unknown';
            if (!grouped.has(day)) {
                grouped.set(day, []);
            }
            grouped.get(day)!.push(game);
        });

        // Define the correct NFL week order
        const dayOrder = ['Thursday', 'Friday', 'Saturday', 'Sunday', 'Monday', 'Tuesday', 'Wednesday'];
        const sortedMap = new Map<string, Game[]>();
        
        // Add days in the correct order, but only if they have games
        dayOrder.forEach(day => {
            if (grouped.has(day)) {
                sortedMap.set(day, grouped.get(day)!);
            }
        });

        return sortedMap;
    }

    function handleWeekChange(event: CustomEvent<{ week: number }>) {
        currentWeek = event.detail.week;
        dispatch('weekChanged', { week: currentWeek });
    }

    async function handlePickChange(event: CustomEvent<{
        gameId: number;
        pickedTeamId?: number | null;
        predictedHomeScore?: number;
        predictedAwayScore?: number;
    }>) {
        try {
            const { gameId, pickedTeamId, predictedHomeScore, predictedAwayScore } = event.detail;
            
            const existingPick = picks.get(gameId);

            let updated: Pick;
            if (existingPick) {
                updated = await picksAPI.update(scenarioId, gameId, {
                    picked_team_id: pickedTeamId === undefined ? null : pickedTeamId,
                    predicted_home_score: predictedHomeScore,
                    predicted_away_score: predictedAwayScore
                });
            } else {
                updated = await picksAPI.create(scenarioId, gameId, {
                    picked_team_id: pickedTeamId === undefined ? null : pickedTeamId,
                    predicted_home_score: predictedHomeScore,
                    predicted_away_score: predictedAwayScore
                });
            }

            picks.set(gameId, updated);
            dispatch('pickUpdated');
            picks = new Map(picks);
        } catch (err: any) {
            console.error('Error saving pick:', err);
            alert('Failed to save pick. Please try again.');
        }
    }
</script>

<div class="bg-neutral border-2 border-primary-700 rounded-lg p-3 sm:p-4 md:p-6 w-full max-w-full">
    <!-- Week Navigator -->
    <WeekNavigator 
        {currentWeek}
        maxWeek={18}
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
        <!-- Games List -->
        <div class="mt-4 md:mt-6 space-y-4 md:space-y-6">
            {#each [...gamesByDay.entries()] as [day, dayGames]}
                <!-- Day Header -->
                <div>
                    <h3 class="text-lg sm:text-xl font-heading font-bold text-primary-700 mb-2 md:mb-3 uppercase tracking-wide">
                        {day}
                    </h3>
                    
                    <!-- Games for this day - Responsive Grid -->
                    <div class="grid grid-cols-2 gap-2 md:gap-3">
                        {#each dayGames as game (game.id)}
                            <GameCard 
                                {game}
                                pick={picks.get(game.id)}
                                on:pickChanged={handlePickChange}
                            />
                        {/each}
                    </div>
                </div>
            {/each}
        </div>

        <!-- Bye Teams -->
        {#if byeTeams.length > 0}
            <ByeTeams teams={byeTeams} />
        {/if}
    {/if}
</div>