<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import type { ConferenceStandings, PlayoffSeed } from '$types';

    export let standings: ConferenceStandings;
    export let conference: 'AFC' | 'NFC';
    export let scenarioId: number;
    export let seasonId: number;

    type ViewMode = 'conference' | 'division';
    export let viewMode: ViewMode = 'conference';

    const dispatch = createEventDispatcher();

    $: divisionWinners = standings.playoff_seeds.slice(0, 4);
    $: wildCardTeams = standings.playoff_seeds.slice(4, 7);
    $: nonPlayoffTeams = standings.playoff_seeds.slice(7);

    $: orderedDivisions = ['North', 'South', 'East', 'West'].filter(div => 
        standings.divisions[`${conference} ${div}`]
    );

    $: teamSeedMap = new Map(
        standings.playoff_seeds.map(seed => [seed.team_id, seed.seed])
    );

    function formatRecord(wins: number, losses: number, ties: number): string {
        return ties > 0 ? `${wins}-${losses}-${ties}` : `${wins}-${losses}`;
    }

    function formatWinPct(winPct: number): string {
        return winPct.toFixed(3);
    }

    function formatGamesBack(gb: number): string {
        if (gb === 0) return '—';
        return gb.toFixed(1);
    }

    function formatPointDiff(diff: number): string {
        if (diff > 0) return `+${diff}`;
        return diff.toString();
    }

    function getTeamSeed(teamId: number): number | undefined {
        return teamSeedMap.get(teamId);
    }

    function handleMouseEnter(e: MouseEvent, primaryColor: string) {
        const target = e.currentTarget as HTMLElement;
        target.style.backgroundColor = `#${primaryColor}90`;
        // Change all text elements to white
        target.querySelectorAll('span, div').forEach(el => {
            (el as HTMLElement).style.color = 'white';
        });
    }

    function handleMouseLeave(e: MouseEvent) {
        const target = e.currentTarget as HTMLElement;
        target.style.backgroundColor = 'transparent';
        // Reset all text colors
        target.querySelectorAll('span, div').forEach(el => {
            (el as HTMLElement).style.color = '';
        });
    }

    function openTeamModal(team: PlayoffSeed) {
        dispatch('openTeamModal', { team });
    }

    function convertToPlayoffSeed(team: any): PlayoffSeed {
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
            point_diff: team.point_diff
        };
    }
</script>

