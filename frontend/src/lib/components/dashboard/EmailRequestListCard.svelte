<script lang="ts">
	import { authService } from '$lib/services/auth';
	import type { EmailChangeRequest } from '$lib/types';
	import { onMount } from 'svelte';
	import Card from '$lib/components/ui/Card.svelte';

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

<Card class="email-requests-card animate-fade-in" style="animation-delay: 0.3s; padding: 24px; margin-bottom: 24px; background: linear-gradient(to bottom right, rgba(255, 255, 255, 0.9), rgba(255, 255, 255, 0.6));">
	<div class="card-header">
		<h2>Permintaan Ganti Email</h2>
		<button class="btn-refresh" onclick={fetchRequests} disabled={loading} title="Refresh">
			🔄
		</button>
	</div>

	{#if error}
		<div class="state-container state-error" style="padding: 16px; margin-bottom: 16px;">{error}</div>
	{/if}
	{#if successMsg}
		<div class="success-message" style="background: rgba(16, 185, 129, 0.1); color: #10b981; padding: 12px; border-radius: var(--radius-sm); margin-bottom: 20px; text-align: center;">{successMsg}</div>
	{/if}

	<div class="data-table-container">
		{#if loading && requests.length === 0}
			<div class="state-container">
				<div class="spinner"></div>
				<p>Memuat daftar permintaan...</p>
			</div>
		{:else if requests.length === 0}
			<div class="state-container">
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
</Card>

<style>
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

	.actions-col {
		text-align: right;
		width: 120px;
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

	.spinner {
		width: 24px;
		height: 24px;
		border: 3px solid rgba(0,0,0,0.1);
		border-top-color: var(--primary-color);
		border-radius: 50%;
		animation: spin 1s linear infinite;
		margin: 0 auto 16px;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}
</style>
