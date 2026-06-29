<script lang="ts">
	import { onMount } from 'svelte';
	import { kelasService } from '$lib/services/kelas';
	import type { PengajuanKelas } from '$lib/types';
	import Card from '$lib/components/ui/Card.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';

	let requests = $state<PengajuanKelas[]>([]);
	let loading = $state(true);
	let error = $state('');

	let searchQuery = $state('');

	let filteredRequests = $derived.by(() => {
		let filtered = requests.filter((r: PengajuanKelas) => r.status === 'approved');
		if (searchQuery.trim() !== '') {
			const q = searchQuery.toLowerCase();
			filtered = filtered.filter(r => 
				r.dosen?.name.toLowerCase().includes(q) ||
				r.dosen?.email.toLowerCase().includes(q) ||
				r.kelas?.name.toLowerCase().includes(q) ||
				r.mata_kuliah?.name.toLowerCase().includes(q)
			);
		}
		return filtered;
	});

	async function loadRequests() {
		loading = true;
		error = '';
		try {
			const res = await kelasService.getAllPengajuan();
			if (res.success && res.data) {
				requests = res.data;
			}
		} catch (err: any) {
			error = err.message || 'Gagal memuat data pengajuan kelas';
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		loadRequests();
	});
</script>

<Card style="padding: 24px; margin-top: 24px;">
	<div class="card-header-flex">
		<div>
			<h3>Daftar Dosen & Kelas Yang Diampu</h3>
			<p class="header-subtitle">Daftar dosen yang telah disetujui untuk mengajar kelas tertentu</p>
		</div>
		<button class="btn-secondary" onclick={loadRequests} disabled={loading}>
			Refresh
		</button>
	</div>

	{#if error}
		<div class="state-container state-error" style="padding: 16px; margin-bottom: 16px;">{error}</div>
	{/if}

	<div class="search-container">
		<input 
			type="text" 
			class="form-input search-input" 
			placeholder="Cari dosen, mata kuliah, atau kelas..." 
			bind:value={searchQuery}
		/>
	</div>

	{#if loading}
		<div class="state-container">Memuat data...</div>
	{:else if filteredRequests.length === 0}
		<div class="state-container">Tidak ada data dosen yang mengampu kelas.</div>
	{:else}
		<div class="data-table-container">
			<table class="data-table">
				<thead>
					<tr>
						<th>Dosen</th>
						<th>Mata Kuliah</th>
						<th>Kelas & Program Studi</th>
						<th>Jadwal</th>
						<th>Kapasitas</th>
					</tr>
				</thead>
				<tbody>
					{#each filteredRequests as req}
						<tr>
							<td>
								<div class="fw-medium">{req.dosen?.name || 'Unknown'}</div>
								<div class="txt-sm">{req.dosen?.email || ''}</div>
							</td>
							<td>
								<div class="fw-medium">
									{req.mata_kuliah?.name || 'Unknown'}
								</div>
								<div class="txt-sm">{req.mata_kuliah?.sks || '-'} SKS</div>
							</td>
							<td>
								<div class="fw-medium">
									{req.kelas?.name || 'Unknown'}
								</div>
								{#if req.kelas?.program_studi}
									<div class="txt-sm" style="color: var(--primary-color);">
										{req.kelas.program_studi.name}
									</div>
								{/if}
							</td>
							<td>
								<Badge type="info">
									{req.kelas?.hari || '-'}
								</Badge>
								<div class="txt-sm" style="margin-top: 6px;">
									{req.kelas?.jam_mulai || ''} - {req.kelas?.jam_selesai || ''}
								</div>
							</td>
							<td>
								<div class="fw-medium">{req.kelas?.capacity || 0} Mahasiswa</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</Card>

<style>
	.card-header-flex {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		margin-bottom: 16px;
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

	.fw-medium {
		font-weight: 500;
		color: var(--text-main);
	}

	.txt-sm {
		font-size: 0.8rem;
		color: var(--text-muted);
		margin-top: 4px;
	}

	.search-container {
		margin-bottom: 20px;
	}
	
	.search-input {
		width: 100%;
		max-width: 400px;
	}
</style>
