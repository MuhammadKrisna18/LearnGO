<script lang="ts">
	import Modal from '$lib/components/ui/Modal.svelte';

	let {
		isOpen = $bindable(false),
		title = "Konfirmasi Tindakan",
		message = "Masukkan kode berikut:",
		expectedCode = "",
		onConfirm,
		onCancel
	} = $props();

	let inputCode = $state("");
	let loading = $state(false);

	$effect(() => {
		if (isOpen) {
			inputCode = "";
		}
	});

	async function handleConfirm(e: Event) {
		e.preventDefault();
		if (inputCode !== expectedCode) return;
		
		loading = true;
		await onConfirm(inputCode);
		loading = false;
		isOpen = false;
	}

	function handleCancel() {
		if (onCancel) onCancel();
		isOpen = false;
	}
</script>

<Modal bind:isOpen {title}>
	<form class="prompt-content" onsubmit={handleConfirm}>
		<p class="prompt-message">
			{message} <strong class="expected-code">{expectedCode}</strong>
		</p>
		
		<div class="form-group">
			<input 
				type="text" 
				bind:value={inputCode} 
				class="form-input" 
				placeholder="Masukkan kode konfirmasi"
				required
				autocomplete="off"
				disabled={loading}
			/>
		</div>

		<div class="modal-actions">
			<button type="button" class="btn-secondary" onclick={handleCancel} disabled={loading}>
				Batal
			</button>
			<button type="submit" class="btn-primary" disabled={inputCode !== expectedCode || loading}>
				{#if loading}Memproses...{:else}Konfirmasi{/if}
			</button>
		</div>
	</form>

	{#snippet footer()}
		<div style="display: none;"></div>
	{/snippet}
</Modal>

<style>
	.prompt-content {
		padding: 10px 0;
	}

	.prompt-message {
		font-size: 1rem;
		color: var(--text-main);
		margin-bottom: 20px;
		line-height: 1.5;
	}

	.expected-code {
		font-family: monospace;
		font-size: 1.2rem;
		padding: 4px 8px;
		background: var(--bg-card);
		border-radius: 4px;
		color: var(--primary-color);
		letter-spacing: 2px;
	}

	.form-group {
		margin-bottom: 24px;
	}

	.modal-actions {
		display: flex;
		justify-content: flex-end;
		gap: 12px;
	}
</style>
