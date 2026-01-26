<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { scenariosAPI } from '$lib/api/scenarios';
    import type { Scenario } from '$types';

    // Props
    export let isOpen = false;
    export let scenario: Scenario;

    // Event dispatcher
    const dispatch = createEventDispatcher();

    // State variables for form fields
    let scenarioName = scenario.name;
    let isPublic = scenario.is_public;

    // Loading and error states
    let loading = false;
    let error = '';

    // Update form fields when scenario or isOpen changes
    $: if (isOpen) {
        scenarioName = scenario.name;
        isPublic = scenario.is_public;
    }

    // Handle form submission
    async function handleSave() {
        error = '';
        loading = true;

        try {
            const updated = await scenariosAPI.update(scenario.id, {
                name: scenarioName,
                is_public: isPublic
            });

            // Merge updated fields with existing scenario to preserve all data
            const mergedScenario = {
                ...scenario,
                ...updated
            }

            dispatch('updated', mergedScenario);
            close();
        } catch (err: any) {
            error = err.response?.data?.error || 'Failed to update scenario';
        } finally {
            loading = false;
        }
    }

    // Handle scenario deletion
    async function handleDelete() {
        if (!confirm('Are you sure you want to delete this scenario? This action cannot be undone.')) {
            return;
        }

        try {
            // Delete scenario and redirect to home page
            await scenariosAPI.delete(scenario.id);
            window.location.href = '/';
        } catch (err: any) {
            error = err.response?.data?.error || 'Failed to delete scenario';
        }
    }

    function close() {
        isOpen = false;
        error = '';
    }

    // Handle backdrop click to close modal
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
        <section class="bg-linear-to-br from-primary-975  to-primary-950 border-4 border-primary-900 rounded-lg max-w-xl w-full overflow-hidden flex flex-col">
            <!-- Header -->
            <div class="bg-primary-900/30 shadow-md p-6 flex justify-between items-center">
                <h2 class="text-2xl font-heading font-bold text-neutral">SCENARIO SETTINGS</h2>
                <button
                    on:click={close}
                    class="p-2 rounded-lg hover:bg-white/10 transition-colors cursor-pointer"
                    aria-label="Close modal"
                >
                    <svg class="w-6 h-6 text-neutral" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                </button>
            </div>

            {#if error}
                <div class="mb-4 p-3 bg-red-900/50 border border-red-600 rounded-md">
                    <p class="text-sm text-red-200 font-sans">{error}</p>
                </div>
            {/if}

            <!-- Form Content -->
            <div class="bg-linear-to-br from-primary-975  to-primary-950 flex-1 overflow-y-auto p-6">
                <form on:submit|preventDefault={handleSave} class="space-y-4">
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
                            class="w-full rounded-md bg-primary-900/60 border-2 border-primary-900 text-neutral placeholder-neutral/50 focus:border-primary-800 focus:ring-2 focus:ring-primary-800 px-4 py-3 font-sans"
                        />
                    </div>

                    <!-- Privacy Toggle -->
                    <div>
                        <label class="flex items-start gap-3 w-full rounded-md bg-primary-900/60 border-2 border-primary-900 px-4 py-3 transition-colors hover:bg-primary-900 cursor-pointer" class:border-primary-400={isPublic}>
                            <div class="relative flex items-center justify-center mt-2">
                                <input
                                    type="checkbox"
                                    id="isPublic"
                                    bind:checked={isPublic}
                                    class="sr-only peer"
                                />
                                <div class="w-6 h-6 bg-primary-900 border-2 border-primary-800 rounded peer-checked:bg-primary-800 peer-checked:border-primary-800 transition-all flex items-center justify-center">
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
                    <div class="flex gap-3 pt-4 pb-3">
                        <button
                            type="button"
                            on:click={close}
                            class="flex-1 bg-primary-900/60 hover:bg-primary-900 border-2 border-primary-900 rounded-lg py-3 font-sans font-semibold text-lg text-neutral transition-colors cursor-pointer"
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

                    <!-- Delete Button -->
                    <div class="pt-6 border-t border-primary-600/30">
                        <button
                            type="button"
                            on:click={handleDelete}
                            class="w-full bg-red-900/80 hover:bg-red-800 border-2 border-red-600 rounded-lg py-3 font-sans font-semibold text-lg text-red-100 hover:text-red-100 transition-colors cursor-pointer"
                        >
                            DELETE SCENARIO
                        </button>
                    </div>
                </form>
            </div>
        </section>
    </div>
{/if}