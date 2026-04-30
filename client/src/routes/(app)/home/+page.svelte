<script lang="ts">
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import {
    me,
    flowers as flowersApi,
    totalFd,
    needsWatering,
    type UserFlower,
    type FlowerTemplate,
    ApiError,
  } from "$lib/api";
  import { userStore } from "$lib/stores/auth";
  import ActionButton from "$lib/components/ActionButton.svelte";
  import HomeHero from "$lib/components/HomeHero.svelte";
  import FlowerCard from "$lib/components/FlowerCard.svelte";

  let user = $derived($userStore);
  let userFlowers = $state<UserFlower[]>([]);
  let templates = $state<FlowerTemplate[]>([]);
  let loading = $state(true);
  let wateringId = $state<number | null>(null);

  const liveFlowers = $derived(userFlowers.filter((f) => !f.is_dried));
  const driedFlowers = $derived(userFlowers.filter((f) => f.is_dried));
  const heroFlower = $derived(
    liveFlowers.length
      ? [...liveFlowers].sort((a, b) => b.day - a.day)[0]
      : null,
  );

  const heroFlowerTemplate = $derived(
    heroFlower ? templates.find((t) => t.id === heroFlower.flower_id) : null,
  );

  function imgPath(f: UserFlower): string | undefined {
    const t = templates.find((t) => t.id === f.flower_id);
    return t ? `/${t.image_path}` : undefined;
  }

  // TODO: почему не используем ssr?
  async function loadData() {
    loading = true;
    try {
      const [meData, tmpl] = await Promise.all([
        me.get(),
        flowersApi.listTemplates(),
      ]);
      $userStore && userStore.set(meData.user);
      userFlowers = meData.flowers ?? [];
      templates = tmpl ?? [];
    } finally {
      loading = false;
    }
  }

  onMount(loadData);
</script>

{#if loading}
  <p>загрузка...</p>
{:else}
  <HomeHero
    username={user?.username}
    fdBalance={user?.fd_balance ?? 0}
    liveFlowersQuantity={liveFlowers.length}
    driedFlowersQuantity={driedFlowers.length}
    flowerImagePath={heroFlowerTemplate?.image_path}
    {heroFlower}
  />

  <section>
    <a href="/seeds">семена</a>
    {#if liveFlowers.length === 0}
      <p>семена начисляются каждые 7 дней с живого цветка</p>
    {/if}
  </section>

  <section>
    {#if liveFlowers.length > 0}
      <p>цветки</p>
      {#each liveFlowers as flower}
        <FlowerCard
          {flower}
          type="flower"
          imagePath={imgPath(flower)}
          link={`/flower/${flower.id}`}
        />
      {/each}
    {/if}
  </section>

  <section>
    <a href="/herbarium">гербарий</a>
    {#if driedFlowers.length === 0}
      <p>пересохшие цветки появятся здесь</p>
    {:else}
      {#each driedFlowers.slice(0, 4) as flower}
        <FlowerCard
          {flower}
          type="herbarium"
          imagePath={imgPath(flower)}
          link={`/herbarium/${flower.id}`}
        />
      {/each}
    {/if}
  </section>
{/if}

<style>
  section {
    max-width: 600px;
  }
</style>
