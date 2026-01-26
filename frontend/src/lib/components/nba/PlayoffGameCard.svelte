<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import type { PlayoffMatchup, PlayoffSeries } from '$types';
    import { NBA_PLAYOFF_ROUND_NAMES } from '$types';
    import ConfirmationModal from '../scenarios/ConfirmationModal.svelte';

    // Props
    export let item: PlayoffMatchup | PlayoffSeries;
    export let hasLaterRounds: boolean = false;

    // Event dispatcher
    const dispatch = createEventDispatcher();

    // State variables for single game matchups
    let predictedHigherScore = '';
    let predictedLowerScore = '';

    // State variables for series matchups
    let predictedHigherWins = '';
    let predictedLowerWins = '';

    // State variable for pending team selection
    let pendingTeamId: number | null = null;

    // State variable for confirmation modal
    let showConfirmation = false;

    // State variables for team tooltips
    let showLowerSeedName = false;
    let showHigherSeedName = false;
    let lowerSeedButton: HTMLButtonElement;
    let higherSeedButton: HTMLButtonElement;
    
    // Determine if this is a series or single game
    $: isSeries = 'best_of' in item;

    // Initialize predicted scores/wins based on item type
    $: {
        if (isSeries) {
            const series = item as PlayoffSeries;
            predictedHigherWins = series.predicted_higher_seed_wins != null ? series.predicted_higher_seed_wins?.toString() : '';
            predictedLowerWins = series.predicted_lower_seed_wins != null ? series.predicted_lower_seed_wins?.toString() : '';
        } else {
            const matchup = item as PlayoffMatchup;
            predictedHigherScore = matchup.predicted_higher_seed_score?.toString() || '';
            predictedLowerScore = matchup.predicted_lower_seed_score?.toString() || '';
        }
    }

    // Determine if higher or lower seed is picked
    $: isHigherSeedPicked = item.picked_team_id === item.higher_seed_team_id;
    $: isLowerSeedPicked = item.picked_team_id === item.lower_seed_team_id;

    // Determine team logo URLs, using alternate if picked
    $: higherTeamLogoURL = (isHigherSeedPicked && item.higher_seed_team?.alternate_logo_url)
        ? item.higher_seed_team.alternate_logo_url
        : (item.higher_seed_team?.logo_url || '');
    $: lowerTeamLogoURL = (isLowerSeedPicked && item.lower_seed_team?.alternate_logo_url)
        ? item.lower_seed_team.alternate_logo_url
        : (item.lower_seed_team?.logo_url || '');

    // Calculate tooltip positions
    $: lowerSeedPosition = showLowerSeedName && lowerSeedButton ? getTeamTooltipPosition(lowerSeedButton) : { top: 0, left: 0 };
    $: higherSeedPosition = showHigherSeedName && higherSeedButton ? getTeamTooltipPosition(higherSeedButton) : { top: 0, left: 0 };

    // Show confirmation modal if later rounds exist
    function selectTeam(teamId: number) {
        if (hasLaterRounds && item.round < 6) {
            pendingTeamId = teamId;
            showConfirmation = true;
            return;
        }

        executeTeamSelection(teamId);
    }

    // Handle team selection logic
    function executeTeamSelection(teamId: number) {
        if (item.picked_team_id === teamId) {
            // Deselect picked team
            if (isSeries) {
                predictedHigherWins = '';
                predictedLowerWins = '';
                dispatch('pickChanged', {
                    itemId: item.id,
                    isSeries: true,
                    pickedTeamId: null,
                    predictedHigherWins: undefined,
                    predictedLowerWins: undefined
                })
            } else {
                predictedHigherScore = '';
                predictedLowerScore = '';
                dispatch('pickChanged', {
                    matchupId: item.id,
                    isSeries: false,
                    pickedTeamId: null,
                    predictedHigherScore: undefined,
                    predictedLowerScore: undefined
                });
            }
        } else {
            // Select new team
            const wasPickMade = item.picked_team_id !== undefined && item.picked_team_id !== null;
            if (wasPickMade) {
                if (isSeries) {
                    predictedHigherWins = '';
                    predictedLowerWins = '';
                } else {
                    predictedHigherScore = '';
                    predictedLowerScore = '';
                }
            }

            if (isSeries) {
                dispatch('pickChanged', {
                    itemId: item.id,
                    isSeries: true,
                    pickedTeamId: teamId,
                    predictedHigherWins: undefined,
                    predictedLowerWins: undefined
                });
            } else {
                dispatch('pickChanged', {
                    matchupId: item.id,
                    isSeries: false,
                    pickedTeamId: teamId,
                    predictedHigherScore: undefined,
                    predictedLowerScore: undefined
                });
            }
        }
    }

    // Handle confirmation modal actions
    function handleConfirm() {
        showConfirmation = false;
        if (pendingTeamId !== null) {
            executeTeamSelection(pendingTeamId);
            pendingTeamId = null;
        }
    }

    // Cancel team selection change
    function handleCancel() {
        showConfirmation = false;
        pendingTeamId = null;
    }

    // Handle score input changes for single game matchups
    function handleScoreChange() {
        if (!isSeries && item.picked_team_id) {
            const higherScore = parseInput(predictedHigherScore);
            const lowerScore = parseInput(predictedLowerScore);

            // Determine winner based on scores
            // If scores are invalid or tied, do not change pick
            let newPickedTeamId = item.picked_team_id;
            if (higherScore !== undefined && lowerScore !== undefined) {
                if (higherScore > lowerScore) {
                    newPickedTeamId = item.higher_seed_team_id;
                } else if (lowerScore > higherScore) {
                    newPickedTeamId = item.lower_seed_team_id;
                }
            }

            dispatch('pickChanged', {
                matchupId: item.id,
                isSeries: false,
                pickedTeamId: newPickedTeamId,
                predictedHigherScore: higherScore,
                predictedLowerScore: lowerScore
            });
        }
    }

    // Handle wins input changes for series matchups
    function handleWinsChange() {
        if (isSeries && item.picked_team_id) {
            const higherWins = parseInput(predictedHigherWins);
            const lowerWins = parseInput(predictedLowerWins);

            // Validate wins input
            if (higherWins !== undefined && lowerWins !== undefined) {
                if (higherWins < 0 || lowerWins < 0 || higherWins > 4 || lowerWins > 4) {
                    return;
                }
                if (higherWins === 4 && lowerWins === 4) {
                    return;
                }
                if (higherWins < 4 && lowerWins < 4) {
                    return;
                }
            }

            // Determine winner based on wins
            let newPickedTeamId = item.picked_team_id;
            if (higherWins !== undefined && lowerWins !== undefined) {
                if (higherWins === 4) {
                    newPickedTeamId = item.higher_seed_team_id;
                } else if (lowerWins === 4) {
                    newPickedTeamId = item.lower_seed_team_id;
                }
            }

            dispatch('pickChanged', {
                itemId: item.id,
                isSeries: true,
                pickedTeamId: newPickedTeamId,
                predictedHigherWins: higherWins,
                predictedLowerWins: lowerWins
            });
        }
    }

    // Parse score or wins input string to number or undefined
    function parseInput(value: string): number | undefined {
        if (value === '') return undefined;
        const parsed = parseInt(value);
        return isNaN(parsed) ? undefined : parsed;
    }

    // Calculate tooltip position for a team button
    function getTeamTooltipPosition(element: HTMLElement) {
        if (!element) return { top: 0, left: 0 };
        const rect = element.getBoundingClientRect();
        return {
            top: rect.top - 8,
            left: rect.left + rect.width / 2
        };
    }
