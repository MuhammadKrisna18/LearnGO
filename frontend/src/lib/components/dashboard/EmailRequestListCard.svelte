<script lang="ts">
	import { authService } from '$lib/services/auth';
	import type { EmailChangeRequest } from '$lib/types';
	import { onMount } from 'svelte';

	let requests = $state<EmailChangeRequest[]>([]);
	let loading = $state(true);
	let error = $state('');
	let successMsg = $state('');

	async function fetchRequests() {
		loading = true;
		try {
			const res = await authService.getEmailRequests();
			if (res.success && res.data) {
				requests = res.data;
			} else {
				error = res.message || 'Gagal mengambil data permintaan email';
			}
		} catch (err: any) {
			error = 'Gagal terhubung ke server';
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		fetchRequests();
	});

	async function handleReview(id: string, approve: boolean) {
		error = '';
		successMsg = '';
		try {
			const res = await authService.reviewEmailRequest(id, approve);
			if (res.success) {
				successMsg = approve ? 'Permintaan disetujui, email telah diganti.' : 'Permintaan berhasil ditolak.';
				fetchRequests();
			} else {
				error = res.message || 'Gagal memproses permintaan';
			}
		} catch (err) {
			error = 'Terjadi kesalahan sistem';
		}
	}
</script>

<div class="email-requests-card glass-panel animate-fade-in" style="animation-delay: 0.3s;">
	<div class="card-header">
		<h2>Permintaan Ganti Email</h2>
		<button class="btn-refresh" onclick={fetchRequests} disabled={loading} title="Refresh">
			🔄
		</button>
	</div>

	{#if error}
		<div class="error-message">{error}</div>
	{/if}
	{#if successMsg}
		<div class="success-message">{successMsg}</div>
	{/if}

	<div class="table-container">
		{#if loading && requests.length === 0}
			<div class="loading-state">
				<div class="spinner"></div>
				<p>Memuat daftar permintaan...</p>
			</div>
		{:else if requests.length === 0}
			<div class="empty-state">
				<p>🎉 Tidak ada permintaan ganti email saat ini.</p>
			</div>
		{:else}
			<table class="data-table">
				<thead>
					<tr>
						<th>NAMA DOSEN</th>
						<th>NID</th>
						<th>EMAIL LAMA</th>
						<th>EMAIL BARU</th>
						<th class="actions-col">AKSI</th>
					</tr>
				</thead>
				<tbody>
					{#each requests as req}
						<tr>
							<td class="font-medium">{req.user?.name || 'Unknown'}</td>
							<td class="mono text-muted">{req.user?.nid || '-'}</td>
							<td class="text-muted">{req.user?.email || '-'}</td>
							<td class="font-medium highlight">{req.new_email}</td>
							<td class="actions-col">
								<button class="btn-icon btn-approve" onclick={() => handleReview(req.id, true)} title="Setujui">
									✅
								</button>
								<button class="btn-icon btn-reject" onclick={() => handleReview(req.id, false)} title="Tolak">
									❌
								</button>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		{/if}
	</div>
</div>

<style>
	.email-requests-card {
		padding: 24px;
		border-radius: var(--radius-lg);
		background: linear-gradient(to bottom right, rgba(255, 255, 255, 0.9), rgba(255, 255, 255, 0.6));
		margin-bottom: 24px;
	}

	.card-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 24px;
	}

	.card-header h2 {
		font-size: 1.25rem;
		color: var(--text-main);
		font-weight: 600;
	}

	.btn-refresh {
		background: none;
		border: 1px solid var(--surface-border);
		border-radius: var(--radius-md);
		padding: 8px;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-refresh:hover:not(:disabled) {
		background: var(--surface-light);
		transform: scale(1.05);
	}

	.table-container {
		overflow-x: auto;
	}

	.data-table {
		width: 100%;
		border-collapse: separate;
		border-spacing: 0;
	}

	.data-table th {
		text-align: left;
		padding: 12px 16px;
		font-size: 0.75rem;
		font-weight: 600;
		color: var(--text-muted);
		text-transform: uppercase;
		letter-spacing: 0.5px;
		border-bottom: 2px solid var(--surface-border);
	}

	.data-table td {
		padding: 16px;
		font-size: 0.95rem;
		color: var(--text-main);
		border-bottom: 1px solid var(--surface-border);
		vertical-align: middle;
	}

	.data-table tbody tr {
		transition: background-color 0.2s;
	}

	.data-table tbody tr:hover {
		background: var(--surface-light);
	}

	.actions-col {
		text-align: right;
		width: 120px;
	}

	.font-medium {
		font-weight: 500;
	}

	.text-muted {
		color: var(--text-muted);
	}

	.mono {
		font-family: monospace;
	}

	.highlight {
		color: var(--primary-color);
	}

	.btn-icon {
		background: none;
		border: 1px solid transparent;
		border-radius: var(--radius-sm);
		padding: 6px;
		cursor: pointer;
		font-size: 1.2rem;
		transition: all 0.2s;
		display: inline-flex;
		align-items: center;
		justify-content: center;
	}

	.btn-approve {
		background: rgba(16, 185, 129, 0.1);
	}
	.btn-approve:hover {
		background: rgba(16, 185, 129, 0.2);
		transform: translateY(-2px);
	}

	.btn-reject {
		background: rgba(239, 68, 68, 0.1);
		margin-left: 8px;
	}
	.btn-reject:hover {
		background: rgba(239, 68, 68, 0.2);
		transform: translateY(-2px);
	}

	.loading-state, .empty-state {
		padding: 48px 0;
		text-align: center;
		color: var(--text-muted);
	}

	.spinner {
		width: 24px;
		height: 24px;
		border: 3px solid rgba(0,0,0,0.1);
		border-top-color: var(--primary-color);
		border-radius: 50%;
		animation: spin 1s linear infinite;
		margin: 0 auto 16px;
	}

	.error-message {
		background: rgba(239, 68, 68, 0.1);
		color: var(--error-color);
		padding: 12px;
		border-radius: var(--radius-sm);
		margin-bottom: 20px;
		text-align: center;
	}

	.success-message {
		background: rgba(16, 185, 129, 0.1);
		color: var(--success-color);
		padding: 12px;
		border-radius: var(--radius-sm);
		margin-bottom: 20px;
		text-align: center;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}
</style>
