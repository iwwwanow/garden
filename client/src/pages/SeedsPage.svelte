<script lang="ts">
  import GardenNav from '../lib/components/GardenNav.svelte';
  import ProfileHeader from '../lib/components/ProfileHeader.svelte';
  import FlowerCard from '../lib/components/FlowerCard.svelte';
  import GreenButton from '../lib/components/GreenButton.svelte';
  import { navigate } from '../lib/router';

  const user = { name: 'John Doe', totalFd: 150, daysGrown: 120, waterings: 2548 };

  // Unplanted seeds in inventory
  const seeds = [
    { id: 1, imageUrl: '/flowers/blue.svg',   dailyFd: '+5 FD/д'  },
    { id: 2, imageUrl: '/flowers/pink.svg',   dailyFd: '+5 FD/д'  },
    { id: 3, imageUrl: '/flowers/violet.svg', dailyFd: '+10 FD/д' },
    { id: 4, imageUrl: '/flowers/amber.svg',  dailyFd: '+10 FD/д' },
  ];
</script>

<div class="page">
  <ProfileHeader
    name={user.name}
    fdBalance={user.totalFd}
    daysGrown={user.daysGrown}
    waterings={user.waterings}
  />

  <GardenNav active="seeds" />

  {#if seeds.length > 0}
    <div class="grid">
      {#each seeds as seed}
        <div class="card-wrap">
          <FlowerCard
            imageUrl={seed.imageUrl}
            fdAmount={seed.dailyFd}
            actionIcon="🌱"
            onClick={() => navigate({ name: 'seed-detail', id: seed.id })}
            onAction={() => navigate({ name: 'seed-detail', id: seed.id })}
          />
        </div>
      {/each}
    </div>
  {:else}
    <div class="empty">
      <p>Нет семян</p>
      <p class="hint">Получи семена от других игроков</p>
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

.grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}

.card-wrap { height: 264px; }

.empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 20px;
  gap: 6px;
}
.empty p {
  font-family: var(--font-display);
  font-size: 18px;
  font-weight: 400;
  color: var(--gray-mid);
  text-align: center;
}
.empty .hint { font-size: 14px; font-weight: 300; color: var(--gray-soft); }
</style>
