<script lang="ts">
  import GlassCard from '../lib/components/GlassCard.svelte';
  import BlueButton from '../lib/components/BlueButton.svelte';
  import BackButton from '../lib/components/BackButton.svelte';

  export let id: number;

  // Mock dried flower data
  const flower = {
    imageUrl: `/flowers/amber.svg`,
    totalFd: 2500,
    plantedAt: '28.02.2026',
    diedAt: '29.03.2026',
    season: 1,
  };
</script>

<div class="page">
  <div class="topbar">
    <BackButton />
    <h2 class="title">Гербарий</h2>
    <div class="spacer" />
  </div>

  <div class="card-wrap">
    <GlassCard>
      <div class="flower-img">
        <img src={flower.imageUrl} alt="Dried flower" />
        <div class="fd-overlay">{flower.totalFd} <span>FD</span></div>
      </div>
    </GlassCard>
  </div>

  <div class="meta">
    <div class="meta-row">
      <span class="meta-label">Принес за всё время</span>
      <span class="meta-val accent">{flower.totalFd} FD</span>
    </div>
    <div class="divider" />
    <div class="meta-row">
      <span class="meta-label">Посажен</span>
      <span class="meta-val">{flower.plantedAt}</span>
    </div>
    <div class="meta-row">
      <span class="meta-label">Засох</span>
      <span class="meta-val">{flower.diedAt}</span>
    </div>
  </div>

  <div class="actions">
    <BlueButton label="Выставить на продажу" full />
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
  height: 400px;
}
.flower-img img {
  width: 100%; height: 100%;
  object-fit: cover;
  display: block;
  filter: grayscale(0.4) sepia(0.15);
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
  text-shadow: 0 2px 8px rgba(0,0,0,0.4);
}
.fd-overlay span { font-size: 20px; }

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

.divider {
  height: 1px;
  background: rgba(0,0,0,0.06);
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

.actions { padding-bottom: 4px; }
</style>
