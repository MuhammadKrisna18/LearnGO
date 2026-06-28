export type ToastType = 'success' | 'error' | 'info';

export interface ToastMessage {
	id: string;
	type: ToastType;
	message: string;
}

function createToastStore() {
	let toasts = $state<ToastMessage[]>([]);

	function add(message: string, type: ToastType = 'info', duration: number = 3000) {
		const id = Math.random().toString(36).substring(2, 9);
		toasts.push({ id, type, message });
		
		setTimeout(() => {
			remove(id);
		}, duration);
	}

	function remove(id: string) {
		const index = toasts.findIndex(t => t.id === id);
		if (index !== -1) {
			toasts.splice(index, 1);
		}
	}

	return {
		get toasts() { return toasts; },
		success: (msg: string, duration?: number) => add(msg, 'success', duration),
		error: (msg: string, duration?: number) => add(msg, 'error', duration),
		info: (msg: string, duration?: number) => add(msg, 'info', duration),
		remove
	};
}

export const toast = createToastStore();
