<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { authState } from '$lib/stores/auth.svelte';
	import DosenRegisterCard from '$lib/components/dashboard/DosenRegisterCard.svelte';
	import DosenListCard from '$lib/components/dashboard/DosenListCard.svelte';

	onMount(() => {
		if (authState.role !== 'admin') {
			goto('/dashboard');
		}
	});
</script>

<svelte:head>
	<title>Manajemen Dosen - Modular Monolith</title>
</svelte:head>

<div class="page-header">
	<h1>Manajemen Dosen</h1>
	<p>Kelola akses dan akun dosen dalam sistem.</p>
</div>

<div class="page-grid two-cols">
	<div class="col-main">
		<DosenListCard />
	</div>
	<div class="col-side">
		<DosenRegisterCard />
	</div>
</div>

<style>
	/* Custom widths for specific layout, overriding standard 1fr 1fr if needed */
	@media (min-width: 1024px) {
		.page-grid.two-cols {
			grid-template-columns: 2fr 1fr;
		}
	}
</style>
