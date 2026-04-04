<script lang="ts">
  import { pageStack } from './lib/router';
  import AeroBackground from './lib/components/AeroBackground.svelte';
  import TabBar from './lib/components/TabBar.svelte';
  import GardenPage from './pages/GardenPage.svelte';
  import FeedPage from './pages/FeedPage.svelte';
  import LeaderboardPage from './pages/LeaderboardPage.svelte';
  import SeedDetailPage from './pages/SeedDetailPage.svelte';
  import HerbariumDetailPage from './pages/HerbariumDetailPage.svelte';
  import FlowerDetailPage from './pages/FlowerDetailPage.svelte';

  $: current = $pageStack[$pageStack.length - 1];
  $: isTab = ['garden', 'feed', 'leaderboard'].includes(current.name);
</script>

<AeroBackground />

<div class="shell">
  <div class="content">
    {#if current.name === 'garden'}
      <GardenPage />
    {:else if current.name === 'feed'}
      <FeedPage />
    {:else if current.name === 'leaderboard'}
      <LeaderboardPage />
    {:else if current.name === 'seed-detail'}
      <SeedDetailPage id={current.id} />
    {:else if current.name === 'herbarium-detail'}
      <HerbariumDetailPage id={current.id} />
    {:else if current.name === 'flower-detail'}
      <FlowerDetailPage userId={current.userId} />
    {/if}
  </div>

  {#if isTab}
    <TabBar activeTab={current.name} />
  {/if}
</div>

<style>
  .shell {
    min-height: 100dvh;
    max-width: 430px;
    margin: 0 auto;
    display: flex;
    flex-direction: column;
    position: relative;
  }

  .content {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow-y: auto;
    overflow-x: hidden;
    padding-bottom: 80px; /* room for TabBar */
  }
</style>
