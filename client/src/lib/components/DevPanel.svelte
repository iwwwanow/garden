<script lang="ts">
	import { dev as devApi, type User } from '$lib/api';
	import { userStore, login } from '$lib/stores/auth';

	let expanded = $state(true);
	let showUsers = $state(false);
	let x = $state(16);
	let y = $state(120);
	let dragging = $state(false);
	let dragStart = $state({ mx: 0, my: 0, px: 0, py: 0 });
	let devUsers = $state<User[]>([]);
	let seedFlowerId = $state(1);
	let status = $state('');

	function onDragStart(e: PointerEvent) {
		dragging = true;
		dragStart = { mx: e.clientX, my: e.clientY, px: x, py: y };
		(e.currentTarget as HTMLElement).setPointerCapture(e.pointerId);
	}

	function onDragMove(e: PointerEvent) {
		if (!dragging) return;
		x = dragStart.px + (e.clientX - dragStart.mx);
		y = dragStart.py + (e.clientY - dragStart.my);
	}

	function onDragEnd() {
		dragging = false;
	}

	function setStatus(msg: string) {
		status = msg;
		setTimeout(() => (status = ''), 2000);
	}

	async function tick() {
		status = '...';
		try {
			await devApi.tick();
			setStatus('+24h ok');
		} catch (e: unknown) {
			status = String(e);
		}
	}

	async function giveSeeds() {
		status = '...';
		try {
			await devApi.giveSeeds(seedFlowerId);
			setStatus('seeds ok');
		} catch (e: unknown) {
			status = String(e);
		}
	}

	async function reset() {
		if (!confirm('Очистить все данные?')) return;
		status = '...';
		try {
			await devApi.reset();
			setStatus('reset ok');
		} catch (e: unknown) {
			status = String(e);
		}
	}

	async function toggleUsers() {
		if (showUsers) {
			showUsers = false;
			return;
		}
		try {
			devUsers = await devApi.listUsers();
			showUsers = true;
		} catch (e: unknown) {
			status = String(e);
		}
	}

	async function switchUser(u: User) {
		try {
			const data = await devApi.token(u.id);
			login(data.token, data.user);
			window.location.reload();
		} catch (e: unknown) {
			status = String(e);
		}
	}
</script>

<div
	class="panel"
	style="left:{x}px; top:{y}px"
	onpointermove={onDragMove}
	onpointerup={onDragEnd}
	role="complementary"
	aria-label="dev panel"
>
	<div
		class="drag-handle"
		onpointerdown={onDragStart}
		role="presentation"
	>
		<span>DEV</span>
		<button
			class="toggle-btn"
			onpointerdown={(e) => e.stopPropagation()}
			onclick={() => (expanded = !expanded)}
		>
			{expanded ? '−' : '+'}
		</button>
	</div>

	{#if expanded}
		<div class="body">
			<button onclick={tick}>+24h</button>

			<div class="seeds-row">
				<button onclick={giveSeeds}>семена</button>
				<select bind:value={seedFlowerId}>
					<option value={1}>🌹1</option>
					<option value={2}>🌷2</option>
					<option value={3}>🌼3</option>
				</select>
			</div>

			<button onclick={reset}>reset</button>

			<button onclick={toggleUsers}>
				users {showUsers ? '▴' : '▾'}
			</button>

			{#if status}
				<span class="status">{status}</span>
			{/if}

			{#if showUsers && devUsers.length}
				<div class="user-list">
					{#each devUsers as u}
						<button class="user-btn" onclick={() => switchUser(u)}>
							{u.username}
						</button>
					{/each}
				</div>
			{/if}
		</div>
	{/if}
</div>

<style>
	.panel {
		position: fixed;
		z-index: 9999;
		background: #000;
		color: #0f0;
		font-family: monospace;
		font-size: 12px;
		min-width: 130px;
		border: 1px solid #0f0;
		user-select: none;
	}

	.drag-handle {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 4px 8px;
		cursor: grab;
		background: #111;
	}

	.drag-handle:active {
		cursor: grabbing;
	}

	.toggle-btn {
		background: none;
		border: none;
		color: #0f0;
		font-family: monospace;
		font-size: 14px;
		cursor: pointer;
		padding: 0 2px;
		line-height: 1;
	}

	.body {
		padding: 6px 8px;
		display: flex;
		flex-direction: column;
		gap: 4px;
	}

	.body > button,
	.seeds-row > button {
		background: none;
		border: 1px solid #0f0;
		color: #0f0;
		font-family: monospace;
		font-size: 12px;
		padding: 2px 6px;
		cursor: pointer;
		text-align: left;
		width: 100%;
	}

	.seeds-row {
		display: flex;
		gap: 4px;
	}

	.seeds-row > button {
		width: auto;
		flex: 1;
	}

	.seeds-row select {
		background: #111;
		color: #0f0;
		border: 1px solid #0f0;
		font-family: monospace;
		font-size: 11px;
		padding: 2px;
	}

	.status {
		color: #ff0;
		font-size: 11px;
	}

	.user-list {
		display: flex;
		flex-direction: column;
		gap: 2px;
		margin-top: 2px;
		max-height: 160px;
		overflow-y: auto;
		border-top: 1px solid #0f0;
		padding-top: 4px;
	}

	.user-btn {
		background: none;
		border: 1px solid #040;
		color: #0f0;
		font-family: monospace;
		font-size: 11px;
		padding: 2px 6px;
		cursor: pointer;
		text-align: left;
		width: 100%;
	}

	.user-btn:hover {
		border-color: #0f0;
	}
</style>
