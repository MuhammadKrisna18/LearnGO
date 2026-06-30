<script lang="ts">
	import { onMount } from 'svelte';
	import { matakuliahService } from '$lib/services/matakuliah';
	import type { MataKuliah } from '$lib/types';
	import { fade } from 'svelte/transition';

	let mataKuliahList = $state<MataKuliah[]>([]);
	let loading = $state(true);
	let error = $state('');

	onMount(async () => {
		try {
			const res = await matakuliahService.getMahasiswaMataKuliah();
			if (res.success && res.data) {
				mataKuliahList = res.data;
			} else {
				error = res.message || 'Gagal memuat daftar mata kuliah';
			}
		} catch (e: any) {
			error = e.message || 'Terjadi kesalahan saat memuat data';
		} finally {
			loading = false;
		}
	});
</script>

<svelte:head>
	<title>Mata Kuliah - SIAKAD Pro</title>
</svelte:head>

<div class="page-header">
	<h1>Mata Kuliah Anda</h1>
	<p>Daftar mata kuliah yang sesuai dengan Program Studi Anda.</p>
</div>

{#if loading}
	<div class="loading-state">
		<div class="spinner"></div>
		<p>Memuat mata kuliah...</p>
	</div>
{:else if error}
	<div class="error-state">
		<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
			<circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>
		</svg>
		<p>{error}</p>
		<button class="btn btn-primary" onclick={() => window.location.reload()}>Coba Lagi</button>
	</div>
{:else if mataKuliahList.length === 0}
	<div class="empty-state glass-panel">
		<div class="empty-icon">
			<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
				<path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1 0-5H20"/>
			</svg>
		</div>
		<h2>Belum Ada Mata Kuliah</h2>
		<p>Program Studi Anda belum memiliki mata kuliah yang terdaftar.</p>
	</div>
{:else}
	<div class="grid-container">
		{#each mataKuliahList as mk (mk.id)}
			<div class="mk-card glass-panel" in:fade>
				<div class="mk-header">
					<h3 class="mk-title">{mk.name}</h3>
					<span class="sks-badge">{mk.sks} SKS</span>
				</div>
				<div class="mk-body">
					{#if mk.pengajuan && mk.pengajuan.length > 0}
						{@const approved = mk.pengajuan.filter(p => p.status === 'approved')}
						{#if approved.length > 0}
							<div class="dosen-list">
								<p class="section-label">Dosen Pengampu:</p>
								{#each approved as peng (peng.id)}
									<div class="dosen-item">
										<div class="avatar">
											{peng.dosen?.name?.charAt(0) || 'D'}
										</div>
										<span class="dosen-name">{peng.dosen?.name || 'Dosen'}</span>
									</div>
								{/each}
							</div>
						{:else}
							<p class="text-muted">Belum ada dosen yang disetujui untuk mata kuliah ini.</p>
						{/if}
					{:else}
						<p class="text-muted">Belum ada dosen pengampu.</p>
					{/if}
				</div>
			</div>
		{/each}
	</div>
{/if}

<style>
	.page-header {
		margin-bottom: 2rem;
	}

	.page-header h1 {
		font-size: 1.75rem;
		font-weight: 700;
		color: var(--text-main);
		margin: 0 0 0.5rem 0;
	}

	.page-header p {
		color: var(--text-muted);
		margin: 0;
		font-size: 1.05rem;
	}

	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 5rem 2rem;
		text-align: center;
		border-radius: 16px;
		border: 1px dashed var(--border);
	}

	.empty-icon {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 80px;
		height: 80px;
		background: white;
		border-radius: 50%;
		color: var(--primary);
		margin-bottom: 1.5rem;
		box-shadow: 0 4px 12px rgba(79, 70, 229, 0.15);
	}

	.grid-container {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
		gap: 1.5rem;
	}

	.mk-card {
		display: flex;
		flex-direction: column;
		padding: 1.5rem;
		border-radius: 12px;
		transition: transform 0.2s, box-shadow 0.2s;
	}

	.mk-card:hover {
		transform: translateY(-4px);
		box-shadow: var(--shadow-md);
	}

	.mk-header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		margin-bottom: 1rem;
		padding-bottom: 1rem;
		border-bottom: 1px solid var(--border);
	}

	.mk-title {
		font-size: 1.25rem;
		font-weight: 600;
		color: var(--text-main);
		margin: 0;
	}

	.sks-badge {
		background: rgba(79, 70, 229, 0.1);
		color: var(--primary);
		padding: 0.25rem 0.75rem;
		border-radius: 20px;
		font-size: 0.875rem;
		font-weight: 600;
	}

	.section-label {
		font-size: 0.875rem;
		font-weight: 600;
		color: var(--text-muted);
		margin: 0 0 0.5rem 0;
	}

	.dosen-list {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.dosen-item {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}

	.avatar {
		width: 32px;
		height: 32px;
		border-radius: 50%;
		background: linear-gradient(135deg, var(--primary), var(--primary-hover));
		color: white;
		display: flex;
		align-items: center;
		justify-content: center;
		font-weight: 600;
		font-size: 0.875rem;
	}

	.dosen-name {
		font-size: 0.95rem;
		color: var(--text-main);
		font-weight: 500;
	}

	.loading-state, .error-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 4rem 0;
		gap: 1rem;
		color: var(--text-muted);
	}

	.error-state {
		color: var(--danger);
	}

	.spinner {
		width: 40px;
		height: 40px;
		border: 3px solid rgba(79, 70, 229, 0.1);
		border-radius: 50%;
		border-top-color: var(--primary);
		animation: spin 1s linear infinite;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}

	.empty-state h2 {
		margin: 0 0 0.75rem 0;
		color: var(--text-main);
		font-size: 1.5rem;
	}

	.empty-state p {
		margin: 0;
		color: var(--text-muted);
		max-width: 500px;
		line-height: 1.6;
	}
</style>
