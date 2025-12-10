import { apiClient } from './client';
import type { Standings } from '$types';

export const standingsAPI = {
    async getByScenario(scenarioId: number): Promise<Standings> {
        const response = await apiClient.get<Standings>(`/scenarios/${scenarioId}/standings`);
        return response.data;
    }
};