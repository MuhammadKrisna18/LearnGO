import { fetchApi } from './api';
import type { ApiResponse, Kelas, CreateKelasPayload, PengajuanKelas } from '$lib/types';

export const kelasService = {
	async getList(): Promise<ApiResponse<Kelas[]>> {
		return await fetchApi<ApiResponse<Kelas[]>>('/kelas', {
			method: 'GET'
		});
	},

	async getById(id: string): Promise<ApiResponse<Kelas>> {
		return await fetchApi<ApiResponse<Kelas>>(`/kelas/${id}`, {
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
	},

	async requestKelas(kelasId: string, mataKuliahId: string): Promise<ApiResponse<PengajuanKelas>> {
		return await fetchApi<ApiResponse<PengajuanKelas>>('/kelas/request', {
			method: 'POST',
			body: JSON.stringify({ kelas_id: kelasId, mata_kuliah_id: mataKuliahId })
		});
	},

	async getMyPengajuan(): Promise<ApiResponse<PengajuanKelas[]>> {
		return await fetchApi<ApiResponse<PengajuanKelas[]>>('/kelas/pengajuan/me', {
			method: 'GET'
		});
	},

	async getAllPengajuan(): Promise<ApiResponse<PengajuanKelas[]>> {
		return await fetchApi<ApiResponse<PengajuanKelas[]>>('/kelas/pengajuan', {
			method: 'GET'
		});
	},

	async approvePengajuan(id: string): Promise<ApiResponse<null>> {
		return await fetchApi<ApiResponse<null>>(`/kelas/pengajuan/${id}/approve`, {
			method: 'POST'
		});
	},

	async rejectPengajuan(id: string): Promise<ApiResponse<null>> {
		return await fetchApi<ApiResponse<null>>(`/kelas/pengajuan/${id}/reject`, {
			method: 'POST'
		});
	}
};
