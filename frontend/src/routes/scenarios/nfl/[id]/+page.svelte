<script lang="ts">
    import { onMount } from 'svelte';
    import { page } from '$app/stores';
    import { scenariosAPI } from '$lib/api/scenarios';
    import { standingsAPI } from '$lib/api/standings';
    import { gamesAPI } from '$lib/api/games';
    import type { Scenario, Standings } from '$types';

    import ScenarioHeader from '$lib/components/scenarios/ScenarioHeader.svelte';
    import ScenarioSettings from '$lib/components/scenarios/ScenarioSettings.svelte';
    import ScenarioInfo from '$lib/components/nfl/ScenarioInfo.svelte';
    import PicksBox from '$lib/components/nfl/PicksBox.svelte';
    import StandingsBox from '$lib/components/nfl/StandingsBox.svelte';
    import StandingsBoxExpanded from '$lib/components/nfl/StandingsBoxExpanded.svelte';
    import DraftOrderBox from '$lib/components/nfl/DraftOrderBox.svelte';
    import TeamModal from '$lib/components/nfl/TeamModal.svelte';
    import { getCurrentNFLWeekFromGames } from '$lib/utils/nfl/dates';
    import type { PlayoffSeed } from '$types';

    let scenarioId: number;
    let scenario: Scenario | null = null;
    let standings: Standings | null = null;
    let loading = true;
    let error = '';

    let showSettings = false;
    let showInfo = false;

    let currentWeek = 1;

    type ViewMode = 'conference' | 'division';
    let standingsViewMode: ViewMode = 'conference';

    let selectedTeam: PlayoffSeed | null = null;

    let saveStatus: 'idle' | 'saving' | 'saved' | 'error' = 'idle';

    $: scenarioId = parseInt($page.params.id ?? '0');

    onMount(async () => {
        await loadScenario();
        await loadStandings();
    });

    async function loadScenario() {
        try {
            loading = true;
            scenario = await scenariosAPI.getById(scenarioId);

            if (scenario.season_id) {
                const allGames = await gamesAPI.getBySeason(scenario.season_id);
                currentWeek = getCurrentNFLWeekFromGames(allGames);
            }
        } catch (err: any) {
            error = err.response?.data?.error || 'Failed to load scenario.';
        } finally {
            loading = false;
        }
    }

    async function loadStandings() {
        try {
            standings = await standingsAPI.getByScenario(scenarioId);
        } catch (err: any) {
            console.error('Failed to load standings:', err);
        }
    }

    function handleScenarioUpdated(event: CustomEvent) {
        scenario = event.detail;
        saveStatus = 'saved';
        setTimeout(() => saveStatus = 'idle', 2000);
    }

    function handleWeekChange(event: CustomEvent) {
        currentWeek = event.detail.week;
    }

    function handlePickUpdated() {
        saveStatus = 'saved';
        setTimeout(() => saveStatus = 'idle', 2000);
        loadStandings();
    }

    function handleOpenTeamModal(event: CustomEvent<{ team: PlayoffSeed }>) {
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
        <a href="/scenarios" class="text-primary-400 hover:text-primary-300">
            ← Back to Scenarios
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

    <!-- Main Content - Responsive Grid -->
    <div class="mt-6 space-y-6 lg:space-y-0">
        <!-- Desktop: 3-column layout -->
        <div class="hidden lg:grid lg:grid-cols-[minmax(250px,1fr)_minmax(700px,2fr)_minmax(250px,1fr)] lg:gap-6">
            <!-- Left: AFC Standings -->
            <div class="min-w-0">
                {#if standings && scenario.season_id}
                    <StandingsBox 
                        standings={standings.afc} 
                        conference="AFC"
                        bind:viewMode={standingsViewMode}
                        on:openTeamModal={handleOpenTeamModal}
                    />
                {/if}
            </div>

            <!-- Center: Picks -->
            <div class="min-w-0">
                <PicksBox 
                    {scenarioId}
                    {currentWeek}
                    on:weekChanged={(e) => currentWeek = e.detail.week}
                    on:pickUpdated={handlePickUpdated}
                />
            </div>

            <!-- Right: NFC Standings -->
            <div class="min-w-0">
                {#if standings && scenario.season_id}
                    <StandingsBox 
                        standings={standings.nfc} 
                        conference="NFC"
                        bind:viewMode={standingsViewMode}
                        on:openTeamModal={handleOpenTeamModal}
                    />
                {/if}
            </div>
        </div>

        <!-- Mobile: Stacked layout -->
        <div class="lg:hidden space-y-6">
            <!-- Picks -->
            <PicksBox 
                {scenarioId}
                {currentWeek}
                on:weekChanged={(e) => currentWeek = e.detail.week}
                on:pickUpdated={handlePickUpdated}
            />

            <!-- Standings -->
            {#if standings && scenario.season_id}
                <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <StandingsBoxExpanded 
                        standings={standings.afc} 
                        conference="AFC"
                        bind:viewMode={standingsViewMode}
                        on:openTeamModal={handleOpenTeamModal}
                    />
                    <StandingsBoxExpanded 
                        standings={standings.nfc} 
                        conference="NFC"
                        bind:viewMode={standingsViewMode}
                        on:openTeamModal={handleOpenTeamModal}
                    />
                </div>
            {/if}
        </div>

        <!-- Draft Order -->
        <div class="w-full">
            <DraftOrderBox />
        </div>
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