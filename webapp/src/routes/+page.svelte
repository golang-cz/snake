<script lang="ts">
	import { SnakeGame } from '$lib/rpc.gen';
	import { Game } from './game';

	let canvas: HTMLCanvasElement;

	const api = new SnakeGame('http://localhost:5252', fetch);

	$effect(() => {
		const game = new Game(canvas.getContext('2d')!, api);
		game.start();
	});
</script>

<div class="wrapper">
	<canvas tabindex="1" bind:this={canvas} height={500} width={500}></canvas>
</div>

<style>
	.wrapper {
		background: -webkit-linear-gradient(top, #7fc5c9, #deabbe);
		text-align: center;
		height: 100vh;
		display: flex;
		flex-direction: column;
		justify-content: center;
	}
	canvas {
		margin: auto;
	}
</style>
