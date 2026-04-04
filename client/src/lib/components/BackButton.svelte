<script lang="ts">
  import { back } from '../router';
</script>

<button class="back" on:click={back}>
  <span class="shine" />
  <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
    <path d="M12.5 15L7.5 10L12.5 5" stroke="white" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"/>
  </svg>
</button>

<style>
.back {
  position: relative;
  width: 38px; height: 38px;
  background: var(--btn-blue-bg);
  box-shadow: var(--btn-blue-shadow);
  border-radius: 50%;
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
.back:active { transform: scale(0.9); }

.shine {
  position: absolute;
  top: 1px; left: 1px; right: 1px;
  height: 55%;
  background: linear-gradient(180deg, rgba(255,255,255,0.4) 0%, rgba(255,255,255,0.04) 100%);
  border-radius: 50%;
  pointer-events: none;
}
</style>
