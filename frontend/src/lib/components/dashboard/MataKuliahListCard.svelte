<script lang="ts">
	import { onMount } from 'svelte';
	import { matakuliahService } from '$lib/services/matakuliah';
	import { programStudiService } from '$lib/services/programstudi';
	import { authState } from '$lib/stores/auth.svelte';
	import type { MataKuliah, ProgramStudi } from '$lib/types';
	import PromptCodeModal from '$lib/components/ui/PromptCodeModal.svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';
	import { toast } from '$lib/stores/toast.svelte';

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
				if (authState.role === 'admin') {
					mkList = resMk.data;
				} else {
					mkList = resMk.data.filter((mk: MataKuliah) => !mk.pengajuan?.some(p => p.status === 'approved'));
				}
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

	let isPromptOpen = $state(false);
	let promptTitle = $state("");
	let promptMessage = $state("");
	let expectedCode = $state("");
	let pendingAction = $state<"delete" | "lepas" | null>(null);
	let selectedMK = $state<MataKuliah | null>(null);

	function promptDelete(mk: MataKuliah) {
		selectedMK = mk;
		pendingAction = "delete";
		expectedCode = Math.floor(100000 + Math.random() * 900000).toString();
		promptTitle = "Hapus Mata Kuliah";
		promptMessage = `Untuk menghapus mata kuliah ${selectedMK?.name || ''}, masukkan kode berikut:`;
		isPromptOpen = true;
	}

	function promptLepas(mk: MataKuliah) {
		selectedMK = mk;
		pendingAction = "lepas";
		expectedCode = Math.floor(100000 + Math.random() * 900000).toString();
		promptTitle = "Lepaskan Mata Kuliah";
		promptMessage = `Untuk melepas mata kuliah ${selectedMK?.name || ''}, masukkan kode berikut:`;
		isPromptOpen = true;
	}

	async function handlePromptConfirm(code: string) {
		if (!selectedMK) return;

		if (code !== expectedCode) {
			toast.error('Kode tidak cocok. Aksi dibatalkan.');
			selectedMK = null;
			isPromptOpen = false;
			return;
		}

		if (pendingAction === "delete") {
			try {
				const res = await matakuliahService.delete(selectedMK.id);
				if (res.success) {
					toast.success('Mata kuliah berhasil dihapus');
					mkList = mkList.filter(m => m.id !== selectedMK!.id);
				} else {
					toast.error(res.message || 'Gagal menghapus mata kuliah');
				}
			} catch (err: any) {
				toast.error(err.message || 'Gagal terhubung ke server');
			} finally {
				selectedMK = null;
				isPromptOpen = false;
			}
		} else if (pendingAction === "lepas") {
			try {
				const res = await matakuliahService.lepasMataKuliah(selectedMK.id);
				if (res.success) {
					toast.success('Mata kuliah berhasil dilepas');
					const resMk = await matakuliahService.getList();
					if (resMk.success && resMk.data) {
						if (authState.role === 'admin') {
							mkList = resMk.data;
						} else {
							mkList = resMk.data.filter((mk: MataKuliah) => !mk.pengajuan?.some(p => p.status === 'approved'));
						}
					}
				} else {
					toast.error(res.message || 'Gagal melepas mata kuliah');
				}
			} catch (err: any) {
				toast.error(err.message || 'Gagal terhubung ke server');
			} finally {
				selectedMK = null;
				isPromptOpen = false;
			}
		}
	}
</script>

