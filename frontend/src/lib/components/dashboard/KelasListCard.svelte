<script lang="ts">
	import { onMount } from 'svelte';
	import { kelasService } from '$lib/services/kelas';
	import { authState } from '$lib/stores/auth.svelte';
	import type { Kelas } from '$lib/types';
	import DeleteConfirmModal from './DeleteConfirmModal.svelte';

	let kelases = $state<Kelas[]>([]);
	let loading = $state(true);
	let error = $state('');

	let isDeleteModalOpen = $state(false);
	let selectedKelas = $state<Kelas | null>(null);

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
		// Listen for event from register card
		window.addEventListener('kelasCreated', fetchKelases);
		return () => {
			window.removeEventListener('kelasCreated', fetchKelases);
		};
	});

	function promptDelete(k: Kelas) {
		selectedKelas = k;
		isDeleteModalOpen = true;
	}

	async function handleDelete() {
		if (!selectedKelas) return;
		try {
			const res = await kelasService.delete(selectedKelas.id);
			if (res.success) {
				kelases = kelases.filter((k) => k.id !== selectedKelas!.id);
			} else {
				alert(res.message || 'Gagal menghapus kelas');
			}
		} catch (err) {
			alert('Terjadi kesalahan saat menghapus kelas');
		} finally {
			isDeleteModalOpen = false;
			selectedKelas = null;
		}
	}
</script>

<div class="card glass-panel animate-fade-in">
	<div class="card-header">
		<h2>Daftar Kelas</h2>
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
		<div class="table-container">
			<table class="data-table">
				<thead>
					<tr>
						<th>Nama Kelas</th>
						<th>Kapasitas</th>
						<th>Program Studi</th>
						{#if authState.role === 'admin'}
							<th>Aksi</th>
						{/if}
					</tr>
				</thead>
				<tbody>
					{#each kelases as k}
						<tr>
							<td><strong>{k.name}</strong></td>
							<td>
								<span class="badge capacity-badge">
									{k.capacity} Mhs
								</span>
							</td>
							<td>{k.program_studi ? k.program_studi.name : '-'}</td>
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
	{/if}
</div>

<DeleteConfirmModal
	bind:isOpen={isDeleteModalOpen}
	title="Hapus Kelas"
	message="Apakah Anda yakin ingin menghapus kelas {selectedKelas?.name}? Tindakan ini tidak dapat dibatalkan."
	onConfirm={handleDelete}
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

	.card-header h2 {
		font-size: 1.25rem;
		font-weight: 600;
		margin: 0;
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
