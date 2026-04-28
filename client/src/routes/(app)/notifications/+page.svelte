<script lang="ts">
	import { onMount } from 'svelte';
	import { notifications as notifApi, formatDate, type Notification } from '$lib/api';

	let items = $state<Notification[]>([]);
	let loading = $state(true);
	let marking = $state(false);

	const unread = $derived(items.filter((n) => !n.is_read).length);

	function notifText(n: Notification): string {
		switch (n.type) {
			case 'flower_watered': return 'Ваш цветок полили';
			case 'flower_died':    return 'Цветок засох';
			case 'seed_awarded':   return 'Получено семя';
			default:               return n.type;
		}
	}

	async function load() {
		loading = true;
		try { items = (await notifApi.list()) ?? []; }
		finally { loading = false; }
	}

	async function markAllRead() {
		marking = true;
		try {
			await notifApi.markAllRead();
			items = items.map((n) => ({ ...n, is_read: true }));
		} finally { marking = false; }
	}

	onMount(load);
</script>

<div class="page">
	<h1>Уведомления {#if unread > 0}({unread}){/if}</h1>
	{#if unread > 0}
		<button onclick={markAllRead} disabled={marking}>Прочитать все</button>
	{/if}

	{#if loading}
		<p>загрузка...</p>
	{:else if items.length === 0}
		<p>нет уведомлений</p>
	{:else}
		<ul>
			{#each items as notif}
				<li class:read={notif.is_read}>
					<span>{notifText(notif)}</span>
					<span>{formatDate(notif.created_at)}</span>
				</li>
			{/each}
		</ul>
	{/if}
</div>

<style>
	.read { opacity: 0.5; }
</style>
