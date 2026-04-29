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

	onMount(loadData);
</script>

{#if loading}
	<p>загрузка...</p>
{:else}
	<!-- Hero -->
	<section>
		<a href="/profile">@{user?.username}</a>
		<p>{user?.fd_balance ?? 0} FD</p>
		<p>{liveFlowers.length} живых цветка</p>
		<p>{driedFlowers.length} гербарий</p>

		<!-- TODO: merge -->
		{#if heroTemplate?.image_path}
			<img src="/{heroTemplate.image_path}" alt="цветок" />
		{/if}
		{#if heroFlower}
			<p>{heroFd} FD</p>
			<p>{heroFlower.day} день</p>
		{/if}

		{#if waterError}<p>{waterError}</p>{/if}

		{#if heroNeedsWater && heroFlower}
			<button>полить</button>
		{/if}

		<!-- TODO: has seeds condition -->
		<button>посадить цветок</button>
	</section>

	<a href="/seeds">семена</a>

	<div class="page">
		<!-- TODO: for what? -->
		{#if liveFlowers.length === 0}
			<p>начисляются каждые 7 дней с живого цветка</p>
		{/if}

		{#if liveFlowers.length > 0}
			<p>цветки</p>
			{#each liveFlowers as flower}
				<FlowerCard {flower} type="flower" imagePath={imgPath(flower)} link={`/flower/${flower.id}`}/>
			{/each}
		{/if}

		<a href="/herbarium">гербарий</a>
		{#if driedFlowers.length === 0}
			<p>пересохшие цветки появятся здесь</p>
		{:else}
			{#each driedFlowers.slice(0, 4) as flower}
				<FlowerCard {flower} type="herbarium" imagePath={imgPath(flower)} link={`/herbarium/${flower.id}`} />
			{/each}
		{/if}
	</div>
{/if}
