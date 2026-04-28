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

	<div>
		<button>FD</button>
		<button disabled>Дней подряд</button>
		<button disabled>Семена</button>
	</div>

	{#if loading}
		<p>загрузка...</p>
	{:else if entries.length === 0}
		<p>список пуст</p>
	{:else}
		<ol>
			{#each entries as user, i}
				<li class:me={currentUser?.id === user.id}>
					<span>{i + 1}.</span>
					<a href="/garden/{user.id}">@{user.username}</a>
					<span>{user.fd_balance} FD</span>
				</li>
			{/each}
		</ol>
	{/if}
</div>

<style>
	.me { background: var(--color-highlight); }
</style>
