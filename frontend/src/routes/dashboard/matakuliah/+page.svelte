<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { authState } from '$lib/stores/auth.svelte';
	import MataKuliahRegisterCard from '$lib/components/dashboard/MataKuliahRegisterCard.svelte';
	import MataKuliahListCard from '$lib/components/dashboard/MataKuliahListCard.svelte';
	onMount(() => {
		if (!authState.token) {
			goto('/login');
		} else if (authState.role !== 'admin' && authState.role !== 'dosen') {
			goto('/dashboard');
		}
	});
</script>

<svelte:head>
	<title>Mata Kuliah - Modular Monolith</title>
</svelte:head>

<div class="page-header">
	<h1>Manajemen Mata Kuliah</h1>
	<p>Tambahkan dan kelola mata kuliah yang tersedia di kampus.</p>
</div>

<div class="page-grid {authState.role === 'admin' ? 'two-cols' : ''}">
	<div class="col-main">
		<MataKuliahListCard />
	</div>
	{#if authState.role === 'admin'}
		<div class="col-side">
			<MataKuliahRegisterCard />
		</div>
	{/if}
</div>


<style>
	/* Custom widths for specific layout, overriding standard 1fr 1fr if needed */
	@media (min-width: 1024px) {
		.page-grid.two-cols {
			grid-template-columns: 2fr 1fr;
		}
	}
</style>
