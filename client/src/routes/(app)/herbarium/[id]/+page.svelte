<script lang="ts">
  import { onMount } from "svelte";
  import { page } from "$app/stores";
  import {
    me,
    flowers as flowersApi,
    totalFd,
    formatDate,
    type UserFlower,
    type FlowerTemplate,
  } from "$lib/api";
  import ActionButton from "$lib/components/ActionButton.svelte";

  const flowerId = $derived(Number($page.params.id));

  let flower = $state<UserFlower | null>(null);
  let template = $state<FlowerTemplate | null>(null);
  let loading = $state(true);
  let isPublic = $state(true);

  const fd = $derived(flower ? totalFd(flower.day) : 0);

  async function load() {
    loading = true;
    try {
      const [data, tmpls] = await Promise.all([
        me.get(),
        flowersApi.listTemplates(),
      ]);
      flower =
        data.flowers.find((f) => f.id === flowerId && f.is_dried) ?? null;
      if (flower)
        template = tmpls.find((t) => t.id === flower!.flower_id) ?? null;
    } finally {
      loading = false;
    }
  }

  function toggleVisibility() {
    // PATCH /api/herbarium/:id/visibility — в бэклоге
    isPublic = !isPublic;
  }

  onMount(load);
</script>

{#if loading}
  <p>загрузка...</p>
{:else if !flower}
  <div class="page">
    <a href="/herbarium">← Гербарий</a>
    <p>Цветок не найден</p>
  </div>
{:else}
  <div class="page">
    <a href="/herbarium">← Гербарий</a>

    {#if template?.image_path}
      <img src="/{template.image_path}" alt="цветок" style="width:100%;" />
    {/if}
    <p>H · {fd} FD</p>

    <table>
      <tbody>
        <tr><td>Принёс всего</td><td>{fd} FD</td></tr>
        <tr><td>Дней прожил</td><td>{flower.day}</td></tr>
        <tr><td>Посажен</td><td>{formatDate(flower.created_at)}</td></tr>
        <tr
          ><td>Последний полив</td><td
            >{flower.last_watered_at
              ? formatDate(flower.last_watered_at)
              : "—"}</td
          ></tr
        >
      </tbody>
    </table>

    {#if isPublic}
      <ActionButton variant="muted" onclick={toggleVisibility}
        >E Скрыть</ActionButton
      >
    {:else}
      <ActionButton variant="primary" onclick={toggleVisibility}
        >E Выставить</ActionButton
      >
    {/if}
  </div>
{/if}
