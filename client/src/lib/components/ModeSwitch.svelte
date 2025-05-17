<script lang="ts">
    import { ModeWatcher,toggleMode } from 'mode-watcher';
    import Button from './ui/button/button.svelte';
</script>

<ModeWatcher />
<Button onclick={toggleMode}>
  <span class="sr-only">Toggle theme</span>
</Button>