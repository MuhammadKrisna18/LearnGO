<script lang="ts">
	import { onMount } from 'svelte';
	import { kelasService } from '$lib/services/kelas';
	import type { PengajuanKelas } from '$lib/types';
	import Card from '$lib/components/ui/Card.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';

	let requests = $state<PengajuanKelas[]>([]);
	let loading = $state(true);
	let error = $state('');

	let filterStatus = $state('all');

	let filteredRequests = $derived.by(() => {
		if (filterStatus === 'all') return requests;
		return requests.filter(r => r.status === filterStatus);
	});

	let stats = $derived.by(() => ({
		total: requests.length,
		pending: requests.filter(r => r.status === 'pending').length,
		approved: requests.filter(r => r.status === 'approved').length,
		rejected: requests.filter(r => r.status === 'rejected').length,
	}));

	async function loadRequests() {
		loading = true;
		error = '';
		try {
			const res = await kelasService.getMyPengajuan();
			if (res.success && res.data) {
				requests = res.data;
			}
		} catch (err: any) {
			error = err.message || 'Gagal memuat pengajuan kelas';
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		loadRequests();
	});

	function getBadgeType(status: string) {
		switch(status) {
			case 'approved': return 'success';
			case 'rejected': return 'error';
			default: return 'warning';
		}
	}

	function getStatusLabel(status: string) {
		switch(status) {
			case 'approved': return 'Disetujui';
			case 'rejected': return 'Ditolak';
			case 'pending': return 'Menunggu';
			default: return status;
		}
	}
</script>

