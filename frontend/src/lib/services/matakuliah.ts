import { fetchApi } from './api';
import type { ApiResponse, MataKuliah, PengajuanMataKuliah } from '$lib/types';

export const matakuliahService = {
	async getList(): Promise<ApiResponse<MataKuliah[]>> {
		return await fetchApi<ApiResponse<MataKuliah[]>>('/matakuliah', {
			method: 'GET'
		});
	},

	async getMahasiswaMataKuliah(): Promise<ApiResponse<MataKuliah[]>> {
		return await fetchApi<ApiResponse<MataKuliah[]>>('/matakuliah/mahasiswa', {
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
	},

	async lepasMataKuliah(id: string): Promise<ApiResponse<null>> {
		return await fetchApi<ApiResponse<null>>(`/matakuliah/${id}/lepas`, {
			method: 'POST'
		});
	},

	async requestMataKuliah(mata_kuliah_id: string): Promise<ApiResponse<PengajuanMataKuliah>> {
		return await fetchApi<ApiResponse<PengajuanMataKuliah>>('/matakuliah/requests', {
			method: 'POST',
			body: JSON.stringify({ mata_kuliah_id })
		});
	},

	async getMyRequests(): Promise<ApiResponse<PengajuanMataKuliah[]>> {
		return await fetchApi<ApiResponse<PengajuanMataKuliah[]>>('/matakuliah/my-requests', {
			method: 'GET'
		});
	},

	async getAllRequests(): Promise<ApiResponse<PengajuanMataKuliah[]>> {
		return await fetchApi<ApiResponse<PengajuanMataKuliah[]>>('/matakuliah/requests', {
			method: 'GET'
		});
	},

	async approveRequest(id: string): Promise<ApiResponse> {
		return await fetchApi<ApiResponse>(`/matakuliah/requests/${id}/approve`, {
			method: 'POST',
		});
	},

	async rejectPengajuan(id: string): Promise<ApiResponse<null>> {
		return await fetchApi<ApiResponse<null>>(`/matakuliah/requests/${id}/reject`, {
			method: 'POST'
		});
	},

	async acceptOffer(id: string): Promise<ApiResponse<null>> {
		return await fetchApi<ApiResponse<null>>(`/matakuliah/requests/${id}/accept-offer`, {
			method: 'POST'
		});
	},

	async rejectOffer(id: string): Promise<ApiResponse<null>> {
		return await fetchApi<ApiResponse<null>>(`/matakuliah/requests/${id}/reject-offer`, {
			method: 'POST'
		});
	}
};
