<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { authState } from '$lib/stores/auth.svelte';
	import MataKuliahRegisterCard from '$lib/components/dashboard/MataKuliahRegisterCard.svelte';

	onMount(() => {
		if (authState.role !== 'admin') {
			goto('/dashboard');
		}
	});
</script>

<svelte:head>
	<title>Mata Kuliah - Modular Monolith</title>
</svelte:head>

<div class="content">
	<div class="header-section">
		<h2>Manajemen Mata Kuliah</h2>
		<p>Tambahkan dan kelola mata kuliah yang tersedia di kampus.</p>
	</div>

	<MataKuliahRegisterCard />
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
</style>
