<script lang="ts">
    import '../app.css';
    import { authStore } from '$stores/auth';
    import { onMount } from 'svelte';
    import { authAPI } from '$lib/api/auth';
	import { goto } from '$app/navigation';
    import ComingSoonModal from "$lib/components/scenarios/ComingSoonModal.svelte";

    // State variables for mobile menu and modals
    let mobileMenuOpen = false;
    let showComingSoonModal=false;

    onMount(async () => {
        if ($authStore.isAuthenticated) {
            // Validate token by fetching current user
            try {
                const user = await authAPI.getCurrentUser();
                authStore.updateUser(user);
            } catch (error) {
                // If token is invalid or expired, log out the user
                authStore.logout();
            }
        }
    });

    function toggleMobileMenu() {
        mobileMenuOpen = !mobileMenuOpen;
    }

    function closeMobileMenu() {
        mobileMenuOpen = false;
    }

    function handleCFBClick(event: Event) {
        event.preventDefault();
        showComingSoonModal = true;
        closeMobileMenu();
    }

    function handleLogout() {
        authStore.logout();
        goto('/');
    }
</script>

<ComingSoonModal 
    bind:isOpen={showComingSoonModal} 
    feature="College Football"
/>

<div class="min-h-screen bg-linear-to-br from-primary-975 to-primary-950 flex flex-col">
    <nav class="bg-primary-900/30 shadow-md">
        <!-- Navbar -->
        <div class="px-4 sm:px-6 lg:px-8">
            <div class="flex justify-between h-16">
                <div class="flex">
                    <a href="/" class="flex items-center">
                        <span class="font-sans font-bold text-2xl sm:text-3xl bg-linear-to-r from-primary-700 via-primary-600 to-primary-500 bg-clip-text text-transparent">GameScript</span>
                    </a>
                    <div class="hidden md:ml-6 md:flex md:space-x-8">
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
                        <button
                            on:click={handleCFBClick}
                            class="inline-flex items-center px-1 pt-1 font-sans font-semibold text-lg text-neutral hover:text-primary-400 transition-colors duration-200 cursor-pointer relative"
                        >
                            CFB
                        </button>
                    </div>
                </div>

                <div class="flex items-center">
                    {#if $authStore.isAuthenticated}
                        <!-- Username -->
                        <span class="text-xs sm:text-sm text-neutral-300 truncate max-w-20 sm:max-w-none mr-2 sm:mr-4">
                            {$authStore.user?.is_admin ? 'ADMIN:' : ''}
                            {$authStore.user?.username}
                        </span>

                        <!-- Profile Button -->
                        <a href="/profile" class="hidden md:block font-sans font-semibold text-sm lg:text-base bg-primary-800/60 hover:bg-primary-600 text-neutral transition-colors duration-200 px-3 py-2 rounded-md mr-2 sm:mr-4">
                            PROFILE
                        </a>

                        <!-- Logout Button -->
                        <button
                            on:click={handleLogout}
                            class="hidden md:block font-sans font-semibold text-sm lg:text-base bg-primary-800/60 hover:bg-primary-600 text-neutral transition-colors duration-200 px-3 py-2 rounded-md"
                        >
                            LOGOUT
                        </button>
                    {:else}
                        <!-- Login Button -->
                        <a href="/auth/login" class="hidden md:block font-sans font-semibold text-sm lg:text-base bg-primary-800/60 hover:bg-primary-600 text-neutral transition-colors duration-200 px-3 py-2 rounded-md mr-2 sm:mr-4">
                            LOGIN
                        </a>

                        <!-- Sign Up Button -->
                        <a
                            href="/auth/register"
                            class="hidden md:block font-sans font-semibold text-sm lg:text-base bg-primary-800/60 hover:bg-primary-600 text-neutral transition-colors duration-200 px-3 py-2 rounded-md"
                        >
                            SIGN UP
                        </a>
                    {/if}

                    <!-- Mobile Menu Button -->
                    <div class="md:hidden">
                        <button
                            on:click={toggleMobileMenu}
                            class="inline-flex items-center justify-center p-2 rounded-md text-neutral hover:text-primary-400 hover:bg-primary-800/30 transition-colors duration-200"
                            aria-expanded={mobileMenuOpen}
                            aria-label="Main menu"
                        >
                            {#if mobileMenuOpen}
                                <!-- Close Icon -->
                                <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                                </svg>
                            {:else}
                                <!-- Hamburger Icon -->
                                <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16m-7 6h7" />
                                </svg>
                            {/if}
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <!-- Mobile Menu -->
        <div
            class="md:hidden fixed inset-y-0 right-0 w-64 bg-primary-900 shadow-xl transform transition-transform duration-300 ease-in-out z-50"
            class:translate-x-0={mobileMenuOpen}
            class:translate-x-full={!mobileMenuOpen}
        >
            <div class="flex flex-col h-full">
                <!-- Close Button -->
                <div class="flex justify-end p-4">
                    <button
                        on:click={closeMobileMenu}
                        class="p-2 rounded-md text-neutral hover:text-primary-400 hover:bg-primary-800/30 transition-colors duration-200"
                        aria-label="Close menu"
                    >
                        <!-- Close icon -->
                        <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </div>

                <!-- Page Links -->
                <div class="px-4 py-2 space-y-1">
                    <a
                        href="/nfl"
                        on:click={closeMobileMenu}
                        class="block px-3 py-2 rounded-md font-sans font-semibold text-base text-neutral hover:text-primary-400 hover:bg-primary-800/30 transition-colors duration-200"
                    >
                        NFL
                    </a>
                    <a
                        href="/nba"
                        on:click={closeMobileMenu}
                        class="block px-3 py-2 rounded-md font-sans font-semibold text-base text-neutral hover:text-primary-400 hover:bg-primary-800/30 transition-colors duration-200"
                    >
                        NBA
                    </a>
                    <button
                        on:click={handleCFBClick}
                        class="w-full text-left px-3 py-2 rounded-md font-sans font-semibold text-base text-neutral hover:text-primary-400 hover:bg-primary-800/30 transition-colors duration-200 cursor-pointer"
                    >
                        CFB
                    </button>
                </div>

                <!-- Divider -->
                <div class="border-t border-primary-700 my-4"></div>

                <!-- Auth Section -->
                <div class="px-4 py-2 space-y-2">
                    {#if $authStore.isAuthenticated}
                        <!-- Profile Button -->
                        <a
                            href="/profile"
                            on:click={closeMobileMenu}
                            class="block w-full text-center px-3 py-2 rounded-md font-sans font-semibold text-sm bg-primary-800/60 hover:bg-primary-600 text-neutral transition-colors duration-200"
                        >
                            Profile
                        </a>

                        <!-- Logout Button -->
                        <button
                            on:click={() => { authStore.logout(); closeMobileMenu(); }}
                            class="block w-full text-center px-3 py-2 rounded-md font-sans font-semibold text-sm bg-primary-800/60 hover:bg-primary-600 text-neutral transition-colors duration-200"                        
                        >
                            Logout
                        </button>
                    {:else}
                        <!-- Login Button -->
                        <a
                            href="/auth/login"
                            on:click={closeMobileMenu}
                            class="block w-full text-center px-3 py-2 rounded-md font-sans font-semibold text-sm bg-primary-800/60 hover:bg-primary-600 text-neutral transition-colors duration-200"
                        >
                            Login
                        </a>

                        <!-- Sign Up Button -->
                        <a
                            href="/auth/register"
                            on:click={closeMobileMenu}
                            class="block w-full text-center px-3 py-2 rounded-md font-sans font-semibold text-sm bg-primary-800/60 hover:bg-primary-600 text-neutral transition-colors duration-200"
                        >
                            Sign Up
                        </a>
                    {/if}
                </div>
            </div>
        </div>

        <!-- Overlay when mobile menu is open -->
        {#if mobileMenuOpen}
            <div
                class="md:hidden fixed inset-0 bg-black/50 z-40"
                on:click={closeMobileMenu}
                on:keydown={(e) => { if (e.key === 'Escape') closeMobileMenu(); }}
                role="button"
                tabindex="0"
                aria-label="Close menu overlay"
            ></div>
        {/if}
    </nav>

    <!-- Main Content -->
    <main class="flex flex-1 flex-col max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <slot />
    </main>

    <!-- Footer -->
    <footer class="bg-primary-900/30 text-center text-neutral py-6 mt-12 shadow-md">
        <p>&copy; {new Date().getFullYear()} GameScript. All rights reserved.</p>
    </footer>
</div>