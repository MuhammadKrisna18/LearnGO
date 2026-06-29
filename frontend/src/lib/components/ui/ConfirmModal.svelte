<script lang="ts">
	let {
		isOpen = $bindable(false),
		title = 'Konfirmasi',
		message = 'Apakah Anda yakin?',
		confirmText = 'Ya, Lanjutkan',
		cancelText = 'Batal',
		variant = 'danger',
		onConfirm
	}: {
		isOpen: boolean;
		title?: string;
		message?: string;
		confirmText?: string;
		cancelText?: string;
		variant?: 'danger' | 'warning';
		onConfirm: () => void;
	} = $props();

	function close() {
		isOpen = false;
	}

	function handleConfirm() {
		onConfirm();
		close();
	}
</script>

{#if isOpen}
	<div class="confirm-backdrop" onclick={close} aria-hidden="true">
		<div class="confirm-panel animate-fade-in" onclick={e => e.stopPropagation()} aria-hidden="true">
			<div class="confirm-icon" class:danger={variant === 'danger'} class:warning={variant === 'warning'}>
				{#if variant === 'danger'}
					<svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/><line x1="10" x2="10" y1="11" y2="17"/><line x1="14" x2="14" y1="11" y2="17"/></svg>
				{:else}
					<svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"/><line x1="12" x2="12" y1="9" y2="13"/><line x1="12" x2="12.01" y1="17" y2="17"/></svg>
				{/if}
			</div>
			<h3>{title}</h3>
			<p>{message}</p>
			<div class="confirm-actions">
				<button class="btn-cancel" onclick={close}>{cancelText}</button>
				<button class="btn-confirm" class:danger={variant === 'danger'} class:warning={variant === 'warning'} onclick={handleConfirm}>{confirmText}</button>
			</div>
		</div>
	</div>
{/if}

<style>
	.confirm-backdrop {
		position: fixed;
		top: 0;
		left: 0;
		width: 100vw;
		height: 100vh;
		background: rgba(0, 0, 0, 0.5);
		backdrop-filter: blur(4px);
		-webkit-backdrop-filter: blur(4px);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 200;
		padding: 24px;
	}

	.confirm-panel {
		background: white;
		border: 1px solid var(--surface-border);
		border-radius: var(--radius-lg);
		padding: 32px;
		width: 100%;
		max-width: 400px;
		text-align: center;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 12px;
	}

	.confirm-icon {
		width: 56px;
		height: 56px;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		margin-bottom: 4px;
	}

	.confirm-icon.danger {
		background: rgba(239, 68, 68, 0.1);
		color: var(--error-color);
	}

	.confirm-icon.warning {
		background: rgba(245, 158, 11, 0.1);
		color: var(--warning-color);
	}

	h3 {
		font-size: 1.1rem;
		font-weight: 700;
		color: var(--text-main);
		margin: 0;
	}

	p {
		font-size: 0.9rem;
		color: var(--text-muted);
		line-height: 1.5;
		margin: 0;
	}

	.confirm-actions {
		display: flex;
		gap: 12px;
		width: 100%;
		margin-top: 8px;
	}

	.btn-cancel {
		flex: 1;
		padding: 10px 16px;
		background: white;
		border: 1px solid var(--surface-border);
		border-radius: var(--radius-sm);
		font-family: inherit;
		font-size: 0.9rem;
		font-weight: 600;
		cursor: pointer;
		color: var(--text-main);
		transition: all 0.2s;
	}

	.btn-cancel:hover {
		background: #f9fafb;
	}

	.btn-confirm {
		flex: 1;
		padding: 10px 16px;
		border: none;
		border-radius: var(--radius-sm);
		font-family: inherit;
		font-size: 0.9rem;
		font-weight: 600;
		cursor: pointer;
		color: white;
		transition: all 0.2s;
	}

	.btn-confirm.danger {
		background: var(--error-color);
	}

	.btn-confirm.danger:hover {
		background: #dc2626;
	}

	.btn-confirm.warning {
		background: var(--warning-color);
	}

	.btn-confirm.warning:hover {
		background: #d97706;
	}
</style>
