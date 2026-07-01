<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { authState } from '$lib/stores/auth.svelte';
	import MahasiswaRegisterCard from '$lib/components/dashboard/MahasiswaRegisterCard.svelte';
	import MahasiswaListCard from '$lib/components/dashboard/MahasiswaListCard.svelte';

	onMount(() => {
		if (authState.role !== 'admin') {
			goto('/dashboard');
		}
	});
</script>

<svelte:head>
	<title>Manajemen Mahasiswa - SIAKAD Pro</title>
</svelte:head>

<div class="page-header">
	<h1>Manajemen Mahasiswa</h1>
	<p>Kelola akses dan akun mahasiswa dalam sistem.</p>
</div>

<div class="page-grid two-cols">
	<div class="col-main">
		<MahasiswaListCard />
	</div>
	<div class="col-side">
		<MahasiswaRegisterCard />
	</div>
</div>

<style>
	
	@media (min-width: 1024px) {
		.page-grid.two-cols {
			grid-template-columns: 2fr 1fr;
		}
	}
</style>
