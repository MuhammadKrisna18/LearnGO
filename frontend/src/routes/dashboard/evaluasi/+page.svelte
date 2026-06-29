<script lang="ts">
	import { onMount } from 'svelte';
	import { authState } from '$lib/stores/auth.svelte';
	import { matakuliahService } from '$lib/services/matakuliah';
	import { dosenService } from '$lib/services/dosen';
	import { kelasService } from '$lib/services/kelas';
	import type { MataKuliah, UserProfile, Kelas, Dosen } from '$lib/types';
	import Card from '$lib/components/ui/Card.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';

	let mkKosong: MataKuliah[] = $state([]);
	let dosenNganggur: Dosen[] = $state([]);
	let kelasKosong: Kelas[] = $state([]);
	
	let loading = $state(true);
	let error = $state('');

	onMount(async () => {
		try {
			loading = true;

			const [resMk, resDosen, resKelas] = await Promise.all([
				matakuliahService.getList(),
				dosenService.getList(),
				kelasService.getList()
			]);

			if (!resMk.success || !resDosen.success || !resKelas.success) {
				error = 'Gagal memuat beberapa data evaluasi dari server.';
			}

			const allMk = resMk.data || [];
			const allDosen = resDosen.data || [];
			const allKelas = resKelas.data || [];

			mkKosong = allMk.filter((mk: MataKuliah) => {
				if (!mk.pengajuan || mk.pengajuan.length === 0) return true;
				return !mk.pengajuan.some(p => p.status === 'approved');
			});

			const dosenWithApprovedMk = new Set<string>();
			for (const mk of allMk) {
				if (mk.pengajuan) {
					for (const p of mk.pengajuan) {
						if (p.status === 'approved') {
							dosenWithApprovedMk.add(p.dosen_id);
						}
					}
				}
			}
			dosenNganggur = allDosen.filter((d: Dosen) => !dosenWithApprovedMk.has(d.id));

			kelasKosong = allKelas.filter((k: Kelas) => {
				if (!k.pengajuan || k.pengajuan.length === 0) return true;
				return !k.pengajuan.some(p => p.status === 'approved');
			});

		} catch (err: any) {
			error = err.message || 'Terjadi kesalahan sistem';
		} finally {
			loading = false;
		}
	});
</script>

