<script lang="ts">
  import GlassCard from '../lib/components/GlassCard.svelte';
  import GlassPill from '../lib/components/GlassPill.svelte';
  import GreenButton from '../lib/components/GreenButton.svelte';
  import BlueButton from '../lib/components/BlueButton.svelte';
  import BackButton from '../lib/components/BackButton.svelte';

  export let userId: number;

  // Mock data for another user's flower
  const data = {
    ownerName: 'Marie Curie',
    imageUrl: 'https://placehold.co/400x480/c5e8f0/0078d7',
    fdEarned: 10,
    day: 3,
    plantedAt: '28.02.2026',
    totalFdAllTime: 2500,
    requestsWatering: true,
  };

  let watered = false;

  function handleWater() {
    watered = true;
  }
</script>

<div class="page">
  <div class="topbar">
    <BackButton />
    <h2 class="title">{data.ownerName}</h2>
    <div class="spacer" />
  </div>

  <GlassCard>
    <div class="flower-img">
      <img src={data.imageUrl} alt="Flower" />

      <!-- Day badge -->
      <div class="day-badge">
        <GlassPill text="День {data.day}" />
      </div>

      <!-- FD overlay bottom -->
      <div class="fd-overlay">
        {data.fdEarned}<span>FD</span>
      </div>
    </div>
  </GlassCard>

  {#if data.requestsWatering}
    <div class="watering-note">
      <span class="drop">💧</span>
      Запрашивает полив
    </div>
  {/if}

  <div class="meta">
    <div class="meta-row">
      <span class="meta-label">Посажен</span>
      <span class="meta-val">{data.plantedAt}</span>
    </div>
    <div class="meta-row">
      <span class="meta-label">Принес за всё время</span>
      <span class="meta-val accent">{data.totalFdAllTime} FD</span>
    </div>
  </div>

  <div class="actions">
    {#if watered}
      <div class="watered-msg">
        <GlassPill text="Полито ✓" />
      </div>
    {:else}
      <GreenButton label="Полить" full onClick={handleWater} disabled={!data.requestsWatering} />
    {/if}
    <BlueButton label="Запросить полив" full />
  </div>
</div>

<style>
.page {
  padding: 8px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.topbar {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 4px 0;
}
.title {
  font-family: var(--font-display);
  font-weight: 500;
  font-size: 26px;
  color: var(--dark);
  flex: 1;
}
.spacer { width: 38px; }

.flower-img {
  position: relative;
  height: 420px;
}
.flower-img img {
  width: 100%; height: 100%;
  object-fit: cover;
  display: block;
}

.day-badge {
  position: absolute;
  top: 12px;
  left: 12px;
  z-index: 4;
}

.fd-overlay {
  position: absolute;
  bottom: 14px;
  left: 14px;
  font-family: var(--font-accent);
  font-style: italic;
  font-weight: 900;
  font-size: 32px;
  color: white;
  text-shadow: 0 2px 8px rgba(0,0,0,0.35);
  display: flex;
  align-items: baseline;
  gap: 4px;
}
.fd-overlay span { font-size: 18px; }

.watering-note {
  display: flex;
  align-items: center;
  gap: 6px;
  font-family: var(--font-body);
  font-size: 13px;
  font-weight: 700;
  color: var(--blue-mid);
  padding: 2px 4px;
}
.drop { font-size: 16px; }

.meta {
  background: rgba(255,255,255,0.5);
  border-radius: 12px;
  padding: 14px 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  outline: 1px rgba(255,255,255,0.4) solid;
  outline-offset: -1px;
  backdrop-filter: blur(8px);
}

.meta-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.meta-label {
  font-family: var(--font-display);
  font-weight: 300;
  font-size: 14px;
  color: var(--gray-mid);
}

.meta-val {
  font-family: var(--font-display);
  font-weight: 500;
  font-size: 16px;
  color: var(--dark);
}
.meta-val.accent {
  font-family: var(--font-accent);
  font-style: italic;
  font-weight: 900;
  font-size: 18px;
  color: var(--blue-mid);
}

.actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding-bottom: 4px;
}

.watered-msg {
  display: flex;
  justify-content: center;
}
</style>
