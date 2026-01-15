<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import type { Game, Pick } from '$types';
    import ConfirmationModal from './ConfirmationModal.svelte';

    export let game: Game;
    export let pick: Pick | undefined = undefined;
    export let compact: boolean = false;
    export let hasPlayoffs: boolean = false;
    export let sportId: number | null = null;

    const dispatch = createEventDispatcher();

    let predictedHomeScore = pick?.predicted_home_score?.toString() || '';
    let predictedAwayScore = pick?.predicted_away_score?.toString() || '';
    let showInfo = false;
    let showConfirmation = false;
    let pendingAction: (() => void) | null = null;
    
    // Reference for positioning fixed tooltips
    let infoButton: HTMLButtonElement;
    let awayButton: HTMLButtonElement;
    let homeButton: HTMLButtonElement;

    $: {
        predictedHomeScore = pick?.predicted_home_score !== null && pick?.predicted_home_score !== undefined 
            ? pick.predicted_home_score.toString() 
            : '';
        predictedAwayScore = pick?.predicted_away_score !== null && pick?.predicted_away_score !== undefined 
            ? pick.predicted_away_score.toString() 
            : '';
    }

    $: isGameCompleted = game.status === 'final';
    
    $: isHomeTeamPicked = pick?.picked_team_id === game.home_team_id;
    $: isAwayTeamPicked = pick?.picked_team_id === game.away_team_id;
    $: isTiePicked = pick !== undefined && pick.picked_team_id === 0;

    $: userMadePick = pick !== undefined;

    $: homeTeamWon = isGameCompleted && game.home_score && game.away_score &&
        game.home_score !== null && 
        game.away_score !== null && 
        game.home_score > game.away_score;
    
    $: awayTeamWon = isGameCompleted && game.away_score && game.home_score &&
        game.home_score !== null && 
        game.away_score !== null && 
        game.away_score > game.home_score;

    $: isActualTie = isGameCompleted && 
        game.home_score && game.home_score !== null && 
        game.away_score && game.away_score !== null && 
        game.home_score === game.away_score;

    $: highlightHomeButton = userMadePick ? isHomeTeamPicked : (isGameCompleted ? homeTeamWon : false);
    $: highlightAwayButton = userMadePick ? isAwayTeamPicked : (isGameCompleted ? awayTeamWon : false);
    $: highlightTieButton = userMadePick ? isTiePicked : (isGameCompleted ? isActualTie : false);

    $: {
        predictedHomeScore = pick?.predicted_home_score !== null && pick?.predicted_home_score !== undefined 
            ? pick.predicted_home_score.toString() 
            : '';
        predictedAwayScore = pick?.predicted_away_score !== null && pick?.predicted_away_score !== undefined 
            ? pick.predicted_away_score.toString() 
            : '';
    }

    $: homeTeamLogoURL = highlightHomeButton && game.home_team.alternate_logo_url 
        ? game.home_team.alternate_logo_url 
        : game.home_team.logo_url;
    
    $: awayTeamLogoURL = highlightAwayButton && game.away_team.alternate_logo_url 
        ? game.away_team.alternate_logo_url 
        : game.away_team.logo_url;

    function selectTeam(teamId: number) {
        if (hasPlayoffs) {
            pendingAction = () => {
                if (pick?.picked_team_id === teamId) {
                    // DELETE the pick
                    predictedHomeScore = '';
                    predictedAwayScore = '';
                    dispatch('pickChanged', {
                        gameId: game.id,
                        deletePick: true
                    });
                } else {
                    // SWITCHING teams - clear scores
                    const wasPickMade = pick !== undefined && pick.picked_team_id !== null;
                    if (wasPickMade) {
                        predictedHomeScore = '';
                        predictedAwayScore = '';
                    }
                    dispatch('pickChanged', {
                        gameId: game.id,
                        pickedTeamId: teamId,
                        predictedHomeScore: undefined,
                        predictedAwayScore: undefined
                    });
                }
            };
            showConfirmation = true;
            return;
        }

        if (pick?.picked_team_id === teamId) {
            // Instead of setting to undefined, DELETE the pick
            predictedHomeScore = '';
            predictedAwayScore = '';
            dispatch('pickChanged', {
                gameId: game.id,
                deletePick: true
            });
        } else {
            // SWITCHING teams - clear scores
            const wasPickMade = pick !== undefined && pick.picked_team_id !== null;
            if (wasPickMade) {
                predictedHomeScore = '';
                predictedAwayScore = '';
            }
            dispatch('pickChanged', {
                gameId: game.id,
                pickedTeamId: teamId,
                predictedHomeScore: undefined,
                predictedAwayScore: undefined
            });
        }
    }

    function selectTie() {
        if (hasPlayoffs) {
            pendingAction = () => {
                if (isTiePicked) {
                    // DELETE the pick
                    predictedHomeScore = '';
                    predictedAwayScore = '';
                    dispatch('pickChanged', {
                        gameId: game.id,
                        deletePick: true
                    });
                } else {
                    // SWITCHING to tie - clear scores
                    const wasPickMade = pick !== undefined && pick.picked_team_id !== null;
                    if (wasPickMade) {
                        predictedHomeScore = '';
                        predictedAwayScore = '';
                    }
                    dispatch('pickChanged', {
                        gameId: game.id,
                        pickedTeamId: 0,
                        predictedHomeScore: undefined,
                        predictedAwayScore: undefined
                    });
                }
            };
            showConfirmation = true;
            return;
        }

        if (isTiePicked) {
            // DELETE the pick
            predictedHomeScore = '';
            predictedAwayScore = '';
            dispatch('pickChanged', {
                gameId: game.id,
                deletePick: true
            });
        } else {
            // SWITCHING to tie - clear scores
            const wasPickMade = pick !== undefined && pick.picked_team_id !== null;
            if (wasPickMade) {
                predictedHomeScore = '';
                predictedAwayScore = '';
            }
            dispatch('pickChanged', {
                gameId: game.id,
                pickedTeamId: 0,
                predictedHomeScore: undefined,
                predictedAwayScore: undefined
            });
        }
    }

    function handleConfirm() {
        showConfirmation = false;
        if (pendingAction) {
            pendingAction();
            pendingAction = null;
        }
    }

    function handleCancel() {
        showConfirmation = false;
        pendingAction = null;
    }

    function handleScoreChange() {
        if (pick) {
            const homeScore = parseScoreInput(predictedHomeScore);
            const awayScore = parseScoreInput(predictedAwayScore);

            let newPickedTeamId = pick.picked_team_id;

            // Determine winner based on scores if both are defined
            if (homeScore !== undefined && awayScore !== undefined) {
                if (homeScore > awayScore) {
                    newPickedTeamId = game.home_team_id;
                } else if (awayScore > homeScore) {
                    newPickedTeamId = game.away_team_id;
                } else {
                    newPickedTeamId = 0; // Tie
                }
            }
            // If scores are not both defined, keep existing pick

            dispatch('pickChanged', {
                gameId: game.id,
                pickedTeamId: newPickedTeamId,
                predictedHomeScore: homeScore,
                predictedAwayScore: awayScore
            });
        }
    }

    function parseScoreInput(value: string): number | undefined {
        if (value === '') return undefined;
        const parsed = parseInt(value);
        return isNaN(parsed) ? undefined : parsed;
    }

    function formatTime(dateString: string): string {
        const date = new Date(dateString);
        const timeStr = date.toLocaleTimeString('en-US', { 
            hour: 'numeric', 
            minute: '2-digit',
        });

        const timeZone = new Intl.DateTimeFormat('en-US', {
            timeZoneName: 'short'
        }).formatToParts(date).find(part => part.type === 'timeZoneName')?.value || '';

        return `${timeStr} ${timeZone}`;
    }

    function formatDate(dateString: string): string {
        const date = new Date(dateString);
        return date.toLocaleDateString('en-US', { 
            weekday: 'short', 
            month: 'short', 
            day: 'numeric' 
        });
    }

    function getTooltipPosition(element: HTMLElement) {
        if (!element) return { top: 0, left: 0 };
        const rect = element.getBoundingClientRect();
        return {
            top: rect.bottom + 8,
            left: rect.left + rect.width / 2
        };
    }

    $: infoPosition = showInfo && infoButton ? getTooltipPosition(infoButton) : { top: 0, left: 0 };
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

