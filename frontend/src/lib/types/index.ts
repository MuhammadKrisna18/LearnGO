export interface ApiResponse<T = any> {
	success: boolean;
	message: string;
	data?: T;
}

export interface User {
	id: number;
	name: string;
	email: string;
	role: 'admin' | 'dosen';
	created_at: string;
}

export interface UserProfile {
	id: string;
	name: string;
	email: string;
	role: 'admin' | 'dosen';
	program_studi_id?: string;
	program_studi?: ProgramStudi;
	created_at: string;
}

export interface Dosen {
	id: number;
	name: string;
	email: string;
	created_at: string;
}

export interface LoginResponse {
	token: string;
	role: 'admin' | 'dosen';
}

export interface ProgramStudi {
	id: string;
	name: string;
	code: string;
}

export interface MataKuliah {
	id: string;
	name: string;
	sks: number;
	program_studi_id: string;
	program_studi?: ProgramStudi;
	created_at: string;
}
