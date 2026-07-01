<script lang="ts">
	import { onMount } from 'svelte';
	import { dosenService } from '$lib/services/dosen';
	import { authService } from '$lib/services/auth';
	import { programStudiService } from '$lib/services/programstudi';
	import type { UserProfile, ProgramStudi } from '$lib/types';
	import Card from '$lib/components/ui/Card.svelte';
	import Modal from '$lib/components/ui/Modal.svelte';
	import PromptCodeModal from '$lib/components/ui/PromptCodeModal.svelte';
	import { toast } from '$lib/stores/toast.svelte';

	let dosenList: UserProfile[] = $state([]);
	let prodiList: ProgramStudi[] = $state([]);
	let loading = $state(true);
	let error = $state('');

	onMount(async () => {
		try {
			const [dosenRes, prodiRes] = await Promise.all([
				dosenService.getList(),
				programStudiService.getList()
			]);

			if (dosenRes.success && dosenRes.data) {
				dosenList = dosenRes.data;
			} else {
				error = dosenRes.message || 'Gagal mengambil daftar dosen';
			}

			if (prodiRes.success && prodiRes.data) {
				prodiList = prodiRes.data;
			}

		} catch (err) {
			error = 'Gagal terhubung ke server';
			console.error(err);
		} finally {
			loading = false;
		}
	});

	function getDosenByProdi(prodiId: string) {
		return dosenList.filter(d => d.program_studi_id === prodiId);
	}

	function getDosenWithoutProdi() {
		return dosenList.filter(d => !d.program_studi_id);
	}

	let isPromptOpen = $state(false);
	let selectedDosen = $state<UserProfile | null>(null);
	let expectedCode = $state("");

	function promptDelete(dosen: UserProfile) {
		selectedDosen = dosen;
		expectedCode = Math.floor(100000 + Math.random() * 900000).toString();
		isPromptOpen = true;
	}

	async function handleDelete(code: string) {
		if (!selectedDosen) return;

		if (code !== expectedCode) {
			toast.error('Kode tidak cocok. Aksi dibatalkan.');
			selectedDosen = null;
			isPromptOpen = false;
			return;
		}

		try {
			const res = await authService.deleteDosen(selectedDosen.id);
			if (res.success) {
				toast.success('Akun dosen berhasil dihapus');
				dosenList = dosenList.filter(d => d.id !== selectedDosen!.id);
				isPromptOpen = false;
			} else {
				toast.error(res.message || 'Gagal menghapus dosen');
			}
		} catch (err) {
			toast.error('Gagal terhubung ke server');
		}
	}

	let isDetailModalOpen = $state(false);
	let detailDosen = $state<UserProfile | null>(null);
	let detailProdiName = $state('');

	function viewDetail(dosen: UserProfile, prodi?: ProgramStudi) {
		detailDosen = dosen;
		detailProdiName = prodi ? prodi.name : '';
		isDetailModalOpen = true;
	}
</script>

