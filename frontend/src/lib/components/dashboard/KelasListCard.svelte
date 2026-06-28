<script lang="ts">
	import { onMount } from 'svelte';
	import { kelasService } from '$lib/services/kelas';
	import { authState } from '$lib/stores/auth.svelte';
	import type { Kelas } from '$lib/types';
	import PromptCodeModal from '$lib/components/ui/PromptCodeModal.svelte';
	import { toast } from '$lib/stores/toast.svelte';

	let kelases = $state<Kelas[]>([]);
	let loading = $state(true);
	let error = $state('');

	let isPromptOpen = $state(false);
	let selectedKelas = $state<Kelas | null>(null);
	let expectedCode = $state("");

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
		isPromptOpen = true;
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
		} catch (err) {
			toast.error('Terjadi kesalahan saat menghapus kelas');
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
			<div class="empty-icon">📝</div>
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
									{#if authState.role === 'admin'}
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
										{#if authState.role === 'admin'}
											<td>
												<button class="btn-delete" aria-label="Hapus Kelas" onclick={() => promptDelete(k)}>
													Hapus
												</button>
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
					<div class="empty-icon">🔍</div>
					<p>Tidak ada data kelas yang sesuai pencarian.</p>
				</div>
			{/each}
		</div>
	{/if}
</div>

<PromptCodeModal
	bind:isOpen={isPromptOpen}
	title="Hapus Kelas"
	message={`Untuk menghapus kelas ${selectedKelas?.name || 'ini'}, masukkan kode berikut:`}
	{expectedCode}
	onConfirm={handleDelete}
	onCancel={() => { isPromptOpen = false; selectedKelas = null; }}
/>

<style>
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
