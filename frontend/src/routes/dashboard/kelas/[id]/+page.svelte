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
	
	import { authState } from '$lib/stores/auth.svelte';

	let activeTab = $state<'pertemuan' | 'mahasiswa' | 'materi' | 'tugas' | 'pengumuman' | 'rekap'>('pertemuan');

	let rekapData = $state<any>(null);
	let loadingRekap = $state(false);
	let errorRekap = $state('');

	async function loadRekap(pengajuanId: string) {
		loadingRekap = true;
		errorRekap = '';
		try {
			const res = await kelasService.getRekapKehadiran(pengajuanId);
			if (res.success && res.data) {
				rekapData = res.data;
			}
		} catch (err: any) {
			errorRekap = err.message || 'Gagal memuat rekap kehadiran';
		} finally {
			loadingRekap = false;
		}
	}
	
	let mahasiswaList = $state<any[]>([]);
	let loadingMahasiswa = $state(false);
	let errorMahasiswa = $state('');

	async function loadMahasiswa(pengajuanId: string) {
		loadingMahasiswa = true;
		errorMahasiswa = '';
		try {
			const res = await kelasService.getMahasiswaInKelas(pengajuanId);
			if (res.success && res.data) {
				mahasiswaList = res.data;
			}
		} catch (err: any) {
			errorMahasiswa = err.message || 'Gagal memuat daftar mahasiswa';
		} finally {
			loadingMahasiswa = false;
		}
	}

	let pertemuanList = $state<any[]>([]);
	let loadingPertemuan = $state(false);
	
	let absensiList = $state<any[]>([]);
	let loadingAbsensi = $state(false);
	let selectedPertemuan = $state<any>(null);

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

	async function mulaiPertemuan() {
		if (!pengajuan) {
			alert("Data pengajuan belum tersedia. Silakan muat ulang halaman.");
			return;
		}
		const judul = prompt("Masukkan Topik/Judul Pertemuan:");
		if (!judul) return;
		try {
			await kelasService.mulaiPertemuan(pengajuan.id, judul);
			loadPertemuan(pengajuan.id);
		} catch (err: any) {
			alert('Gagal memulai pertemuan: ' + (err.message || 'Error tidak diketahui'));
		}
	}

	async function akhiriPertemuan(id: string) {
		if (!confirm('Akhiri pertemuan ini? Absensi akan dikunci.')) return;
		try {
			await kelasService.akhiriPertemuan(id);
			if (pengajuan) loadPertemuan(pengajuan.id);
			selectedPertemuan = null;
		} catch (err) {
			alert('Gagal mengakhiri pertemuan');
		}
	}

	async function openAbsensi(pertemuan: any) {
		selectedPertemuan = pertemuan;
		loadingAbsensi = true;
		try {
			const res = await kelasService.getAbsensi(pertemuan.id);
			if (res.success && res.data && res.data.length > 0) {
				absensiList = res.data;
			} else {
				// Initialize based on mahasiswaList
				absensiList = mahasiswaList.map(m => ({
					mahasiswa_id: m.id,
					mahasiswa: m,
					status_kehadiran: 'alpa' // Default
				}));
			}
		} catch (err) {
			console.error('Failed to load absensi:', err);
		} finally {
			loadingAbsensi = false;
		}
	}

	async function submitAbsensi() {
		if (!selectedPertemuan) return;
		try {
			const payload = absensiList.map(a => ({
				mahasiswa_id: a.mahasiswa_id,
				status_kehadiran: a.status_kehadiran
			}));
			await kelasService.submitAbsensi(selectedPertemuan.id, payload);
			alert('Absensi berhasil disimpan!');
			openAbsensi(selectedPertemuan); // Reload to get real IDs
		} catch (err) {
			alert('Gagal menyimpan absensi');
		}
	}

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
						if (authState.role === 'admin') {
							activeTab = 'rekap';
							loadRekap(approved.id);
						} else {
							loadPertemuan(approved.id);
							loadMahasiswa(approved.id);
							loadRekap(approved.id);
						}
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

		<!-- Navigation Tabs -->
		<div class="tabs-container">
			{#if authState.role === 'dosen'}
			<button class="tab-btn" class:active={activeTab === 'pertemuan'} onclick={() => activeTab = 'pertemuan'}>
				Pertemuan
			</button>
			<button class="tab-btn" class:active={activeTab === 'mahasiswa'} onclick={() => activeTab = 'mahasiswa'}>
				Daftar Mahasiswa
			</button>
			{/if}
			<button class="tab-btn" class:active={activeTab === 'rekap'} onclick={() => activeTab = 'rekap'}>
				Rekap Kehadiran
			</button>
			{#if authState.role === 'dosen'}
			<button class="tab-btn" class:active={activeTab === 'materi'} onclick={() => activeTab = 'materi'}>
				Materi
			</button>
			<button class="tab-btn" class:active={activeTab === 'tugas'} onclick={() => activeTab = 'tugas'}>
				Tugas
			</button>
			<button class="tab-btn" class:active={activeTab === 'pengumuman'} onclick={() => activeTab = 'pengumuman'}>
				Pengumuman
			</button>
			{/if}
		</div>

		<div class="content-area">
			<Card style="padding: 24px;">
				{#if activeTab === 'pertemuan'}
					<div class="placeholder-section">
						<div class="section-header">
							<h3>Sesi Pertemuan</h3>
							<button class="btn-primary-sm" onclick={mulaiPertemuan}>Mulai Pertemuan Baru</button>
						</div>
						
						{#if loadingPertemuan}
							<div class="state-container">
								<div class="spinner"></div>
								<p>Memuat daftar pertemuan...</p>
							</div>
						{:else if pertemuanList.length === 0}
							<div class="empty-state">
								<div class="empty-icon">📅</div>
								<p>Belum ada sesi pertemuan yang dimulai.</p>
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
											{#if p.status === 'berlangsung'}
												<div style="margin-top: 12px;">
													<span style="font-size: 0.9rem; color: var(--text-muted);">Kode Absensi:</span>
													<strong style="font-size: 1.5rem; font-family: monospace; background: #f1f5f9; padding: 4px 12px; border-radius: 4px; color: var(--primary-color); letter-spacing: 2px;">{p.kode_absensi}</strong>
												</div>
											{/if}
										</div>
										<div class="pertemuan-actions">
											<Badge type={p.status === 'berlangsung' ? 'success' : 'neutral'}>
												{p.status === 'berlangsung' ? 'Sedang Berlangsung' : 'Selesai'}
											</Badge>
											<button class="btn-secondary-sm" onclick={() => openAbsensi(p)}>Lihat Absensi</button>
											{#if p.status === 'berlangsung'}
												<button class="btn-danger-sm" onclick={() => akhiriPertemuan(p.id)}>Akhiri Kelas</button>
											{/if}
										</div>
									</div>
								{/each}
							</div>
						{/if}

						<!-- Absensi Section (Inline Modal) -->
						{#if selectedPertemuan}
							<div class="modal-backdrop">
								<div class="modal-content absensi-modal">
									<div class="modal-header">
										<h3>Absensi - {selectedPertemuan.judul}</h3>
										<button class="close-btn" onclick={() => selectedPertemuan = null}>✕</button>
									</div>
									<div class="modal-body">
										{#if loadingAbsensi}
											<div class="state-container">
												<div class="spinner"></div>
												<p>Memuat data absensi...</p>
											</div>
										{:else}
											<table class="data-table">
												<thead>
													<tr>
														<th>NRP</th>
														<th>Nama</th>
														<th>Status Kehadiran</th>
													</tr>
												</thead>
												<tbody>
													{#each absensiList as a}
														<tr>
															<td style="font-family: monospace;">{a.mahasiswa?.nrp}</td>
															<td><strong>{a.mahasiswa?.name}</strong></td>
															<td>
																{#if a.status_kehadiran === 'hadir'}
																	<Badge type="success">Hadir</Badge>
																{:else if a.status_kehadiran === 'izin'}
																	<Badge type="warning">Izin</Badge>
																{:else if a.status_kehadiran === 'sakit'}
																	<Badge type="info">Sakit</Badge>
																{:else}
																	<Badge type="danger">Alpa</Badge>
																{/if}
															</td>
														</tr>
													{/each}
												</tbody>
											</table>
										{/if}
									</div>
									<div class="modal-footer">
										<button class="btn-secondary" onclick={() => selectedPertemuan = null}>Tutup</button>
									</div>
								</div>
							</div>
						{/if}
					</div>
				{:else if activeTab === 'mahasiswa'}
					<div class="placeholder-section">
						<div class="section-header">
							<h3>Daftar Mahasiswa</h3>
							<Badge type="info">{mahasiswaList.length} / {kelasInfo.capacity} Mahasiswa</Badge>
						</div>
						
						{#if loadingMahasiswa}
							<div class="state-container">
								<div class="spinner"></div>
								<p>Memuat daftar mahasiswa...</p>
							</div>
						{:else if errorMahasiswa}
							<div class="state-container state-error">
								<p>{errorMahasiswa}</p>
								<button class="btn-primary-sm" onclick={() => pengajuan && loadMahasiswa(pengajuan.id)}>Coba Lagi</button>
							</div>
						{:else if mahasiswaList.length === 0}
							<div class="empty-state">
								<div class="empty-icon">👥</div>
								<p>Belum ada mahasiswa yang terdaftar di kelas ini.</p>
							</div>
						{:else}
							<div class="table-container">
								<table class="data-table">
									<thead>
										<tr>
											<th>No</th>
											<th>NRP</th>
											<th>Nama Lengkap</th>
											<th>Email</th>
										</tr>
									</thead>
									<tbody>
										{#each mahasiswaList as mhs, i}
											<tr>
												<td>{i + 1}</td>
												<td style="font-family: monospace;">{mhs.nrp}</td>
												<td><strong>{mhs.name}</strong></td>
												<td>{mhs.email}</td>
											</tr>
										{/each}
									</tbody>
								</table>
							</div>
						{/if}
					</div>
				{:else if activeTab === 'rekap'}
					<div class="placeholder-section">
						<div class="section-header">
							<h3>Rekap Kehadiran (16 Pertemuan)</h3>
							<button class="btn-primary-sm" onclick={() => pengajuan && loadRekap(pengajuan.id)}>Refresh</button>
						</div>
						
						{#if loadingRekap}
							<div class="state-container">
								<div class="spinner"></div>
								<p>Memuat rekap kehadiran...</p>
							</div>
						{:else if errorRekap}
							<div class="state-container state-error">
								<p>{errorRekap}</p>
								<button class="btn-primary-sm" onclick={() => pengajuan && loadRekap(pengajuan.id)}>Coba Lagi</button>
							</div>
						{:else if !rekapData || !rekapData.mahasiswa || rekapData.mahasiswa.length === 0}
							<div class="empty-state">
								<div class="empty-icon">📊</div>
								<p>Belum ada data rekap kehadiran.</p>
							</div>
						{:else}
							<div class="table-container rekap-container">
								<table class="data-table rekap-table">
									<thead>
										<tr>
											<th class="sticky-col sticky-col-1">NRP</th>
											<th class="sticky-col sticky-col-2">Nama Lengkap</th>
											{#if authState.role === 'admin'}
												<th class="center-col">Total Kehadiran</th>
												<th class="center-col">Persentase</th>
											{:else}
												{#each Array(16) as _, i}
													<th class="center-col">P{i + 1}</th>
												{/each}
											{/if}
										</tr>
									</thead>
									<tbody>
										{#if authState.role === 'admin' && rekapData.summary}
											{#each rekapData.summary as sum}
												<tr>
													<td class="sticky-col sticky-col-1" style="font-family: monospace;">{sum.nrp}</td>
													<td class="sticky-col sticky-col-2"><strong>{sum.name}</strong></td>
													<td class="center-col">
														<Badge type={sum.total_hadir >= 12 ? 'success' : 'danger'}>{sum.total_hadir} / {sum.total_pertemuan} Kali</Badge>
													</td>
													<td class="center-col">
														<strong>{sum.total_pertemuan > 0 ? Math.round((sum.total_hadir / sum.total_pertemuan) * 100) : 0}%</strong>
													</td>
												</tr>
											{/each}
										{:else}
											{#each rekapData.mahasiswa as mhs}
												<tr>
													<td class="sticky-col sticky-col-1" style="font-family: monospace;">{mhs.nrp}</td>
													<td class="sticky-col sticky-col-2"><strong>{mhs.name}</strong></td>
													{#each Array(16) as _, i}
														{@const pId = rekapData.pertemuan[i]?.id}
														{@const status = pId ? mhs.kehadiran[pId] : null}
														<td class="center-col">
															{#if status === 'hadir'}
																<span title="Hadir" style="color: var(--success-color); font-weight: bold;">H</span>
															{:else if status === 'alpa'}
																<span title="Alpa" style="color: var(--danger-color); font-weight: bold;">A</span>
															{:else if status === 'izin'}
																<span title="Izin" style="color: var(--warning-color); font-weight: bold;">I</span>
															{:else if status === 'sakit'}
																<span title="Sakit" style="color: var(--warning-color); font-weight: bold;">S</span>
															{:else}
																<span style="color: var(--text-muted); opacity: 0.5;">-</span>
															{/if}
														</td>
													{/each}
												</tr>
											{/each}
										{/if}
									</tbody>
								</table>
							</div>
						{/if}
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

	.state-error p {
		color: #ef4444;
		margin: 0;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}

	.table-container {
		width: 100%;
		overflow-x: auto;
		background: white;
		border-radius: var(--radius-md);
		border: 1px solid rgba(0,0,0,0.05);
	}

	.data-table {
		width: 100%;
		border-collapse: collapse;
		text-align: left;
	}

	.data-table th, .data-table td {
		padding: 16px;
		border-bottom: 1px solid rgba(0,0,0,0.05);
	}

	.data-table th {
		background: #f8fafc;
		font-weight: 600;
		color: var(--text-muted);
		font-size: 0.85rem;
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	.data-table tr:last-child td {
		border-bottom: none;
	}

	.data-table tbody tr:hover {
		background: #f8fafc;
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

	.btn-danger-sm {
		background: #ef4444;
		color: white;
		border: none;
		padding: 8px 16px;
		border-radius: var(--radius-sm);
		font-size: 0.85rem;
		font-weight: 500;
		cursor: pointer;
		transition: background 0.2s;
	}
	.btn-danger-sm:hover {
		background: #dc2626;
	}

	.btn-secondary-sm {
		background: transparent;
		color: var(--text-main);
		border: 1px solid rgba(0,0,0,0.1);
		padding: 8px 16px;
		border-radius: var(--radius-sm);
		font-size: 0.85rem;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s;
	}
	.btn-secondary-sm:hover {
		background: rgba(0,0,0,0.05);
	}

	.absensi-modal {
		max-width: 700px;
		width: 90%;
		max-height: 90vh;
		display: flex;
		flex-direction: column;
	}
	.absensi-modal .modal-body {
		overflow-y: auto;
		padding: 0;
	}

	@media (max-width: 768px) {
		.kelas-title {
			font-size: 2rem;
		}
		
		.kelas-banner {
			padding: 24px;
		}
	}

	.rekap-container {
		overflow-x: auto;
		max-width: 100%;
	}

	.rekap-table {
		min-width: 1200px;
	}

	.rekap-table th, .rekap-table td {
		white-space: nowrap;
	}

	.sticky-col {
		position: sticky;
		background-color: white;
		z-index: 2;
	}

	.rekap-table thead .sticky-col {
		z-index: 3;
	}

	.sticky-col-1 {
		left: 0;
		border-right: 1px solid rgba(0,0,0,0.05);
	}

	.sticky-col-2 {
		left: 120px; /* Adjust based on NRP width */
		border-right: 2px solid rgba(0,0,0,0.1);
		box-shadow: 2px 0 5px rgba(0,0,0,0.02);
	}

	.center-col {
		text-align: center;
		min-width: 40px;
	}
</style>
