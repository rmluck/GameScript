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

<div class="min-h-screen bg-linear-to-br from-primary-975  to-primary-950 flex flex-col">
    <nav class="bg-primary-900/30 shadow-md">
        <div class="px-4 sm:px-6 lg:px-8">
            <div class="flex justify-between h-16">
                <div class="flex">
                    <a href="/" class="flex items-center">
                        <span class="font-sans font-bold text-3xl bg-linear-to-r from-primary-700 via-primary-600 to-primary-500 bg-clip-text text-transparent">GameScript</span>
                    </a>
                    <div class="ml-6 flex space-x-8">
                        <a
                            href="/nfl"
                            class="inline-flex items-center px-1 pt-1 font-sans font-semibold text-lg text-neutral hover:text-primary-400 transition-colors duration-200"
                        >
                            NFL
                        </a>
                        <a
                            href="/nba"
                            class="inline-flex items-center px-1 pt-1 font-sans font-semibold text-lg text-neutral hover:text-primary-400 transition-colors duration-200"
                        >
                            NBA
                        </a>
                        <a
                            href="/cfb"
                            class="inline-flex items-center px-1 pt-1 font-sans font-semibold text-lg text-neutral hover:text-primary-400 transition-colors duration-200"
                        >
                            CFB
                        </a>
                    </div>
                </div>
                <div class="flex items-center">
                    {#if $authStore.isAuthenticated}
                        <span class="text-sm text-neutral-300 mr-4">
                            {$authStore.user?.is_admin ? 'ADMIN:' : ''}
                            {$authStore.user?.username}
                        </span>
                        <a href="/profile" class="font-sans font-semibold text-lg bg-primary-800/60 hover:bg-primary-600 text-neutral transition-colors duration-200 px-4 py-2 rounded-md mr-4">
                            PROFILE
                        </a>
                        <button
                            on:click={() => authStore.logout()}
                            class="font-sans font-semibold text-lg bg-primary-800/60 hover:bg-primary-600 text-neutral transition-colors duration-200 px-4 py-2 rounded-md"
                        >
                            LOGOUT
                        </button>
                    {:else}
                        <a href="/auth/login" class="font-sans font-semibold text-lg bg-primary-800/60 hover:bg-primary-600 text-neutral transition-colors duration-200 px-4 py-2 rounded-md mr-4">
                            LOGIN
                        </a>
                        <a
                            href="/auth/register"
                            class="font-sans font-semibold text-lg bg-primary-800/60 hover:bg-primary-600 text-neutral transition-colors duration-200 px-4 py-2 rounded-md"
                        >
                            SIGN UP
                        </a>
                    {/if}
                </div>
            </div>
        </div>
    </nav>

    <main class="flex flex-1 flex-col max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <slot />
    </main>

    <footer class="bg-primary-900/30 text-center text-neutral py-6 mt-12 shadow-md">
        <p>&copy; {new Date().getFullYear()} GameScript. All rights reserved.</p>
    </footer>
</div>