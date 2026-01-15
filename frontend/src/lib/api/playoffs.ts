import { apiClient } from './client';
import type { PlayoffState, PlayoffMatchup, PlayoffSeries } from '$types';

export const playoffsAPI = {
    async getState(scenarioId: number): Promise<{
        playoff_state: PlayoffState | null;
        can_enable: boolean;
    }> {
        const response = await apiClient.get(`/playoffs/scenarios/${scenarioId}/state`);
        return response.data;
    },

    async enable(scenarioId: number): Promise<void> {
        await apiClient.post(`/playoffs/scenarios/${scenarioId}/enable`);
    },

    async getMatchups(scenarioId: number, round: number): Promise<PlayoffMatchup[] | PlayoffSeries[]> {
        const response = await apiClient.get(`/playoffs/scenarios/${scenarioId}/rounds/${round}`);
        return response.data;
    },

    async updatePick(scenarioId: number, matchupId: number, data: {
        picked_team_id?: number | null;
        predicted_higher_seed_score?: number;
        predicted_lower_seed_score?: number;
        predicted_higher_seed_wins?: number;
        predicted_lower_seed_wins?: number;
    }): Promise<PlayoffMatchup | PlayoffSeries> {
        const response = await apiClient.put(
            `/playoffs/scenarios/${scenarioId}/matchups/${matchupId}`,
            data
        );
        return response.data;
    },

    async generateNextRound(scenarioId: number): Promise<void> {
        await apiClient.post(`/playoffs/scenarios/${scenarioId}/generate`);
    },

    async deletePick(scenarioId: number, matchupId: number): Promise<void> {
        await apiClient.delete(`/playoffs/scenarios/${scenarioId}/matchups/${matchupId}`);
    }
};