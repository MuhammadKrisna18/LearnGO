import { authState } from '$lib/stores/auth.svelte';

const API_BASE_URL = 'http://localhost:8080/api/v1';

export async function fetchApi<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
	const headers = new Headers(options.headers);
	
	if (!headers.has('Content-Type') && !(options.body instanceof FormData)) {
		headers.set('Content-Type', 'application/json');
	}

	if (authState.token) {
		headers.set('Authorization', `Bearer ${authState.token}`);
	}

	const res = await fetch(`${API_BASE_URL}${endpoint}`, {
		...options,
		headers
	});

	const data = await res.json();

	if (!res.ok) {
		if (res.status === 401) {
			authState.logout();
		}
		throw new Error(data.message || 'An error occurred during the request.');
	}

	return data;
}
