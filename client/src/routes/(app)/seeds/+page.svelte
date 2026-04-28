<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { seeds as seedsApi, flowers as flowersApi, type Seed, type FlowerTemplate } from '$lib/api';

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
		<div class="grid">
			{#each seedList as seed}
				<button onclick={() => goto(`/seeds/${seed.id}`)}
					style="width:100%; aspect-ratio:190/288; overflow:hidden; display:flex; flex-direction:column;">
					<div style="flex:1; position:relative;">
						{#if imgPath(seed)}
							<img src={imgPath(seed)} alt="" style="width:100%; height:100%; object-fit:cover; display:block;" />
						{/if}
						<span>S</span>
					</div>
					<div>×{seed.quantity}</div>
				</button>
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
