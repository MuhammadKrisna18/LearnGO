<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { kelasService } from '$lib/services/kelas';
	import { fetchApi } from '$lib/services/api';
	import Card from '$lib/components/ui/Card.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';

	let id = $page.params.id;
	let loading = $state(true);
	let error = $state('');

	let activeTab = $state<'pertemuan' | 'materi' | 'tugas' | 'pengumuman'>('pertemuan');
	
	let kelasInfo = $state<any>(null);
	let pertemuanList = $state<any[]>([]);
	let loadingPertemuan = $state(false);

	let absensiCode = $state('');
	let submittingAbsensi = $state(false);

	async function loadKelas() {
		loading = true;
		error = '';
		try {
			
			
			const res = await kelasService.getMyJadwal();
			if (res.success && res.data) {
				const myClass = res.data.find(p => p.id === id);
				if (myClass) {
					kelasInfo = myClass;
					loadPertemuan(myClass.id);
				} else {
					error = 'Kelas tidak ditemukan atau Anda tidak terdaftar';
				}
			}
		} catch (err: any) {
			error = err.message || 'Gagal memuat data kelas';
		} finally {
			loading = false;
		}
	}

	async function loadPertemuan(pengajuanId: string) {
		loadingPertemuan = true;
		try {
			const res = await kelasService.getPertemuan(pengajuanId);
			if (res.success && res.data) {
				pertemuanList = res.data;
			}
		} catch (err) {
			console.error('Failed to load pertemuan:', err);
		} finally {
			loadingPertemuan = false;
		}
	}

	async function submitAbsensi(pertemuanId: string) {
		if (!absensiCode || absensiCode.length !== 6) {
			alert('Masukkan 6 digit kode absensi dengan benar');
			return;
		}
		submittingAbsensi = true;
		try {
			const res = await fetchApi<any>(`/pertemuan/${pertemuanId}/attend`, {
				method: 'POST',
				body: JSON.stringify({ kode: absensiCode })
			});
			if (res.success) {
				alert('Berhasil melakukan absensi!');
				absensiCode = '';
				loadPertemuan(kelasInfo.id);
			} else {
				alert(res.message || 'Gagal melakukan absensi');
			}
		} catch (err: any) {
			alert(err.message || 'Gagal melakukan absensi');
		} finally {
			submittingAbsensi = false;
		}
	}

	onMount(() => {
		loadKelas();
	});
</script>

<svelte:head>
	<title>{kelasInfo ? kelasInfo.kelas?.name : 'Detail Kelas'} - SIAKAD Pro</title>
</svelte:head>

