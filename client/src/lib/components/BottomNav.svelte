<script lang="ts">
	import { page } from '$app/stores';

	const tabs = [
		{ href: '/home', label: 'Сад' },
		{ href: '/leaderboard', label: 'Рейтинг' },
		{ href: '/profile', label: 'Профиль' }
	];

	const active = $derived($page.url.pathname);
</script>

<nav class="bottom-nav">
	{#each tabs as tab}
		<a
			href={tab.href}
			class:active={active === tab.href || (tab.href === '/home' && active === '/')}
		>
			{tab.label}
		</a>
	{/each}
</nav>

<style>
	.bottom-nav {
		position: fixed;
		bottom: 0;
		left: 0;
		right: 0;
	}
</style>
