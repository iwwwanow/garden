import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import { me, setToken, clearToken, type User } from '$lib/api';

export const userStore = writable<User | null>(null);

export async function initAuth(): Promise<void> {
	if (!browser) return;
	const token = localStorage.getItem('token');
	if (!token) return;
	try {
		const data = await me.get();
		userStore.set(data.user);
	} catch {
		clearToken();
		userStore.set(null);
	}
}

export function login(token: string, user: User): void {
	setToken(token);
	userStore.set(user);
}

export function logout(): void {
	clearToken();
	userStore.set(null);
}
