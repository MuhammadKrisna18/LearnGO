<script lang="ts">
	import { onMount } from 'svelte';
	import { authState } from '$lib/stores/auth.svelte';
	import { authService } from '$lib/services/auth';
	import ProfileCard from '$lib/components/dashboard/ProfileCard.svelte';
	import DosenListCard from '$lib/components/dashboard/DosenListCard.svelte';

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

<div class="content">
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
		<ProfileCard profile={authState.profile} />

		{#if authState.profile.role === 'admin'}
			<DosenListCard />
		{:else if authState.profile.role === 'dosen'}
			<div class="dosen-empty-state glass-panel animate-fade-in" style="animation-delay: 0.2s;">
				<div class="empty-icon">🎓</div>
				<h3>Selamat Datang, Dosen {authState.profile.name}!</h3>
				<p>Belum ada jadwal atau kelas yang ditugaskan kepada Anda saat ini.</p>
				<button class="btn-primary empty-action-btn">Cek Notifikasi</button>
			</div>
		{/if}
	{/if}
</div>

<style>
	.content {
		max-width: 800px;
		display: flex;
		flex-direction: column;
		gap: 24px;
	}

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
		margin-bottom: 20px;
		font-size: 0.875rem;
		text-align: center;
	}

	.dosen-empty-state {
		padding: 48px 32px;
		text-align: center;
		border-radius: var(--radius-lg);
		background: linear-gradient(135deg, rgba(37, 99, 235, 0.05), rgba(56, 189, 248, 0.05));
		border: 1px solid rgba(37, 99, 235, 0.1);
	}

	.empty-icon {
		font-size: 48px;
		margin-bottom: 16px;
		filter: drop-shadow(0 4px 12px rgba(37, 99, 235, 0.2));
	}

	.dosen-empty-state h3 {
		font-size: 1.5rem;
		color: var(--text-main);
		margin-bottom: 12px;
	}

	.dosen-empty-state p {
		color: var(--text-muted);
		margin-bottom: 24px;
		max-width: 400px;
		margin-left: auto;
		margin-right: auto;
	}

	.empty-action-btn {
		background: var(--primary-color);
		color: white;
		padding: 10px 24px;
		border-radius: var(--radius-md);
		font-weight: 500;
		border: none;
		cursor: pointer;
		transition: all 0.2s;
	}

	.empty-action-btn:hover {
		background: var(--primary-hover);
		transform: translateY(-1px);
		box-shadow: 0 4px 12px rgba(37, 99, 235, 0.2);
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}
</style>
