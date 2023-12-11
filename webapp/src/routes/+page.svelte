<script lang="ts">
	import { SnakeGame, WebrpcError, type JoinGameReturn, type State, Direction } from '$lib/rpc.gen';
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
		Object.values(state.snakes).forEach((snake) => {
			for (let i = 0; i < snake.body.length; i++) {
				drawSquare(snake.body[i].x, snake.body[i].y, snake.color);
			}
		});
	}
	function drawItems() {
		Object.values(state.items).forEach((item) => {
			drawSquare(item.body.x, item.body.y, item.color);
		});
	}
	function drawSquare(x: number, y: number, color: string) {
		ctx.fillStyle = color;
		ctx.fillRect(x * cellSize, y * cellSize, cellSize, cellSize);
	}

	let id: number;

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

		let resp = await api.createSnake({
			username: 'test'
		});

		id = resp.snakeId;

		api.joinGame({
			onMessage: (resp: JoinGameReturn) => {
				state = resp.state;
				drawGrid();
				drawSnakes();
				drawItems();
			},
			onError: (error: WebrpcError) => {
				// TODO: reconnect()
				//if (error instanceof WebrpcStreamLostError) {
				setInterval(() => {
					location.reload();
				}, 100);
				//}

				console.error('onError()', error);
				if (error.message == 'AbortError') {
					//log.value = [...log.value, { type: 'warn', log: 'Connection closed by abort signal' }];
				} else {
					//log.value = [...log.value, { type: 'error', log: String(error) }];
				}
			}
		});
	});

	const arrowMap: Record<string, Direction> = {
		ArrowUp: Direction.up,
		ArrowDown: Direction.down,
		ArrowLeft: Direction.left,
		ArrowRight: Direction.right
	};

	const handleKeyDown = (e: KeyboardEvent) => {
		const key = e.key;
		if (key in arrowMap) {
			e.preventDefault();
			api.turnSnake({ snakeId: id, direction: arrowMap[key] });
		}
	};
</script>

<svelte:window on:keydown={handleKeyDown} />

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