<div class="page-container">
	{#if loading}
		<div class="state-container">
			<div class="spinner"></div>
			<p>Memuat informasi kelas...</p>
		</div>
	{:else if error}
		<div class="state-container state-error">
			<p>{error}</p>
			<button class="btn-primary" onclick={loadKelas}>Coba Lagi</button>
		</div>
	{:else if kelasInfo}
		
		<div class="kelas-header">
			<div class="header-content">
				<h1 class="kelas-title">{kelasInfo.kelas?.name}</h1>
				<div class="mata-kuliah-info">
					<span class="mk-name">{kelasInfo.mata_kuliah?.name}</span>
					<Badge type="primary">{kelasInfo.mata_kuliah?.sks} SKS</Badge>
				</div>
				<p class="dosen-name">Dosen Pengampu: {kelasInfo.dosen?.name || 'Belum diatur'}</p>
			</div>
			
			<div class="schedule-card">
				<div class="schedule-icon">📅</div>
				<div class="schedule-details">
					<span class="day">{kelasInfo.kelas?.hari}</span>
					<span class="time">{kelasInfo.kelas?.jam_mulai} - {kelasInfo.kelas?.jam_selesai}</span>
					<span class="room">Ruang: {kelasInfo.kelas?.ruangan || 'TBA'}</span>
				</div>
			</div>
		</div>

		
		<div class="tabs-container">
			<button class="tab-btn" class:active={activeTab === 'pertemuan'} onclick={() => activeTab = 'pertemuan'}>
				Pertemuan
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
				{#if activeTab === 'pertemuan'}
					<div class="placeholder-section">
						<div class="section-header">
							<h3>Sesi Pertemuan</h3>
						</div>
						
						{#if loadingPertemuan}
							<div class="state-container">
								<div class="spinner"></div>
								<p>Memuat daftar pertemuan...</p>
							</div>
						{:else if pertemuanList.length === 0}
							<div class="empty-state">
								<div class="empty-icon">📅</div>
								<p>Belum ada sesi pertemuan.</p>
							</div>
						{:else}
							<div class="pertemuan-list">
								{#each pertemuanList as p, i}
									<div class="pertemuan-card">
										<div class="pertemuan-info">
											<h4>{p.judul || `Pertemuan ${i + 1}`}</h4>
											<div class="pertemuan-meta">
												<span>{new Date(p.tanggal).toLocaleDateString('id-ID', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })}</span>
												<span>{new Date(p.waktu_mulai).toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })} - {p.waktu_selesai ? new Date(p.waktu_selesai).toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' }) : 'Berlangsung'}</span>
											</div>
										</div>
										<div class="pertemuan-actions">
											{#if p.status === 'berlangsung'}
												<div class="absensi-input-group">
													<input type="text" class="input-kode" placeholder="Kode" maxlength="6" bind:value={absensiCode} disabled={submittingAbsensi} />
													<button class="btn-primary-sm" onclick={() => submitAbsensi(p.id)} disabled={submittingAbsensi || absensiCode.length !== 6}>
														{submittingAbsensi ? 'Wait..' : 'Hadir'}
													</button>
												</div>
											{:else}
												<Badge type="neutral">Selesai</Badge>
											{/if}
										</div>
									</div>
								{/each}
							</div>
						{/if}
					</div>
				{:else}
					<div class="placeholder-section">
						<div class="empty-state">
							<div class="empty-icon">🚧</div>
							<p>Fitur sedang dalam pengembangan</p>
						</div>
					</div>
				{/if}
			</Card>
		</div>
	{/if}
</div>

<style>
	.page-container {
		max-width: 1200px;
		margin: 0 auto;
		padding: 24px;
	}

	.kelas-header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		margin-bottom: 32px;
		background: white;
		padding: 32px;
		border-radius: var(--radius-lg);
		box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);
	}

	.kelas-title {
		margin: 0 0 12px 0;
		font-size: 2.5rem;
		font-weight: 700;
		color: var(--text-main);
		letter-spacing: -0.02em;
	}

	.mata-kuliah-info {
		display: flex;
		align-items: center;
		gap: 12px;
		margin-bottom: 8px;
	}

	.mk-name {
		font-size: 1.25rem;
		color: var(--text-muted);
		font-weight: 500;
	}

	.dosen-name {
		margin: 0;
		color: var(--text-muted);
		font-size: 1rem;
	}

	.schedule-card {
		display: flex;
		align-items: center;
		gap: 16px;
		background: #f8fafc;
		padding: 20px 24px;
		border-radius: var(--radius-md);
		border: 1px solid rgba(0,0,0,0.05);
	}

	.schedule-icon {
		font-size: 2.5rem;
	}

	.schedule-details {
		display: flex;
		flex-direction: column;
		gap: 4px;
	}

	.schedule-details .day {
		font-weight: 600;
		color: var(--text-main);
		font-size: 1.1rem;
	}

	.schedule-details .time, .schedule-details .room {
		color: var(--text-muted);
		font-size: 0.9rem;
	}

	.tabs-container {
		display: flex;
		gap: 8px;
		margin-bottom: 24px;
		border-bottom: 2px solid #e2e8f0;
		padding-bottom: 1px;
	}

	.tab-btn {
		padding: 12px 24px;
		background: transparent;
		border: none;
		font-size: 1rem;
		font-weight: 600;
		color: var(--text-muted);
		cursor: pointer;
		position: relative;
		transition: all 0.2s;
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
		bottom: -3px;
		left: 0;
		right: 0;
		height: 3px;
		background: var(--primary-color);
		border-radius: 3px 3px 0 0;
	}

	.state-container {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 64px 20px;
		color: var(--text-muted);
		text-align: center;
	}

	.spinner {
		width: 40px;
		height: 40px;
		border: 3px solid #f3f3f3;
		border-top: 3px solid var(--primary-color);
		border-radius: 50%;
		animation: spin 1s linear infinite;
		margin-bottom: 16px;
	}

	.empty-state {
		text-align: center;
		padding: 48px 20px;
		color: var(--text-muted);
	}

	.empty-icon {
		font-size: 3rem;
		margin-bottom: 16px;
		opacity: 0.5;
	}

	.btn-primary {
		background-color: var(--primary-color);
		color: white;
		border: none;
		padding: 10px 20px;
		border-radius: var(--radius-sm);
		font-size: 1rem;
		font-weight: 500;
		cursor: pointer;
		transition: background 0.2s;
	}

	.btn-primary:hover {
		background-color: var(--primary-hover);
	}

	.btn-primary-sm {
		background-color: var(--primary-color);
		color: white;
		border: none;
		padding: 8px 16px;
		border-radius: var(--radius-sm);
		font-size: 0.9rem;
		font-weight: 500;
		cursor: pointer;
		transition: background 0.2s;
	}

	.btn-primary-sm:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.pertemuan-list {
		display: flex;
		flex-direction: column;
		gap: 16px;
	}

	.pertemuan-card {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 20px;
		background: white;
		border: 1px solid rgba(0,0,0,0.05);
		border-radius: var(--radius-md);
		box-shadow: 0 2px 4px rgba(0,0,0,0.02);
	}

	.pertemuan-info h4 {
		margin: 0 0 8px 0;
		font-size: 1.1rem;
		color: var(--text-main);
	}

	.pertemuan-meta {
		display: flex;
		gap: 16px;
		font-size: 0.9rem;
		color: var(--text-muted);
	}

	.pertemuan-actions {
		display: flex;
		align-items: center;
		gap: 12px;
	}

	.absensi-input-group {
		display: flex;
		gap: 8px;
		align-items: center;
	}

	.input-kode {
		padding: 8px 12px;
		border: 1px solid rgba(0,0,0,0.2);
		border-radius: var(--radius-sm);
		font-size: 1rem;
		font-family: monospace;
		width: 100px;
		text-align: center;
		letter-spacing: 2px;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}

	@media (max-width: 768px) {
		.kelas-header {
			flex-direction: column;
			gap: 24px;
		}
		
		.schedule-card {
			width: 100%;
			box-sizing: border-box;
		}

		.pertemuan-card {
			flex-direction: column;
			align-items: flex-start;
			gap: 16px;
		}
		
		.pertemuan-actions {
			width: 100%;
			justify-content: flex-end;
		}
	}
</style>
