<script lang="ts">
    import { createEventDispatcher, onMount, onDestroy } from 'svelte';
    import { getNFLWeekDateRangesFromGames, formatNFLDateRange, type NFLWeekDateRange } from '$lib/utils/nfl/dates';
    import { getNBAWeekDateRangesFromGames, formatNBADateRange, type NBAWeekDateRange } from '$lib/utils/nba/dates';
    import { NFL_PLAYOFF_ROUND_NAMES, NBA_PLAYOFF_ROUND_NAMES } from '$types';
    import type { Game, PlayoffState } from '$types';

    // Props
    export let currentWeek: number;
    export let allGames: Game[] = [];
    export let playoffState: PlayoffState | null = null;
    export let canEnablePlayoffs: boolean = false;
    export let sportId: number | null = null;
    export let isCurrentRoundComplete: boolean = false;

    // Event dispatcher
    const dispatch = createEventDispatcher();
    
    // State variables for dropdown and week calculations
    let weekDropdownOpen = false;
    let weekDateRanges: Map<number, NFLWeekDateRange | NBAWeekDateRange> = new Map();

    // Constants based on sport
    $: NUM_WEEKS = sportId === 1 ? 18 : sportId === 2 ? 25 : 0;
    $: FINAL_PLAYOFF_ROUND = sportId === 1 ? 4 : sportId === 2 ? 6 : 0;

    // Calculate max available week based on playoff state
    $: maxAvailableWeek = getMaxAvailableWeek(sportId, playoffState, isCurrentRoundComplete, currentWeek);

    // Calculate playoff round if applicable
    $: playoffRound = currentWeek > NUM_WEEKS ? currentWeek - NUM_WEEKS : 0;

    // Get week date ranges based on sport
    $: if (sportId && allGames.length > 0) {
        weekDateRanges = sportId === 1
            ? getNFLWeekDateRangesFromGames(allGames)
            : getNBAWeekDateRangesFromGames(allGames);
    }

    // Current label for the button
    $: currentLabel = playoffRound === 0
        ? `Week ${currentWeek}` : (sportId === 1 ? NFL_PLAYOFF_ROUND_NAMES[playoffRound] : NBA_PLAYOFF_ROUND_NAMES[playoffRound]) || `Playoff Round ${playoffRound}`;

    function getMaxAvailableWeek(sport: number | null, playoffs: PlayoffState | null, roundComplete: boolean, week: number): number {
        if (!sport) return 1000;
        
        // If playoffs are not enabled, max week is end of regular season
        if (!playoffs?.is_enabled) {
            return NUM_WEEKS;
        }

        // Don't allow going past final round
        if (playoffs.current_round >= FINAL_PLAYOFF_ROUND) {
            return NUM_WEEKS + FINAL_PLAYOFF_ROUND;
        }

        // If current round is complete, enable next round
        if (roundComplete && week === NUM_WEEKS + playoffs.current_round) {
            console.log('Round complete! Enabling next round: ', NUM_WEEKS + playoffs.current_round + 1);
            return NUM_WEEKS + playoffs.current_round + 1;
        }

        // Otherwise, max week is current playoff round
        return NUM_WEEKS + playoffs.current_round;
    }

    function previousWeek() {
        if (currentWeek > 1) {
            dispatch('weekChanged', { week: currentWeek - 1 });
        }
    }

    function nextWeek() {
        if (currentWeek < maxAvailableWeek) {
            dispatch('weekChanged', { week: currentWeek + 1 });
        }
    }

    function selectWeek(week: number) {
        dispatch('weekChanged', { week });
        weekDropdownOpen = false;
    }

    function getWeekLabel(week: number): string {
        if (!sportId) return `WEEK ${week}`;
        
        const round = week > NUM_WEEKS ? week - NUM_WEEKS : 0;

        if (round > 0) {
            // Playoff round labels
            return sportId === 1
                ? NFL_PLAYOFF_ROUND_NAMES[round] || `Playoff Round ${round}`
                : NBA_PLAYOFF_ROUND_NAMES[round] || `Playoff Round ${round}`;
        }

        // Regular season week labels with date ranges
        const dateRange = weekDateRanges.get(week);
        if (!dateRange) return `WEEK ${week}`;
        return `WEEK ${week} (${sportId === 1 ? formatNFLDateRange(dateRange.startDate, dateRange.endDate) : formatNBADateRange(dateRange.startDate, dateRange.endDate)})`;
    }

    function handleKeydown(event: KeyboardEvent) {
        // Don't interfere if user is typing in an input or textarea
        if (event.target instanceof HTMLInputElement || event.target instanceof HTMLTextAreaElement) {
            return;
        }

        // Close dropdown with Escape
        if (event.key === 'Escape' && weekDropdownOpen) {
            weekDropdownOpen = false;
            event.preventDefault();
            return;
        }

        // Navigate weeks with Arrow keys
        if (event.key === 'ArrowLeft') {
            previousWeek();
            event.preventDefault();
        } else if (event.key === 'ArrowRight') {
            nextWeek();
            event.preventDefault();
        }
    }

    // Setup global keydown listener
    onMount(() => {
        window.addEventListener('keydown', handleKeydown);
    });

    // Cleanup listener on destroy
    onDestroy(() => {
        window.removeEventListener('keydown', handleKeydown);
    });
