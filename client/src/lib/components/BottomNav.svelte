<script lang="ts">
	import { page } from '$app/stores';

	const tabs = [
		{ href: '/home', label: 'Сад', icon: '⬡' },
		{ href: '/leaderboard', label: 'Рейтинг', icon: '▲' },
		{ href: '/profile', label: 'Профиль', icon: '○' }
	];

	const active = $derived($page.url.pathname);
</script>

<nav class="bottom-nav">
	{#each tabs as tab}
		<a
			href={tab.href}
			class="tab"
			class:tab--active={active === tab.href || (tab.href === '/home' && active === '/')}
		>
			<span class="tab-icon">{tab.icon}</span>
			<span class="tab-label">{tab.label}</span>
		</a>
	{/each}
</nav>

<style>
	.bottom-nav {
		position: fixed;
		bottom: 0;
		left: 0;
		right: 0;
		height: var(--nav-height);
		background: var(--color-nav);
		display: flex;
		align-items: center;
		justify-content: space-around;
		z-index: 100;
	}

	.tab {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 4px;
		flex: 1;
		padding: 12px 0;
		color: var(--color-text);
		opacity: 0.45;
		text-decoration: none;
		font-size: 12px;
	}

	.tab--active {
		opacity: 1;
	}

	.tab-icon {
		font-size: 22px;
	}

	.tab-label {
		font-size: 11px;
		font-weight: 500;
	}
</style>
