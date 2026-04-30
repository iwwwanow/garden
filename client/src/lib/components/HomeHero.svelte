<script lang="ts">
  import { totalFd } from "$lib/api";
  import { needsWatering } from "$lib/api";

  import type { UserFlower } from "$lib/api";

  interface Props {
    username?: string;
    fdBalance: number;
    liveFlowersQuantity: number;
    driedFlowersQuantity: number;
    flowerImagePath?: string;
    heroFlower: UserFlower | null;
  }

  let {
    username,
    fdBalance,
    liveFlowersQuantity,
    driedFlowersQuantity,
    flowerImagePath,
    heroFlower,
  }: Props = $props();

  // TODO: needs logic
  let waterError = $state("");

  const heroFlowerTotalFd = $derived(heroFlower ? totalFd(heroFlower.day) : 0);
  const heroFlowerNeedsWater = $derived(
    heroFlower ? needsWatering(heroFlower) : false,
  );
</script>

<section>
  <a href="/profile">@{username}</a>
  <p>{fdBalance} FD</p>
  <p>{liveFlowersQuantity} живых цветка</p>
  <p>{driedFlowersQuantity} гербарий</p>

  <!-- TODO: merge -->
  {#if flowerImagePath}
    <img src="/{flowerImagePath}" alt="цветок" />
  {/if}

  {#if heroFlower}
    <p>totla: {heroFlowerTotalFd} FD</p>
    <p>{heroFlower.day} день</p>
  {/if}

  {#if waterError}
    <p>{waterError}</p>
  {/if}

  {#if heroFlowerNeedsWater && heroFlower}
    <a href="#TODO:">полить</a>
  {/if}

  <!-- TODO: has seeds condition -->
  <button>посадить цветок</button>
</section>
