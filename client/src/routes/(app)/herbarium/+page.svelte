<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { me, totalFd, flowerEmoji, type UserFlower } from '$lib/api';
	import FlowerCard from '$lib/components/FlowerCard.svelte';

	let flowers = $state<UserFlower[]>([]);
	let loading = $state(true);

	async function load() {
		loading = true;
		try {
			const data = await me.get();
			flowers = (data.flowers ?? []).filter((f) => f.is_dried);
		} finally {
			loading = false;
		}
	}

	onMount(load);
</script>

<div class="page">
	<h1>Гербарий</h1>

	{#if loading}
		<p class="loading">загрузка...</p>
	{:else if flowers.length === 0}
		<div class="empty">
			<p>Ваш гербарий пуст</p>
			<p class="hint">Здесь будут жить воспоминания о цветках</p>
		</div>
	{:else}
		<div class="feature-flower">
			<span class="type-badge">H</span>
			<div class="feature-img">{flowerEmoji(flowers[0].flower_id)}</div>
			<div class="feature-fd">{totalFd(flowers[0].day)} FD</div>
		</div>

		<div class="grid">
			{#each flowers as flower}
				<FlowerCard
					{flower}
					type="herbarium"
					flowerTemplateId={flower.flower_id}
					onclick={() => goto(`/herbarium/${flower.id}`)}
				/>
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

	.feature-flower {
		position: relative;
		display: flex;
		flex-direction: column;
		align-items: center;
		background: var(--color-card);
		padding: 32px;
		gap: 8px;
	}

	.type-badge {
		position: absolute;
		top: 12px;
		left: 12px;
		font-size: 14px;
		font-weight: 600;
		background: rgba(255, 255, 255, 0.8);
		padding: 2px 6px;
	}

	.feature-img {
		font-size: 100px;
		line-height: 1;
	}

	.feature-fd {
		font-size: 24px;
		font-weight: 600;
	}

	.grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 8px;
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
