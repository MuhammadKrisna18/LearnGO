<script>
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { authState } from '$lib/stores/auth.svelte';

	let visible = true;

	onMount(() => {
		
		setTimeout(() => {
			visible = false;
			setTimeout(() => {
				const user = authState.user;
				if (user) {
					if (user.role === 'mahasiswa') {
						goto('/mahasiswa');
					} else {
						goto('/dashboard');
					}
				} else {
					goto('/login');
				}
			}, 500); 
		}, 2500);
	});
</script>

<svelte:head>
	<title>SIAKAD Pro</title>
</svelte:head>

<div class="splash-container" class:fade-out={!visible}>
	<div class="logo-wrapper">
		<div class="logo-icon">
			<svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
				<path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1 0-5H20"/>
			</svg>
		</div>
		<h1 class="logo-text">SIAKAD Pro</h1>
		<div class="loading-bar">
			<div class="loading-progress"></div>
		</div>
	</div>
</div>

<style>
	.splash-container {
		display: flex;
		justify-content: center;
		align-items: center;
		height: 100vh;
		background: linear-gradient(135deg, var(--primary) 0%, #312e81 100%);
		transition: opacity 0.5s ease-out;
		opacity: 1;
	}

	.fade-out {
		opacity: 0;
	}

	.logo-wrapper {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 1.5rem;
	}

	.logo-icon {
		color: white;
		animation: float 3s ease-in-out infinite;
	}

	.logo-text {
		color: white;
		font-family: 'Outfit', sans-serif;
		font-size: 3rem;
		font-weight: 700;
		margin: 0;
		letter-spacing: -0.02em;
	}

	.loading-bar {
		width: 200px;
		height: 4px;
		background: rgba(255, 255, 255, 0.2);
		border-radius: 4px;
		overflow: hidden;
		margin-top: 1rem;
	}

	.loading-progress {
		width: 0%;
		height: 100%;
		background: white;
		animation: load 2s cubic-bezier(0.4, 0, 0.2, 1) forwards;
	}

	@keyframes float {
		0%, 100% { transform: translateY(0); }
		50% { transform: translateY(-10px); }
	}

	@keyframes load {
		0% { width: 0%; }
		50% { width: 70%; }
		100% { width: 100%; }
	}
</style>
