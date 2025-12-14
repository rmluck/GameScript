<script lang="ts">
    import { createEventDispatcher } from 'svelte';

    export let isOpen = false;
    export let sport: string | undefined = undefined;

    console.log('sport in ScenarioInfo:', sport);

    const dispatch = createEventDispatcher();

    function close() {
        isOpen = false;
        dispatch('close');
    }

    function handleClickOutside(event: MouseEvent) {
        if (event.target === event.currentTarget) {
            close();
        }
    }
</script>

{#if isOpen}
    <div 
        class="fixed inset-0 bg-black/70 flex items-center justify-center z-50 p-4"
        on:click={handleClickOutside}
        on:keydown={(e) => { if (e.key === 'Escape') close(); }}
        role="dialog"
        aria-modal="true"
        tabindex="-1"
    >
        <section
            class="bg-neutral border-4 border-primary-600 rounded-lg max-w-2xl w-full max-h-[90vh] overflow-hidden flex flex-col"
            role="document"
            tabindex="-1"
            aria-label="Scenario information modal"
        >
            <!-- Header -->
            <div class="border-b-4 border-primary-600 p-6 bg-primary-600/70">
                <div class="flex items-center justify-between">
                    <h2 class="text-2xl font-heading font-bold text-white">
                        Scenario Information
                    </h2>
                    <button
                        on:click={close}
                        class="p-2 rounded-lg hover:bg-white/10 transition-colors cursor-pointer"
                        aria-label="Close modal"
                    >
                        <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </div>
            </div>

            <!-- Content -->
            <div class="bg-primary-600/50 flex-1 overflow-y-auto p-6">
                {#if sport === 'NFL'}
                    <div class="space-y-4">
                        <div>
                            <h3 class="text-xl font-heading font-bold text-neutral mb-2">NFL Scenario Rules</h3>
                            <p class="text-neutral/95">
                                This scenario allows you to predict the outcomes of NFL games and see how they affect playoff standings.
                            </p>
                        </div>
                        
                        <div>
                            <h4 class="text-lg font-sans font-semibold text-neutral mb-2">How to Use</h4>
                            <ul class="list-disc list-inside space-y-1 text-neutral/95">
                                <li>Click on team buttons to select a winner</li>
                                <li>Enter predicted scores for tiebreaker scenarios</li>
                                <li>Changes are saved automatically</li>
                                <li>View updated standings in real-time</li>
                            </ul>
                        </div>

                        <div>
                            <h4 class="text-lg font-sans font-semibold text-neutral mb-2">Playoff Seeding</h4>
                            <ul class="list-disc list-inside space-y-1 text-neutral/95">
                                <li>Seeds 1-7 represent playoff positions</li>
                                <li>Seed 1 gets a first-round bye</li>
                                <li>Top 4 seeds are division winners</li>
                                <li>Seeds 5-7 are wild card teams</li>
                            </ul>
                        </div>
                    </div>
                {:else}
                    <div class="space-y-4">
                        <p class="text-neutral/95">
                            Information about this sport's scenario rules will appear here.
                        </p>
                    </div>
                {/if}
            </div>

            <!-- Footer -->
            <div class="bg-primary-600/50 p-4">
                <button
                    on:click={close}
                    class="w-full px-4 py-2 bg-primary-600 hover:bg-primary-500 text-white font-sans font-semibold rounded-lg transition-colors cursor-pointer"
                >
                    Close
                </button>
            </div>
        </section>
    </div>
{/if}