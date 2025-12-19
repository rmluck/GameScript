<script lang="ts">
    import { createEventDispatcher, onMount } from 'svelte';
    import { gamesAPI } from '$lib/api/games';
    import { picksAPI } from '$lib/api/picks';
    import { standingsAPI } from '$lib/api/standings';
    import { playoffsAPI } from '$lib/api/playoffs';
    import GamePickerRow from './GamePickerRow.svelte';
    import type { PlayoffSeed, Game, Pick, PlayoffState } from '$types';

    export let team: PlayoffSeed;
    export let scenarioId: number;
    export let seasonId: number;

    const dispatch = createEventDispatcher();

    let teamGames: Game[] = [];
    let picks: Map<number, Pick> = new Map();
    let loading = true;
    let error = '';

    let byeWeek: number | null = null;
    let gamesWithBye: Array<Game | { isByeWeek: true; week: number }> = [];

    // Add this to store updated team data
    let currentTeam: PlayoffSeed = team;
    let teamDivision: string = '';

    let playoffState: PlayoffState | null = null;

    onMount(async () => {
        await loadTeamGames();
        await loadPicks();
        await loadPlayoffState();
    });

    async function loadTeamGames() {
        try {
            loading = true;
            const allGames = await gamesAPI.getBySeason(seasonId);
            teamGames = allGames.filter(
                game => game.home_team_id === team.team_id || game.away_team_id === team.team_id
            );

            await loadTeamDivision();

            // Find bye week
            const playedWeeks = new Set(teamGames.map(game => game.week));
            for (let week = 1; week <= 18; week++) {
                if (!playedWeeks.has(week)) {
                    byeWeek = week;
                    break;
                }
            }

            // Insert bye week into the schedule
            if (byeWeek !== null) {
                gamesWithBye = [
                    ...teamGames.filter(game => game.week < byeWeek!),
                    { isByeWeek: true, week: byeWeek },
                    ...teamGames.filter(game => game.week > byeWeek!)
                ];
            } else {
                gamesWithBye = teamGames;
            }
        } catch (err: any) {
            error = 'Failed to load team games';
            console.error(err);
        } finally {
            loading = false;
        }
    }

    async function loadPicks() {
        try {
            const allPicks = await picksAPI.getByScenario(scenarioId);
            picks = new Map(allPicks.map(pick => [pick.game_id, pick]));
        } catch (err: any) {
            console.error('Error loading picks:', err);
        }
    }

    async function loadTeamDivision() {
        try {
            // Get team's division from the first game
            if (teamGames.length > 0) {
                const firstGame = teamGames[0];
                const teamData = firstGame.home_team_id === team.team_id 
                    ? firstGame.home_team 
                    : firstGame.away_team;
                teamDivision = teamData.division;
            }
        } catch (err: any) {
            console.error('Error loading team division:', err);
        }
    }

    async function loadPlayoffState() {
        try {
            const response = await playoffsAPI.getState(scenarioId);
            playoffState = response.playoff_state;
        } catch (err: any) {
            console.error('Failed to load playoff state:', err);
        }
    }

    // Add function to reload team data
    async function reloadTeamData() {
        try {
            const standings = await standingsAPI.getByScenario(scenarioId);
            
            // Find updated team data in standings
            const allSeeds = [
                ...standings.afc.playoff_seeds,
                ...standings.nfc.playoff_seeds
            ];
            
            const updatedTeam = allSeeds.find(seed => seed.team_id === team.team_id);
            
            if (updatedTeam) {
                currentTeam = updatedTeam;
            }
        } catch (err: any) {
            console.error('Error reloading team data:', err);
        }
    }

    function closeModal() {
        dispatch('close');
    }

    function handleClickOutside(event: MouseEvent) {
        if (event.target === event.currentTarget) {
            closeModal();
        }
    }

    async function handlePickChange(event: CustomEvent) {
        try {
            const { gameId, pickedTeamId, predictedHomeScore, predictedAwayScore, deletePick } = event.detail;
            const existingPick = picks.get(gameId);

            if (deletePick && existingPick) {
                await picksAPI.delete(scenarioId, gameId);
                picks.delete(gameId);
            } else if (existingPick) {
                const updated = await picksAPI.update(scenarioId, gameId, {
                    picked_team_id: pickedTeamId === undefined ? null : pickedTeamId,
                    predicted_home_score: predictedHomeScore,
                    predicted_away_score: predictedAwayScore
                });
                picks.set(gameId, updated);
            } else {
                const updated = await picksAPI.create(scenarioId, gameId, {
                    picked_team_id: pickedTeamId === undefined ? null : pickedTeamId,
                    predicted_home_score: predictedHomeScore,
                    predicted_away_score: predictedAwayScore
                });
                picks.set(gameId, updated);
            }
            
            picks = new Map(picks);
            
            // Reload team data to get updated stats
            await reloadTeamData();
            
            dispatch('pickUpdated');
        } catch (err: any) {
            console.error('Error saving pick:', err);
            alert('Failed to save pick. Please try again.');
        }
    }

    function formatRecord(wins: number, losses: number, ties: number): string {
        return ties > 0 ? `${wins}-${losses}-${ties}` : `${wins}-${losses}`;
    }

    function formatWinPct(winPct: number): string {
        return winPct.toFixed(3);
    }

    function formatGamesBack(gb: number): string {
        if (gb === 0) return '-';
        return gb.toFixed(1);
    }

    function formatPointDiff(diff: number): string {
        if (diff > 0) return `+${diff}`;
        return diff.toString();
    }

    function formatDate(dateString: string): string {
        const date = new Date(dateString);
        return date.toLocaleDateString('en-US', { 
            weekday: 'short', 
            month: 'short', 
            day: 'numeric' 
        });
    }

    function formatTime(dateString: string): string {
        const date = new Date(dateString);
        return date.toLocaleTimeString('en-US', { 
            hour: 'numeric', 
            minute: '2-digit',
        });
    }

    function isHomeGame(game: Game): boolean {
        return game.home_team_id === team.team_id;
    }

    function getOpponent(game: Game) {
        return isHomeGame(game) ? game.away_team : game.home_team;
    }

    function getGameResult(game: Game): 'win' | 'loss' | 'tie' | 'upcoming' {
        const pick = picks.get(game.id);
        
        if (pick?.picked_team_id !== undefined && pick.picked_team_id !== null) {
            if (pick.picked_team_id === 0) return 'tie';
            return pick.picked_team_id === team.team_id ? 'win' : 'loss';
        }
        
        if (game.status === 'final' && game.home_score !== null && game.away_score !== null) {
            if (game.home_score === game.away_score) return 'tie';
            const teamScore = isHomeGame(game) ? game.home_score : game.away_score;
            const oppScore = isHomeGame(game) ? game.away_score : game.home_score;
            return (teamScore ?? 0) > (oppScore ?? 0) ? 'win' : 'loss';
        }
        
        return 'upcoming';
    }
