<script lang="ts">
	import { onMount } from 'svelte';
	import { matakuliahService } from '$lib/services/matakuliah';
	import { authState } from '$lib/stores/auth.svelte';
	import type { MataKuliah, PengajuanMataKuliah } from '$lib/types';

	let mkList = $state<MataKuliah[]>([]);
	let myRequests = $state<PengajuanMataKuliah[]>([]);
	let loading = $state(true);
	let error = $state('');

	async function loadData() {
		loading = true;
		error = '';
		try {
			const [mkRes, reqRes] = await Promise.all([
				matakuliahService.getList(),
				matakuliahService.getMyRequests()
			]);
			
			if (mkRes.success && mkRes.data) {
				mkList = mkRes.data;
			}
			if (reqRes.success && reqRes.data) {
				myRequests = reqRes.data;
			}
		} catch (err: any) {
			error = err.message || 'Gagal memuat data';
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		loadData();
	});

	async function requestMk(mkId: string) {
		const randomCode = Math.floor(100000 + Math.random() * 900000).toString();
		const userInput = prompt(`Untuk melanjutkan pengajuan, masukkan kode berikut: ${randomCode}`);
		if (!userInput) return;
		
		if (userInput !== randomCode) {
			alert('Kode tidak cocok. Pengajuan dibatalkan.');
			return;
		}

		try {
			const res = await matakuliahService.requestMataKuliah(mkId);
			if (res.success) {
				await loadData();
			} else {
				alert(res.message);
			}
		} catch (err: any) {
			alert(err.message || 'Gagal mengajukan mata kuliah');
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

<div class="glass-panel p-6" style="margin-top: 24px;">
	<div class="header-flex">
		<h2 style="margin-bottom: 0;">Pengajuan Mata Kuliah</h2>
		<button class="btn-secondary" onclick={loadData} disabled={loading}>
			Refresh
		</button>
	</div>

	{#if error}
		<div class="error-message">{error}</div>
	{/if}

	{#if loading}
		<div style="text-align: center; padding: 20px;">Memuat data...</div>
	{:else}
		<!-- Mata Kuliah yang Sudah Diambil -->
		<div class="mk-diampu-section">
			<h3>Mata Kuliah Yang Diampu</h3>
			{#if myRequests.filter(req => req.status === 'approved').length === 0}
				<p class="empty-text">Anda belum mengampu mata kuliah apapun.</p>
			{:else}
				<div class="grid-container diampu-grid">
					{#each myRequests.filter(req => req.status === 'approved') as req}
						<div class="mk-card">
							<div class="mk-card-icon">📚</div>
							<div class="mk-card-content">
								<div class="mk-title">{req.mata_kuliah?.name || 'Mata Kuliah Dihapus'}</div>
								<div class="mk-subtitle">{req.mata_kuliah?.sks || 0} SKS - {req.mata_kuliah?.program_studi?.name || 'Unknown'}</div>
								<div class="mk-date">Disetujui: {new Date(req.created_at).toLocaleDateString()}</div>
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</div>

		<hr class="section-divider" />

		<div class="grid-container">
			<!-- Available MKs -->
			<div class="grid-column">
				<h3>Mata Kuliah Tersedia</h3>
				{#if mkList.length === 0}
					<p class="empty-text">Belum ada mata kuliah.</p>
				{:else}
					<ul class="list-container">
						{#each mkList as mk}
							<li class="list-item">
								<div>
									<div class="item-title">{mk.name}</div>
									<div class="item-subtitle">{mk.sks} SKS - {mk.program_studi?.name || 'Unknown'}</div>
								</div>
								{#if mk.pengajuan && mk.pengajuan.some(p => p.status === 'pending' || p.status === 'approved')}
									{#if mk.pengajuan.some(p => (p.status === 'pending' || p.status === 'approved') && p.dosen_id === authState.profile?.id)}
										<span class="status-text text-warning">Sudah Anda Ajukan</span>
									{:else}
										<span class="status-text text-muted">Diambil Dosen Lain</span>
									{/if}
								{:else}
									<button class="btn-primary-sm" onclick={() => requestMk(mk.id)}>
										Ajukan
									</button>
								{/if}
							</li>
						{/each}
					</ul>
				{/if}
			</div>

			<!-- My Requests -->
			<div class="grid-column">
				<h3>Riwayat Pengajuan</h3>
				{#if myRequests.length === 0}
					<p class="empty-text">Belum ada pengajuan.</p>
				{:else}
					<ul class="list-container">
						{#each myRequests as req}
							<li class="list-item flex-col">
								<div class="req-header" style="margin-bottom: 0;">
									<div class="item-title">{req.mata_kuliah?.name || 'Mata Kuliah Dihapus'}</div>
									<span class={`status-badge ${getStatusBadgeClass(req.status)}`}>
										{req.status.toUpperCase()}
									</span>
								</div>
							</li>
						{/each}
					</ul>
				{/if}
			</div>
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

	.header-flex {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 24px;
	}

	.error-message {
		color: var(--error-color);
		font-size: 0.875rem;
		background: rgba(239, 68, 68, 0.1);
		padding: 10px;
		border-radius: var(--radius-sm);
		margin-bottom: 16px;
	}

	.grid-container {
		display: grid;
		gap: 24px;
		grid-template-columns: 1fr;
	}
	@media (min-width: 768px) {
		.grid-container {
			grid-template-columns: 1fr 1fr;
		}
		.diampu-grid {
			grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
		}
	}

	h3 {
		font-size: 1.1rem;
		margin-bottom: 16px;
		color: var(--text-main);
	}

	.list-container {
		list-style: none;
		padding: 0;
		margin: 0;
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.list-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 12px;
		background: rgba(0, 0, 0, 0.02);
		border-radius: var(--radius-md);
		border: 1px solid rgba(0, 0, 0, 0.05);
	}

	.flex-col {
		flex-direction: column;
		align-items: stretch;
	}

	.item-title {
		font-weight: 500;
		color: var(--text-main);
	}

	.item-subtitle {
		font-size: 0.75rem;
		color: var(--text-muted);
		margin-top: 4px;
	}

	.btn-primary-sm {
		background: var(--primary-color);
		color: white;
		border: none;
		padding: 6px 12px;
		border-radius: var(--radius-sm);
		font-size: 0.75rem;
		font-weight: 500;
		cursor: pointer;
		transition: background 0.2s;
	}
	.btn-primary-sm:hover {
		background: var(--primary-hover);
	}

	.status-text {
		font-size: 0.75rem;
		font-weight: 600;
		padding: 4px 8px;
		border-radius: var(--radius-sm);
		background: rgba(0,0,0,0.05);
	}
	.text-warning { color: #854d0e; background: #fef9c3; }
	.text-muted { color: var(--text-muted); }

	.req-header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		margin-bottom: 12px;
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

	.mk-diampu-section {
		margin-bottom: 24px;
	}
	.section-divider {
		border: 0;
		border-top: 1px solid rgba(0,0,0,0.05);
		margin: 32px 0;
	}
	.mk-card {
		background: rgba(255,255,255,0.9);
		border: 1px solid var(--surface-border);
		border-radius: var(--radius-md);
		padding: 16px;
		display: flex;
		gap: 16px;
		align-items: center;
		transition: transform 0.2s, box-shadow 0.2s;
	}
	.mk-card:hover {
		transform: translateY(-2px);
		box-shadow: var(--shadow-md);
	}
	.mk-card-icon {
		font-size: 2rem;
		background: rgba(79, 70, 229, 0.1);
		width: 56px;
		height: 56px;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: var(--radius-md);
	}
	.mk-title {
		font-weight: 600;
		font-size: 1.05rem;
		color: var(--text-main);
	}
	.mk-subtitle {
		font-size: 0.85rem;
		color: var(--text-muted);
		margin-top: 4px;
	}
	.mk-date {
		font-size: 0.75rem;
		color: var(--primary-color);
		margin-top: 8px;
		font-weight: 500;
	}
</style>
