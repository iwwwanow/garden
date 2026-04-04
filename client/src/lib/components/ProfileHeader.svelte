<script lang="ts">
  export let name: string;
  export let fdBalance: number;
  export let daysGrown: number;
  export let waterings: number;
</script>

<header class="header">
  <div class="name-row">
    <h1 class="name">{name}</h1>
    <slot name="action" />
  </div>
  <div class="stats">
    <span class="stat fd">
      <span class="val">{fdBalance}</span>
      <span class="unit">FD</span>
    </span>
    <span class="dot">·</span>
    <span class="stat">
      <span class="val">{daysGrown}</span>
      <span class="unit">дней роста</span>
    </span>
    <span class="dot">·</span>
    <span class="stat">
      <span class="val">{waterings}</span>
      <span class="unit">поливов</span>
    </span>
  </div>
</header>

<style>
.header {
  padding: 16px 16px 8px;
}

.name-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 6px;
}

.name {
  font-family: var(--font-display);
  font-weight: 500;
  font-size: 28px;
  color: var(--dark);
  line-height: 1.1;
}

.stats {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.stat {
  display: flex;
  align-items: baseline;
  gap: 3px;
}

.val {
  font-family: var(--font-accent);
  font-style: italic;
  font-weight: 900;
  font-size: 15px;
  color: var(--dark);
}

.stat.fd .val {
  font-size: 17px;
  color: var(--blue-mid);
}

.unit {
  font-family: var(--font-display);
  font-weight: 300;
  font-size: 12px;
  color: var(--gray-mid);
}

.dot {
  color: var(--gray-light);
  font-size: 12px;
}
</style>
