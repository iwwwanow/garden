<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { me, flowers as flowersApi, totalFd, needsWatering, flowerEmoji, formatDate, type UserFlower, type FlowerTemplate, ApiError } from '$lib/api';
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
	const seedsPreview = $derived<UserFlower[]>([]); // seeds come from separate endpoint

	// Hero flower: max day among live flowers
	const heroFlower = $derived(
		liveFlowers.length ? [...liveFlowers].sort((a, b) => b.day - a.day)[0] : null
	);

	const heroFd = $derived(heroFlower ? totalFd(heroFlower.day) : 0);
	const heroNeedsWater = $derived(heroFlower ? needsWatering(heroFlower) : false);
	const heroEmoji = $derived(heroFlower ? flowerEmoji(heroFlower.flower_id) : '');

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

	function templateForFlower(f: UserFlower): FlowerTemplate | undefined {
		return templates.find((t) => t.id === f.flower_id);
	}

	onMount(loadData);
</script>

{#if loading}
	<div class="loading">загрузка...</div>
{:else}
	<!-- Hero -->
	<section class="hero">
		<div class="hero-header">
			<div>
				<div class="username">@{user?.username ?? ''}</div>
				<div class="fd-balance">{user?.fd_balance ?? 0} FD</div>
				<div class="meta">
					{liveFlowers.length} живых · {driedFlowers.length} в гербарии
				</div>
			</div>
			<a href="/profile" class="edit-btn" title="редактировать профиль">PE</a>
		</div>

		{#if heroFlower}
			<div class="hero-flower">
				<div class="hero-emoji">{heroEmoji}</div>
				<div class="hero-fd">{heroFd} FD</div>
				<div class="hero-meta">День {heroFlower.day} · Посажен {formatDate(heroFlower.created_at)}</div>
			</div>

			{#if waterError}
				<p class="water-error">{waterError}</p>
			{/if}

			{#if heroNeedsWater}
				<div class="hero-action">
					<ActionButton variant="confirm" loading={wateringId === heroFlower.id} onclick={waterHero}>
						P Полить
					</ActionButton>
				</div>
			{/if}
		{:else}
			<div class="no-flowers">У вас нет живых цветков</div>
		{/if}

		<div class="plant-action">
			<ActionButton variant="confirm" onclick={() => goto('/seeds')}>
				Посадить цветок
			</ActionButton>
		</div>
	</section>

	<!-- Seeds preview -->
	<section class="section">
		<div class="section-header">
			<h2>семена</h2>
			<a href="/seeds">Смотреть все →</a>
		</div>
		{#if liveFlowers.length === 0}
			<p class="empty">Семян нет. Они начисляются каждые 7 дней с живого цветка.</p>
		{/if}
	</section>

	<!-- Live flowers grid -->
	{#if liveFlowers.length > 0}
		<section class="section">
			<h2>цветки</h2>
			<div class="grid">
				{#each liveFlowers as flower}
					<FlowerCard
						{flower}
						type="flower"
						flowerTemplateId={flower.flower_id}
						onclick={() => goto(`/flower/${flower.id}`)}
					/>
				{/each}
			</div>
		</section>
	{/if}

	<!-- Herbarium preview -->
	{#if driedFlowers.length > 0}
		<section class="section">
			<div class="section-header">
				<h2>гербарий</h2>
				<a href="/herbarium">Смотреть все →</a>
			</div>
			<div class="grid">
				{#each driedFlowers.slice(0, 3) as flower}
					<FlowerCard
						{flower}
						type="herbarium"
						flowerTemplateId={flower.flower_id}
						onclick={() => goto(`/herbarium/${flower.id}`)}
					/>
				{/each}
			</div>
		</section>
	{/if}
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

	.hero {
		min-height: 100dvh;
		display: flex;
		flex-direction: column;
		padding: 24px var(--page-padding);
		max-width: 440px;
		margin: 0 auto;
	}

	.hero-header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
	}

	.username {
		font-size: 18px;
		font-weight: 500;
	}

	.fd-balance {
		font-size: 48px;
		font-weight: 600;
		line-height: 1.1;
	}

	.meta {
		font-size: 13px;
		opacity: 0.5;
		margin-top: 4px;
	}

	.edit-btn {
		font-size: 14px;
		font-weight: 600;
		padding: 6px 10px;
		background: var(--color-card);
		color: var(--color-text);
	}

	.hero-flower {
		flex: 1;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 8px;
	}

	.hero-emoji {
		font-size: 128px;
		line-height: 1;
	}

	.hero-fd {
		font-size: 28px;
		font-weight: 600;
	}

	.hero-meta {
		font-size: 13px;
		opacity: 0.5;
	}

	.hero-action,
	.plant-action {
		display: flex;
		justify-content: center;
		margin-top: 12px;
	}

	.no-flowers {
		flex: 1;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 16px;
		opacity: 0.4;
	}

	.water-error {
		color: #c00;
		font-size: 13px;
		text-align: center;
		margin-top: 8px;
	}

	.section {
		padding: 20px var(--page-padding);
		max-width: 440px;
		margin: 0 auto;
		width: 100%;
	}

	.section-header {
		display: flex;
		justify-content: space-between;
		align-items: baseline;
		margin-bottom: 12px;
	}

	.section-header a {
		font-size: 13px;
		opacity: 0.6;
	}

	.grid {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: 8px;
		margin-top: 12px;
	}

	.empty {
		font-size: 14px;
		opacity: 0.5;
	}
</style>
