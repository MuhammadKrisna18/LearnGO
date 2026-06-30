<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { authState } from '$lib/stores/auth.svelte';

	let { children } = $props();

	onMount(() => {
		if (!authState.token || authState.role !== 'mahasiswa') {
			goto('/login');
		}
	});

	function handleLogout(e: Event) {
		e.preventDefault();
		authState.logout();
	}
</script>

<div class="layout-wrapper">
	<nav class="top-nav glass-panel">
		<div class="nav-brand">
			<div class="brand-logo">
				<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m12 3-1.912 5.813a2 2 0 0 1-1.275 1.275L3 12l5.813 1.912a2 2 0 0 1 1.275 1.275L12 21l1.912-5.813a2 2 0 0 1 1.275-1.275L21 12l-5.813-1.912a2 2 0 0 1-1.275-1.275L12 3Z"/></svg>
			</div>
			<div class="brand-text">
				<h1>SIAKAD</h1>
				<p>Mahasiswa</p>
			</div>
		</div>
		<div class="nav-actions">
			<span class="user-info">
				<span class="user-name">{authState.user?.name || 'Mahasiswa'}</span>
				<span class="badge badge-primary">{authState.user?.nrp || 'NRP'}</span>
			</span>
			<button class="btn btn-outline btn-sm logout-btn" onclick={handleLogout}>
				<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
				Keluar
			</button>
		</div>
	</nav>

	<main class="main-content">
		{@render children()}
	</main>
</div>

<style>
	.layout-wrapper {
		min-height: 100vh;
		background-color: var(--bg-body);
		display: flex;
		flex-direction: column;
	}

	.top-nav {
		position: sticky;
		top: 0;
		z-index: 40;
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 1rem 2rem;
		border-bottom: 1px solid var(--border);
		background: rgba(255, 255, 255, 0.85);
		backdrop-filter: blur(12px);
	}

	.nav-brand {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}

	.brand-logo {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 40px;
		height: 40px;
		background: linear-gradient(135deg, var(--primary) 0%, #312e81 100%);
		border-radius: 12px;
		color: white;
		box-shadow: 0 4px 12px rgba(79, 70, 229, 0.25);
	}

	.brand-text h1 {
		font-family: 'Outfit', sans-serif;
		font-size: 1.25rem;
		font-weight: 700;
		color: var(--text-main);
		margin: 0;
		line-height: 1.1;
		letter-spacing: -0.02em;
	}

	.brand-text p {
		font-size: 0.75rem;
		color: var(--primary);
		font-weight: 600;
		margin: 0;
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}

	.nav-actions {
		display: flex;
		align-items: center;
		gap: 1.5rem;
	}

	.user-info {
		display: flex;
		align-items: center;
		gap: 0.75rem;
	}

	.user-name {
		font-weight: 500;
		color: var(--text-main);
	}

	.logout-btn {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.5rem 1rem;
	}

	.main-content {
		flex: 1;
		padding: 2rem;
		max-width: 1200px;
		margin: 0 auto;
		width: 100%;
		box-sizing: border-box;
	}
</style>
