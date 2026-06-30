<script lang="ts">
	let { children, class: className = '', style = '' } = $props();
</script>

<div class={`card ${className}`} {style}>
	{@render children()}
</div>
