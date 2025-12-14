<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import type { Game, Pick } from '$types';

    export let game: Game;
    export let pick: Pick | undefined = undefined;
    export let compact: boolean = false;

    const dispatch = createEventDispatcher();

    let predictedHomeScore = pick?.predicted_home_score?.toString() || '';
    let predictedAwayScore = pick?.predicted_away_score?.toString() || '';

    $: {
        predictedHomeScore = pick?.predicted_home_score?.toString() || '';
        predictedAwayScore = pick?.predicted_away_score?.toString() || '';
    }

    $: isGameCompleted = game.status === 'final';
    
    $: isHomeTeamPicked = pick?.picked_team_id === game.home_team_id;
    $: isAwayTeamPicked = pick?.picked_team_id === game.away_team_id;
    $: isTiePicked = pick !== undefined && pick.picked_team_id === 0;

    $: userMadePick = pick !== undefined;

    $: homeTeamWon = isGameCompleted && 
        game.home_score !== null && 
        game.away_score !== null && 
        game.home_score > game.away_score;
    
    $: awayTeamWon = isGameCompleted && 
        game.home_score !== null && 
        game.away_score !== null && 
        game.away_score > game.home_score;

    $: highlightHomeButton = userMadePick ? isHomeTeamPicked : (isGameCompleted ? homeTeamWon : false);
    $: highlightAwayButton = userMadePick ? isAwayTeamPicked : (isGameCompleted ? awayTeamWon : false);

    $: homeTeamLogoURL = highlightHomeButton && game.home_team.alternate_logo_url 
        ? game.home_team.alternate_logo_url 
        : game.home_team.logo_url;
    
    $: awayTeamLogoURL = highlightAwayButton && game.away_team.alternate_logo_url 
        ? game.away_team.alternate_logo_url 
        : game.away_team.logo_url;

    function selectTeam(teamId: number) {
        if (pick?.picked_team_id === teamId) {
            dispatch('pickChanged', {
                gameId: game.id,
                pickedTeamId: undefined,
                predictedHomeScore: predictedHomeScore ? parseInt(predictedHomeScore) : undefined,
                predictedAwayScore: predictedAwayScore ? parseInt(predictedAwayScore) : undefined
            });
        } else {
            dispatch('pickChanged', {
                gameId: game.id,
                pickedTeamId: teamId,
                predictedHomeScore: predictedHomeScore ? parseInt(predictedHomeScore) : undefined,
                predictedAwayScore: predictedAwayScore ? parseInt(predictedAwayScore) : undefined
            });
        }
    }

    function selectTie() {
        if (isTiePicked) {
            dispatch('pickChanged', {
                gameId: game.id,
                pickedTeamId: undefined,
                predictedHomeScore: predictedHomeScore ? parseInt(predictedHomeScore) : undefined,
                predictedAwayScore: predictedAwayScore ? parseInt(predictedAwayScore) : undefined
            });
        } else {
            dispatch('pickChanged', {
                gameId: game.id,
                pickedTeamId: 0,
                predictedHomeScore: predictedHomeScore ? parseInt(predictedHomeScore) : undefined,
                predictedAwayScore: predictedAwayScore ? parseInt(predictedAwayScore) : undefined
            });
        }
    }

    function handleScoreChange() {
        if (pick || predictedHomeScore || predictedAwayScore) {
            dispatch('pickChanged', {
                gameId: game.id,
                pickedTeamId: pick?.picked_team_id,
                predictedHomeScore: predictedHomeScore ? parseInt(predictedHomeScore) : undefined,
                predictedAwayScore: predictedAwayScore ? parseInt(predictedAwayScore) : undefined
            });
        }
    }
</script>

<style>
    input[type='number']::-webkit-inner-spin-button,
    input[type='number']::-webkit-outer-spin-button {
        -webkit-appearance: none;
        margin: 0;
    }
    input[type='number'] {
        -moz-appearance: textfield;
        appearance: textfield;
    }
</style>

