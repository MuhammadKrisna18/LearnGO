<div class="skeleton-container">
	<div class="skeleton-line wide"></div>
	<div class="skeleton-line medium"></div>
	<div class="skeleton-line short"></div>
	<div class="skeleton-line wide"></div>
	<div class="skeleton-line medium"></div>
</div>

<style>
	.skeleton-container {
		display: flex;
		flex-direction: column;
		gap: 16px;
		padding: 24px;
	}

	.skeleton-line {
		height: 16px;
		border-radius: 6px;
		background: linear-gradient(90deg, #f0f0f0 25%, #e0e0e0 50%, #f0f0f0 75%);
		background-size: 200% 100%;
		animation: shimmer 1.5s ease-in-out infinite;
	}

	.skeleton-line.wide { width: 100%; }
	.skeleton-line.medium { width: 70%; }
	.skeleton-line.short { width: 40%; }

	@keyframes shimmer {
		0% { background-position: 200% 0; }
		100% { background-position: -200% 0; }
	}
</style>
