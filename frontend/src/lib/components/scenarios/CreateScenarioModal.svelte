<script lang="ts">
    import { createEventDispatcher, onMount } from 'svelte';
    import { scenariosAPI } from '$lib/api/scenarios';
    import type { Sport, Season } from '$types';

    export let isOpen = false;

    const dispatch = createEventDispatcher();

    let sports: Sport[] = [];
    let seasons: Season[] = [];
    let selectedSportId: number | null = null;
    let selectedSeasonId: number | null = null;
    let scenarioName = '';
    let isPublic = false;
    let loading = false;
    let error = '';
    
    // Custom dropdown state
    let sportDropdownOpen = false;
    let seasonDropdownOpen = false;

    onMount(async () => {
        // Load sports
        const response = await fetch('/api/sports');
        sports = await response.json();
        console.log('Loaded sports:', sports);
    });

    function formatDate() {
        const today = new Date();
        const month = String(today.getMonth() + 1).padStart(2, '0');
        const day = String(today.getDate()).padStart(2, '0');
        const year = today.getFullYear();
        return `${month}/${day}/${year}`;
    }

    async function loadSeasons() {
        if (!selectedSportId) return;
        
        const response = await fetch(`/api/sports/${selectedSportId}/seasons`);
        seasons = await response.json();
        
        // Auto-select active season
        const activeSeason = seasons.find(s => s.is_active);
        if (activeSeason) {
            selectedSeasonId = activeSeason.id;
            
            // Auto-generate name if empty
            if (!scenarioName) {
                const sport = sports.find(s => s.id === selectedSportId);
                const todayDate = formatDate();

                if (sport && activeSeason.end_year) {
                    scenarioName = `${sport.short_name} ${activeSeason.start_year}-${activeSeason.end_year} Season - ${todayDate}`;
                } else if (sport) {
                    scenarioName = `${sport.short_name} ${activeSeason.start_year} Season - ${todayDate}`;
                }
            }
        }
    }

    $: if (selectedSportId) {
        loadSeasons();
    }
    
    function selectSport(sportId: number) {
        selectedSportId = sportId;
        sportDropdownOpen = false;
    }
    
    function selectSeason(seasonId: number) {
        selectedSeasonId = seasonId;
        seasonDropdownOpen = false;
    }
    
    $: selectedSportName = selectedSportId 
        ? sports.find(s => s.id === selectedSportId)?.name 
        : null;
    
    $: selectedSeasonName = selectedSeasonId
        ? seasons.find(s => s.id === selectedSeasonId)
        : null;

    async function handleSubmit() {
        error = '';
        
        if (!scenarioName || !selectedSportId || !selectedSeasonId) {
            error = 'Please fill in all fields';
            return;
        }

        loading = true;

        try {
            const scenario = await scenariosAPI.create({
                name: scenarioName,
                sport_id: selectedSportId,
                season_id: selectedSeasonId,
                is_public: isPublic
            });

            dispatch('created', scenario);
            close();
        } catch (err: any) {
            error = err.response?.data?.error || 'Failed to create scenario';
        } finally {
            loading = false;
        }
    }

    function close() {
        isOpen = false;
        dispatch('close');
        
        // Reset form
        selectedSportId = null;
        selectedSeasonId = null;
        scenarioName = '';
        isPublic = false;
        error = '';
        sportDropdownOpen = false;
        seasonDropdownOpen = false;
    }

    function handleBackdropClick(event: MouseEvent) {
        if (event.target === event.currentTarget) {
            close();
        }
    }

    function handleKeydown(event: KeyboardEvent) {
        if (event.key === 'Escape' && isOpen) {
            close();
        }
    }
</script>

