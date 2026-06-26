<script lang="ts">
	import { fade, scale } from 'svelte/transition';

	let {
		isOpen = $bindable(false),
		title = 'Konfirmasi Hapus',
		itemName = '',
		onConfirm,
		onCancel
	} = $props<{
		isOpen: boolean;
		title?: string;
		itemName: string;
		onConfirm: () => void;
		onCancel: () => void;
	}>();

	let generatedCode = $state('');
	let inputCode = $state('');

	$effect(() => {
		if (isOpen) {
			generatedCode = Math.floor(100000 + Math.random() * 900000).toString();
			inputCode = '';
		}
	});

	function handleConfirm() {
		if (inputCode === generatedCode) {
			onConfirm();
			isOpen = false;
		}
	}

	function handleCancel() {
		onCancel();
		isOpen = false;
	}
</script>

{#if isOpen}
	<div class="modal-overlay" in:fade={{ duration: 200 }} out:fade={{ duration: 150 }}>
		<div class="modal-content glass-panel" in:scale={{ duration: 200, start: 0.95 }} out:scale={{ duration: 150, start: 0.95 }}>
			<div class="modal-header">
				<h3>{title}</h3>
				<button class="close-btn" onclick={handleCancel}>&times;</button>
			</div>

			<div class="modal-body">
				<p>Anda yakin ingin menghapus <strong>{itemName}</strong>?</p>
				<p class="warning-text">Tindakan ini tidak dapat dibatalkan!</p>
				
				<div class="code-section">
					<p>Untuk melanjutkan, ketik 6 digit kode di bawah ini:</p>
					<div class="generated-code">{generatedCode}</div>
					
					<input 
						type="text" 
						class="form-input code-input" 
						bind:value={inputCode} 
						placeholder="Masukkan 6 digit kode" 
						maxlength="6"
					/>
				</div>
			</div>

			<div class="modal-footer">
				<button class="btn-cancel" onclick={handleCancel}>Batal</button>
				<button 
					class="btn-danger" 
					onclick={handleConfirm}
					disabled={inputCode !== generatedCode}
				>
					Hapus Permanen
				</button>
			</div>
		</div>
	</div>
{/if}

<style>
	.modal-overlay {
		position: fixed;
		top: 0;
		left: 0;
		width: 100vw;
		height: 100vh;
		background: rgba(15, 23, 42, 0.6);
		backdrop-filter: blur(4px);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
	}

	.modal-content {
		width: 100%;
		max-width: 400px;
		background: rgba(255, 255, 255, 0.95);
		border-radius: var(--radius-lg);
		box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
		overflow: hidden;
	}

	.modal-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 20px 24px;
		border-bottom: 1px solid var(--surface-border);
	}

	.modal-header h3 {
		font-size: 1.25rem;
		font-weight: 600;
		color: var(--error-color);
		margin: 0;
	}

	.close-btn {
		background: none;
		border: none;
		font-size: 1.5rem;
		color: var(--text-muted);
		cursor: pointer;
		padding: 0;
		line-height: 1;
	}
	
	.close-btn:hover {
		color: var(--text-main);
	}

	.modal-body {
		padding: 24px;
	}

	.modal-body p {
		margin-bottom: 8px;
		color: var(--text-main);
	}

	.warning-text {
		color: var(--error-color) !important;
		font-weight: 500;
		font-size: 0.9rem;
	}

	.code-section {
		margin-top: 24px;
		background: rgba(241, 245, 249, 0.5);
		padding: 16px;
		border-radius: var(--radius-md);
		border: 1px dashed var(--surface-border);
		text-align: center;
	}

	.code-section p {
		font-size: 0.85rem;
		color: var(--text-muted);
		margin-bottom: 12px;
	}

	.generated-code {
		font-size: 2rem;
		font-weight: 800;
		letter-spacing: 4px;
		color: var(--text-main);
		margin-bottom: 16px;
		font-family: monospace;
	}

	.code-input {
		text-align: center;
		letter-spacing: 2px;
		font-size: 1.1rem;
		font-weight: 600;
	}

	.modal-footer {
		padding: 16px 24px;
		background: rgba(248, 250, 252, 0.8);
		border-top: 1px solid var(--surface-border);
		display: flex;
		justify-content: flex-end;
		gap: 12px;
	}

	.btn-cancel {
		background: white;
		border: 1px solid var(--surface-border);
		padding: 10px 16px;
		border-radius: var(--radius-md);
		font-weight: 500;
		color: var(--text-main);
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-cancel:hover {
		background: #f1f5f9;
	}

	.btn-danger {
		background: var(--error-color);
		color: white;
		border: none;
		padding: 10px 16px;
		border-radius: var(--radius-md);
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-danger:hover:not(:disabled) {
		background: #b91c1c;
	}

	.btn-danger:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}
</style>
