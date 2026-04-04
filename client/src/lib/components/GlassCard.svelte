<script lang="ts">
  export let padding: string = '0';
  export let gap: string = '0';
  export let column: boolean = false;
  export let center: boolean = false;
  export let onClick: (() => void) | undefined = undefined;
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<div
  class="card"
  class:col={column}
  class:center
  class:clickable={!!onClick}
  style:padding
  style:gap
  on:click={onClick}
>
  <div class="shine" />
  <slot />
</div>

<style>
.card {
  position: relative;
  background: var(--glass-bg);
  box-shadow: var(--glass-shadow);
  border-radius: var(--glass-radius);
  outline: var(--glass-border);
  outline-offset: -1px;
  backdrop-filter: var(--glass-blur);
  -webkit-backdrop-filter: var(--glass-blur);
  overflow: hidden;
}

.card.col   { display: flex; flex-direction: column; }
.card.center { display: flex; align-items: center; justify-content: center; }
.card.clickable { cursor: pointer; transition: transform 0.12s ease, box-shadow 0.12s ease; }
.card.clickable:active { transform: scale(0.97); }

/* Top-of-card shine (dew on glass) */
.shine {
  position: absolute;
  top: 1px; left: 1px; right: 1px;
  height: 45%;
  background: linear-gradient(
    180deg,
    rgba(255,255,255,0.38) 0%,
    rgba(255,255,255,0) 100%
  );
  border-top-left-radius: calc(var(--glass-radius) - 1px);
  border-top-right-radius: calc(var(--glass-radius) - 1px);
  pointer-events: none;
  z-index: 1;
}
</style>
