<script lang="ts">
	let { stats }: { stats: { label: string; value: number; icon: string; color: string }[] } = $props();
</script>

<div class="stats-row">
	{#each stats as stat}
		<div class="stat-card card hover-lift">
			<div class="stat-icon" style="background: {stat.color}10; color: {stat.color};">
				{@html stat.icon}
			</div>
			<div class="stat-info">
				<span class="stat-value">{stat.value}</span>
				<span class="stat-label">{stat.label}</span>
			</div>
		</div>
	{/each}
</div>

<style>
	.stats-row {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
		gap: 16px;
		margin-bottom: 24px;
	}

	.stat-card {
		display: flex;
		align-items: center;
		gap: 16px;
		padding: 20px;
		border-radius: var(--radius-md);
	}

	.stat-icon {
		width: 48px;
		height: 48px;
		border-radius: var(--radius-sm);
		display: flex;
		align-items: center;
		justify-content: center;
		flex-shrink: 0;
	}

	.stat-info {
		display: flex;
		flex-direction: column;
	}

	.stat-value {
		font-size: 1.5rem;
		font-weight: 700;
		color: var(--text-main);
		line-height: 1;
	}

	.stat-label {
		font-size: 0.8rem;
		color: var(--text-muted);
		margin-top: 4px;
	}
</style>
