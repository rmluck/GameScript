<script lang="ts">
    import { createEventDispatcher } from 'svelte';

    // Props
    export let isOpen = false;
    export let feature = '';

    // Event dispatcher
    const dispatch = createEventDispatcher();

    function closeModal() {
        isOpen = false;
        dispatch('close');
    }

    // Close modal when clicking outside the content
    function handleBackdropClick(event: MouseEvent) {
        if (event.target === event.currentTarget) {
            closeModal();
        }
    }

    // Close modal on Escape key press
    function handleKeydown(event: KeyboardEvent) {
        if (event.key === 'Escape' && isOpen) {
            closeModal();
        }
    }
</script>

{#if isOpen}
    <!-- Modal Overlay -->
    <div 
        class="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50" 
        on:click={handleBackdropClick} 
        on:keydown={handleKeydown} 
        role="dialog" 
        aria-modal="true" 
        tabindex="-1"
    >
        <div class="bg-primary-900 border-2 border-primary-700 rounded-lg shadow-xl max-w-md w-full mx-4 p-6">
            <div class="flex justify-between items-center mb-4">
                <h2 class="text-2xl font-heading font-bold text-neutral">COMING SOON</h2>
                <button 
                    on:click={closeModal} 
                    class="text-neutral hover:text-primary-400 transition-colors cursor-pointer" 
                    aria-label="Close modal"
                >
                    <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                </button>
            </div>

            <div class="text-center py-6">
                <div class="text-6xl mb-4">🚧</div>
                <p class="font-sans text-lg text-neutral mb-4">
                    {feature} is currently under development and will be available soon!
                </p>
                <p class="font-sans text-sm text-neutral/70">
                    Check back later for updates.
                </p>
            </div>

            <button
                on:click={closeModal}
                class="w-full bg-primary-600 hover:bg-primary-500 border-2 border-primary-500 hover:border-primary-400 rounded-lg py-3 font-sans font-semibold text-xl text-neutral transition-all hover:scale-105 cursor-pointer"
            >
                GOT IT
            </button>
        </div>
    </div>
{/if}