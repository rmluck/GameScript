// API functions related to user authentication

import { apiClient } from './client';
import type { User, AuthResponse } from '$types';

export const authAPI = {
    async register(email: string, username: string, password: string): Promise<AuthResponse> {
        const response = await apiClient.post<AuthResponse>('/auth/register', {
            email,
            username,
            password
        });
        return response.data;
    },

    async login(email: string, password: string): Promise<AuthResponse> {
        const response = await apiClient.post<AuthResponse>('/auth/login', {
            email,
            password
        });
        return response.data;
    },

    async getCurrentUser(): Promise<User> {
        const response = await apiClient.get<User>('/auth/me');
        return response.data;
    },

    async updateProfile(data: {
        username?: string;
        email?: string;
        current_password?: string;
        new_password?: string;
    }): Promise<User> {
        const response = await apiClient.put<User>('/auth/profile', data);
        return response.data;
    },

    async logout(): Promise<void> {
        // Could call backend logout endpoint if you add one
        // For now, just clear local state
    }
};