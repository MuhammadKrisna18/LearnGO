<div class="splash-content admin-splash">
	<div class="role-icon">
		<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
			<path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
		</svg>
	</div>
	<h2>Welcome, Administrator</h2>
	<p>Preparing your dashboard...</p>
	<div class="progress-bar">
		<div class="progress-fill"></div>
	</div>
</div>

<style>
	.splash-content {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
	}

	.role-icon {
		width: 80px;
		height: 80px;
		background: linear-gradient(135deg, var(--primary-color), #3b82f6);
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		margin-bottom: 24px;
		color: white;
		box-shadow: 0 8px 24px rgba(37, 99, 235, 0.4);
		animation: pulse 2s infinite;
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

	.progress-bar {
		width: 100%;
		height: 6px;
		background: rgba(255, 255, 255, 0.1);
		border-radius: 3px;
		overflow: hidden;
	}

	.progress-fill {
		height: 100%;
		background: var(--primary-color);
		width: 0%;
		animation: load 3s ease-out forwards;
	}

	@keyframes pulse {
		0% { transform: scale(1); box-shadow: 0 8px 24px rgba(37, 99, 235, 0.4); }
		50% { transform: scale(1.05); box-shadow: 0 12px 32px rgba(37, 99, 235, 0.6); }
		100% { transform: scale(1); box-shadow: 0 8px 24px rgba(37, 99, 235, 0.4); }
	}

	@keyframes load {
		0% { width: 0%; }
		20% { width: 40%; }
		80% { width: 90%; }
		100% { width: 100%; }
	}
</style>
