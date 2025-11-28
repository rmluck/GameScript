<script lang="ts">
    import { authStore } from '$lib/stores/auth';
</script>

<svelte:head>
    <title>GameScript</title>
</svelte:head>

<div class="text-center">
    <h1 class="font-display font-bold text-neutral-200 text-4xl sm:text-5xl md:text-6xl">
        Welcome to <span class="font-sans text-primary-400 text-5xl sm:text-6xl md:text-7xl">GameScript</span>
    </h1>
    <p class="mt-6 max-w-md mx-auto font-sans text-neutral-300 sm:text-xl md:mt-5 md:text-2xl md:max-w-4xl">
        Create custom playoff scenarios and see how your picks affect the standings, playoff seeding, and draft order.
    </p>

    <div class="mt-10 flex justify-center gap-x-6">
        {#if $authStore.isAuthenticated}
            <a
                href="/scenarios/new"
                class="rounded-md bg-primary-800 hover:bg-accent-500 border-2 border-primary-600 px-6 py-3 text-base font-semibold text-neutral-200 hover:text-accent-100 hover:border-accent-800 shadow-sm"
            >
                Create Scenario
            </a>
            <a
                href="/scenarios"
                class="rounded-md bg-primary-800 hover:bg-accent-500 border-2 border-primary-600 px-6 py-3 text-base font-semibold text-neutral-200 hover:text-accent-100 hover:border-accent-800 shadow-sm"
            >
                View My Scenarios
            </a>
        {:else}
            <a
                href="/auth/register"
                class="rounded-md bg-primary-800 hover:bg-accent-500 border-2 border-primary-600 px-6 py-3 font-sans font-semibold text-xl text-neutral-200 hover:text-accent-100 hover:border-accent-800 shadow-sm"
            >
                GET STARTED
            </a>
            <a
                href="/scenarios"
                class="rounded-md bg-primary-800 hover:bg-accent-500 border-2 border-primary-600 px-6 py-3 font-sans font-semibold text-xl text-neutral-200 hover:text-accent-100 hover:border-accent-800 shadow-sm"
            >
                BROWSE AS GUEST
            </a>
        {/if}
    </div>

    <div class="mt-20">
        <h2 class="font-heading font-bold text-2xl text-neutral-200">HOW IT WORKS</h2>
        <div class="mt-8 grid grid-cols-1 gap-8 sm:grid-cols-3">
            <div class="p-6 bg-white rounded-lg shadow">
                <div class="text-4xl font-bold text-primary-600 mb-4">1</div>
                <h3 class="text-lg font-semibold text-gray-900 mb-2">Create a Scenario</h3>
                <p class="text-gray-600">Start with the current NFL season and customize the name and settings.</p>
            </div>
            <div class="p-6 bg-white rounded-lg shadow">
                <div class="text-4xl font-bold text-primary-600 mb-4">2</div>
                <h3 class="text-lg font-semibold text-gray-900 mb-2">Make Your Picks</h3>
                <p class="text-gray-600">Pick winners and predict scores for upcoming games.</p>
            </div>
            <div class="p-6 bg-white rounded-lg shadow">
                <div class="text-4xl font-bold text-primary-600 mb-4">3</div>
                <h3 class="text-lg font-semibold text-gray-900 mb-2">See the Results</h3>
                <p class="text-gray-600">View updated standings, playoff seeding, and draft order based on your picks.</p>
            </div>
        </div>
    </div>
</div>