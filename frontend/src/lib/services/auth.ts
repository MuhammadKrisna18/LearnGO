import { fetchApi } from './api';
import type { ApiResponse, LoginResponse, User } from '$lib/types';
import { authState } from '$lib/stores/auth.svelte';

export const authService = {
	async login(email: string, password: string): Promise<ApiResponse<LoginResponse>> {
		const data = await fetchApi<ApiResponse<LoginResponse>>('/auth/login', {
			method: 'POST',
			body: JSON.stringify({ email, password })
		});
		
		if (data.success && data.data) {
			authState.setAuth(data.data.token, data.data.role);
		}
		
		return data;
	},

	async getProfile(): Promise<ApiResponse<User>> {
		const data = await fetchApi<ApiResponse<User>>('/auth/me', {
			method: 'GET'
		});
		
		if (data.success && data.data) {
			authState.setProfile(data.data);
		}
		
		return data;
	},

	async getDosenList(): Promise<ApiResponse<UserProfile[]>> {
		return await fetchApi<ApiResponse<UserProfile[]>>('/auth/dosen', {
			method: 'GET'
		});
	},

	async deleteDosen(id: string): Promise<ApiResponse<null>> {
		return await fetchApi<ApiResponse<null>>(`/auth/dosen/${id}`, {
			method: 'DELETE'
		});
	}
};
