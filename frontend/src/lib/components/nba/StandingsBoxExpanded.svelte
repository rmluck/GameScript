<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import type { NBAConferenceStandings, NBAPlayoffSeed } from '$types';

    // Props
    export let standings: NBAConferenceStandings;
    export let conference: 'East' | 'West';
    type ViewMode = 'conference' | 'division';
    export let viewMode: ViewMode = 'conference';

    // Event dispatcher
    const dispatch = createEventDispatcher();

    // Sort teams into groups
    $: playoffTeams = standings.playoff_seeds.slice(0, 6);
    $: playInTeams = standings.playoff_seeds.slice(6, 10);
    $: nonPlayoffTeams = standings.playoff_seeds.slice(10);

    $: orderedDivisions = conference === "East"
        ? ['Atlantic', 'Central', 'Southeast'].filter(div => standings.divisions[`${div}`])
        : ['Northwest', 'Pacific', 'Southwest'].filter(div => standings.divisions[`${div}`]);

    // Map for quick seed lookup
    $: teamSeedMap = new Map(
        standings.playoff_seeds.map(seed => [seed.team_id, seed.seed])
    );

    function formatRecord(wins: number, losses: number): string {
        return `${wins}-${losses}`;
    }

    function formatWinPct(winPct: number): string {
        if (winPct === -1.0) return '.000'; // Handle 0-0 case
        return winPct.toFixed(3);
    }

    function formatGamesBack(gb: number): string {
        if (gb === 0) return '—';
        return gb.toFixed(1);
    }

    function formatPoints(points: number, games_with_scores: number): string {
        const points_average = points / (games_with_scores);
        return points_average.toFixed(1);
    }

    function formatPointDiff(pointsFor: number, pointsAgainst: number, gamesWithScores: number): string {
        const diff_average = (pointsFor - pointsAgainst) / gamesWithScores;
        if (diff_average > 0) return `+${diff_average.toFixed(1)}`;
        return diff_average.toFixed(1);
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

    function openTeamModal(team: NBAPlayoffSeed) {
        dispatch('openTeamModal', { team });
    }

    function convertToPlayoffSeed(team: any): NBAPlayoffSeed {
        return {
            seed: team.seed,
            team_id: team.team_id,
            team_name: team.team_name,
            team_city: team.team_city,
            team_abbr: team.team_abbr,
            wins: team.wins,
            losses: team.losses,
            win_pct: team.win_pct,
            home_wins: team.home_wins,
            home_losses: team.home_losses,
            away_wins: team.away_wins,
            away_losses: team.away_losses,
            division_wins: team.division_wins,
            division_losses: team.division_losses,
            conference_wins: team.conference_wins,
            conference_losses: team.conference_losses,
            division_games_back: team.division_games_back,
            conference_games_back: team.conference_games_back,
            points_for: team.points_for,
            points_against: team.points_against,
            games_with_scores: team.games_with_scores,
            strength_of_schedule: team.strength_of_schedule,
            strength_of_victory: team.strength_of_victory,
            is_division_winner: team.is_division_winner,
            logo_url: team.logo_url,
            team_primary_color: team.team_primary_color,
            team_secondary_color: team.team_secondary_color
        };
    }
</script>

<div class="bg-neutral border-2 border-primary-700 rounded-lg px-4 py-6 w-full">
    <!-- Header -->
    <div class="flex items-center justify-between pb-6 mb-4 border-b-2 border-primary-700">
        <h2 class="text-2xl font-heading font-bold uppercase tracking-wide"
            style="color: {conference === 'West' ? '#C8102E' : '#013369'}">
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
            <div class="space-y-4 min-w-[800px]">
                <!-- Playoff Teams (Seeds 1-6) -->
                <div>
                    <h3 class="text-lg font-sans font-bold text-primary-700 uppercase tracking-wide mb-2 px-2">
                        Playoff Teams
                    </h3>
                    
                    <!-- Header Row -->
                    <div class="grid grid-cols-15 gap-1 px-2 border-b border-primary-700/30 text-xs font-sans font-bold text-black uppercase">
                        <div class="col-span-1" title="Playoff Seed">Seed</div>
                        <div class="col-span-2" title="Team">Team</div>
                        <div class="col-span-1 text-center" title="Overall Record">Record</div>
                        <div class="col-span-1 text-center" title="Win Percentage">PCT</div>
                        <div class="col-span-1 text-center" title="Conference Record">Conf</div>
                        <div class="col-span-1 text-center" title="Division Record">Div</div>
                        <div class="col-span-1 text-center" title="Home Record">Home</div>
                        <div class="col-span-1 text-center" title="Away Record">Away</div>
                        <div class="col-span-1 text-center" title="Games Back">GB</div>
                        <div class="col-span-1 text-center" title="Point Differential">Diff</div>
                        <div class="col-span-1 text-center" title="Points Per Game">PPG</div>
                        <div class="col-span-1 text-center" title="Opponent Points Per Game">OPP PPG</div>
                        <div class="col-span-1 text-center" title="Strength of Schedule">SOS</div>
                        <div class="col-span-1 text-center" title="Strength of Victory">SOV</div>
                    </div>

                    <!-- Data Rows -->
                    <div class="space-y-1 mt-2">
                        {#each playoffTeams as seed}
                            <button class="w-full grid grid-cols-15 gap-1 px-2 py-2 rounded transition-colors cursor-pointer"
                                on:mouseenter={(e) => handleMouseEnter(e, seed.team_primary_color)}
                                on:mouseleave={handleMouseLeave}
                                on:click={() => openTeamModal(seed)}
                            >
                                <div class="col-span-1">
                                    <span class="text-sm font-heading font-bold text-primary-600">
                                        {seed.seed}
                                    </span>
                                </div>
                                <div class="col-span-2 flex items-center gap-2">
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
                                <div class="col-span-1 text-center">
                                    <span class="text-sm font-heading font-bold text-black">
                                        {formatRecord(seed.wins, seed.losses)}
                                    </span>
                                </div>
                                <div class="col-span-1 text-center">
                                    <span class="text-sm font-sans text-black">
                                        {formatWinPct(seed.win_pct)}
                                    </span>
                                </div>
                                <div class="col-span-1 text-center">
                                    <span class="text-xs font-sans text-black">
                                        {formatRecord(seed.conference_wins, seed.conference_losses)}
                                    </span>
                                </div>
                                <div class="col-span-1 text-center">
                                    <span class="text-xs font-sans text-black">
                                        {formatRecord(seed.division_wins, seed.division_losses)}
                                    </span>
                                </div>
                                <div class="col-span-1 text-center">
                                    <span class="text-xs font-sans text-black">
                                        {formatRecord(seed.home_wins, seed.home_losses)}
                                    </span>
                                </div>
                                <div class="col-span-1 text-center">
                                    <span class="text-xs font-sans text-black">
                                        {formatRecord(seed.away_wins, seed.away_losses)}
                                    </span>
                                </div>
                                <div class="col-span-1 text-center">
                                    <span class="text-sm font-sans text-black">
                                        {formatGamesBack(seed.conference_games_back)}
                                    </span>
                                </div>
                                <div class="col-span-1 text-center">
                                    <span class="text-sm font-sans text-black">
                                        {formatPointDiff(seed.points_for, seed.points_against, seed.games_with_scores)}
                                    </span>
                                </div>
                                <div class="col-span-1 text-center">
                                    <span class="text-sm font-sans text-black">
                                        {formatPoints(seed.points_for, seed.games_with_scores)}
                                    </span>
                                </div>
                                <div class="col-span-1 text-center">
                                    <span class="text-sm font-sans text-black">
                                        {formatPoints(seed.points_against, seed.games_with_scores)}
                                    </span>
                                </div>
                                <div class="col-span-1 text-center">
                                    <span class="text-sm font-sans text-black">
                                        {formatWinPct(seed.strength_of_schedule)}
                                    </span>
                                </div>
                                <div class="col-span-1 text-center">
                                    <span class="text-sm font-sans text-black">
                                        {formatWinPct(seed.strength_of_victory)}
                                    </span>
                                </div>
                            </button>
                        {/each}
                    </div>
                </div>

                <!-- Play-In Teams (Seeds 7-10) -->
                {#if playInTeams.length > 0}
                    <div>
                        <h3 class="text-lg font-sans font-bold text-primary-700 uppercase tracking-wide mb-2 px-2">
                            Play-In Teams
                        </h3>
                        
                        <!-- Header Row -->
                        <div class="grid grid-cols-15 gap-1 px-2 border-b border-primary-700/30 text-xs font-sans font-bold text-black uppercase">
                            <div class="col-span-1" title="Playoff Seed">Seed</div>
                            <div class="col-span-2" title="Team">Team</div>
                            <div class="col-span-1 text-center" title="Overall Record">Record</div>
                            <div class="col-span-1 text-center" title="Win Percentage">PCT</div>
                            <div class="col-span-1 text-center" title="Conference Record">Conf</div>
                            <div class="col-span-1 text-center" title="Division Record">Div</div>
                            <div class="col-span-1 text-center" title="Home Record">Home</div>
                            <div class="col-span-1 text-center" title="Away Record">Away</div>
                            <div class="col-span-1 text-center" title="Games Back">GB</div>
                            <div class="col-span-1 text-center" title="Point Differential">Diff</div>
                            <div class="col-span-1 text-center" title="Points Per Game">PPG</div>
                            <div class="col-span-1 text-center" title="Opponent Points Per Game">OPP PPG</div>
                            <div class="col-span-1 text-center" title="Strength of Schedule">SOS</div>
                            <div class="col-span-1 text-center" title="Strength of Victory">SOV</div>
                        </div>

                        <div class="space-y-1 mt-2">
                            {#each playInTeams as seed}
                                <button class="w-full grid grid-cols-15 gap-1 px-2 py-2 rounded transition-colors cursor-pointer"
                                    on:mouseenter={(e) => handleMouseEnter(e, seed.team_primary_color)}
                                    on:mouseleave={handleMouseLeave}
                                    on:click={() => openTeamModal(seed)}
                                >
                                    <div class="col-span-1">
                                        <span class="text-sm font-heading font-bold text-primary-600">
                                            {seed.seed}
                                        </span>
                                    </div>
                                    <div class="col-span-2 flex items-center gap-2">
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
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-heading font-bold text-black">
                                            {formatRecord(seed.wins, seed.losses)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatWinPct(seed.win_pct)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-xs font-sans text-black">
                                            {formatRecord(seed.conference_wins, seed.conference_losses)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-xs font-sans text-black">
                                            {formatRecord(seed.division_wins, seed.division_losses)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-xs font-sans text-black">
                                            {formatRecord(seed.home_wins, seed.home_losses)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-xs font-sans text-black">
                                            {formatRecord(seed.away_wins, seed.away_losses)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatGamesBack(seed.conference_games_back)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatPointDiff(seed.points_for, seed.points_against, seed.games_with_scores)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatPoints(seed.points_for, seed.games_with_scores)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatPoints(seed.points_against, seed.games_with_scores)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatWinPct(seed.strength_of_schedule)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatWinPct(seed.strength_of_victory)}
                                        </span>
                                    </div>
                                </button>
                            {/each}
                        </div>
                    </div>
                {/if}

                <!-- Non-Playoff Teams (Seeds 11-15) -->
                {#if nonPlayoffTeams.length > 0}
                    <div class="opacity-60">
                        <h3 class="text-lg font-sans font-bold text-primary-700 uppercase tracking-wide mb-2 px-2">
                            Out of Playoffs
                        </h3>
                        
                        <div class="grid grid-cols-15 gap-1 px-2 border-b border-primary-700/30 text-xs font-sans font-bold text-black uppercase">
                            <div class="col-span-1" title="Playoff Seed">Seed</div>
                            <div class="col-span-2" title="Team">Team</div>
                            <div class="col-span-1 text-center" title="Overall Record">Record</div>
                            <div class="col-span-1 text-center" title="Win Percentage">PCT</div>
                            <div class="col-span-1 text-center" title="Conference Record">Conf</div>
                            <div class="col-span-1 text-center" title="Division Record">Div</div>
                            <div class="col-span-1 text-center" title="Home Record">Home</div>
                            <div class="col-span-1 text-center" title="Away Record">Away</div>
                            <div class="col-span-1 text-center" title="Games Back">GB</div>
                            <div class="col-span-1 text-center" title="Point Differential">Diff</div>
                            <div class="col-span-1 text-center" title="Points Per Game">PPG</div>
                            <div class="col-span-1 text-center" title="Opponent Points Per Game">OPP PPG</div>
                            <div class="col-span-1 text-center" title="Strength of Schedule">SOS</div>
                            <div class="col-span-1 text-center" title="Strength of Victory">SOV</div>
                        </div>

                        <div class="space-y-1 mt-2">
                            {#each nonPlayoffTeams as seed}
                                <button class="w-full grid grid-cols-15 gap-1 px-2 py-2 rounded transition-colors cursor-pointer"
                                    on:mouseenter={(e) => handleMouseEnter(e, seed.team_primary_color)}
                                    on:mouseleave={handleMouseLeave}
                                    on:click={() => openTeamModal(seed)}
                                >
                                    <div class="col-span-1">
                                        <span class="text-sm font-heading font-bold text-primary-600">
                                            {seed.seed}
                                        </span>
                                    </div>
                                    <div class="col-span-2 flex items-center gap-2">
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
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-heading font-bold text-black">
                                            {formatRecord(seed.wins, seed.losses)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatWinPct(seed.win_pct)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-xs font-sans text-black">
                                            {formatRecord(seed.conference_wins, seed.conference_losses)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-xs font-sans text-black">
                                            {formatRecord(seed.division_wins, seed.division_losses)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-xs font-sans text-black">
                                            {formatRecord(seed.home_wins, seed.home_losses)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-xs font-sans text-black">
                                            {formatRecord(seed.away_wins, seed.away_losses)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatGamesBack(seed.conference_games_back)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatPointDiff(seed.points_for, seed.points_against, seed.games_with_scores)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatPoints(seed.points_for, seed.games_with_scores)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatPoints(seed.points_against, seed.games_with_scores)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatWinPct(seed.strength_of_schedule)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatWinPct(seed.strength_of_victory)}
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
            <div class="space-y-4 min-w-[800px]">
                {#each orderedDivisions as divisionName}
                    {@const divisionTeams = standings.divisions[divisionName] || []}
                    <div>
                        <h3 class="text-lg font-sans font-bold text-primary-700 uppercase tracking-wide mb-2 px-2">
                            {divisionName}
                        </h3>

                        <div class="grid grid-cols-15 gap-1 px-2 border-b border-primary-700/30 text-xs font-sans font-bold text-black uppercase">
                            <div class="col-span-1" title="Playoff Seed">Seed</div>
                            <div class="col-span-2" title="Team">Team</div>
                            <div class="col-span-1 text-center" title="Overall Record">Record</div>
                            <div class="col-span-1 text-center" title="Win Percentage">PCT</div>
                            <div class="col-span-1 text-center" title="Conference Record">Conf</div>
                            <div class="col-span-1 text-center" title="Division Record">Div</div>
                            <div class="col-span-1 text-center" title="Home Record">Home</div>
                            <div class="col-span-1 text-center" title="Away Record">Away</div>
                            <div class="col-span-1 text-center" title="Games Back">GB</div>
                            <div class="col-span-1 text-center" title="Point Differential">Diff</div>
                            <div class="col-span-1 text-center" title="Points Per Game">PPG</div>
                            <div class="col-span-1 text-center" title="Opponent Points Per Game">OPP PPG</div>
                            <div class="col-span-1 text-center" title="Strength of Schedule">SOS</div>
                            <div class="col-span-1 text-center" title="Strength of Victory">SOV</div>
                        </div>

                        <div class="space-y-1 mt-2">
                            {#each divisionTeams as team}
                                {@const teamSeed = getTeamSeed(team.team_id)}
                                {@const isPlayoffTeam = teamSeed && teamSeed <= 7}
                                <button class="w-full grid grid-cols-15 gap-1 px-2 py-2 rounded transition-colors cursor-pointer"
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
                                    <div class="col-span-2 flex items-center gap-2">
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
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-heading font-bold text-black">
                                            {formatRecord(team.wins, team.losses)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatWinPct(team.win_pct)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-xs font-sans text-black">
                                            {formatRecord(team.conference_wins, team.conference_losses)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-xs font-sans text-black">
                                            {formatRecord(team.division_wins, team.division_losses)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-xs font-sans text-black">
                                            {formatRecord(team.home_wins, team.home_losses)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-xs font-sans text-black">
                                            {formatRecord(team.away_wins, team.away_losses)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatGamesBack(team.division_games_back)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatPointDiff(team.points_for, team.points_against, team.games_with_scores)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatPoints(team.points_for, team.games_with_scores)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatPoints(team.points_against, team.games_with_scores)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatWinPct(team.strength_of_schedule)}
                                        </span>
                                    </div>
                                    <div class="col-span-1 text-center">
                                        <span class="text-sm font-sans text-black">
                                            {formatWinPct(team.strength_of_victory)}
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