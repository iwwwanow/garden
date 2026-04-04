<script lang="ts">
  import GlassCard from '../lib/components/GlassCard.svelte';
  import GlassPill from '../lib/components/GlassPill.svelte';
  import FlowerCard from '../lib/components/FlowerCard.svelte';
  import GreenButton from '../lib/components/GreenButton.svelte';
  import BlueButton from '../lib/components/BlueButton.svelte';
  import SubTabs from '../lib/components/SubTabs.svelte';
  import ProfileHeader from '../lib/components/ProfileHeader.svelte';
  import { navigate } from '../lib/router';

  // Mock data — replace with real API
  const user = {
    name: 'John Doe',
    fdBalance: 150,
    daysGrown: 120,
    waterings: 2548,
  };

  const activeFlower = {
    imageUrl: '/assets/flowers/season-1.png',
    day: 4,
    fdEarned: 10,
    wateredToday: false,
    isDried: false,
  };

  const seeds = [
    { id: 1, imageUrl: 'https://placehold.co/189x284/a8d8ea/4da6ff', fdAmount: '5FD' },
    { id: 2, imageUrl: 'https://placehold.co/189x284/c5e8f0/0078d7', fdAmount: '10FD' },
    { id: 3, imageUrl: 'https://placehold.co/189x284/d4f1e8/2ecc71', fdAmount: '5FD' },
    { id: 4, imageUrl: 'https://placehold.co/189x284/b8dff0/40f4c1', fdAmount: '10FD' },
  ];

  const herbarium = [
    { id: 1, imageUrl: 'https://placehold.co/189x284/adadad/4f4f4f', fdAmount: '32FD', isDried: true },
    { id: 2, imageUrl: 'https://placehold.co/189x284/d0d0d0/4f4f4f', fdAmount: '64FD', isDried: true },
  ];

  let subTab = 'seeds';
  const subTabs = [
    { id: 'seeds',     label: 'Семена' },
    { id: 'herbarium', label: 'Гербарий' },
  ];

  function handleWater() {
    // TODO: API call
    console.log('Water flower');
  }
</script>

<div class="page">
  <!-- Profile header -->
  <ProfileHeader
    name={user.name}
    fdBalance={user.fdBalance}
    daysGrown={user.daysGrown}
    waterings={user.waterings}
  />

  <!-- Main flower card -->
  <div class="main-card-wrap">
    <GlassCard column>
      <div class="flower-img-wrap">
        <img
          src={activeFlower.imageUrl}
          alt="Your flower"
          class:grayscale={activeFlower.isDried}
        />

        <!-- Day counter top-left -->
        <div class="day-badge">
          <GlassPill text="День {activeFlower.day}" />
        </div>
      </div>

      <!-- Bottom bar: FD + action -->
      <div class="flower-bar">
        <span class="fd-amount">
          {activeFlower.fdEarned}<span class="fd-unit">FD</span>
        </span>
        <div class="flower-actions">
          {#if !activeFlower.wateredToday && !activeFlower.isDried}
            <GreenButton label="Полить" onClick={handleWater} />
          {:else if activeFlower.isDried}
            <BlueButton label="Посадить цветок" />
          {:else}
            <GlassPill text="Полито ✓" />
          {/if}
        </div>
      </div>
    </GlassCard>
  </div>

  <!-- Sub-tabs -->
  <div class="section">
    <SubTabs tabs={subTabs} active={subTab} onChange={(id) => (subTab = id)} />
  </div>

  <!-- Seeds / Herbarium grid -->
  {#if subTab === 'seeds'}
    <div class="grid">
      {#each seeds as seed}
        <div class="card-wrap">
          <FlowerCard
            imageUrl={seed.imageUrl}
            fdAmount={seed.fdAmount}
            actionIcon="🤝"
            onClick={() => navigate({ name: 'seed-detail', id: seed.id })}
          />
        </div>
      {/each}
    </div>

    <div class="plant-btn">
      <GreenButton label="Посадить цветок" full />
    </div>

  {:else}
    <div class="grid">
      {#each herbarium as h}
        <div class="card-wrap">
          <FlowerCard
            imageUrl={h.imageUrl}
            fdAmount={h.fdAmount}
            actionIcon="📤"
            isDried={h.isDried}
            onClick={() => navigate({ name: 'herbarium-detail', id: h.id })}
          />
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
.page {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 8px;
}

/* Active flower */
.main-card-wrap { padding: 0; }

.flower-img-wrap {
  position: relative;
  height: 360px;
  overflow: hidden;
}

.flower-img-wrap img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}
.flower-img-wrap img.grayscale { filter: grayscale(1) opacity(0.55); }

.day-badge {
  position: absolute;
  top: 12px;
  left: 12px;
  z-index: 4;
}

.flower-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 14px 14px;
}

.fd-amount {
  font-family: var(--font-accent);
  font-style: italic;
  font-weight: 900;
  font-size: 26px;
  color: var(--dark);
  display: flex;
  align-items: baseline;
  gap: 3px;
}
.fd-unit {
  font-size: 14px;
  font-weight: 700;
  color: var(--blue-mid);
}

.flower-actions { display: flex; gap: 8px; align-items: center; }

/* Section spacing */
.section { padding: 0 4px; }

/* Card grid 2 columns */
.grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}

.card-wrap {
  height: 260px;
}
.card-wrap > :global(.card) {
  height: 100%;
}

/* Plant button */
.plant-btn { padding: 4px 0 8px; }
</style>