<div class="flex items-center gap-2" class:z-50={showInfo} class:opacity-80={isGameCompleted && !userMadePick}>
    <!-- Away Team Section -->
    <div class="flex-1 flex items-stretch">
        <!-- Away Score -->
        <div class="w-10 sm:w-12 shrink-0">
            {#if isGameCompleted && (game.away_score !== null && game.away_score !== undefined) && !userMadePick}
                <div class="h-full flex items-center justify-center font-heading text-xl font-bold rounded-l-lg border-2 border-r-0"
                     style={`background-color: #${game.away_team.primary_color}90; border-color: #${game.away_team.primary_color}; color: #${game.away_team.primary_color};`}>
                    {game.away_score}
                </div>
            {:else}
                <input
                    type="number"
                    min="0"
                    max="99"
                    bind:value={predictedAwayScore}
                    on:change={handleScoreChange}
                    placeholder="--"
                    class="h-full w-full text-center border-2 border-r-0 rounded-l-lg px-1 py-2 font-heading text-lg font-bold placeholder-neutral/40 transition-colors focus:outline-none"
                    style={`background-color: #${game.away_team.primary_color}90; border-color: #${game.away_team.primary_color}; color: #${game.away_team.primary_color};`}
                />
            {/if}
        </div>

        <!-- Away Team Button -->
        <button
            bind:this={awayButton}
            on:click={() => selectTeam(game.away_team_id)}
            on:mouseover={() => {
                if (!highlightAwayButton) {
                    awayButton.style.backgroundColor = `#${game.away_team.primary_color}50`;
                }
            }}
            on:focus={() => {
                if (!highlightAwayButton) {
                    awayButton.style.backgroundColor = `#${game.away_team.primary_color}50`;
                }
            }}
            on:mouseout={() => {
                if (!highlightAwayButton) {
                    awayButton.style.backgroundColor = 'transparent';
                }
            }}
            on:blur={() => {
                if (!highlightAwayButton) {
                    awayButton.style.backgroundColor = 'transparent';
                }
            }}
            class="flex-1 p-2 rounded-r-lg border-2 transition-all cursor-pointer"
            class:border-primary-600={!highlightAwayButton}
            class:hover:border-primary-400={!highlightAwayButton}
            style={highlightAwayButton 
                ? `background-color: #${game.away_team.primary_color}; border-color: #${game.away_team.primary_color};` 
                : `background-color: transparent; border-color: #${game.away_team.primary_color};`}
        >
            <div class="flex items-center justify-center gap-2">
                {#if awayTeamLogoURL}
                    <img 
                        src={awayTeamLogoURL}
                        alt={game.away_team.abbreviation}
                        class="w-8 h-8 object-contain"
                    />
                {/if}
                {#if compact}
                    <span class="hidden sm:inline sm:text-sm font-sans font-semibold transition-colors"
                          class:text-white={highlightAwayButton}
                          class:text-black={!highlightAwayButton}>
                        {game.away_team.abbreviation}
                    </span>
                {/if}
            </div>
        </button>
    </div>

    <!-- Center: Info/Tie -->
    <div class="flex flex-col items-center gap-1 shrink-0">
        <!-- Info Button -->
        <button
            bind:this={infoButton}
            on:mouseenter={() => showInfo = true}
            on:mouseleave={() => showInfo = false}
            class="p-1 sm:p-1.5 rounded-full bg-primary-700/50 hover:bg-primary-600 transition-colors cursor-pointer"
            title="Game Info"
        >
            <svg class="w-3 sm:w-4 h-3 sm:h-4 text-neutral" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
        </button>

        <!-- Tie Button -->
        {#if sportId === 1}
            <button
                on:click={selectTie}
                class="px-1.5 sm:px-2 py-0.5 rounded border text-xs font-sans font-semibold transition-all cursor-pointer"
                class:bg-primary-600={highlightTieButton}
                class:border-primary-500={highlightTieButton}
                class:text-neutral={highlightTieButton}
                class:bg-transparent={!highlightTieButton}
                class:border-primary-600={!highlightTieButton}
                class:hover:border-primary-400={!highlightTieButton}
            >
                TIE
            </button>
        {/if}
    </div>

    <!-- Home Team Section -->
    <div class="flex-1 flex items-stretch">
        <!-- Home Team Button -->
        <button
            bind:this={homeButton}
            on:click={() => selectTeam(game.home_team_id)}
            on:mouseover={() => {
                if (!highlightHomeButton) {
                    homeButton.style.backgroundColor = `#${game.home_team.primary_color}50`;
                }
            }}
            on:focus={() => {
                if (!highlightHomeButton) {
                    homeButton.style.backgroundColor = `#${game.home_team.primary_color}50`;
                }
            }}
            on:mouseout={() => {
                if (!highlightHomeButton) {
                    homeButton.style.backgroundColor = 'transparent';
                }
            }}
            on:blur={() => {
                if (!highlightHomeButton) {
                    homeButton.style.backgroundColor = 'transparent';
                }
            }}
            on:mouseout={() => {
                if (!highlightHomeButton) {
                    homeButton.style.backgroundColor = 'transparent';
                }
            }}
            on:blur={() => {
                if (!highlightHomeButton) {
                    homeButton.style.backgroundColor = 'transparent';
                }
            }}
            class="flex-1 p-2 rounded-l-lg border-2 transition-all cursor-pointer"
            class:border-primary-600={!highlightHomeButton}
            class:hover:border-primary-400={!highlightHomeButton}
            style={highlightHomeButton 
                ? `background-color: #${game.home_team.primary_color}; border-color: #${game.home_team.primary_color};` 
                : `background-color: transparent; border-color: #${game.home_team.primary_color};`}
        >
            <div class="flex items-center justify-center gap-2">
                {#if homeTeamLogoURL}
                    <img 
                        src={homeTeamLogoURL}
                        alt={game.home_team.abbreviation}
                        class="w-8 h-8 object-contain"
                    />
                {/if}
                {#if compact}
                    <span class="hidden sm:inline sm:text-sm font-sans font-semibold transition-colors"
                          class:text-white={highlightHomeButton}
                          class:text-black={!highlightHomeButton}>
                        {game.home_team.abbreviation}
                    </span>
                {/if}
            </div>
        </button>

        <!-- Home Score -->
        <div class="w-12 shrink-0">
            {#if isGameCompleted && (game.home_score !== null && game.home_score !== undefined) && !userMadePick}
                <div class="h-full flex items-center justify-center font-heading text-xl font-bold rounded-r-lg border-2 border-l-0"
                     style={`background-color: #${game.home_team.primary_color}90; border-color: #${game.home_team.primary_color}; color: #${game.home_team.primary_color};`}>
                    {game.home_score}
                </div>
            {:else}
                <input
                    type="number"
                    min="0"
                    max="99"
                    bind:value={predictedHomeScore}
                    on:change={handleScoreChange}
                    placeholder="--"
                    class="h-full w-full text-center border-2 border-l-0 rounded-r-lg px-1 py-2 font-heading text-lg font-bold placeholder-neutral/40 transition-colors focus:outline-none"
                    style={`background-color: #${game.home_team.primary_color}90; border-color: #${game.home_team.primary_color}; color: #${game.home_team.primary_color};`}
                />
            {/if}
        </div>
    </div>
</div>

<!-- Tooltip rendered outside parent (not affected by opacity) -->
{#if showInfo}
    <div 
        class="fixed z-50 w-48 bg-primary-950 border border-primary-600 rounded-lg p-3 text-sm text-neutral text-center shadow-xl pointer-events-none"
        style="top: {infoPosition.top}px; left: {infoPosition.left}px; transform: translateX(-50%);"
    >
        <div class="space-y-1 font-sans">
            <div class="font-semibold">{formatDate(game.start_time)}</div>
            <div>{formatTime(game.start_time)}</div>
            {#if game.location}
                <div class="text-neutral/70">{game.location}</div>
            {/if}
            {#if game.network}
                <div class="text-primary-400">{game.network}</div>
            {/if}
            {#if isGameCompleted}
                <div class="pt-2 mt-2 border-t border-primary-600">
                    <div class="text-green-400 font-semibold mb-1">FINAL</div>
                    <div class="text-neutral">{game.away_team.abbreviation} {game.away_score} - {game.home_team.abbreviation} {game.home_score}</div>
                </div>
            {/if}
        </div>
    </div>
{/if}

{#if showConfirmation}
    <ConfirmationModal
        title="Reset Playoff Matchups?"
        message="Changing this regular season pick will reset all playoff picks and regenerate playoff matchups. This action cannot be undone."
        confirmText="Change Pick"
        cancelText="Cancel"
        warningType="regular"
        on:confirm={handleConfirm}
        on:cancel={handleCancel}
    />
{/if}