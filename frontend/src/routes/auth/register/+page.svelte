<script lang="ts">
    import { authStore } from '$lib/stores/auth';
    import { authAPI } from '$lib/api/auth';
    import { goto } from '$app/navigation';
    import { validatePassword, validateEmail, validateUsername } from '$lib/utils/validation';

    // State variables for account details
    let email = '';
    let username = '';
    let password = '';
    let confirmPassword = '';

    // Loading and error states
    let errors: string[] = [];
    let loading = false;

    async function handleRegister() {
        errors = [];

        // Validate inputs
        if (!email || !username || !password || !confirmPassword) {
            errors = ['All fields are required'];
            return;
        }

        if (!validateEmail(email)) {
            errors.push('Please enter a valid email address');
        }

        const usernameValidation = validateUsername(username);
        if (!usernameValidation.isValid) {
            errors.push(...usernameValidation.errors);
        }

        const passwordValidation = validatePassword(password);
        if (!passwordValidation.isValid) {
            errors.push(...passwordValidation.errors);
        }

        if (password !== confirmPassword) {
            errors.push('Passwords do not match');
        }

        if (errors.length > 0) {
            return;
        }

        loading = true;

        try {
            const response = await authAPI.register(email, username, password);
            authStore.login(response.user, response.token);
            goto('/profile');
        } catch (err: any) {
            errors.push(err.response?.data?.error || 'Registration failed. Please try again.');
        } finally {
            loading = false;
        }
    }
</script>

<svelte:head>
    <title>Sign Up - GameScript</title>
</svelte:head>

<div class="flex flex-1 items-center justify-center min-h-full">
    <div class="min-w-xs sm:min-w-lg">
        <div class="bg-primary-900/60 border-2 border-primary-700 py-6 sm:py-8 px-6 shadow rounded-lg">
            <h2 class="text-3xl sm:text-4xl font-heading font-bold text-neutral mb-4 sm:mb-6 text-center">SIGN UP</h2>

            {#if errors.length > 0}
                <div class="mb-4 p-4 bg-red-900/50 border-2 border-red-600 rounded-md">
                    {#each errors as error}
                        <p class="text-sm text-red-200 font-sans">{error}</p>
                    {/each}
                </div>
            {/if}

            <!-- Registration Form -->
            <form on:submit|preventDefault={handleRegister} class="space-y-4 sm:space-y-6">
                <div>
                    <!-- Email Input -->
                    <label for="email" class="block text-lg font-semibold font-sans text-neutral mb-2">
                        Email
                    </label>
                    <input
                        type="email"
                        id="email"
                        bind:value={email}
                        required
                        class="mt-1 block w-full rounded-md bg-primary-800/60 border-2 border-primary-600 text-neutral placeholder-neutral/50 focus:border-primary-400 focus:ring-2 focus:ring-primary-400 px-4 py-3 font-sans transition-colors"
                        placeholder="your@email.com"
                    />
                </div>

                <!-- Username Input -->
                <div>
                    <label for="username" class="block text-lg font-semibold font-sans text-neutral mb-2">
                        Username
                    </label>
                    <input
                        type="text"
                        id="username"
                        bind:value={username}
                        required
                        minlength="3"
                        maxlength="50"
                        class="mt-1 block w-full rounded-md bg-primary-800/60 border-2 border-primary-600 text-neutral placeholder-neutral/50 focus:border-primary-400 focus:ring-2 focus:ring-primary-400 px-4 py-3 font-sans transition-colors"
                        placeholder="username"
                    />
                </div>

                <!-- Password Input -->
                <div>
                    <label for="password" class="block text-lg font-semibold font-sans text-neutral mb-2">
                        Password
                    </label>
                    <input
                        type="password"
                        id="password"
                        bind:value={password}
                        required
                        minlength="8"
                        class="mt-1 block w-full rounded-md bg-primary-800/60 border-2 border-primary-600 text-neutral placeholder-neutral/50 focus:border-primary-400 focus:ring-2 focus:ring-primary-400 px-4 py-3 font-sans transition-colors"
                        placeholder="********"
                    />
                    <p class="mt-1 text-xs text-neutral/70 font-sans">Must be at least 8 characters</p>
                </div>

                <!-- Confirm Password Input -->
                <div>
                    <label for="confirmPassword" class="block text-lg font-semibold font-sans text-neutral mb-2">
                        Confirm Password
                    </label>
                    <input
                        type="password"
                        id="confirmPassword"
                        bind:value={confirmPassword}
                        required
                        minlength="8"
                        class="mt-1 block w-full rounded-md bg-primary-800/60 border-2 border-primary-600 text-neutral placeholder-neutral/50 focus:border-primary-400 focus:ring-2 focus:ring-primary-400 px-4 py-3 font-sans transition-colors"
                        placeholder="********"
                    />
                </div>

                <!-- Submit Button -->
                <button
                    type="submit"
                    disabled={loading}
                    class="w-full bg-primary-600 hover:bg-primary-500 border-2 border-primary-500 hover:border-primary-400 rounded-lg shadow-lg transition-all hover:scale-105 py-2 sm:py-3 font-sans font-semibold text-lg sm:text-xl text-neutral disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:scale-100 justify-center focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
                >
                    {loading ? 'CREATING ACCOUNT...' : 'CREATE ACCOUNT'}
                </button>
            </form>

            <!-- Login Link -->
            <p class="mt-6 font-sans text-center text-lg text-neutral">
                Already have an account?
                <a href="/auth/login" class="font-semibold text-primary-300 hover:text-primary-200 hover:underline transition-all duration-200">
                    Login
                </a>
            </p>
        </div>
    </div>
</div>