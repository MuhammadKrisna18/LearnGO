<script lang="ts">
	import Modal from '$lib/components/ui/Modal.svelte';

	let {
		isOpen = $bindable(false),
		title = "Konfirmasi Hapus",
		itemName = "",
		onConfirm,
		onCancel
	} = $props();

	let loading = $state(false);

	async function handleConfirm() {
		loading = true;
		await onConfirm();
		loading = false;
		isOpen = false;
	}

	function handleCancel() {
		if (onCancel) onCancel();
		isOpen = false;
	}
</script>

<Modal bind:isOpen {title}>
	<div class="confirm-content">
		<div class="warning-icon">⚠️</div>
		<p class="confirm-message">
			Apakah Anda yakin ingin menghapus <strong>{itemName}</strong>?
		</p>
		<p class="confirm-subtext">
			Tindakan ini tidak dapat dibatalkan. Data yang telah dihapus akan hilang secara permanen.
		</p>
	</div>

	{#snippet footer()}
		<button class="btn-secondary" onclick={handleCancel} disabled={loading}>
			Batal
		</button>
		<button class="btn-delete" onclick={handleConfirm} disabled={loading}>
			{#if loading}Menghapus...{:else}Ya, Hapus{/if}
		</button>
	{/snippet}
</Modal>

<style>
	.confirm-content {
		text-align: center;
		padding: 12px 0;
	}

	.warning-icon {
		font-size: 3rem;
		margin-bottom: 16px;
	}

	.confirm-message {
		font-size: 1.1rem;
		color: var(--text-main);
		margin-bottom: 8px;
		line-height: 1.5;
	}

	.confirm-subtext {
		font-size: 0.9rem;
		color: var(--text-muted);
	}
</style>
