<script lang="ts">
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';
    import { scenariosAPI } from '$lib/api/scenarios';
    import { apiClient } from '$lib/api/client';

    // Create NFL scenario on mount
    onMount(async () => {
        try {
            const sports = await apiClient.get('/sports').then(r => r.data);
            const nflSport = sports.find((s: any) => s.short_name === 'NFL');
            
            if (!nflSport) {
                throw new Error('NFL sport not found');
            }

            // Get active season for NFL
            const seasons = await fetch(`/api/sports/${nflSport.id}/seasons`).then(r => r.json());
            const activeSeason = seasons.find((s: any) => s.is_active);

            if (!activeSeason) {
                throw new Error('No active NFL season found');
            }

            const todayDate = formatDate();

            // Create scenario with default settings
            const scenario = await scenariosAPI.create({
                name: `NFL ${activeSeason.start_year}-${activeSeason.end_year} Season - ${todayDate}`,
                sport_id: nflSport.id,
                season_id: activeSeason.id,
                is_public: false // Default to private
            });

            // Redirect to scenario page
            goto(`/scenarios/nfl/${scenario.id}`);
        } catch (error) {
            console.error('Failed to create NFL scenario:', error);
        }
    });

    function formatDate() {
        const today = new Date();
        const month = String(today.getMonth() + 1).padStart(2, '0');
        const day = String(today.getDate()).padStart(2, '0');
        const year = today.getFullYear();
        return `${month}/${day}/${year}`;
    }
</script>

<div class="flex items-center justify-center min-h-full">
    <div class="text-center">
        <div class="animate-spin rounded-full h-16 w-16 border-b-2 border-primary-400 mx-auto mb-4"></div>
        <p class="text-neutral text-xl font-sans">Creating your NFL scenario...</p>
    </div>
</div>