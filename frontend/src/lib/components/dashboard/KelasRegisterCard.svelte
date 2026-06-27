<script lang="ts">
	import { onMount } from 'svelte';
	import { kelasService } from '$lib/services/kelas';
	import { programStudiService } from '$lib/services/programstudi';
	import type { ProgramStudi } from '$lib/types';

	let name = $state('');
	let capacity = $state(25);
	let program_studi_id = $state('');

	let prodiList = $state<ProgramStudi[]>([]);
	let loadingProdi = $state(true);
	let loading = $state(false);
	let error = $state('');
	let success = $state('');

	onMount(async () => {
		try {
			const res = await programStudiService.getList();
			if (res.success && res.data) {
				prodiList = res.data;
			}
		} catch (err) {
			console.error('Gagal mengambil data program studi:', err);
		} finally {
			loadingProdi = false;
		}
	});

	async function handleSubmit(e: Event) {
		e.preventDefault();
		error = '';
		success = '';
		loading = true;

		try {
			const res = await kelasService.create({
				name,
				capacity,
				program_studi_id
			});

			if (res.success) {
				success = 'Kelas berhasil ditambahkan!';
				name = '';
				capacity = 25;
				program_studi_id = '';
				// Dispatch event to refresh list
				window.dispatchEvent(new CustomEvent('kelasCreated'));
			} else {
				error = res.message || 'Gagal menambahkan kelas';
			}
		} catch (err: any) {
			error = err.message || 'Terjadi kesalahan sistem';
		} finally {
			loading = false;
		}
	}
</script>

<div class="card glass-panel animate-fade-in" style="animation-delay: 0.1s;">
	<div class="card-header">
		<h2>Tambah Kelas</h2>
	</div>

	<div class="card-body">
		{#if error}
			<div class="alert error">{error}</div>
		{/if}
		{#if success}
			<div class="alert success">{success}</div>
		{/if}

		<form onsubmit={handleSubmit} class="form-grid">
			<div class="form-group">
				<label for="name">Nama Kelas (Format: IF-101)</label>
				<input
					type="text"
					id="name"
					class="form-input"
					bind:value={name}
					placeholder="Cth: IF-101"
					required
					disabled={loading}
					pattern="^IF-[1-3]0[1-7]$"
					title="Format wajib: IF-101 s/d IF-107, IF-201 s/d IF-207, atau IF-301 s/d IF-307"
				/>
			</div>

			<div class="form-group">
				<label for="capacity">Kapasitas (25 - 50)</label>
				<input
					type="number"
					id="capacity"
					class="form-input"
					bind:value={capacity}
					min="25"
					max="50"
					required
					disabled={loading}
				/>
			</div>

			<div class="form-group">
				<label for="program_studi">Program Studi</label>
				<select
					id="program_studi"
					class="form-input"
					bind:value={program_studi_id}
					required
					disabled={loading || loadingProdi}
				>
					<option value="" disabled>-- Pilih Program Studi --</option>
					{#each prodiList as prodi}
						<option value={prodi.id}>{prodi.name}</option>
					{/each}
				</select>
			</div>

			<button type="submit" class="btn-primary" disabled={loading}>
				{loading ? 'Menambahkan...' : 'Tambah Kelas'}
			</button>
		</form>
	</div>
</div>

<style>
	.card {
		background: var(--surface-light);
		border-radius: var(--radius-lg);
		border: 1px solid var(--surface-border);
		position: sticky;
		top: 24px;
	}

	.card-header {
		padding: 20px 24px;
		border-bottom: 1px solid var(--surface-border);
		background: rgba(255, 255, 255, 0.02);
	}

	.card-header h2 {
		font-size: 1.25rem;
		font-weight: 600;
		margin: 0;
	}

	.card-body {
		padding: 24px;
	}

	.form-grid {
		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	.form-group {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.form-group label {
		font-size: 0.85rem;
		font-weight: 600;
		color: var(--text-muted);
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	.form-input {
		width: 100%;
		padding: 10px 14px;
		background: rgba(255, 255, 255, 0.03);
		border: 1px solid var(--surface-border);
		border-radius: var(--radius-md);
		color: var(--text-main);
		font-size: 0.95rem;
		transition: all 0.2s;
	}

	.form-input:focus {
		outline: none;
		border-color: var(--primary-color);
		box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.2);
	}

	.form-input:disabled {
		opacity: 0.7;
		cursor: not-allowed;
	}

	.alert {
		padding: 12px;
		border-radius: var(--radius-md);
		margin-bottom: 20px;
		font-size: 0.9rem;
	}

	.alert.error {
		background: rgba(239, 68, 68, 0.1);
		color: var(--error-color);
		border: 1px solid rgba(239, 68, 68, 0.2);
	}

	.alert.success {
		background: rgba(16, 185, 129, 0.1);
		color: var(--success-color);
		border: 1px solid rgba(16, 185, 129, 0.2);
	}

	.btn-primary {
		width: 100%;
		margin-top: 8px;
		padding: 12px;
		font-size: 1rem;
		font-weight: 500;
	}
</style>
