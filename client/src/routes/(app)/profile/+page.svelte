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
		editing = true; error = ''; success = '';
	}

	async function saveEdit() {
		saving = true; error = '';
		try {
			const updated = await meApi.update(firstName);
			userStore.set(updated);
			editing = false;
			success = 'Сохранено';
			setTimeout(() => (success = ''), 2000);
		} catch (e: unknown) {
			error = e instanceof ApiError ? e.message : 'Ошибка';
		} finally { saving = false; }
	}

	function doLogout() { logout(); goto('/auth'); }
</script>

<div class="page">
	<h1>Профиль</h1>

	{#if user}
		<table><tbody>
			<tr><td>Имя пользователя</td><td>@{user.username}</td></tr>
			<tr><td>Имя</td><td>{user.first_name || '—'}</td></tr>
			<tr><td>FD баланс</td><td>{user.fd_balance} FD</td></tr>
		</tbody></table>

		{#if !editing}
			<ActionButton variant="primary" onclick={startEdit}>PE Редактировать имя</ActionButton>
		{:else}
			<label for="firstname">Имя</label>
			<input class="input" id="firstname" type="text" bind:value={firstName} />
			{#if error}<p>{error}</p>{/if}
			<ActionButton variant="confirm" loading={saving} onclick={saveEdit}>Сохранить</ActionButton>
			<ActionButton variant="muted" onclick={() => (editing = false)}>Отмена</ActionButton>
		{/if}

		{#if success}<p>{success}</p>{/if}

		<hr />
		<ActionButton variant="muted" onclick={doLogout}>Выйти</ActionButton>
	{/if}
</div>
