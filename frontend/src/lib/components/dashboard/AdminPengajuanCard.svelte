<script lang="ts">
	import { onMount } from 'svelte';
	import { matakuliahService } from '$lib/services/matakuliah';
	import type { PengajuanMataKuliah } from '$lib/types';

	let requests = $state<PengajuanMataKuliah[]>([]);
	let loading = $state(true);
	let error = $state('');

	async function loadRequests() {
		loading = true;
		error = '';
		try {
			const res = await matakuliahService.getAllRequests();
			if (res.success && res.data) {
				requests = res.data;
			}
		} catch (err: any) {
			error = err.message || 'Gagal memuat pengajuan';
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		loadRequests();
	});

	async function handleApprove(id: string) {
		const randomCode = Math.floor(100000 + Math.random() * 900000).toString();
		const code = prompt(`Untuk Menerima pengajuan, masukkan kode berikut: ${randomCode}`);
		if (!code) return;
		
		if (code !== randomCode) {
			alert('Kode tidak cocok. Aksi dibatalkan.');
			return;
		}

		try {
			const res = await matakuliahService.approveRequest(id);
			if (res.success) {
				await loadRequests();
			} else {
				alert(res.message);
			}
		} catch (err: any) {
			alert(err.message || 'Gagal menyetujui pengajuan');
		}
	}

	async function handleReject(id: string) {
		const randomCode = Math.floor(100000 + Math.random() * 900000).toString();
		const code = prompt(`Untuk Menolak pengajuan, masukkan kode berikut: ${randomCode}`);
		if (!code) return;

		if (code !== randomCode) {
			alert('Kode tidak cocok. Aksi dibatalkan.');
			return;
		}

		try {
			const res = await matakuliahService.rejectRequest(id);
			if (res.success) {
				await loadRequests();
			} else {
				alert(res.message);
			}
		} catch (err: any) {
			alert(err.message || 'Gagal menolak pengajuan');
		}
	}

	function getStatusBadgeClass(status: string) {
		switch(status) {
			case 'approved': return 'badge-success';
			case 'rejected': return 'badge-error';
			default: return 'badge-warning';
		}
	}
</script>

<div class="glass-panel p-6 mt-6">
	<div class="flex justify-between items-center mb-6">
		<h2 class="text-xl font-bold">Daftar Pengajuan Mata Kuliah Dosen</h2>
		<button class="btn-secondary" onclick={loadRequests} disabled={loading}>
			Refresh
		</button>
	</div>

	{#if error}
		<div class="error-message mb-4">{error}</div>
	{/if}

	{#if loading}
		<div class="text-center py-4">Memuat data...</div>
	{:else if requests.length === 0}
		<p class="text-center text-gray-500 py-4">Belum ada pengajuan mata kuliah.</p>
	{:else}
		<div class="overflow-x-auto">
			<table class="data-table">
				<thead>
					<tr>
						<th>Dosen</th>
						<th>Mata Kuliah</th>
						<th>Status</th>
						<th>Tanggal</th>
						<th>Aksi</th>
					</tr>
				</thead>
				<tbody>
					{#each requests as req}
						<tr>
							<td>
								<div class="font-medium">{req.dosen?.name || 'Unknown'}</div>
								<div class="text-xs text-gray-500">{req.dosen?.email || ''}</div>
							</td>
							<td>
								<div class="font-medium">{req.mata_kuliah?.name || 'Unknown'}</div>
								<div class="text-xs text-gray-500">{req.mata_kuliah?.sks || '-'} SKS</div>
							</td>
							<td>
								<span class={`status-badge ${getStatusBadgeClass(req.status)}`}>
									{req.status.toUpperCase()}
								</span>
							</td>
							<td>{new Date(req.created_at).toLocaleDateString('id-ID')}</td>
							<td>
								{#if req.status === 'pending'}
									<div class="flex gap-2">
										<button class="btn-approve" onclick={() => handleApprove(req.id)}>Terima</button>
										<button class="btn-reject" onclick={() => handleReject(req.id)}>Tolak</button>
									</div>
								{:else}
									<span class="text-gray-400 text-sm">Selesai</span>
								{/if}
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>

<style>
	.glass-panel {
		background: rgba(255, 255, 255, 0.8);
		backdrop-filter: blur(12px);
		border-radius: var(--radius-lg);
		border: 1px solid rgba(255, 255, 255, 0.5);
		box-shadow: var(--shadow-md);
	}

	.error-message {
		color: var(--error-color);
		font-size: 0.875rem;
		background: rgba(239, 68, 68, 0.1);
		padding: 10px;
		border-radius: var(--radius-sm);
	}

	.data-table {
		width: 100%;
		border-collapse: separate;
		border-spacing: 0;
	}

	.data-table th {
		text-align: left;
		padding: 12px 16px;
		font-size: 0.875rem;
		font-weight: 600;
		color: var(--text-muted);
		border-bottom: 2px solid rgba(0, 0, 0, 0.05);
	}

	.data-table td {
		padding: 16px;
		vertical-align: middle;
		border-bottom: 1px solid rgba(0, 0, 0, 0.05);
	}

	.data-table tr:last-child td {
		border-bottom: none;
	}

	.status-badge {
		font-size: 0.7rem;
		padding: 4px 8px;
		border-radius: 9999px;
		font-weight: 600;
	}
	.badge-success { background: #dcfce7; color: #166534; }
	.badge-warning { background: #fef9c3; color: #854d0e; }
	.badge-error { background: #fee2e2; color: #991b1b; }

	.btn-approve {
		background: #10b981;
		color: white;
		border: none;
		padding: 6px 12px;
		border-radius: var(--radius-sm);
		font-size: 0.75rem;
		font-weight: 500;
		cursor: pointer;
	}
	.btn-approve:hover { background: #059669; }

	.btn-reject {
		background: #ef4444;
		color: white;
		border: none;
		padding: 6px 12px;
		border-radius: var(--radius-sm);
		font-size: 0.75rem;
		font-weight: 500;
		cursor: pointer;
	}
	.btn-reject:hover { background: #dc2626; }
</style>
