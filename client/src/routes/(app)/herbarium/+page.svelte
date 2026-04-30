<script lang="ts">
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import {
    me,
    flowers as flowersApi,
    totalFd,
    type UserFlower,
    type FlowerTemplate,
  } from "$lib/api";
  import FlowerCard from "$lib/components/FlowerCard.svelte";

  let flowers = $state<UserFlower[]>([]);
  let templates = $state<FlowerTemplate[]>([]);
  let loading = $state(true);

  function imgPath(f: UserFlower): string | undefined {
    const t = templates.find((t) => t.id === f.flower_id);
    return t ? `/${t.image_path}` : undefined;
  }

  async function load() {
    loading = true;
    try {
      const [data, tmpl] = await Promise.all([
        me.get(),
        flowersApi.listTemplates(),
      ]);
      flowers = (data.flowers ?? []).filter((f) => f.is_dried);
      templates = tmpl ?? [];
    } finally {
      loading = false;
    }
  }

  onMount(load);
</script>

<div class="page">
  <h1>Гербарий</h1>

  {#if loading}
    <p>загрузка...</p>
  {:else if flowers.length === 0}
    <p>пусто — здесь живут воспоминания о пересохших цветках</p>
  {:else}
    {#if flowers[0]}
      <button
        onclick={() => goto(`/herbarium/${flowers[0].id}`)}
        style="display:block; width:100%;"
      >
        {#if imgPath(flowers[0])}
          <img src={imgPath(flowers[0])} alt="цветок" style="width:100%;" />
        {/if}
        <p>H · {totalFd(flowers[0].day)} FD · День {flowers[0].day}</p>
      </button>
    {/if}

    <h2>Все цветки</h2>
    <div class="grid">
      {#each flowers as flower}
        <FlowerCard
          {flower}
          type="herbarium"
          imagePath={imgPath(flower)}
          onclick={() => goto(`/herbarium/${flower.id}`)}
        />
      {/each}
    </div>
  {/if}
</div>

<style>
  .grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 8px;
  }
</style>
