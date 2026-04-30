<script lang="ts">
  import { onMount } from "svelte";
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";
  import {
    seeds as seedsApi,
    flowers as flowersApi,
    users as usersApi,
    ApiError,
    type Seed,
    type FlowerTemplate,
    type User,
  } from "$lib/api";
  import ActionButton from "$lib/components/ActionButton.svelte";

  const seedId = $derived(Number($page.params.id));

  let seed = $state<Seed | null>(null);
  let template = $state<FlowerTemplate | null>(null);
  let loading = $state(true);
  let planting = $state(false);
  let sharing = $state(false);
  let shareUsername = $state("");
  let shareQty = $state(1);
  let resolvedUser = $state<User | null>(null);
  let resolving = $state(false);
  let showShare = $state(false);
  let error = $state("");
  let success = $state("");

  async function load() {
    loading = true;
    try {
      const [list, tmpls] = await Promise.all([
        seedsApi.list(),
        flowersApi.listTemplates(),
      ]);
      seed = list.find((s) => s.id === seedId) ?? null;
      if (seed) template = tmpls.find((t) => t.id === seed!.flower_id) ?? null;
    } finally {
      loading = false;
    }
  }

  // TODO: to lib
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

  async function resolveUsername() {
    if (!shareUsername.trim()) {
      resolvedUser = null;
      return;
    }
    resolving = true;
    error = "";
    try {
      resolvedUser = await usersApi.getByUsername(shareUsername.trim());
    } catch {
      resolvedUser = null;
      error = "Пользователь не найден";
    } finally {
      resolving = false;
    }
  }

  async function share() {
    if (!seed || !resolvedUser) return;
    sharing = true;
    error = "";
    try {
      await seedsApi.share(resolvedUser.id, seed.flower_id, shareQty);
      success = `Передано @${resolvedUser.username}`;
      showShare = false;
      shareUsername = "";
      resolvedUser = null;
      await load();
    } catch (e: unknown) {
      error = e instanceof ApiError ? e.message : "Ошибка передачи";
    } finally {
      sharing = false;
    }
  }

  onMount(load);
</script>

{#if loading}
  <p>загрузка...</p>
{:else if !seed}
  <div class="page">
    <a href="/seeds">← Семена</a>
    <p>Семя не найдено</p>
  </div>
{:else}
  <div class="page">
    <a href="/seeds">← Семена</a>

    {#if template?.image_path}
      <img src="/{template.image_path}" alt="цветок" style="width:100%;" />
    {/if}
    <p>S · ×{seed.quantity}</p>

    {#if error}<p>{error}</p>{/if}
    {#if success}<p>{success}</p>{/if}

    <ActionButton variant="confirm" loading={planting} onclick={plant}>
      Посадить цветок
    </ActionButton>
    <ActionButton
      variant="primary"
      onclick={() => {
        showShare = !showShare;
        error = "";
      }}
    >
      Поделиться
    </ActionButton>

    {#if showShare}
      <div>
        <label for="share-to">Имя пользователя</label>
        <input
          class="input"
          id="share-to"
          type="text"
          bind:value={shareUsername}
          placeholder="@username"
          onblur={resolveUsername}
        />
        {#if resolving}<span>...</span>
        {:else if resolvedUser}<span
            >✓ {resolvedUser.first_name || resolvedUser.username}</span
          >
        {/if}
      </div>
      <div>
        <label for="share-qty">Количество (макс. {seed.quantity})</label>
        <input
          class="input"
          id="share-qty"
          type="number"
          bind:value={shareQty}
          min="1"
          max={seed.quantity}
        />
      </div>
      <ActionButton
        variant="confirm"
        loading={sharing}
        disabled={!resolvedUser}
        onclick={share}
      >
        Передать
      </ActionButton>
    {/if}
  </div>
{/if}
