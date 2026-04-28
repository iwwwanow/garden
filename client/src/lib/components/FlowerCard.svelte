<script lang="ts">
	import { totalFd, needsWatering, type UserFlower } from '$lib/api';

	interface Props {
		flower: UserFlower;
		type: 'flower' | 'seed' | 'herbarium';
		imagePath?: string;
		onclick?: () => void;
	}

	let { flower, type, imagePath, onclick }: Props = $props();

	const fd = $derived(totalFd(flower.day));
	const watering = $derived(type === 'flower' && needsWatering(flower));
	const typeLabel = $derived(type === 'seed' ? 'S' : type === 'herbarium' ? 'H' : null);
</script>

<button {onclick} aria-label="цветок" style="width:100%; aspect-ratio:190/288; overflow:hidden; display:flex; flex-direction:column;">
	<div style="flex:1; position:relative;">
		{#if imagePath}
			<img src={imagePath} alt="" style="width:100%; height:100%; object-fit:cover; display:block;" />
		{/if}
		{#if typeLabel}<span>{typeLabel}</span>{/if}
		{#if watering}<span>P</span>{/if}
	</div>
	<div>
		<span>{fd} FD</span>
		<span>д.{flower.day}</span>
	</div>
</button>
