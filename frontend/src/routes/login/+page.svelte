<script lang="ts">
	import { goto } from '$app/navigation';

	let email = '';
	let password = '';
	let loading = false;
	let error = '';

	async function handleLogin() {
		error = '';
		loading = true;

		try {
			const res = await fetch('http://localhost:8080/api/v1/auth/login', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ email, password })
			});

			const data = await res.json();

			if (!res.ok) {
				error = data.message || 'Login failed. Please try again.';
			} else if (data.success && data.data && data.data.token) {
				// Store the token and role in localStorage
				localStorage.setItem('token', data.data.token);
				localStorage.setItem('role', data.data.role);
				
				// Redirect to dashboard
				goto('/dashboard');
			} else {
				error = 'Invalid response from server.';
			}
		} catch (err) {
			error = 'Network error. Is the backend server running?';
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Login - Modular Monolith</title>
</svelte:head>

<main class="login-container">
	<div class="glass-panel login-card animate-fade-in">
		<div class="login-header">
			<div class="logo">
				<span class="logo-icon">✦</span>
			</div>
			<h1>Welcome Back</h1>
			<p>Please enter your details to sign in.</p>
		</div>

		<form on:submit|preventDefault={handleLogin}>
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
</style>
