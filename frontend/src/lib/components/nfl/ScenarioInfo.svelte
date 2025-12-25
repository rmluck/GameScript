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
            class="bg-linear-to-br from-primary-975  to-primary-950 border-4 border-primary-900 rounded-lg max-w-2xl w-full max-h-[90vh] overflow-hidden flex flex-col"
            role="document"
            tabindex="-1"
            aria-label="Scenario information modal"
        >
            <!-- Header -->
            <div class="shadow-md p-6 bg-primary-900/30">
                <div class="flex items-center justify-between">
                    <h2 class="text-2xl font-heading font-bold text-neutral">
                        SCENARIO INFORMATION
                    </h2>
                    <button
                        on:click={close}
                        class="p-2 rounded-lg hover:bg-white/10 transition-colors cursor-pointer"
                        aria-label="Close modal"
                    >
                        <svg class="w-6 h-6 text-neutral" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </div>
            </div>

            <!-- Content -->
            <div class="bg-linear-to-br from-primary-975  to-primary-950 flex-1 overflow-y-auto p-6">
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
                                <li>Simulate playoff scenarios</li>
                            </ul>
                        </div>

                        <div>
                            <h4 class="text-lg font-sans font-semibold text-neutral mb-2">Playoff Seeding</h4>
                            <ul class="list-disc list-inside space-y-1 text-neutral/95">
                                <li>Seeds 1-7 represent playoff positions</li>
                                <li>Seed 1 gets a first-round bye</li>
                                <li>Seeds 1-4 are division winners</li>
                                <li>Seeds 5-7 are wild card teams</li>
                            </ul>
                        </div>

                        <div class="space-y-3">
                            <h4 class="text-lg font-sans font-semibold text-neutral mb-2">Tiebreakers</h4>

                            <ul class="list-disc list-inside space-y-1 text-neutral/95">
                                <li>Only one team advances in any tiebreaking step. Remaining tied teams revert back to first step of applicable procedure.</li>
                                <li>In comparing records against common opponents amongst tied teams, best win percentage is deciding factor because the teams may have played an unequal number of games against common opponents.</li>
                                <li>To determine tiebreakers among division winners, apply inter-conference tiebreakers.</li>
                                <li>To determine tiebreakers among wild card teams, apply division tiebreakers if the involved teams are from the same division or conference tiebreakers otherwise.</li>
                                <li>Treats a 0-0 record as less than 0.000. So, a team that is 0-0 will be lower than a team that is 0-1. This is particularly important early in the season when few games have been played.</li>
                                <li>A game must have an outcome for it to count towards record against common opponents, as well as a team's strength of schedule and strength of victory.</li>
                            </ul>

                            <div>
                                <h5 class="text-md font-sans font-semibold text-neutral mb-1">Inter-Division Tiebreakers</h5>
                                <p class="text-neutral">If two teams in the same division have an identical win percentage, these are the tiebreaking steps to be taken until a winner is determined:</p>
                                <ol class="list-decimal list-inside space-y-1 text-neutral/95 ml-5">
                                    <li>Head-to-head record</li>
                                    <li>Division record</li>
                                    <li>Record vs. common opponents</li>
                                    <li>Conference record</li>
                                    <li>Strength of victory</li>
                                    <li>Strength of schedule</li>
                                    <li>Point differential</li>
                                    <li>Points for</li>
                                    <li>Points allowed</li>
                                    <li>Coin toss (random choice)</li>
                                </ol>
                                <p class="text-neutral">If three or more teams in the same division have an identical win percentage, these are the tiebreaking steps to be taken until a winner is determined:</p>
                                <ol class="list-decimal list-inside space-y-1 text-neutral/95 ml-5">
                                    <p class="text-neutral">Note: If at some point during these steps, at least one team is eliminated and there are only two teams left that remain tied, the tiebreaker restarts at Step 1 of the two-team format above. If at some point during these steps, at least one team is eliminated but there are still at least three teams left that remain tied, the tiebreaker restarts at Step 1 of these steps.</p>
                                    <li>Head-to-head record</li>
                                    <li>Division record</li>
                                    <li>Record vs. common opponents</li>
                                    <li>Conference record</li>
                                    <li>Strength of victory</li>
                                    <li>Strength of schedule</li>
                                    <li>Coin toss (random choice)</li>
                                </ol>
                            </div>

                            <h5 class="text-md font-sans font-semibold text-neutral mb-1">Inter-Conference Tiebreakers</h5>
                            <p class="text-neutral">If two teams in the same conference have an identical win percentage, these are the tiebreaking steps to be taken until a winner is determined:</p>
                            <ol class="list-decimal list-inside space-y-1 text-neutral/95 ml-5">
                                <p>If the tied teams are from the same division, apply the inter-division tiebreaker. Otherwise, continue with these steps.</p>
                                <li>Head-to-head record</li>
                                <li>Conference record</li>
                                <li>Record vs. common opponents (minimum 4 games)</li>
                                <li>Strength of victory</li>
                                <li>Strength of schedule</li>
                                <li>Point differential</li>
                                <li>Points for</li>
                                <li>Points allowed</li>
                                <li>Coin toss (random choice)</li>
                            </ol>
                            <p class="text-neutral">If three or more teams in the same division have an identical win percentage, these are the tiebreaking steps to be taken until a winner is determined:</p>
                            <ol class="list-decimal list-inside space-y-1 text-neutral/95 ml-5">
                                <p class="text-neutral">Note: If at some point during these steps, at least one team is eliminated and there are only two teams left that remain tied, the tiebreaker restarts at Step 1 of the two-team format above. If at some point during these steps, at least one team is eliminated but there are still at least three teams left that remain tied, the tiebreaker restarts at Step 1 of these steps.</p>
                                <li>Apply inter-division tiebreaker to eliminate all but the highest ranked team in each division involved in tiebreaker. The original seeding within the division upon application of the division tiebreaker remains the same for all subsequent applications of the procedure that are necessary.</li>
                                <li>Head-to-head sweep (applicable only if one team has defeated each of the others or if one team has lost to each of the others)</li>
                                <li>Conference record</li>
                                <li>Record vs. common opponents (minimum 4 games)</li>
                                <li>Strength of victory</li>
                                <li>Strength of schedule</li>
                                <li>Coin toss (random choice)</li>
                            </ol>
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
            <div class="bg-primary-900/30 p-4 shadow-md">
                <button
                    on:click={close}
                    class="w-full px-4 py-2 bg-primary-900 hover:bg-primary-500 text-lg text-white font-sans font-semibold rounded-lg transition-colors cursor-pointer"
                >
                    Close
                </button>
            </div>
        </section>
    </div>
{/if}