<div class="flex items-center gap-2">
    <!-- Away Team Section -->
    <div class="flex-1 flex items-stretch">
        <!-- Away Score -->
        <div class="w-12 shrink-0">
            {#if isGameCompleted && game.away_score !== null && !userMadePick}
                <!-- Show actual score as non-editable -->
                <div class="h-full flex items-center justify-center font-heading text-xl font-bold rounded-l-lg border-2 border-r-0"
                     style={`background-color: #${game.away_team.primary_color}90; border-color: #${game.away_team.primary_color}; color: #${game.away_team.primary_color};`}>
                    {game.away_score}
                </div>
            {:else}
                <!-- Editable input - always enabled -->
                <input
                    type="number"
                    min="0"
                    max="99"
                    bind:value={predictedAwayScore}
                    on:change={handleScoreChange}
                    placeholder="--"
                    class="h-full w-full text-center border-2 border-r-0 rounded-l-lg px-1 py-2 font-heading text-lg font-bold placeholder-neutral/40 transition-colors focus:outline-none"
                    style={`background-color: #${game.away_team.primary_color}90; border-color: #${game.away_team.primary_color}; color: #${game.away_team.primary_color};`}
                />
            {/if}
        </div>

        <!-- Away Team Button -->
        <button
            on:click={() => selectTeam(game.away_team_id)}
            class="flex-1 p-2 rounded-r-lg border-2 transition-all cursor-pointer"
            class:border-primary-600={!highlightAwayButton}
            class:hover:border-primary-400={!highlightAwayButton}
            style={highlightAwayButton 
                ? `background-color: #${game.away_team.primary_color}; border-color: #${game.away_team.primary_color};` 
                : `background-color: transparent; border-color: #${game.away_team.primary_color};`}
        >
            <div class="flex items-center justify-center gap-2">
                {#if awayTeamLogoURL}
                    <img 
                        src={awayTeamLogoURL}
                        alt={game.away_team.abbreviation}
                        class="w-8 h-8 object-contain"
                    />
                {/if}
                {#if compact}
                    <span class="text-sm font-sans font-semibold transition-colors"
                          class:text-white={highlightAwayButton}
                          class:text-black={!highlightAwayButton}>
                        {game.away_team.abbreviation}
                    </span>
                {/if}
            </div>
        </button>
    </div>

    <!-- Tie Button -->
    <button
        on:click={selectTie}
        class="px-2 sm:px-3 py-1 sm:py-2 rounded border text-xs font-sans font-semibold transition-all shrink-0 cursor-pointer"
        class:bg-primary-600={isTiePicked}
        class:border-primary-500={isTiePicked}
        class:text-neutral={isTiePicked}
        class:bg-transparent={!isTiePicked}
        class:border-primary-600={!isTiePicked}
        class:hover:border-primary-400={!isTiePicked}
    >
        TIE
    </button>

    <!-- Home Team Section -->
    <div class="flex-1 flex items-stretch">
        <!-- Home Team Button -->
        <button
            on:click={() => selectTeam(game.home_team_id)}
            class="flex-1 p-2 rounded-l-lg border-2 transition-all cursor-pointer"
            class:border-primary-600={!highlightHomeButton}
            class:hover:border-primary-400={!highlightHomeButton}
            style={highlightHomeButton 
                ? `background-color: #${game.home_team.primary_color}; border-color: #${game.home_team.primary_color};` 
                : `background-color: transparent; border-color: #${game.home_team.primary_color};`}
        >
            <div class="flex items-center justify-center gap-2">
                {#if homeTeamLogoURL}
                    <img 
                        src={homeTeamLogoURL}
                        alt={game.home_team.abbreviation}
                        class="w-8 h-8 object-contain"
                    />
                {/if}
                {#if compact}
                    <span class="text-sm font-sans font-semibold transition-colors"
                          class:text-white={highlightHomeButton}
                          class:text-black={!highlightHomeButton}>
                        {game.home_team.abbreviation}
                    </span>
                {/if}
            </div>
        </button>

        <!-- Home Score -->
        <div class="w-12 shrink-0">
            {#if isGameCompleted && game.home_score !== null && !userMadePick}
                <!-- Show actual score as non-editable -->
                <div class="h-full flex items-center justify-center font-heading text-xl font-bold rounded-r-lg border-2 border-l-0"
                     style={`background-color: #${game.home_team.primary_color}90; border-color: #${game.home_team.primary_color}; color: #${game.home_team.primary_color};`}>
                    {game.home_score}
                </div>
            {:else}
                <!-- Editable input - always enabled -->
                <input
                    type="number"
                    min="0"
                    max="99"
                    bind:value={predictedHomeScore}
                    on:change={handleScoreChange}
                    placeholder="--"
                    class="h-full w-full text-center border-2 border-l-0 rounded-r-lg px-1 py-2 font-heading text-lg font-bold placeholder-neutral/40 transition-colors focus:outline-none"
                    style={`background-color: #${game.home_team.primary_color}90; border-color: #${game.home_team.primary_color}; color: #${game.home_team.primary_color};`}
                />
            {/if}
        </div>
    </div>
</div>