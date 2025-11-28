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

<div class="min-h-screen bg-primary-900">
    <nav class="bg-primary-800 shadow-sm">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div class="flex justify-between h-16">
                <div class="flex">
                    <a href="/" class="flex items-center">
                        <span class="font-sans font-bold text-3xl text-primary-400">GameScript</span>
                    </a>
                    <div class="ml-6 flex space-x-8">
                        <a
                            href="/nfl"
                            class="inline-flex items-center px-1 pt-1 font-sans font-semibold text-lg text-neutral-100 hover:text-neutral-200"
                        >
                            NFL
                        </a>
                        <a
                            href="/nba"
                            class="inline-flex items-center px-1 pt-1 font-sans font-semibold text-lg text-neutral-100 hover:text-neutral-200"
                        >
                            NBA
                        </a>
                        <a
                            href="/cfb"
                            class="inline-flex items-center px-1 pt-1 font-sans font-semibold text-lg text-neutral-100 hover:text-neutral-200"
                        >
                            CFB
                        </a>
                    </div>
                </div>
                <div class="flex items-center">
                    {#if $authStore.isAuthenticated}
                        <span class="text-sm text-neutral-300 mr-4">
                            {$authStore.user?.username}
                        </span>
                        <a href="/auth/profile" class="font-sans font-semibold text-lg text-neutral-100 hover:text-neutral-200 mr-4">
                            PROFILE
                        </a>
                        <button
                            on:click={() => authStore.logout()}
                            class="font-sans font-semibold text-lg text-neutral-100 bg-primary-500 hover:bg-accent-500 px-4 py-2 rounded-md"
                        >
                            LOGOUT
                        </button>
                    {:else}
                        <a href="/auth/login" class="font-sans font-semibold text-lg text-neutral-100 bg-primary-500 hover:bg-accent-500 px-4 py-2 rounded-md mr-4">
                            LOGIN
                        </a>
                        <a
                            href="/auth/register"
                            class="font-sans font-semibold text-lg text-neutral-100 bg-primary-500 hover:bg-accent-500 px-4 py-2 rounded-md"
                        >
                            SIGN UP
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