</script>

<style>
    /* Hide spin buttons for number inputs */
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
        <!-- Lower Seed Section -->
        <div class="flex-1 flex items-stretch">
            <!-- Lower Seed Score -->
            <div class="w-10 shrink-0">
                {#if isSeries}
                    <input
                        type="number"
                        min="0"
                        max="4"
                        bind:value={predictedLowerWins}
                        on:change={handleWinsChange}
                        placeholder="--"
                        class="h-full w-full text-center border-2 border-r-0 rounded-l-lg px-1 py-2 font-heading text-lg font-bold placeholder-neutral/40 transition-colors focus:outline-none"
                        style={`background-color: #${item.lower_seed_team?.primary_color}90; border-color: #${item.lower_seed_team?.primary_color}; color: #${item.lower_seed_team?.primary_color};`}
                    />
                {:else}
                    <input
                        type="number"
                        min="0"
                        max="200"
                        bind:value={predictedLowerScore}
                        on:change={handleScoreChange}
                        placeholder="--"
                        class="h-full w-full text-center border-2 border-r-0 rounded-l-lg px-1 py-2 font-heading text-lg font-bold placeholder-neutral/40 transition-colors focus:outline-none"
                        style={`background-color: #${item.lower_seed_team?.primary_color}90; border-color: #${item.lower_seed_team?.primary_color}; color: #${item.lower_seed_team?.primary_color};`}
                    />
                {/if}
            </div>

            <!-- Lower Seed Button -->
            <button
                bind:this={lowerSeedButton}
                on:click={() => selectTeam(item.lower_seed_team_id)}
                on:mouseenter={() => showLowerSeedName = true}
                on:mouseleave={() => showLowerSeedName = false}
                on:mouseover={() => {
                    if (!isLowerSeedPicked) {
                        lowerSeedButton.style.backgroundColor = `#${item.lower_seed_team?.primary_color}50`;
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
                    ? `background-color: #${item.lower_seed_team?.primary_color}; border-color: #${item.lower_seed_team?.primary_color};` 
                    : `background-color: transparent; border-color: #${item.lower_seed_team?.primary_color};`}
            >
                <!-- Seed Badge -->
                <div
                    class="absolute top-0.5 right-0 text-xs font-bold px-1.5 py-0 rounded"
                    style={isLowerSeedPicked 
                        ? `color: white;` 
                        : `color: #${item.lower_seed_team?.primary_color};`}
                >
                    {item.lower_seed}
                </div>

                <div class="flex items-center justify-center">
                    {#if lowerTeamLogoURL}
                        <img 
                            src={lowerTeamLogoURL}
                            alt={item.lower_seed_team?.abbreviation}
                            class="w-8 h-8 object-contain"
                        />
                    {/if}
                </div>
            </button>
        </div>

        <!-- VS Divider -->
        <div class="text-xs font-sans font-bold text-black/50 shrink-0">VS</div>

        <!-- Higher Seed Section -->
        <div class="flex-1 flex items-stretch">
            <!-- Higher Seed Button -->
            <button
                bind:this={higherSeedButton}
                on:click={() => selectTeam(item.higher_seed_team_id)}
                on:mouseenter={() => showHigherSeedName = true}
                on:mouseleave={() => showHigherSeedName = false}
                on:mouseover={() => {
                    if (!isHigherSeedPicked) {
                        higherSeedButton.style.backgroundColor = `#${item.higher_seed_team?.primary_color}50`;
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
                    ? `background-color: #${item.higher_seed_team?.primary_color}; border-color: #${item.higher_seed_team?.primary_color};` 
                    : `background-color: transparent; border-color: #${item.higher_seed_team?.primary_color};`}
            >
                <!-- Seed Badge -->
                <div
                    class="absolute top-0.5 left-0 text-xs font-bold px-1.5 py-0 rounded"
                    style={isHigherSeedPicked 
                        ? `color: white;` 
                        : `color: #${item.higher_seed_team?.primary_color};`}
                >
                    {item.higher_seed}
                </div>
                
                <div class="flex items-center justify-center">
                    {#if higherTeamLogoURL}
                        <img 
                            src={higherTeamLogoURL}
                            alt={item.higher_seed_team?.abbreviation}
                            class="w-8 h-8 object-contain"
                        />
                    {/if}
                </div>
            </button>

            <!-- Higher Seed Score -->
            <div class="w-10 shrink-0">
                {#if isSeries}
                    <input
                        type="number"
                        min="0"
                        max="4"
                        bind:value={predictedHigherWins}
                        on:change={handleWinsChange}
                        placeholder="--"
                        class="h-full w-full text-center border-2 border-l-0 rounded-r-lg px-1 py-2 font-heading text-lg font-bold placeholder-neutral/40 transition-colors focus:outline-none"
                        style={`background-color: #${item.higher_seed_team?.primary_color}90; border-color: #${item.higher_seed_team?.primary_color}; color: #${item.higher_seed_team?.primary_color};`}
                    />
                {:else}
                    <input
                        type="number"
                        min="0"
                        max="200"
                        bind:value={predictedHigherScore}
                        on:change={handleScoreChange}
                        placeholder="--"
                        class="h-full w-full text-center border-2 border-l-0 rounded-r-lg px-1 py-2 font-heading text-lg font-bold placeholder-neutral/40 transition-colors focus:outline-none"
                        style={`background-color: #${item.higher_seed_team?.primary_color}90; border-color: #${item.higher_seed_team?.primary_color}; color: #${item.higher_seed_team?.primary_color};`}
                    />
                {/if}
            </div>
        </div>
    </div>
</div>

<!-- Lower Seed Tooltip -->
{#if showLowerSeedName && item.lower_seed_team}
    <div
        class="fixed z-50 px-3 py-1.5 bg-primary-950 border border-primary-600 rounded-lg shadow-xl whitespace-nowrap pointer-events-none"
        style="top: {lowerSeedPosition.top}px; left: {lowerSeedPosition.left}px; transform: translate(-50%, -100%);"
    >
        <span class="text-sm font-sans font-semibold text-neutral">
            {item.lower_seed_team.city} {item.lower_seed_team.name}
        </span>
    </div>
{/if}

<!-- Higher Seed Tooltip -->
{#if showHigherSeedName && item.higher_seed_team}
    <div
        class="fixed z-50 px-3 py-1.5 bg-primary-950 border border-primary-600 rounded-lg shadow-xl whitespace-nowrap pointer-events-none"
        style="top: {higherSeedPosition.top}px; left: {higherSeedPosition.left}px; transform: translate(-50%, -100%);"
    >
        <span class="text-sm font-sans font-semibold text-neutral">
            {item.higher_seed_team.city} {item.higher_seed_team.name}
        </span>
    </div>
{/if}

<!-- Confirmation Modal -->
{#if showConfirmation}
    <ConfirmationModal
        title="Reset Later Playoff Rounds?"
        message={`Changing this ${NBA_PLAYOFF_ROUND_NAMES[item.round]} pick will reset all subsequent playoff rounds and regenerate future matchups. This action cannot be undone.`}
        warningType="playoff"
        confirmText="Change Pick"
        cancelText="Cancel"
        on:confirm={handleConfirm}
        on:cancel={handleCancel}
    />
{/if}