{#if isOpen}
    <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
    <div class="fixed inset-0 bg-black/30 backdrop-blur-md flex items-center justify-center z-50" on:click={handleBackdropClick} on:keydown={handleKeydown} role="dialog" aria-modal="true" aria-labelledby="modal-title" tabindex="-1">
        <div class="bg-primary-900 border-2 border-primary-700 rounded-lg shadow-xl max-w-md w-full mx-4 p-6">
            <div class="flex justify-between items-center mb-6">
                <h2 class="text-2xl font-heading font-bold text-neutral">CREATE SCENARIO</h2>
                <button on:click={close} class="text-neutral hover:text-primary-400 transition-colors cursor-pointer" aria-label="Close modal">
                    <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                </button>
            </div>

            {#if error}
                <div class="mb-4 p-3 bg-red-900/50 border border-red-600 rounded-md">
                    <p class="text-sm text-red-200 font-sans">{error}</p>
                </div>
            {/if}

            <form on:submit|preventDefault={handleSubmit} class="space-y-4">
                <!-- Custom Sport Dropdown -->
                <div class="relative">
                    <label for="sport" class="block text-lg font-semibold font-sans text-neutral mb-2">
                        Sport
                    </label>
                    <button
                        type="button"
                        on:click={() => sportDropdownOpen = !sportDropdownOpen}
                        class="w-full rounded-md bg-primary-800/60 border-2 border-primary-600 focus:border-primary-400 focus:ring-2 focus:ring-primary-400 px-4 py-3 font-sans text-md text-left transition-colors hover:bg-primary-800 flex justify-between items-center {selectedSportId ? 'text-neutral' : 'text-neutral/50'}"
                        class:border-primary-400={sportDropdownOpen}
                    >
                        <span>
                            {#if selectedSportId && selectedSportName}
                                {selectedSportName} ({sports.find(s => s.id === selectedSportId)?.short_name})
                            {:else}
                                Select a sport...
                            {/if}
                        </span>
                        <svg class="w-5 h-5 transition-transform" class:rotate-180={sportDropdownOpen} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                        </svg>
                    </button>
                    
                    {#if sportDropdownOpen}
                        <div class="absolute z-10 w-full mt-1 bg-primary-800 border-2 border-primary-600 rounded-md shadow-lg max-h-60 overflow-auto">
                            {#each sports as sport}
                                <button
                                    type="button"
                                    on:click={() => selectSport(sport.id)}
                                    class="w-full px-4 py-3 text-left text-neutral hover:bg-primary-700 transition-colors font-sans"
                                    class:bg-primary-700={selectedSportId === sport.id}
                                >
                                    {sport.name} ({sport.short_name})
                                </button>
                            {/each}
                        </div>
                    {/if}
                </div>

                {#if seasons.length > 0}
                    <!-- Custom Season Dropdown -->
                    <div class="relative">
                        <label for="season" class="block text-lg font-semibold font-sans text-neutral mb-2">
                            Season
                        </label>
                        <button
                            type="button"
                            on:click={() => seasonDropdownOpen = !seasonDropdownOpen}
                            class="w-full rounded-md bg-primary-800/60 border-2 border-primary-600 focus:border-primary-400 focus:ring-2 focus:ring-primary-400 px-4 py-3 font-sans text-md text-left transition-colors hover:bg-primary-800 flex justify-between items-center {selectedSeasonId ? 'text-neutral' : 'text-neutral/50'}"
                            class:border-primary-400={seasonDropdownOpen}
                        >
                            <span>
                                {#if selectedSeasonId && selectedSeasonName}
                                    {selectedSeasonName.start_year}{selectedSeasonName.end_year ? `-${selectedSeasonName.end_year}` : ''} 
                                    {selectedSeasonName.is_active ? '(Current)' : ''}
                                {:else}
                                    Select a season...
                                {/if}
                            </span>
                            <svg class="w-5 h-5 transition-transform" class:rotate-180={seasonDropdownOpen} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                            </svg>
                        </button>
                        
                        {#if seasonDropdownOpen}
                            <div class="absolute z-10 w-full mt-1 bg-primary-800 border-2 border-primary-600 rounded-md shadow-lg max-h-60 overflow-auto">
                                {#each seasons as season}
                                    <button
                                        type="button"
                                        on:click={() => selectSeason(season.id)}
                                        class="w-full px-4 py-3 text-left text-neutral hover:bg-primary-700 transition-colors font-sans"
                                        class:bg-primary-700={selectedSeasonId === season.id}
                                    >
                                        {season.start_year}{season.end_year ? `-${season.end_year}` : ''} 
                                        {season.is_active ? '(Current)' : ''}
                                    </button>
                                {/each}
                            </div>
                        {/if}
                    </div>
                {/if}

                <div>
                    <label for="name" class="block text-lg font-semibold font-sans text-neutral mb-2">
                        Scenario Name
                    </label>
                    <input
                        type="text"
                        id="name"
                        bind:value={scenarioName}
                        required
                        maxlength="100"
                        class="w-full rounded-md bg-primary-800/60 border-2 border-primary-600 text-neutral placeholder-neutral/50 focus:border-primary-400 focus:ring-2 focus:ring-primary-400 px-4 py-3 font-sans text-md"
                        placeholder="My 2024 Predictions"
                    />
                </div>

                <!-- Custom Checkbox -->
                <div class="relative">
                    <label class="flex items-center gap-3 w-full rounded-md bg-primary-800/60 border-2 border-primary-600 px-4 py-3 transition-colors hover:bg-primary-800 cursor-pointer" class:border-primary-400={isPublic}>
                        <div class="relative flex items-center justify-center">
                            <input
                                type="checkbox"
                                id="isPublic"
                                bind:checked={isPublic}
                                class="sr-only peer"
                            />
                            <div class="w-6 h-6 bg-primary-900 border-2 border-primary-600 rounded peer-checked:bg-primary-500 peer-checked:border-primary-400 transition-all flex items-center justify-center">
                                {#if isPublic}
                                    <svg class="w-4 h-4 text-neutral" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7" />
                                    </svg>
                                {/if}
                            </div>
                        </div>
                        <span class="text-md font-sans text-neutral select-none">
                            Make this scenario public
                        </span>
                    </label>
                </div>

                <div class="flex gap-3 mt-6">
                    <button
                        type="button"
                        on:click={close}
                        class="flex-1 bg-primary-800/60 hover:bg-primary-700 border-2 border-primary-600 rounded-lg py-3 font-sans font-semibold text-xl text-neutral transition-colors cursor-pointer"
                    >
                        CANCEL
                    </button>
                    <button
                        type="submit"
                        disabled={loading}
                        class="flex-1 bg-primary-600 hover:bg-primary-500 border-2 border-primary-500 hover:border-primary-400 rounded-lg py-3 font-sans font-semibold text-xl text-neutral transition-all hover:scale-105 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:scale-100 cursor-pointer"
                    >
                        {loading ? 'CREATING...' : 'CREATE'}
                    </button>
                </div>
            </form>
        </div>
    </div>
{/if}