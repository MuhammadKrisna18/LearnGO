<script lang="ts">
	let { children, class: className = '', style = '' } = $props();
</script>

<div class={`glass-panel ${className}`} {style}>
	{@render children()}
</div>
