<script lang="ts">
	import type { User } from '$lib/types';
	let { profile } = $props<{ profile: User }>();
</script>

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