<Card style="padding: 24px; margin-top: 24px;">
	<div class="card-header-flex">
		<div>
			<h3>Riwayat Pengajuan Kelas</h3>
			<p class="header-subtitle">Pantau status permintaan kelas Anda</p>
		</div>
		<button class="btn-secondary" onclick={loadRequests} disabled={loading}>
			Refresh
		</button>
	</div>

	{#if !loading && requests.length > 0}
		<div class="stats-row">
			<button class="stat-chip" class:active={filterStatus === 'all'} onclick={() => filterStatus = 'all'}>
				<span class="stat-count">{stats.total}</span>
				<span class="stat-label">Semua</span>
			</button>
			<button class="stat-chip chip-warning" class:active={filterStatus === 'pending'} onclick={() => filterStatus = 'pending'}>
				<span class="stat-count">{stats.pending}</span>
				<span class="stat-label">Menunggu</span>
			</button>
			<button class="stat-chip chip-success" class:active={filterStatus === 'approved'} onclick={() => filterStatus = 'approved'}>
				<span class="stat-count">{stats.approved}</span>
				<span class="stat-label">Disetujui</span>
			</button>
			<button class="stat-chip chip-error" class:active={filterStatus === 'rejected'} onclick={() => filterStatus = 'rejected'}>
				<span class="stat-count">{stats.rejected}</span>
				<span class="stat-label">Ditolak</span>
			</button>
		</div>
	{/if}

	{#if error}
		<div class="state-container state-error" style="padding: 16px; margin-bottom: 16px;">{error}</div>
	{/if}

	{#if loading}
		<div class="state-container">Memuat data...</div>
	{:else if requests.length === 0}
		<div class="empty-state">
			<div class="empty-icon">📋</div>
			<p>Anda belum pernah mengajukan kelas.</p>
			<p class="empty-hint">Kunjungi halaman <strong>Daftar Kelas</strong> untuk mengajukan kelas.</p>
		</div>
	{:else if filteredRequests.length === 0}
		<div class="state-container">Tidak ada pengajuan dengan status ini.</div>
	{:else}
		<div class="request-list">
			{#each filteredRequests as req}
				<div class="request-card" class:card-approved={req.status === 'approved'} class:card-rejected={req.status === 'rejected'} class:card-pending={req.status === 'pending'}>
					<div class="request-main">
						<div class="request-info">
							<div class="kelas-name">{req.kelas?.name || 'Unknown'}</div>
							<div class="kelas-detail">
								<span class="detail-chip">
									📅 {req.kelas?.hari || '-'}
								</span>
								<span class="detail-chip">
									🕐 {req.kelas?.jam_mulai || ''} - {req.kelas?.jam_selesai || ''}
								</span>
								{#if req.kelas?.program_studi}
									<span class="detail-chip prodi-chip">
										🏫 {req.kelas.program_studi.name}
									</span>
								{/if}
							</div>
							{#if req.mata_kuliah}
								<div class="mk-info">
									📚 Mata Kuliah: <strong>{req.mata_kuliah.name}</strong>
								</div>
							{/if}
						</div>
						<div class="request-meta">
							<Badge type={getBadgeType(req.status)}>
								{getStatusLabel(req.status)}
							</Badge>
							<div class="date-text">{new Date(req.created_at).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })}</div>
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}
</Card>

<style>
	.card-header-flex {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		margin-bottom: 24px;
	}

	.card-header-flex h3 {
		font-size: 1.25rem;
		font-weight: 600;
		color: var(--text-main);
		margin: 0 0 4px 0;
	}

	.header-subtitle {
		font-size: 0.85rem;
		color: var(--text-muted);
		margin: 0;
	}

	.stats-row {
		display: flex;
		gap: 10px;
		margin-bottom: 20px;
		flex-wrap: wrap;
	}

	.stat-chip {
		display: flex;
		align-items: center;
		gap: 8px;
		padding: 8px 16px;
		border-radius: var(--radius-full);
		border: 1px solid rgba(0,0,0,0.08);
		background: rgba(0,0,0,0.02);
		cursor: pointer;
		transition: all 0.2s;
		font-family: inherit;
	}

	.stat-chip:hover {
		background: rgba(0,0,0,0.05);
	}

	.stat-chip.active {
		background: var(--primary-color);
		color: white;
		border-color: var(--primary-color);
	}

	.stat-chip.active .stat-count,
	.stat-chip.active .stat-label {
		color: white;
	}

	.chip-warning.active {
		background: #f59e0b;
		border-color: #f59e0b;
	}

	.chip-success.active {
		background: #10b981;
		border-color: #10b981;
	}

	.chip-error.active {
		background: #ef4444;
		border-color: #ef4444;
	}

	.stat-count {
		font-weight: 700;
		font-size: 1rem;
		color: var(--text-main);
	}

	.stat-label {
		font-size: 0.8rem;
		color: var(--text-muted);
	}

	.request-list {
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.request-card {
		padding: 16px 20px;
		border-radius: var(--radius-md);
		border: 1px solid rgba(0,0,0,0.06);
		background: rgba(255,255,255,0.8);
		transition: transform 0.15s, box-shadow 0.15s;
		border-left: 4px solid transparent;
	}

	.request-card:hover {
		transform: translateY(-1px);
		box-shadow: var(--shadow-sm);
	}

	.card-pending {
		border-left-color: #f59e0b;
		background: linear-gradient(135deg, rgba(245, 158, 11, 0.03), transparent);
	}

	.card-approved {
		border-left-color: #10b981;
		background: linear-gradient(135deg, rgba(16, 185, 129, 0.03), transparent);
	}

	.card-rejected {
		border-left-color: #ef4444;
		background: linear-gradient(135deg, rgba(239, 68, 68, 0.03), transparent);
	}

	.request-main {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		gap: 16px;
	}

	.request-info {
		flex: 1;
	}

	.kelas-name {
		font-weight: 600;
		font-size: 1.05rem;
		color: var(--text-main);
		margin-bottom: 8px;
	}

	.kelas-detail {
		display: flex;
		gap: 8px;
		flex-wrap: wrap;
		margin-bottom: 8px;
	}

	.detail-chip {
		font-size: 0.78rem;
		padding: 3px 10px;
		border-radius: var(--radius-full);
		background: rgba(0,0,0,0.04);
		color: var(--text-muted);
	}

	.prodi-chip {
		background: rgba(79, 70, 229, 0.08);
		color: var(--primary-color);
	}

	.mk-info {
		font-size: 0.85rem;
		color: var(--text-muted);
	}

	.mk-info strong {
		color: var(--primary-color);
	}

	.request-meta {
		display: flex;
		flex-direction: column;
		align-items: flex-end;
		gap: 8px;
		flex-shrink: 0;
	}

	.date-text {
		font-size: 0.75rem;
		color: var(--text-muted);
	}

	.empty-state {
		text-align: center;
		padding: 40px 20px;
	}

	.empty-icon {
		font-size: 3rem;
		margin-bottom: 12px;
	}

	.empty-state p {
		color: var(--text-muted);
		margin: 0 0 4px 0;
	}

	.empty-hint {
		font-size: 0.85rem;
	}

	@media (max-width: 640px) {
		.request-main {
			flex-direction: column;
		}

		.request-meta {
			flex-direction: row;
			align-items: center;
			justify-content: space-between;
			width: 100%;
		}

		.stats-row {
			gap: 6px;
		}

		.stat-chip {
			padding: 6px 12px;
		}
	}
</style>
