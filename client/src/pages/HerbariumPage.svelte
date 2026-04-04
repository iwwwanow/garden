<script lang="ts">
  import GardenNav from '../lib/components/GardenNav.svelte';
  import ProfileHeader from '../lib/components/ProfileHeader.svelte';
  import FlowerCard from '../lib/components/FlowerCard.svelte';
  import { navigate } from '../lib/router';

  const user = { name: 'John Doe', totalFd: 150, daysGrown: 120, waterings: 2548 };

  // Dried flower collection
  const herbarium = [
    { id: 1, imageUrl: '/flowers/amber.svg',  totalFd: '32 FD', isDried: true },
    { id: 2, imageUrl: '/flowers/pink.svg',   totalFd: '64 FD', isDried: true },
    { id: 3, imageUrl: '/flowers/violet.svg', totalFd: '15 FD', isDried: true },
  ];
</script>

<div class="page">
  <ProfileHeader
    name={user.name}
    fdBalance={user.totalFd}
    daysGrown={user.daysGrown}
    waterings={user.waterings}
  />

  <GardenNav active="herbarium" />

  {#if herbarium.length > 0}
    <div class="grid">
      {#each herbarium as h}
        <div class="card-wrap">
          <FlowerCard
            imageUrl={h.imageUrl}
            fdAmount={h.totalFd}
            actionIcon="📤"
            isDried={h.isDried}
            onClick={() => navigate({ name: 'herbarium-detail', id: h.id })}
          />
        </div>
      {/each}
    </div>
  {:else}
    <div class="empty">
      <p>Гербарий пуст</p>
      <p class="hint">Засохшие цветы попадут сюда</p>
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