</script>

<div class="flex items-center justify-between gap-4 pb-4 border-b-2 border-primary-700">
    <!-- Previous Week Button -->
    <button
        on:click={previousWeek}
        disabled={currentWeek === 1}
        class="p-2 rounded-lg bg-primary-800 hover:bg-primary-600 border-2 border-primary-600 transition-colors cursor-pointer disabled:opacity-30 disabled:cursor-not-allowed"
        title="Previous Week"
    >
        <svg class="w-6 h-6 text-neutral" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
    </button>

    <!-- Week Selector -->
    <div class="flex-1 relative">
        <!-- Week Dropdown Button -->
        <button
            type="button"
            on:click={() => weekDropdownOpen = !weekDropdownOpen}
            class="w-full text-center text-md sm:text-2xl font-heading font-bold bg-primary-800 border-2 border-primary-600 rounded-lg px-3 sm:px-4 py-2 text-neutral transition-colors hover:bg-primary-600 flex justify-between items-center cursor-pointer"
            class:border-primary-400={weekDropdownOpen}
        >
            <span class="flex-1 truncate">{currentLabel}</span>
            <svg class="w-5 h-5 ml-2 shrink-0 transition-transform" class:rotate-180={weekDropdownOpen} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
        </button>
        
        {#if weekDropdownOpen}
            <div class="absolute z-10 w-full mt-1 bg-primary-800 border-2 border-primary-600 rounded-md shadow-lg max-h-60 overflow-auto">
                <!-- Regular Season Weeks -->
                {#if sportId}
                    {#each Array(NUM_WEEKS) as _, i}
                        <button
                            type="button"
                            on:click={() => selectWeek(i + 1)}
                            class="w-full px-4 py-3 text-center text-neutral hover:bg-primary-700 transition-colors font-heading font-bold text-base sm:text-lg cursor-pointer"
                            class:bg-primary-700={currentWeek === i + 1}
                        >
                            {getWeekLabel(i + 1)}
                        </button>
                    {/each}
                {/if}

                <!-- Playoff Rounds (if enabled) -->
                {#if playoffState?.is_enabled && sportId}
                    <div class="border-t-2 border-primary-600 my-1"></div>
                    {#each Array(playoffState.current_round) as _, i}
                        {@const week = NUM_WEEKS + 1 + i}
                        <button
                            type="button"
                            on:click={() => selectWeek(week)}
                            class="w-full px-4 py-3 text-center text-neutral hover:bg-primary-700 transition-colors font-heading font-bold text-base sm:text-lg cursor-pointer"
                            class:bg-primary-700={currentWeek === week}
                        >
                            {getWeekLabel(week)}
                        </button>
                    {/each}
                {/if}

                <!-- Message for enabling playoffs -->
                {#if canEnablePlayoffs && !playoffState?.is_enabled}
                    <div class="border-t-2 border-primary-600 my-1"></div>
                    <div class="px-4 py-3 text-center text-neutral/50 font-sans text-sm">
                        Start playoffs to unlock
                    </div>
                {/if}
            </div>
        {/if}
    </div>

    <!-- Next Week Button -->
    <button
        on:click={nextWeek}
        disabled={sportId !== null && currentWeek >= maxAvailableWeek}
        class="p-2 rounded-lg bg-primary-800 hover:bg-primary-600 border-2 border-primary-600 transition-colors cursor-pointer disabled:opacity-30 disabled:cursor-not-allowed"
        title="Next Week"
    >
        <svg class="w-6 h-6 text-neutral" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
    </button>
</div>