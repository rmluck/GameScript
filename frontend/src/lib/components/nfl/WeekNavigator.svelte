<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { getWeekDateRangesFromGames, formatDateRange, type WeekDateRange } from '$lib/utils/nfl/dates';
    import { PLAYOFF_ROUNDS, PLAYOFF_ROUND_NAMES } from '$types';
    import type { Game, PlayoffState } from '$types';

    export let currentWeek: number;
    export let allGames: Game[] = [];
    export let playoffState: PlayoffState | null = null;
    export let canEnablePlayoffs: boolean = false;

    const dispatch = createEventDispatcher();
    
    let weekDropdownOpen = false;
    let weekDateRanges: Map<number, WeekDateRange> = new Map();

    $: weekDateRanges = getWeekDateRangesFromGames(allGames);
    $: maxAvailableWeek = getMaxAvailableWeek();

    // Make currentLabel properly reactive to currentWeek changes
    $: currentLabel = currentWeek > 18 
        ? (PLAYOFF_ROUND_NAMES[currentWeek - 18] || `Playoff Round ${currentWeek - 18}`)
        : `Week ${currentWeek}`;

    function getMaxAvailableWeek(): number {
        if (!playoffState?.is_enabled) {
            return 18;
        }
        return 18 + playoffState.current_round;
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
        if (week > 18) {
            const round = week - 18;
            return PLAYOFF_ROUND_NAMES[round] || `Playoff Round ${round}`;
        }

        const dateRange = weekDateRanges.get(week);
        if (!dateRange) return `WEEK ${week}`;
        return `WEEK ${week} (${formatDateRange(dateRange.startDate, dateRange.endDate)})`;
    }
</script>

<div class="flex items-center justify-between gap-4 pb-4 border-b-2 border-primary-700">
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

    <div class="flex-1 relative">
        <button
            type="button"
            on:click={() => weekDropdownOpen = !weekDropdownOpen}
            class="w-full text-center text-xl sm:text-2xl font-heading font-bold bg-primary-800 border-2 border-primary-600 rounded-lg px-3 sm:px-4 py-2 text-neutral transition-colors hover:bg-primary-600 flex justify-between items-center cursor-pointer"
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
                {#each Array(18) as _, i}
                    <button
                        type="button"
                        on:click={() => selectWeek(i + 1)}
                        class="w-full px-4 py-3 text-center text-neutral hover:bg-primary-700 transition-colors font-heading font-bold text-base sm:text-lg cursor-pointer"
                        class:bg-primary-700={currentWeek === i + 1}
                    >
                        {getWeekLabel(i + 1)}
                    </button>
                {/each}

                <!-- Playoff Rounds (if enabled) -->
                {#if playoffState?.is_enabled}
                    <div class="border-t-2 border-primary-600 my-1"></div>
                    {#each Array(playoffState.current_round) as _, i}
                        {@const week = 19 + i}
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

                <!-- Show locked playoff rounds if all regular season is done but playoffs not enabled -->
                {#if canEnablePlayoffs && !playoffState?.is_enabled}
                    <div class="border-t-2 border-primary-600 my-1"></div>
                    <div class="px-4 py-3 text-center text-neutral/50 font-sans text-sm">
                        Start playoffs to unlock
                    </div>
                {/if}
            </div>
        {/if}
    </div>

    <button
        on:click={nextWeek}
        disabled={currentWeek === maxAvailableWeek}
        class="p-2 rounded-lg bg-primary-800 hover:bg-primary-600 border-2 border-primary-600 transition-colors cursor-pointer disabled:opacity-30 disabled:cursor-not-allowed"
        title="Next Week"
    >
        <svg class="w-6 h-6 text-neutral" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
    </button>
</div>