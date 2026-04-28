<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { me, flowers as flowersApi, totalFd, needsWatering, formatDate, ApiError, type UserFlower, type FlowerTemplate } from '$lib/api';
	import { userStore } from '$lib/stores/auth';
	import ActionButton from '$lib/components/ActionButton.svelte';

	const flowerId = $derived(Number($page.params.id));
	const currentUser = $derived($userStore);

	let flower = $state<UserFlower | null>(null);
	let template = $state<FlowerTemplate | null>(null);
	let loading = $state(true);
	let watering = $state(false);
	let error = $state('');
	let wateredOk = $state(false);

	const fd = $derived(flower ? totalFd(flower.day) : 0);
	const canWater = $derived(flower ? needsWatering(flower) : false);

	async function load() {
		loading = true; error = '';
		try {
			const [data, tmpls] = await Promise.all([me.get(), flowersApi.listTemplates()]);
			flower = data.flowers.find((f) => f.id === flowerId) ?? null;
			if (flower) template = tmpls.find((t) => t.id === flower!.flower_id) ?? null;
		} finally {
			loading = false;
		}
	}

	async function water() {
		if (!flower) return;
		watering = true; error = '';
		try { await flowersApi.water(flower.id); wateredOk = true; await load(); }
		catch (e: unknown) { error = e instanceof ApiError ? e.message : 'Ошибка полива'; }
		finally { watering = false; }
	}

	onMount(load);
</script>

{#if loading}
	<p>загрузка...</p>
{:else if !flower}
	<div class="page">
		<a href="/home">← Сад</a>
		<p>Цветок не найден</p>
	</div>
{:else}
	<div class="page">
		<a href="/garden/{flower.user_id}">
			owner: @{currentUser?.id === flower.user_id ? currentUser?.username : flower.user_id}
		</a>

		{#if template?.image_path}
			<img src="/{template.image_path}" alt="цветок" style="width:100%;" />
		{/if}
		<p>{fd} FD · День {flower.day}</p>

		<table><tbody>
			<tr><td>Посажен</td><td>{formatDate(flower.created_at)}</td></tr>
			<tr><td>День роста</td><td>{flower.day}</td></tr>
			<tr><td>Принёс всего</td><td>{fd} FD</td></tr>
		</tbody></table>

		{#if error}<p>{error}</p>{/if}
		{#if wateredOk}<p>Полито!</p>{/if}

		{#if canWater}
			<ActionButton variant="confirm" loading={watering} onclick={water}>P Полить</ActionButton>
		{:else}
			<ActionButton variant="muted" disabled>Уже полит сегодня</ActionButton>
		{/if}
		<ActionButton variant="primary">Z Запросить полив</ActionButton>
	</div>
{/if}
