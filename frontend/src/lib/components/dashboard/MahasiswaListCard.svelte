<script lang="ts">
	import { onMount } from 'svelte';
	import { authService } from '$lib/services/auth';
	import { programStudiService } from '$lib/services/programstudi';
	import type { UserProfile, ProgramStudi } from '$lib/types';
	import Card from '$lib/components/ui/Card.svelte';
	import PromptCodeModal from '$lib/components/ui/PromptCodeModal.svelte';
	import { toast } from '$lib/stores/toast.svelte';

	let mahasiswaList: UserProfile[] = $state([]);
	let prodiList: ProgramStudi[] = $state([]);
	let loading = $state(true);
	let error = $state('');

	onMount(async () => {
		await fetchData();
		
		// Listen for the custom event to refresh data
		window.addEventListener('mahasiswa-added', fetchData as EventListener);
		
		return () => {
			window.removeEventListener('mahasiswa-added', fetchData as EventListener);
		};
	});

	async function fetchData() {
		loading = true;
		try {
			const [mhsRes, prodiRes] = await Promise.all([
				authService.getMahasiswaList(),
				programStudiService.getList()
			]);

			if (mhsRes.success && mhsRes.data) {
				mahasiswaList = mhsRes.data;
			} else {
				error = mhsRes.message || 'Gagal mengambil daftar mahasiswa';
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
	}

	function getMahasiswaByProdi(prodiId: string) {
		return mahasiswaList.filter(m => m.program_studi_id === prodiId);
	}

	function getMahasiswaWithoutProdi() {
		return mahasiswaList.filter(m => !m.program_studi_id);
	}

</script>

<Card title="Daftar Mahasiswa Aktif">
	<div class="header-actions" slot="header-actions">
		<button class="btn btn-secondary btn-sm" onclick={fetchData} disabled={loading}>
			{#if loading}
				<span class="spinner"></span> Memuat...
			{:else}
				<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
					<path d="M21 2v6h-6"/><path d="M3 12a9 9 0 0 1 15-6.7L21 8"/><path d="M3 22v-6h6"/><path d="M21 12a9 9 0 0 1-15 6.7L3 16"/>
				</svg> Refresh
			{/if}
		</button>
	</div>

	{#if loading && mahasiswaList.length === 0}
		<div class="loading-state">
			<div class="spinner large"></div>
			<p>Memuat data mahasiswa...</p>
		</div>
	{:else if error}
		<div class="error-state">
			<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
				<circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>
			</svg>
			<p>{error}</p>
			<button class="btn btn-secondary" onclick={fetchData}>Coba Lagi</button>
		</div>
	{:else if mahasiswaList.length === 0}
		<div class="empty-state">
			<div class="empty-icon">
				<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
					<path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/>
					<circle cx="9" cy="7" r="4"/>
					<line x1="19" y1="8" x2="19" y2="14"/>
					<line x1="16" y1="11" x2="22" y2="11"/>
				</svg>
			</div>
			<h3>Belum ada mahasiswa</h3>
			<p>Data mahasiswa yang terdaftar akan muncul di sini. Silakan daftarkan mahasiswa pertama menggunakan form di samping.</p>
		</div>
	{:else}
		<div class="prodi-groups">
			{#each prodiList as prodi}
				{@const prodiMhs = getMahasiswaByProdi(prodi.id)}
				{#if prodiMhs.length > 0}
					<div class="prodi-section">
						<div class="prodi-header">
							<div class="prodi-info">
								<h3 class="prodi-title">{prodi.name}</h3>
								<span class="badge badge-primary">{prodiMhs.length} Mahasiswa</span>
							</div>
						</div>
						
						<div class="list-container">
							{#each prodiMhs as mhs, index}
								<div class="list-item">
									<div class="item-number">{index + 1}</div>
									<div class="item-main">
										<div class="item-header">
											<h4 class="item-title">{mhs.name}</h4>
											<span class="badge badge-outline">NRP: {mhs.nrp || '-'}</span>
										</div>
										<div class="item-meta">
											<span class="meta-item">
												<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
													<rect width="20" height="16" x="2" y="4" rx="2"/><path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7"/>
												</svg>
												{mhs.email}
											</span>
										</div>
									</div>
								</div>
							{/each}
						</div>
					</div>
				{/if}
			{/each}

			{#if getMahasiswaWithoutProdi().length > 0}
				{@const withoutProdi = getMahasiswaWithoutProdi()}
				<div class="prodi-section">
					<div class="prodi-header">
						<div class="prodi-info">
							<h3 class="prodi-title text-warning">Belum Ada Program Studi</h3>
							<span class="badge badge-warning">{withoutProdi.length} Mahasiswa</span>
						</div>
					</div>
					
					<div class="list-container">
						{#each withoutProdi as mhs, index}
							<div class="list-item warning-border">
								<div class="item-number">{index + 1}</div>
								<div class="item-main">
									<div class="item-header">
										<h4 class="item-title">{mhs.name}</h4>
										<span class="badge badge-outline">NRP: {mhs.nrp || '-'}</span>
									</div>
									<div class="item-meta">
										<span class="meta-item">
											<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
												<rect width="20" height="16" x="2" y="4" rx="2"/><path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7"/>
											</svg>
											{mhs.email}
										</span>
									</div>
								</div>
							</div>
						{/each}
					</div>
				</div>
			{/if}
		</div>
	{/if}
</Card>

<style>
	.header-actions {
		display: flex;
		gap: 0.5rem;
	}

	.btn-sm {
		padding: 0.35rem 0.75rem;
		font-size: 0.85rem;
		display: flex;
		align-items: center;
		gap: 0.35rem;
	}

	.prodi-groups {
		display: flex;
		flex-direction: column;
		gap: 2rem;
	}

	.prodi-section {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.prodi-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding-bottom: 0.75rem;
		border-bottom: 2px solid var(--border);
	}

	.prodi-info {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}

	.prodi-title {
		margin: 0;
		font-size: 1.1rem;
		color: var(--text-main);
	}

	.text-warning {
		color: #d97706;
	}

	.badge-warning {
		background-color: #fef3c7;
		color: #d97706;
	}

	.list-container {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.list-item {
		display: flex;
		align-items: center;
		gap: 1.25rem;
		padding: 1.25rem;
		background: white;
		border: 1px solid var(--border);
		border-radius: 12px;
		transition: all 0.2s ease;
	}

	.list-item:hover {
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
		border-color: var(--primary-light);
	}

	.warning-border {
		border-left: 4px solid #f59e0b;
	}

	.item-number {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 36px;
		height: 36px;
		background-color: var(--surface);
		color: var(--text-muted);
		border-radius: 50%;
		font-weight: 600;
		font-size: 0.9rem;
		flex-shrink: 0;
	}

	.item-main {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: 0.4rem;
	}

	.item-header {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		flex-wrap: wrap;
	}

	.item-title {
		margin: 0;
		font-size: 1.05rem;
		font-weight: 600;
		color: var(--text-main);
	}

	.item-meta {
		display: flex;
		align-items: center;
		gap: 1.25rem;
		flex-wrap: wrap;
	}

	.meta-item {
		display: flex;
		align-items: center;
		gap: 0.4rem;
		font-size: 0.85rem;
		color: var(--text-muted);
	}

	.loading-state, .error-state, .empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 4rem 2rem;
		text-align: center;
		background: var(--surface);
		border-radius: 12px;
		border: 1px dashed var(--border);
	}

	.empty-icon {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 80px;
		height: 80px;
		border-radius: 50%;
		background: white;
		color: var(--text-muted);
		margin-bottom: 1.5rem;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
	}

	.empty-state h3 {
		margin: 0 0 0.5rem 0;
		font-size: 1.25rem;
		color: var(--text-main);
	}

	.empty-state p {
		margin: 0;
		color: var(--text-muted);
		max-width: 400px;
	}
</style>
