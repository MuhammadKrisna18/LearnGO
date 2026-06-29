<script lang="ts">
	import { goto } from '$app/navigation';
	import AdminSplash from '$lib/components/splash/AdminSplash.svelte';
	import DefaultSplash from '$lib/components/splash/DefaultSplash.svelte';
	import DosenSplash from '$lib/components/splash/DosenSplash.svelte';
	import { authService } from '$lib/services/auth';
	import { authState } from '$lib/stores/auth.svelte';

	let email = $state('');
	let password = $state('');
	let loading = $state(false);
	let error = $state('');
	let showSplash = $state(false);

	async function handleLogin(e: Event) {
		if (e) e.preventDefault();
		error = '';
		loading = true;

		try {
			const res = await authService.login(email, password);

			if (res.success) {
				// Fetch profile so that the splash screen can display the user's name
				await authService.getProfile();

				showSplash = true;
				
				setTimeout(() => {
					goto('/dashboard');
				}, 1000);
			} else {
				error = res.message || 'Login failed. Please try again.';
			}
		} catch (err: any) {
			error = err.message || 'Network error. Is the backend server running?';
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Login - SIAKAD Pro</title>
</svelte:head>

<main class="login-container">
	{#if showSplash}
		<div class="glass-panel splash-card animate-fade-in">
			{#if authState.role?.toLowerCase() === 'admin'}
				<AdminSplash />
			{:else if authState.role?.toLowerCase() === 'dosen'}
				<DosenSplash />
			{:else}
				<DefaultSplash />
			{/if}
		</div>
	{:else}
		<div class="glass-panel login-card animate-fade-in">
		<div class="login-header">
			<div class="logo">
				<span class="logo-icon">✦</span>
			</div>
			<h1>Welcome Back</h1>
			<p>Please enter your details to sign in.</p>
		</div>

		<form onsubmit={handleLogin}>
			<div class="form-group">
				<label class="form-label" for="email">Email</label>
				<input
					class="form-input"
					type="email"
					id="email"
					bind:value={email}
					placeholder="adminGO@golang.id"
					required
					disabled={loading}
				/>
			</div>

			<div class="form-group">
				<label class="form-label" for="password">Password</label>
				<input
					class="form-input"
					type="password"
					id="password"
					bind:value={password}
					placeholder="••••••••"
					required
					disabled={loading}
				/>
			</div>

			{#if error}
				<div class="error-message">
					{error}
				</div>
			{/if}

			<button type="submit" class="btn-primary" disabled={loading}>
				{#if loading}
					<span class="spinner"></span> Signing in...
				{:else}
					Sign in to Dashboard
				{/if}
			</button>
		</form>
		</div>
	{/if}
</main>

<style>
	.login-container {
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 100vh;
		padding: 20px;
	}

	.login-card {
		width: 100%;
		max-width: 420px;
		padding: 40px;
		border-radius: var(--radius-lg);
	}

	.login-header {
		text-align: center;
		margin-bottom: 32px;
	}

	.logo {
		width: 48px;
		height: 48px;
		background: linear-gradient(135deg, var(--primary-color), #3b82f6);
		border-radius: 12px;
		display: flex;
		align-items: center;
		justify-content: center;
		margin: 0 auto 16px;
		box-shadow: 0 4px 12px rgba(37, 99, 235, 0.3);
	}

	.logo-icon {
		color: white;
		font-size: 24px;
	}

	.login-header h1 {
		font-size: 1.75rem;
		font-weight: 600;
		margin-bottom: 8px;
		color: var(--text-main);
	}

	.login-header p {
		color: var(--text-muted);
		font-size: 0.95rem;
	}

	.error-message {
		background: rgba(239, 68, 68, 0.1);
		color: var(--error-color);
		padding: 12px;
		border-radius: var(--radius-sm);
		border: 1px solid rgba(239, 68, 68, 0.2);
		margin-bottom: 20px;
		font-size: 0.875rem;
		text-align: center;
	}

	.spinner {
		display: inline-block;
		width: 16px;
		height: 16px;
		border: 2px solid rgba(255, 255, 255, 0.3);
		border-radius: 50%;
		border-top-color: #ffffff;
		animation: spin 1s ease-in-out infinite;
		margin-right: 8px;
		vertical-align: middle;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}

	.splash-card {
		width: 100%;
		max-width: 420px;
		padding: 40px;
		border-radius: var(--radius-lg);
		text-align: center;
	}
</style>
