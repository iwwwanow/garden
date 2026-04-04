<script lang="ts">
  export let tabs: { id: string; label: string }[];
  export let active: string;
  export let onChange: (id: string) => void;
</script>

<div class="strip">
  {#each tabs as tab}
    <button
      class="tab"
      class:active={active === tab.id}
      on:click={() => onChange(tab.id)}
    >
      {#if active === tab.id}
        <span class="tab-bg">
          <span class="tab-shine" />
        </span>
      {/if}
      <span class="lbl">{tab.label}</span>
    </button>
  {/each}
</div>

<style>
.strip {
  display: flex;
  gap: 4px;
  padding: 4px;
  background: rgba(255,255,255,0.55);
  border-radius: 12px;
  outline: 1px rgba(255,255,255,0.35) solid;
  outline-offset: -1px;
}

.tab {
  flex: 1;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 8px 12px;
  border-radius: 9px;
  cursor: pointer;
  transition: transform 0.1s ease;
}
.tab:active { transform: scale(0.96); }

.tab-bg {
  position: absolute;
  inset: 0;
  background: linear-gradient(180deg, var(--blue-light) 0%, var(--blue-mid) 100%);
  box-shadow:
    0px 1px 0px rgba(255,255,255,0.3) inset,
    0px 2px 6px rgba(0,120,215,0.3);
  border-radius: 9px;
  overflow: hidden;
}

.tab-shine {
  position: absolute;
  top: 1px; left: 1px; right: 1px;
  height: 55%;
  background: linear-gradient(180deg, rgba(255,255,255,0.28) 0%, rgba(255,255,255,0) 100%);
  border-top-left-radius: 8px;
  border-top-right-radius: 8px;
}

.lbl {
  position: relative;
  font-family: var(--font-body);
  font-weight: 700;
  font-size: 13px;
  color: rgba(26,26,46,0.55);
}

.tab.active .lbl { color: white; }
</style>
