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

<div class="layout-wrapper">
	<!-- Sidebar -->
	<aside class="sidebar glass-panel">
		<div class="sidebar-brand">
			<div class="brand-logo">
				<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m12 3-1.912 5.813a2 2 0 0 1-1.275 1.275L3 12l5.813 1.912a2 2 0 0 1 1.275 1.275L12 21l1.912-5.813a2 2 0 0 1 1.275-1.275L21 12l-5.813-1.912a2 2 0 0 1-1.275-1.275L12 3Z"/></svg>
			</div>
			<div class="brand-text">
				<h1>Modular</h1>
				<p>Monolith</p>
			</div>
		</div>

		<nav class="sidebar-nav">
			<a href="/dashboard" class="nav-item {$page.url.pathname === '/dashboard' ? 'active' : ''}">
				<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="7" height="9" x="3" y="3" rx="1"/><rect width="7" height="5" x="14" y="3" rx="1"/><rect width="7" height="9" x="14" y="12" rx="1"/><rect width="7" height="5" x="3" y="16" rx="1"/></svg>
				<span>Dashboard</span>
			</a>
			
			{#if authState.role === 'admin'}
				<div class="nav-section-title">ADMINISTRASI</div>
				<a href="/dashboard/dosen" class="nav-item {$page.url.pathname.includes('/dashboard/dosen') ? 'active' : ''}">
					<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M22 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
					<span>Dosen</span>
				</a>
				<a href="/dashboard/matakuliah" class="nav-item {$page.url.pathname.includes('/dashboard/matakuliah') ? 'active' : ''}">
					<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1 0-5H20"/></svg>
					<span>Mata Kuliah</span>
				</a>
			{/if}
		</nav>

		<div class="sidebar-footer">
			<button class="btn-logout" onclick={handleLogout}>
				<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" x2="9" y1="12" y2="12"/></svg>
				<span>Log out</span>
			</button>
		</div>
	</aside>

	<!-- Main Content Area -->
	<div class="main-wrapper">
		<!-- Top Header -->
		<header class="top-header glass-panel">
			<div class="header-breadcrumb">
				<span class="text-muted">Home</span>
				<span class="separator">/</span>
				<span class="current-page">
					{#if $page.url.pathname === '/dashboard'}
						Dashboard
					{:else if $page.url.pathname.includes('/dashboard/dosen')}
						Manajemen Dosen
					{:else if $page.url.pathname.includes('/dashboard/matakuliah')}
						Mata Kuliah
					{/if}
				</span>
			</div>
			
			<a href="/dashboard/profile" class="header-profile">
				<div class="user-role-badge">
					{authState.role?.toUpperCase() || 'USER'}
				</div>
				<div class="user-avatar">
					{authState.profile?.name?.charAt(0).toUpperCase() || 'U'}
				</div>
			</a>
		</header>

		<!-- Main Content Slot -->
		<main class="content-area">
			{@render children()}
		</main>
	</div>
</div>

<style>
	.layout-wrapper {
		display: flex;
		min-height: 100vh;
		background: 
			radial-gradient(circle at 0% 0%, rgba(79, 70, 229, 0.08), transparent 40%),
			radial-gradient(circle at 100% 100%, rgba(16, 185, 129, 0.08), transparent 40%);
		background-attachment: fixed;
	}

	/* Sidebar */
	.sidebar {
		width: 280px;
		display: flex;
		flex-direction: column;
		border-radius: 0;
		border-top: none;
		border-bottom: none;
		border-left: none;
		position: sticky;
		top: 0;
		height: 100vh;
		z-index: 40;
		background: rgba(255, 255, 255, 0.7);
	}

	.sidebar-brand {
		padding: 32px 24px;
		display: flex;
		align-items: center;
		gap: 16px;
	}

	.brand-logo {
		width: 44px;
		height: 44px;
		background: linear-gradient(135deg, var(--primary-color), var(--primary-hover));
		border-radius: 12px;
		display: flex;
		align-items: center;
		justify-content: center;
		color: white;
		box-shadow: 0 8px 16px rgba(79, 70, 229, 0.3);
	}

	.brand-text h1 {
		font-size: 1.25rem;
		font-weight: 700;
		color: var(--text-main);
		line-height: 1.1;
		letter-spacing: -0.025em;
	}

	.brand-text p {
		font-size: 0.75rem;
		color: var(--primary-color);
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.1em;
		margin-top: 2px;
	}

	.sidebar-nav {
		flex: 1;
		padding: 16px;
		display: flex;
		flex-direction: column;
		gap: 4px;
	}

	.nav-section-title {
		font-size: 0.7rem;
		font-weight: 600;
		color: var(--text-muted);
		text-transform: uppercase;
		letter-spacing: 0.1em;
		padding: 24px 12px 8px 12px;
	}

	.nav-item {
		display: flex;
		align-items: center;
		gap: 12px;
		padding: 12px 16px;
		border-radius: var(--radius-md);
		color: var(--text-muted);
		text-decoration: none;
		font-weight: 500;
		transition: all 0.2s ease;
	}

	.nav-item:hover {
		background: rgba(79, 70, 229, 0.05);
		color: var(--primary-color);
	}

	.nav-item.active {
		background: rgba(79, 70, 229, 0.1);
		color: var(--primary-color);
		font-weight: 600;
	}

	.sidebar-footer {
		padding: 24px;
	}

	.btn-logout {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 8px;
		width: 100%;
		background: transparent;
		color: var(--text-muted);
		border: 1px solid var(--surface-border);
		padding: 12px;
		border-radius: var(--radius-md);
		font-family: inherit;
		font-size: 0.9rem;
		font-weight: 600;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-logout:hover {
		background: rgba(239, 68, 68, 0.05);
		color: var(--error-color);
		border-color: rgba(239, 68, 68, 0.2);
	}

	/* Main Area */
	.main-wrapper {
		flex: 1;
		display: flex;
		flex-direction: column;
		min-width: 0; /* Important for flex children overflow */
	}

	/* Top Header */
	.top-header {
		height: 72px;
		padding: 0 40px;
		display: flex;
		align-items: center;
		justify-content: space-between;
		border-radius: 0;
		border-top: none;
		border-right: none;
		border-left: none;
		position: sticky;
		top: 0;
		z-index: 30;
		background: rgba(255, 255, 255, 0.6);
	}

	.header-breadcrumb {
		display: flex;
		align-items: center;
		gap: 8px;
		font-size: 0.95rem;
		font-weight: 500;
	}

	.separator {
		color: var(--surface-border);
	}

	.current-page {
		color: var(--text-main);
		font-weight: 600;
	}

	.header-profile {
		display: flex;
		align-items: center;
		gap: 16px;
		text-decoration: none;
		cursor: pointer;
		padding: 4px 8px;
		border-radius: var(--radius-full);
		transition: background-color 0.2s;
	}

	.header-profile:hover {
		background-color: rgba(0, 0, 0, 0.03);
	}

	.user-role-badge {
		font-size: 0.75rem;
		font-weight: 700;
		color: var(--primary-color);
		background: rgba(79, 70, 229, 0.1);
		padding: 4px 12px;
		border-radius: var(--radius-full);
		letter-spacing: 0.05em;
	}

	.user-avatar {
		width: 40px;
		height: 40px;
		border-radius: 50%;
		background: linear-gradient(135deg, var(--text-main), var(--text-muted));
		color: white;
		display: flex;
		align-items: center;
		justify-content: center;
		font-weight: 600;
		font-size: 1.1rem;
		box-shadow: 0 4px 6px rgba(0,0,0,0.1);
	}

	/* Content Area */
	.content-area {
		flex: 1;
		padding: 40px;
		max-width: 1200px;
		width: 100%;
		margin: 0 auto;
	}
</style>
