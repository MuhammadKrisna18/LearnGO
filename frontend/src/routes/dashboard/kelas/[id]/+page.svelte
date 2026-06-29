<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { kelasService } from '$lib/services/kelas';
	import type { Kelas, PengajuanKelas } from '$lib/types';
	import Card from '$lib/components/ui/Card.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';

	let id = $derived($page.params.id);
	let kelasInfo = $state<Kelas | null>(null);
	let pengajuan = $state<PengajuanKelas | null>(null);
	let loading = $state(true);
	let error = $state('');

	let activeTab = $state<'mahasiswa' | 'materi' | 'tugas' | 'pengumuman'>('mahasiswa');

	async function loadKelas() {
		loading = true;
		error = '';
		try {
			const res = await kelasService.getById(id);
			if (res.success && res.data) {
				kelasInfo = res.data;
				
				if (kelasInfo.pengajuan) {

					const approved = kelasInfo.pengajuan.find(p => p.status === 'approved');
					if (approved) {
						pengajuan = approved;
					}
				}
			} else {
				error = res.message || 'Gagal memuat kelas';
			}
		} catch (err: any) {
			error = err.message || 'Terjadi kesalahan saat memuat kelas';
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		loadKelas();
	});
</script>

<div class="kelas-portal animate-fade-in">
	{#if loading}
		<div class="state-container">
			<div class="spinner"></div>
			<p>Memasuki kelas...</p>
		</div>
	{:else if error}
		<div class="state-container state-error">
			<p>{error}</p>
			<a href="/dashboard" class="btn-secondary" style="margin-top: 16px;">Kembali ke Dashboard</a>
		</div>
	{:else if kelasInfo}

		<div class="kelas-banner">
			<div class="banner-content">
				<a href="/dashboard" class="back-link">← Kembali ke Dashboard</a>
				<h1 class="kelas-title">{kelasInfo.name}</h1>
				<div class="kelas-meta">
					{#if pengajuan?.mata_kuliah}
						<span class="meta-item">
							<span class="icon">📚</span>
							{pengajuan.mata_kuliah.name} ({pengajuan.mata_kuliah.sks} SKS)
						</span>
					{/if}
					<span class="meta-item">
						<span class="icon">📅</span>
						{kelasInfo.hari}, {kelasInfo.jam_mulai} - {kelasInfo.jam_selesai}
					</span>
					{#if kelasInfo.program_studi}
						<span class="meta-item">
							<span class="icon">🏫</span>
							{kelasInfo.program_studi.name}
						</span>
					{/if}
					{#if pengajuan?.dosen}
						<span class="meta-item">
							<span class="icon">👨‍🏫</span>
							Dosen: {pengajuan.dosen.name}
						</span>
					{/if}
				</div>
			</div>
			<div class="banner-decoration"></div>
		</div>

		<div class="tabs-container">
			<button class="tab-btn" class:active={activeTab === 'mahasiswa'} onclick={() => activeTab = 'mahasiswa'}>
				Daftar Mahasiswa
			</button>
			<button class="tab-btn" class:active={activeTab === 'materi'} onclick={() => activeTab = 'materi'}>
				Materi
			</button>
			<button class="tab-btn" class:active={activeTab === 'tugas'} onclick={() => activeTab = 'tugas'}>
				Tugas
			</button>
			<button class="tab-btn" class:active={activeTab === 'pengumuman'} onclick={() => activeTab = 'pengumuman'}>
				Pengumuman
			</button>
		</div>

		<div class="content-area">
			<Card style="padding: 24px;">
				{#if activeTab === 'mahasiswa'}
					<div class="placeholder-section">
						<h3>Daftar Mahasiswa</h3>
						<div class="empty-state">
							<div class="empty-icon">👥</div>
							<p>Belum ada mahasiswa yang terdaftar di kelas ini.</p>
							<p class="empty-hint">Kapasitas maksimal: {kelasInfo.capacity} Mahasiswa</p>
						</div>
					</div>
				{:else if activeTab === 'materi'}
					<div class="placeholder-section">
						<div class="section-header">
							<h3>Materi Kuliah</h3>
							<button class="btn-primary-sm">Tambah Materi</button>
						</div>
						<div class="empty-state">
							<div class="empty-icon">📁</div>
							<p>Belum ada materi yang diunggah.</p>
						</div>
					</div>
				{:else if activeTab === 'tugas'}
					<div class="placeholder-section">
						<div class="section-header">
							<h3>Tugas</h3>
							<button class="btn-primary-sm">Buat Tugas</button>
						</div>
						<div class="empty-state">
							<div class="empty-icon">📝</div>
							<p>Belum ada tugas yang diberikan.</p>
						</div>
					</div>
				{:else if activeTab === 'pengumuman'}
					<div class="placeholder-section">
						<div class="section-header">
							<h3>Pengumuman Kelas</h3>
							<button class="btn-primary-sm">Buat Pengumuman</button>
						</div>
						<div class="empty-state">
							<div class="empty-icon">📢</div>
							<p>Belum ada pengumuman.</p>
						</div>
					</div>
				{/if}
			</Card>
		</div>
	{/if}
</div>

<style>
	.kelas-portal {
		display: flex;
		flex-direction: column;
		gap: 24px;
	}

	.kelas-banner {
		background: linear-gradient(135deg, var(--primary-color), var(--primary-hover));
		border-radius: var(--radius-lg);
		padding: 32px;
		color: white;
		position: relative;
		overflow: hidden;
		box-shadow: var(--shadow-md);
	}

	.banner-content {
		position: relative;
		z-index: 2;
	}

	.back-link {
		color: rgba(255, 255, 255, 0.8);
		text-decoration: none;
		font-size: 0.9rem;
		display: inline-block;
		margin-bottom: 16px;
		transition: color 0.2s;
	}

	.back-link:hover {
		color: white;
	}

	.kelas-title {
		font-size: 2.5rem;
		font-weight: 700;
		margin: 0 0 16px 0;
		color: white;
		text-shadow: 0 2px 4px rgba(0,0,0,0.1);
	}

	.kelas-meta {
		display: flex;
		flex-wrap: wrap;
		gap: 20px;
		font-size: 0.95rem;
	}

	.meta-item {
		display: flex;
		align-items: center;
		gap: 8px;
		background: rgba(0, 0, 0, 0.15);
		padding: 6px 12px;
		border-radius: var(--radius-full);
		backdrop-filter: blur(4px);
	}

	.meta-item .icon {
		font-size: 1.1rem;
	}

	.banner-decoration {
		position: absolute;
		top: -50%;
		right: -10%;
		width: 300px;
		height: 300px;
		background: radial-gradient(circle, rgba(255,255,255,0.15) 0%, rgba(255,255,255,0) 70%);
		border-radius: 50%;
		z-index: 1;
	}

	.tabs-container {
		display: flex;
		gap: 8px;
		border-bottom: 2px solid rgba(0,0,0,0.05);
		padding-bottom: 0;
		overflow-x: auto;
	}

	.tab-btn {
		background: none;
		border: none;
		padding: 12px 24px;
		font-size: 1rem;
		font-weight: 600;
		color: var(--text-muted);
		cursor: pointer;
		position: relative;
		transition: color 0.2s;
		white-space: nowrap;
	}

	.tab-btn:hover {
		color: var(--text-main);
	}

	.tab-btn.active {
		color: var(--primary-color);
	}

	.tab-btn.active::after {
		content: '';
		position: absolute;
		bottom: -2px;
		left: 0;
		right: 0;
		height: 3px;
		background: var(--primary-color);
		border-radius: 3px 3px 0 0;
	}

	.section-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 24px;
	}

	.section-header h3 {
		font-size: 1.25rem;
		font-weight: 600;
		margin: 0;
		color: var(--text-main);
	}

	.btn-primary-sm {
		background: var(--primary-color);
		color: white;
		border: none;
		padding: 8px 16px;
		border-radius: var(--radius-sm);
		font-size: 0.85rem;
		font-weight: 500;
		cursor: pointer;
		transition: background 0.2s;
	}

	.btn-primary-sm:hover {
		background: var(--primary-hover);
	}

	.empty-state {
		text-align: center;
		padding: 60px 20px;
		background: rgba(0, 0, 0, 0.02);
		border-radius: var(--radius-md);
		border: 1px dashed rgba(0, 0, 0, 0.1);
	}

	.empty-icon {
		font-size: 3.5rem;
		margin-bottom: 16px;
		opacity: 0.8;
	}

	.empty-state p {
		color: var(--text-muted);
		margin: 0 0 8px 0;
		font-size: 1.1rem;
	}

	.empty-hint {
		font-size: 0.9rem !important;
		opacity: 0.7;
	}

	.state-container {
		padding: 60px;
		text-align: center;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 16px;
		color: var(--text-muted);
	}

	.spinner {
		width: 40px;
		height: 40px;
		border: 4px solid rgba(0, 0, 0, 0.1);
		border-top-color: var(--primary-color);
		border-radius: 50%;
		animation: spin 1s linear infinite;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}

	@media (max-width: 768px) {
		.kelas-title {
			font-size: 2rem;
		}
		
		.kelas-banner {
			padding: 24px;
		}
	}
</style>
