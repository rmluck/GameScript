<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import type { NFLDraftPick, NBADraftPick } from '$types';

    // Props
    export let draftOrder: NFLDraftPick[] | NBADraftPick[];

    // Event dispatcher (if needed in the future)
    const dispatch = createEventDispatcher();

    function handleMouseEnter(e: MouseEvent, primaryColor: string) {
        const target = e.currentTarget as HTMLElement;
        target.style.backgroundColor = `#${primaryColor}90`;
        target.querySelectorAll('span, div').forEach(el => {
            (el as HTMLElement).style.color = 'white';
        });
    }

    function handleMouseLeave(e: MouseEvent) {
        const target = e.currentTarget as HTMLElement;
        target.style.backgroundColor = 'transparent';
        target.querySelectorAll('span, div').forEach(el => {
            (el as HTMLElement).style.color = '';
        });
    }
</script>

<div class="bg-neutral border-2 border-primary-700 rounded-lg px-4 py-6 w-full">
    <div class="flex items-center justify-between pb-6 mb-4 border-b-2 border-primary-700">
        <h2 class="text-2xl font-heading font-bold uppercase tracking-wide text-primary-700">
            Draft Order
        </h2>
    </div>

    <!-- Teams List -->
    {#if draftOrder.length === 0}
        <div class="text-center py-8">
            <p class="text-neutral/70 text-lg font-sans">No draft order available yet</p>
        </div>
    {:else}
        <div class="space-y-1">
            {#each draftOrder as pick}
                <div
                    class="flex items-center gap-3 px-2 py-2 rounded transition-colors"
                    on:mouseenter={(e) => handleMouseEnter(e, pick.team_primary_color)}
                    on:mouseleave={handleMouseLeave}
                    role="listitem"
                >
                    <div class="w-10 text-center">
                        <span class="text-lg font-heading font-bold text-primary-600">
                            {pick.pick}
                        </span>
                    </div>

                    {#if pick.logo_url}
                        <img 
                            src={pick.logo_url} 
                            alt={pick.team_abbr}
                            class="w-8 h-8 object-contain"
                        />
                    {/if}

                    <div class="flex-1 min-w-0">
                        <div class="text-sm font-sans font-semibold text-black truncate">
                            {pick.team_name}
                        </div>
                    </div>

                    <div class="text-right">
                        <span class="text-sm font-heading font-bold text-black whitespace-nowrap">
                            {pick.record}
                        </span>
                    </div>
                </div>
            {/each}
        </div>
    {/if}
</div>