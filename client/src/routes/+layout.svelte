<script lang="ts">
  import "../app.css";
  import { dev } from "$app/environment";
  import { onMount } from "svelte";
  import { initAuth } from "$lib/stores/auth";
  import DevPanel from "$lib/components/DevPanel.svelte";

  let { children } = $props();

  onMount(() => {
    initAuth();
  });
</script>

{@render children()}
{#if dev}
  <DevPanel />
{/if}
