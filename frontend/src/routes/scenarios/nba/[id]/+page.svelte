<script lang="ts">
    import { onMount } from 'svelte';
    import { page } from '$app/stores';
    import { scenariosAPI } from '$lib/api/scenarios';
    import { standingsAPI } from '$lib/api/standings';
    import { gamesAPI } from '$lib/api/games';
    import { playoffsAPI } from '$lib/api/playoffs';
    import { getCurrentNBAWeekFromGames } from '$lib/utils/nba/dates';
    import type { Scenario, NBAStandings, NBAPlayoffSeed, PlayoffState, Game } from '$types';
    import ScenarioHeader from '$lib/components/scenarios/ScenarioHeader.svelte';
    import ScenarioSettings from '$lib/components/scenarios/ScenarioSettings.svelte';
    import ScenarioInfo from '$lib/components/scenarios/ScenarioInfo.svelte';
    import PicksBox from '$lib/components/scenarios/PicksBox.svelte';
    import StandingsBox from '$lib/components/nba/StandingsBox.svelte';
    import StandingsBoxExpanded from '$lib/components/nba/StandingsBoxExpanded.svelte';
    import DraftOrderBox from '$lib/components/scenarios/DraftOrderBox.svelte';
    import PlayoffPicksBox from '$lib/components/nba/PlayoffPicksBox.svelte';
    import TeamModal from '$lib/components/nba/TeamModal.svelte';

    // State variables for scenario
    let scenarioId: number;
    let scenario: Scenario | null = null;

    // State variable for standings
    let standings: NBAStandings | null = null;

    // State variables for playoffs
    let playoffState: PlayoffState | null = null;
    let canEnablePlayoffs: boolean = false;

    // State variables for modals
    let showSettings = false;
    let showInfo = false;

    // State variables for current week/round
    let currentWeek = 1;
    let currentPlayoffRound = 1;

    // State variable for games
    let allGames: Game[] = [];
    
    // Loading and error states
    let loading = true;
    let error = '';

    // State variables for view modes
    // Single source of truth for view mode
    let viewKey = 'regular-1'; // Format: 'regular-{week}' or 'playoff-{round}'
    type ViewMode = 'conference' | 'division';
    let standingsViewMode: ViewMode = 'conference';

    // State variable for selected team in modal
    let selectedTeam: NBAPlayoffSeed | null = null;

    // State variable for save status
    let saveStatus: 'idle' | 'saving' | 'saved' | 'error' = 'idle';

    // State variables for picks box references
    let desktopPicksBox: PicksBox | PlayoffPicksBox;
    let mobilePicksBox: PicksBox | PlayoffPicksBox;

    // Get scenario ID from URL params
    $: scenarioId = parseInt($page.params.id ?? '0');

    // Load data on mount
    onMount(async () => {
        await loadScenario();
        await loadStandings();
        await loadPlayoffState();
    });

    async function loadScenario() {
        try {
            loading = true;
            scenario = await scenariosAPI.getById(scenarioId);

            // Load all games for the scenario's season to determine current week
            if (scenario.season_id) {
                allGames = await gamesAPI.getBySeason(scenario.season_id);
                currentWeek = getCurrentNBAWeekFromGames(allGames);
                viewKey = `regular-${currentWeek}`;
            }
        } catch (err: any) {
            error = err.response?.data?.error || 'Failed to load scenario.';
        } finally {
            loading = false;
        }
    }

    async function loadStandings() {
        try {
            standings = await standingsAPI.getByNBAScenario(scenarioId);
        } catch (err: any) {
            console.error('Failed to load standings:', err);
        }
    }

    async function loadPlayoffState() {
        try {
            const response = await playoffsAPI.getState(scenarioId);
            playoffState = response.playoff_state;
            canEnablePlayoffs = response.can_enable;

            // If playoffs are enabled, switch to playoff view
            if (playoffState?.is_enabled) {
                currentPlayoffRound = playoffState.current_round;
                viewKey = `playoff-${currentPlayoffRound}`;
            }
        } catch (err: any) {
            console.error('Failed to load playoff state:', err);
        }
    }

    async function handleEnablePlayoffs() {
        try {
            saveStatus = 'saving';
            await playoffsAPI.enable(scenarioId);
            await loadPlayoffState();
            await loadStandings();
            saveStatus = 'saved';
            setTimeout(() => saveStatus = 'idle', 2000);
        } catch (err: any) {
            saveStatus = 'error';
            alert('Failed to enable playoffs: ' + (err.response?.data?.error || err.message));
            setTimeout(() => saveStatus = 'idle', 2000);
        }
    }

    function handleScenarioUpdated(event: CustomEvent) {
        scenario = event.detail;
        saveStatus = 'saved';
        setTimeout(() => saveStatus = 'idle', 2000);
    }

    async function handleWeekChange(event: CustomEvent<{ week: number }>) {
        const newWeek = event.detail.week;
        
        // If week is greater than 25, it's a playoff round
        if (newWeek > 25) {
            const newRound = newWeek - 25;
            
            // If trying to advance to a round that doesn't exist yet, generate it
            if (playoffState && newRound > playoffState.current_round) {
                try {
                    await playoffsAPI.generateNextRound(scenarioId);
                    await loadPlayoffState(); // Reload to get updated current_round
                } catch (err: any) {
                    alert('Failed to generate next round: ' + (err.response?.data?.error || err.message));
                    return;
                }
            }
            
            // Update current playoff round and view key
            currentPlayoffRound = newRound;
            viewKey = `playoff-${currentPlayoffRound}`;
        } else {
            // Update current week and view key
            currentWeek = newWeek;
            viewKey = `regular-${currentWeek}`;
        }
    }

    async function handlePickUpdated() {
        saveStatus = 'saved';
        setTimeout(() => saveStatus = 'idle', 2000);
        await loadPlayoffState();
        await loadStandings();

        // Reload picks in both PicksBox instances
        if (desktopPicksBox && 'reloadPicks' in desktopPicksBox) {
            await desktopPicksBox.reloadPicks();
        }
        if (mobilePicksBox && 'reloadPicks' in mobilePicksBox) {
            await mobilePicksBox.reloadPicks();
        }
    }

    function handleOpenTeamModal(event: CustomEvent<{ team: NBAPlayoffSeed }>) {
        selectedTeam = event.detail.team;
    }

    function handleCloseTeamModal() {
        selectedTeam = null;
    }
