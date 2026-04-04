<script lang="ts">
  import GlassCard from '../lib/components/GlassCard.svelte';
  import GlassPill from '../lib/components/GlassPill.svelte';
  import GreenButton from '../lib/components/GreenButton.svelte';
  import BlueButton from '../lib/components/BlueButton.svelte';
  import GardenNav from '../lib/components/GardenNav.svelte';
  import ProfileHeader from '../lib/components/ProfileHeader.svelte';

  const user = {
    name: 'John Doe',
    totalFd: 150,
    daysGrown: 120,
    waterings: 2548,
  };

  const activeFlower = {
    imageUrl: '/flowers/mint.svg',
    day: 4,
    maxDay: 7,
    dailyFd: 8,
    wateredToday: false,
    isDried: false,
  };

  let watered = activeFlower.wateredToday;

  function handleWater() {
    watered = true;
  }
</script>

<div class="page">
  <ProfileHeader
    name={user.name}
    fdBalance={user.totalFd}
    daysGrown={user.daysGrown}
    waterings={user.waterings}
  />

  <GardenNav active="garden" />

  <div class="flower-card">
    <GlassCard>
      <div class="flower-img-wrap">
        <img
          src={activeFlower.imageUrl}
          alt="Твой цветок"
          class:grayscale={activeFlower.isDried}
        />

        <div class="badge top-left">
          <GlassPill text="День {activeFlower.day} / {activeFlower.maxDay}" />
        </div>

        <div class="flower-overlay">
          <GlassPill text="+{activeFlower.dailyFd} FD сегодня" />
          {#if activeFlower.isDried}
            <BlueButton label="Посадить" small />
          {:else if watered}
            <GlassPill text="Полито ✓" />
          {:else}
            <GreenButton label="Полить" small onClick={handleWater} />
          {/if}
        </div>
      </div>
    </GlassCard>
  </div>
</div>

<style>
.page {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 8px;
}

.flower-card { flex: 1; }

.flower-img-wrap {
  position: relative;
  height: 420px;
  overflow: hidden;
}

.flower-img-wrap img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}
.flower-img-wrap img.grayscale { filter: grayscale(1) opacity(0.5); }

.badge.top-left {
  position: absolute;
  top: 12px;
  left: 12px;
  z-index: 4;
}

.flower-overlay {
  position: absolute;
  bottom: 0; left: 0; right: 0;
  padding: 32px 14px 14px;
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  background: linear-gradient(0deg, rgba(0,0,0,0.22) 0%, transparent 100%);
  z-index: 4;
}
</style>
