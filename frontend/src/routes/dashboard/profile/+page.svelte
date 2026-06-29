<script lang="ts">
	import { onMount } from 'svelte';
	import { authState } from '$lib/stores/auth.svelte';
	import { authService } from '$lib/services/auth';
	import ProfileCard from '$lib/components/dashboard/ProfileCard.svelte';
	import DosenProfileCard from '$lib/components/dashboard/DosenProfileCard.svelte';
	import Skeleton from '$lib/components/ui/Skeleton.svelte';

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
	<title>Profil Anda - SIAKAD Pro</title>
</svelte:head>

<div class="page-header">
	<h1>Profil Anda</h1>
	<p>Kelola informasi akun Anda.</p>
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
			<ProfileCard profile={authState.profile} />
		{:else if authState.profile.role === 'dosen'}
			<DosenProfileCard profile={authState.profile} />
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

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}
</style>
