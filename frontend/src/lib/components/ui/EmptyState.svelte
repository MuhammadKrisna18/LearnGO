<script lang="ts">
	let { message = 'Tidak ada data.' }: { message?: string } = $props();
</script>

<div class="empty-state">
	<div class="empty-icon">
		<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
			<path d="M21 8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16Z"/>
			<path d="m3.3 7 8.7 5 8.7-5"/>
			<path d="M12 22V12"/>
		</svg>
	</div>
	<p>{message}</p>
</div>

<style>
	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 48px 24px;
		gap: 16px;
	}

	.empty-icon {
		width: 80px;
		height: 80px;
		border-radius: 50%;
		background: #f9fafb;
		border: 1px dashed var(--surface-border);
		display: flex;
		align-items: center;
		justify-content: center;
		color: #d1d5db;
	}

	p {
		color: var(--text-muted);
		font-size: 0.9rem;
		margin: 0;
	}
</style>
