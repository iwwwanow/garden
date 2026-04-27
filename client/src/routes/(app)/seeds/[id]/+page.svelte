<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { seeds as seedsApi, flowers as flowersApi, users as usersApi, flowerEmoji, formatDate, ApiError, type Seed, type FlowerTemplate, type User } from '$lib/api';
	import ActionButton from '$lib/components/ActionButton.svelte';

	const seedId = $derived(Number($page.params.id));

	let seed = $state<Seed | null>(null);
	let template = $state<FlowerTemplate | null>(null);
	let loading = $state(true);
	let planting = $state(false);
	let sharing = $state(false);
	let shareUsername = $state('');
	let shareQty = $state(1);
	let resolvedUser = $state<User | null>(null);
	let resolving = $state(false);
	let showShare = $state(false);
	let error = $state('');
	let success = $state('');

	async function load() {
		loading = true;
		try {
			const [list, tmpls] = await Promise.all([seedsApi.list(), flowersApi.listTemplates()]);
			seed = list.find((s) => s.id === seedId) ?? null;
			if (seed) {
				template = tmpls.find((t) => t.id === seed!.flower_id) ?? null;
			}
		} finally {
			loading = false;
		}
	}

	async function plant() {
		if (!seed) return;
		planting = true;
		error = '';
		try {
			await flowersApi.plant(seed.flower_id);
			goto('/home');
		} catch (e: unknown) {
			error = e instanceof ApiError ? e.message : 'Ошибка посадки';
		} finally {
			planting = false;
		}
	}

	async function resolveUsername() {
		if (!shareUsername.trim()) { resolvedUser = null; return; }
		resolving = true;
		error = '';
		try {
			resolvedUser = await usersApi.getByUsername(shareUsername.trim());
		} catch {
			resolvedUser = null;
			error = 'Пользователь не найден';
		} finally {
			resolving = false;
		}
	}

	async function share() {
		if (!seed || !resolvedUser) return;
		sharing = true;
		error = '';
		try {
			await seedsApi.share(resolvedUser.id, seed.flower_id, shareQty);
			success = `Передано @${resolvedUser.username}`;
			showShare = false;
			shareUsername = '';
			resolvedUser = null;
			await load();
		} catch (e: unknown) {
			error = e instanceof ApiError ? e.message : 'Ошибка передачи';
		} finally {
			sharing = false;
		}
	}

	onMount(load);
</script>

{#if loading}
	<div class="loading">загрузка...</div>
{:else if !seed}
	<div class="page">
		<p>Семя не найдено</p>
		<a href="/seeds">← назад</a>
	</div>
{:else}
	<div class="page">
		<a href="/seeds" class="back">← Семена</a>

		<div class="flower-wrap">
			<span class="type-badge">S</span>
			<div class="flower-img">{flowerEmoji(seed.flower_id)}</div>
			<div class="qty">x{seed.quantity}</div>
		</div>

		{#if error}
			<p class="error">{error}</p>
		{/if}
		{#if success}
			<p class="success">{success}</p>
		{/if}

		<div class="actions">
			<ActionButton variant="confirm" loading={planting} onclick={plant}>
				Посадить цветок
			</ActionButton>
			<ActionButton variant="primary" onclick={() => { showShare = !showShare; error = ''; }}>
				Поделиться
			</ActionButton>
		</div>

		{#if showShare}
			<div class="share-form">
				<div class="field">
					<label for="share-to">Имя пользователя</label>
					<div class="username-row">
						<input
							id="share-to"
							type="text"
							bind:value={shareUsername}
							placeholder="@username"
							onblur={resolveUsername}
						/>
						{#if resolving}
							<span class="resolve-hint">...</span>
						{:else if resolvedUser}
							<span class="resolve-ok">✓ {resolvedUser.first_name || resolvedUser.username}</span>
						{/if}
					</div>
				</div>
				<div class="field">
					<label for="share-qty">Количество (макс. {seed.quantity})</label>
					<input id="share-qty" type="number" bind:value={shareQty} min="1" max={seed.quantity} />
				</div>
				<ActionButton variant="confirm" loading={sharing} disabled={!resolvedUser} onclick={share}>
					Передать
				</ActionButton>
			</div>
		{/if}
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
		gap: 12px;
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
		font-size: 100px;
		line-height: 1;
	}

	.qty {
		font-size: 24px;
		font-weight: 600;
	}

	.actions {
		display: flex;
		flex-direction: column;
		gap: 10px;
		align-items: center;
	}

	.share-form {
		display: flex;
		flex-direction: column;
		gap: 12px;
		background: var(--color-card);
		padding: 16px;
	}

	.field {
		display: flex;
		flex-direction: column;
		gap: 6px;
	}

	label {
		font-size: 14px;
		font-weight: 500;
	}

	input {
		border: 1px solid #ccc;
		padding: 10px;
		font-size: 16px;
		font-family: inherit;
		border-radius: 0;
		width: 100%;
	}

	.username-row {
		display: flex;
		align-items: center;
		gap: 8px;
	}

	.username-row input {
		flex: 1;
	}

	.resolve-hint {
		font-size: 13px;
		opacity: 0.5;
	}

	.resolve-ok {
		font-size: 13px;
		color: var(--color-confirm);
		white-space: nowrap;
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
