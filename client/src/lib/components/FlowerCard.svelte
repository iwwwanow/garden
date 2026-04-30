<script lang="ts">
  import {
    totalFd,
    needsWatering,
    type UserFlower,
    type Seed,
    flowers as flowersApi,
    ApiError,
  } from "$lib/api";

  interface Props {
    flower?: UserFlower;
    seed?: Seed;
    // TODO: enum
    type: "flower" | "seed" | "herbarium";
    imagePath?: string;
    link: string;
    quantity?: number;
  }

  let planting = $state(false);
  let error = $state("");

  let { flower, type, imagePath, link, quantity, seed }: Props = $props();

  const fd = $derived(flower && type === "flower" && totalFd(flower.day));
  const floweringDays = $derived(flower && type === "flower" && flower.day);
  const watering = $derived(
    flower && type === "flower" && needsWatering(flower),
  );

  // TODO: enum
  const typeLabel = $derived(
    type === "seed" ? "S" : type === "herbarium" ? "H" : null,
  );

  // TODO: to lib
  // TODO: refactor to fetch and a tag??
  async function plant() {
    if (!seed) return;
    planting = true;
    error = "";
    try {
      await flowersApi.plant(seed.flower_id);
      goto("/home");
    } catch (e: unknown) {
      error = e instanceof ApiError ? e.message : "Ошибка посадки";
    } finally {
      planting = false;
    }
  }
</script>

<a href={link}>
  {#if imagePath}
    <img
      src={imagePath}
      alt=""
      style="width:100%; height:100%; object-fit:cover; display:block;"
    />
  {/if}
  {#if typeLabel}<span>{type}</span>{/if}
  {#if watering}<span>полить</span>{/if}
  {#if fd}<span>{fd} FD</span>{/if}
  {#if floweringDays}<span>{floweringDays} дней</span>{/if}
  {#if quantity}<span>{quantity}</span>{/if}
</a>

<!-- TODO: enum -->
<!-- TODO: to components? -->
{#if type === "seed"}
  <button onclick={plant}> посадить </button>
{/if}
