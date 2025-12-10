<script lang="ts">
    import type { Team } from '$types';

    export let teams: Team[];

    let hoveredTeamId: number | null = null;
</script>

<div class="mt-6 pt-6 border-t-2 border-primary-700">
    <h3 class="text-xl font-heading font-bold text-primary-700 mb-3 uppercase tracking-wide">
        Teams on Bye
    </h3>
    
    <div class="grid grid-cols-3 md:grid-cols-6 lg:grid-cols-8 gap-3">
        {#each teams as team}
            <div
                on:mouseenter={() => hoveredTeamId = team.id}
                on:mouseleave={() => hoveredTeamId = null}
                aria-label={`Team on bye: ${team.city} ${team.name}`}
                role="tooltip"
                class="relative p-2 rounded-lg border-2 transition-all"
                style={`background-color: transparent; border-color: #${team.primary_color};`}
            >
                <div class="flex items-center justify-center">
                    {#if team.logo_url}
                        <img 
                            src={team.logo_url}
                            alt={team.abbreviation}
                            class="w-8 h-8 object-contain"
                        />
                    {/if}
                </div>

                <!-- Hover Tooltip -->
                {#if hoveredTeamId === team.id}
                    <div class="absolute z-10 left-1/2 -translate-x-1/2 bottom-full mb-2 px-3 py-1.5 bg-primary-950 border border-primary-600 rounded-lg shadow-xl whitespace-nowrap">
                        <span class="text-sm font-sans font-semibold text-neutral">
                            {team.city} {team.name}
                        </span>
                    </div>
                {/if}
            </div>
        {/each}
    </div>
</div>