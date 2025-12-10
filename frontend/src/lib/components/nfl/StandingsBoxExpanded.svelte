<script lang="ts">
    import type { ConferenceStandings } from '$types';

    export let standings: ConferenceStandings;
    export let conference: 'AFC' | 'NFC';

    type ViewMode = 'conference' | 'division';
    let viewMode: ViewMode = 'conference';

    $: divisionWinners = standings.playoff_seeds.slice(0, 4);
    $: wildCardTeams = standings.playoff_seeds.slice(4, 7);
    $: nonPlayoffTeams = standings.playoff_seeds.slice(7);

    $: orderedDivisions = ['North', 'South', 'East', 'West'].filter(div => 
        standings.divisions[`${conference} ${div}`]
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
</script>

<div class="bg-neutral border-2 border-primary-700 rounded-lg p-4 w-full">
    <!-- Header -->
    <div class="flex items-center justify-between mb-4">
        <h2 class="text-2xl font-heading font-bold uppercase tracking-wide"
            style="color: {conference === 'AFC' ? '#C8102E' : '#013369'}">
            {conference}
        </h2>

        <!-- View Toggle -->
        <div class="flex bg-primary-900 rounded-lg p-1 gap-1">
            <button
                on:click={() => viewMode = 'conference'}
                class="px-4 py-2 text-sm font-sans font-semibold rounded transition-colors"
                class:bg-primary-600={viewMode === 'conference'}
                class:text-neutral={viewMode === 'conference'}
                class:hover:text-neutral={viewMode !== 'conference'}
            >
                Conference
            </button>
            <button
                on:click={() => viewMode = 'division'}
                class="px-4 py-2 text-sm font-sans font-semibold rounded transition-colors"
                class:bg-primary-600={viewMode === 'division'}
                class:text-neutral={viewMode === 'division'}
                class:hover:text-neutral={viewMode !== 'division'}
            >
                Division
            </button>
        </div>
    </div>

    <!-- Conference View with Stats -->
    {#if viewMode === 'conference'}
        <div class="space-y-4">
            <!-- Division Winners (Seeds 1-4) -->
            <div>
                <h3 class="text-xs font-sans font-bold text-neutral/60 uppercase tracking-wide mb-3 px-2">
                    Division Winners
                </h3>
                
                <!-- Header Row -->
                <div class="grid grid-cols-12 gap-2 px-2 pb-2 border-b border-primary-700/30 text-xs font-sans font-bold text-neutral/60 uppercase">
                    <div class="col-span-1">Seed</div>
                    <div class="col-span-3">Team</div>
                    <div class="col-span-2 text-center">Record</div>
                    <div class="col-span-1 text-center">PCT</div>
                    <div class="col-span-2 text-center">Conf</div>
                    <div class="col-span-2 text-center">Div</div>
                    <div class="col-span-1 text-center">GB</div>
                </div>

                <!-- Data Rows -->
                <div class="space-y-1 mt-2">
                    {#each divisionWinners as seed}
                        {@const teamData = standings.playoff_seeds.find(s => s.seed === seed.seed)?.team}
                        <div class="grid grid-cols-12 gap-2 px-2 py-2 rounded hover:bg-primary-900/30 transition-colors">
                            <div class="col-span-1">
                                <span class="text-sm font-heading font-bold text-primary-400">
                                    {seed.seed}
                                </span>
                            </div>
                            <div class="col-span-3">
                                <div class="text-sm font-sans font-semibold text-neutral">
                                    {seed.team_abbr}
                                </div>
                            </div>
                            <div class="col-span-2 text-center">
                                <span class="text-sm font-heading font-bold text-neutral">
                                    {formatRecord(seed.wins, seed.losses, seed.ties)}
                                </span>
                            </div>
                            <div class="col-span-1 text-center">
                                <span class="text-sm font-sans text-neutral">
                                    {teamData ? formatWinPct(teamData.win_pct) : '—'}
                                </span>
                            </div>
                            <div class="col-span-2 text-center">
                                <span class="text-sm font-sans text-neutral">
                                    {teamData?.conference_record || '—'}
                                </span>
                            </div>
                            <div class="col-span-2 text-center">
                                <span class="text-sm font-sans text-neutral">
                                    {teamData?.division_record || '—'}
                                </span>
                            </div>
                            <div class="col-span-1 text-center">
                                <span class="text-sm font-sans text-neutral">
                                    {teamData ? formatGamesBack(teamData.conference_games_back) : '—'}
                                </span>
                            </div>
                        </div>
                    {/each}
                </div>
            </div>

            <!-- Wild Card -->
            {#if wildCardTeams.length > 0}
                <div class="border-t border-primary-700/50 pt-4">
                    <h3 class="text-xs font-sans font-bold text-neutral/60 uppercase tracking-wide mb-3 px-2">
                        Wild Card
                    </h3>
                    
                    <!-- Header Row -->
                    <div class="grid grid-cols-12 gap-2 px-2 pb-2 border-b border-primary-700/30 text-xs font-sans font-bold text-neutral/60 uppercase">
                        <div class="col-span-1">Seed</div>
                        <div class="col-span-3">Team</div>
                        <div class="col-span-2 text-center">Record</div>
                        <div class="col-span-1 text-center">PCT</div>
                        <div class="col-span-2 text-center">Conf</div>
                        <div class="col-span-2 text-center">Div</div>
                        <div class="col-span-1 text-center">GB</div>
                    </div>

                    <div class="space-y-1 mt-2">
                        {#each wildCardTeams as seed}
                            {@const teamData = standings.playoff_seeds.find(s => s.seed === seed.seed)?.team}
                            <div class="grid grid-cols-12 gap-2 px-2 py-2 rounded hover:bg-primary-900/30 transition-colors">
                                <div class="col-span-1">
                                    <span class="text-sm font-heading font-bold text-primary-400">
                                        {seed.seed}
                                    </span>
                                </div>
                                <div class="col-span-3">
                                    <div class="text-sm font-sans font-semibold text-neutral">
                                        {seed.team_abbr}
                                    </div>
                                </div>
                                <div class="col-span-2 text-center">
                                    <span class="text-sm font-heading font-bold text-neutral">
                                        {formatRecord(seed.wins, seed.losses, seed.ties)}
                                    </span>
                                </div>
                                <div class="col-span-1 text-center">
                                    <span class="text-sm font-sans text-neutral">
                                        {teamData ? formatWinPct(teamData.win_pct) : '—'}
                                    </span>
                                </div>
                                <div class="col-span-2 text-center">
                                    <span class="text-sm font-sans text-neutral">
                                        {teamData?.conference_record || '—'}
                                    </span>
                                </div>
                                <div class="col-span-2 text-center">
                                    <span class="text-sm font-sans text-neutral">
                                        {teamData?.division_record || '—'}
                                    </span>
                                </div>
                                <div class="col-span-1 text-center">
                                    <span class="text-sm font-sans text-neutral">
                                        {teamData ? formatGamesBack(teamData.conference_games_back) : '—'}
                                    </span>
                                </div>
                            </div>
                        {/each}
                    </div>
                </div>
            {/if}

            <!-- Non-Playoff Teams -->
            {#if nonPlayoffTeams.length > 0}
                <div class="border-t border-primary-700/50 pt-4 opacity-60">
                    <h3 class="text-xs font-sans font-bold text-neutral/60 uppercase tracking-wide mb-3 px-2">
                        Out of Playoffs
                    </h3>
                    
                    <div class="grid grid-cols-12 gap-2 px-2 pb-2 border-b border-primary-700/30 text-xs font-sans font-bold text-neutral/60 uppercase">
                        <div class="col-span-1">Seed</div>
                        <div class="col-span-3">Team</div>
                        <div class="col-span-2 text-center">Record</div>
                        <div class="col-span-1 text-center">PCT</div>
                        <div class="col-span-2 text-center">Conf</div>
                        <div class="col-span-2 text-center">Div</div>
                        <div class="col-span-1 text-center">GB</div>
                    </div>

                    <div class="space-y-1 mt-2">
                        {#each nonPlayoffTeams as seed}
                            {@const teamData = standings.playoff_seeds.find(s => s.seed === seed.seed)?.team}
                            <div class="grid grid-cols-12 gap-2 px-2 py-2 rounded hover:bg-primary-900/30 transition-colors">
                                <div class="col-span-1">
                                    <span class="text-sm font-heading font-bold text-primary-400">
                                        {seed.seed}
                                    </span>
                                </div>
                                <div class="col-span-3">
                                    <div class="text-sm font-sans font-semibold text-neutral">
                                        {seed.team_abbr}
                                    </div>
                                </div>
                                <div class="col-span-2 text-center">
                                    <span class="text-sm font-heading font-bold text-neutral">
                                        {formatRecord(seed.wins, seed.losses, seed.ties)}
                                    </span>
                                </div>
                                <div class="col-span-1 text-center">
                                    <span class="text-sm font-sans text-neutral">
                                        {teamData ? formatWinPct(teamData.win_pct) : '—'}
                                    </span>
                                </div>
                                <div class="col-span-2 text-center">
                                    <span class="text-sm font-sans text-neutral">
                                        {teamData?.conference_record || '—'}
                                    </span>
                                </div>
                                <div class="col-span-2 text-center">
                                    <span class="text-sm font-sans text-neutral">
                                        {teamData?.division_record || '—'}
                                    </span>
                                </div>
                                <div class="col-span-1 text-center">
                                    <span class="text-sm font-sans text-neutral">
                                        {teamData ? formatGamesBack(teamData.conference_games_back) : '—'}
                                    </span>
                                </div>
                            </div>
                        {/each}
                    </div>
                </div>
            {/if}
        </div>
    {/if}

    <!-- Division View with Stats -->
    {#if viewMode === 'division'}
        <div class="space-y-4">
            {#each orderedDivisions as divisionName}
                {@const fullDivisionName = `${conference} ${divisionName}`}
                {@const divisionTeams = standings.divisions[fullDivisionName] || []}
                
                <div>
                    <h3 class="text-xs font-sans font-bold text-neutral/60 uppercase tracking-wide mb-3 px-2">
                        {conference} {divisionName}
                    </h3>

                    <!-- Header Row -->
                    <div class="grid grid-cols-12 gap-2 px-2 pb-2 border-b border-primary-700/30 text-xs font-sans font-bold text-neutral/60 uppercase">
                        <div class="col-span-1">Rk</div>
                        <div class="col-span-3">Team</div>
                        <div class="col-span-2 text-center">Record</div>
                        <div class="col-span-1 text-center">PCT</div>
                        <div class="col-span-2 text-center">Conf</div>
                        <div class="col-span-2 text-center">Div</div>
                        <div class="col-span-1 text-center">GB</div>
                    </div>

                    <!-- Data Rows -->
                    <div class="space-y-1 mt-2">
                        {#each divisionTeams as team, index}
                            <div class="grid grid-cols-12 gap-2 px-2 py-2 rounded hover:bg-primary-900/30 transition-colors"
                                 class:opacity-60={index > 0}>
                                <div class="col-span-1">
                                    <span class="text-sm font-heading font-bold text-primary-400">
                                        {team.rank}
                                    </span>
                                </div>
                                <div class="col-span-3">
                                    <div class="text-sm font-sans font-semibold text-neutral">
                                        {team.team_abbr}
                                    </div>
                                </div>
                                <div class="col-span-2 text-center">
                                    <span class="text-sm font-heading font-bold text-neutral">
                                        {formatRecord(team.wins, team.losses, team.ties)}
                                    </span>
                                </div>
                                <div class="col-span-1 text-center">
                                    <span class="text-sm font-sans text-neutral">
                                        {formatWinPct(team.win_pct)}
                                    </span>
                                </div>
                                <div class="col-span-2 text-center">
                                    <span class="text-sm font-sans text-neutral">
                                        {team.conference_record}
                                    </span>
                                </div>
                                <div class="col-span-2 text-center">
                                    <span class="text-sm font-sans text-neutral">
                                        {team.division_record}
                                    </span>
                                </div>
                                <div class="col-span-1 text-center">
                                    <span class="text-sm font-sans text-neutral">
                                        {formatGamesBack(team.division_games_back)}
                                    </span>
                                </div>
                            </div>
                        {/each}
                    </div>
                </div>
            {/each}
        </div>
    {/if}
</div>