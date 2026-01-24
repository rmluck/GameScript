<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import type { Scenario } from '$types';

    export let scenario: Scenario;
    export let saveStatus: 'idle' | 'saving' | 'saved' | 'error' = 'idle';

    const dispatch = createEventDispatcher();

    let showCopiedMessage = false;
    let copiedTimeout: ReturnType<typeof setTimeout>;

    async function handleShare() {
        const url = window.location.href;

        // Trye native share API first (mobile devices)
        if (navigator.share) {
            try {
                await navigator.share({
                    title: `GameScript - ${scenario.name}`,
                    text: `Check out my scenario "${scenario.name}" on GameScript!`,
                    url: url
                });
                return;
            } catch (err) {
                // User cancelled or share failed, fall back to clipboard
                if ((err as Error).name !== 'AbortError') {
                    console.error('Native share failed:', err);
                }
            }
        }

        // Fallback to clipboard
        try {
            await navigator.clipboard.writeText(url);
            showCopiedMessage = true;

            // Clear any existing timeout
            if (copiedTimeout) clearTimeout(copiedTimeout);

            // Hide message after 2 seconds
            copiedTimeout = setTimeout(() => {
                showCopiedMessage = false;
            }, 2000);
        } catch (err) {
            console.error('Failed to copy to clipboard:', err);
            // Fallback for older browsers
            fallbackCopyToClipboard(url);
        }
    }

    function fallbackCopyToClipboard(text: string) {
        const textArea = document.createElement('textarea');
        textArea.value = text;
        textArea.style.position = 'fixed';
        textArea.style.left = '-999999px';
        document.body.appendChild(textArea);
        textArea.select();

        try {
            document.execCommand('copy');
            showCopiedMessage = true;

            if (copiedTimeout) clearTimeout(copiedTimeout);

            copiedTimeout = setTimeout(() => {
                showCopiedMessage = false;
            }, 2000);
        } catch (err) {
            console.error('Fallback copy failed:', err);
        }

        document.body.removeChild(textArea);
    }
</script>

<div class="flex items-center justify-between gap-2 sm:gap-4 lg:gap-8 mb-6">
    <!-- Left: Scenario Name & Breadcrumb -->
    <div class="flex-1 min-w-0">
        <!-- <nav class="text-sm mb-2">
            <a href="/scenarios" class="text-primary-400 hover:text-primary-300 hover:underline transition-all duration-200">
                ← Scenarios
            </a>
        </nav> -->
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

        <!-- Share Button -->
         <div class="relative">
            <button
                on:click={handleShare}
                class="p-1.5 sm:p-2 rounded-lg bg-primary-800/60 hover:bg-primary-700 border-2 border-primary-600 transition-colors cursor-pointer"
                title="Share Scenario"
            >
                <svg class="w-4 h-4 sm:w-5 sm:h-5 text-neutral" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z" />
                </svg>
            </button>

            <!-- Copied Message Tooltip -->
             {#if showCopiedMessage}
                <div class="absolute top-full right-0 mt-2 px-3 py-2 bg-green-500 text-white text-sm font-sans font-semibold rounded-md shadow-lg whitespace-nowrap z-50 animate-fade-in">
                    Link copied!
                    <div class="absolute -top-1 right-4 w-2 h-2 bg-green-500 transform rotate-45"></div>
                </div>
             {/if}
         </div>

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

<style>
    @keyframes fade-in {
        from {
            opacity: 0;
            transform: translateY(-4px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }

    .animate-fade-in {
        animation: fade-in 0.2s ease-out;
    }
</style>