</script>

<!-- Modal Overlay -->
<div 
    class="fixed inset-0 bg-black/80 flex items-center justify-center z-50 p-4"
    on:click={handleClickOutside}
    on:keydown={(e) => { if (e.key === 'Escape') closeModal(); }}
    role="dialog"
    aria-modal="true"
    tabindex="-1"
>
    <section
        class="bg-neutral border-4 rounded-lg max-w-5xl w-full max-h-[90vh] overflow-hidden flex flex-col"
        style="border-color: #{currentTeam.team_primary_color};"
        role="document"
        tabindex="-1"
        aria-label="Team details modal"
    >
        <!-- Header -->
        <div
            class="border-b-4 p-6"
            style={`background-color: #${currentTeam.team_primary_color}70; border-color: #${currentTeam.team_primary_color};`}
        >
            <div class="flex items-center justify-between mb-4">
                <div class="flex items-center gap-4">
                    {#if currentTeam.logo_url}
                        <img 
                            src={currentTeam.logo_url} 
                            alt={currentTeam.team_abbr}
                            class="w-16 h-16 object-contain"
                        />
                    {/if}
                    <div>
                        <p class="text-lg font-sans text-black/70">
                            {currentTeam.team_city}
                        </p>
                        <h2 class="text-3xl font-heading font-bold text-black mb-1">
                            {currentTeam.team_name}
                        </h2>
                        {#if teamDivision}
                            <p class="text-sm font-sans text-black/60 uppercase">
                                {teamDivision}
                            </p>
                        {/if}
                    </div>
                </div>
                
                <button
                    on:click={closeModal}
                    class="p-2 rounded-lg cursor-pointer transition-colors bg-transparent"
                    style={`color: #${currentTeam.team_primary_color}E6;`}
                    on:mouseenter={(e) => e.currentTarget.style.backgroundColor = `#${currentTeam.team_primary_color}80`}
                    on:mouseleave={(e) => e.currentTarget.style.backgroundColor = `transparent`}
                    aria-label="Close modal"
                >
                    <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                </button>
            </div>

            <!-- Team Stats - 2 rows of 4 -->
            <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
                <!-- Row 1 -->
                <div class="rounded-lg p-3 bg-white/50">
                    <div class="text-xs font-sans font-semibold text-black/60 uppercase">Record</div>
                    <div class="text-xl font-heading font-bold text-black">
                        {formatRecord(currentTeam.wins, currentTeam.losses, currentTeam.ties)}
                    </div>
                </div>
                <div class="bg-white/50 rounded-lg p-3">
                    <div class="text-xs font-sans font-semibold text-black/60 uppercase">Seed</div>
                    <div class="text-xl font-heading font-bold text-black">#{currentTeam.seed}</div>
                </div>
                <div class="bg-white/50 rounded-lg p-3">
                    <div class="text-xs font-sans font-semibold text-black/60 uppercase">Conf Record</div>
                    <div class="flex flex-row">
                        <div class="text-xl font-heading font-bold text-black">
                            {formatRecord(currentTeam.conference_wins, currentTeam.conference_losses, currentTeam.conference_ties)}
                        </div>
                        <div class="text-sm font-sans text-black/60 ml-2 mt-2">
                            ({formatGamesBack(currentTeam.conference_games_back)} GB)
                        </div>
                    </div>
                </div>
                <div class="bg-white/50 rounded-lg p-3">
                    <div class="text-xs font-sans font-semibold text-black/60 uppercase">Div Record</div>
                    <div class="flex flex-row">
                        <div class="text-xl font-heading font-bold text-black">
                            {formatRecord(currentTeam.division_wins, currentTeam.division_losses, currentTeam.division_ties)}
                        </div>
                        <div class="text-sm font-sans text-black/60 ml-2 mt-2">
                            ({formatGamesBack(currentTeam.division_games_back)} GB)
                        </div>
                    </div>
                </div>

                <!-- Row 2 -->
                <div class="bg-white/50 rounded-lg p-3">
                    <div class="text-xs font-sans font-semibold text-black/60 uppercase">Win %</div>
                    <div class="text-xl font-heading font-bold text-black">
                        {formatWinPct(currentTeam.win_pct)}
                    </div>
                </div>
                <div class="bg-white/50 rounded-lg p-3">
                    <div class="text-xs font-sans font-semibold text-black/60 uppercase">Point Diff</div>
                    <div class="text-xl font-heading font-bold text-black">
                        {formatPointDiff(currentTeam.point_diff)}
                    </div>
                </div>
                <div class="bg-white/50 rounded-lg p-3">
                    <div class="text-xs font-sans font-semibold text-black/60 uppercase">Points For</div>
                    <div class="text-xl font-heading font-bold text-black">
                        {currentTeam.points_for}
                    </div>
                </div>
                <div class="bg-white/50 rounded-lg p-3">
                    <div class="text-xs font-sans font-semibold text-black/60 uppercase">Points Against</div>
                    <div class="text-xl font-heading font-bold text-black">
                        {currentTeam.points_against}
                    </div>
                </div>
            </div>
        </div>

        <!-- Schedule Content -->
        <div class="flex-1 overflow-y-auto p-6">
            <h3
                class="text-xl font-heading font-bold mb-4 uppercase tracking-wide"
                style={`color: #${currentTeam.team_primary_color};`}
            >
                Season Schedule
            </h3>

            {#if loading}
                <div class="flex items-center justify-center py-12">
                    <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-400"></div>
                </div>
            {:else if error}
                <div class="text-center py-8">
                    <p class="text-red-400 text-lg">{error}</p>
                </div>
            {:else if teamGames.length === 0}
                <div class="text-center py-8">
                    <p class="text-neutral/70 text-lg">No games found for this team</p>
                </div>
            {:else}
                <div class="space-y-2">
                    {#each gamesWithBye as item ('isByeWeek' in item ? `bye-${item.week}` : item.id)}
                        {#if 'isByeWeek' in item}
                            <!-- Bye Week Row -->
                            <div
                                class="border-2 rounded-lg transition-colors"
                                style={`border-color: #${currentTeam.team_primary_color}40; background-color: #${currentTeam.team_primary_color}10;`}
                                role="region"
                                on:mouseenter={(e) => e.currentTarget.style.borderColor = `#${currentTeam.team_primary_color}90`}
                                on:mouseleave={(e) => e.currentTarget.style.borderColor = `#${currentTeam.team_primary_color}40`}
                            >
                                <div class="grid grid-cols-3 grid-rows-3">
                                    <div class="flex items-center justify-left pl-4 pt-4">
                                        <span class="text-sm font-sans font-semibold text-black">
                                            Week {item.week}
                                        </span>
                                    </div>
                                    <div></div>
                                    <div></div>
                                    
                                    <div></div>
                                    <div class="flex items-center justify-center">
                                        <span class="text-lg font-heading font-bold uppercase tracking-wide"
                                            style={`color: #${currentTeam.team_primary_color};`}>
                                            BYE WEEK
                                        </span>
                                    </div>
                                    <div></div>
                                    
                                    <div></div>
                                    <div></div>
                                    <div></div>
                                </div>
                            </div>
                        {:else}
                            {@const game = item}
                            {@const result = getGameResult(game)}
                            {@const pick = picks.get(game.id)}
                            {@const isHome = isHomeGame(game)}
                            
                            <div
                                class="border-2 rounded-lg p-4 transition-colors"
                                style={`border-color: #${currentTeam.team_primary_color}60;`}
                                on:mouseenter={(e) => e.currentTarget.style.borderColor = `#${currentTeam.team_primary_color}90`}
                                on:mouseleave={(e) => e.currentTarget.style.borderColor = `#${currentTeam.team_primary_color}60`}
                                role="region"
                            >
                                <!-- Game Header -->
                                <div class="flex items-center justify-between mb-3">
                                    <div class="flex items-center gap-2">
                                        <span class="text-sm font-sans font-semibold text-black">
                                            Week {game.week}
                                        </span>
                                        <span class="text-xs font-sans text-black/60">
                                            {formatDate(game.start_time)} • {formatTime(game.start_time)}
                                        </span>
                                    </div>
                                    <div class="flex items-center gap-2">
                                        {#if result === 'win'}
                                            <span class="px-2 py-1 bg-green-500/20 text-green-700 text-xs font-sans font-bold rounded">W</span>
                                        {:else if result === 'loss'}
                                            <span class="px-2 py-1 bg-red-500/20 text-red-700 text-xs font-sans font-bold rounded">L</span>
                                        {:else if result === 'tie'}
                                            <span class="px-2 py-1 bg-gray-500/20 text-gray-700 text-xs font-sans font-bold rounded">T</span>
                                        {/if}
                                        <span class="text-xs font-sans text-black/60">
                                            {isHome ? 'vs' : '@'}
                                        </span>
                                    </div>
                                </div>

                                <!-- Game Picker Row (Reusable Component) -->
                                <GamePickerRow
                                    {game}
                                    {pick}
                                    compact={true}
                                    hasPlayoffs={playoffState?.is_enabled || false}
                                    on:pickChanged={handlePickChange}
                                />
                            </div>
                        {/if}
                    {/each}
                </div>
            {/if}
        </div>
    </section>
</div>