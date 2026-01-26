<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import type { NFLConferenceStandings, NFLPlayoffSeed } from '$types';

    // Props
    export let standings: NFLConferenceStandings;
    export let conference: 'AFC' | 'NFC';
    type ViewMode = 'conference' | 'division';
    export let viewMode: ViewMode = 'conference';

    // Event dispatcher
    const dispatch = createEventDispatcher();

    // Sort teams into groups
    $: divisionWinners = standings.playoff_seeds.slice(0, 4);
    $: wildCardTeams = standings.playoff_seeds.slice(4, 7);
    $: nonPlayoffTeams = standings.playoff_seeds.slice(7);

    $: orderedDivisions = ['North', 'South', 'East', 'West'].filter(div => 
        standings.divisions[`${conference} ${div}`]
    );

    // Map for quick seed lookup
    $: teamSeedMap = new Map(
        standings.playoff_seeds.map(seed => [seed.team_id, seed.seed])
    );

    function formatRecord(wins: number, losses: number, ties: number): string {
        return ties > 0 ? `${wins}-${losses}-${ties}` : `${wins}-${losses}`;
    }

    function getTeamSeed(teamId: number): number | undefined {
        return teamSeedMap.get(teamId);
    }

    function handleMouseEnter(e: MouseEvent, primaryColor: string) {
        // Set background and text colors on hover
        const target = e.currentTarget as HTMLElement;
        target.style.backgroundColor = `#${primaryColor}90`;
        target.querySelectorAll('span, div').forEach(el => {
            (el as HTMLElement).style.color = 'white';
        });
    }

    function handleMouseLeave(e: MouseEvent) {
        // Reset background and text colors on hover leave
        const target = e.currentTarget as HTMLElement;
        target.style.backgroundColor = 'transparent';
        target.querySelectorAll('span, div').forEach(el => {
            (el as HTMLElement).style.color = '';
        });
    }

    function openTeamModal(team: NFLPlayoffSeed) {
        dispatch('openTeamModal', { team });
    }

    function convertToPlayoffSeed(team: any): NFLPlayoffSeed {
        return {
            seed: team.seed,
            team_id: team.team_id,
            team_name: team.team_name,
            team_city: team.team_city,
            team_abbr: team.team_abbr,
            wins: team.wins,
            losses: team.losses,
            ties: team.ties,
            win_pct: team.win_pct,
            is_division_winner: team.is_division_winner,
            logo_url: team.logo_url,
            team_primary_color: team.team_primary_color,
            team_secondary_color: team.team_secondary_color,
            home_wins: team.home_wins,
            home_losses: team.home_losses,
            home_ties: team.home_ties,
            away_wins: team.away_wins,
            away_losses: team.away_losses,
            away_ties: team.away_ties,
            conference_wins: team.conference_wins,
            conference_losses: team.conference_losses,
            conference_ties: team.conference_ties,
            division_wins: team.division_wins,
            division_losses: team.division_losses,
            division_ties: team.division_ties,
            conference_games_back: team.conference_games_back,
            division_games_back: team.division_games_back,
            points_for: team.points_for,
            points_against: team.points_against,
            point_diff: team.point_diff,
            strength_of_schedule: team.strength_of_schedule,
            strength_of_victory: team.strength_of_victory
        };
    }
</script>

