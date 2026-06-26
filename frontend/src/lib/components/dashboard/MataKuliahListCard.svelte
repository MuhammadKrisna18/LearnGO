<script lang="ts">
	import { onMount } from 'svelte';
	import { matakuliahService } from '$lib/services/matakuliah';
	import { programStudiService } from '$lib/services/programstudi';
	import type { MataKuliah, ProgramStudi } from '$lib/types';
	import DeleteConfirmModal from './DeleteConfirmModal.svelte';

	let mkList: MataKuliah[] = $state([]);
	let prodiList: ProgramStudi[] = $state([]);
	let loading = $state(true);
	let error = $state('');

	onMount(async () => {
		try {
			const [resMk, resProdi] = await Promise.all([
				matakuliahService.getList(),
				programStudiService.getList()
			]);
			
			if (resMk.success && resMk.data) {
				mkList = resMk.data;
			} else {
				error = resMk.message || 'Gagal mengambil daftar mata kuliah';
			}
			
			if (resProdi.success && resProdi.data) {
				prodiList = resProdi.data;
			}
		} catch (err) {
			error = 'Gagal terhubung ke server';
			console.error(err);
		} finally {
			loading = false;
		}
	});

	let isDeleteModalOpen = $state(false);
	let selectedMK = $state<MataKuliah | null>(null);

	function promptDelete(mk: MataKuliah) {
		selectedMK = mk;
		isDeleteModalOpen = true;
	}

	async function handleDelete() {
		if (!selectedMK) return;
		
		try {
			const res = await matakuliahService.delete(selectedMK.id);
			if (res.success) {
				// Remove from list
				mkList = mkList.filter(m => m.id !== selectedMK!.id);
			} else {
				error = res.message || 'Gagal menghapus mata kuliah';
			}
		} catch (err) {
			error = 'Gagal terhubung ke server';
		}
	}
</script>

<div class="mk-list-card glass-panel animate-fade-in" style="animation-delay: 0.4s;">
	<div class="card-header">
		<h3>Daftar Mata Kuliah</h3>
		<p>Daftar seluruh mata kuliah yang tersedia di dalam sistem.</p>
	</div>

	{#if loading}
		<div class="table-loading">Memuat data...</div>
	{:else if error}
		<div class="table-error">{error}</div>
	{:else if mkList.length === 0}
		<div class="table-empty">Belum ada mata kuliah yang didaftarkan.</div>
	{:else}
		<div class="prodi-groups">
			{#each prodiList as prodi}
				{@const prodiMks = mkList.filter(mk => mk.program_studi_id === prodi.id)}
				<div class="prodi-section">
					<h4 class="prodi-title">{prodi.name} ({prodi.code})</h4>
					
					{#if prodiMks.length === 0}
						<div class="table-empty prodi-empty">Belum ada mata kuliah.</div>
					{:else}
						<div class="table-container">
							<table class="mk-table">
								<thead>
									<tr>
										<th>Nama Mata Kuliah</th>
										<th>SKS</th>
										<th>Tgl Ditambahkan</th>
										<th>Aksi</th>
									</tr>
								</thead>
								<tbody>
									{#each prodiMks as mk}
										<tr>
											<td>{mk.name}</td>
											<td><span class="badge sks-badge">{mk.sks} SKS</span></td>
											<td>{new Date(mk.created_at).toLocaleDateString()}</td>
											<td>
												<button class="btn-delete" aria-label="Hapus Mata Kuliah" onclick={() => promptDelete(mk)}>
													Hapus
												</button>
											</td>
										</tr>
									{/each}
								</tbody>
							</table>
						</div>
					{/if}
				</div>
			{/each}
		</div>
	{/if}
</div>

<DeleteConfirmModal 
	bind:isOpen={isDeleteModalOpen}
	title="Hapus Mata Kuliah"
	itemName={selectedMK?.name || ''}
	onConfirm={handleDelete}
	onCancel={() => selectedMK = null}
/>

<style>
	.mk-list-card {
		padding: 32px;
		border-radius: var(--radius-lg);
		margin-top: 24px;
	}

	.card-header {
		margin-bottom: 24px;
	}

	.card-header h3 {
		font-size: 1.25rem;
		font-weight: 600;
		color: var(--text-main);
		margin-bottom: 4px;
	}

	.card-header p {
		font-size: 0.9rem;
		color: var(--text-muted);
	}

	.table-container {
		overflow-x: auto;
		border: 1px solid var(--surface-border);
		border-radius: var(--radius-md);
	}

	.mk-table {
		width: 100%;
		border-collapse: collapse;
		text-align: left;
	}

	.mk-table th {
		background: rgba(241, 245, 249, 0.5);
		padding: 16px;
		font-size: 0.85rem;
		text-transform: uppercase;
		letter-spacing: 0.5px;
		color: var(--text-muted);
		border-bottom: 1px solid var(--surface-border);
	}

	.mk-table td {
		padding: 16px;
		border-bottom: 1px solid var(--surface-border);
		font-size: 0.95rem;
	}

	.mk-table tr:last-child td {
		border-bottom: none;
	}

	.mk-table tr:hover {
		background: rgba(248, 250, 252, 0.5);
	}

	.table-loading,
	.table-empty,
	.table-error {
		padding: 32px;
		text-align: center;
		color: var(--text-muted);
		background: rgba(248, 250, 252, 0.5);
		border-radius: var(--radius-md);
		border: 1px dashed var(--surface-border);
	}
	
	.table-error {
		color: var(--error-color);
		background: rgba(239, 68, 68, 0.05);
		border-color: rgba(239, 68, 68, 0.2);
	}

	.badge {
		font-size: 0.75rem;
		padding: 4px 10px;
		border-radius: 12px;
		font-weight: 600;
		letter-spacing: 0.5px;
	}

	.sks-badge {
		background: rgba(16, 185, 129, 0.1);
		color: #10b981;
		border: 1px solid rgba(16, 185, 129, 0.2);
	}
	
	.prodi-groups {
		display: flex;
		flex-direction: column;
		gap: 32px;
	}
	
	.prodi-title {
		font-size: 1.1rem;
		font-weight: 600;
		color: var(--primary-color);
		margin-bottom: 16px;
		padding-bottom: 8px;
		border-bottom: 2px solid rgba(99, 102, 241, 0.2);
	}
	
	.prodi-empty {
		padding: 20px;
		background: transparent;
		border: 1px dashed var(--surface-border);
	}
	
	.btn-delete {
		background: rgba(239, 68, 68, 0.1);
		color: var(--error-color);
		border: 1px solid rgba(239, 68, 68, 0.2);
		cursor: pointer;
		transition: all 0.2s;
		font-size: 0.85rem;
		padding: 6px 12px;
		border-radius: var(--radius-sm);
		font-weight: 500;
	}
	
	.btn-delete:hover {
		background: rgba(239, 68, 68, 0.2);
	}
</style>
