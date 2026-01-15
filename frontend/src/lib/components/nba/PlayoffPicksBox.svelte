<script lang="ts">
    import { onMount, createEventDispatcher } from 'svelte';
    import { playoffsAPI } from '$lib/api/playoffs';
    import { standingsAPI } from '$lib/api/standings';
    import { gamesAPI } from '$lib/api/games';
    import type { PlayoffMatchup, PlayoffSeries, PlayoffState, NBAStandings, Game } from '$types';
    import { NBA_PLAYOFF_ROUND_NAMES, NBA_PLAYOFF_ROUNDS } from '$types';

    import WeekNavigator from '../scenarios/WeekNavigator.svelte';
    import PlayoffGameCard from './PlayoffGameCard.svelte';

    export let scenarioId: number;
    export let playoffState: PlayoffState;
    export let currentRound: number;
    export let seasonId: number;

    const dispatch = createEventDispatcher();

    let items: (PlayoffMatchup | PlayoffSeries)[] = [];
    let standings: NBAStandings | null = null;
    let allGames: Game[] = [];
    let loading = true;
    let error = '';

    // Convert round to week for WeekNavigator
    $: currentWeek = 25 + currentRound;

    // Determine if current round uses series or single games
    $: isSeries = currentRound >= NBA_PLAYOFF_ROUNDS.CONFERENCE_QUARTERFINALS;

    $: isCurrentRoundComplete = items.length > 0 && items.every(m => m.picked_team_id != null);

    $: eastItems = items.filter(m => m.conference === 'Eastern');
    $: westItems = items.filter(m => m.conference === 'Western');
    $: finalsItem = currentRound === NBA_PLAYOFF_ROUNDS.NBA_FINALS ? items[0] : null;

    onMount(async () => {
        await loadData();
    });

    $: if (currentRound) {
        loadMatchups();
    }

    async function loadData() {
        await loadStandings();
        await loadGames();
        await loadMatchups();
    }

    async function loadStandings() {
        try {
            standings = await standingsAPI.getByNBAScenario(scenarioId);
        } catch (err: any) {
            console.error('Failed to load standings:', err);
        }
    }

    async function loadGames() {
        try {
            allGames = await gamesAPI.getBySeason(seasonId);
        } catch (err: any) {
            console.error('Failed to load games:', err);
        }
    }

    async function loadMatchups() {
        try {
            loading = true;
            error = '';
            items = await playoffsAPI.getMatchups(scenarioId, currentRound) as (PlayoffMatchup | PlayoffSeries)[];
        } catch (err: any) {
            error = err.response?.data?.error || 'Failed to load playoff matchups';
            console.error('Error loading matchups:', err);
        } finally {
            loading = false;
        }
    }

    // Export function to reload picks from parent
    export async function reloadPicks() {
        await loadMatchups();
    }

    function handleWeekChange(event: CustomEvent<{ week: number }>) {
        const newWeek = event.detail.week;
        // Dispatch up to parent which will handle switching between regular season and playoffs
        dispatch('weekChanged', { week: newWeek });
    }

    async function handlePickChange(event: CustomEvent<{
        itemId?: number;
        matchupId?: number;
        isSeries?: boolean;
        pickedTeamId?: number | null;
        predictedHigherScore?: number;
        predictedLowerScore?: number;
        predictedHigherWins?: number;
        predictedLowerWins?: number;
    }>) {
        try {
            const {
                itemId,
                matchupId,
                isSeries: isSeriesPick,
                pickedTeamId,
                predictedHigherScore,
                predictedLowerScore,
                predictedHigherWins,
                predictedLowerWins
            } = event.detail;

            const id = itemId || matchupId;
            if (!id) return;

            // Optimistically update local items
            items = items.map(item => {
                if (item.id === id) {
                    if (isSeriesPick && 'best_of' in item) {
                        return {
                            ...item,
                            picked_team_id: pickedTeamId === undefined ? null : pickedTeamId,
                            predicted_higher_seed_wins: predictedHigherWins,
                            predicted_lower_seed_wins: predictedLowerWins
                        } as PlayoffSeries;
                    } else {
                        return {
                            ...item,
                            picked_team_id: pickedTeamId === undefined ? null : pickedTeamId,
                            predicted_higher_seed_score: predictedHigherScore,
                            predicted_lower_seed_score: predictedLowerScore
                        } as PlayoffMatchup;
                    }
                }
                return item;
            });

            // Update on server
            if (isSeriesPick) {
                await playoffsAPI.updatePick(scenarioId, id, {
                    picked_team_id: pickedTeamId === undefined ? null : pickedTeamId,
                    predicted_higher_seed_wins: predictedHigherWins,
                    predicted_lower_seed_wins: predictedLowerWins
                });
            } else {
                await playoffsAPI.updatePick(scenarioId, id, {
                    picked_team_id: pickedTeamId === undefined ? null : pickedTeamId,
                    predicted_higher_seed_score: predictedHigherScore,
                    predicted_lower_seed_score: predictedLowerScore
                });
            }

            // Only dispatch pickUpdated (DON'T auto-advance to next round)
            dispatch('pickUpdated');
        } catch (err: any) {
            console.error('Error saving pick:', err);
            // Revert optimistic update on error
            await loadMatchups();
            alert('Failed to save pick. Please try again.');
        }
    }
