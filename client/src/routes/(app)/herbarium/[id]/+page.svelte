<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { me, totalFd, flowerEmoji, formatDate, type UserFlower } from '$lib/api';
	import ActionButton from '$lib/components/ActionButton.svelte';

	const flowerId = $derived(Number($page.params.id));

	let flower = $state<UserFlower | null>(null);
	let loading = $state(true);

	// Herbarium visibility is a backlog feature — toggle is UI-only for now
	let isPublic = $state(true);

	async function load() {
		loading = true;
		try {
			const data = await me.get();
			flower = data.flowers.find((f) => f.id === flowerId && f.is_dried) ?? null;
		} finally {
			loading = false;
		}
	}

	const fd = $derived(flower ? totalFd(flower.day) : 0);
	const emoji = $derived(flower ? flowerEmoji(flower.flower_id) : '');

	function toggleVisibility() {
		// PATCH /api/herbarium/:id/visibility — в бэклоге
		isPublic = !isPublic;
	}

	onMount(load);
</script>

{#if loading}
	<div class="loading">загрузка...</div>
{:else if !flower}
	<div class="page">
		<p>Цветок не найден</p>
		<a href="/herbarium">← назад</a>
	</div>
{:else}
	<div class="page">
		<a href="/herbarium" class="back">← Гербарий</a>

		<div class="flower-wrap">
			<span class="type-badge">H</span>
			<div class="flower-img">{emoji}</div>
		</div>

		<div class="stats">
			<div class="stat">
				<span class="stat-label">Принес за все время</span>
				<span class="stat-value">{fd} FD</span>
			</div>
			<div class="stat">
				<span class="stat-label">Посажен</span>
				<span class="stat-value">{formatDate(flower.created_at)}</span>
			</div>
			<div class="stat">
				<span class="stat-label">Умер</span>
				<span class="stat-value">{flower.last_watered_at ? formatDate(flower.last_watered_at) : '—'}</span>
			</div>
			<div class="stat">
				<span class="stat-label">Дней прожил</span>
				<span class="stat-value">{flower.day}</span>
			</div>
		</div>

		<div class="actions">
			{#if isPublic}
				<ActionButton variant="muted" onclick={toggleVisibility}>
					E Скрыть
				</ActionButton>
			{:else}
				<ActionButton variant="primary" onclick={toggleVisibility}>
					E Выставить
				</ActionButton>
			{/if}
		</div>
	</div>
{/if}

<style>
	.loading {
		display: flex;
		align-items: center;
		justify-content: center;
		height: 100dvh;
		font-size: 18px;
		opacity: 0.4;
	}

	.page {
		padding: 20px var(--page-padding);
		padding-bottom: calc(var(--nav-height) + 20px);
		max-width: 440px;
		margin: 0 auto;
		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	.back {
		font-size: 14px;
		opacity: 0.6;
	}

	.flower-wrap {
		position: relative;
		display: flex;
		flex-direction: column;
		align-items: center;
		background: var(--color-card);
		padding: 32px;
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

	.flower-img {
		font-size: 120px;
		line-height: 1;
		opacity: 0.6;
	}

	.stats {
		display: flex;
		flex-direction: column;
		gap: 8px;
		background: var(--color-card);
		padding: 16px;
	}

	.stat {
		display: flex;
		justify-content: space-between;
		font-size: 14px;
	}

	.stat-label {
		opacity: 0.6;
	}

	.stat-value {
		font-weight: 500;
	}

	.actions {
		display: flex;
		justify-content: center;
	}
</style>
