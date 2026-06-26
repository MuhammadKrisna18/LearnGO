import type { User } from '$lib/types';

function createAuthStore() {
	// Only run localStorage logic if we're in the browser
	const isBrowser = typeof window !== 'undefined';
	
	let token = $state(isBrowser ? localStorage.getItem('token') : null);
	let role = $state(isBrowser ? localStorage.getItem('role') : null);
	let profile = $state<User | null>(null);

	function setAuth(newToken: string, newRole: string) {
		token = newToken;
		role = newRole;
		if (isBrowser) {
			localStorage.setItem('token', newToken);
			localStorage.setItem('role', newRole);
		}
	}

	function setProfile(newProfile: User) {
		profile = newProfile;
	}

	function logout() {
		token = null;
		role = null;
		profile = null;
		if (isBrowser) {
			localStorage.removeItem('token');
			localStorage.removeItem('role');
			window.location.href = '/login';
		}
	}

	return {
		get token() { return token; },
		get role() { return role; },
		get profile() { return profile; },
		setAuth,
		setProfile,
		logout
	};
}

export const authState = createAuthStore();
