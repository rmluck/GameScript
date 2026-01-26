// API functions related to teams

import { apiClient } from './client';
import type { Team } from '$types';

export const teamsAPI = {
    async getBySeason(seasonId: number): Promise<Team[]> {
        const response = await apiClient.get<Team[]>(`/seasons/${seasonId}/teams`);
        return response.data;
    },

    async getById(teamId: number): Promise<Team> {
        const response = await apiClient.get<Team>(`/teams/${teamId}`);
        return response.data;
    }
};