<div class="evaluasi-container animate-fade-in">
	<div class="header">
		<h2>Laporan & Evaluasi Akademik</h2>
		<p class="text-muted">Pantau mata kuliah, dosen, dan kelas yang belum dialokasikan secara optimal.</p>
	</div>

	{#if loading}
		<Card class="state-container">
			<div class="spinner"></div>
			<p>Menganalisis data...</p>
		</Card>
	{:else if error}
		<Card class="state-container state-error">{error}</Card>
	{:else}
		<div class="eval-grid">

			<Card class="eval-card">
				<div class="eval-card-header">
					<h3>Mata Kuliah Belum Diambil</h3>
					<Badge type={mkKosong.length > 0 ? "warning" : "success"}>{mkKosong.length} MK</Badge>
				</div>
				<div class="eval-card-body">
					{#if mkKosong.length === 0}
						<div class="empty-list">Semua mata kuliah sudah memiliki pengampu.</div>
					{:else}
						<ul class="item-list">
							{#each mkKosong as mk}
								<li>
									<div class="item-title">{mk.name}</div>
									<div class="item-meta">
										<Badge type="default">{mk.sks} SKS</Badge>
										<span class="prodi-name">{mk.program_studi?.name || 'Unknown Prodi'}</span>
									</div>
								</li>
							{/each}
						</ul>
					{/if}
				</div>
			</Card>

			<Card class="eval-card">
				<div class="eval-card-header">
					<h3>Dosen Belum Mengambil MK</h3>
					<Badge type={dosenNganggur.length > 0 ? "warning" : "success"}>{dosenNganggur.length} Dosen</Badge>
				</div>
				<div class="eval-card-body">
					{#if dosenNganggur.length === 0}
						<div class="empty-list">Semua dosen sudah memiliki mata kuliah.</div>
					{:else}
						<ul class="item-list">
							{#each dosenNganggur as dosen}
								<li>
									<div class="dosen-info">
										<div class="avatar">
											{#if dosen.photo_url}
												<img src={`http://localhost:8080${dosen.photo_url}`} alt={dosen.name} class="avatar-img-small" />
											{:else}
												{dosen.name.charAt(0).toUpperCase()}
											{/if}
										</div>
										<div class="dosen-details">
											<div class="item-title">{dosen.name}</div>
											<div class="item-meta">{dosen.email}</div>
										</div>
									</div>
								</li>
							{/each}
						</ul>
					{/if}
				</div>
			</Card>

			<Card class="eval-card">
				<div class="eval-card-header">
					<h3>Kelas Belum Digunakan</h3>
					<Badge type="info">{kelasKosong.length} Kelas</Badge>
				</div>
				<div class="eval-card-body">
					{#if kelasKosong.length === 0}
						<div class="empty-list">Belum ada kelas yang dibuat.</div>
					{:else}
						<ul class="item-list">
							{#each kelasKosong as k}
								<li>
									<div class="item-title">{k.name}</div>
									<div class="item-meta">
										<span class="schedule">{k.hari}, {k.jam_mulai} - {k.jam_selesai}</span>
										<Badge type="default">{k.capacity} Mhs</Badge>
									</div>
								</li>
							{/each}
						</ul>
					{/if}
				</div>
			</Card>
		</div>
	{/if}
</div>

<style>
	.evaluasi-container {
		display: flex;
		flex-direction: column;
		gap: 24px;
	}

	.header h2 {
		font-size: 1.5rem;
		font-weight: 700;
		color: var(--text-main);
		margin-bottom: 8px;
	}

	.eval-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
		gap: 24px;
	}

	.eval-card {
		display: flex;
		flex-direction: column;
		max-height: 500px;
	}

	.eval-card-header {
		padding: 20px;
		border-bottom: 1px solid var(--surface-border);
		display: flex;
		justify-content: space-between;
		align-items: center;
		background: rgba(255, 255, 255, 0.02);
	}

	.eval-card-header h3 {
		font-size: 1.1rem;
		font-weight: 600;
		margin: 0;
		color: var(--primary-color);
	}

	.eval-card-body {
		flex: 1;
		padding: 20px;
		overflow-y: auto;
	}

	.item-list {
		list-style: none;
		padding: 0;
		margin: 0;
		display: flex;
		flex-direction: column;
		gap: 16px;
	}

	.item-list li {
		padding-bottom: 16px;
		border-bottom: 1px dashed var(--surface-border);
	}

	.item-list li:last-child {
		border-bottom: none;
		padding-bottom: 0;
	}

	.item-title {
		font-weight: 600;
		color: var(--text-main);
		margin-bottom: 4px;
	}

	.item-meta {
		display: flex;
		align-items: center;
		gap: 12px;
		font-size: 0.85rem;
		color: var(--text-muted);
	}

	.dosen-info {
		display: flex;
		align-items: center;
		gap: 12px;
	}

	.dosen-details {
		display: flex;
		flex-direction: column;
	}

	.avatar {
		width: 32px;
		height: 32px;
		background: var(--primary-color);
		color: white;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		font-weight: 600;
		font-size: 0.85rem;
		overflow: hidden;
		flex-shrink: 0;
	}

	.avatar-img-small {
		width: 100%;
		height: 100%;
		object-fit: cover;
		aspect-ratio: 1/1;
	}

	.prodi-name {
		background: rgba(0,0,0,0.04);
		padding: 2px 6px;
		border-radius: 4px;
	}

	.schedule {
		font-weight: 500;
	}

	.empty-list {
		text-align: center;
		padding: 32px 16px;
		color: var(--text-muted);
		font-style: italic;
	}

	.state-container {
		padding: 48px;
		text-align: center;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 16px;
		color: var(--text-muted);
	}

	.spinner {
		width: 32px;
		height: 32px;
		border: 3px solid rgba(0, 0, 0, 0.1);
		border-top-color: var(--primary-color);
		border-radius: 50%;
		animation: spin 1s linear infinite;
	}

	.alert-info {
		background: rgba(59, 130, 246, 0.1);
		color: #2563eb;
		padding: 12px;
		border-radius: var(--radius-sm);
		font-size: 0.85rem;
		margin-bottom: 16px;
		border: 1px solid rgba(59, 130, 246, 0.2);
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}
</style>
