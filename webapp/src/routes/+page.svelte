<script lang="ts">
	import { SnakeGame, WebrpcError, type JoinGameReturn, type State } from '$lib/rpc.gen';
	import { onMount } from 'svelte';

	let canvas: HTMLCanvasElement;
	let ctx: CanvasRenderingContext2D;

	const api = new SnakeGame('http://localhost:5252', fetch);

	// Mock data
	let state: State = { width: 70, height: 70, snakes: [], items: [] };

	// Constants
	const cellSize = 10;
	const pxHeight = `${state.width * cellSize}px`;
	const pxWidth = `${state.height * cellSize}px`;

	function drawSnakes() {
		state.snakes.forEach((snake) => {
			for (let i = 0; i < snake.body.length; i++) {
				drawSquare(snake.body[i].x, snake.body[i].y, snake.color);
			}
		});
	}
	function drawItems() {
		state.items.forEach((item) => {
			drawSquare(item.body.x, item.body.y, item.color);
		});
	}
	function drawSquare(x: number, y: number, color: string) {
		ctx.fillStyle = color;
		ctx.fillRect(x * cellSize, y * cellSize, cellSize, cellSize);
	}

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

		api.joinGame({
			onMessage: (resp: JoinGameReturn) => {
				state = resp.state;
				drawGrid();
				drawSnakes();
				drawItems();

				console.log(resp);
			},
			onError: (error: WebrpcError) => {
				// TODO: reconnect()
				setInterval(() => {
					location.reload();
				}, 250);

				console.error('onError()', error);
				if (error.message == 'AbortError') {
					//log.value = [...log.value, { type: 'warn', log: 'Connection closed by abort signal' }];
				} else {
					//log.value = [...log.value, { type: 'error', log: String(error) }];
				}
			}
		});
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
