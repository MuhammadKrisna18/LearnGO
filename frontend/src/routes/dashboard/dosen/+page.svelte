<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { authState } from '$lib/stores/auth.svelte';
	import DosenRegisterCard from '$lib/components/dashboard/DosenRegisterCard.svelte';

	onMount(() => {
		if (authState.role !== 'admin') {
			goto('/dashboard');
		}
	});
</script>

<svelte:head>
	<title>Manajemen Dosen - Modular Monolith</title>
</svelte:head>

<div class="content">
	<div class="header-section">
		<h2>Manajemen Dosen</h2>
		<p>Kelola akses dan akun dosen dalam sistem.</p>
	</div>

	<DosenRegisterCard />
</div>

<style>
	.content {
		max-width: 800px;
		display: flex;
		flex-direction: column;
		gap: 24px;
	}

	.header-section {
		margin-bottom: 8px;
	}

	.header-section h2 {
		font-size: 1.75rem;
		font-weight: 700;
		color: var(--text-main);
	}

	.header-section p {
		color: var(--text-muted);
		margin-top: 4px;
	}

	.mono {
		font-family: monospace;
		color: var(--text-muted);
	}
</style>
