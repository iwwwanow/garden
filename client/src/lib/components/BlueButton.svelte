<script lang="ts">
  export let label: string;
  export let full: boolean = false;
  export let small: boolean = false;
  export let disabled: boolean = false;
  export let onClick: () => void = () => {};
</script>

<button
  class="btn"
  class:full
  class:small
  {disabled}
  on:click={onClick}
>
  <span class="shine" />
  {label}
</button>

<style>
.btn {
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 11px 24px;
  background: var(--btn-blue-bg);
  box-shadow: var(--btn-blue-shadow);
  border-radius: 8px;
  outline: var(--btn-blue-border);
  outline-offset: -1px;
  overflow: hidden;
  font-family: var(--font-body);
  font-weight: 700;
  font-size: 14px;
  color: white;
  text-shadow: 0px 1px 2px rgba(0,0,0,0.18);
  transition: transform 0.1s ease, filter 0.1s ease;
  user-select: none;
  cursor: pointer;
}

.btn.full  { width: 100%; border-radius: 10px; padding: 14px 24px; font-size: 15px; }
.btn.small { padding: 7px 16px; font-size: 12px; border-radius: 6px; }
.btn:active:not(:disabled) { transform: scale(0.96); filter: brightness(0.95); }
.btn:disabled { opacity: 0.5; cursor: not-allowed; }

.shine {
  position: absolute;
  top: 1px; left: 1px;
  right: 1px;
  height: 48%;
  background: linear-gradient(180deg, rgba(255,255,255,0.38) 0%, rgba(255,255,255,0.04) 100%);
  border-top-left-radius: 7px;
  border-top-right-radius: 7px;
  pointer-events: none;
}
</style>