</script>

<svelte:head>
    <title>{scenario?.name || 'Scenario'} - GameScript</title>
</svelte:head>

{#if loading}
    <div class="flex items-center justify-center min-h-full">
        <div class="animate-spin rounded-full h-16 w-16 border-b-2 border-primary-400"></div>
    </div>
{:else if error}
    <div class="text-center">
        <p class="text-red-400 text-xl mb-4">{error}</p>
        <a href="/" class="text-primary-400 hover:text-primary-300">
            ← Back to Home
        </a>
    </div>
{:else if scenario}
    <!-- Modals -->
    <ScenarioSettings 
        bind:isOpen={showSettings}
        {scenario}
        on:updated={handleScenarioUpdated}
    />
    <ScenarioInfo bind:isOpen={showInfo} sport={scenario.sport_short_name} />

    <!-- Header -->
    <ScenarioHeader 
        {scenario}
        {saveStatus}
        on:openSettings={() => showSettings = true}
        on:openInfo={() => showInfo = true}
    />

    <!-- Enable Playoffs Banner -->
    {#if canEnablePlayoffs && !playoffState?.is_enabled}
        <div class="mt-4 bg-primary-600 border-2 border-primary-500 rounded-lg p-3 sm:p-4">
            <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3 sm:gap-4">
                <div class="flex-1 min-w-0">
                    <h3 class="text-base sm:text-lg font-heading font-bold text-neutral mb-1">
                        Regular Season Complete!
                    </h3>
                    <p class="text-sm sm:text-base text-neutral/80 font-sans">
                        All regular season games have been completed. Ready to start the playoffs?
                    </p>
                </div>
                <button
                    on:click={handleEnablePlayoffs}
                    class="w-full sm:w-auto shrink-0 px-4 sm:px-6 py-2 sm:py-3 bg-neutral hover:bg-neutral/90 text-primary-950 font-heading font-bold text-base sm:text-lg rounded-lg transition-colors cursor-pointer whitespace-nowrap"
                >
                    Start Playoffs
                </button>
            </div>
        </div>
    {/if}

    <!-- Main Content -->
    <div class="mt-6 space-y-6 lg:space-y-0">
        <!-- Desktop: 3-column layout with fixed-width standings -->
        <div class="hidden lg:grid lg:grid-cols-[200px_1fr_200px] lg:gap-6">
            <!-- Left: East Standings -->
            <div class="w-[200px]">
                {#if standings && scenario.season_id}
                    <StandingsBox 
                        standings={standings.eastern} 
                        conference="East"
                        bind:viewMode={standingsViewMode}
                        on:openTeamModal={handleOpenTeamModal}
                    />
                {/if}
            </div>

            <!-- Center: Picks -->
            <div class="min-w-[500px]">
                {#key viewKey}
                    {#if viewKey.startsWith('playoff-') && playoffState?.is_enabled && scenario.season_id}
                        <PlayoffPicksBox
                            bind:this={desktopPicksBox}
                            {scenarioId}
                            {playoffState}
                            currentRound={currentPlayoffRound}
                            seasonId={scenario.season_id}
                            on:weekChanged={handleWeekChange}
                            on:pickUpdated={handlePickUpdated}
                        />
                    {:else}
                        <PicksBox 
                            bind:this={desktopPicksBox}
                            {scenarioId}
                            {currentWeek}
                            {playoffState}
                            {canEnablePlayoffs}
                            on:weekChanged={handleWeekChange}
                            on:pickUpdated={handlePickUpdated}
                        />
                    {/if}
                {/key}
            </div>

            <!-- Right: West Standings -->
            <div class="w-[200px]">
                {#if standings && scenario.season_id}
                    <StandingsBox 
                        standings={standings.western} 
                        conference="West"
                        bind:viewMode={standingsViewMode}
                        on:openTeamModal={handleOpenTeamModal}
                    />
                {/if}
            </div>
        </div>

        <!-- Mobile: Stacked layout -->
        <div class="lg:hidden space-y-6">
            <!-- Picks -->
            {#key viewKey}
                {#if viewKey.startsWith('playoff-') && playoffState?.is_enabled && scenario.season_id}
                    <PlayoffPicksBox
                        bind:this={mobilePicksBox}
                        {scenarioId}
                        {playoffState}
                        currentRound={currentPlayoffRound}
                        seasonId={scenario.season_id}
                        on:weekChanged={handleWeekChange}
                        on:pickUpdated={handlePickUpdated}
                    />
                {:else}
                    <PicksBox 
                        bind:this={mobilePicksBox}
                        {scenarioId}
                        {currentWeek}
                        {playoffState}
                        {canEnablePlayoffs}
                        on:weekChanged={handleWeekChange}
                        on:pickUpdated={handlePickUpdated}
                    />
                {/if}
            {/key}

            <!-- Standings -->
            {#if standings && scenario.season_id}
                <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <StandingsBoxExpanded 
                        standings={standings.eastern} 
                        conference="East"
                        bind:viewMode={standingsViewMode}
                        on:openTeamModal={handleOpenTeamModal}
                    />
                    <StandingsBoxExpanded 
                        standings={standings.western} 
                        conference="West"
                        bind:viewMode={standingsViewMode}
                        on:openTeamModal={handleOpenTeamModal}
                    />
                </div>
            {/if}
        </div>

        <!-- Draft Order -->
        <!-- {#if standings}
            <div class="w-full mt-6">
                <DraftOrderBox 
                    draftOrder={standings.draft_order}
                />
            </div>
        {/if} -->
    </div>
{/if}

<!-- Team Modal -->
{#if selectedTeam && scenario}
    <TeamModal 
        team={selectedTeam}
        {scenarioId}
        seasonId={scenario.season_id}
        on:close={handleCloseTeamModal}
        on:pickUpdated={handlePickUpdated}
    />
{/if}