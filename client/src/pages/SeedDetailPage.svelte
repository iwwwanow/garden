<script lang="ts">
  import GlassCard from '../lib/components/GlassCard.svelte';
  import BlueButton from '../lib/components/BlueButton.svelte';
  import GreenButton from '../lib/components/GreenButton.svelte';
  import BackButton from '../lib/components/BackButton.svelte';

  export let id: number;

  // Mock seed data
  const seed = {
    imageUrl: `https://placehold.co/400x480/a8d8ea/4da6ff`,
    fdAmount: '10FD',
    acquiredAt: '28.02.2026',
    fromUser: '@iwwwanow',
  };
</script>

<div class="page">
  <div class="topbar">
    <BackButton />
    <h2 class="title">Семя</h2>
    <div class="spacer" />
  </div>

  <div class="card-wrap">
    <GlassCard>
      <div class="flower-img">
        <img src={seed.imageUrl} alt="Seed flower" />
        <div class="fd-overlay">{seed.fdAmount}</div>
      </div>
    </GlassCard>
  </div>

  <div class="meta">
    <div class="meta-row">
      <span class="meta-label">Приобретено</span>
      <span class="meta-val">{seed.acquiredAt}</span>
    </div>
    <div class="meta-row">
      <span class="meta-label">От</span>
      <span class="meta-val accent">{seed.fromUser}</span>
    </div>
  </div>

  <div class="actions">
    <GreenButton label="Посадить" full />
    <BlueButton label="Поделиться" full />
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

.card-wrap { flex: 1; }

.flower-img {
  position: relative;
  height: 400px;
}
.flower-img img {
  width: 100%; height: 100%;
  object-fit: cover;
  display: block;
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
}

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
.meta-val.accent { color: var(--blue-mid); }

.actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>