<Card class="dosen-list-card animate-fade-in" style="animation-delay: 0.3s; padding: 32px; margin-top: 24px;">
	<div class="card-header" style="margin-bottom: 24px;">
		<h3 style="font-size: 1.25rem; font-weight: 600; color: var(--text-main); margin-bottom: 4px;">Daftar Akun Dosen</h3>
		<p class="text-muted" style="font-size: 0.9rem;">Daftar seluruh dosen yang terdaftar di dalam sistem.</p>
	</div>

	{#if loading}
		<div class="state-container">Memuat data...</div>
	{:else if error}
		<div class="state-container state-error">{error}</div>
	{:else}
		{#each prodiList as prodi}
			{@const prodiDosen = getDosenByProdi(prodi.id)}
			<div class="prodi-section">
				<h4 class="prodi-title">{prodi.name}</h4>
				{#if prodiDosen.length === 0}
					<div class="state-container" style="padding: 20px; background: transparent;">Belum ada akun dosen di program studi ini.</div>
				{:else}
					<div class="data-table-container">
						<table class="data-table">
							<thead>
								<tr>
									<th>Nama Dosen</th>
									<th>Email</th>
									<th>Tgl Bergabung</th>
									<th>Aksi</th>
								</tr>
							</thead>
							<tbody>
								{#each prodiDosen as dosen}
									<tr>
										<td>
											<div class="dosen-info">
												<div class="avatar">
													{#if dosen.photo_url}
														<img src={`http://localhost:8080${dosen.photo_url}`} alt={dosen.name} class="avatar-img-small" />
													{:else}
														{dosen.name.charAt(0).toUpperCase()}
													{/if}
												</div>
												<span class="font-medium">{dosen.name}</span>
											</div>
										</td>
										<td><span class="mono">{dosen.email}</span></td>
										<td>{new Date(dosen.created_at).toLocaleDateString()}</td>
										<td>
											<div class="action-buttons">
												<button class="btn-detail" aria-label="Lihat Detail Dosen" onclick={() => viewDetail(dosen, prodi)}>
													Detail
												</button>
												<button class="btn-delete" aria-label="Hapus Dosen" onclick={() => promptDelete(dosen)}>
													Hapus
												</button>
											</div>
										</td>
									</tr>
								{/each}
							</tbody>
						</table>
					</div>
				{/if}
			</div>
		{/each}

		{@const noProdiDosen = getDosenWithoutProdi()}
		{#if noProdiDosen.length > 0}
			<div class="prodi-section">
				<h4 class="prodi-title unassigned">Belum Memiliki Program Studi</h4>
				<div class="data-table-container">
					<table class="data-table">
						<thead>
							<tr>
								<th>Nama Dosen</th>
								<th>Email</th>
								<th>Tgl Bergabung</th>
								<th>Aksi</th>
							</tr>
						</thead>
						<tbody>
							{#each noProdiDosen as dosen}
								<tr>
									<td>
										<div class="dosen-info">
											<div class="avatar">
												{#if dosen.photo_url}
													<img src={`http://localhost:8080${dosen.photo_url}`} alt={dosen.name} class="avatar-img-small" />
												{:else}
													{dosen.name.charAt(0).toUpperCase()}
												{/if}
											</div>
											<span class="font-medium">{dosen.name}</span>
										</div>
									</td>
									<td><span class="mono">{dosen.email}</span></td>
									<td>{new Date(dosen.created_at).toLocaleDateString()}</td>
									<td>
										<div class="action-buttons">
											<button class="btn-detail" aria-label="Lihat Detail Dosen" onclick={() => viewDetail(dosen)}>
												Detail
											</button>
											<button class="btn-delete" aria-label="Hapus Dosen" onclick={() => promptDelete(dosen)}>
												Hapus
											</button>
										</div>
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			</div>
		{/if}
	{/if}
</Card>

<PromptCodeModal 
	bind:isOpen={isPromptOpen}
	title="Hapus Akun Dosen"
	message={`Untuk menghapus akun dosen ${selectedDosen?.name || ''}, masukkan kode berikut:`}
	{expectedCode}
	onConfirm={handleDelete}
	onCancel={() => selectedDosen = null}
/>

{#if isDetailModalOpen && detailDosen}
	<Modal bind:isOpen={isDetailModalOpen} title="Detail Profil Dosen">
		<div class="detail-avatar">
			{#if detailDosen.photo_url}
				<img src={`http://localhost:8080${detailDosen.photo_url}`} alt={detailDosen.name} class="avatar-img-large" />
			{:else}
				{detailDosen.name.charAt(0).toUpperCase()}
			{/if}
		</div>
		<div class="detail-group">
			<span class="detail-label">Nama Lengkap</span>
			<span class="detail-value">{detailDosen.name}</span>
		</div>
		<div class="detail-group">
			<span class="detail-label">Nama Panggilan</span>
			<span class="detail-value">{detailDosen.nickname || '-'}</span>
		</div>
		<div class="detail-group">
			<span class="detail-label">Nomor Induk Dosen (NID)</span>
			<span class="detail-value">
				{#if detailDosen.nid}
					<span class="mono">{detailDosen.nid}</span>
				{:else}
					-
				{/if}
			</span>
		</div>
		<div class="detail-group">
			<span class="detail-label">Email</span>
			<span class="detail-value mono">{detailDosen.email}</span>
		</div>
		<div class="detail-group">
			<span class="detail-label">Program Studi</span>
			<span class="detail-value">{detailProdiName || 'Belum ada program studi'}</span>
		</div>
		<div class="detail-group">
			<span class="detail-label">Tanggal Bergabung</span>
			<span class="detail-value">{new Date(detailDosen.created_at).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric' })}</span>
		</div>

		{#snippet footer()}
			<button class="btn-secondary" onclick={() => isDetailModalOpen = false}>Tutup</button>
		{/snippet}
	</Modal>
{/if}

<style>
	.prodi-section {
		margin-bottom: 32px;
	}

	.prodi-section:last-child {
		margin-bottom: 0;
	}

	.prodi-title {
		font-size: 1.1rem;
		font-weight: 600;
		color: var(--primary-color);
		margin-bottom: 12px;
		display: flex;
		align-items: center;
		gap: 8px;
	}

	.prodi-title::before {
		content: '';
		display: block;
		width: 4px;
		height: 16px;
		background: var(--primary-color);
		border-radius: 4px;
	}

	.prodi-title.unassigned {
		color: var(--warning-color, #f59e0b);
	}

	.prodi-title.unassigned::before {
		background: var(--warning-color, #f59e0b);
	}
	
	.dosen-info {
		display: flex;
		align-items: center;
		gap: 12px;
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
	}

	.avatar-img-small {
		width: 100%;
		height: 100%;
		object-fit: cover;
		aspect-ratio: 1/1;
	}

	
	.detail-avatar {
		width: 64px;
		height: 64px;
		background: linear-gradient(135deg, var(--primary-color), var(--primary-hover));
		color: white;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		font-weight: 700;
		font-size: 1.75rem;
		margin: 0 auto 16px auto;
		box-shadow: 0 4px 10px rgba(79, 70, 229, 0.3);
		overflow: hidden;
	}

	.avatar-img-large {
		width: 100%;
		height: 100%;
		object-fit: cover;
		aspect-ratio: 1/1;
	}

	.detail-group {
		display: flex;
		flex-direction: column;
		gap: 4px;
	}

	.detail-label {
		font-size: 0.8rem;
		font-weight: 500;
		color: var(--text-muted);
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}

	.detail-value {
		font-size: 1rem;
		color: var(--text-main);
		font-weight: 500;
	}
</style>
