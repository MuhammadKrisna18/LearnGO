<script lang="ts">
	let { children, type = 'default' } = $props();

	function getBadgeClass(t: string) {
		switch(t) {
			case 'success': return 'badge-success';
			case 'error': return 'badge-error';
			case 'warning': return 'badge-warning';
			case 'info': return 'badge-info';
			default: return 'badge-default';
		}
	}
</script>

<span class={`badge ${getBadgeClass(type)}`}>
	{@render children()}
</span>

<style>
	.badge {
		font-size: 0.75rem;
		padding: 4px 8px;
		border-radius: 9999px;
		font-weight: 600;
		letter-spacing: 0.5px;
		display: inline-block;
	}

	.badge-default {
		background: rgba(0,0,0,0.05);
		color: var(--text-muted);
	}

	.badge-success {
		background: #dcfce7;
		color: #166534;
	}

	.badge-warning {
		background: #fef9c3;
		color: #854d0e;
	}

	.badge-error {
		background: #fee2e2;
		color: #991b1b;
	}

	.badge-info {
		background: rgba(16, 185, 129, 0.1);
		color: #10b981;
		border: 1px solid rgba(16, 185, 129, 0.2);
	}
</style>