<div class="bg-neutral border-2 border-primary-700 rounded-lg px-4 py-6 w-full">
    <!-- Header -->
    <div class="flex items-center justify-between pb-6 mb-4 border-b-2 border-primary-700">
        <h2 class="text-2xl font-heading font-bold uppercase tracking-wide"
            style="color: {conference === 'AFC' ? '#C8102E' : '#013369'}">
            {conference}
        </h2>

        <!-- View Toggle -->
        <div class="flex bg-primary-800 border-2 border-primary-600 rounded-lg p-1 gap-1">
            <button
                on:click={() => viewMode = 'conference'}
                class="px-3 py-1.5 text-xs sm:text-sm font-sans font-semibold text-neutral rounded transition-colors cursor-pointer"
                class:bg-primary-600={viewMode === 'conference'}
                class:hover:bg-primary-700={viewMode !== 'conference'}
            >
                Conference
            </button>
            <button
                on:click={() => viewMode = 'division'}
                class="px-3 py-1.5 text-xs sm:text-sm font-sans font-semibold text-neutral rounded transition-colors cursor-pointer"
                class:bg-primary-600={viewMode === 'division'}
                class:hover:bg-primary-700={viewMode !== 'division'}
            >
                Division
            </button>
        </div>
    </div>

    <!-- Conference View with Stats -->
    <div class="overflow-x-auto">
        {#if viewMode === 'conference'}
            <div class="space-y-4 min-w-[500px]">
                <!-- Division Winners (Seeds 1-4) -->
                <div>
                    <h3 class="text-lg font-sans font-bold text-primary-700 uppercase tracking-wide mb-2 px-2">
                        Division Winners
                    </h3>
                    
                    <!-- Header Row -->
                    <div class="grid grid-cols-15 gap-2 px-2 border-b border-primary-700/30 text-sm font-sans font-bold text-black uppercase">
                        <div class="col-span-1">Seed</div>
                        <div class="col-span-3">Team</div>
                        <div class="col-span-2 text-center">Record</div>
                        <div class="col-span-1 text-center">PCT</div>
                        <div class="col-span-2 text-center">Conf</div>
                        <div class="col-span-2 text-center">Div</div>
                        <div class="col-span-1 text-center">GB</div>
                        <div class="col-span-1 text-center">Diff</div>
                        <div class="col-span-1 text-center">PF</div>
                        <div class="col-span-1 text-center">PA</div>
                    </div>

                    <!-- Data Rows -->
                    <div class="space-y-1 mt-2">
                        {#each divisionWinners as seed}
                            <button class="w-full grid grid-cols-15 gap-2 px-2 py-2 rounded transition-colors cursor-pointer"
                                on:mouseenter={(e) => handleMouseEnter(e, seed.team_primary_color)}
                                on:mouseleave={handleMouseLeave}
                                on:click={() => openTeamModal(seed)}
                            >
                                <div class="col-span-1">
                                    <span class="text-sm font-heading font-bold text-primary-600">
                                        {seed.seed}
                                    </span>
                                </div>
                                <div class="col-span-3 flex items-center gap-2">
                                    {#if seed.logo_url}
                                        <img 
                                            src={seed.logo_url} 
                                            alt={seed.team_abbr}
                                            class="w-6 h-6 object-contain"
                                        />
                                    {/if}
                                    <div class="text-sm font-sans font-semibold text-black truncate">
                                        {seed.team_abbr}
                                    </div>
                                </div>
                                <div class="col-span-2 text-center">
                                    <span class="text-sm font-heading font-bold text-black">
                                        {formatRecord(seed.wins, seed.losses, seed.ties)}
                                    </span>
                                </div>
                                <div class="col-span-1 text-center">
                                    <span class="text-sm font-sans text-black">
                                        {formatWinPct(seed.win_pct)}
                                    </span>
                                </div>
                                <div class="col-span-2 text-center">
                                    <span class="text-sm font-sans text-black">
                                        {formatRecord(seed.conference_wins, seed.conference_losses, seed.conference_ties)}
                                    </span>
                                </div>
                                <div class="col-span-2 text-center">
                                    <span class="text-sm font-sans text-black">
                                        {formatRecord(seed.division_wins, seed.division_losses, seed.division_ties)}
                                    </span>
                                </div>
                                <div class="col-span-1 text-center">
                                    <span class="text-sm font-sans text-black">
                                        {formatGamesBack(seed.conference_games_back)}
                                    </span>
                                </div>
                                <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatPointDiff(seed.point_diff)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {seed.points_for}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {seed.points_against}
                                        </span>
                                    </div>
                            </button>
                        {/each}
                    </div>
                </div>

                <!-- Wild Card -->
                {#if wildCardTeams.length > 0}
                    <div>
                        <h3 class="text-lg font-sans font-bold text-primary-700 uppercase tracking-wide mb-2 px-2">
                            Wild Card
                        </h3>
                        
                        <!-- Header Row -->
                        <div class="grid grid-cols-15 gap-2 px-2 border-b border-primary-700/30 text-sm font-sans font-bold text-black uppercase">
                            <div class="col-span-1">Seed</div>
                            <div class="col-span-3">Team</div>
                            <div class="col-span-2 text-center">Record</div>
                            <div class="col-span-1 text-center">PCT</div>
                            <div class="col-span-2 text-center">Conf</div>
                            <div class="col-span-2 text-center">Div</div>
                            <div class="col-span-1 text-center">GB</div>
                            <div class="col-span-1 text-center">Diff</div>
                            <div class="col-span-1 text-center">PF</div>
                            <div class="col-span-1 text-center">PA</div>
                        </div>

                        <div class="space-y-1 mt-2">
                            {#each wildCardTeams as seed}
                                <button class="w-full grid grid-cols-15 gap-2 px-2 py-2 rounded transition-colors cursor-pointer"
                                    on:mouseenter={(e) => handleMouseEnter(e, seed.team_primary_color)}
                                    on:mouseleave={handleMouseLeave}
                                    on:click={() => openTeamModal(seed)}
                                >
                                    <div class="col-span-1">
                                        <span class="text-sm font-heading font-bold text-primary-600">
                                            {seed.seed}
                                        </span>
                                    </div>
                                    <div class="col-span-3 flex items-center gap-2">
                                        {#if seed.logo_url}
                                            <img 
                                                src={seed.logo_url} 
                                                alt={seed.team_abbr}
                                                class="w-6 h-6 object-contain"
                                            />
                                        {/if}
                                        <div class="text-sm font-sans font-semibold text-black truncate">
                                            {seed.team_abbr}
                                        </div>
                                    </div>
                                    <div class="col-span-2 text-center">
                                        <span class="text-sm font-heading font-bold text-black">
                                            {formatRecord(seed.wins, seed.losses, seed.ties)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatWinPct(seed.win_pct)}
                                        </span>
                                    </div>
                                    <div class="col-span-2 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatRecord(seed.conference_wins, seed.conference_losses, seed.conference_ties)}
                                        </span>
                                    </div>
                                    <div class="col-span-2 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatRecord(seed.division_wins, seed.division_losses, seed.division_ties)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatGamesBack(seed.conference_games_back)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatPointDiff(seed.point_diff)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {seed.points_for}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {seed.points_against}
                                        </span>
                                    </div>
                                </button>
                            {/each}
                        </div>
                    </div>
                {/if}

                <!-- Non-Playoff Teams -->
                {#if nonPlayoffTeams.length > 0}
                    <div class="opacity-60">
                        <h3 class="text-lg font-sans font-bold text-primary-700 uppercase tracking-wide mb-2 px-2">
                            Out of Playoffs
                        </h3>
                        
                        <div class="grid grid-cols-15 gap-2 px-2 border-b border-primary-700/30 text-sm font-sans font-bold text-black uppercase">
                            <div class="col-span-1">Seed</div>
                            <div class="col-span-3">Team</div>
                            <div class="col-span-2 text-center">Record</div>
                            <div class="col-span-1 text-center">PCT</div>
                            <div class="col-span-2 text-center">Conf</div>
                            <div class="col-span-2 text-center">Div</div>
                            <div class="col-span-1 text-center">GB</div>
                            <div class="col-span-1 text-center">Diff</div>
                            <div class="col-span-1 text-center">PF</div>
                            <div class="col-span-1 text-center">PA</div>
                        </div>

                        <div class="space-y-1 mt-2">
                            {#each nonPlayoffTeams as seed}
                                <button class="w-full grid grid-cols-15 gap-2 px-2 py-2 rounded transition-colors cursor-pointer"
                                    on:mouseenter={(e) => handleMouseEnter(e, seed.team_primary_color)}
                                    on:mouseleave={handleMouseLeave}
                                    on:click={() => openTeamModal(seed)}
                                >
                                    <div class="col-span-1">
                                        <span class="text-sm font-heading font-bold text-primary-600">
                                            {seed.seed}
                                        </span>
                                    </div>
                                    <div class="col-span-3 flex items-center gap-2">
                                        {#if seed.logo_url}
                                            <img 
                                                src={seed.logo_url} 
                                                alt={seed.team_abbr}
                                                class="w-6 h-6 object-contain"
                                            />
                                        {/if}
                                        <div class="text-sm font-sans font-semibold text-black truncate">
                                            {seed.team_abbr}
                                        </div>
                                    </div>
                                    <div class="col-span-2 text-center">
                                        <span class="text-sm font-heading font-bold text-black">
                                            {formatRecord(seed.wins, seed.losses, seed.ties)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatWinPct(seed.win_pct)}
                                        </span>
                                    </div>
                                    <div class="col-span-2 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatRecord(seed.conference_wins, seed.conference_losses, seed.conference_ties)}
                                        </span>
                                    </div>
                                    <div class="col-span-2 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatRecord(seed.division_wins, seed.division_losses, seed.division_ties)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatGamesBack(seed.conference_games_back)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatPointDiff(seed.point_diff)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {seed.points_for}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {seed.points_against}
                                        </span>
                                    </div>
                                </button>
                            {/each}
                        </div>
                    </div>
                {/if}
            </div>
        {/if}

        <!-- Division View with Stats -->
        {#if viewMode === 'division'}
            <div class="space-y-4 min-w-[500px]">
                {#each orderedDivisions as divisionName}
                    {@const fullDivisionName = `${conference} ${divisionName}`}
                    {@const divisionTeams = standings.divisions[fullDivisionName] || []}
                    
                    <div>
                        <h3 class="text-lg font-sans font-bold text-primary-700 uppercase tracking-wide mb-2 px-2">
                            {conference} {divisionName}
                        </h3>

                        <!-- Header Row -->
                        <div class="grid grid-cols-15 gap-2 px-2 border-b border-primary-700/30 text-sm font-sans font-bold text-black uppercase">
                            <div class="col-span-1">Seed</div>
                            <div class="col-span-3">Team</div>
                            <div class="col-span-2 text-center">Record</div>
                            <div class="col-span-1 text-center">PCT</div>
                            <div class="col-span-2 text-center">Conf</div>
                            <div class="col-span-2 text-center">Div</div>
                            <div class="col-span-1 text-center">GB</div>
                            <div class="col-span-1 text-center">Diff</div>
                            <div class="col-span-1 text-center">PF</div>
                            <div class="col-span-1 text-center">PA</div>
                        </div>

                        <!-- Data Rows -->
                        <div class="space-y-1 mt-2">
                            {#each divisionTeams as team}
                                {@const teamSeed = getTeamSeed(team.team_id)}
                                {@const isPlayoffTeam = teamSeed && teamSeed <= 7}
                                <button class="w-full grid grid-cols-15 gap-2 px-2 py-2 rounded transition-colors cursor-pointer"
                                    class:opacity-60={!isPlayoffTeam}
                                    on:mouseenter={(e) => handleMouseEnter(e, team.team_primary_color)}
                                    on:mouseleave={handleMouseLeave}
                                    on:click={() => openTeamModal(convertToPlayoffSeed(team))}
                                >
                                    <div class="col-span-1">
                                        <span class="text-sm font-heading font-bold text-primary-600">
                                            {teamSeed ?? '—'}
                                        </span>
                                    </div>
                                    <div class="col-span-3 flex items-center gap-2">
                                        {#if team.logo_url}
                                            <img 
                                                src={team.logo_url} 
                                                alt={team.team_abbr}
                                                class="w-6 h-6 object-contain"
                                            />
                                        {/if}
                                        <div class="text-sm font-sans font-semibold text-black truncate">
                                            {team.team_abbr}
                                        </div>
                                    </div>
                                    <div class="col-span-2 text-center">
                                        <span class="text-sm font-heading font-bold text-black">
                                            {formatRecord(team.wins, team.losses, team.ties)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatWinPct(team.win_pct)}
                                        </span>
                                    </div>
                                    <div class="col-span-2 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatRecord(team.conference_wins, team.conference_losses, team.conference_ties)}
                                        </span>
                                    </div>
                                    <div class="col-span-2 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatRecord(team.division_wins, team.division_losses, team.division_ties)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatGamesBack(team.division_games_back)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatPointDiff(team.point_diff)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {team.points_for}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {team.points_against}
                                        </span>
                                    </div>
                                </button>
                            {/each}
                        </div>
                    </div>
                {/each}
            </div>
        {/if}
    </div>
</div>