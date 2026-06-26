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

export interface MataKuliah {
	id: string;
	name: string;
	sks: number;
	created_at: string;
}
