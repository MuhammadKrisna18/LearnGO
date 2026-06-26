import { fetchApi } from './api';
import type { ApiResponse, MataKuliah } from '$lib/types';

export const matakuliahService = {
	async getList(): Promise<ApiResponse<MataKuliah[]>> {
		return await fetchApi<ApiResponse<MataKuliah[]>>('/matakuliah', {
			method: 'GET'
		});
	},

	async create(name: string, sks: number, program_studi_id: string): Promise<ApiResponse<MataKuliah>> {
		return await fetchApi<ApiResponse<MataKuliah>>('/matakuliah', {
			method: 'POST',
			body: JSON.stringify({ name, sks, program_studi_id })
		});
	},

	async delete(id: string): Promise<ApiResponse<null>> {
		return await fetchApi<ApiResponse<null>>(`/matakuliah/${id}`, {
			method: 'DELETE'
		});
	}
};
