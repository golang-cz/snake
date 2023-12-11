<script lang="ts">
	import { onMount } from 'svelte';
	//import { Snake } from './rpc.gen';

	let canvas: HTMLCanvasElement;
	let ctx: CanvasRenderingContext2D;

	//const api = new Snake('http://localhost:5174', fetch);

	// Mock data
	let state = {
		snakes: [
			{
				id: 1,
				color: 'blue',
				body: [
					{ x: 10, y: 20 },
					{ x: 10, y: 21 },
					{ x: 11, y: 21 },
					{ x: 12, y: 21 },
					{ x: 13, y: 21 },
					{ x: 14, y: 21 }
				],
				direction: 'up'
			},
			{
				id: 1,
				color: 'green',
				body: [
					{ x: 50, y: 41 },
					{ x: 50, y: 42 },
					{ x: 50, y: 43 },
					{ x: 50, y: 44 }
				],
				direction: 'down'
			}
		],
		items: [
			{
				id: 1,
				color: 'red',
				body: [{ x: 10, y: 10 }]
			}
		]
	};

	const width = 70;
	const height = 70;

	// Constants
	const cellSize = 10;
	const pxHeight = `${width * cellSize}px`;
	const pxWidth = `${height * cellSize}px`;

	function drawSnakes() {
		state.snakes.forEach((snake) => {
			for (let i = 0; i < snake.body.length; i++) {
				drawSquare(snake.body[i].x, snake.body[i].y, snake.color);
			}
		});
	}
	function drawItems() {
		state.items.forEach((item) => {
			for (let i = 0; i < item.body.length; i++) {
				drawSquare(item.body[i].x, item.body[i].y, item.color);
			}
		});
	}
	function drawSquare(x: number, y: number, color: string) {
		ctx.fillStyle = color;
		ctx.fillRect(x * cellSize, y * cellSize, cellSize, cellSize);
	}

	// api.joinGame(
	// 	{
	// 		onMessage: (msg) => {
	// 			console.log(snake);
	// 		}
	// 	}
	// );

	onMount(async () => {
		ctx = canvas.getContext('2d')!;
		canvas.focus();

		function drawGrid() {
			ctx.fillStyle = '#fff';
			ctx.strokeStyle = '#eee';

			ctx.fillRect(0, 0, canvas.height, canvas.width);

			for (let x = 0.5; x < canvas.width; x += cellSize) {
				ctx.moveTo(x, 0);
				ctx.lineTo(x, canvas.height);
			}
			for (let y = 0.5; y < canvas.height; y += cellSize) {
				ctx.moveTo(0, y);
				ctx.lineTo(canvas.width, y);
			}

			ctx.stroke();
		}

		ctx.beginPath();
		drawGrid();
		drawSnakes();
		drawItems();
	});
</script>

<div class="wrapper">
	<canvas tabindex="1" bind:this={canvas} height={pxHeight} width={pxWidth}></canvas>
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
