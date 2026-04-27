<script lang="ts">
	import { onMount } from 'svelte';
	import { leaderboard as lbApi, type User } from '$lib/api';
	import { userStore } from '$lib/stores/auth';

	let entries = $state<User[]>([]);
	let loading = $state(true);
	const currentUser = $derived($userStore);

	async function load() {
		loading = true;
		try {
			entries = (await lbApi.get()) ?? [];
		} finally {
			loading = false;
		}
	}

	onMount(load);
</script>

<div class="page">
	<h1>Рейтинг</h1>

	<div class="tabs">
		<button class="tab tab--active">FD</button>
		<button class="tab tab--disabled" title="скоро">Дней подряд</button>
		<button class="tab tab--disabled" title="скоро">Семена</button>
	</div>

	{#if loading}
		<p class="loading">загрузка...</p>
	{:else if entries.length === 0}
		<p class="empty">Список пуст</p>
	{:else}
		<ol class="list">
			{#each entries as user, i}
				<li
					class="row"
					class:row--me={currentUser?.id === user.id}
					aria-current={currentUser?.id === user.id ? 'true' : undefined}
				>
					<span class="rank">{i + 1}</span>
					<a href="/garden/{user.id}" class="username">@{user.username}</a>
					<span class="score">{user.fd_balance} FD</span>
				</li>
			{/each}
		</ol>
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

	.tabs {
		display: flex;
		gap: 0;
		border-bottom: 2px solid var(--color-card);
	}

	.tab {
		flex: 1;
		padding: 8px 4px;
		font-size: 14px;
		font-weight: 500;
		background: none;
		border: none;
		cursor: pointer;
		opacity: 0.4;
		font-family: inherit;
	}

	.tab--active {
		opacity: 1;
		border-bottom: 2px solid var(--color-text);
		margin-bottom: -2px;
	}

	.tab--disabled {
		cursor: not-allowed;
	}

	.list {
		list-style: none;
		display: flex;
		flex-direction: column;
		gap: 0;
	}

	.row {
		display: flex;
		align-items: center;
		gap: 12px;
		padding: 12px 8px;
		border-bottom: 1px solid var(--color-card);
	}

	.row--me {
		background: var(--color-highlight);
	}

	.rank {
		width: 28px;
		font-size: 14px;
		font-weight: 600;
		flex-shrink: 0;
	}

	.username {
		flex: 1;
		font-size: 15px;
	}

	.score {
		font-size: 15px;
		font-weight: 600;
	}

	.loading,
	.empty {
		opacity: 0.4;
		font-size: 16px;
	}
</style>
