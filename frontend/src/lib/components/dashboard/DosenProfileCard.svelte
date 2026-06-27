<script lang="ts">
	import { authService } from '$lib/services/auth';
	import { programStudiService } from '$lib/services/programstudi';
	import type { UserProfile, ProgramStudi } from '$lib/types';
	import { onMount } from 'svelte';
	import RequestEmailModal from './RequestEmailModal.svelte';

	let { profile } = $props<{ profile: UserProfile }>();

	let isEditing = $state(false);
	let isEmailModalOpen = $state(false);

	// Form states
	let editName = $state(profile.name);
	let editNickname = $state(profile.nickname || '');

	let loading = $state(false);
	let error = $state('');
	let success = $state('');

	let prodiList = $state<ProgramStudi[]>([]);
	let loadingProdi = $state(true);

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

	function toggleEdit() {
		isEditing = !isEditing;
		if (!isEditing) {
			// reset forms
			editName = profile.name;
			editNickname = profile.nickname || '';
			error = '';
			success = '';
		}
	}

	async function handleSave() {
		loading = true;
		error = '';
		success = '';

		try {
			const res = await authService.updateProfile({
				name: editName,
				nickname: editNickname
			});

			if (res.success) {
				success = 'Profil berhasil diperbarui!';
				isEditing = false;
			} else {
				error = res.message || 'Gagal memperbarui profil';
			}
		} catch (err: any) {
			error = 'Gagal terhubung ke server';
		} finally {
			loading = false;
		}
	}

	async function handleEmailRequest(newEmail: string) {
		error = '';
		success = '';
		try {
			const res = await authService.requestEmailChange(newEmail);
			if (res.success) {
				success = 'Permintaan ganti email berhasil diajukan dan menunggu persetujuan.';
			} else {
				error = res.message || 'Gagal mengajukan ganti email';
			}
		} catch (err) {
			error = 'Terjadi kesalahan sistem';
		}
	}
</script>

