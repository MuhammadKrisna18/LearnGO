<script lang="ts">
	import { onMount } from 'svelte';
	import { authState } from '$lib/stores/auth.svelte';
	import { authService } from '$lib/services/auth';
	import EmailRequestListCard from '$lib/components/dashboard/EmailRequestListCard.svelte';
	import DosenMataKuliahRequestCard from '$lib/components/dashboard/DosenMataKuliahRequestCard.svelte';
	import AdminPengajuanCard from '$lib/components/dashboard/AdminPengajuanCard.svelte';
	import AdminPengajuanKelasCard from '$lib/components/dashboard/AdminPengajuanKelasCard.svelte';
	import AdminDosenKelasListCard from '$lib/components/dashboard/AdminDosenKelasListCard.svelte';
	import DosenPengajuanKelasCard from '$lib/components/dashboard/DosenPengajuanKelasCard.svelte';
	import DosenKelasSayaCard from '$lib/components/dashboard/DosenKelasSayaCard.svelte';

	let loading = $state(true);
	let error = $state('');

	onMount(async () => {
		try {
			if (!authState.profile) {
				await authService.getProfile();
			}
		} catch (err: any) {
			error = err.message || 'Failed to connect to the server.';
		} finally {
			loading = false;
		}
	});
</script>

<svelte:head>
	<title>Dashboard - Modular Monolith</title>
</svelte:head>

<div class="page-header">
	<h1>Beranda</h1>
	<p>Ringkasan informasi akun dan notifikasi Anda.</p>
</div>

<div class="page-grid">
	{#if loading}
		<div class="loading-state glass-panel animate-fade-in">
			<div class="spinner"></div>
			<p>Loading your profile...</p>
		</div>
	{:else if error}
		<div class="error-message animate-fade-in">
			{error}
		</div>
	{:else if authState.profile}
		{#if authState.profile.role === 'admin'}
			<div class="page-grid">
				<EmailRequestListCard />
				<AdminPengajuanCard />
				<AdminPengajuanKelasCard />
				<AdminDosenKelasListCard />
			</div>
		{:else if authState.profile.role === 'dosen'}
			<DosenKelasSayaCard />
			<DosenMataKuliahRequestCard />
			<DosenPengajuanKelasCard />
		{/if}
	{/if}
</div>

<style>
	.loading-state {
		padding: 48px;
		text-align: center;
		border-radius: var(--radius-lg);
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 16px;
		color: var(--text-muted);
	}

	.spinner {
		width: 32px;
		height: 32px;
		border: 3px solid rgba(0, 0, 0, 0.1);
		border-top-color: var(--primary-color);
		border-radius: 50%;
		animation: spin 1s linear infinite;
	}

	.error-message {
		background: rgba(239, 68, 68, 0.1);
		color: var(--error-color);
		padding: 12px;
		border-radius: var(--radius-sm);
		border: 1px solid rgba(239, 68, 68, 0.2);
		font-size: 0.875rem;
		text-align: center;
	}

	.dosen-empty-state {
		padding: 48px 32px;
		text-align: center;
		border-radius: var(--radius-lg);
		background: linear-gradient(135deg, rgba(79, 70, 229, 0.03), rgba(16, 185, 129, 0.03));
	}

	.empty-icon {
		margin-bottom: 20px;
		display: flex;
		justify-content: center;
		filter: drop-shadow(0 4px 12px rgba(79, 70, 229, 0.2));
	}

	.dosen-empty-state h3 {
		font-size: 1.5rem;
		color: var(--text-main);
		margin-bottom: 12px;
		font-weight: 700;
		letter-spacing: -0.025em;
	}

	.dosen-empty-state p {
		color: var(--text-muted);
		margin-bottom: 24px;
		max-width: 400px;
		margin-left: auto;
		margin-right: auto;
	}

	.empty-action-btn {
		width: auto;
		padding: 10px 24px;
		border-radius: var(--radius-full);
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}
</style>
