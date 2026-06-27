import { fetchApi } from './api';
import type { ApiResponse, Kelas, CreateKelasPayload } from '$lib/types';

export const kelasService = {
	async getList(): Promise<ApiResponse<Kelas[]>> {
		return await fetchApi<ApiResponse<Kelas[]>>('/kelas', {
			method: 'GET'
		});
	},

	async create(payload: CreateKelasPayload): Promise<ApiResponse<Kelas>> {
		return await fetchApi<ApiResponse<Kelas>>('/kelas', {
			method: 'POST',
			body: JSON.stringify(payload)
		});
	},

	async delete(id: string): Promise<ApiResponse<null>> {
		return await fetchApi<ApiResponse<null>>(`/kelas/${id}`, {
			method: 'DELETE'
		});
	}
};
