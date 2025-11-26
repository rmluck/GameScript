import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import type { User } from '$types';

interface AuthState {
    user: User | null;
    token: string | null;
    isAuthenticated: boolean;
}

function createAuthStore() {
    const initialState: AuthState = {
        user: null,
        token: null,
        isAuthenticated: false
    };

    // Load from localStorage on init
    if (browser) {
        const storedToken = localStorage.getItem('token');
        const storedUser = localStorage.getItem('user');
        if (storedToken && storedUser) {
            initialState.token = storedToken;
            initialState.user = JSON.parse(storedUser);
            initialState.isAuthenticated = true;
        }
    }

    const { subscribe, set, update } = writable<AuthState>(initialState);

    return {
        subscribe,
        login: (user: User, token: string) => {
            if (browser) {
                localStorage.setItem('token', token);
                localStorage.setItem('user', JSON.stringify(user));
            }
            set({ user, token, isAuthenticated: true });
        },
        logout: () => {
            if (browser) {
                localStorage.removeItem('token');
                localStorage.removeItem('user');
            }
            set({ user: null, token: null, isAuthenticated: false });
        },
        updateUser: (user: User) => {
            if (browser) {
                localStorage.setItem('user', JSON.stringify(user));
            }
            update((state) => ({ ...state, user }));
        }
    };
}

export const authStore = createAuthStore();