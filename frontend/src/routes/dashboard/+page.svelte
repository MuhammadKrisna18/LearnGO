<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	let profile: any = null;
	let loading = true;
	let error = '';

	onMount(async () => {
		const token = localStorage.getItem('token');
		
		if (!token) {
			goto('/login');
			return;
		}

		try {
			const res = await fetch('http://localhost:8080/api/v1/auth/me', {
				method: 'GET',
				headers: {
					'Authorization': `Bearer ${token}`,
					'Content-Type': 'application/json'
				}
			});

			const data = await res.json();

			if (res.ok && data.success) {
				profile = data.data;
			} else {
				localStorage.removeItem('token');
				localStorage.removeItem('role');
				goto('/login');
			}
		} catch (err) {
			error = 'Failed to connect to the server.';
		} finally {
			loading = false;
		}
	});

	function handleLogout() {
		localStorage.removeItem('token');
		localStorage.removeItem('role');
		goto('/login');
	}
</script>

<svelte:head>
	<title>Dashboard - Modular Monolith</title>
</svelte:head>

<main class="dashboard-container">
	<nav class="glass-panel navbar animate-fade-in">
		<div class="nav-brand">
			<span class="logo-icon">✦</span>
			<span class="brand-name">Modular Monolith</span>
		</div>
		<button class="btn-logout" on:click={handleLogout}>Logout</button>
	</nav>

	<div class="content">
		{#if loading}
			<div class="loading-state glass-panel animate-fade-in">
				<div class="spinner"></div>
				<p>Loading your profile...</p>
			</div>
		{:else if error}
			<div class="error-message animate-fade-in">
				{error}
			</div>
		{:else if profile}
			<div class="profile-card glass-panel animate-fade-in" style="animation-delay: 0.1s;">
				<div class="profile-header">
					<div class="avatar">
						{profile.name.charAt(0).toUpperCase()}
					</div>
					<div class="profile-title">
						<h2>{profile.name}</h2>
						<span class="badge {profile.role}">{profile.role.toUpperCase()}</span>
					</div>
				</div>
				
				<div class="profile-details">
					<div class="detail-group">
						<label>Email Address</label>
						<p>{profile.email}</p>
					</div>
					<div class="detail-group">
						<label>Account ID</label>
						<p class="mono">{profile.id}</p>
					</div>
					<div class="detail-group">
						<label>Member Since</label>
						<p>{new Date(profile.created_at).toLocaleDateString()}</p>
					</div>
				</div>
			</div>
			
			<div class="welcome-banner glass-panel animate-fade-in" style="animation-delay: 0.2s;">
				<h3>Dashboard Overview</h3>
				<p>Welcome to your secure dashboard! This page successfully fetched your data from the Go API using your JWT token.</p>
			</div>
		{/if}
	</div>
</main>

<style>
	.dashboard-container {
		min-height: 100vh;
		padding: 24px;
	}

	.navbar {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 16px 24px;
		border-radius: var(--radius-md);
		margin-bottom: 32px;
	}

	.nav-brand {
		display: flex;
		align-items: center;
		gap: 12px;
	}

	.logo-icon {
		color: var(--primary-color);
		font-size: 20px;
	}

	.brand-name {
		font-weight: 600;
		font-size: 1.1rem;
		letter-spacing: 0.5px;
	}

	.btn-logout {
		background: #ffffff;
		color: var(--text-main);
		border: 1px solid var(--surface-border);
		padding: 8px 16px;
		border-radius: var(--radius-sm);
		font-family: inherit;
		font-size: 0.875rem;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-logout:hover {
		background: rgba(239, 68, 68, 0.1);
		border-color: rgba(239, 68, 68, 0.3);
		color: #ef4444;
	}

	.content {
		max-width: 800px;
		margin: 0 auto;
		display: flex;
		flex-direction: column;
		gap: 24px;
	}

	.loading-state {
		padding: 48px;
		text-align: center;
		border-radius: var(--radius-lg);
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 16px;
		color: var(--text-muted);
	}

	.spinner {
		width: 32px;
		height: 32px;
		border: 3px solid rgba(0, 0, 0, 0.1);
		border-top-color: var(--primary-color);
		border-radius: 50%;
		animation: spin 1s linear infinite;
	}

	.profile-card {
		padding: 32px;
		border-radius: var(--radius-lg);
	}

	.profile-header {
		display: flex;
		align-items: center;
		gap: 20px;
		margin-bottom: 32px;
		padding-bottom: 24px;
		border-bottom: 1px solid var(--surface-border);
	}

	.avatar {
		width: 64px;
		height: 64px;
		background: linear-gradient(135deg, var(--primary-color), #3b82f6);
		color: #ffffff;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 24px;
		font-weight: 700;
		box-shadow: 0 4px 12px rgba(37, 99, 235, 0.3);
	}

	.profile-title h2 {
		margin-bottom: 4px;
		font-weight: 600;
	}

	.badge {
		font-size: 0.7rem;
		padding: 2px 8px;
		border-radius: 12px;
		font-weight: 600;
		letter-spacing: 0.5px;
	}

	.badge.admin {
		background: rgba(37, 99, 235, 0.1);
		color: #2563eb;
		border: 1px solid rgba(37, 99, 235, 0.2);
	}

	.profile-details {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
		gap: 24px;
	}

	.detail-group label {
		display: block;
		font-size: 0.8rem;
		color: var(--text-muted);
		margin-bottom: 4px;
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	.detail-group p {
		font-size: 1.05rem;
		font-weight: 500;
	}

	.mono {
		font-family: monospace;
		color: var(--text-muted);
	}

	.welcome-banner {
		padding: 24px;
		border-radius: var(--radius-lg);
		background: linear-gradient(to right, rgba(37, 99, 235, 0.05), rgba(56, 189, 248, 0.05));
		border-left: 4px solid var(--primary-color);
	}

	.welcome-banner h3 {
		margin-bottom: 8px;
		color: var(--text-main);
	}

	.welcome-banner p {
		color: var(--text-muted);
		line-height: 1.5;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}
</style>
