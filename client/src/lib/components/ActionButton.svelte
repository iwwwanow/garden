<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		variant?: 'primary' | 'confirm' | 'muted';
		type?: 'button' | 'submit';
		disabled?: boolean;
		loading?: boolean;
		onclick?: () => void;
		children: Snippet;
	}

	let {
		variant = 'primary',
		type = 'button',
		disabled = false,
		loading = false,
		onclick,
		children
	}: Props = $props();
</script>

<button
	{type}
	data-variant={variant}
	class="action-btn"
	disabled={disabled || loading}
	{onclick}
>
	{#if loading}
		<span class="spinner spinner-sm"></span>
	{:else}
		{@render children()}
	{/if}
</button>

<style>
	/* структурное: размер из токенов */
	.action-btn {
		width: var(--btn-width, 363px);
		height: var(--btn-height, 55px);
	}
</style>