<div class="dosen-profile-card glass-panel animate-fade-in">
	<div class="profile-header">
		<div class="avatar-large">{profile.name.charAt(0).toUpperCase()}</div>
		<div class="header-info">
			<h2>{profile.name}</h2>
			<div class="badges">
				<span class="role-badge">Dosen</span>
				{#if profile.nid}
					<span class="nid-badge">NID: {profile.nid}</span>
				{/if}
			</div>
		</div>
		<div class="header-actions">
			<button class="btn-edit" onclick={toggleEdit}>
				{isEditing ? 'Batal' : 'Edit Profil'}
			</button>
		</div>
	</div>

	{#if error}
		<div class="error-message">{error}</div>
	{/if}
	{#if success}
		<div class="success-message">{success}</div>
	{/if}

	<div class="profile-content">
		<div class="info-group">
			<label>Nama Lengkap</label>
			{#if isEditing}
				<input class="form-input" type="text" bind:value={editName} disabled={loading} required/>
			{:else}
				<p>{profile.name}</p>
			{/if}
		</div>

		<div class="info-group">
			<label>Nama Panggilan</label>
			{#if isEditing}
				<input class="form-input" type="text" bind:value={editNickname} disabled={loading} placeholder="Belum diatur"/>
			{:else}
				<p>{profile.nickname || '-'}</p>
			{/if}
		</div>

		<div class="info-group">
			<label>Nomer Induk Dosen (NID)</label>
			<p class="mono">{profile.nid || 'Belum di-generate'}</p>
		</div>

		<div class="info-group email-group">
			<label>Email</label>
			<div class="email-display">
				<p>{profile.email}</p>
				<button class="btn-link" onclick={() => isEmailModalOpen = true}>Ajukan Ganti Email</button>
			</div>
			{#if profile.pending_email}
				<div class="pending-notice">
					Menunggu persetujuan Admin untuk mengganti email ke: <strong>{profile.pending_email}</strong>
				</div>
			{/if}
		</div>

		<div class="info-group">
			<label>Program Studi</label>
			<p>{profile.program_studi ? profile.program_studi.name : 'Belum diatur'}</p>
		</div>

		{#if isEditing}
			<div class="form-actions">
				<button class="btn-primary" onclick={handleSave} disabled={loading}>
					{loading ? 'Menyimpan...' : 'Simpan Perubahan'}
				</button>
			</div>
		{/if}
	</div>
</div>

<RequestEmailModal 
	bind:isOpen={isEmailModalOpen}
	onSubmit={handleEmailRequest}
/>

<style>
	.dosen-profile-card {
		padding: 32px;
		border-radius: var(--radius-lg);
	}

	.profile-header {
		display: flex;
		align-items: center;
		gap: 24px;
		margin-bottom: 32px;
		padding-bottom: 24px;
		border-bottom: 1px solid var(--surface-border);
		flex-wrap: wrap;
	}

	.avatar-large {
		width: 80px;
		height: 80px;
		border-radius: 50%;
		background: linear-gradient(135deg, var(--primary-color), var(--secondary-color));
		color: white;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 2.5rem;
		font-weight: 700;
		box-shadow: 0 8px 16px rgba(37, 99, 235, 0.2);
	}

	.header-info {
		flex: 1;
	}

	.header-info h2 {
		font-size: 1.5rem;
		color: var(--text-main);
		margin-bottom: 8px;
	}

	.badges {
		display: flex;
		gap: 8px;
	}

	.role-badge {
		background: rgba(37, 99, 235, 0.1);
		color: var(--primary-color);
		padding: 4px 12px;
		border-radius: var(--radius-full);
		font-size: 0.85rem;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.5px;
	}

	.nid-badge {
		background: rgba(16, 185, 129, 0.1);
		color: #10b981;
		padding: 4px 12px;
		border-radius: var(--radius-full);
		font-size: 0.85rem;
		font-weight: 600;
		letter-spacing: 0.5px;
	}

	.header-actions {
		display: flex;
		gap: 12px;
	}

	.btn-edit {
		background: var(--surface-light);
		border: 1px solid var(--surface-border);
		padding: 8px 16px;
		border-radius: var(--radius-md);
		cursor: pointer;
		font-weight: 500;
		transition: all 0.2s;
	}

	.btn-edit:hover {
		background: var(--surface-border);
	}

	.btn-delete-account {
		background: rgba(239, 68, 68, 0.1);
		color: var(--error-color);
		border: 1px solid rgba(239, 68, 68, 0.2);
		padding: 8px 16px;
		border-radius: var(--radius-md);
		cursor: pointer;
		font-weight: 500;
		transition: all 0.2s;
	}

	.btn-delete-account:hover {
		background: rgba(239, 68, 68, 0.2);
	}

	.profile-content {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 24px;
	}

	.info-group {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.info-group label {
		font-size: 0.85rem;
		color: var(--text-muted);
		text-transform: uppercase;
		letter-spacing: 0.5px;
		font-weight: 600;
	}

	.info-group p {
		font-size: 1.05rem;
		color: var(--text-main);
		font-weight: 500;
	}

	.mono {
		font-family: monospace;
	}

	.email-group {
		grid-column: 1 / -1;
	}

	.email-display {
		display: flex;
		align-items: center;
		gap: 16px;
	}

	.btn-link {
		background: none;
		border: none;
		color: var(--primary-color);
		font-size: 0.85rem;
		font-weight: 600;
		cursor: pointer;
		text-decoration: underline;
		padding: 0;
	}

	.pending-notice {
		margin-top: 8px;
		font-size: 0.85rem;
		color: #eab308;
		background: rgba(234, 179, 8, 0.1);
		padding: 8px 12px;
		border-radius: var(--radius-sm);
		border: 1px solid rgba(234, 179, 8, 0.2);
	}

	.form-actions {
		grid-column: 1 / -1;
		display: flex;
		justify-content: flex-end;
		margin-top: 16px;
	}

	.error-message {
		background: rgba(239, 68, 68, 0.1);
		color: var(--error-color);
		padding: 12px;
		border-radius: var(--radius-sm);
		border: 1px solid rgba(239, 68, 68, 0.2);
		margin-bottom: 20px;
		text-align: center;
	}

	.success-message {
		background: rgba(16, 185, 129, 0.1);
		color: var(--success-color);
		padding: 12px;
		border-radius: var(--radius-sm);
		border: 1px solid rgba(16, 185, 129, 0.2);
		margin-bottom: 20px;
		text-align: center;
	}

	@media (max-width: 640px) {
		.profile-content {
			grid-template-columns: 1fr;
		}
	}
</style>
