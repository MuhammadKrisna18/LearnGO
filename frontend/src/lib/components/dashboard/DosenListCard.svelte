<script lang="ts">
	import { onMount } from 'svelte';
	import { dosenService } from '$lib/services/dosen';
	import { authService } from '$lib/services/auth';
	import { programStudiService } from '$lib/services/programstudi';
	import type { UserProfile, ProgramStudi } from '$lib/types';
	import DeleteConfirmModal from './DeleteConfirmModal.svelte';

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

	let isDeleteModalOpen = $state(false);
	let selectedDosen = $state<UserProfile | null>(null);

	function promptDelete(dosen: UserProfile) {
		selectedDosen = dosen;
		isDeleteModalOpen = true;
	}

	async function handleDelete() {
		if (!selectedDosen) return;
		
		try {
			const res = await authService.deleteDosen(selectedDosen.id);
			if (res.success) {

				dosenList = dosenList.filter(d => d.id !== selectedDosen!.id);
				isDeleteModalOpen = false;
			} else {
				error = res.message || 'Gagal menghapus dosen';
			}
		} catch (err) {
			error = 'Gagal terhubung ke server';
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

<div class="dosen-list-card glass-panel animate-fade-in" style="animation-delay: 0.3s;">
	<div class="card-header">
		<h3>Daftar Akun Dosen</h3>
		<p>Daftar seluruh dosen yang terdaftar di dalam sistem.</p>
	</div>

	{#if loading}
		<div class="table-loading">Memuat data...</div>
	{:else if error}
		<div class="table-error">{error}</div>
	{:else}
		{#each prodiList as prodi}
			{@const prodiDosen = getDosenByProdi(prodi.id)}
			<div class="prodi-section">
				<h4 class="prodi-title">{prodi.name}</h4>
				{#if prodiDosen.length === 0}
					<div class="table-empty">Belum ada akun dosen di program studi ini.</div>
				{:else}
					<div class="table-container">
						<table class="dosen-table">
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
												<div class="avatar">{dosen.name.charAt(0).toUpperCase()}</div>
												<span class="dosen-name">{dosen.name}</span>
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
				<div class="table-container">
					<table class="dosen-table">
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
											<div class="avatar">{dosen.name.charAt(0).toUpperCase()}</div>
											<span class="dosen-name">{dosen.name}</span>
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
</div>

<DeleteConfirmModal 
	bind:isOpen={isDeleteModalOpen}
	title="Hapus Akun Dosen"
	itemName={selectedDosen?.name || ''}
	onConfirm={handleDelete}
	onCancel={() => selectedDosen = null}
/>

{#if isDetailModalOpen && detailDosen}
	<div class="modal-backdrop" onclick={() => isDetailModalOpen = false}>
		<div class="modal-content glass-panel animate-fade-in" onclick={e => e.stopPropagation()}>
			<div class="modal-header">
				<h3>Detail Profil Dosen</h3>
				<button class="btn-close" onclick={() => isDetailModalOpen = false}>✕</button>
			</div>
			<div class="modal-body">
				<div class="detail-avatar">
					{detailDosen.name.charAt(0).toUpperCase()}
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
							<span class="dosen-nid">{detailDosen.nid}</span>
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
			</div>
			<div class="modal-footer">
				<button class="btn-secondary" onclick={() => isDetailModalOpen = false}>Tutup</button>
			</div>
		</div>
	</div>
{/if}

<style>
	.dosen-list-card {
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

	.table-container {
		overflow-x: auto;
		border: 1px solid var(--surface-border);
		border-radius: var(--radius-md);
	}

	.dosen-table {
		width: 100%;
		border-collapse: collapse;
		text-align: left;
	}

	.dosen-table th {
		background: rgba(241, 245, 249, 0.5);
		padding: 16px;
		font-size: 0.85rem;
		text-transform: uppercase;
		letter-spacing: 0.5px;
		color: var(--text-muted);
		border-bottom: 1px solid var(--surface-border);
	}

	.dosen-table td {
		padding: 16px;
		border-bottom: 1px solid var(--surface-border);
		font-size: 0.95rem;
	}

	.dosen-table tr:last-child td {
		border-bottom: none;
	}

	.dosen-table tr:hover {
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

	.mono {
		font-family: monospace;
		color: var(--text-muted);
	}

	.mono {
		font-family: monospace;
		color: var(--text-muted);
	}

	.action-buttons {
		display: flex;
		align-items: center;
		gap: 8px;
	}

	.btn-detail {
		background: rgba(79, 70, 229, 0.1);
		color: var(--primary-color);
		border: 1px solid rgba(79, 70, 229, 0.2);
		cursor: pointer;
		transition: all 0.2s;
		font-size: 0.85rem;
		padding: 6px 12px;
		border-radius: var(--radius-sm);
		font-weight: 500;
	}

	.btn-detail:hover {
		background: rgba(79, 70, 229, 0.2);
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

	/* Modal Styles */
	.modal-backdrop {
		position: fixed;
		top: 0;
		left: 0;
		width: 100vw;
		height: 100vh;
		background: rgba(15, 23, 42, 0.6);
		backdrop-filter: blur(4px);
		-webkit-backdrop-filter: blur(4px);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 100;
		padding: 24px;
	}

	.modal-content {
		width: 100%;
		max-width: 480px;
		padding: 0;
		display: flex;
		flex-direction: column;
		border-radius: var(--radius-lg);
		overflow: hidden;
	}

	.modal-header {
		padding: 20px 24px;
		border-bottom: 1px solid var(--surface-border);
		display: flex;
		align-items: center;
		justify-content: space-between;
	}

	.modal-header h3 {
		font-size: 1.25rem;
		font-weight: 600;
		color: var(--text-main);
	}

	.btn-close {
		background: transparent;
		border: none;
		font-size: 1.25rem;
		color: var(--text-muted);
		cursor: pointer;
		transition: color 0.2s;
	}

	.btn-close:hover {
		color: var(--error-color);
	}

	.modal-body {
		padding: 24px;
		display: flex;
		flex-direction: column;
		gap: 16px;
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

	.modal-footer {
		padding: 16px 24px;
		border-top: 1px solid var(--surface-border);
		display: flex;
		justify-content: flex-end;
		background: rgba(248, 250, 252, 0.5);
	}

	.btn-secondary {
		padding: 10px 20px;
		background: white;
		border: 1px solid var(--surface-border);
		color: var(--text-main);
		border-radius: var(--radius-sm);
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-secondary:hover {
		background: #f1f5f9;
	}
</style>
