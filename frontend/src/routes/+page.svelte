<script lang="ts">
    import { goto } from '$app/navigation'
    import { authStore } from '$lib/stores/auth';
    import CreateScenarioModal from '$lib/components/scenarios/CreateScenarioModal.svelte';
    import ComingSoonModal from '$lib/components/scenarios/ComingSoonModal.svelte';

    let showCreateModal = false;
    let showComingSoonModal = false;

    function handleScenarioCreated(event: CustomEvent) {
        const scenario = event.detail;
        goto(`/scenarios/${scenario.id}`);
    }

    function handleCFBClick(event: Event) {
        event.preventDefault();
        showComingSoonModal = true;
    }
</script>

<svelte:head>
    <title>GameScript</title>
    <meta name="description" content="Create custom playoff scenarios and see how your picks affect the standings, playoff seeding, and draft order." />
</svelte:head>

<CreateScenarioModal 
    bind:isOpen={showCreateModal} 
    on:created={handleScenarioCreated}
/>

<ComingSoonModal 
    bind:isOpen={showComingSoonModal} 
    feature="College Football"
/>

<div class="text-center px-4 sm:px-6 md:px-8">
    <h1 class="font-display font-bold text-neutral text-4xl sm:text-5xl md:text-6xl mt-16">
        Welcome to
        <span class="bg-linear-to-r from-primary-700 via-primary-600 to-primary-500 bg-clip-text text-transparent font-sans text-5xl sm:text-6xl md:text-7xl">GameScript</span>
    </h1>
    <p class="mt-6 max-w-md mx-auto font-sans text-neutral text-xl sm:text-2xl md:max-w-4xl">
        Create custom playoff scenarios and see how your picks affect the standings, playoff seeding, and draft order.
    </p>

    <!-- CTA Buttons -->
    <div class="mt-8 flex flex-col sm:flex-row justify-center gap-4 sm:gap-x-6">
        {#if $authStore.isAuthenticated}
            <button
                on:click={() => showCreateModal = true}
                class="bg-primary-900/60 hover:bg-primary-600 border-2 border-primary-900 hover:border-primary-500 rounded-lg shadow-lg transition-all hover:scale-105 px-6 py-3 font-sans font-semibold text-lg sm:text-xl text-neutral cursor-pointer"
            >
                CREATE SCENARIO
            </button>
            <a
                href="/scenarios"
                class="bg-primary-900/60 hover:bg-primary-600 border-2 border-primary-900 hover:border-primary-500 rounded-lg shadow-lg transition-all hover:scale-105 px-6 py-3 font-sans font-semibold text-lg sm:text-xl text-neutral"
            >
                VIEW SCENARIOS
            </a>
        {:else}
            <a
                href="/auth/register"
                class="bg-primary-900/60 hover:bg-primary-600 border-2 border-primary-900 hover:border-primary-500 rounded-lg shadow-lg transition-all hover:scale-105 px-6 py-3 font-sans font-semibold text-lg sm:text-xl text-neutral"
            >
                GET STARTED
            </a>
            <a
                href="/scenarios"
                class="bg-primary-900/60 hover:bg-primary-600 border-2 border-primary-900 hover:border-primary-500 rounded-lg shadow-lg transition-all hover:scale-105 px-6 py-3 font-sans font-semibold text-lg sm:text-xl text-neutral"
            >
                BROWSE AS GUEST
            </a>
        {/if}
    </div>

    <!-- How It Works Section -->
    <div class="mt-20">
        <h2 class="font-heading font-bold text-2xl sm:text-3xl text-neutral mb-4">HOW IT WORKS</h2>
        <div class="mt-10 grid grid-cols-1 gap-6 sm:gap-8 sm:grid-cols-3">
            <div class="p-6 sm:p-8 bg-neutral rounded-lg shadow">
                <div class="text-4xl sm:text-5xl font-bold text-primary-400 mb-4">1</div>
                <h3 class="text-lg sm:text-xl font-heading font-semibold text-primary-600 mb-2">Create a Scenario</h3>
                <p class="font-sans text-base sm:text-lg text-gray-600">Choose a sport and customize the name and settings of your season simulation.</p>
            </div>
            <div class="p-6 sm:p-8 bg-neutral rounded-lg shadow">
                <div class="text-4xl sm:text-5xl font-bold text-primary-400 mb-4">2</div>
                <h3 class="text-lg sm:text-xl font-heading font-semibold text-primary-600 mb-2">Make Your Picks</h3>
                <p class="font-sans text-base sm:text-lg text-gray-600">Pick winners and predict scores for upcoming games.</p>
            </div>
            <div class="p-6 sm:p-8 bg-neutral rounded-lg shadow">
                <div class="text-4xl sm:text-5xl font-bold text-primary-400 mb-4">3</div>
                <h3 class="text-lg sm:text-xl font-heading font-semibold text-primary-600 mb-2">See the Results</h3>
                <p class="font-sans text-base sm:text-lg text-gray-600">View updated standings, playoff seeding, and draft order based on your picks.</p>
            </div>
        </div>
    </div>

    <!-- Sports Selection Section -->
    <div class="mt-24">
        <h2 class="font-heading font-bold text-2xl sm:text-3xl text-neutral mb-4">CHOOSE YOUR SPORT</h2>
        <p class="font-sans text-neutral text-lg sm:text-xl mb-8 max-w-2xl mx-auto">
            Create scenarios for your favorite leagues and see how the playoff picture changes.
        </p>
        <div class="grid grid-cols-3 gap-3 sm:gap-6 max-w-5xl mx-auto">
            <a href="/nfl" class="group p-4 sm:p-8 bg-primary-900/60 hover:bg-primary-600 border-2 border-primary-900 hover:border-primary-500 rounded-lg shadow-lg transition-all hover:scale-105">
                <h3 class="font-display text-2xl sm:text-4xl text-neutral  sm:mb-2">NFL</h3>
                <p class="hidden sm:block font-sans text-neutral text-lg">National Football League</p>
            </a>
            <a href="/nba" class="group p-4 sm:p-8 bg-primary-900/60 hover:bg-primary-600 border-2 border-primary-900 hover:border-primary-500 rounded-lg shadow-lg transition-all hover:scale-105">
                <h3 class="font-display text-2xl sm:text-4xl text-neutral  sm:mb-2">NBA</h3>
                <p class="hidden sm:block font-sans text-neutral text-lg">National Basketball Association</p>
            </a>
            <button on:click={handleCFBClick} class="group p-4 sm:p-8 bg-primary-900/60 hover:bg-primary-600 border-2 border-primary-900 hover:border-primary-500 rounded-lg shadow-lg transition-all hover:scale-105 relative cursor-pointer">
                <div class="absolute top-2 right-2 bg-yellow-500/90 text-primary-900 text-xs font-sans font-bold px-2 py-1 rounded">SOON</div>
                <h3 class="font-display text-2xl sm:text-4xl text-neutral  sm:mb-2">CFB</h3>
                <p class="hidden sm:block font-sans text-neutral text-lg">College Football</p>
            </button>
        </div>
    </div>

    <!-- Features Section -->
    <div class="mt-24 max-w-4xl mx-auto">
        <h2 class="font-heading font-bold text-2xl sm:text-3xl text-neutral mb-8">KEY FEATURES</h2>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 sm:gap-6 text-left">
            <div class="p-4 sm:p-6 bg-primary-900/60 border-l-4 border-primary-500 rounded-r-lg">
                <h3 class="font-heading text-base sm:text-lg font-semibold text-neutral mb-2">📝 Custom Scenarios</h3>
                <p class="font-sans text-neutral text-sm">Create and name your own season scenarios with personalized settings.</p>
            </div>
            <div class="p-4 sm:p-6 bg-primary-900/60 border-l-4 border-primary-500 rounded-r-lg">
                <h3 class="font-heading text-base sm:text-lg font-semibold text-neutral mb-2">📊 Live Standings</h3>
                <p class="font-sans text-neutral text-sm">See real-time updates to division standings as you make picks.</p>
            </div>
            <div class="p-4 sm:p-6 bg-primary-900/60 border-l-4 border-primary-500 rounded-r-lg">
                <h3 class="font-heading text-base sm:text-lg font-semibold text-neutral mb-2">🏆 Playoff Seeding</h3>
                <p class="font-sans text-neutral text-sm">Watch playoff brackets update automatically using official tiebreaker rules.</p>
            </div>
            <div class="p-4 sm:p-6 bg-primary-900/60 border-l-4 border-primary-500 rounded-r-lg">
                <h3 class="font-heading text-base sm:text-lg font-semibold text-neutral mb-2">📋 Draft Order</h3>
                <p class="font-sans text-neutral text-sm">Track how teams move in the draft order based on records.</p>
            </div>
            <div class="p-4 sm:p-6 bg-primary-900/60 border-l-4 border-primary-500 rounded-r-lg">
                <h3 class="font-heading text-base sm:text-lg font-semibold text-neutral mb-2">⚙️ User-Friendly Interface</h3>
                <p class="font-sans text-neutral text-sm">Intuitive design for easy navigation and scenario management.</p>
            </div>
            <div class="p-4 sm:p-6 bg-primary-900/60 border-l-4 border-primary-500 rounded-r-lg">
                <h3 class="font-heading text-base sm:text-lg font-semibold text-neutral mb-2">💾 Save & Share</h3>
                <p class="font-sans text-neutral text-sm">Save multiple scenarios and share them with friends.</p>
            </div>
        </div>
    </div>
</div>