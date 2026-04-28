<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth as authApi } from '$lib/api';
	import { login } from '$lib/stores/auth';
	import ActionButton from '$lib/components/ActionButton.svelte';

	let mode = $state<'login' | 'register'>('login');
	let username = $state('');
	let password = $state('');
	let firstName = $state('');
	let error = $state('');
	let loading = $state(false);

	async function submit() {
		error = '';
		loading = true;
		try {
			if (mode === 'login') {
				const res = await authApi.login(username, password);
				login(res.token, res.user);
				goto('/home');
			} else {
				await authApi.register(username, password, firstName);
				const res = await authApi.login(username, password);
				login(res.token, res.user);
				goto('/home');
			}
		} catch (e: unknown) {
			error = e instanceof Error ? e.message : 'Ошибка';
		} finally {
			loading = false;
		}
	}
</script>

<div class="page" style="padding-bottom: var(--page-padding);">
	<h1>Сад богини цветов</h1>

	<div>
		<button onclick={() => (mode = 'login')}>Войти</button>
		<button onclick={() => (mode = 'register')}>Регистрация</button>
	</div>

	<form onsubmit={(e) => { e.preventDefault(); submit(); }}>
		<div>
			<label for="username">Имя пользователя</label>
			<input id="username" type="text" bind:value={username} autocomplete="username" required />
		</div>

		{#if mode === 'register'}
			<div>
				<label for="firstname">Имя</label>
				<input id="firstname" type="text" bind:value={firstName} />
			</div>
		{/if}

		<div>
			<label for="password">Пароль</label>
			<input id="password" type="password" bind:value={password}
				autocomplete={mode === 'login' ? 'current-password' : 'new-password'} required />
		</div>

		{#if error}<p>{error}</p>{/if}

		<ActionButton variant="confirm" type="submit" {loading}>
			{mode === 'login' ? 'Войти' : 'Создать аккаунт'}
		</ActionButton>
	</form>
</div>
