<script lang="ts">
  import { switchTab, type TabName } from '../router';

  export let activeTab: string;

  const tabs: { id: TabName; label: string; icon: string }[] = [
    { id: 'garden',      label: 'Сад',     icon: '🌸' },
    { id: 'feed',        label: 'Лента',   icon: '🌿' },
    { id: 'leaderboard', label: 'Таблица', icon: '🏆' },
  ];
</script>

<nav class="bar">
  <div class="inner">
    <div class="shine" />
    {#each tabs as tab}
      <button
        class="item"
        class:active={activeTab === tab.id}
        on:click={() => switchTab(tab.id)}
      >
        {#if activeTab === tab.id}
          <span class="item-bg">
            <span class="item-shine" />
          </span>
        {/if}
        <span class="icon">{tab.icon}</span>
        <span class="label">{tab.label}</span>
      </button>
    {/each}
  </div>
</nav>

<style>
.bar {
  position: fixed;
  bottom: 8px;
  left: 50%;
  transform: translateX(-50%);
  width: calc(100% - 16px);
  max-width: 414px;
  z-index: 100;
  padding: 0 8px;
}

.inner {
  position: relative;
  display: flex;
  align-items: center;
  background: linear-gradient(
    180deg,
    rgba(255,255,255,0.4) 0%,
    rgba(255,255,255,0.15) 100%
  );
  box-shadow:
    0px 1px 0px 1px rgba(255,255,255,0.55) inset,
    0px 4px 20px rgba(0,0,0,0.1);
  border-radius: 14px;
  outline: 1px rgba(255,255,255,0.5) solid;
  outline-offset: -1px;
  backdrop-filter: blur(14px);
  -webkit-backdrop-filter: blur(14px);
  padding: 6px;
  gap: 4px;
  overflow: hidden;
}

/* Top glass highlight strip */
.shine {
  position: absolute;
  top: 1px; left: 6px; right: 6px;
  height: 52%;
  background: linear-gradient(180deg, rgba(255,255,255,0.32) 0%, rgba(255,255,255,0) 100%);
  border-top-left-radius: 10px;
  border-top-right-radius: 10px;
  pointer-events: none;
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
  border-radius: 9px;
  cursor: pointer;
  transition: transform 0.1s ease;
  z-index: 1;
}
.item:active { transform: scale(0.93); }

/* Active background panel */
.item-bg {
  position: absolute;
  inset: 0;
  background: linear-gradient(180deg, var(--blue-light) 0%, var(--blue-mid) 100%);
  box-shadow:
    0px 1px 0px rgba(255,255,255,0.3) inset,
    0px 2px 8px rgba(0,120,215,0.35);
  border-radius: 9px;
  overflow: hidden;
}

.item-shine {
  position: absolute;
  top: 1px; left: 1px; right: 1px;
  height: 55%;
  background: linear-gradient(180deg, rgba(255,255,255,0.28) 0%, rgba(255,255,255,0) 100%);
  border-top-left-radius: 8px;
  border-top-right-radius: 8px;
}

.icon {
  font-size: 20px;
  line-height: 1;
  position: relative;
}

.label {
  font-family: var(--font-body);
  font-weight: 700;
  font-size: 11px;
  position: relative;
  transition: color 0.15s ease;
  color: rgba(26,26,46,0.55);
}

.item.active .label { color: white; }
.item.active .icon  { filter: drop-shadow(0 1px 2px rgba(0,0,0,0.2)); }
</style>