<div class="bg-neutral border-2 border-primary-700 rounded-lg px-4 py-6 w-full">
    <!-- Header -->
    <div class="flex items-center justify-between pb-6 mb-4 border-b-2 border-primary-700">
        <h2 class="text-xl font-heading font-bold uppercase tracking-wide"
            style="color: {conference === 'AFC' ? '#C8102E' : '#013369'}">
            {conference}
        </h2>

        <!-- View Toggle -->
        <div class="flex bg-primary-800 border-2 border-primary-600 rounded-lg p-1 gap-1">
            <button
                on:click={() => viewMode = 'conference'}
                class="p-1.5 text-xs sm:text-sm font-sans font-semibold text-neutral rounded transition-colors cursor-pointer"
                class:bg-primary-600={viewMode === 'conference'}
                class:hover:bg-primary-700={viewMode !== 'conference'}
            >
                Conference
            </button>
            <button
                on:click={() => viewMode = 'division'}
                class="p-1.5 text-xs sm:text-sm font-sans font-semibold text-neutral rounded transition-colors cursor-pointer"
                class:bg-primary-600={viewMode === 'division'}
                class:hover:bg-primary-700={viewMode !== 'division'}
            >
                Division
            </button>
        </div>
    </div>

    <!-- Conference View -->
    {#if viewMode === 'conference'}
        <div class="space-y-4">
            <!-- Division Winners (Seeds 1-4) -->
            <div>
                <h3 class="text-lg font-sans font-bold text-primary-700 uppercase tracking-wide mb-2 px-2">
                    Division Winners
                </h3>

                <div class="space-y-1">
                    {#each divisionWinners as seed}
                        <button class="w-full flex items-center gap-2 px-2 py-2 rounded transition-colors cursor-pointer"
                            on:mouseenter={(e) => handleMouseEnter(e, seed.team_primary_color)}
                            on:mouseleave={handleMouseLeave}
                            on:click={() => openTeamModal(seed)}
                        >
                            <span class="text-sm font-heading font-bold text-primary-600 w-6">
                                {seed.seed}
                            </span>
                            <img 
                                src={seed.logo_url} 
                                alt={seed.team_abbr}
                                class="w-6 h-6 object-contain"
                            />
                            <div class="flex-1 min-w-0">
                                <div class="text-left text-sm font-sans font-semibold text-black truncate">
                                    {seed.team_name}
                                </div>
                            </div>
                            <span class="text-sm font-heading font-bold text-black whitespace-nowrap">
                                {formatRecord(seed.wins, seed.losses, seed.ties)}
                            </span>
                        </button>
                    {/each}
                </div>
            </div>

            <!-- Wild Card (Seeds 5-7) -->
            {#if wildCardTeams.length > 0}
                <div class="border-t border-primary-700/50 pt-4">
                    <h3 class="text-lg font-sans font-bold text-primary-700 uppercase tracking-wide mb-2 px-2">
                        Wild Card
                    </h3>

                    <div class="space-y-1">
                        {#each wildCardTeams as seed}
                            <button class="w-full flex items-center gap-2 px-2 py-2 rounded transition-colors cursor-pointer"
                                on:mouseenter={(e) => handleMouseEnter(e, seed.team_primary_color)}
                                on:mouseleave={handleMouseLeave}
                                on:click={() => openTeamModal(seed)}
                            >
                                <span class="text-sm font-heading font-bold text-primary-600 w-6">
                                    {seed.seed}
                                </span>
                                <img 
                                    src={seed.logo_url} 
                                    alt={seed.team_abbr}
                                    class="w-6 h-6 object-contain"
                                />
                                <div class="flex-1 min-w-0">
                                    <div class="text-left text-sm font-sans font-semibold text-black truncate">
                                        {seed.team_name}
                                    </div>
                                </div>
                                <span class="text-sm font-heading font-bold text-black whitespace-nowrap">
                                    {formatRecord(seed.wins, seed.losses, seed.ties)}
                                </span>
                            </button>
                        {/each}
                    </div>
                </div>
            {/if}

            <!-- Non-Playoff Teams (Seeds 8-16) -->
            {#if nonPlayoffTeams.length > 0}
                <div class="border-t border-primary-700/50 pt-4 opacity-60">
                    <h3 class="text-lg font-sans font-bold text-primary-700 uppercase tracking-wide mb-2 px-2">
                        Out of Playoffs
                    </h3>

                    <div class="space-y-1">
                        {#each nonPlayoffTeams as seed}
                            <button class="w-full flex items-center gap-1 px-1 py-2 rounded transition-colors cursor-pointer"
                                on:mouseenter={(e) => handleMouseEnter(e, seed.team_primary_color)}
                                on:mouseleave={handleMouseLeave}
                                on:click={() => openTeamModal(seed)}
                            >
                                <span class="text-sm font-heading font-bold text-primary-600 w-6">
                                    {seed.seed}
                                </span>
                                <img 
                                    src={seed.logo_url} 
                                    alt={seed.team_abbr}
                                    class="w-6 h-6 object-contain"
                                />
                                <div class="flex-1 min-w-0">
                                    <div class="text-left text-sm font-sans font-semibold text-black truncate">
                                        {seed.team_name}
                                    </div>
                                </div>
                                <span class="text-sm font-heading font-bold text-black whitespace-nowrap">
                                    {formatRecord(seed.wins, seed.losses, seed.ties)}
                                </span>
                            </button>
                        {/each}
                    </div>
                </div>
            {/if}
        </div>
    {/if}

    <!-- Division View -->
    {#if viewMode === 'division'}
        <div class="space-y-4">
            {#each orderedDivisions as divisionName}
                {@const fullDivisionName = `${conference} ${divisionName}`}
                {@const divisionTeams = standings.divisions[fullDivisionName] || []}
                
                <div>
                    <h3 class="text-lg font-sans font-bold text-primary-700 uppercase tracking-wide mb-2 px-2">
                        {conference} {divisionName}
                    </h3>
                    
                    <div class="space-y-1">
                        {#each divisionTeams as team, index}
                            <button class="w-full flex items-center gap-2 px-2 py-2 rounded transition-colors cursor-pointer"
                                on:mouseenter={(e) => handleMouseEnter(e, team.team_primary_color)}
                                on:mouseleave={handleMouseLeave}
                                on:click={() => openTeamModal(convertToPlayoffSeed(team))}
                                class:opacity-60={(getTeamSeed(team.team_id) !== undefined) && getTeamSeed(team.team_id)! > 7}
                            >
                                <span class="text-sm font-heading font-bold text-primary-600 w-6">
                                    {getTeamSeed(team.team_id) ?? '-'}
                                </span>
                                {#if team.logo_url}
                                    <img 
                                        src={team.logo_url} 
                                        alt={team.team_abbr}
                                        class="w-6 h-6 object-contain"
                                    />
                                {/if}
                                <div class="flex-1 min-w-0">
                                    <div class="text-left text-sm font-sans font-semibold text-black truncate">
                                        {team.team_name}
                                    </div>
                                </div>
                                <span class="text-sm font-heading font-bold text-black whitespace-nowrap">
                                    {formatRecord(team.wins, team.losses, team.ties)}
                                </span>
                            </button>
                        {/each}
                    </div>
                </div>
            {/each}
        </div>
    {/if}
</div>

<!-- Expanded View (stacked layout - more detailed stats) -->
<style>
    @media (max-width: 1279px) {
        /* Show expanded stats when stacked */
        :global(.standings-expanded) {
            display: block;
        }
    }
</style>