<script lang="ts">
	import { onMount } from 'svelte';
	import { Chat } from './rpc';

	let canvas: HTMLCanvasElement;
	let ctx: CanvasRenderingContext2D;

	const api = new Chat('http://localhost:5174', fetch);

	//Mock data
	let snake = [
		{ x: 0, y: 0 },
		{ x: 1, y: 0 },
		{ x: 2, y: 0 },
		{ x: 3, y: 0 },
		{ x: 4, y: 0 }
	];

	let food = [{ x: 10, y: 10 }];

	const width = 70;
	const height = 70;

	//Constants
	const cellSize = 10;
	const pxHeight = `${width * cellSize}px`;
	const pxWidth = `${height * cellSize}px`;

	function drawSnake() {
		for (let i = 0; i < snake.length; i++) {
			drawSquare(snake[i].x, snake[i].y, 'green');
		}
	}
	function drawFood() {
		for (let i = 0; i < food.length; i++) {
			drawSquare(food[i].x, food[i].y, 'red');
		}
	}
	function drawSquare(x: number, y: number, color: string) {
		ctx.fillStyle = color;
		ctx.fillRect(x * cellSize, y * cellSize, cellSize, cellSize);
	}

	// api.streamSnake(
	// 	{},
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
		drawSnake();
		drawFood();
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
