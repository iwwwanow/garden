<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { me, flowers as flowersApi, totalFd, needsWatering, flowerEmoji, formatDate, ApiError, type UserFlower } from '$lib/api';
	import { userStore } from '$lib/stores/auth';
	import ActionButton from '$lib/components/ActionButton.svelte';

	const flowerId = $derived(Number($page.params.id));
	const currentUser = $derived($userStore);

	let flower = $state<UserFlower | null>(null);
	let loading = $state(true);
	let watering = $state(false);
	let error = $state('');
	let wateredOk = $state(false);

	const fd = $derived(flower ? totalFd(flower.day) : 0);
	const canWater = $derived(flower ? needsWatering(flower) : false);
	const emoji = $derived(flower ? flowerEmoji(flower.flower_id) : '');

	async function load() {
		loading = true;
		error = '';
		try {
			const data = await me.get();
			flower = data.flowers.find((f) => f.id === flowerId) ?? null;
		} finally {
			loading = false;
		}
	}

	async function water() {
		if (!flower) return;
		watering = true;
		error = '';
		try {
			await flowersApi.water(flower.id);
			wateredOk = true;
			await load();
		} catch (e: unknown) {
			error = e instanceof ApiError ? e.message : 'Ошибка полива';
		} finally {
			watering = false;
		}
	}

	onMount(load);
</script>

{#if loading}
	<div class="loading">загрузка...</div>
{:else if !flower}
	<div class="page">
		<p>Цветок не найден</p>
		<a href="/home">← назад</a>
	</div>
{:else}
	<div class="page">
		<a href="/garden/{flower.user_id}" class="owner-tag">
			owner: @{currentUser?.id === flower.user_id ? currentUser?.username : flower.user_id}
		</a>

		<div class="flower-wrap">
			<span class="type-badge">{emoji}</span>
			<div class="flower-img">{emoji}</div>
			<div class="flower-fd">{fd} FD</div>
		</div>

		<div class="stats">
			<div class="stat">
				<span class="stat-label">Посажен</span>
				<span class="stat-value">{formatDate(flower.created_at)}</span>
			</div>
			<div class="stat">
				<span class="stat-label">День роста</span>
				<span class="stat-value">{flower.day}</span>
			</div>
			<div class="stat">
				<span class="stat-label">Принес за все время</span>
				<span class="stat-value">{fd} FD</span>
			</div>
		</div>

		{#if error}
			<p class="error">{error}</p>
		{/if}
		{#if wateredOk}
			<p class="success">Полито!</p>
		{/if}

		<div class="actions">
			{#if canWater}
				<ActionButton variant="confirm" loading={watering} onclick={water}>
					P Полить
				</ActionButton>
			{:else}
				<ActionButton variant="muted" disabled>
					Уже полит сегодня
				</ActionButton>
			{/if}
			<ActionButton variant="primary">
				Z Запросить полив
			</ActionButton>
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

	.owner-tag {
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
		gap: 12px;
	}

	.type-badge {
		position: absolute;
		top: 12px;
		left: 12px;
		font-size: 18px;
	}

	.flower-img {
		font-size: 120px;
		line-height: 1;
	}

	.flower-fd {
		font-size: 28px;
		font-weight: 600;
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
		flex-direction: column;
		gap: 10px;
		align-items: center;
	}

	.error {
		color: #c00;
		font-size: 14px;
	}

	.success {
		color: var(--color-confirm);
		font-size: 14px;
	}
</style>
