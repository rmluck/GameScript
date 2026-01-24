<script lang="ts">
    import { onMount } from 'svelte';
    import { authStore } from '$lib/stores/auth';
    import { authAPI } from '$lib/api/auth';
    import { scenariosAPI } from '$lib/api/scenarios';
    import { goto } from '$app/navigation';
    import CreateScenarioModal from '$lib/components/scenarios/CreateScenarioModal.svelte';
    import type { User, Scenario } from '$types';

    let user: User | null = null;
    let scenarios: Scenario[] = [];
    let loading = true;
    let error = '';

    // Settings state
    let showSettings = false;
    let username = '';
    let email = '';
    let currentPassword = '';
    let newPassword = '';
    let confirmPassword = '';
    let settingsSaving = false;
    let settingsError = '';
    let settingsSuccess = '';

    // Scenario management
    let deletingScenarioId: number | null = null;
    let showCreateModal = false;

    // Share functionality
    let copiedScenarioId: number | null = null;
    let copiedTimeout: ReturnType<typeof setTimeout>;

    onMount(async () => {
        // Check if user is authenticated
        if (!$authStore.isAuthenticated) {
            goto('/auth/login');
            return;
        }

        await loadUserData();
        await loadScenarios();
    });

    async function loadUserData() {
        try {
            user = await authAPI.getCurrentUser();
            username = user.username;
            email = user.email;
        } catch (err: any) {
            error = 'Failed to load user data';
            console.error(err);
        }
    }

    async function loadScenarios() {
        try {
            loading = true;
            scenarios = await scenariosAPI.getAll();
        } catch (err: any) {
            error = 'Failed to load scenarios';
            console.error(err);
        } finally {
            loading = false;
        }
    }

    async function handleUpdateProfile() {
        settingsSaving = true;
        settingsError = '';
        settingsSuccess = '';

        try {
            // Validate inputs
            if (username.length < 3) {
                settingsError = 'Username must be at least 3 characters';
                settingsSaving = false;
                return;
            }

            if (newPassword && newPassword !== confirmPassword) {
                settingsError = 'New passwords do not match';
                settingsSaving = false;
                return;
            }

            const updateData: any = { username, email };
            
            if (newPassword) {
                updateData.current_password = currentPassword;
                updateData.new_password = newPassword;
            }

            const updatedUser = await authAPI.updateProfile(updateData);
            authStore.updateUser(updatedUser);
            
            settingsSuccess = 'Profile updated successfully!';
            currentPassword = '';
            newPassword = '';
            confirmPassword = '';
            
            setTimeout(() => {
                settingsSuccess = '';
                showSettings = false;
            }, 2000);
        } catch (err: any) {
            settingsError = err.response?.data?.error || 'Failed to update profile';
        } finally {
            settingsSaving = false;
        }
    }

    async function handleDeleteScenario(scenarioId: number) {
        if (!confirm('Are you sure you want to delete this scenario? This action cannot be undone.')) {
            return;
        }

        deletingScenarioId = scenarioId;
        try {
            await scenariosAPI.delete(scenarioId);
            scenarios = scenarios.filter(s => s.id !== scenarioId);
        } catch (err: any) {
            alert('Failed to delete scenario');
            console.error(err);
        } finally {
            deletingScenarioId = null;
        }
    }

    async function handleShareScenario(scenario: Scenario) {
        const url = `${window.location.origin}/scenarios/${scenario.sport_short_name?.toLowerCase()}/${scenario.id}`;

        // Try native share API first (mobile devices)
        if (navigator.share) {
            try {
                await navigator.share({
                    title: `GameScript - ${scenario.name}`,
                    text: `Check out my scenario ${scenario.name} on GameScript!`,
                    url: url
                });
                return;
            } catch (err) {
                // User cancelled or share failed, fall back to clipboard
                if ((err as Error).name !== 'AbortError') {
                    console.error('Native share failed:', err);
                }
            }
        }

        // Fallback to clipboard
        try {
            await navigator.clipboard.writeText(url);
            copiedScenarioId = scenario.id;

            // Clear any existing timeout
            if (copiedTimeout) clearTimeout(copiedTimeout);

            // Hide message after 2 seconds
            copiedTimeout = setTimeout(() => {
                copiedScenarioId = null;
            }, 2000);
        } catch (err) {
            console.error('Failed to copy to clipboard:', err);
            // Fallback for older browsers
            fallbackCopyToClipboard(url, scenario.id);
        }
    }

    function fallbackCopyToClipboard(text: string, scenarioId: number) {
        const textArea = document.createElement('textarea');
        textArea.value = text;
        textArea.style.position = 'fixed';
        textArea.style.left = '-999999px';
        document.body.appendChild(textArea);
        textArea.select();

        try {
            document.execCommand('copy');
            copiedScenarioId = scenarioId;

            if (copiedTimeout) clearTimeout(copiedTimeout);

            copiedTimeout = setTimeout(() => {
                copiedScenarioId = null;
            }, 2000);
        } catch (err) {
            console.error('Fallback copy failed:', err);
        } finally {
            document.body.removeChild(textArea);
        }
    }

    function handleScenarioCreated(event: CustomEvent) {
        const scenario = event.detail;
        console.log("NEW SCENARIO CREATED:", scenario);
        goto(`/scenarios/${scenario.sport_short_name?.toLowerCase()}/${scenario.id}`)
    }

    function formatDate(dateString: string): string {
        return new Date(dateString).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'short',
            day: 'numeric'
        });
    }

    function getRelativeTime(dateString: string): string {
        const date = new Date(dateString);
        const now = new Date();
        const diffInMs = now.getTime() - date.getTime();
        const diffInMins = Math.floor(diffInMs / (1000 * 60));
        const diffInHours = Math.floor(diffInMs / (1000 * 60 * 60));
        const diffInDays = Math.floor(diffInMs / (1000 * 60 * 60 * 24));

        if (diffInMins < 1) return 'just now';
        if (diffInMins < 60) return `${diffInMins}m ago`;
        if (diffInHours < 24) return `${diffInHours}h ago`;
        if (diffInDays < 7) return `${diffInDays}d ago`;
        
        return formatDate(dateString);
    }

    function getScenarioSeasonLabel(scenario: Scenario): string {
        if (scenario.season_end_year) {
            return `${scenario.season_start_year}-${scenario.season_end_year}`;
        }
        return `${scenario.season_start_year}`;
    }
