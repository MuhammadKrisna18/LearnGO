<svelte:head>
	<title>Dashboard Mahasiswa - SIAKAD Pro</title>
</svelte:head>

<div class="page-header">
	<h1>Selamat Datang</h1>
	<p>Ini adalah halaman portal mahasiswa. Fitur-fitur akademik akan segera hadir di sini.</p>
</div>

<div class="empty-state glass-panel">
	<div class="empty-icon">
		<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
			<path d="M12 2v20"/><path d="m17 7-5-5-5 5"/><path d="m17 17-5 5-5-5"/>
		</svg>
	</div>
	<h2>Sedang Dalam Pengembangan</h2>
	<p>Halaman ini saat ini masih kosong dan sedang dalam tahap pengembangan oleh tim.</p>
</div>

<style>
	.page-header {
		margin-bottom: 2rem;
	}

	.page-header h1 {
		font-size: 1.75rem;
		font-weight: 700;
		color: var(--text-main);
		margin: 0 0 0.5rem 0;
	}

	.page-header p {
		color: var(--text-muted);
		margin: 0;
		font-size: 1.05rem;
	}

	.empty-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 5rem 2rem;
		text-align: center;
		border-radius: 16px;
		border: 1px dashed var(--border);
	}

	.empty-icon {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 80px;
		height: 80px;
		background: white;
		border-radius: 50%;
		color: var(--primary);
		margin-bottom: 1.5rem;
		box-shadow: 0 4px 12px rgba(79, 70, 229, 0.15);
	}

	.empty-state h2 {
		margin: 0 0 0.75rem 0;
		color: var(--text-main);
		font-size: 1.5rem;
	}

	.empty-state p {
		margin: 0;
		color: var(--text-muted);
		max-width: 500px;
		line-height: 1.6;
	}
</style>
