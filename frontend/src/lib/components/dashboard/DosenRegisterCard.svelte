<script lang="ts">
	import { dosenService } from '$lib/services/dosen';

	let dosenName = $state('');
	let dosenUsername = $state('');
	let dosenPassword = $state('');
	let registerLoading = $state(false);
	let registerError = $state('');
	let registerSuccess = $state('');

	async function handleRegisterDosen(e: Event) {
		if (e) e.preventDefault();
		registerError = '';
		registerSuccess = '';
		registerLoading = true;

		try {
			const res = await dosenService.register(dosenName, dosenUsername, dosenPassword);

			if (res.success) {
				registerSuccess = `Berhasil membuat akun Dosen! Email: ${dosenUsername}@DosenGO.id`;
				dosenName = '';
				dosenUsername = '';
				dosenPassword = '';
			} else {
				registerError = res.message || 'Gagal membuat akun Dosen';
			}
		} catch (err) {
			registerError = 'Gagal terhubung ke server';
		} finally {
			registerLoading = false;
		}
	}
</script>

<div class="register-card glass-panel animate-fade-in" style="animation-delay: 0.1s;">
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

<style>
	.register-card {
		padding: 32px;
		border-radius: var(--radius-lg);
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
</style>
