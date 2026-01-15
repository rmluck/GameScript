<script lang="ts">
    import { createEventDispatcher } from 'svelte';

    export let title: string;
    export let message: string;
    export let confirmText: string = 'Continue';
    export let cancelText: string = 'Cancel';
    export let warningType: 'regular' | 'playoff' = 'regular';

    const dispatch = createEventDispatcher();

    function confirm() {
        dispatch('confirm');
    }

    function cancel() {
        dispatch('cancel');
    }

    function handleClickOutside(event: MouseEvent) {
        if (event.target === event.currentTarget) {
            cancel();
        }
    }
</script>

<div 
    class="fixed inset-0 bg-black/80 flex items-center justify-center z-100 p-4"
    on:click={handleClickOutside}
    on:keydown={(e) => { if (e.key === 'Escape') cancel(); }}
    role="dialog"
    aria-modal="true"
    tabindex="-1"
>
    <div
        class="bg-neutral border-4 rounded-lg max-w-md w-full shadow-2xl"
        class:border-red-600={warningType === 'playoff'}
        class:border-yellow-600={warningType === 'regular'}
        role="document"
        tabindex="-1"
    >
        <!-- Header -->
        <div
            class="border-b-4 p-6"
            class:bg-red-600={warningType === 'playoff'}
            class:border-red-600={warningType === 'playoff'}
            class:bg-yellow-600={warningType === 'regular'}
            class:border-yellow-600={warningType === 'regular'}
        >
            <div class="flex items-center gap-3">
                <!-- Warning Icon -->
                <div class="shrink-0">
                    <svg class="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                    </svg>
                </div>
                
                <h2 class="text-2xl font-heading font-bold text-white">
                    {title}
                </h2>
            </div>
        </div>

        <!-- Body -->
        <div class="p-6">
            <p class="text-lg font-sans text-black leading-relaxed">
                {message}
            </p>
        </div>

        <!-- Footer -->
        <div class="flex gap-3 p-6 pt-0">
            <button
                on:click={cancel}
                class="flex-1 px-6 py-3 rounded-lg border-2 border-primary-600 bg-transparent hover:bg-primary-800 text-black hover:text-neutral font-heading font-bold text-lg transition-colors cursor-pointer"
            >
                {cancelText}
            </button>
            <button
                on:click={confirm}
                class="flex-1 px-6 py-3 rounded-lg border-2 font-heading font-bold text-lg transition-colors cursor-pointer"
                class:bg-red-600={warningType === 'playoff'}
                class:border-red-600={warningType === 'playoff'}
                class:hover:bg-red-700={warningType === 'playoff'}
                class:bg-yellow-600={warningType === 'regular'}
                class:border-yellow-600={warningType === 'regular'}
                class:hover:bg-yellow-700={warningType === 'regular'}
                class:text-white={true}
            >
                {confirmText}
            </button>
        </div>
    </div>
</div>