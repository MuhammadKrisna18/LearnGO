<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { authState } from '$lib/stores/auth.svelte';

	let { children } = $props();

	onMount(() => {
		if (!authState.token) {
			goto('/login');
		}
	});

	function handleLogout(e: Event) {
		e.preventDefault();
		authState.logout();
	}
</script>

<div class="layout-container">
	<aside class="sidebar glass-panel">
		<div class="sidebar-header">
			<div class="logo">
				<span class="logo-icon">✦</span>
			</div>
			<div class="brand-text">
				<h2>Modular</h2>
				<p>Monolith</p>
			</div>
		</div>

		<div class="sidebar-menu">
			<a href="/dashboard" class="menu-item {$page.url.pathname === '/dashboard' ? 'active' : ''}">
				<span class="menu-icon">🏠</span>
				<span class="menu-text">Dashboard</span>
			</a>
			{#if authState.role === 'admin'}
				<a href="/dashboard/dosen" class="menu-item {$page.url.pathname.includes('/dashboard/dosen') ? 'active' : ''}">
					<span class="menu-icon">👥</span>
					<span class="menu-text">Manajemen Dosen</span>
				</a>
			{/if}
		</div>

		<div class="sidebar-footer">
			<div class="user-role-badge">
				{authState.role?.toUpperCase() || 'USER'}
			</div>
			<button class="btn-logout" onclick={handleLogout}>
				<span class="logout-icon">🚪</span>
				Logout
			</button>
		</div>
	</aside>

	<main class="main-content">
		{@render children()}
	</main>
</div>

<style>
	.layout-container {
		display: flex;
		min-height: 100vh;
		background-color: transparent;
	}

	.sidebar {
		width: 260px;
		display: flex;
		flex-direction: column;
		border-radius: 0;
		border-right: 1px solid var(--surface-border);
		border-top: none;
		border-bottom: none;
		border-left: none;
		position: sticky;
		top: 0;
		height: 100vh;
		z-index: 10;
		background: rgba(255, 255, 255, 0.85);
	}

	.sidebar-header {
		padding: 32px 24px;
		display: flex;
		align-items: center;
		gap: 16px;
		border-bottom: 1px solid rgba(0, 0, 0, 0.05);
	}

	.logo {
		width: 40px;
		height: 40px;
		background: linear-gradient(135deg, var(--primary-color), #3b82f6);
		border-radius: 10px;
		display: flex;
		align-items: center;
		justify-content: center;
		box-shadow: 0 4px 12px rgba(37, 99, 235, 0.3);
		flex-shrink: 0;
	}

	.logo-icon {
		color: white;
		font-size: 20px;
	}

	.brand-text h2 {
		font-size: 1.1rem;
		font-weight: 700;
		color: var(--text-main);
		line-height: 1.2;
	}

	.brand-text p {
		font-size: 0.8rem;
		color: var(--text-muted);
		text-transform: uppercase;
		letter-spacing: 1px;
	}

	.sidebar-menu {
		flex: 1;
		padding: 24px 16px;
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.menu-item {
		display: flex;
		align-items: center;
		gap: 12px;
		padding: 12px 16px;
		border-radius: var(--radius-sm);
		color: var(--text-muted);
		text-decoration: none;
		font-weight: 500;
		transition: all 0.2s ease;
	}

	.menu-item:hover {
		background: rgba(37, 99, 235, 0.05);
		color: var(--primary-color);
	}

	.menu-item.active {
		background: rgba(37, 99, 235, 0.1);
		color: var(--primary-color);
	}

	.menu-icon {
		font-size: 1.2rem;
	}

	.sidebar-footer {
		padding: 24px;
		border-top: 1px solid rgba(0, 0, 0, 0.05);
		display: flex;
		flex-direction: column;
		gap: 16px;
	}

	.user-role-badge {
		text-align: center;
		font-size: 0.75rem;
		font-weight: 700;
		color: var(--primary-color);
		background: rgba(37, 99, 235, 0.1);
		padding: 6px 12px;
		border-radius: 20px;
		letter-spacing: 1px;
	}

	.btn-logout {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 8px;
		width: 100%;
		background: #ffffff;
		color: var(--error-color);
		border: 1px solid rgba(239, 68, 68, 0.2);
		padding: 10px 16px;
		border-radius: var(--radius-sm);
		font-family: inherit;
		font-size: 0.9rem;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-logout:hover {
		background: rgba(239, 68, 68, 0.1);
		border-color: rgba(239, 68, 68, 0.3);
	}

	.main-content {
		flex: 1;
		padding: 32px 40px;
		max-width: 1000px;
		margin: 0 auto;
		width: 100%;
		overflow-y: auto;
	}
</style>
