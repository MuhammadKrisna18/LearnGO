<script lang="ts">
	import { onMount } from 'svelte';

	let visible = $state(false);

	onMount(() => {
		setTimeout(() => {
			visible = true;
		}, 100);
	});
</script>

<div class="splash-container" class:visible>
	<div class="icon-container">
		<span class="dosen-icon">🎓</span>
	</div>
	<h2>Portal Dosen</h2>
	<p>Menyiapkan ruang kerja Anda...</p>
	<div class="progress-bar">
		<div class="progress-fill"></div>
	</div>
</div>

<style>
	.splash-container {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		opacity: 0;
		transform: scale(0.95);
		transition: all 0.6s cubic-bezier(0.16, 1, 0.3, 1);
	}

	.splash-container.visible {
		opacity: 1;
		transform: scale(1);
	}

	.icon-container {
		width: 80px;
		height: 80px;
		background: linear-gradient(135deg, rgba(37, 99, 235, 0.1), rgba(56, 189, 248, 0.1));
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		margin-bottom: 24px;
		box-shadow: 0 8px 32px rgba(37, 99, 235, 0.15);
		border: 1px solid rgba(37, 99, 235, 0.2);
		animation: float 3s ease-in-out infinite;
	}

	.dosen-icon {
		font-size: 40px;
	}

	@keyframes float {
		0% { transform: translateY(0px); }
		50% { transform: translateY(-10px); }
		100% { transform: translateY(0px); }
	}

	h2 {
		font-size: 1.5rem;
		font-weight: 700;
		color: var(--text-main);
		margin-bottom: 8px;
		background: linear-gradient(135deg, var(--primary-color), #3b82f6);
		-webkit-background-clip: text;
		-webkit-text-fill-color: transparent;
	}

	p {
		color: var(--text-muted);
		font-size: 0.95rem;
		margin-bottom: 32px;
	}

	.progress-bar {
		width: 200px;
		height: 6px;
		background: rgba(37, 99, 235, 0.1);
		border-radius: 10px;
		overflow: hidden;
	}

	.progress-fill {
		width: 0%;
		height: 100%;
		background: linear-gradient(90deg, var(--primary-color), #3b82f6);
		border-radius: 10px;
		animation: load 2.5s cubic-bezier(0.4, 0, 0.2, 1) forwards;
	}

	@keyframes load {
		0% { width: 0%; }
		50% { width: 70%; }
		100% { width: 100%; }
	}
</style>
