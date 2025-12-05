<script lang="ts">
    import { onMount } from 'svelte';
    import { page } from '$app/stores';
    import { scenariosAPI } from '$lib/api/scenarios';
    import type { Scenario } from '$types';

    import ScenarioHeader from '$lib/components/scenarios/ScenarioHeader.svelte';
    import ScenarioSettings from '$lib/components/scenarios/ScenarioSettings.svelte';
    import ScenarioInfo from '$lib/components/nfl/ScenarioInfo.svelte';
    import PicksBox from '$lib/components/nfl/PicksBox.svelte';
    import StandingsBox from '$lib/components/nfl/StandingsBox.svelte';
    import DraftOrderBox from '$lib/components/nfl/DraftOrderBox.svelte';

    let scenarioId: number;
    let scenario: Scenario | null = null;
    let loading = true;
    let error = '';

    let showSettings = false;
    let showInfo = false;

    // Current week being viewed
    let currentWeek = 1;

    // Save status indicator
    let saveStatus: 'idle' | 'saving' | 'saved' | 'error' = 'idle';

    $: scenarioId = parseInt($page.params.id);

    onMount(async () => {
        await loadScenario();
    });

    async function loadScenario() {
        try {
            loading = true;
            scenario = await scenariosAPI.getById(scenarioId);

            // Determine current week (could be based on current date or last edited week)
            currentWeek = getCurrentWeek(scenario);
        } catch (err: any) {
            error = err.response?.data?.error || 'Failed to load scenario.';
        } finally {
            loading = false;
        }
    }

    function getCurrentWeek(scenario: Scenario): number {
        // TODO: Calculate current NFL week based on today's date
        // For now, default to week 1
        return 1;
    }

    function handleScenarioUpdated(event: CustomEvent) {
        scenario = event.detail;
        saveStatus = 'saved';
        setTimeout(() => saveStatus = 'idle', 2000);
    }

    function handlePickUpdated() {
        saveStatus = 'saved';
        setTimeout(() => saveStatus = 'idle', 2000);
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
        <div class="hidden lg:grid lg:grid-cols-12 lg:gap-6">
            <!-- Left: AFC Standings -->
            <div class="lg:col-span-3">
                <StandingsBox conference="AFC" {scenarioId} {currentWeek} />
            </div>

            <!-- Center: Picks -->
            <div class="lg:col-span-6">
                <PicksBox 
                    {scenarioId}
                    {currentWeek}
                    on:weekChanged={(e) => currentWeek = e.detail.week}
                    on:pickUpdated={handlePickUpdated}
                />
            </div>

            <!-- Right: NFC Standings -->
            <div class="lg:col-span-3">
                <StandingsBox conference="NFC" {scenarioId} {currentWeek} />
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
            <StandingsBox conference="AFC" {scenarioId} {currentWeek} expanded={true} />
            <StandingsBox conference="NFC" {scenarioId} {currentWeek} expanded={true} />
        </div>

        <!-- Draft Order -->
        <div class="w-full">
            <DraftOrderBox {scenarioId} />
        </div>
    </div>
{/if}