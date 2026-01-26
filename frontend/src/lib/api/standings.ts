// API functions related to standings

import { apiClient } from './client';
import type { NFLStandings, NBAStandings } from '$types';

export const standingsAPI = {
    async getByNFLScenario(scenarioId: number): Promise<NFLStandings> {
        const response = await apiClient.get<NFLStandings>(`/scenarios/${scenarioId}/standings`);
        return response.data;
    },

    async getByNBAScenario(scenarioId: number): Promise<NBAStandings> {
        const response = await apiClient.get<NBAStandings>(`/scenarios/${scenarioId}/standings`);
        return response.data;
    }
};