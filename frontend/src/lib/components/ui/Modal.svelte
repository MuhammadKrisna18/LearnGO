<script lang="ts">
	let {
		isOpen = $bindable(false),
		title,
		children,
		footer
	} = $props();

	function close() {
		isOpen = false;
	}
</script>

{#if isOpen}
	<div class="modal-backdrop" onclick={close} aria-hidden="true">
		<div class="modal-content glass-panel animate-fade-in" onclick={e => e.stopPropagation()} aria-hidden="true">
			<div class="modal-header">
				<h3>{title}</h3>
				<button class="btn-close" onclick={close}>✕</button>
			</div>
			
			<div class="modal-body">
				{@render children()}
			</div>
			
			{#if footer}
				<div class="modal-footer">
					{@render footer()}
				</div>
			{/if}
		</div>
	</div>
{/if}

<style>
	.modal-backdrop {
		position: fixed;
		top: 0;
		left: 0;
		width: 100vw;
		height: 100vh;
		background: rgba(15, 23, 42, 0.6);
		backdrop-filter: blur(4px);
		-webkit-backdrop-filter: blur(4px);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 100;
		padding: 24px;
	}

	.modal-content {
		width: 100%;
		max-width: 480px;
		padding: 0;
		display: flex;
		flex-direction: column;
		border-radius: var(--radius-lg);
		overflow: hidden;
	}

	.modal-header {
		padding: 20px 24px;
		border-bottom: 1px solid var(--surface-border);
		display: flex;
		align-items: center;
		justify-content: space-between;
	}

	.modal-header h3 {
		font-size: 1.25rem;
		font-weight: 600;
		color: var(--text-main);
		margin: 0;
	}

	.btn-close {
		background: transparent;
		border: none;
		font-size: 1.25rem;
		color: var(--text-muted);
		cursor: pointer;
		transition: color 0.2s;
	}

	.btn-close:hover {
		color: var(--error-color);
	}

	.modal-body {
		padding: 24px;
		display: flex;
		flex-direction: column;
		gap: 16px;
	}

	.modal-footer {
		padding: 16px 24px;
		border-top: 1px solid var(--surface-border);
		display: flex;
		justify-content: flex-end;
		gap: 12px;
		background: rgba(248, 250, 252, 0.5);
	}
</style>
