<script lang="ts">
	import { onMount } from 'svelte';
	import { kelasService } from '$lib/services/kelas';
	import { programStudiService } from '$lib/services/programstudi';
	import type { ProgramStudi, Kelas } from '$lib/types';
	import { toast } from '$lib/stores/toast.svelte';

	let name = $state('');
	let capacity = $state(25);
	let program_studi_id = $state('');
	let hari = $state('');
	let waktu = $state('');

	const hariOptions = ['Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat'];
	const waktuOptions = [
		{ label: '07:00 - 09:00', start: '07:00', end: '09:00' },
		{ label: '09:00 - 11:00', start: '09:00', end: '11:00' },
		{ label: '11:00 - 13:00', start: '11:00', end: '13:00' },
		{ label: '13:00 - 15:00', start: '13:00', end: '15:00' },
		{ label: '15:00 - 17:00', start: '15:00', end: '17:00' },
		{ label: '17:00 - 19:00', start: '17:00', end: '19:00' },
		{ label: '19:00 - 21:00', start: '19:00', end: '21:00' }
	];

	const nameOptions: string[] = [];
	for (let i = 1; i <= 3; i++) {
		for (let j = 1; j <= 7; j++) {
			nameOptions.push(`IF-${i}0${j}`);
		}
	}

	let prodiList = $state<ProgramStudi[]>([]);
	let loadingProdi = $state(true);
	let loading = $state(false);
	let error = $state('');
	let success = $state('');

	let kelases = $state<Kelas[]>([]);

	onMount(async () => {
		try {
			const resProdi = await programStudiService.getList();
			if (resProdi.success && resProdi.data) {
				prodiList = resProdi.data;
			}
			
			const resKelas = await kelasService.getList();
			if (resKelas.success && resKelas.data) {
				kelases = resKelas.data;
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

			const selectedWaktu = waktuOptions.find(w => w.label === waktu);
			
			const res = await kelasService.create({
				name,
				capacity,
				hari,
				jam_mulai: selectedWaktu?.start || '',
				jam_selesai: selectedWaktu?.end || '',
				program_studi_id
			});

			if (res.success) {
				toast.success('Berhasil menambahkan kelas');
				name = '';
				capacity = 25;
				program_studi_id = '';
				hari = '';
				waktu = '';

				const resKelas = await kelasService.getList();
				if (resKelas.success && resKelas.data) {
					kelases = resKelas.data;
				}
				
				window.dispatchEvent(new CustomEvent('kelasCreated'));
			} else {
				toast.error(res.message || 'Gagal menambahkan kelas');
			}
		} catch (err: any) {
			toast.error(err.message || 'Terjadi kesalahan sistem');
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
				<label for="name">Nama Kelas</label>
				<select
					id="name"
					class="form-input"
					bind:value={name}
					required
					disabled={loading}
				>
					<option value="" disabled>-- Pilih Kelas --</option>
					{#each nameOptions as n}
						<option value={n}>{n}</option>
					{/each}
				</select>
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
				<label for="hari">Hari</label>
				<select
					id="hari"
					class="form-input"
					bind:value={hari}
					required
					disabled={loading}
				>
					<option value="" disabled>-- Pilih Hari --</option>
					{#each hariOptions as h}
						<option value={h}>{h}</option>
					{/each}
				</select>
			</div>

			<div class="form-group">
				<label for="waktu">Jam Kelas</label>
				<select
					id="waktu"
					class="form-input"
					bind:value={waktu}
					required
					disabled={loading || !name || !hari}
				>
					<option value="" disabled>-- Pilih Jam Kelas --</option>
					{#each waktuOptions as w}
						{@const isUsed = kelases.some(k => k.name === name && k.hari === hari && k.jam_mulai === w.start)}
						<option value={w.label} disabled={isUsed}>
							{w.label} {isUsed ? '(Penuh)' : ''}
						</option>
					{/each}
				</select>
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
