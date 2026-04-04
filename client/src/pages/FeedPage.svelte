<script lang="ts">
  import GlassCard from '../lib/components/GlassCard.svelte';
  import GlassPill from '../lib/components/GlassPill.svelte';
  import BlueButton from '../lib/components/BlueButton.svelte';
  import { navigate } from '../lib/router';

  // Mock feed entries
  const entries = [
    {
      userId: 1,
      name: 'Marie Curie',
      fdBalance: 420,
      daysGrown: 87,
      waterings: 1240,
      requestsWatering: true,
      flowers: [
        { id: 1, imageUrl: 'https://placehold.co/84x84/a8d8ea/4da6ff', fdAmount: '10FD' },
        { id: 2, imageUrl: 'https://placehold.co/84x84/c5e8f0/0078d7', fdAmount: '5FD' },
      ],
    },
    {
      userId: 2,
      name: 'Nikola Tesla',
      fdBalance: 256,
      daysGrown: 42,
      waterings: 630,
      requestsWatering: false,
      flowers: [
        { id: 3, imageUrl: 'https://placehold.co/84x84/d4f1e8/2ecc71', fdAmount: '5FD' },
        { id: 4, imageUrl: 'https://placehold.co/84x84/b8dff0/40f4c1', fdAmount: '10FD' },
      ],
    },
    {
      userId: 3,
      name: 'Ada Lovelace',
      fdBalance: 711,
      daysGrown: 155,
      waterings: 3102,
      requestsWatering: true,
      flowers: [
        { id: 5, imageUrl: 'https://placehold.co/84x84/ffe4b5/d4a017', fdAmount: '32FD' },
        { id: 6, imageUrl: 'https://placehold.co/84x84/f0d4e8/a020f0', fdAmount: '8FD' },
      ],
    },
  ];
</script>

<div class="page">
  <h2 class="title">Лента</h2>

  {#each entries as e}
    <GlassCard column padding="12px" gap="10px">
      <!-- User row -->
      <div class="user-row">
        <!-- Avatar placeholder -->
        <div class="avatar" role="button" tabindex="0"
          on:click={() => navigate({ name: 'flower-detail', userId: e.userId })}
          on:keydown={() => {}}
        >
          <span class="avatar-initial">{e.name[0]}</span>
        </div>

        <div class="user-info">
          <span class="user-name">{e.name}</span>
          <div class="user-stats">
            <span>{e.fdBalance} FD</span>
            <span class="sep">·</span>
            <span>{e.daysGrown} дней</span>
            <span class="sep">·</span>
            <span>{e.waterings} поливов</span>
          </div>
          {#if e.requestsWatering}
            <span class="request-badge">Запрашивает полив</span>
          {/if}
        </div>
      </div>

      <!-- Flower mini-grid + water button -->
      <div class="flowers-row">
        <div class="mini-grid">
          {#each e.flowers as f}
            <div class="mini-card">
              <div class="mini-shine" />
              <img src={f.imageUrl} alt="flower" />
              <span class="mini-fd">{f.fdAmount}</span>
            </div>
          {/each}
        </div>

        {#if e.requestsWatering}
          <BlueButton label="Полить" onClick={() => navigate({ name: 'flower-detail', userId: e.userId })} />
        {/if}
      </div>
    </GlassCard>
  {/each}
</div>

<style>
.page {
  padding: 8px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.title {
  font-family: var(--font-display);
  font-weight: 500;
  font-size: 32px;
  padding: 8px 4px 0;
  color: var(--dark);
}

/* User row */
.user-row {
  display: flex;
  align-items: flex-start;
  gap: 10px;
}

.avatar {
  width: 44px; height: 44px;
  border-radius: 50%;
  background: linear-gradient(180deg, var(--blue-light), var(--blue-mid));
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(0,120,215,0.25);
  transition: transform 0.1s;
}
.avatar:active { transform: scale(0.93); }
.avatar-initial {
  font-family: var(--font-display);
  font-size: 22px;
  font-weight: 500;
  color: white;
  line-height: 1;
}

.user-info { flex: 1; display: flex; flex-direction: column; gap: 2px; }

.user-name {
  font-family: var(--font-display);
  font-weight: 500;
  font-size: 18px;
  color: var(--dark);
}

.user-stats {
  display: flex;
  gap: 4px;
  font-family: var(--font-body);
  font-size: 11px;
  color: var(--gray-mid);
}

.sep { color: var(--gray-light); }

.request-badge {
  display: inline-block;
  font-size: 11px;
  font-weight: 700;
  color: var(--blue-mid);
  font-family: var(--font-body);
}

/* Flowers + water button row */
.flowers-row {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 10px;
}

.mini-grid {
  display: flex;
  gap: 6px;
}

.mini-card {
  position: relative;
  width: 80px; height: 80px;
  border-radius: 10px;
  overflow: hidden;
  background: var(--glass-bg);
  outline: 1px rgba(255,255,255,0.4) solid;
  outline-offset: -1px;
}

.mini-shine {
  position: absolute;
  top: 1px; left: 1px; right: 1px;
  height: 45%;
  background: linear-gradient(180deg, rgba(255,255,255,0.35) 0%, transparent 100%);
  border-top-left-radius: 9px;
  border-top-right-radius: 9px;
  pointer-events: none;
  z-index: 2;
}

.mini-card img {
  width: 100%; height: 100%;
  object-fit: cover;
}

.mini-fd {
  position: absolute;
  bottom: 4px;
  left: 5px;
  font-family: var(--font-accent);
  font-style: italic;
  font-weight: 900;
  font-size: 11px;
  color: white;
  text-shadow: 0 1px 3px rgba(0,0,0,0.5);
  z-index: 3;
}
</style>