</script>

<div class="bg-neutral border-2 border-primary-700 rounded-lg p-3 sm:p-4 md:p-6 w-full max-w-full">
    <!-- Week Navigator -->
    {#if allGames.length > 0}
        <WeekNavigator 
            currentWeek={currentWeek}
            {allGames}
            {playoffState}
            canEnablePlayoffs={false}
            sportId={2}
            {isCurrentRoundComplete}
            on:weekChanged={handleWeekChange}
        />
    {/if}

    {#if loading}
        <div class="flex items-center justify-center py-12">
            <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-400"></div>
        </div>
    {:else if error}
        <div class="text-center py-8">
            <p class="text-red-400 text-lg">{error}</p>
        </div>
    {:else if items.length === 0}
        <div class="text-center py-8">
            <p class="text-neutral/70 text-lg">Complete all picks from the previous round to unlock this round.</p>
        </div>
    {:else}
        <!-- Matchups -->
        <div class="mt-4 md:mt-6 space-y-4 md:space-y-6">
            <!-- NBA Finals (no conferences) -->
            {#if currentRound === NBA_PLAYOFF_ROUNDS.NBA_FINALS && finalsItem}
                <div>
                    <h3 class="text-lg sm:text-xl font-heading font-bold text-primary-700 mb-2 md:mb-3 uppercase tracking-wide text-center">
                        NBA Finals
                    </h3>
                    
                    <div class="max-w-md mx-auto">
                        <PlayoffGameCard 
                            item={finalsItem}
                            hasLaterRounds={false}
                            on:pickChanged={handlePickChange}
                        />
                    </div>
                </div>
            {:else}
                <!-- Eastern Conference -->
                {#if eastItems.length > 0}
                    <div>
                        <h3 class="text-lg sm:text-xl font-heading font-bold mb-2 md:mb-3 uppercase tracking-wide"
                            style="color: #C8102E">
                            East {NBA_PLAYOFF_ROUND_NAMES[currentRound]}
                        </h3>
                        
                        <div class="grid grid-cols-1 sm:grid-cols-2 gap-2 md:gap-3">
                            {#each eastItems as item (item.id)}
                                <PlayoffGameCard 
                                    {item}
                                    hasLaterRounds={playoffState.current_round > currentRound}
                                    on:pickChanged={handlePickChange}
                                />
                            {/each}
                        </div>
                    </div>
                {/if}

                <!-- Western Conference -->
                {#if westItems.length > 0}
                    <div>
                        <h3 class="text-lg sm:text-xl font-heading font-bold mb-2 md:mb-3 uppercase tracking-wide"
                            style="color: #013369">
                            West {NBA_PLAYOFF_ROUND_NAMES[currentRound]}
                        </h3>
                        
                        <div class="grid grid-cols-1 sm:grid-cols-2 gap-2 md:gap-3">
                            {#each westItems as item (item.id)}
                                <PlayoffGameCard 
                                    {item}
                                    hasLaterRounds={playoffState.current_round > currentRound}
                                    on:pickChanged={handlePickChange}
                                />
                            {/each}
                        </div>
                    </div>
                {/if}
            {/if}
        </div>
    {/if}
</div>