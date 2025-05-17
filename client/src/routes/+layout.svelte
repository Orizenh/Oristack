<script lang="ts">
    import ModeSwitch from '$lib/components/ModeSwitch.svelte';
	import '../app.css';
	let { children } = $props();
</script>

<ModeSwitch />
{@render children?.()}
