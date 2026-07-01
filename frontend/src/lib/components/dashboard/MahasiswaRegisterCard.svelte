<script lang="ts">
	import { onMount } from 'svelte';
	import { authService } from '$lib/services/auth';
	import { programStudiService } from '$lib/services/programstudi';
	import type { ProgramStudi } from '$lib/types';
	import { toast } from '$lib/stores/toast.svelte';

	let mhsName = $state('');
	let mhsNRP = $state('');
	let mhsPassword = $state('');
	let mhsProdiId = $state('');

	let prodiList = $state<ProgramStudi[]>([]);
	let loadingProdi = $state(true);

	let registerLoading = $state(false);
	let registerError = $state('');
	let registerSuccess = $state('');

	function getNRPPrefix(prodiName: string): string {
		if (!prodiName) return '';
		const name = prodiName.toLowerCase();
		if (name.includes('rpl') || name.includes('perangkat lunak')) return '50532410';
		if (name.includes('rka') || name.includes('artificial') || name.includes('artifisial') || name.includes('kecerdasan')) return '50542410';
		if (name.includes('informatika') || name.includes('tc')) return '50252410';
		return '';
	}

	let selectedProdi = $derived(prodiList.find(p => p.id === mhsProdiId));
	let nrpPrefix = $derived(selectedProdi ? getNRPPrefix(selectedProdi.name) : '');

	
	$effect(() => {
		if (nrpPrefix && !mhsNRP.startsWith(nrpPrefix)) {
			mhsNRP = nrpPrefix;
		}
	});

	onMount(async () => {
		await loadProdi();
	});

	async function loadProdi() {
		loadingProdi = true;
		try {
			const res = await programStudiService.getList();
			if (res.success && res.data) {
				prodiList = res.data;
			}
		} catch (err) {
			console.error("Gagal memuat daftar prodi:", err);
		} finally {
			loadingProdi = false;
		}
	}

	async function handleRegister(e: Event) {
		e.preventDefault();
		registerError = '';
		registerSuccess = '';
		
		if (!mhsName || !mhsNRP || !mhsPassword || !mhsProdiId) {
			registerError = 'Mohon lengkapi semua data';
			return;
		}

		if (mhsNRP.length !== 10 || !/^\d+$/.test(mhsNRP)) {
			registerError = 'NRP harus terdiri dari tepat 10 digit angka';
			return;
		}

		registerLoading = true;

		try {
			const res = await authService.registerMahasiswa({
				name: mhsName,
				nrp: mhsNRP,
				password: mhsPassword,
				program_studi_id: mhsProdiId
			});

			if (res.success) {
				registerSuccess = 'Akun Mahasiswa berhasil dibuat!';
				mhsName = '';
				mhsNRP = '';
				mhsPassword = '';
				mhsProdiId = '';
				toast.success('Akun Mahasiswa berhasil dibuat!');
				
				
				window.dispatchEvent(new CustomEvent('mahasiswa-added'));
			} else {
				registerError = res.message || 'Gagal mendaftarkan mahasiswa';
				toast.error(registerError);
			}
		} catch (err: any) {
			registerError = err.message || 'Terjadi kesalahan jaringan';
			toast.error(registerError);
		} finally {
			registerLoading = false;
			setTimeout(() => {
				registerSuccess = '';
				registerError = '';
			}, 3000);
		}
	}
</script>

<div class="card hover-lift register-card">
	<div class="card-header">
		<div class="header-icon">
			<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
				<path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/>
				<circle cx="9" cy="7" r="4"/>
				<line x1="19" y1="8" x2="19" y2="14"/>
				<line x1="16" y1="11" x2="22" y2="11"/>
			</svg>
		</div>
		<div>
			<h3>Buat Akun Mahasiswa</h3>
			<p>Daftarkan akun mahasiswa baru ke dalam sistem.</p>
		</div>
	</div>

	<form onsubmit={handleRegister} class="card-body">
		<div class="form-group">
			<label class="form-label" for="mhsProdi">Program Studi</label>
			<select 
				class="form-input" 
				id="mhsProdi" 
				bind:value={mhsProdiId} 
				required 
				disabled={registerLoading || loadingProdi}
			>
				<option value="" disabled selected>-- Pilih Program Studi --</option>
				{#each prodiList as prodi}
					<option value={prodi.id}>{prodi.name}</option>
				{/each}
			</select>
		</div>

		<div class="form-group">
			<label class="form-label" for="mhsName">Nama Lengkap</label>
			<input
				class="form-input"
				type="text"
				id="mhsName"
				bind:value={mhsName}
				placeholder="Cth: Budi Santoso"
				required
				disabled={registerLoading}
			/>
		</div>

		<div class="form-group">
			<label class="form-label" for="mhsNRP">NRP</label>
			<div class="input-with-suffix">
				<input
					class="form-input username-input"
					type="text"
					id="mhsNRP"
					bind:value={mhsNRP}
					placeholder={nrpPrefix ? `${nrpPrefix}**` : "50532410**"}
					required
					maxlength="10"
					disabled={registerLoading}
				/>
				<span class="email-suffix">@student.its.golang</span>
			</div>
		</div>

		<div class="form-group">
			<label class="form-label" for="mhsPassword">Password</label>
			<input
				class="form-input"
				type="password"
				id="mhsPassword"
				bind:value={mhsPassword}
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

		<button type="submit" class="btn btn-primary w-full" disabled={registerLoading || loadingProdi}>
			{#if registerLoading}
				<span class="spinner"></span> Memproses...
			{:else}
				Daftarkan Mahasiswa
			{/if}
		</button>
	</form>
</div>

<style>
	.register-card {
		position: sticky;
		top: 1.5rem;
		border: 1px solid var(--border);
	}

	.card-header {
		display: flex;
		align-items: flex-start;
		gap: 1rem;
		padding: 1.25rem 1.5rem;
		border-bottom: 1px solid var(--border);
		background: rgba(255, 255, 255, 0.5);
	}

	.header-icon {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 40px;
		height: 40px;
		border-radius: 10px;
		background: var(--primary-light);
		color: var(--primary);
		flex-shrink: 0;
	}

	.card-header h3 {
		margin: 0 0 0.25rem 0;
		font-size: 1.1rem;
		color: var(--text-main);
	}

	.card-header p {
		margin: 0;
		font-size: 0.85rem;
		color: var(--text-muted);
	}

	.card-body {
		padding: 1.5rem;
		display: flex;
		flex-direction: column;
		gap: 1.25rem;
	}

	.input-with-suffix {
		display: flex;
		align-items: center;
	}

	.username-input {
		border-top-right-radius: 0;
		border-bottom-right-radius: 0;
		border-right: none;
	}

	.email-suffix {
		display: inline-flex;
		align-items: center;
		padding: 0 1rem;
		background-color: var(--surface);
		border: 1px solid var(--border);
		border-left: none;
		border-top-right-radius: 8px;
		border-bottom-right-radius: 8px;
		color: var(--text-muted);
		font-size: 0.9rem;
		height: 40px;
		box-sizing: border-box;
		white-space: nowrap;
	}

	.error-message, .success-message {
		padding: 0.75rem 1rem;
		border-radius: 8px;
		font-size: 0.9rem;
	}

	.error-message {
		background-color: #fee2e2;
		color: #b91c1c;
		border: 1px solid #fecaca;
	}

	.success-message {
		background-color: #dcfce3;
		color: #15803d;
		border: 1px solid #bbf7d0;
	}

	.w-full {
		width: 100%;
	}
</style>
