<script lang="ts">
  import GlassPill from './GlassPill.svelte';

  export let imageUrl: string = 'https://placehold.co/189x284/a8d8ea/ffffff';
  export let fdAmount: string = '5FD';
  export let actionIcon: string = '🤝'; // share/water icon
  export let isDried: boolean = false;
  export let onClick: () => void = () => {};
  export let onAction: () => void = () => {};
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<div class="card" class:dried={isDried} on:click={onClick}>
  <div class="shine" />
  <img src={imageUrl} alt="Flower" class:grayscale={isDried} />

  <div class="overlay">
    <GlassPill text={fdAmount} />
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <button
      class="action-btn"
      on:click|stopPropagation={onAction}
    >
      <span class="btn-shine" />
      <span class="icon">{actionIcon}</span>
    </button>
  </div>
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
  cursor: pointer;
  transition: transform 0.12s ease;
}
.card:active { transform: scale(0.97); }

.shine {
  position: absolute;
  top: 1px; left: 1px; right: 1px;
  height: 42%;
  background: linear-gradient(180deg, rgba(255,255,255,0.38) 0%, rgba(255,255,255,0) 100%);
  border-top-left-radius: 15px;
  border-top-right-radius: 15px;
  pointer-events: none;
  z-index: 2;
}

img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}
img.grayscale { filter: grayscale(1) opacity(0.6); }

/* Bottom overlay: FD + action */
.overlay {
  position: absolute;
  bottom: 0; left: 0; right: 0;
  padding: 10px 10px 12px;
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  background: linear-gradient(0deg, rgba(0,0,0,0.18) 0%, transparent 100%);
  z-index: 3;
}

/* Action icon button (blue circle/rect) */
.action-btn {
  position: relative;
  width: 40px; height: 40px;
  background: var(--btn-blue-bg);
  box-shadow: var(--btn-blue-shadow);
  border-radius: 20px;
  outline: var(--btn-blue-border);
  outline-offset: -1px;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: transform 0.1s ease;
  flex-shrink: 0;
}
.action-btn:active { transform: scale(0.9); }

.btn-shine {
  position: absolute;
  top: 1px; left: 1px; right: 1px;
  height: 50%;
  background: linear-gradient(180deg, rgba(255,255,255,0.4) 0%, rgba(255,255,255,0.04) 100%);
  border-top-left-radius: 20px;
  border-top-right-radius: 20px;
  pointer-events: none;
}

.icon { font-size: 18px; line-height: 1; }
</style>
