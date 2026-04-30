<script lang="ts">
  import BottomNav from "$lib/components/BottomNav.svelte";

  let { children } = $props();
</script>

{@render children()}
<BottomNav />
