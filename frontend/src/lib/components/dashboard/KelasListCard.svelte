<script lang="ts">
	import { onMount } from 'svelte';
	import { kelasService } from '$lib/services/kelas';
	import { matakuliahService } from '$lib/services/matakuliah';
	import { authState } from '$lib/stores/auth.svelte';
	import type { Kelas, PengajuanMataKuliah } from '$lib/types';
	import PromptCodeModal from '$lib/components/ui/PromptCodeModal.svelte';
	import { toast } from '$lib/stores/toast.svelte';

	let kelases = $state<Kelas[]>([]);
	let loading = $state(true);
	let error = $state('');

	let isPromptOpen = $state(false);
	let selectedKelas = $state<Kelas | null>(null);
	let expectedCode = $state("");
	let promptAction = $state<'delete' | 'request'>('delete');
	
	let myApprovedMataKuliah = $state<PengajuanMataKuliah[]>([]);
	let myActiveKelasCount = $state(0);
	let availableMataKuliahForKelas = $state<PengajuanMataKuliah[]>([]);
	let selectedMataKuliahId = $state("");

	let searchQuery = $state('');
	let sortBy = $state('kelas');

	const hariOrder: Record<string, number> = {
		'Senin': 1,
		'Selasa': 2,
		'Rabu': 3,
		'Kamis': 4,
		'Jumat': 5
	};

	let groupedKelases = $derived.by(() => {
		let filtered = [...kelases];
		
		if (searchQuery.trim() !== '') {
			const q = searchQuery.toLowerCase();
			filtered = filtered.filter(k => 
				k.name.toLowerCase().includes(q) || 
				(k.program_studi?.name || 'Lainnya').toLowerCase().includes(q)
			);
		}

		filtered.sort((a, b) => {
			if (sortBy === 'kelas') {
				const diff = a.name.localeCompare(b.name);
				if (diff !== 0) return diff;
				return (hariOrder[a.hari] || 99) - (hariOrder[b.hari] || 99);
			} else if (sortBy === 'hari') {
				const diff = (hariOrder[a.hari] || 99) - (hariOrder[b.hari] || 99);
				if (diff !== 0) return diff;
				return a.jam_mulai.localeCompare(b.jam_mulai);
			} else if (sortBy === 'jam') {
				const diff = a.jam_mulai.localeCompare(b.jam_mulai);
				if (diff !== 0) return diff;
				return (hariOrder[a.hari] || 99) - (hariOrder[b.hari] || 99);
			}
			return 0;
		});

		const groups: Record<string, Kelas[]> = {};
		for (const k of filtered) {
			const psName = k.program_studi ? k.program_studi.name : 'Lainnya';
			if (!groups[psName]) {
				groups[psName] = [];
			}
			groups[psName].push(k);
		}
		
		return Object.entries(groups)
			.sort((a, b) => a[0].localeCompare(b[0]))
			.map(([name, classes]) => ({
				programStudiName: name,
				classes
			}));
	});

	async function fetchKelases() {
		try {
			loading = true;
			const res = await kelasService.getList();
			if (res.success && res.data) {
				kelases = res.data;
			} else {
				error = res.message || 'Gagal mengambil data kelas';
			}
			
			if (authState.role === 'dosen') {
				const [mkRes, kelasReqRes] = await Promise.all([
					matakuliahService.getMyRequests(),
					kelasService.getMyPengajuan()
				]);
				
				if (mkRes.success && mkRes.data) {
					myApprovedMataKuliah = mkRes.data.filter((p: any) => p.status === 'approved');
				}
				
				if (kelasReqRes.success && kelasReqRes.data) {
					myActiveKelasCount = kelasReqRes.data.filter((p: any) => p.status === 'approved' || p.status === 'pending').length;
				}
			}
		} catch (err: any) {
			error = err.message || 'Terjadi kesalahan sistem';
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		fetchKelases();

		window.addEventListener('kelasCreated', fetchKelases);
		return () => {
			window.removeEventListener('kelasCreated', fetchKelases);
		};
	});

	function promptDelete(k: Kelas) {
		selectedKelas = k;
		expectedCode = Math.floor(100000 + Math.random() * 900000).toString();
		promptAction = 'delete';
		isPromptOpen = true;
	}

	function promptRequest(k: Kelas) {
		selectedKelas = k;
		expectedCode = Math.floor(100000 + Math.random() * 900000).toString();
		promptAction = 'request';
		
		availableMataKuliahForKelas = myApprovedMataKuliah.filter(mk => mk.mata_kuliah?.program_studi_id === k.program_studi_id);
		if (availableMataKuliahForKelas.length > 0) {
			selectedMataKuliahId = availableMataKuliahForKelas[0].mata_kuliah_id;
		} else {
			selectedMataKuliahId = "";
		}

		isPromptOpen = true;
	}

	async function handleConfirmCode(code: string) {
		if (promptAction === 'delete') {
			await handleDelete(code);
		} else if (promptAction === 'request') {
			await handleRequest(code);
		}
	}

	async function handleRequest(code: string) {
		if (!selectedKelas) return;

		if (code !== expectedCode) {
			toast.error('Kode tidak cocok. Aksi dibatalkan.');
			selectedKelas = null;
			isPromptOpen = false;
			return;
		}

		try {
			const res = await kelasService.requestKelas(selectedKelas.id, selectedMataKuliahId);
			if (res.success) {
				toast.success('Berhasil mengajukan request kelas');
				fetchKelases();
			} else {
				toast.error(res.message || 'Gagal mengajukan kelas');
			}
		} catch (err: any) {
			toast.error(err.message || 'Terjadi kesalahan saat mengajukan kelas');
		} finally {
			isPromptOpen = false;
			selectedKelas = null;
		}
	}

	async function handleDelete(code: string) {
		if (!selectedKelas) return;

		if (code !== expectedCode) {
			toast.error('Kode tidak cocok. Aksi dibatalkan.');
			selectedKelas = null;
			isPromptOpen = false;
			return;
		}

		try {
			const res = await kelasService.delete(selectedKelas.id);
			if (res.success) {
				toast.success('Kelas berhasil dihapus');
				kelases = kelases.filter((k) => k.id !== selectedKelas!.id);
			} else {
				toast.error(res.message || 'Gagal menghapus kelas');
			}
		} catch (err: any) {
			toast.error(err.message || 'Terjadi kesalahan saat menghapus kelas');
		} finally {
			isPromptOpen = false;
			selectedKelas = null;
		}
	}
