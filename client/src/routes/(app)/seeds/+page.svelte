<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { seeds as seedsApi, flowers as flowersApi, flowerEmoji, type Seed, type FlowerTemplate } from '$lib/api';

	let seedList = $state<Seed[]>([]);
	let templates = $state<FlowerTemplate[]>([]);
	let loading = $state(true);

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

	function templateFor(seed: Seed): FlowerTemplate | undefined {
		return templates.find((t) => t.id === seed.flower_id);
	}

	onMount(load);
</script>

<div class="page">
	<h1>Семена</h1>

	{#if loading}
		<p class="loading">загрузка...</p>
	{:else if seedList.length === 0}
		<div class="empty">
			<p>У вас пока нет семян</p>
			<p class="hint">Семена начисляются каждые 7 дней с живого цветка</p>
		</div>
	{:else}
		<div class="grid">
			{#each seedList as seed}
				<button class="seed-card" onclick={() => goto(`/seeds/${seed.id}`)}>
					<span class="type-icon">S</span>
					<span class="flower-img">{flowerEmoji(seed.flower_id)}</span>
					<div class="card-footer">
						<span>x{seed.quantity}</span>
					</div>
				</button>
			{/each}
		</div>
	{/if}
</div>

<style>
	.page {
		padding: 24px var(--page-padding);
		padding-bottom: calc(var(--nav-height) + 24px);
		max-width: 440px;
		margin: 0 auto;
		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	.grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 8px;
	}

	.seed-card {
		position: relative;
		display: flex;
		flex-direction: column;
		background: var(--color-card);
		width: 100%;
		aspect-ratio: 149 / 230;
		cursor: pointer;
		border: none;
		padding: 0;
		font-family: inherit;
	}

	.type-icon {
		position: absolute;
		top: 8px;
		left: 8px;
		font-size: 14px;
		font-weight: 600;
		background: rgba(255, 255, 255, 0.8);
		padding: 2px 5px;
	}

	.flower-img {
		flex: 1;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 64px;
	}

	.card-footer {
		display: flex;
		justify-content: space-between;
		padding: 6px 8px;
		background: rgba(0, 0, 0, 0.05);
		font-size: 14px;
		font-weight: 600;
	}

	.loading {
		opacity: 0.4;
		font-size: 16px;
	}

	.empty {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.hint {
		font-size: 14px;
		opacity: 0.5;
	}
</style>
