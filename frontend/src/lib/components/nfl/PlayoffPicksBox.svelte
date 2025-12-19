<script lang="ts">
    import { onMount, createEventDispatcher } from 'svelte';
    import { playoffsAPI } from '$lib/api/playoffs';
    import { standingsAPI } from '$lib/api/standings';
    import { gamesAPI } from '$lib/api/games';
    import type { PlayoffMatchup, PlayoffState, Standings, Game } from '$types';
    import { PLAYOFF_ROUND_NAMES } from '$types';
    
    import WeekNavigator from './WeekNavigator.svelte';
    import PlayoffGameCard from './PlayoffGameCard.svelte';

    export let scenarioId: number;
    export let playoffState: PlayoffState;
    export let currentRound: number;
    export let seasonId: number;

    const dispatch = createEventDispatcher();

    let matchups: PlayoffMatchup[] = [];
    let standings: Standings | null = null;
    let allGames: Game[] = [];
    let loading = true;
    let error = '';

    // Convert round to week for WeekNavigator (Week 19 = Round 1, Week 20 = Round 2, etc.)
    $: currentWeek = 18 + currentRound;

    $: afcMatchups = matchups.filter(m => m.conference === 'AFC');
    $: nfcMatchups = matchups.filter(m => m.conference === 'NFC');
    $: superBowlMatchup = currentRound === 4 ? matchups[0] : null;

    // Get bye teams for Wild Card round
    $: byeTeams = currentRound === 1 && standings ? [
        standings.afc.playoff_seeds[0], // AFC #1 seed
        standings.nfc.playoff_seeds[0]  // NFC #1 seed
    ] : [];

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
            standings = await standingsAPI.getByScenario(scenarioId);
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
            matchups = await playoffsAPI.getMatchups(scenarioId, currentRound);
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
        matchupId: number;
        pickedTeamId?: number | null;
        predictedHigherScore?: number;
        predictedLowerScore?: number;
    }>) {
        try {
            const { matchupId, pickedTeamId, predictedHigherScore, predictedLowerScore } = event.detail;
            
            // Optimistically update local state first
            matchups = matchups.map(m => {
                if (m.id === matchupId) {
                    return {
                        ...m,
                        picked_team_id: pickedTeamId,
                        predicted_higher_seed_score: predictedHigherScore,
                        predicted_lower_seed_score: predictedLowerScore
                    };
                }
                return m;
            });

            // Then update on server
            await playoffsAPI.updatePick(scenarioId, matchupId, {
                picked_team_id: pickedTeamId === undefined ? null : pickedTeamId,
                predicted_higher_seed_score: predictedHigherScore,
                predicted_lower_seed_score: predictedLowerScore
            });

            // Only dispatch pickUpdated (parent will handle reloading standings/playoff state)
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
    {:else if matchups.length === 0}
        <div class="text-center py-8">
            <p class="text-neutral/70 text-lg">Complete all picks from the previous round to unlock this round.</p>
        </div>
    {:else}
        <!-- Matchups -->
        <div class="mt-4 md:mt-6 space-y-4 md:space-y-6">
            <!-- Super Bowl (no conferences) -->
            {#if currentRound === 4 && superBowlMatchup}
                <div>
                    <h3 class="text-lg sm:text-xl font-heading font-bold text-primary-700 mb-2 md:mb-3 uppercase tracking-wide text-center">
                        Championship
                    </h3>
                    
                    <div class="max-w-md mx-auto">
                        <PlayoffGameCard 
                            matchup={superBowlMatchup}
                            hasLaterRounds={false}
                            on:pickChanged={handlePickChange}
                        />
                    </div>
                </div>
            {:else}
                <!-- AFC Conference -->
                {#if afcMatchups.length > 0}
                    <div>
                        <h3 class="text-lg sm:text-xl font-heading font-bold mb-2 md:mb-3 uppercase tracking-wide"
                            style="color: #C8102E">
                            AFC {PLAYOFF_ROUND_NAMES[currentRound]}
                        </h3>
                        
                        <div class="grid grid-cols-1 sm:grid-cols-2 gap-2 md:gap-3">
                            {#each afcMatchups as matchup (matchup.id)}
                                <PlayoffGameCard 
                                    {matchup}
                                    hasLaterRounds={playoffState.current_round > currentRound}
                                    on:pickChanged={handlePickChange}
                                />
                            {/each}
                        </div>
                    </div>
                {/if}

                <!-- NFC Conference -->
                {#if nfcMatchups.length > 0}
                    <div>
                        <h3 class="text-lg sm:text-xl font-heading font-bold mb-2 md:mb-3 uppercase tracking-wide"
                            style="color: #013369">
                            NFC {PLAYOFF_ROUND_NAMES[currentRound]}
                        </h3>
                        
                        <div class="grid grid-cols-1 sm:grid-cols-2 gap-2 md:gap-3">
                            {#each nfcMatchups as matchup (matchup.id)}
                                <PlayoffGameCard 
                                    {matchup}
                                    hasLaterRounds={playoffState.current_round > currentRound}
                                    on:pickChanged={handlePickChange}
                                />
                            {/each}
                        </div>
                    </div>
                {/if}
            {/if}

            <!-- Bye Teams (Wild Card Round Only) -->
            {#if currentRound === 1 && byeTeams.length > 0}
                <div class="mt-6 pt-6 border-t-2 border-primary-700">
                    <h3 class="text-xl font-heading font-bold text-primary-700 mb-3 uppercase tracking-wide">
                        First Round Bye
                    </h3>
                    
                    <div class="grid grid-cols-2 gap-4">
                        {#each byeTeams as team}
                            <div class="border-2 rounded-lg p-4"
                                style="background-color: #{team.team_primary_color}90; border-color: #{team.team_primary_color};">
                                <div class="flex items-center gap-3">
                                    {#if team.logo_url}
                                        <img 
                                            src={team.logo_url} 
                                            alt={team.team_abbr}
                                            class="w-12 h-12 object-contain"
                                        />
                                    {/if}
                                    <div class="flex-1">
                                        <div
                                            class="text-sm font-sans font-semibold"
                                            style="color: #{team.team_primary_color}"
                                        >
                                            {team.team_city}
                                        </div>
                                        <div
                                            class="text-lg font-heading font-bold"
                                            style="color: #{team.team_primary_color}"
                                        >
                                            {team.team_name}
                                        </div>
                                    </div>
                                </div>
                            </div>
                        {/each}
                    </div>
                </div>
            {/if}
        </div>
    {/if}
</div>