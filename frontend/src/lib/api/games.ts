import { apiClient } from './client';
import type { Game } from '$types';

export const gamesAPI = {
    async getBySeasonAndWeek(seasonId: number, week: number): Promise<Game[]> {
        const response = await apiClient.get<Game[]>(`/seasons/${seasonId}/weeks/${week}/games`);
        return response.data;
    },

    async getBySeason(seasonId: number): Promise<Game[]> {
        const response = await apiClient.get<Game[]>(`/seasons/${seasonId}/games`);
        return response.data;
    },

    async getById(gameId: number): Promise<Game> {
        const response = await apiClient.get<Game>(`/games/${gameId}`);
        return response.data;
    }
};