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
	class="btn btn--{variant}"
	disabled={disabled || loading}
	{onclick}
>
	{#if loading}
		...
	{:else}
		{@render children()}
	{/if}
</button>

<style>
	.btn {
		display: block;
		width: var(--btn-width, 363px);
		height: var(--btn-height, 55px);
		font-size: 24px;
		font-weight: 400;
		color: #ffffff;
		text-align: center;
		cursor: pointer;
		border: none;
		border-radius: 0;
		transition: opacity 0.1s;
	}

	.btn:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}

	.btn--primary {
		background: var(--color-primary);
	}

	.btn--confirm {
		background: var(--color-confirm);
	}

	.btn--muted {
		background: var(--color-muted);
	}
</style>
