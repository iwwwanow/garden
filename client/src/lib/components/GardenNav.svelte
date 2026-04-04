<script lang="ts">
  import { switchGardenTab } from '../router';

  export let active: 'garden' | 'seeds' | 'herbarium';

  const items = [
    { id: 'garden'    as const, label: 'Сад',      icon: '🌸' },
    { id: 'seeds'     as const, label: 'Семена',   icon: '🌱' },
    { id: 'herbarium' as const, label: 'Гербарий', icon: '🍂' },
  ];
</script>

<div class="nav">
  {#each items as item}
    <button
      class="item"
      class:active={active === item.id}
      on:click={() => switchGardenTab(item.id)}
    >
      {#if active === item.id}
        <span class="bg"><span class="bg-shine" /></span>
      {/if}
      <span class="icon">{item.icon}</span>
      <span class="lbl">{item.label}</span>
    </button>
  {/each}
</div>

<style>
.nav {
  display: flex;
  gap: 4px;
  padding: 5px;
  background: rgba(255,255,255,0.55);
  border-radius: 14px;
  outline: 1px rgba(255,255,255,0.4) solid;
  outline-offset: -1px;
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
}

.item {
  flex: 1;
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2px;
  padding: 8px 4px 9px;
  border-radius: 10px;
  cursor: pointer;
  transition: transform 0.1s ease;
  z-index: 1;
}
.item:active { transform: scale(0.94); }

.bg {
  position: absolute;
  inset: 0;
  background: linear-gradient(180deg, var(--blue-light) 0%, var(--blue-mid) 100%);
  box-shadow: 0px 1px 0px rgba(255,255,255,0.3) inset, 0px 2px 8px rgba(0,120,215,0.3);
  border-radius: 10px;
  overflow: hidden;
}

.bg-shine {
  position: absolute;
  top: 1px; left: 1px; right: 1px;
  height: 55%;
  background: linear-gradient(180deg, rgba(255,255,255,0.28) 0%, rgba(255,255,255,0) 100%);
  border-top-left-radius: 9px;
  border-top-right-radius: 9px;
}

.icon {
  font-size: 18px;
  line-height: 1;
  position: relative;
}

.lbl {
  position: relative;
  font-family: var(--font-body);
  font-weight: 700;
  font-size: 11px;
  color: rgba(26,26,46,0.5);
  transition: color 0.12s;
}

.item.active .lbl { color: white; }
</style>
