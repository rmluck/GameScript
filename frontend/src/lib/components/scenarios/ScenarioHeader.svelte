<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import type { Scenario } from '$types';

    export let scenario: Scenario;
    export let saveStatus: 'idle' | 'saving' | 'saved' | 'error' = 'idle';

    const dispatch = createEventDispatcher();
</script>

<div class="flex items-center justify-between gap-2 sm:gap-4 lg:gap-8 mb-6">
    <!-- Left: Scenario Name & Breadcrumb -->
    <div class="flex-1 min-w-0">
        <nav class="text-sm mb-2">
            <a href="/scenarios" class="text-primary-400 hover:text-primary-300 hover:underline transition-all duration-200">
                ← Scenarios
            </a>
        </nav>
        <h1 class="text-xl sm:text-2xl lg:text-3xl font-heading font-bold text-neutral truncate">
            {scenario.name}
        </h1>
        <p class="text-neutral/70 text-xs sm:text-sm mt-1">
            {scenario.sport_short_name} - {scenario.season_start_year}{scenario.season_end_year ? `-${scenario.season_end_year}` : ''} Season
        </p>
    </div>

    <!-- Right: Save Status & Action Buttons -->
    <div class="flex items-center gap-2 sm:gap-3 lg:gap-4 shrink-0">
        <!-- Save Status Indicator -->
        {#if saveStatus != 'idle'}
            <div class="hidden sm:flex items-center gap-2 text-sm">
                {#if saveStatus === 'saving'}
                    <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-primary-400"></div>
                    <span class="text-neutral/70">Saving...</span>
                {:else if saveStatus === 'saved'}
                    <svg class="w-4 h-4 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                    </svg>
                    <span class="text-green-400">Saved</span>
                {:else if saveStatus === 'error'}
                    <svg class="w-4 h-4 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                    <span class="text-red-400">Error</span>
                {/if}
            </div>
        {/if}

        <!-- Info Button -->
        <button
            on:click={() => dispatch('openInfo')}
            class="p-1.5 sm:p-2 rounded-lg bg-primary-800/60 hover:bg-primary-700 border-2 border-primary-600 transition-colors cursor-pointer"
            title="Help & Information"
        >
            <svg class="w-4 h-4 sm:w-5 sm:h-5 text-neutral" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
        </button>

        <!-- Settings Button -->
        <button
            on:click={() => dispatch('openSettings')}
            class="p-1.5 sm:p-2 rounded-lg bg-primary-800/60 hover:bg-primary-700 border-2 border-primary-600 transition-colors cursor-pointer"
            title="Scenario Settings"
        >
            <svg class="w-4 h-4 sm:w-5 sm:h-5 text-neutral" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
        </button>
    </div>
</div>