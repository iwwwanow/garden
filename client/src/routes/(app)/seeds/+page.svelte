<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { seeds as seedsApi, flowers as flowersApi, type Seed, type FlowerTemplate } from '$lib/api';
	import FlowerCard from '$lib/components/FlowerCard.svelte';

	let seedList = $state<Seed[]>([]);
	let templates = $state<FlowerTemplate[]>([]);
	let loading = $state(true);

	function imgPath(seed: Seed): string | undefined {
		const t = templates.find((t) => t.id === seed.flower_id);
		return t ? `/${t.image_path}` : undefined;
	}

	async function load() {
		loading = true;
		try {
			const [s, t] = await Promise.all([seedsApi.list(), flowersApi.listTemplates()]);
			seedList = (s ?? []).filter((s) => s.quantity > 0);
			templates = t ?? [];
		} finally {
			loading = false;
		}
	}

	onMount(load);
</script>

<div class="page">
	<h1>Семена</h1>

	{#if loading}
		<p>загрузка...</p>
	{:else if seedList.length === 0}
		<p>нет семян — начисляются каждые 7 дней с живого цветка</p>
	{:else}
			{#each seedList as seed}
				<FlowerCard flower={seed} type="seed" imagePath={imgPath(seed)} link={`/seeds/${seed.id}`} quantity={seed.quantity}/>
			{/each}
	{/if}
</div>
