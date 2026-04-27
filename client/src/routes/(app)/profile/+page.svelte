<script lang="ts">
	import { goto } from '$app/navigation';
	import { me as meApi, ApiError } from '$lib/api';
	import { userStore, logout } from '$lib/stores/auth';
	import ActionButton from '$lib/components/ActionButton.svelte';

	const user = $derived($userStore);

	let editing = $state(false);
	let firstName = $state('');
	let saving = $state(false);
	let error = $state('');
	let success = $state('');

	function startEdit() {
		firstName = user?.first_name ?? '';
		editing = true;
		error = '';
		success = '';
	}

	async function saveEdit() {
		saving = true;
		error = '';
		try {
			const updated = await meApi.update(firstName);
			userStore.set(updated);
			editing = false;
			success = 'Сохранено';
			setTimeout(() => (success = ''), 2000);
		} catch (e: unknown) {
			error = e instanceof ApiError ? e.message : 'Ошибка';
		} finally {
			saving = false;
		}
	}

	function doLogout() {
		logout();
		goto('/auth');
	}
</script>

<div class="page">
	<h1>Профиль</h1>

	{#if user}
		<div class="profile-card">
			<div class="profile-row">
				<span class="label">Имя пользователя</span>
				<span class="value">@{user.username}</span>
			</div>
			<div class="profile-row">
				<span class="label">Имя</span>
				<span class="value">{user.first_name || '—'}</span>
			</div>
			<div class="profile-row">
				<span class="label">FD баланс</span>
				<span class="value fd">{user.fd_balance} FD</span>
			</div>
		</div>

		{#if !editing}
			<div class="actions">
				<ActionButton variant="primary" onclick={startEdit}>
					PE Редактировать имя
				</ActionButton>
			</div>
		{:else}
			<div class="edit-form">
				<div class="field">
					<label for="firstname">Имя</label>
					<input id="firstname" type="text" bind:value={firstName} />
				</div>
				{#if error}
					<p class="error">{error}</p>
				{/if}
				<div class="edit-actions">
					<ActionButton variant="confirm" loading={saving} onclick={saveEdit}>
						Сохранить
					</ActionButton>
					<ActionButton variant="muted" onclick={() => (editing = false)}>
						Отмена
					</ActionButton>
				</div>
			</div>
		{/if}

		{#if success}
			<p class="success">{success}</p>
		{/if}

		<div class="logout-wrap">
			<ActionButton variant="muted" onclick={doLogout}>
				Выйти
			</ActionButton>
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

	.profile-card {
		display: flex;
		flex-direction: column;
		gap: 0;
		background: var(--color-card);
	}

	.profile-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 14px 16px;
		border-bottom: 1px solid rgba(0, 0, 0, 0.06);
	}

	.label {
		font-size: 14px;
		opacity: 0.6;
	}

	.value {
		font-size: 15px;
		font-weight: 500;
	}

	.fd {
		font-size: 20px;
		font-weight: 600;
	}

	.actions,
	.logout-wrap {
		display: flex;
		justify-content: center;
	}

	.edit-form {
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

	.edit-actions {
		display: flex;
		flex-direction: column;
		gap: 8px;
		align-items: center;
	}

	.error {
		color: #c00;
		font-size: 14px;
	}

	.success {
		color: var(--color-confirm);
		font-size: 14px;
		text-align: center;
	}
</style>
