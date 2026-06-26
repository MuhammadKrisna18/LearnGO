<script lang="ts">
	import { onMount } from 'svelte';
	import { matakuliahService } from '$lib/services/matakuliah';
	import { programStudiService } from '$lib/services/programstudi';
	import type { ProgramStudi } from '$lib/types';

	let mkName = $state('');
	let mkSks = $state<number | ''>('');
	let mkProdiId = $state('');
	
	let prodiList = $state<ProgramStudi[]>([]);
	let loadingProdi = $state(true);
	
	let registerLoading = $state(false);
	let registerError = $state('');
	let registerSuccess = $state('');

	onMount(async () => {
		try {
			const res = await programStudiService.getList();
			if (res.success && res.data) {
				prodiList = res.data;
			}
		} catch (err) {
			console.error("Gagal mengambil data prodi:", err);
		} finally {
			loadingProdi = false;
		}
	});

	async function handleRegisterMataKuliah(e: Event) {
		if (e) e.preventDefault();
		registerError = '';
		registerSuccess = '';
		
		if (typeof mkSks !== 'number' || mkSks < 1) {
			registerError = 'SKS minimal adalah 1';
			return;
		}
		
		if (!mkProdiId) {
			registerError = 'Program Studi wajib dipilih';
			return;
		}
		
		registerLoading = true;

		try {
			const res = await matakuliahService.create(mkName, mkSks, mkProdiId);

			if (res.success) {
				registerSuccess = `Berhasil membuat Mata Kuliah baru: ${mkName} (${mkSks} SKS)`;
				mkName = '';
				mkSks = '';
				mkProdiId = '';
			} else {
				registerError = res.message || 'Gagal membuat Mata Kuliah';
			}
		} catch (err: any) {
			registerError = err.message || 'Gagal terhubung ke server';
		} finally {
			registerLoading = false;
		}
	}
</script>

<div class="register-card glass-panel animate-fade-in" style="animation-delay: 0.1s;">
	<div class="card-header">
		<h3>Tambah Mata Kuliah</h3>
		<p>Tambahkan mata kuliah baru ke dalam sistem.</p>
	</div>

	<form onsubmit={handleRegisterMataKuliah}>
		<div class="form-group">
			<label class="form-label" for="mkName">Nama Mata Kuliah</label>
			<input
				class="form-input"
				type="text"
				id="mkName"
				bind:value={mkName}
				placeholder="Contoh: Pemrograman Web"
				required
				disabled={registerLoading}
			/>
		</div>

		<div class="form-group">
			<label class="form-label" for="mkSks">SKS (Satuan Kredit Semester)</label>
			<input
				class="form-input"
				type="number"
				id="mkSks"
				bind:value={mkSks}
				placeholder="Contoh: 3"
				min="1"
				max="6"
				required
				disabled={registerLoading}
			/>
		</div>

		<div class="form-group">
			<label class="form-label" for="mkProdi">Program Studi</label>
			<select 
				class="form-input" 
				id="mkProdi" 
				bind:value={mkProdiId} 
				required 
				disabled={registerLoading || loadingProdi}
			>
				<option value="" disabled selected>-- Pilih Program Studi --</option>
				{#each prodiList as prodi}
					<option value={prodi.id}>{prodi.name}</option>
				{/each}
			</select>
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
				Tambah Mata Kuliah
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
