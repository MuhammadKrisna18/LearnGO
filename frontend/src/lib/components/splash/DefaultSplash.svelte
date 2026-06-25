<div class="splash-content">
	<div class="spinner-large"></div>
	<h2>Welcome</h2>
	<p>Loading your workspace...</p>
</div>

<style>
	.splash-content {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
	}

	.splash-content h2 {
		font-size: 1.75rem;
		font-weight: 600;
		color: var(--text-main);
		margin-bottom: 8px;
	}

	.splash-content p {
		color: var(--text-muted);
		margin-bottom: 32px;
	}

	.spinner-large {
		width: 48px;
		height: 48px;
		border: 4px solid rgba(37, 99, 235, 0.2);
		border-radius: 50%;
		border-top-color: var(--primary-color);
		animation: spin 1s ease-in-out infinite;
		margin-bottom: 24px;
	}

	@keyframes spin {
		to { transform: rotate(360deg); }
	}
</style>
