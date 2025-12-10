<script lang="ts">
    import { createEventDispatcher } from 'svelte';

    export let currentWeek: number;
    export let maxWeek: number = 18; // Regular season weeks

    const dispatch = createEventDispatcher();
    
    let weekDropdownOpen = false;

    function previousWeek() {
        if (currentWeek > 1) {
            dispatch('weekChanged', { week: currentWeek - 1 });
        }
    }

    function nextWeek() {
        if (currentWeek < maxWeek) {
            dispatch('weekChanged', { week: currentWeek + 1 });
        }
    }

    function selectWeek(week: number) {
        dispatch('weekChanged', { week });
        weekDropdownOpen = false;
    }
</script>

<div class="flex items-center justify-between gap-4 pb-4 border-b-2 border-primary-700">
    <!-- Left: Previous Button -->
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

    <!-- Center: Custom Week Dropdown -->
    <div class="flex-1 relative">
        <button
            type="button"
            on:click={() => weekDropdownOpen = !weekDropdownOpen}
            class="w-full text-center text-2xl font-heading font-bold bg-primary-800 border-2 border-primary-600 rounded-lg px-4 py-2 text-neutral transition-colors hover:bg-primary-600 flex justify-between items-center cursor-pointer"
            class:border-primary-400={weekDropdownOpen}
        >
            <span class="flex-1">WEEK {currentWeek}</span>
            <svg class="w-5 h-5 transition-transform" class:rotate-180={weekDropdownOpen} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
        </button>
        
        {#if weekDropdownOpen}
            <div class="absolute z-10 w-full mt-1 bg-primary-800 border-2 border-primary-600 rounded-md shadow-lg max-h-60 overflow-auto">
                {#each Array(maxWeek) as _, i}
                    <button
                        type="button"
                        on:click={() => selectWeek(i + 1)}
                        class="w-full px-4 py-3 text-center text-neutral hover:bg-primary-700 transition-colors font-heading font-bold text-xl cursor-pointer"
                        class:bg-primary-700={currentWeek === i + 1}
                    >
                        WEEK {i + 1}
                    </button>
                {/each}
            </div>
        {/if}
    </div>

    <!-- Right: Next Button -->
    <button
        on:click={nextWeek}
        disabled={currentWeek === maxWeek}
        class="p-2 rounded-lg bg-primary-800 hover:bg-primary-600 border-2 border-primary-600 transition-colors cursor-pointer disabled:opacity-30 disabled:cursor-not-allowed"
        title="Next Week"
    >
        <svg class="w-6 h-6 text-neutral" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
    </button>
</div>