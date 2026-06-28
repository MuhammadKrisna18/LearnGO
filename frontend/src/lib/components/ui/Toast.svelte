<script lang="ts">
	import { toast } from '$lib/stores/toast.svelte';
	import { fly } from 'svelte/transition';
</script>

<div class="toast-container">
	{#each toast.toasts as t (t.id)}
		<div
			class="toast {t.type}"
			transition:fly={{ y: 20, duration: 300 }}
			role="alert"
		>
			<div class="toast-icon">
				{#if t.type === 'success'}
					<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
						<path d="M22 11.08V12a10 10 10 0 1 1-5.93-9.14"></path>
						<polyline points="22 4 12 14.01 9 11.01"></polyline>
					</svg>
				{:else if t.type === 'error'}
					<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
						<circle cx="12" cy="12" r="10"></circle>
						<line x1="15" y1="9" x2="9" y2="15"></line>
						<line x1="9" y1="9" x2="15" y2="15"></line>
					</svg>
				{:else}
					<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
						<circle cx="12" cy="12" r="10"></circle>
						<line x1="12" y1="16" x2="12" y2="12"></line>
						<line x1="12" y1="8" x2="12.01" y2="8"></line>
					</svg>
				{/if}
			</div>
			<div class="toast-message">{t.message}</div>
			<button class="toast-close" onclick={() => toast.remove(t.id)} aria-label="Close">
				&times;
			</button>
		</div>
	{/each}
</div>

<style>
	.toast-container {
		position: fixed;
		bottom: 24px;
		right: 24px;
		display: flex;
		flex-direction: column;
		gap: 12px;
		z-index: 9999;
		pointer-events: none;
	}

	.toast {
		pointer-events: auto;
		display: flex;
		align-items: center;
		gap: 12px;
		padding: 12px 16px;
		border-radius: var(--radius-md);
		background: rgba(255, 255, 255, 0.95);
		backdrop-filter: blur(8px);
		box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
		min-width: 300px;
		max-width: 400px;
		border: 1px solid var(--surface-border);
	}

	.toast.success {
		border-left: 4px solid var(--success-color, #10b981);
	}

	.toast.success .toast-icon {
		color: var(--success-color, #10b981);
	}

	.toast.error {
		border-left: 4px solid var(--danger-color, #ef4444);
	}

	.toast.error .toast-icon {
		color: var(--danger-color, #ef4444);
	}

	.toast.info {
		border-left: 4px solid var(--primary-color, #4f46e5);
	}

	.toast.info .toast-icon {
		color: var(--primary-color, #4f46e5);
	}

	.toast-icon {
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.toast-message {
		flex: 1;
		font-size: 0.9rem;
		font-weight: 500;
		color: var(--text-main);
	}

	.toast-close {
		background: none;
		border: none;
		font-size: 1.25rem;
		color: var(--text-muted);
		cursor: pointer;
		padding: 0 4px;
		line-height: 1;
		opacity: 0.7;
		transition: opacity 0.2s;
	}

	.toast-close:hover {
		opacity: 1;
		color: var(--text-main);
	}
</style>