<div class="mk-list-container animate-fade-in" style="animation-delay: 0.4s;">
	<div class="section-header">
		<h3>Daftar Mata Kuliah per Program Studi</h3>
		<p>Daftar seluruh mata kuliah yang tersedia di dalam sistem dikelompokkan berdasarkan program studi.</p>
	</div>

	{#if loading}
		<Card class="state-container">Memuat data...</Card>
	{:else if error}
		<Card class="state-container state-error">{error}</Card>
	{:else if mkList.length === 0}
		<Card class="state-container">Belum ada mata kuliah yang didaftarkan.</Card>
	{:else}
		<div class="prodi-groups">
			{#each prodiList as prodi}
				{@const prodiMks = mkList.filter(mk => mk.program_studi_id === prodi.id)}
				<Card class="prodi-section" style="padding: 24px;">
					<div class="prodi-header">
						<h4 class="prodi-title">{prodi.name} ({prodi.code})</h4>
						<span class="prodi-count">{prodiMks.length} Mata Kuliah</span>
					</div>
					
					{#if prodiMks.length === 0}
						<div class="state-container prodi-empty">Belum ada mata kuliah untuk program studi ini.</div>
					{:else}
						<div class="data-table-container">
							<table class="data-table">
								<thead>
									<tr>
										<th>Nama Mata Kuliah</th>
										<th>SKS</th>
										<th>Dosen Pengampu</th>
										<th>Tgl Ditambahkan</th>
										{#if authState.role === 'admin'}
											<th>Aksi</th>
										{/if}
									</tr>
								</thead>
								<tbody>
									{#each prodiMks as mk}
										<tr>
											<td><span class="font-medium">{mk.name}</span></td>
											<td><Badge type="info">{mk.sks} SKS</Badge></td>
											<td>
												{#if mk.pengajuan && mk.pengajuan.some(p => p.status === 'approved')}
													{@const approvedP = mk.pengajuan.find(p => p.status === 'approved')}
													<div class="dosen-info">
														<div class="dosen-name-wrapper">
															<span class="dosen-name">{approvedP?.dosen?.name || 'Unknown'}</span>
															{#if approvedP?.dosen?.program_studi_id === mk.program_studi_id}
																<Badge type="success">Sesuai Prodi</Badge>
															{:else if approvedP?.dosen}
																<Badge type="warning">Lintas Jurusan</Badge>
															{/if}
														</div>
														{#if approvedP?.dosen?.email}
															<span class="dosen-email">{approvedP.dosen.email}</span>
														{/if}
													</div>
												{:else}
													<span class="text-muted italic">Belum diambil</span>
												{/if}
											</td>
											<td>{new Date(mk.created_at).toLocaleDateString()}</td>
											{#if authState.role === 'admin'}
												<td>
													{#if mk.pengajuan && mk.pengajuan.some(p => p.status === 'approved' || p.status === 'pending')}
														<button class="btn-lepas" aria-label="Lepas Mata Kuliah" onclick={() => promptLepas(mk)}>
															Lepas
														</button>
													{:else}
														<button class="btn-delete" aria-label="Hapus Mata Kuliah" onclick={() => promptDelete(mk)}>
															Hapus
														</button>
													{/if}
												</td>
											{/if}
										</tr>
									{/each}
								</tbody>
							</table>
						</div>
					{/if}
				</Card>
			{/each}
		</div>
	{/if}
</div>

<PromptCodeModal
	bind:isOpen={isPromptOpen}
	title={promptTitle}
	message={promptMessage}
	{expectedCode}
	onConfirm={handlePromptConfirm}
	onCancel={() => { isPromptOpen = false; selectedMK = null; }}
/>

<style>
	.mk-list-container {
		margin-top: 24px;
	}

	.section-header {
		margin-bottom: 24px;
	}

	.section-header h3 {
		font-size: 1.25rem;
		font-weight: 600;
		color: var(--text-main);
		margin-bottom: 4px;
	}

	.section-header p {
		font-size: 0.9rem;
		color: var(--text-muted);
	}

	.prodi-section {
		border-radius: var(--radius-lg);
	}

	.prodi-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 16px;
		padding-bottom: 12px;
		border-bottom: 2px solid rgba(79, 70, 229, 0.1);
	}

	.prodi-title {
		font-size: 1.1rem;
		font-weight: 600;
		color: var(--primary-color);
		margin: 0;
	}

	.prodi-count {
		font-size: 0.85rem;
		color: var(--text-muted);
		background: rgba(0,0,0,0.05);
		padding: 4px 10px;
		border-radius: 9999px;
		font-weight: 500;
	}
	
	.dosen-info {
		display: flex;
		flex-direction: column;
		gap: 2px;
	}
	
	.dosen-name-wrapper {
		display: flex;
		align-items: center;
		gap: 8px;
		flex-wrap: wrap;
	}
	
	.dosen-name {
		font-weight: 500;
		color: var(--text-main);
	}
	
	.dosen-email {
		font-size: 0.75rem;
		color: var(--text-muted);
	}
	
	.italic {
		font-style: italic;
	}
	
	.prodi-groups {
		display: flex;
		flex-direction: column;
		gap: 24px;
	}
	
	.prodi-empty {
		padding: 32px;
		background: transparent;
		border: 1px dashed var(--surface-border);
		color: var(--text-muted);
	}

	.btn-lepas {
		padding: 6px 12px;
		background: rgba(245, 158, 11, 0.1);
		color: #d97706;
		border-radius: var(--radius-sm);
		font-weight: 500;
		font-size: 0.85rem;
		transition: all 0.2s;
	}

	.btn-lepas:hover {
		background: rgba(245, 158, 11, 0.2);
	}
</style>
