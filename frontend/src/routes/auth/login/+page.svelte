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
            goto('/scenarios');
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

<div class="max-w-md mx-auto">
    <div class="bg-white py-8 px-6 shadow rounded-lg">
        <h2 class="text-3xl font-bold text-gray-900 mb-6">Login</h2>

        {#if error}
            <div class="mb-4 p-4 bg-red-50 border border-red-200 rounded-md">
                <p class="text-sm text-red-800">{error}</p>
            </div>
        {/if}

        <form on:submit|preventDefault={handleLogin} class="space-y-6">
            <div>
                <label for="email" class="block text-sm font-medium text-gray-700">Email</label>
                <input
                    type="email"
                    id="email"
                    bind:value={email}
                    required
                    class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 px-4 py-2 border"
                />
            </div>

            <div>
                <label for="password" class="block text-sm font-medium text-gray-700">Password</label>
                <input
                    type="password"
                    id="password"
                    bind:value={password}
                    required
                    class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 px-4 py-2 border"
                />
            </div>

            <button
                type="submit"
                disabled={loading}
                class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50"
            >
                {loading ? 'Logging in...' : 'Login'}
            </button>
        </form>

        <p class="mt-4 text-center text-sm text-gray-600">
            Don't have an account?
            <a href="/auth/register" class="font-medium text-primary-600 hover:text-primary-500">
                Sign up
            </a>
        </p>
    </div>
</div>