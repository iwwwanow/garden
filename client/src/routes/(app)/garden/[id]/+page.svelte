<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { users, flowers as flowersApi, totalFd, needsWatering, ApiError, type UserFlower, type User, type FlowerTemplate } from '$lib/api';
	import { userStore } from '$lib/stores/auth';
	import ActionButton from '$lib/components/ActionButton.svelte';
	import FlowerCard from '$lib/components/FlowerCard.svelte';

	const gardenUserId = $derived(Number($page.params.id));
	const currentUser = $derived($userStore);

	let owner = $state<User | null>(null);
	let userFlowers = $state<UserFlower[]>([]);
	let templates = $state<FlowerTemplate[]>([]);
	let loading = $state(true);
	let wateringId = $state<number | null>(null);
	let waterError = $state('');

	const liveFlowers = $derived(userFlowers.filter((f) => !f.is_dried));
	const driedFlowers = $derived(userFlowers.filter((f) => f.is_dried));
	const heroFlower = $derived(liveFlowers.length ? [...liveFlowers].sort((a, b) => b.day - a.day)[0] : null);
	const heroFd = $derived(heroFlower ? totalFd(heroFlower.day) : 0);
	const heroNeedsWater = $derived(heroFlower ? needsWatering(heroFlower) : false);
	const heroTemplate = $derived(heroFlower ? templates.find((t) => t.id === heroFlower.flower_id) : null);

	function imgPath(f: UserFlower): string | undefined {
		const t = templates.find((t) => t.id === f.flower_id);
		return t ? `/${t.image_path}` : undefined;
	}

	async function loadData() {
		loading = true;
		try {
			const [ownerData, flowerData, tmpl] = await Promise.all([
				users.getById(gardenUserId),
				flowersApi.getByUser(gardenUserId),
				flowersApi.listTemplates()
			]);
			owner = ownerData;
			userFlowers = flowerData ?? [];
			templates = tmpl ?? [];
		} finally {
			loading = false;
		}
	}

	async function waterHero() {
		if (!heroFlower) return;
		wateringId = heroFlower.id;
		waterError = '';
		try {
			await flowersApi.water(heroFlower.id);
			await loadData();
		} catch (e: unknown) {
			waterError = e instanceof ApiError ? e.message : 'Ошибка полива';
		} finally {
			wateringId = null;
		}
	}

	onMount(loadData);
</script>

{#if loading}
	<p>загрузка...</p>
{:else if owner}
	<div class="page">
		<p>owner: <a href="/garden/{gardenUserId}">@{owner.username}</a></p>
		<h1>{owner.fd_balance} FD</h1>
		<p>{liveFlowers.length} живых · {driedFlowers.length} гербарий</p>

		{#if heroTemplate?.image_path}
			<img src="/{heroTemplate.image_path}" alt="цветок" style="width:100%;" />
		{/if}
		{#if heroFlower}
			<p>{heroFd} FD · День {heroFlower.day}</p>
		{/if}

		{#if waterError}<p>{waterError}</p>{/if}

		{#if heroNeedsWater && heroFlower}
			<ActionButton variant="confirm" loading={wateringId === heroFlower.id} onclick={waterHero}>
				P Полить
			</ActionButton>
		{/if}

		{#if liveFlowers.length > 0}
			<h2>Цветки</h2>
			<div class="grid">
				{#each liveFlowers as flower}
					<FlowerCard {flower} type="flower" imagePath={imgPath(flower)}
						onclick={() => goto(`/flower/${flower.id}`)} />
				{/each}
			</div>
		{/if}

		{#if driedFlowers.length > 0}
			<h2>Гербарий</h2>
			<div class="grid">
				{#each driedFlowers.slice(0, 4) as flower}
					<FlowerCard {flower} type="herbarium" imagePath={imgPath(flower)}
						onclick={() => goto(`/herbarium/${flower.id}`)} />
				{/each}
			</div>
		{/if}
	</div>
{/if}

<style>
	.grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 8px;
	}
</style>
