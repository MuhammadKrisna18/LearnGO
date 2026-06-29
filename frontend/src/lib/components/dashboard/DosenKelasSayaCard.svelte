<script lang="ts">
	import { onMount } from 'svelte';
	import { kelasService } from '$lib/services/kelas';
	import type { PengajuanKelas } from '$lib/types';
	import Card from '$lib/components/ui/Card.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';

	let myClasses = $state<PengajuanKelas[]>([]);
	let loading = $state(true);
	let error = $state('');

	async function loadMyClasses() {
		loading = true;
		error = '';
		try {
			const res = await kelasService.getMyPengajuan();
			if (res.success && res.data) {
				myClasses = res.data.filter((req: PengajuanKelas) => req.status === 'approved');
			}
		} catch (err: any) {
			error = err.message || 'Gagal memuat kelas saya';
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		loadMyClasses();
	});
</script>

{#if !loading && myClasses.length > 0}
<Card style="padding: 24px; margin-top: 24px; border-left: 4px solid var(--primary-color);">
	<div class="card-header-flex">
		<div style="display: flex; align-items: center; gap: 12px;">
			<div class="header-icon">👨‍🏫</div>
			<div>
				<h3>Kelas Saya</h3>
				<p class="header-subtitle">Daftar kelas yang saat ini Anda ampu</p>
			</div>
		</div>
		<button class="btn-secondary" onclick={loadMyClasses} disabled={loading}>
			Refresh
		</button>
	</div>

	{#if error}
		<div class="state-container state-error" style="padding: 16px; margin-bottom: 16px;">{error}</div>
	{/if}

	<div class="kelas-grid">
		{#each myClasses as req}
			<div class="kelas-card">
				<div class="kelas-header">
					<div class="kelas-title">{req.kelas?.name || 'Unknown'}</div>
					{#if req.kelas?.program_studi}
						<Badge type="primary">{req.kelas.program_studi.name}</Badge>
					{/if}
				</div>
				<div class="kelas-body">
					{#if req.mata_kuliah}
						<div class="info-row">
							<span class="info-icon">📚</span>
							<span class="info-text"><strong>{req.mata_kuliah.name}</strong> ({req.mata_kuliah.sks} SKS)</span>
						</div>
					{/if}
					<div class="info-row">
						<span class="info-icon">📅</span>
						<span class="info-text">{req.kelas?.hari || '-'}</span>
					</div>
					<div class="info-row">
						<span class="info-icon">🕐</span>
						<span class="info-text">{req.kelas?.jam_mulai || ''} - {req.kelas?.jam_selesai || ''}</span>
					</div>
					<div class="info-row">
						<span class="info-icon">👥</span>
						<span class="info-text">Kapasitas: {req.kelas?.capacity || 0} Mahasiswa</span>
					</div>
				</div>
			</div>
		{/each}
	</div>
</Card>
{/if}

<style>
	.card-header-flex {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		margin-bottom: 24px;
	}

	.header-icon {
		font-size: 2rem;
		background: rgba(79, 70, 229, 0.1);
		width: 48px;
		height: 48px;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: var(--radius-md);
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

	.kelas-grid {
		display: grid;
		gap: 16px;
		grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
	}

	.kelas-card {
		background: rgba(255,255,255,0.9);
		border: 1px solid rgba(0,0,0,0.06);
		border-radius: var(--radius-md);
		padding: 16px;
		transition: transform 0.2s, box-shadow 0.2s;
	}

	.kelas-card:hover {
		transform: translateY(-2px);
		box-shadow: var(--shadow-sm);
	}

	.kelas-header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		margin-bottom: 12px;
		padding-bottom: 12px;
		border-bottom: 1px solid rgba(0,0,0,0.05);
	}

	.kelas-title {
		font-weight: 600;
		font-size: 1.1rem;
		color: var(--text-main);
	}

	.kelas-body {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.info-row {
		display: flex;
		align-items: center;
		gap: 8px;
	}

	.info-icon {
		font-size: 0.9rem;
	}

	.info-text {
		font-size: 0.85rem;
		color: var(--text-muted);
	}

	.info-text strong {
		color: var(--text-main);
	}
</style>
