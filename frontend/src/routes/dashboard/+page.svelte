<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	let profile: any = $state(null);
	let loading = $state(true);
	let error = $state('');

	let dosenName = $state('');
	let dosenUsername = $state('');
	let dosenPassword = $state('');
	let registerLoading = $state(false);
	let registerError = $state('');
	let registerSuccess = $state('');

	let dosenList: any[] = $state([]);
	let loadingDosen = $state(false);

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
				if (profile.role === 'admin') {
					fetchDosenList(token);
				}
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

	async function handleRegisterDosen(e: Event) {
		if (e) e.preventDefault();
		registerError = '';
		registerSuccess = '';
		registerLoading = true;

		const token = localStorage.getItem('token');

		try {
			const res = await fetch('http://localhost:8080/api/v1/auth/dosen', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					'Authorization': `Bearer ${token}`
				},
				body: JSON.stringify({
					name: dosenName,
					username: dosenUsername,
					password: dosenPassword
				})
			});

			const data = await res.json();

			if (res.ok && data.success) {
				registerSuccess = `Berhasil membuat akun Dosen! Email: ${dosenUsername}@DosenGO.id`;
				dosenName = '';
				dosenUsername = '';
				dosenPassword = '';
				fetchDosenList(token); // Refresh the list after successful registration
			} else {
				registerError = data.message || 'Gagal membuat akun Dosen';
			}
		} catch (err) {
			registerError = 'Gagal terhubung ke server';
		} finally {
			registerLoading = false;
		}
	}

	async function fetchDosenList(token: string) {
		loadingDosen = true;
		try {
			const res = await fetch('http://localhost:8080/api/v1/auth/dosen', {
				method: 'GET',
				headers: {
					'Authorization': `Bearer ${token}`
				}
			});
			const data = await res.json();
			if (res.ok && data.success) {
				dosenList = data.data || [];
			}
		} catch (err) {
			console.error("Failed to fetch dosen list", err);
		} finally {
			loadingDosen = false;
		}
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
		<button class="btn-logout" onclick={handleLogout}>Logout</button>
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
						<span class="detail-label">Email Address</span>
						<p>{profile.email}</p>
					</div>
					<div class="detail-group">
						<span class="detail-label">Account ID</span>
						<p class="mono">{profile.id}</p>
					</div>
					<div class="detail-group">
						<span class="detail-label">Member Since</span>
						<p>{new Date(profile.created_at).toLocaleDateString()}</p>
					</div>
				</div>
			</div>
			
			{#if profile.role === 'admin'}
				<div class="welcome-banner glass-panel animate-fade-in" style="animation-delay: 0.2s;">
					<h3>Dashboard Overview</h3>
					<p>Welcome to your secure dashboard! This page successfully fetched your data from the Go API using your JWT token.</p>
				</div>
			{:else if profile.role === 'dosen'}
				<div class="dosen-empty-state glass-panel animate-fade-in" style="animation-delay: 0.2s;">
					<div class="empty-icon">🎓</div>
					<h3>Selamat Datang, Dosen {profile.name}!</h3>
					<p>Belum ada jadwal atau kelas yang ditugaskan kepada Anda saat ini.</p>
					<button class="btn-primary empty-action-btn">Cek Notifikasi</button>
				</div>
			{/if}

			{#if profile.role === 'admin'}
				<div class="register-card glass-panel animate-fade-in" style="animation-delay: 0.3s;">
					<div class="card-header">
						<h3>Buat Akun Dosen</h3>
						<p>Tambahkan akun dosen baru ke dalam sistem.</p>
					</div>

					<form onsubmit={handleRegisterDosen}>
						<div class="form-group">
							<label class="form-label" for="dosenName">Nama Lengkap</label>
							<input
								class="form-input"
								type="text"
								id="dosenName"
								bind:value={dosenName}
								placeholder="Budi Santoso"
								required
								disabled={registerLoading}
							/>
						</div>

						<div class="form-group">
							<label class="form-label" for="dosenUsername">Username Email</label>
							<div class="input-with-suffix">
								<input
									class="form-input username-input"
									type="text"
									id="dosenUsername"
									bind:value={dosenUsername}
									placeholder="budi"
									required
									disabled={registerLoading}
								/>
								<span class="email-suffix">@DosenGO.id</span>
							</div>
						</div>

						<div class="form-group">
							<label class="form-label" for="dosenPassword">Password</label>
							<input
								class="form-input"
								type="password"
								id="dosenPassword"
								bind:value={dosenPassword}
								placeholder="••••••••"
								required
								disabled={registerLoading}
							/>
						</div>

						{#if registerError}
							<div class="error-message">
								{registerError}
							</div>
						{/if}

						{#if registerSuccess}
							<div class="success-message">
								{registerSuccess}
							</div>
						{/if}

						<button type="submit" class="btn-primary" disabled={registerLoading}>
							{#if registerLoading}
								Menyimpan...
							{:else}
								Buat Akun
							{/if}
						</button>
					</form>
				</div>

				<div class="dosen-list-card glass-panel animate-fade-in" style="animation-delay: 0.4s;">
					<div class="card-header">
						<h3>Daftar Akun Dosen</h3>
						<p>Daftar seluruh dosen yang terdaftar di dalam sistem.</p>
					</div>

					{#if loadingDosen}
						<div class="table-loading">Memuat data...</div>
					{:else if dosenList.length === 0}
						<div class="table-empty">Belum ada akun dosen yang didaftarkan.</div>
					{:else}
						<div class="table-container">
							<table class="dosen-table">
								<thead>
									<tr>
										<th>Nama</th>
										<th>Email</th>
										<th>Tgl Bergabung</th>
									</tr>
								</thead>
								<tbody>
									{#each dosenList as dosen}
										<tr>
											<td>{dosen.name}</td>
											<td><span class="mono">{dosen.email}</span></td>
											<td>{new Date(dosen.created_at).toLocaleDateString()}</td>
										</tr>
									{/each}
								</tbody>
							</table>
						</div>
					{/if}
				</div>
			{/if}
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

	.detail-group .detail-label {
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

	.register-card {
		padding: 32px;
		border-radius: var(--radius-lg);
		margin-top: 12px;
	}

	.card-header {
		margin-bottom: 24px;
	}

	.card-header h3 {
		font-size: 1.25rem;
		font-weight: 600;
		color: var(--text-main);
		margin-bottom: 4px;
	}

	.card-header p {
		font-size: 0.9rem;
		color: var(--text-muted);
	}

	.input-with-suffix {
		display: flex;
		align-items: center;
		background: #ffffff;
		border: 1px solid #cbd5e1;
		border-radius: var(--radius-sm);
		overflow: hidden;
		transition: all 0.3s ease;
	}

	.input-with-suffix:focus-within {
		border-color: var(--primary-color);
		box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.15);
	}

	.username-input {
		border: none;
		border-radius: 0;
		flex-grow: 1;
		outline: none;
		box-shadow: none !important;
	}

	.username-input:focus {
		box-shadow: none !important;
		border-color: transparent;
	}

	.email-suffix {
		padding: 12px 16px;
		background: #f1f5f9;
		color: var(--text-muted);
		font-weight: 500;
		border-left: 1px solid #cbd5e1;
		font-size: 0.9rem;
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

	.success-message {
		background: rgba(16, 185, 129, 0.1);
		color: var(--success-color);
		padding: 12px;
		border-radius: var(--radius-sm);
		border: 1px solid rgba(16, 185, 129, 0.2);
		margin-bottom: 20px;
		font-size: 0.875rem;
		text-align: center;
	}

	.dosen-list-card {
		padding: 32px;
		border-radius: var(--radius-lg);
		margin-top: 24px;
	}

	.table-container {
		overflow-x: auto;
		border: 1px solid var(--surface-border);
		border-radius: var(--radius-md);
	}

	.dosen-table {
		width: 100%;
		border-collapse: collapse;
		text-align: left;
	}

	.dosen-table th {
		background: rgba(241, 245, 249, 0.5);
		padding: 16px;
		font-size: 0.85rem;
		text-transform: uppercase;
		letter-spacing: 0.5px;
		color: var(--text-muted);
		border-bottom: 1px solid var(--surface-border);
	}

	.dosen-table td {
		padding: 16px;
		border-bottom: 1px solid var(--surface-border);
		font-size: 0.95rem;
	}

	.dosen-table tr:last-child td {
		border-bottom: none;
	}

	.dosen-table tr:hover {
		background: rgba(248, 250, 252, 0.5);
	}

	.table-loading, .table-empty {
		padding: 32px;
		text-align: center;
		color: var(--text-muted);
		background: rgba(248, 250, 252, 0.5);
		border-radius: var(--radius-md);
		border: 1px dashed var(--surface-border);
	}

	.dosen-empty-state {
		padding: 48px 32px;
		text-align: center;
		border-radius: var(--radius-lg);
		background: linear-gradient(135deg, rgba(37, 99, 235, 0.05), rgba(56, 189, 248, 0.05));
		border: 1px solid rgba(37, 99, 235, 0.1);
	}

	.empty-icon {
		font-size: 48px;
		margin-bottom: 16px;
		filter: drop-shadow(0 4px 12px rgba(37, 99, 235, 0.2));
	}

	.dosen-empty-state h3 {
		font-size: 1.5rem;
		color: var(--text-main);
		margin-bottom: 12px;
	}

	.dosen-empty-state p {
		color: var(--text-muted);
		margin-bottom: 24px;
		max-width: 400px;
		margin-left: auto;
		margin-right: auto;
	}

	.empty-action-btn {
		background: var(--primary-color);
		color: white;
		padding: 10px 24px;
		border-radius: var(--radius-md);
		font-weight: 500;
		border: none;
		cursor: pointer;
		transition: all 0.2s;
	}

	.empty-action-btn:hover {
		background: var(--primary-hover);
		transform: translateY(-1px);
		box-shadow: 0 4px 12px rgba(37, 99, 235, 0.2);
	}
</style>
