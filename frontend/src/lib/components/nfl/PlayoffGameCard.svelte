<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import type { PlayoffMatchup } from '$types';
    import { PLAYOFF_ROUND_NAMES } from '$types';
    import ConfirmationModal from './ConfirmationModal.svelte';

    export let matchup: PlayoffMatchup;
    export let hasLaterRounds: boolean = false;

    const dispatch = createEventDispatcher();

    let predictedHigherScore = matchup.predicted_higher_seed_score?.toString() || '';
    let predictedLowerScore = matchup.predicted_lower_seed_score?.toString() || '';
    let showConfirmation = false;
    let pendingTeamId: number | null = null;

    // Update scores when matchup changes
    $: {
        predictedHigherScore = matchup.predicted_higher_seed_score?.toString() || '';
        predictedLowerScore = matchup.predicted_lower_seed_score?.toString() || '';
    }

    $: isHigherSeedPicked = matchup.picked_team_id === matchup.higher_seed_team_id;
    $: isLowerSeedPicked = matchup.picked_team_id === matchup.lower_seed_team_id;

    $: higherTeamLogoURL = isHigherSeedPicked && matchup.higher_seed_team?.alternate_logo_url 
        ? matchup.higher_seed_team.alternate_logo_url 
        : matchup.higher_seed_team?.logo_url;
    
    $: lowerTeamLogoURL = isLowerSeedPicked && matchup.lower_seed_team?.alternate_logo_url 
        ? matchup.lower_seed_team.alternate_logo_url 
        : matchup.lower_seed_team?.logo_url;

    let showLowerSeedName = false;
    let showHigherSeedName = false;
    let lowerSeedButton: HTMLButtonElement;
    let higherSeedButton: HTMLButtonElement;

    function selectTeam(teamId: number) {
        // Show warning if later rounds exist
        if (hasLaterRounds && matchup.round < 4) {
            pendingTeamId = teamId;
            showConfirmation = true;
            return;
        }

        executeTeamSelection(teamId);
    }

    function executeTeamSelection(teamId: number) {
        if (matchup.picked_team_id === teamId) {
            // DELETE the pick
            dispatch('pickChanged', {
                matchupId: matchup.id,
                pickedTeamId: null,
                predictedHigherScore: parseScoreInput(predictedHigherScore),
                predictedLowerScore: parseScoreInput(predictedLowerScore)
            });
        } else {
            dispatch('pickChanged', {
                matchupId: matchup.id,
                pickedTeamId: teamId,
                predictedHigherScore: parseScoreInput(predictedHigherScore),
                predictedLowerScore: parseScoreInput(predictedLowerScore)
            });
        }
    }

    function handleConfirm() {
        showConfirmation = false;
        if (pendingTeamId !== null) {
            executeTeamSelection(pendingTeamId);
            pendingTeamId = null;
        }
    }

    function handleCancel() {
        showConfirmation = false;
        pendingTeamId = null;
    }

    function handleScoreChange() {
        if (matchup.picked_team_id) {
            dispatch('pickChanged', {
                matchupId: matchup.id,
                pickedTeamId: matchup.picked_team_id,
                predictedHigherScore: parseScoreInput(predictedHigherScore),
                predictedLowerScore: parseScoreInput(predictedLowerScore)
            });
        }
    }

    function parseScoreInput(value: string): number | undefined {
        if (value === '') return undefined;
        const parsed = parseInt(value);
        return isNaN(parsed) ? undefined : parsed;
    }

    function getTeamTooltipPosition(element: HTMLElement) {
        if (!element) return { top: 0, left: 0 };
        const rect = element.getBoundingClientRect();
        return {
            top: rect.top - 8,
            left: rect.left + rect.width / 2
        };
    }

    $: lowerSeedPosition = showLowerSeedName && lowerSeedButton ? getTeamTooltipPosition(lowerSeedButton) : { top: 0, left: 0 };
    $: higherSeedPosition = showHigherSeedName && higherSeedButton ? getTeamTooltipPosition(higherSeedButton) : { top: 0, left: 0 };
