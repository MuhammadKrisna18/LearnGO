import { fetchApi } from './api';
import type { ApiResponse, ProgramStudi } from '../types';

export const programStudiService = {
	async getList(): Promise<ApiResponse<ProgramStudi[]>> {
		return await fetchApi<ApiResponse<ProgramStudi[]>>('/program-studi', {
			method: 'GET'
		});
	}
};
