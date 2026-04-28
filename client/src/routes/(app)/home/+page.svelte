<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { me, flowers as flowersApi, totalFd, needsWatering, type UserFlower, type FlowerTemplate, ApiError } from '$lib/api';
	import { userStore } from '$lib/stores/auth';
	import ActionButton from '$lib/components/ActionButton.svelte';
	import FlowerCard from '$lib/components/FlowerCard.svelte';

	let user = $derived($userStore);
	let userFlowers = $state<UserFlower[]>([]);
	let templates = $state<FlowerTemplate[]>([]);
	let loading = $state(true);
	let wateringId = $state<number | null>(null);
	let waterError = $state('');

	const liveFlowers = $derived(userFlowers.filter((f) => !f.is_dried));
	const driedFlowers = $derived(userFlowers.filter((f) => f.is_dried));
	const heroFlower = $derived(liveFlowers.length ? [...liveFlowers].sort((a, b) => b.day - a.day)[0] : null);
	const heroTemplate = $derived(heroFlower ? templates.find((t) => t.id === heroFlower.flower_id) : null);
	const heroFd = $derived(heroFlower ? totalFd(heroFlower.day) : 0);
	const heroNeedsWater = $derived(heroFlower ? needsWatering(heroFlower) : false);

	function imgPath(f: UserFlower): string | undefined {
		const t = templates.find((t) => t.id === f.flower_id);
		return t ? `/${t.image_path}` : undefined;
	}

	async function loadData() {
		loading = true;
		try {
			const [meData, tmpl] = await Promise.all([me.get(), flowersApi.listTemplates()]);
			$userStore && userStore.set(meData.user);
			userFlowers = meData.flowers ?? [];
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
{:else}
	<!-- Hero -->
	<section>
		<p>@{user?.username}</p>
		<h1>{user?.fd_balance ?? 0} FD</h1>
		<p>{liveFlowers.length} живых · {driedFlowers.length} гербарий</p>
		<a href="/profile">PE редактировать</a>

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

		<ActionButton variant="confirm" onclick={() => goto('/seeds')}>
			Посадить цветок
		</ActionButton>
	</section>

	<!-- Цветки -->
	<div class="page">
		<h2><a href="/seeds">Семена →</a></h2>
		{#if liveFlowers.length === 0}
			<p>начисляются каждые 7 дней с живого цветка</p>
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

		<h2><a href="/herbarium">Гербарий →</a></h2>
		{#if driedFlowers.length === 0}
			<p>пересохшие цветки появятся здесь</p>
		{:else}
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
