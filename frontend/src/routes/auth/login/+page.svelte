<script lang="ts">
    import { authStore } from '$stores/auth';
    import { authAPI } from '$api/auth';
    import { goto } from '$app/navigation';

    let email = '';
    let password = '';
    let error = '';
    let loading = false;

    async function handleLogin() {
        error = '';
        loading = true;

        try {
            const response = await authAPI.login(email, password);
            authStore.login(response.user, response.token);
            goto('/profile');
        } catch (err: any) {
            error = err.response?.data?.error || 'Login failed. Please try again.';
        } finally {
            loading = false;
        }
    }
</script>

<svelte:head>
    <title>Login - GameScript</title>
</svelte:head>

<div class="flex flex-1 items-center justify-center min-h-full">
    <div class="min-w-xs sm:min-w-lg">
        <div class="bg-primary-900/60 border-2 border-primary-700 py-8 sm:py-12 px-6 shadow rounded-lg">
            <h2 class="text-3xl sm:text-4xl font-heading font-bold text-neutral mb-4 sm:mb-6 text-center">LOGIN</h2>

            {#if error}
                <div class="mb-4 p-4 bg-red-900/50 border-2 border-red-600 rounded-md">
                    <p class="text-sm text-red-200 font-sans">{error}</p>
                </div>
            {/if}

            <form on:submit|preventDefault={handleLogin} class="space-y-4 sm:space-y-6">
                <div>
                    <label for="email" class="block text-lg font-semibold font-sans text-neutral mb-2">Email</label>
                    <input
                        type="email"
                        id="email"
                        bind:value={email}
                        required
                        class="mt-1 block w-full rounded-md bg-primary-800/60 border-2 border-primary-600 text-neutral placeholder-neutral/50 focus:border-primary-400 focus:ring-2 focus:ring-primary-400 px-4 py-3 font-sans transition-colors"
                        placeholder="your@email.com"
                    />
                </div>

                <div>
                    <label for="password" class="block text-lg font-semibold font-sans text-neutral mb-2">Password</label>
                    <input
                        type="password"
                        id="password"
                        bind:value={password}
                        required
                        class="mt-1 block w-full rounded-md bg-primary-800/60 border-2 border-primary-600 text-neutral placeholder-neutral/50 focus:border-primary-400 focus:ring-2 focus:ring-primary-400 px-4 py-3 font-sans transition-colors"
                        placeholder="********"
                    />
                </div>

                <button
                    type="submit"
                    disabled={loading}
                    class="w-full bg-primary-600 hover:bg-primary-500 border-2 border-primary-500 hover:border-primary-400 rounded-lg shadow-lg transition-all hover:scale-105 py-2 sm:py-3 font-sans font-semibold text-lg sm:text-xl text-neutral disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:scale-100 justify-center focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 cursor-pointer"
                >
                    {loading ? 'LOGGING IN...' : 'LOGIN'}
                </button>
            </form>

            <p class="mt-6 font-sans text-center text-lg text-neutral">
                Don't have an account?
                <a href="/auth/register" class="font-semibold text-primary-300 hover:text-primary-200 hover:underline transition-all duration-200">
                    Sign up
                </a>
            </p>
        </div>
    </div>
</div>