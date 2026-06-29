<script lang="ts">
	import { onMount } from 'svelte';
	import { kelasService } from '$lib/services/kelas';
	import type { PengajuanKelas } from '$lib/types';
	import Card from '$lib/components/ui/Card.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';
	import PromptCodeModal from '$lib/components/ui/PromptCodeModal.svelte';
	import { toast } from '$lib/stores/toast.svelte';

	let requests = $state<PengajuanKelas[]>([]);
	let loading = $state(true);
	let error = $state('');

	let isPromptOpen = $state(false);
	let promptTitle = $state("");
	let promptMessage = $state("");
	let expectedCode = $state("");
	let pendingAction = $state<"approve" | "reject" | null>(null);
	let pendingId = $state("");

	async function loadRequests() {
		loading = true;
		error = '';
		try {
			const res = await kelasService.getAllPengajuan();
			if (res.success && res.data) {
				requests = res.data.filter((r: PengajuanKelas) => r.status === 'pending');
			}
		} catch (err: any) {
			error = err.message || 'Gagal memuat pengajuan';
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		loadRequests();
	});

	function openApprovePrompt(id: string) {
		pendingId = id;
		pendingAction = "approve";
		expectedCode = Math.floor(100000 + Math.random() * 900000).toString();
		promptTitle = "Konfirmasi Terima Pengajuan";
		promptMessage = "Untuk menerima pengajuan, masukkan kode berikut:";
		isPromptOpen = true;
	}

	function openRejectPrompt(id: string) {
		pendingId = id;
		pendingAction = "reject";
		expectedCode = Math.floor(100000 + Math.random() * 900000).toString();
		promptTitle = "Konfirmasi Tolak Pengajuan";
		promptMessage = "Untuk menolak pengajuan, masukkan kode berikut:";
		isPromptOpen = true;
	}

	async function handlePromptConfirm(code: string) {
		if (code !== expectedCode) {
			toast.error('Kode tidak cocok. Aksi dibatalkan.');
			return;
		}

		if (pendingAction === "approve") {
			try {
				const res = await kelasService.approvePengajuan(pendingId);
				if (res.success) {
					toast.success('Pengajuan berhasil disetujui');
					await loadRequests();
				} else {
					toast.error(res.message);
				}
			} catch (err: any) {
				toast.error(err.message || 'Gagal menyetujui pengajuan');
			}
		} else if (pendingAction === "reject") {
			try {
				const res = await kelasService.rejectPengajuan(pendingId);
				if (res.success) {
					toast.success('Pengajuan berhasil ditolak');
					await loadRequests();
				} else {
					toast.error(res.message);
				}
			} catch (err: any) {
				toast.error(err.message || 'Gagal menolak pengajuan');
			}
		}
	}

	function getBadgeType(status: string) {
		switch(status) {
			case 'approved': return 'success';
			case 'rejected': return 'error';
			default: return 'warning';
		}
	}
</script>

<Card style="padding: 24px; margin-top: 24px;">
	<div class="card-header-flex">
		<h3>Daftar Pengajuan Kelas</h3>
		<button class="btn-secondary" onclick={loadRequests} disabled={loading}>
			Refresh
		</button>
	</div>

	{#if error}
		<div class="state-container state-error" style="padding: 16px; margin-bottom: 16px;">{error}</div>
	{/if}

	{#if loading}
		<div class="state-container">Memuat data...</div>
	{:else if requests.length === 0}
		<div class="state-container">Tidak ada pengajuan kelas yang menunggu persetujuan.</div>
	{:else}
		<div class="data-table-container">
			<table class="data-table">
				<thead>
					<tr>
						<th>Dosen</th>
						<th>Kelas</th>
						<th>Status</th>
						<th>Tanggal</th>
						<th>Aksi</th>
					</tr>
				</thead>
				<tbody>
					{#each requests as req}
						<tr>
							<td>
								<div class="fw-medium">{req.dosen?.name || 'Unknown'}</div>
								<div class="txt-sm">{req.dosen?.email || ''}</div>
							</td>
							<td>
								<div class="fw-medium">
									{req.kelas?.name || 'Unknown'}
								</div>
								<div class="txt-sm">{req.kelas?.hari || ''} {req.kelas?.jam_mulai || ''}-{req.kelas?.jam_selesai || ''}</div>
								{#if req.mata_kuliah}
									<div class="txt-sm" style="color: var(--primary-color);">MK: {req.mata_kuliah.name}</div>
								{/if}
							</td>
							<td>
								<Badge type={getBadgeType(req.status)}>
									{req.status.toUpperCase()}
								</Badge>
							</td>
							<td>{new Date(req.created_at).toLocaleDateString('id-ID')}</td>
							<td>
								<div class="action-buttons">
									<button class="btn-approve" onclick={() => openApprovePrompt(req.id)}>Terima</button>
									<button class="btn-reject" onclick={() => openRejectPrompt(req.id)}>Tolak</button>
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</Card>

<PromptCodeModal
	bind:isOpen={isPromptOpen}
	title={promptTitle}
	message={promptMessage}
	{expectedCode}
	onConfirm={handlePromptConfirm}
/>

<style>
	.btn-approve {
		background: #10b981;
		color: white;
		border: none;
		padding: 6px 12px;
		border-radius: var(--radius-sm);
		font-size: 0.75rem;
		font-weight: 500;
		cursor: pointer;
	}
	.btn-approve:hover { background: #059669; }

	.btn-reject {
		background: #ef4444;
		color: white;
		border: none;
		padding: 6px 12px;
		border-radius: var(--radius-sm);
		font-size: 0.75rem;
		font-weight: 500;
		cursor: pointer;
	}
	.btn-reject:hover { background: #dc2626; }

	.card-header-flex {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 24px;
	}

	.card-header-flex h3 {
		font-size: 1.25rem;
		font-weight: 600;
		color: var(--text-main);
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

	.action-buttons {
		display: flex;
		gap: 8px;
	}
</style>
