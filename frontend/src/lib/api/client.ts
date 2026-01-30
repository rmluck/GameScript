// API client for making requests to the backend

import axios, { type AxiosInstance, type AxiosError } from 'axios';
import { browser } from '$app/environment';
import { env } from '$env/dynamic/public';

// Use environment variable for API URL
const BASE_URL = env.PUBLIC_API_URL || 'https://gamescript.onrender.com/api';

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