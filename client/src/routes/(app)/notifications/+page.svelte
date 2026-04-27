<script lang="ts">
	import { onMount } from 'svelte';
	import { notifications as notifApi, formatDate, type Notification } from '$lib/api';
	import ActionButton from '$lib/components/ActionButton.svelte';

	let items = $state<Notification[]>([]);
	let loading = $state(true);
	let marking = $state(false);

	const unread = $derived(items.filter((n) => !n.is_read).length);

	function notifText(n: Notification): string {
		const payload = n.payload as Record<string, unknown>;
		switch (n.type) {
			case 'flower_watered':
				return `Ваш цветок полили`;
			case 'flower_died':
				return `Цветок засох`;
			case 'seed_awarded':
				return `Получено семя`;
			default:
				return n.type;
		}
	}

	async function load() {
		loading = true;
		try {
			items = (await notifApi.list()) ?? [];
		} finally {
			loading = false;
		}
	}

	async function markAllRead() {
		marking = true;
		try {
			await notifApi.markAllRead();
			items = items.map((n) => ({ ...n, is_read: true }));
		} finally {
			marking = false;
		}
	}

	onMount(load);
</script>

<div class="page">
	<div class="header">
		<h1>Уведомления</h1>
		{#if unread > 0}
			<ActionButton variant="muted" loading={marking} onclick={markAllRead}>
				Прочитать все
			</ActionButton>
		{/if}
	</div>

	{#if loading}
		<p class="loading">загрузка...</p>
	{:else if items.length === 0}
		<p class="empty">Нет уведомлений</p>
	{:else}
		<ul class="list">
			{#each items as notif}
				<li class="item" class:item--unread={!notif.is_read}>
					<span class="item-text">{notifText(notif)}</span>
					<span class="item-date">{formatDate(notif.created_at)}</span>
				</li>
			{/each}
		</ul>
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

	.header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 12px;
	}

	.header :global(.btn) {
		width: auto;
		height: auto;
		font-size: 14px;
		padding: 8px 14px;
	}

	.list {
		list-style: none;
		display: flex;
		flex-direction: column;
		gap: 0;
	}

	.item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 14px 8px;
		border-bottom: 1px solid var(--color-card);
		opacity: 0.5;
	}

	.item--unread {
		opacity: 1;
	}

	.item-text {
		font-size: 14px;
	}

	.item-date {
		font-size: 12px;
		opacity: 0.5;
		flex-shrink: 0;
	}

	.loading,
	.empty {
		opacity: 0.4;
		font-size: 16px;
	}
</style>
