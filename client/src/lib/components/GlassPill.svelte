<script lang="ts">
  export let text: string;
  export let small: boolean = false;
</script>

<span class="pill" class:small>
  <span class="shine" />
  {text}
</span>

<style>
.pill {
  position: relative;
  display: inline-flex;
  align-items: center;
  padding: 6px 14px;
  background: var(--pill-bg);
  box-shadow: var(--pill-shadow);
  border-radius: 20px;
  outline: var(--pill-border);
  outline-offset: -1px;
  backdrop-filter: var(--pill-blur);
  -webkit-backdrop-filter: var(--pill-blur);
  font-family: var(--font-accent);
  font-style: italic;
  font-weight: 900;
  font-size: 14px;
  color: var(--dark);
  overflow: hidden;
  user-select: none;
}

.pill.small {
  padding: 4px 10px;
  font-size: 11px;
}

.shine {
  position: absolute;
  inset: 1px;
  height: 55%;
  background: linear-gradient(180deg, rgba(255,255,255,0.45) 0%, rgba(255,255,255,0) 100%);
  border-radius: 20px;
  pointer-events: none;
}
</style>
