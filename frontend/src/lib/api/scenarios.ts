import { apiClient } from './client';
import type { Scenario } from '$types';

export const scenariosAPI = {
    async getAll(): Promise<Scenario[]> {
        const response = await apiClient.get<Scenario[]>('/scenarios');
        return response.data;
    },

    async getById(id: number): Promise<Scenario> {
        const response = await apiClient.get<Scenario>(`/scenarios/${id}`);
        return response.data;
    },

    async create(data: {
        name: string;
        sport_id: number;
        season_id: number;
        is_public: boolean;
    }): Promise<Scenario> {
        const response = await apiClient.post<Scenario>('/scenarios', data);
        return response.data;
    },

    async update(id: number, data: { name?: string; is_public?: boolean }): Promise<Scenario> {
        const response = await apiClient.put<Scenario>(`/scenarios/${id}`, data);
        return response.data;
    },

    async delete(id: number): Promise<void> {
        await apiClient.delete(`/scenarios/${id}`);
    }
};