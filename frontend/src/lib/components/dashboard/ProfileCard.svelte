<script lang="ts">
	import { authService } from '$lib/services/auth';
	import { toast } from '$lib/stores/toast.svelte';
	import type { UserProfile } from '$lib/types';
	let { profile } = $props<{ profile: UserProfile }>();

	let fileInput: HTMLInputElement;
	let uploading = $state(false);

	async function handleFileChange(event: Event) {
		const target = event.target as HTMLInputElement;
		if (target.files && target.files.length > 0) {
			const file = target.files[0];
			if (file.size > 2 * 1024 * 1024) {
				toast.error('Ukuran file maksimal 2MB');
				return;
			}
			
			uploading = true;
			try {
				const res = await authService.uploadProfilePhoto(file);
				if (res.success) {
					toast.success('Foto profil berhasil diperbarui');
				} else {
					toast.error(res.message || 'Gagal mengunggah foto profil');
				}
			} catch (err: any) {
				toast.error(err.message || 'Gagal terhubung ke server');
			} finally {
				uploading = false;
				target.value = ''; 
			}
		}
	}

	function triggerUpload() {
		if (!uploading) {
			fileInput.click();
		}
	}
</script>

<div class="profile-card glass-panel animate-fade-in" style="animation-delay: 0.1s;">
	<div class="profile-header">

		<div class="avatar" class:uploading onclick={triggerUpload} title="Klik untuk mengubah foto profil">
			{#if profile.photo_url}
				<img src={`http://localhost:8080${profile.photo_url}`} alt={profile.name} class="avatar-img" />
			{:else}
				{profile.name.charAt(0).toUpperCase()}
			{/if}
			
			<div class="avatar-overlay">
				{#if uploading}
					<span class="spinner-small"></span>
				{:else}
					<span>📷 UBAH</span>
				{/if}
			</div>
		</div>
		<input 
			type="file" 
			accept=".jpg,.jpeg,.png" 
			bind:this={fileInput} 
			onchange={handleFileChange} 
			style="display: none;" 
		/>
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

<style>
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
		position: relative;
		overflow: hidden;
		cursor: pointer;
	}

	.avatar.uploading {
		opacity: 0.7;
		pointer-events: none;
	}

	.avatar-img {
		width: 100%;
		height: 100%;
		object-fit: cover;
		aspect-ratio: 1/1;
	}

	.avatar-overlay {
		position: absolute;
		inset: 0;
		background: rgba(0, 0, 0, 0.5);
		display: flex;
		align-items: center;
		justify-content: center;
		opacity: 0;
		transition: opacity 0.2s ease;
		color: white;
		font-size: 0.75rem;
		font-weight: 600;
		letter-spacing: 0.5px;
	}

	.avatar:hover .avatar-overlay {
		opacity: 1;
	}

	.spinner-small {
		width: 20px;
		height: 20px;
		border: 2px solid rgba(255,255,255,0.3);
		border-top-color: white;
		border-radius: 50%;
		animation: spin 1s linear infinite;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
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
</style>
