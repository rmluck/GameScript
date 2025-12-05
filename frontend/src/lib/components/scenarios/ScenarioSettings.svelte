<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { scenariosAPI } from '$lib/api/scenarios';
    import type { Scenario } from '$types';

    export let isOpen = false;
    export let scenario: Scenario;

    const dispatch = createEventDispatcher();

    let scenarioName = scenario.name;
    let isPublic = scenario.is_public;
    let loading = false;
    let error = '';

    $: if (isOpen) {
        scenarioName = scenario.name;
        isPublic = scenario.is_public;
    }

    async function handleSave() {
        error = '';
        loading = true;

        try {
            const updated = await scenariosAPI.update(scenario.id, {
                name: scenarioName,
                is_public: isPublic
            });

            dispatch('updated', updated);
            close();
        } catch (err: any) {
            error = err.response?.data?.error || 'Failed to update scenario';
        } finally {
            loading = false;
        }
    }

    async function handleDelete() {
        if (!confirm('Are you sure you want to delete this scenario? This action cannot be undone.')) {
            return;
        }

        try {
            await scenariosAPI.delete(scenario.id);
            window.location.href = '/scenarios';
        } catch (err: any) {
            error = err.response?.data?.error || 'Failed to delete scenario';
        }
    }

    function close() {
        isOpen = false;
        error = '';
    }

    function handleBackdropClick(event: MouseEvent) {
        if (event.target === event.currentTarget) {
            close();
        }
    }
</script>

{#if isOpen}
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <div 
        class="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50" 
        on:click={handleBackdropClick}
        on:keydown={(e) => e.key === 'Escape' && close()}
    >
        <div class="bg-primary-900 border-2 border-primary-700 rounded-lg shadow-xl max-w-lg w-full mx-4 p-6">
            <div class="flex justify-between items-center mb-6">
                <h2 class="text-2xl font-heading font-bold text-neutral">SCENARIO SETTINGS</h2>
                <button on:click={close} class="text-neutral hover:text-primary-400 transition-colors cursor-pointer" title="Close Settings">
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

            <form on:submit|preventDefault={handleSave} class="space-y-6">
                <!-- Scenario Name -->
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
                        class="w-full rounded-md bg-primary-800/60 border-2 border-primary-600 text-neutral placeholder-neutral/50 focus:border-primary-400 focus:ring-2 focus:ring-primary-400 px-4 py-3 font-sans"
                    />
                </div>

                <!-- Privacy Toggle -->
                <div>
                    <label class="flex items-start gap-3 w-full rounded-md bg-primary-800/60 border-2 border-primary-600 px-4 py-3 transition-colors hover:bg-primary-800 cursor-pointer" class:border-primary-400={isPublic}>
                        <div class="relative flex items-center justify-center mt-2">
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
                        <div class="flex-1">
                            <span class="block text-md font-sans text-neutral select-none">
                                Make this scenario public
                            </span>
                            <p class="text-sm text-neutral/70 font-sans">
                                {isPublic ? 'Anyone can view this scenario' : 'Only you can view this scenario'}
                            </p>
                        </div>
                    </label>
                </div>

                <!-- Action Buttons -->
                <div class="flex gap-3 pt-4">
                    <button
                        type="button"
                        on:click={close}
                        class="flex-1 bg-primary-800/60 hover:bg-primary-700 border-2 border-primary-600 rounded-lg py-3 font-sans font-semibold text-lg text-neutral transition-colors cursor-pointer"
                    >
                        CANCEL
                    </button>
                    <button
                        type="submit"
                        disabled={loading}
                        class="flex-1 bg-primary-600 hover:bg-primary-500 border-2 border-primary-500 hover:border-primary-400 rounded-lg py-3 font-sans font-semibold text-lg text-neutral transition-all hover:scale-105 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:scale-100 cursor-pointer"
                    >
                        {loading ? 'SAVING...' : 'SAVE CHANGES'}
                    </button>
                </div>

                <!-- Danger Zone -->
                <div class="pt-6 border-t border-red-600/30">
                    <button
                        type="button"
                        on:click={handleDelete}
                        class="w-full bg-red-900/80 hover:bg-red-800/40 border-2 border-red-600 rounded-lg py-3 font-sans font-semibold text-lg text-red-200 hover:text-red-100 transition-colors cursor-pointer"
                    >
                        DELETE SCENARIO
                    </button>
                </div>
            </form>
        </div>
    </div>
{/if}