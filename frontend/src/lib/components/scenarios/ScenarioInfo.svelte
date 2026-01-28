<script lang="ts">
    import { createEventDispatcher } from 'svelte';

    // Props
    export let isOpen = false;
    export let sport: string | undefined = undefined;

    // Event dispatcher
    const dispatch = createEventDispatcher();

    function close() {
        isOpen = false;
        dispatch('close');
    }

    // Close modal when clicking outside content
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
                                    <li>Points scored</li>
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
                                <li>Points scored</li>
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
                {:else if sport === 'NBA'}
                    <div class="space-y-4">
                        <div>
                            <h3 class="text-xl font-heading font-bold text-neutral mb-2">NBA Scenario Rules</h3>
                            <p class="text-neutral/95">
                                This scenario allows you to predict the outcomes of NBA games and see how they affect playoff standings.
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
                                <li>Seeds 1-10 represent postseason positions</li>
                                <li>Seeds 1-6 get automatic playoff berths</li>
                                <li>Seeds 7-10 participate in single-elimination Play-In Tournament</li>
                                <li>Playoff series are best-of-7</li>
                            </ul>
                        </div>

                        <div class="space-y-3">
                            <h4 class="text-lg font-sans font-semibold text-neutral mb-2">Tiebreakers</h4>

                            <div class="space-y-2 text-neutral ml-5">
                                <div>
                                    <p class="font-semibold text-neutral">Division Winners:</p>
                                    <ul class="list-disc list-inside space-y-1 ml-4">
                                        <li>Ties to determine division winners must be broken before any other ties</li>
                                        <li>When a tie must be broken to determine a division winner, the results of the tie-break shall be used to determine only the division winner, and not for any other purpose</li>
                                    </ul>
                                </div>

                                <div>
                                    <p class="font-semibold text-neutral">Multi-Team Tie Breaking Process:</p>
                                    <ul class="list-disc list-inside space-y-1 ml-4">
                                        <li><span class="font-semibold">Complete Breaking:</span> If each tied team has a different win percentage or point differential under the applicable criterion, teams are ranked accordingly and no further tiebreakers are needed</li>
                                        <li><span class="font-semibold">Partial Breaking:</span> If one or more (but not all) teams have different performance under a criterion:
                                            <ul class="list-circle list-inside ml-6 mt-1">
                                                <li>Better performing team(s) get higher playoff position(s)</li>
                                                <li>Remaining tied teams restart the tiebreaker process from the beginning using two-team criteria (if two teams remain) or multi-team criteria (if three or more remain)</li>
                                            </ul>
                                        </li>
                                        <li><span class="font-semibold">Random Drawing:</span> If application of all criteria does not break the tie, playoff positions are determined by random drawing</li>
                                    </ul>
                                </div>
                            </div>

                            <div>
                                <h5 class="text-md font-sans font-semibold text-neutral mb-1">Two-Way Tiebreakers</h5>
                                <p class="text-neutral">In the case of a tie in regular season records involving only two teams, the following criteria will be utilized in the following order:</p>
                                <ol class="list-decimal list-inside space-y-1 text-neutral/95 ml-5">
                                    <li>Head-to-head record</li>
                                    <li>Division winner (this criterion is applied regardless of whether the tied teams are in the same division)</li>
                                    <li>Division record (only if tied teams are in same division)</li>
                                    <li>Conference record</li>
                                    <li>Win percentage vs. teams eligible for postseason in own conference</li>
                                    <li>Win percentage vs. teams eligible for postseason in opposing conference</li>
                                    <li>Point differential</li>
                                </ol>
                            </div>

                            <div>
                                <h5 class="text-md font-sans font-semibold text-neutral mb-1">Multi-Way Tiebreakers</h5>
                                <p class="text-neutral">In the case of a tie in regular season records involving more than two teams, the following criteria will be utilized in the following order:</p>
                                <ol class="list-decimal list-inside space-y-1 text-neutral/95 ml-5">
                                    <li>Division winner (this criterion is applied regardless of whether the tied teams are in the same division)</li>
                                    <li>Head-to-head record among tied teams</li>
                                    <li>Division record (only if all tied teams are in same division)</li>
                                    <li>Conference record</li>
                                    <li>Win percentage vs. teams eligible for postseason in own conference</li>
                                    <li>Point differential</li>
                                </ol>
                            </div>
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