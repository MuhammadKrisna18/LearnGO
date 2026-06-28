import { fetchApi } from './api';
import type { ApiResponse, LoginResponse, UserProfile, User, EmailChangeRequest } from '$lib/types';
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

	async getProfile(): Promise<ApiResponse<UserProfile>> {
		const data = await fetchApi<ApiResponse<UserProfile>>('/auth/me', {
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
	},

	async updateProfile(data: { name: string, nickname?: string, program_studi_id?: string }): Promise<ApiResponse<UserProfile>> {
		const res = await fetchApi<ApiResponse<UserProfile>>('/auth/me', {
			method: 'PUT',
			body: JSON.stringify(data)
		});
		
		if (res.success && res.data) {
			authState.setProfile(res.data as any);
		}
		
		return res;
	},

	async requestEmailChange(new_email: string): Promise<ApiResponse<null>> {
		const res = await fetchApi<ApiResponse<null>>('/auth/email-request', {
			method: 'POST',
			body: JSON.stringify({ new_email })
		});
		if (res.success) {
			// refresh profile to show pending email
			await this.getProfile();
		}
		return res;
	},

	async getEmailRequests(): Promise<ApiResponse<EmailChangeRequest[]>> {
		return await fetchApi<ApiResponse<EmailChangeRequest[]>>('/auth/email-request', {
			method: 'GET'
		});
	},

	async reviewEmailRequest(id: string, approve: boolean): Promise<ApiResponse<null>> {
		return await fetchApi<ApiResponse<null>>(`/auth/email-request/${id}`, {
			method: 'PUT',
			body: JSON.stringify({ approve })
		});
	}
};
