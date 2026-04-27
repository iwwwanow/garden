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
				const user = await authApi.register(username, password, firstName);
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

<div class="page page--no-nav auth-page">
	<h1 class="title">Сад богини цветов</h1>

	<div class="tabs">
		<button class:active={mode === 'login'} onclick={() => (mode = 'login')}>Войти</button>
		<button class:active={mode === 'register'} onclick={() => (mode = 'register')}>Регистрация</button>
	</div>

	<form onsubmit={(e) => { e.preventDefault(); submit(); }}>
		<div class="field">
			<label for="username">Имя пользователя</label>
			<input id="username" type="text" bind:value={username} autocomplete="username" required />
		</div>

		{#if mode === 'register'}
			<div class="field">
				<label for="firstname">Имя</label>
				<input id="firstname" type="text" bind:value={firstName} />
			</div>
		{/if}

		<div class="field">
			<label for="password">Пароль</label>
			<input id="password" type="password" bind:value={password} autocomplete={mode === 'login' ? 'current-password' : 'new-password'} required />
		</div>

		{#if error}
			<p class="error">{error}</p>
		{/if}

		<div class="btn-wrap">
			<ActionButton variant="confirm" type="submit" {loading}>
				{mode === 'login' ? 'Войти' : 'Создать аккаунт'}
			</ActionButton>
		</div>
	</form>
</div>

<style>
	.auth-page {
		display: flex;
		flex-direction: column;
		gap: 32px;
		padding-top: 64px;
	}

	.title {
		font-size: 28px;
		font-weight: 600;
		text-align: center;
	}

	.tabs {
		display: flex;
		gap: 0;
		border-bottom: 2px solid var(--color-card);
	}

	.tabs button {
		flex: 1;
		padding: 10px;
		font-size: 16px;
		font-weight: 500;
		background: none;
		border: none;
		cursor: pointer;
		opacity: 0.4;
	}

	.tabs button.active {
		opacity: 1;
		border-bottom: 2px solid var(--color-text);
		margin-bottom: -2px;
	}

	form {
		display: flex;
		flex-direction: column;
		gap: 16px;
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
		border: 1px solid var(--color-card);
		padding: 12px;
		font-size: 16px;
		font-family: inherit;
		border-radius: 0;
		outline: none;
		width: 100%;
	}

	input:focus {
		border-color: var(--color-text);
	}

	.error {
		color: #c00;
		font-size: 14px;
	}

	.btn-wrap {
		display: flex;
		justify-content: center;
		margin-top: 8px;
	}
</style>
