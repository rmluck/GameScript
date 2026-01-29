// API client for making requests to the backend

import axios, { type AxiosInstance, type AxiosError } from 'axios';
import { browser } from '$app/environment';

// Use environment variable for API URL
const BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

class APIClient {
    private client: AxiosInstance;

    constructor() {
        // Create Axios instance
        this.client = axios.create({
            baseURL: BASE_URL,
            headers: {
                'Content-Type': 'application/json'
            },
            withCredentials: true
        });

        // Request interceptor to add auth token
        this.client.interceptors.request.use(
            (config) => {
                if (browser) {
                    const token = localStorage.getItem('token');
                    if (token) {
                        config.headers.Authorization = `Bearer ${token}`;
                    }
                }
                return config;
            },
            (error) => {
                return Promise.reject(error);
            }
        );

        // Response interceptor for error handling
        this.client.interceptors.response.use(
            (response) => response,
            (error: AxiosError) => {
                if (error.response?.status === 401 && browser) {
                    // Unauthorized - clear token and redirect to login
                    localStorage.removeItem('token');
                    localStorage.removeItem('user');
                    window.location.href = '/auth/login';
                }
                return Promise.reject(error);
            }
        );
    }

    public getClient(): AxiosInstance {
        return this.client;
    }
}

export const apiClient = new APIClient().getClient();