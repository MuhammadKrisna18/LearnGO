<script>
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	onMount(() => {
		goto('/login');
	});
</script>

<svelte:head>
	<title>Redirecting...</title>
</svelte:head>

<div style="display: flex; justify-content: center; align-items: center; height: 100vh; color: white;">
	<p>Redirecting to login...</p>
</div>