</script>

<style>
    input[type='number']::-webkit-inner-spin-button,
    input[type='number']::-webkit-outer-spin-button {
        -webkit-appearance: none;
        margin: 0;
    }
    input[type='number'] {
        -moz-appearance: textfield;
        appearance: textfield;
    }
</style>

<div
    class="relative bg-neutral rounded-lg p-3"
    class:z-50={showLowerSeedName || showHigherSeedName}
>
    <div class="flex items-center gap-2">
        <!-- Lower Seed Section (Score on LEFT, Button on RIGHT like away team) -->
        <div class="flex-1 flex items-stretch">
            <!-- Lower Seed Score -->
            <div class="w-10 shrink-0">
                <input
                    type="number"
                    min="0"
                    max="99"
                    bind:value={predictedLowerScore}
                    on:change={handleScoreChange}
                    placeholder="--"
                    class="h-full w-full text-center border-2 border-r-0 rounded-l-lg px-1 py-2 font-heading text-lg font-bold placeholder-neutral/40 transition-colors focus:outline-none"
                    style={`background-color: #${matchup.lower_seed_team?.primary_color}90; border-color: #${matchup.lower_seed_team?.primary_color}; color: #${matchup.lower_seed_team?.primary_color};`}
                />
            </div>

            <!-- Lower Seed Button -->
            <button
                bind:this={lowerSeedButton}
                on:click={() => selectTeam(matchup.lower_seed_team_id)}
                on:mouseenter={() => showLowerSeedName = true}
                on:mouseleave={() => showLowerSeedName = false}
                on:mouseover={() => {
                    if (!isLowerSeedPicked) {
                        lowerSeedButton.style.backgroundColor = `#${matchup.lower_seed_team?.primary_color}50`;
                    }
                }}
                on:mouseout={() => {
                    if (!isLowerSeedPicked) {
                        lowerSeedButton.style.backgroundColor = 'transparent';
                    }
                }}
                on:focus={() => showLowerSeedName = true}
                on:blur={() => showLowerSeedName = false}
                class="flex-1 py-2 px-3 rounded-r-lg border-2 transition-all cursor-pointer relative"
                class:border-primary-600={!isLowerSeedPicked}
                class:hover:border-primary-400={!isLowerSeedPicked}
                style={isLowerSeedPicked 
                    ? `background-color: #${matchup.lower_seed_team?.primary_color}; border-color: #${matchup.lower_seed_team?.primary_color};` 
                    : `background-color: transparent; border-color: #${matchup.lower_seed_team?.primary_color};`}
            >
                <!-- Seed Badge -->
                <div
                    class="absolute top-0.5 right-0 text-xs font-bold px-1.5 py-0 rounded"
                    style={isLowerSeedPicked 
                        ? `color: white;` 
                        : `color: #${matchup.lower_seed_team?.primary_color};`}
                >
                    {matchup.lower_seed}
                </div>

                <div class="flex items-center justify-center">
                    {#if lowerTeamLogoURL}
                        <img 
                            src={lowerTeamLogoURL}
                            alt={matchup.lower_seed_team?.abbreviation}
                            class="w-8 h-8 object-contain"
                        />
                    {/if}
                </div>
            </button>
        </div>

        <!-- VS Divider -->
        <div class="text-xs font-sans font-bold text-black/50 shrink-0">VS</div>

        <!-- Higher Seed Section (Button on LEFT, Score on RIGHT like home team) -->
        <div class="flex-1 flex items-stretch">
            <!-- Higher Seed Button -->
            <button
                bind:this={higherSeedButton}
                on:click={() => selectTeam(matchup.higher_seed_team_id)}
                on:mouseenter={() => showHigherSeedName = true}
                on:mouseleave={() => showHigherSeedName = false}
                on:mouseover={() => {
                    if (!isHigherSeedPicked) {
                        higherSeedButton.style.backgroundColor = `#${matchup.higher_seed_team?.primary_color}50`;
                    }
                }}
                on:mouseout={() => {
                    if (!isHigherSeedPicked) {
                        higherSeedButton.style.backgroundColor = 'transparent';
                    }
                }}
                on:focus={() => showHigherSeedName = true}
                on:blur={() => showHigherSeedName = false}
                class="flex-1 py-2 px-3 rounded-l-lg border-2 transition-all cursor-pointer relative"
                class:border-primary-600={!isHigherSeedPicked}
                class:hover:border-primary-400={!isHigherSeedPicked}
                style={isHigherSeedPicked 
                    ? `background-color: #${matchup.higher_seed_team?.primary_color}; border-color: #${matchup.higher_seed_team?.primary_color};` 
                    : `background-color: transparent; border-color: #${matchup.higher_seed_team?.primary_color};`}
            >
                <!-- Seed Badge -->
                <div
                    class="absolute top-0.5 left-0 text-xs font-bold px-1.5 py-0 rounded"
                    style={isHigherSeedPicked 
                        ? `color: white;` 
                        : `color: #${matchup.higher_seed_team?.primary_color};`}
                >
                    {matchup.higher_seed}
                </div>
                
                <div class="flex items-center justify-center">
                    {#if higherTeamLogoURL}
                        <img 
                            src={higherTeamLogoURL}
                            alt={matchup.higher_seed_team?.abbreviation}
                            class="w-8 h-8 object-contain"
                        />
                    {/if}
                </div>
            </button>

            <!-- Higher Seed Score -->
            <div class="w-10 shrink-0">
                <input
                    type="number"
                    min="0"
                    max="99"
                    bind:value={predictedHigherScore}
                    on:change={handleScoreChange}
                    placeholder="--"
                    class="h-full w-full text-center border-2 border-l-0 rounded-r-lg px-1 py-2 font-heading text-lg font-bold placeholder-neutral/40 transition-colors focus:outline-none"
                    style={`background-color: #${matchup.higher_seed_team?.primary_color}90; border-color: #${matchup.higher_seed_team?.primary_color}; color: #${matchup.higher_seed_team?.primary_color};`}
                />
            </div>
        </div>
    </div>
</div>

<!-- Tooltips rendered outside parent -->
{#if showLowerSeedName && matchup.lower_seed_team}
    <div
        class="fixed z-50 px-3 py-1.5 bg-primary-950 border border-primary-600 rounded-lg shadow-xl whitespace-nowrap pointer-events-none"
        style="top: {lowerSeedPosition.top}px; left: {lowerSeedPosition.left}px; transform: translate(-50%, -100%);"
    >
        <span class="text-sm font-sans font-semibold text-neutral">
            {matchup.lower_seed_team.city} {matchup.lower_seed_team.name}
        </span>
    </div>
{/if}

{#if showHigherSeedName && matchup.higher_seed_team}
    <div
        class="fixed z-50 px-3 py-1.5 bg-primary-950 border border-primary-600 rounded-lg shadow-xl whitespace-nowrap pointer-events-none"
        style="top: {higherSeedPosition.top}px; left: {higherSeedPosition.left}px; transform: translate(-50%, -100%);"
    >
        <span class="text-sm font-sans font-semibold text-neutral">
            {matchup.higher_seed_team.city} {matchup.higher_seed_team.name}
        </span>
    </div>
{/if}

{#if showConfirmation}
    <ConfirmationModal
        title="Reset Later Playoff Rounds?"
        message={`Changing this ${PLAYOFF_ROUND_NAMES[matchup.round]} pick will reset all subsequent playoff rounds and regenerate future matchups. This action cannot be undone.`}
        warningType="playoff"
        confirmText="Change Pick"
        cancelText="Cancel"
        on:confirm={handleConfirm}
        on:cancel={handleCancel}
    />
{/if}