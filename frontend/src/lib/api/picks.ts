// API functions related to picks

import { apiClient } from './client';
import type { Pick } from '$types';

export const picksAPI = {
    async getByScenario(scenarioId: number): Promise<Pick[]> {
        const response = await apiClient.get<Pick[]>(`/picks/scenarios/${scenarioId}`);
        return response.data;
    },

    async getByGame(scenarioId: number, gameId: number): Promise<Pick> {
        const response = await apiClient.get<Pick>(`/picks/scenarios/${scenarioId}/games/${gameId}`);
        return response.data;
    },

    async create(scenarioId: number, gameId: number, data: {
        picked_team_id?: number | null;
        predicted_home_score?: number;
        predicted_away_score?: number;
    }): Promise<Pick> {
        const response = await apiClient.post<Pick>(
            `/picks/scenarios/${scenarioId}/games/${gameId}`,
            {
                picked_team_id: data.picked_team_id === undefined ? null : data.picked_team_id,
                predicted_home_score: data.predicted_home_score,
                predicted_away_score: data.predicted_away_score
            }
        );
        return response.data;
    },

    async update(scenarioId: number, gameId: number, data: {
        picked_team_id?: number | null;
        predicted_home_score?: number;
        predicted_away_score?: number;
    }): Promise<Pick> {
        const response = await apiClient.put<Pick>(
            `/picks/scenarios/${scenarioId}/games/${gameId}`,
            {
                picked_team_id: data.picked_team_id === undefined ? null : data.picked_team_id,
                predicted_home_score: data.predicted_home_score,
                predicted_away_score: data.predicted_away_score
            }
        );
        return response.data;
    },

    async delete(scenarioId: number, gameId: number): Promise<void> {
        await apiClient.delete(`/picks/scenarios/${scenarioId}/games/${gameId}`);
    }
};