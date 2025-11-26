<script lang="ts">
    import '../app.css';
    import { authStore } from '$stores/auth';
    import { onMount } from 'svelte';
    import { authAPI } from '$lib/api/auth';

    onMount(async () => {
        // Validate token on mount
        if ($authStore.isAuthenticated) {
            try {
                const user = await authAPI.getCurrentUser();
                authStore.updateUser(user);
            } catch (error) {
                // Token is invalid
                authStore.logout();
            }
        }
    });
</script>

<div class="min-h-screen bg-gray-50">
    <nav class="bg-white shadow-sm">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div class="flex justify-between h-16">
                <div class="flex">
                    <a href="/" class="flex items-center">
                        <span class="text-2xl font-bold text-primary-600">GameScript</span>
                    </a>
                    <div class="hidden sm:ml-6 sm:flex sm:space-x-8">
                        <a
                            href="/scenarios"
                            class="inline-flex items-center px-1 pt-1 text-sm font-medium text-gray-900"
                        >
                            Scenarios
                        </a>
                        <a
                            href="/games"
                            class="inline-flex items-center px-1 pt-1 text-sm font-medium text-gray-500 hover:text-gray-900"
                        >
                            Games
                        </a>
                    </div>
                </div>
                <div class="flex items-center">
                    {#if $authStore.isAuthenticated}
                        <span class="text-sm text-gray-700 mr-4">
                            {$authStore.user?.username}
                        </span>
                        <button
                            on:click={() => authStore.logout()}
                            class="text-sm font-medium text-gray-500 hover:text-gray-900"
                        >
                            Logout
                        </button>
                    {:else}
                        <a href="/auth/login" class="text-sm font-medium text-gray-500 hover:text-gray-900 mr-4">
                            Login
                        </a>
                        <a
                            href="/auth/register"
                            class="text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 px-4 py-2 rounded-md"
                        >
                            Sign Up
                        </a>
                    {/if}
                </div>
            </div>
        </div>
    </nav>

    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <slot />
    </main>
</div>