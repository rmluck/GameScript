<script lang="ts">
    import { authStore } from '$lib/stores/auth';
</script>

<svelte:head>
    <title>GameScript</title>
</svelte:head>

<div class="text-center">
    <h1 class="text-4xl font-bold text-gray-900 sm:text-5xl md:text-6xl">
        Welcome to <span class="text-primary-600">GameScript</span>
    </h1>
    <p class="mt-3 max-w-md mx-auto text-base text-gray-500 sm:text-lg md:mt-5 md:text-xl md:max-w-3xl">
        Create custom playoff scenarios and see how your picks affect the standings, playoff seeding, and draft order.
    </p>

    <div class="mt-10 flex justify-center gap-x-6">
        {#if $authStore.isAuthenticated}
            <a
                href="/scenarios/new"
                class="rounded-md bg-primary-600 px-6 py-3 text-base font-semibold text-white shadow-sm hover:bg-primary-700 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600"
            >
                Create Scenario
            </a>
            <a
                href="/scenarios"
                class="rounded-md bg-white px-6 py-3 text-base font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50"
            >
                View My Scenarios
            </a>
        {:else}
            <a
                href="/auth/register"
                class="rounded-md bg-primary-600 px-6 py-3 text-base font-semibold text-gray-900 shadow-sm hover:bg-primary-700"
            >
                Get Started
            </a>
            <a
                href="/scenarios"
                class="rounded-md bg-white px-6 py-3 text-base font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50"
            >
                Browse as Guest
            </a>
        {/if}
    </div>

    <div class="mt-20">
        <h2 class="text-2xl font-bold text-gray-900">How it works</h2>
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