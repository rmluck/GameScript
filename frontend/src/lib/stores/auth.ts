import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import type { User } from '$types';

interface AuthState {
    user: User | null;
    token: string | null;
    isAuthenticated: boolean;
}

function createAuthStore() {
    const { subscribe, set, update } = writable<AuthState>({
        user: null,
        token: null,
        isAuthenticated: false
    });

    // Initialize from localStorage with expiry check
    if (browser) {
        const storedToken = localStorage.getItem('token');
        const tokenExpiry = localStorage.getItem('token_expiry');
        const storedUser = localStorage.getItem('user');

        if (storedToken && tokenExpiry && storedUser) {
            const expiry = new Date(tokenExpiry);
            if (expiry > new Date()) {
                set({
                    user: JSON.parse(storedUser),
                    token: storedToken,
                    isAuthenticated: true
                });
            } else {
                // Token expired, clear everything
                localStorage.removeItem('token');
                localStorage.removeItem('token_expiry');
                localStorage.removeItem('user');
            }
        }
    }

    return {
        subscribe,
        login: (user: User, token: string) => {
            if (browser) {
                localStorage.setItem('token', token);
                // Set expiry for 7 days
                const expiry = new Date();
                expiry.setDate(expiry.getDate() + 7);
                localStorage.setItem('token_expiry', expiry.toISOString());
                localStorage.setItem('user', JSON.stringify(user));
            }
            set({ user, token, isAuthenticated: true });
        },
        logout: () => {
            if (browser) {
                localStorage.removeItem('token');
                localStorage.removeItem('token_expiry');
                localStorage.removeItem('user');
            }
            set({ user: null, token: null, isAuthenticated: false });
        },
        updateUser: (user: User) => {
            if (browser) {
                localStorage.setItem('user', JSON.stringify(user));
            }
            update(state => ({ ...state, user }));
        }
    };
}

export const authStore = createAuthStore();