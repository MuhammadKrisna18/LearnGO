<script lang="ts">
	import { onMount } from 'svelte';
	import { authState } from '$lib/stores/auth.svelte';
	import { authService } from '$lib/services/auth';
	import { matakuliahService } from '$lib/services/matakuliah';
	import { kelasService } from '$lib/services/kelas';
	import { dosenService } from '$lib/services/dosen';
	import EmailRequestListCard from '$lib/components/dashboard/EmailRequestListCard.svelte';
	import DosenMataKuliahRequestCard from '$lib/components/dashboard/DosenMataKuliahRequestCard.svelte';
	import AdminPengajuanCard from '$lib/components/dashboard/AdminPengajuanCard.svelte';
	import AdminPengajuanKelasCard from '$lib/components/dashboard/AdminPengajuanKelasCard.svelte';
	import AdminDosenKelasListCard from '$lib/components/dashboard/AdminDosenKelasListCard.svelte';
	import DosenPengajuanKelasCard from '$lib/components/dashboard/DosenPengajuanKelasCard.svelte';
	import DosenKelasSayaCard from '$lib/components/dashboard/DosenKelasSayaCard.svelte';
	import StatsRow from '$lib/components/dashboard/StatsRow.svelte';
	import Skeleton from '$lib/components/ui/Skeleton.svelte';

	let loading = $state(true);
	let error = $state('');
	let adminStats = $state<{ label: string; value: number; icon: string; color: string }[]>([]);
	let dosenStats = $state<{ label: string; value: number; icon: string; color: string }[]>([]);

	onMount(async () => {
		try {
			if (!authState.profile) {
				await authService.getProfile();
			}

			if (authState.profile?.role === 'admin') {
				const [resDosen, resMk, resKelas] = await Promise.all([
					dosenService.getList(),
					matakuliahService.getList(),
					kelasService.getList()
				]);
				const allDosen = resDosen.data || [];
				const allMk = resMk.data || [];
				const allKelas = resKelas.data || [];

				const pendingMk = allMk.filter((mk: any) =>
					mk.pengajuan?.some((p: any) => p.status === 'pending')
				).length;

				adminStats = [
					{ label: 'Total Dosen', value: allDosen.length, icon: '<svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M22 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>', color: '#6366f1' },
					{ label: 'Total Mata Kuliah', value: allMk.length, icon: '<svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1 0-5H20"/></svg>', color: '#10b981' },
					{ label: 'Total Kelas', value: allKelas.length, icon: '<svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"/><line x1="3" y1="9" x2="21" y2="9"/><line x1="9" y1="21" x2="9" y2="9"/></svg>', color: '#f59e0b' },
					{ label: 'Pengajuan Pending', value: pendingMk, icon: '<svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>', color: '#ef4444' }
				];
			} else if (authState.profile?.role === 'dosen') {
				const [resMk, resKelas] = await Promise.all([
					matakuliahService.getMyRequests(),
					kelasService.getMyPengajuan()
				]);
				const myMkRequests = resMk.data || [];
				const myKelasRequests = resKelas.data || [];

				const approvedMk = myMkRequests.filter((r: any) => r.status === 'approved').length;
				const approvedKelas = myKelasRequests.filter((r: any) => r.status === 'approved').length;

				dosenStats = [
					{ label: 'MK Diampu', value: approvedMk, icon: '<svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1 0-5H20"/></svg>', color: '#10b981' },
					{ label: 'Kelas Aktif', value: approvedKelas, icon: '<svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"/><line x1="3" y1="9" x2="21" y2="9"/><line x1="9" y1="21" x2="9" y2="9"/></svg>', color: '#6366f1' },
					{ label: 'Pengajuan MK', value: myMkRequests.length, icon: '<svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>', color: '#f59e0b' },
					{ label: 'Pengajuan Kelas', value: myKelasRequests.length, icon: '<svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>', color: '#ef4444' }
				];
			}
		} catch (err: any) {
			error = err.message || 'Failed to connect to the server.';
		} finally {
			loading = false;
		}
	});
</script>

<svelte:head>
	<title>Dashboard - SIAKAD Pro</title>
</svelte:head>

<div class="page-header">
	<h1>Beranda</h1>
	<p>Ringkasan informasi akun dan notifikasi Anda.</p>
</div>

<div class="page-grid">
	{#if loading}
		<Skeleton />
	{:else if error}
		<div class="error-message animate-fade-in">
			{error}
		</div>
	{:else if authState.profile}
		{#if authState.profile.role === 'admin'}
			<StatsRow stats={adminStats} />
			<div class="page-grid">
				<EmailRequestListCard />
				<AdminPengajuanCard />
				<AdminPengajuanKelasCard />
				<AdminDosenKelasListCard />
			</div>
		{:else if authState.profile.role === 'dosen'}
			<StatsRow stats={dosenStats} />
			<DosenKelasSayaCard />
			<DosenMataKuliahRequestCard />
			<DosenPengajuanKelasCard />
		{/if}
	{/if}
</div>

<style>
	.error-message {
		background: rgba(239, 68, 68, 0.1);
		color: var(--error-color);
		padding: 12px;
		border-radius: var(--radius-sm);
		border: 1px solid rgba(239, 68, 68, 0.2);
		font-size: 0.875rem;
		text-align: center;
	}
</style>
