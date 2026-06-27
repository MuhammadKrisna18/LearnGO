<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { authState } from '$lib/stores/auth.svelte';
	import KelasListCard from '$lib/components/dashboard/KelasListCard.svelte';
	import KelasRegisterCard from '$lib/components/dashboard/KelasRegisterCard.svelte';

	onMount(() => {
		if (!authState.token) {
			goto('/login');
		} else if (authState.role !== 'admin' && authState.role !== 'dosen') {
			goto('/dashboard');
		}
	});
</script>

<svelte:head>
	<title>Manajemen Kelas - Modular Monolith</title>
</svelte:head>

<div class="page-header">
	<h1>Manajemen Kelas</h1>
	<p>Daftar dan kelola kelas yang tersedia.</p>
</div>

<div class="page-grid {authState.role === 'admin' ? 'two-cols' : ''}">
	<div class="col-main">
		<KelasListCard />
	</div>
	{#if authState.role === 'admin'}
		<div class="col-side">
			<KelasRegisterCard />
		</div>
	{/if}
</div>

<style>
	.page-header {
		margin-bottom: 24px;
	}

	.page-header h1 {
		font-size: 1.8rem;
		font-weight: 700;
		color: var(--text-main);
		margin-bottom: 8px;
	}

	.page-header p {
		color: var(--text-muted);
	}

	.page-grid {
		display: grid;
		gap: 24px;
		grid-template-columns: 1fr;
		align-items: start;
	}

	.two-cols {
		grid-template-columns: 1fr 350px;
	}

	@media (max-width: 1024px) {
		.two-cols {
			grid-template-columns: 1fr;
		}
		
		.col-side {
			grid-row: 1;
		}
	}
</style>
