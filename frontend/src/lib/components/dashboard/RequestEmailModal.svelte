<script lang="ts">
	let { isOpen = $bindable(), onSubmit } = $props<{
		isOpen: boolean;
		onSubmit: (newEmail: string) => void;
	}>();

	let newEmail = $state('');

	function close() {
		isOpen = false;
		newEmail = '';
	}

	function handleSubmit(e: Event) {
		e.preventDefault();
		if (newEmail.trim()) {
			onSubmit(newEmail.trim());
			close();
		}
	}
</script>

{#if isOpen}
	<div class="modal-backdrop" onclick={close} role="presentation">
		<div class="modal-content animate-scale" onclick={(e) => e.stopPropagation()} role="dialog">
			<div class="modal-header">
				<h3>Ajukan Ganti Email</h3>
				<button class="btn-close" onclick={close}>&times;</button>
			</div>
			
			<form onsubmit={handleSubmit}>
				<div class="form-group">
					<label for="newEmail">Email Baru</label>
					<input 
						id="newEmail"
						type="email" 
						bind:value={newEmail} 
						class="form-input" 
						placeholder="Masukkan email baru..."
						required
					/>
					<p class="help-text">Email baru Anda membutuhkan persetujuan dari Admin sebelum aktif.</p>
				</div>
				
				<div class="modal-actions">
					<button type="button" class="btn-secondary" onclick={close}>Batal</button>
					<button type="submit" class="btn-primary" disabled={!newEmail}>Ajukan Perubahan</button>
				</div>
			</form>
		</div>
	</div>
{/if}

<style>
	.modal-backdrop {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(15, 23, 42, 0.4);
		backdrop-filter: blur(4px);
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1000;
	}

	.modal-content {
		background: white;
		width: 100%;
		max-width: 450px;
		border-radius: var(--radius-lg);
		box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
		border: 1px solid rgba(255, 255, 255, 0.2);
		overflow: hidden;
	}

	.modal-header {
		padding: 20px 24px;
		border-bottom: 1px solid var(--surface-border);
		display: flex;
		justify-content: space-between;
		align-items: center;
		background: var(--surface-light);
	}

	.modal-header h3 {
		font-size: 1.15rem;
		font-weight: 600;
		color: var(--text-main);
	}

	.btn-close {
		background: none;
		border: none;
		font-size: 1.5rem;
		color: var(--text-muted);
		cursor: pointer;
		line-height: 1;
		transition: color 0.2s;
	}

	.btn-close:hover {
		color: var(--text-main);
	}

	form {
		padding: 24px;
	}

	.form-group {
		margin-bottom: 24px;
	}

	.form-group label {
		display: block;
		font-size: 0.875rem;
		font-weight: 600;
		color: var(--text-main);
		margin-bottom: 8px;
	}

	.form-input {
		width: 100%;
		padding: 10px 14px;
		border: 1px solid var(--surface-border);
		border-radius: var(--radius-md);
		font-size: 0.95rem;
		transition: all 0.2s;
	}

	.form-input:focus {
		outline: none;
		border-color: var(--primary-color);
		box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
	}

	.help-text {
		font-size: 0.8rem;
		color: var(--text-muted);
		margin-top: 8px;
	}

	.modal-actions {
		display: flex;
		justify-content: flex-end;
		gap: 12px;
	}

	.btn-secondary {
		padding: 10px 20px;
		background: var(--surface-light);
		color: var(--text-main);
		border: 1px solid var(--surface-border);
		border-radius: var(--radius-md);
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-secondary:hover {
		background: var(--surface-border);
	}

	.btn-primary {
		padding: 10px 20px;
		background: var(--primary-color);
		color: white;
		border: none;
		border-radius: var(--radius-md);
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s;
	}

	.btn-primary:hover:not(:disabled) {
		background: var(--primary-hover);
		transform: translateY(-1px);
	}

	.btn-primary:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.animate-scale {
		animation: scaleIn 0.3s cubic-bezier(0.16, 1, 0.3, 1) forwards;
	}

	@keyframes scaleIn {
		from {
			opacity: 0;
			transform: scale(0.95) translateY(10px);
		}
		to {
			opacity: 1;
			transform: scale(1) translateY(0);
		}
	}
</style>
