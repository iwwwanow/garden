<script lang="ts">
	import { totalFd, needsWatering, flowerEmoji, type UserFlower } from '$lib/api';

	interface Props {
		flower: UserFlower;
		type: 'flower' | 'seed' | 'herbarium';
		flowerTemplateId?: number;
		onclick?: () => void;
	}

	let { flower, type, flowerTemplateId, onclick }: Props = $props();

	const fd = $derived(totalFd(flower.day));
	const watering = $derived(type === 'flower' && needsWatering(flower));
	const emoji = $derived(flowerEmoji(flowerTemplateId ?? flower.flower_id));
	const typeIcon = $derived(type === 'seed' ? 'S' : type === 'herbarium' ? 'H' : emoji);
</script>

<button class="card" {onclick} aria-label="цветок">
	<span class="type-icon">{typeIcon}</span>
	<span class="flower-img">{emoji}</span>
	<div class="card-footer">
		<span class="fd">{fd} FD</span>
		{#if watering}
			<span class="needs-water" title="нуждается в поливе">P</span>
		{/if}
	</div>
</button>

<style>
	.card {
		position: relative;
		display: flex;
		flex-direction: column;
		align-items: stretch;
		background: var(--color-card);
		width: 100%;
		aspect-ratio: 149 / 230;
		cursor: pointer;
		border: none;
		padding: 0;
		font-family: inherit;
	}

	.type-icon {
		position: absolute;
		top: 8px;
		left: 8px;
		font-size: 14px;
		font-weight: 600;
		background: rgba(255, 255, 255, 0.8);
		padding: 2px 5px;
		line-height: 1;
	}

	.flower-img {
		flex: 1;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 64px;
	}

	.card-footer {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 6px 8px;
		background: rgba(0, 0, 0, 0.05);
		font-size: 12px;
		font-weight: 500;
	}

	.fd {
		font-size: 12px;
	}

	.needs-water {
		font-size: 12px;
		font-weight: 600;
		color: var(--color-confirm);
	}
</style>