</script>

<div class="card glass-panel animate-fade-in">
	<div class="card-header">
		<div class="header-content">
			<h2>Daftar Kelas</h2>
			<div class="controls">
				<input type="text" class="search-input" bind:value={searchQuery} placeholder="Cari nama kelas atau prodi..." />
				<select class="sort-select" bind:value={sortBy}>
					<option value="kelas">Urutkan: Kelas</option>
					<option value="hari">Urutkan: Hari</option>
					<option value="jam">Urutkan: Jam</option>
				</select>
			</div>
		</div>
	</div>

	{#if loading}
		<div class="loading-state">
			<div class="spinner"></div>
			<p>Memuat data...</p>
		</div>
	{:else if error}
		<div class="error-message">
			{error}
			<button class="btn-retry" onclick={fetchKelases}>Coba Lagi</button>
		</div>
	{:else if kelases.length === 0}
		<div class="empty-state">
			<p>Belum ada data kelas.</p>
		</div>
	{:else}
		<div class="groups-container">
			{#each groupedKelases as group}
				<div class="group-section">
					<h3 class="group-title">{group.programStudiName}</h3>
					<div class="table-container">
						<table class="data-table">
							<thead>
								<tr>
									<th>Nama Kelas</th>
									<th>Kapasitas</th>
									<th>Jadwal</th>
									<th>Status</th>
									{#if authState.role === 'admin' || authState.role === 'dosen'}
										<th>Aksi</th>
									{/if}
								</tr>
							</thead>
							<tbody>
								{#each group.classes as k}
									<tr>
										<td><strong>{k.name}</strong></td>
										<td>
											<span class="badge capacity-badge">
												{k.capacity} Mhs
											</span>
										</td>
										<td>
											<div class="jadwal-info">
												<span class="hari">{k.hari}</span>
												<span class="jam">{k.jam_mulai} - {k.jam_selesai}</span>
											</div>
										</td>
										<td>
											{#if k.pengajuan && k.pengajuan.some(p => p.status === 'approved')}
												{@const approvedP = k.pengajuan.find(p => p.status === 'approved')}
												<span class="badge" style="background: rgba(239, 68, 68, 0.1); color: var(--error-color);">
													Diambil oleh {approvedP?.dosen?.name || 'Dosen'}
												</span>
											{:else if authState.role === 'dosen' && k.pengajuan && k.pengajuan.some(p => p.status === 'pending' && p.dosen_id === authState.profile?.id)}
												<span class="badge" style="background: rgba(245, 158, 11, 0.1); color: #f59e0b;">
													Menunggu Persetujuan
												</span>
											{:else}
												<span class="badge" style="background: rgba(16, 185, 129, 0.1); color: #10b981;">
													Tersedia
												</span>
											{/if}
										</td>
										{#if authState.role === 'admin'}
											<td>
												<div class="controls" style="gap: 8px;">
													{#if k.pengajuan && k.pengajuan.some(p => p.status === 'approved')}
														<a href="/dashboard/kelas/{k.id}" class="btn-request" style="text-decoration: none;">
															Detail
														</a>
													{/if}
													<button class="btn-delete" aria-label="Hapus Kelas" onclick={() => promptDelete(k)}>
														Hapus
													</button>
												</div>
											</td>
										{:else if authState.role === 'dosen'}
											<td>
												{#if k.pengajuan && k.pengajuan.some(p => p.status === 'approved')}
													<button class="btn-request disabled" disabled>
														Tidak Tersedia
													</button>
												{:else if k.pengajuan && k.pengajuan.some(p => p.status === 'pending' && p.dosen_id === authState.profile?.id)}
													<button class="btn-request disabled" disabled>
														Requested
													</button>
												{:else if myActiveKelasCount >= myApprovedMataKuliah.length}
													<button class="btn-request disabled" disabled title="Batas pengambilan kelas telah tercapai">
														Limit Tercapai
													</button>
												{:else}
													<button class="btn-request" onclick={() => promptRequest(k)}>
														Request
													</button>
												{/if}
											</td>
										{/if}
									</tr>
								{/each}
							</tbody>
						</table>
					</div>
				</div>
			{:else}
				<div class="empty-state">
					<p>Tidak ada data kelas yang sesuai pencarian.</p>
				</div>
			{/each}
		</div>
	{/if}
</div>

<PromptCodeModal
	bind:isOpen={isPromptOpen}
	title={promptAction === 'delete' ? "Hapus Kelas" : "Request Kelas"}
	message={promptAction === 'delete' ? `Untuk menghapus kelas ${selectedKelas?.name || 'ini'}, masukkan kode berikut:` : `Masukkan kode berikut untuk mengajukan request kelas ${selectedKelas?.name || ''}:`}
	{expectedCode}
	disableConfirm={promptAction === 'request' && availableMataKuliahForKelas.length === 0}
	onConfirm={handleConfirmCode}
	onCancel={() => { isPromptOpen = false; selectedKelas = null; }}
>
	{#snippet children()}
		{#if promptAction === 'request'}
			<div class="mk-selector">
				<label for="mk-select">Pilih Mata Kuliah:</label>
				{#if availableMataKuliahForKelas.length > 0}
					<select id="mk-select" bind:value={selectedMataKuliahId} class="form-input">
						{#each availableMataKuliahForKelas as mk}
							<option value={mk.mata_kuliah_id}>{mk.mata_kuliah?.name}</option>
						{/each}
					</select>
				{:else}
					<div style="color: var(--error-color); font-size: 0.9rem; margin-top: 4px;">
						Anda tidak memiliki mata kuliah yang disetujui di program studi kelas ini.
					</div>
				{/if}
			</div>
		{/if}
	{/snippet}
</PromptCodeModal>

<style>
	.mk-selector {
		margin-bottom: 16px;
		display: flex;
		flex-direction: column;
		gap: 6px;
	}
	.mk-selector label {
		font-size: 0.9rem;
		font-weight: 500;
		color: var(--text-main);
	}
	.card {
		background: var(--surface-light);
		border-radius: var(--radius-lg);
		border: 1px solid var(--surface-border);
		overflow: hidden;
	}

	.card-header {
		padding: 20px 24px;
		border-bottom: 1px solid var(--surface-border);
		background: rgba(255, 255, 255, 0.02);
	}

	.header-content {
		display: flex;
		justify-content: space-between;
		align-items: center;
		flex-wrap: wrap;
		gap: 16px;
	}

	.card-header h2 {
		font-size: 1.25rem;
		font-weight: 600;
		margin: 0;
	}

	.controls {
		display: flex;
		gap: 12px;
		flex-wrap: wrap;
	}

	.search-input, .sort-select {
		padding: 8px 12px;
		border: 1px solid var(--surface-border);
		border-radius: var(--radius-sm);
		background: var(--surface-light);
		color: var(--text-main);
		font-size: 0.9rem;
	}

	.search-input {
		width: 200px;
	}

	.search-input:focus, .sort-select:focus {
		outline: none;
		border-color: var(--primary-color);
	}

	.groups-container {
		display: flex;
		flex-direction: column;
		gap: 0;
	}

	.group-section {
		border-bottom: 2px solid var(--surface-border);
	}
	
	.group-section:last-child {
		border-bottom: none;
	}

	.group-title {
		padding: 16px 24px;
		margin: 0;
		background: rgba(var(--primary-color-rgb), 0.03);
		color: var(--primary-color);
		font-size: 1rem;
		font-weight: 600;
	}

	.table-container {
		overflow-x: auto;
	}

	.data-table {
		width: 100%;
		border-collapse: collapse;
		text-align: left;
	}

	.data-table th,
	.data-table td {
		padding: 16px 24px;
		border-bottom: 1px solid var(--surface-border);
	}

	.data-table th {
		font-weight: 600;
		color: var(--text-muted);
		font-size: 0.85rem;
		text-transform: uppercase;
		letter-spacing: 0.5px;
		background: rgba(255, 255, 255, 0.02);
	}

	.data-table tr:last-child td {
		border-bottom: none;
	}

	.data-table tbody tr {
		transition: background-color 0.2s;
	}

	.data-table tbody tr:hover {
		background: rgba(255, 255, 255, 0.03);
	}

	.capacity-badge {
		background: rgba(16, 185, 129, 0.1);
		color: #10b981;
		padding: 4px 10px;
		border-radius: 12px;
		font-size: 0.8rem;
		font-weight: 600;
		border: 1px solid rgba(16, 185, 129, 0.2);
	}

	.jadwal-info {
		display: flex;
		flex-direction: column;
		gap: 4px;
	}

	.hari {
		font-weight: 600;
		color: var(--text-main);
	}

	.jam {
		font-size: 0.85rem;
		color: var(--text-muted);
		background: rgba(0, 0, 0, 0.04);
		padding: 2px 6px;
		border-radius: 4px;
		display: inline-block;
		width: max-content;
	}

	.btn-delete {
		background: rgba(239, 68, 68, 0.1);
		color: var(--error-color);
		border: 1px solid rgba(239, 68, 68, 0.2);
		padding: 6px 12px;
		border-radius: var(--radius-sm);
		cursor: pointer;
		font-size: 0.85rem;
		font-weight: 500;
		transition: all 0.2s;
	}

	.btn-delete:hover {
		background: rgba(239, 68, 68, 0.2);
	}

	.btn-request {
		background: rgba(79, 70, 229, 0.1);
		color: var(--primary-color);
		border: 1px solid rgba(79, 70, 229, 0.2);
		padding: 6px 12px;
		border-radius: var(--radius-sm);
		cursor: pointer;
		font-size: 0.85rem;
		font-weight: 500;
		transition: all 0.2s;
	}

	.btn-request:hover {
		background: rgba(79, 70, 229, 0.2);
	}

	.btn-request.disabled {
		background: rgba(0, 0, 0, 0.05);
		color: var(--text-muted);
		border-color: transparent;
		cursor: not-allowed;
	}

	.loading-state,
	.empty-state {
		padding: 48px;
		text-align: center;
		color: var(--text-muted);
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 16px;
	}

	.spinner {
		width: 32px;
		height: 32px;
		border: 3px solid rgba(0, 0, 0, 0.1);
		border-top-color: var(--primary-color);
		border-radius: 50%;
		animation: spin 1s linear infinite;
	}

	.empty-icon {
		font-size: 3rem;
		opacity: 0.5;
	}

	.error-message {
		margin: 24px;
		padding: 16px;
		background: rgba(239, 68, 68, 0.1);
		border: 1px solid rgba(239, 68, 68, 0.2);
		border-radius: var(--radius-md);
		color: var(--error-color);
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 12px;
	}

	.btn-retry {
		background: var(--error-color);
		color: white;
		border: none;
		padding: 6px 16px;
		border-radius: var(--radius-sm);
		cursor: pointer;
		font-weight: 500;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}
</style>
