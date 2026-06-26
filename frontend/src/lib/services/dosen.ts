import { fetchApi } from './api';
import type { ApiResponse, Dosen } from '$lib/types';

export const dosenService = {
	async getList(): Promise<ApiResponse<Dosen[]>> {
		return await fetchApi<ApiResponse<Dosen[]>>('/auth/dosen', {
			method: 'GET'
		});
	},

	async register(name: string, username: string, password: string, program_studi_id: string): Promise<ApiResponse<UserProfile>> {
		return await fetchApi<ApiResponse<UserProfile>>('/auth/dosen', {
			method: 'POST',
			body: JSON.stringify({ name, username, password, program_studi_id })
		});
	},
};