</script>

<svelte:head>
    <title>Profile - GameScript</title>
</svelte:head>

<CreateScenarioModal 
    bind:isOpen={showCreateModal}
    on:created={handleScenarioCreated}
/>

<div class="max-w-6xl mx-auto px-4 sm:px-0">
    <!-- Profile Header -->
    <div class="rounded-lg p-4 sm:p-6 mb-6">
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
            <div class="flex items-center gap-3 sm:gap-4">
                {#if user?.avatar_url}
                    <img 
                        src={user.avatar_url} 
                        alt={user.username}
                        class="w-16 h-16 sm:w-20 sm:h-20 rounded-full object-cover border-2 border-primary-600"
                    />
                {:else}
                    <div class="w-16 h-16 sm:w-20 sm:h-20 rounded-full bg-primary-700 flex items-center justify-center border-2 border-primary-600">
                        <span class="text-2xl sm:text-3xl font-heading font-bold text-neutral">
                            {user?.username.charAt(0).toUpperCase()}
                        </span>
                    </div>
                {/if}
                <div>
                    <h1 class="text-2xl sm:text-3xl font-heading font-bold text-neutral">
                        {user?.username}
                    </h1>
                    <p class="text-sm sm:text-base text-neutral/70">{user?.email}</p>
                    {#if user?.is_admin}
                        <span class="inline-block mt-1 px-2 py-1 bg-primary-600 text-neutral text-xs font-sans font-semibold rounded">
                            ADMIN
                        </span>
                    {/if}
                </div>
            </div>
            <button
                on:click={() => showSettings = !showSettings}
                class="w-full sm:w-auto px-4 py-2 bg-primary-700 hover:bg-primary-600 text-neutral font-sans font-semibold rounded-lg transition-colors cursor-pointer"
            >
                {showSettings ? 'Cancel' : 'Edit Profile'}
            </button>
        </div>
    </div>

    <!-- Settings Panel -->
    {#if showSettings}
        <div class="border-2 border-primary-700 rounded-lg p-4 sm:p-6 mb-6">
            <h2 class="text-xl sm:text-2xl font-heading font-bold text-neutral mb-4">Profile Settings</h2>
            
            {#if settingsError}
                <div class="mb-4 p-3 bg-red-500/20 border border-red-500 rounded text-sm sm:text-base text-red-400">
                    {settingsError}
                </div>
            {/if}
            
            {#if settingsSuccess}
                <div class="mb-4 p-3 bg-green-500/20 border border-green-500 rounded text-sm sm:text-base text-green-400">
                    {settingsSuccess}
                </div>
            {/if}

            <form on:submit|preventDefault={handleUpdateProfile} class="space-y-4">
                <!-- Username -->
                <div>
                    <label for="username" class="block text-sm font-sans font-semibold text-neutral mb-2">
                        Username
                    </label>
                    <input
                        id="username"
                        type="text"
                        bind:value={username}
                        class="w-full px-3 sm:px-4 py-2 bg-primary-900 border border-primary-700 rounded-lg text-neutral text-sm sm:text-base focus:outline-none focus:border-primary-500"
                        required
                    />
                </div>

                <!-- Email -->
                <div class="pb-2">
                    <label for="email" class="block text-sm font-sans font-semibold text-neutral mb-2">
                        Email
                    </label>
                    <input
                        id="email"
                        type="email"
                        bind:value={email}
                        class="w-full px-3 sm:px-4 py-2 bg-primary-900 border border-primary-700 rounded-lg text-neutral text-sm sm:text-base focus:outline-none focus:border-primary-500"
                        required
                    />
                </div>

                <!-- Change Password Section -->
                <div class="pt-4 border-t border-primary-700">
                    <h3 class="text-base sm:text-lg font-heading font-bold text-neutral mb-3">Change Password</h3>
                    
                    <div class="space-y-4">
                        <div>
                            <label for="current-password" class="block text-sm font-sans font-semibold text-neutral mb-2">
                                Current Password
                            </label>
                            <input
                                id="current-password"
                                type="password"
                                bind:value={currentPassword}
                                class="w-full px-3 sm:px-4 py-2 bg-primary-900 border border-primary-700 rounded-lg text-neutral text-sm sm:text-base focus:outline-none focus:border-primary-500"
                                placeholder="Leave blank to keep current password"
                            />
                        </div>

                        <div>
                            <label for="new-password" class="block text-sm font-sans font-semibold text-neutral mb-2">
                                New Password
                            </label>
                            <input
                                id="new-password"
                                type="password"
                                bind:value={newPassword}
                                class="w-full px-3 sm:px-4 py-2 bg-primary-900 border border-primary-700 rounded-lg text-neutral text-sm sm:text-base focus:outline-none focus:border-primary-500"
                                placeholder="Leave blank to keep current password"
                            />
                        </div>

                        <div>
                            <label for="confirm-password" class="block text-sm font-sans font-semibold text-neutral mb-2">
                                Confirm New Password
                            </label>
                            <input
                                id="confirm-password"
                                type="password"
                                bind:value={confirmPassword}
                                class="w-full px-3 sm:px-4 py-2 bg-primary-900 border border-primary-700 rounded-lg text-neutral text-sm sm:text-base focus:outline-none focus:border-primary-500"
                                placeholder="Leave blank to keep current password"
                            />
                        </div>
                    </div>
                </div>

                <!-- Submit Button -->
                <div class="flex flex-col sm:flex-row justify-end gap-3 pt-4">
                    <button
                        type="button"
                        on:click={() => showSettings = false}
                        class="w-full sm:w-auto px-6 py-2 bg-primary-800 hover:bg-primary-700 text-neutral font-sans font-semibold rounded-lg transition-colors cursor-pointer"
                    >
                        Cancel
                    </button>
                    <button
                        type="submit"
                        disabled={settingsSaving}
                        class="w-full sm:w-auto px-6 py-2 bg-primary-600 hover:bg-primary-500 text-neutral font-sans font-semibold rounded-lg transition-colors cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                        {settingsSaving ? 'Saving...' : 'Save Changes'}
                    </button>
                </div>
            </form>
        </div>
    {/if}

    <!-- Scenarios Section -->
    <div class="border-2 border-primary-700 rounded-lg p-4 sm:p-6">
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-6">
            <h2 class="text-xl sm:text-2xl font-heading font-bold text-neutral">My Scenarios</h2>
            <button 
                on:click={() => showCreateModal = true}
                class="w-full sm:w-auto px-4 py-2 bg-primary-600 hover:bg-primary-500 text-neutral font-sans font-semibold rounded-lg transition-colors cursor-pointer"
            >
                Create New Scenario
            </button>
        </div>

        {#if loading}
            <div class="flex items-center justify-center py-12">
                <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-400"></div>
            </div>
        {:else if error}
            <div class="text-center py-8">
                <p class="text-red-400 text-base sm:text-lg">{error}</p>
            </div>
        {:else if scenarios.length === 0}
            <div class="text-center py-12">
                <p class="text-neutral/70 text-base sm:text-lg mb-4">You haven't created any scenarios yet.</p>
                <a 
                    href="/scenarios/create"
                    class="inline-block px-6 py-3 bg-primary-600 hover:bg-primary-500 text-neutral font-sans font-semibold rounded-lg transition-colors"
                >
                    Create Your First Scenario
                </a>
            </div>
        {:else}
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {#each scenarios as scenario}
                    <div class="bg-primary-900/30 border border-primary-700 rounded-lg p-4 hover:border-primary-500 transition-colors">
                        <div class="flex items-start justify-between mb-1">
                            <div class="flex-1 min-w-0">
                                <h3 class="text-base sm:text-lg font-heading font-bold text-neutral mb-1 truncate" title={scenario.name}>
                                    {scenario.name}
                                </h3>
                                <div class="flex items-center gap-1 text-xs sm:text-sm text-neutral/70">
                                    <span class="uppercase font-sans font-semibold">{scenario.sport_short_name}</span>
                                    <span>-</span>
                                    <span>{getScenarioSeasonLabel(scenario)} Season</span>
                                </div>
                            </div>
                            {#if scenario.is_public}
                                <span class="ml-2 mt-0.5 px-2 py-1 bg-green-500/20 text-green-400 text-xs font-sans font-semibold rounded shrink-0">
                                    PUBLIC
                                </span>
                            {:else}
                                <span class="ml-2 mt-0.5 px-2 py-1 bg-gray-500/20 text-gray-400 text-xs font-sans font-semibold rounded shrink-0">
                                    PRIVATE
                                </span>
                            {/if}
                        </div>

                        <!-- Timestamps on one line with bullet separator -->
                        <div class="flex items-center gap-2 text-xs text-neutral/50 mb-4">
                            <span title={`Created ${formatDate(scenario.created_at)}`}>
                                Created {formatDate(scenario.created_at)}
                            </span>
                            <span>•</span>
                            <span 
                                class="text-primary-400 font-semibold"
                                title={`Last updated ${formatDate(scenario.updated_at)}`}
                            >
                                Updated {getRelativeTime(scenario.updated_at)}
                            </span>
                        </div>

                        <div class="flex gap-2">
                            <a
                                href="/scenarios/{scenario.sport_short_name?.toLowerCase()}/{scenario.id}"
                                class="flex-1 px-3 py-2 bg-primary-700 hover:bg-primary-600 text-neutral text-center text-sm sm:text-base font-sans font-semibold rounded transition-colors"
                            >
                                Open
                            </a>
                            <div class="relative">
                                <button
                                    on:click={() => handleShareScenario(scenario)}
                                    class="px-3 py-3 bg-primary-800/60 hover:bg-primary-700 border border-primary-600 text-neutral text-sm sm:text-base font-sans font-semibold rounded transition-colors cursor-pointer"
                                    title="Share Scenario"
                                >
                                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z" />
                                    </svg>
                                </button>

                                <!-- Copied Message Tooltip -->
                                {#if copiedScenarioId === scenario.id}
                                    <div class="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 px-3 py-2 bg-green-500 text-white text-xs font-sans font-semibold rounded-md shadow-lg whitespace-nowrap z-50 animate-fade-in">
                                        Link copied!
                                        <div class="absolute top-full left-1/2 -translate-x-1/2 -mt-1 w-2 h-2 bg-green-500 transform rotate-45"></div>
                                    </div>
                                {/if}
                            </div>
                            <button
                                on:click={() => handleDeleteScenario(scenario.id)}
                                disabled={deletingScenarioId === scenario.id}
                                class="px-3 py-2 bg-red-900/80 hover:bg-red-800 text-neutral border border-red-600 text-sm sm:text-base font-sans font-semibold rounded transition-colors cursor-pointer disabled:opacity-50"
                            >
                                {deletingScenarioId === scenario.id ? '...' : 'Delete'}
                            </button>
                        </div>
                    </div>
                {/each}
            </div>
        {/if}
    </div>
</div>

<style>
    @keyframes fade-in {
        from {
            opacity: 0;
            transform: translateY(4px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }

    .animate-fade-in {
        animation: fade-in 0.2s ease-out;
    }
</style>