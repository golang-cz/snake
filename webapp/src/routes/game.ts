import {
	Direction,
	WebrpcError,
	type JoinGameReturn,
	type SnakeGame,
	type State
} from '$lib/rpc.gen';

const arrowMap: Record<string, Direction> = {
	ArrowUp: Direction.up,
	ArrowDown: Direction.down,
	ArrowLeft: Direction.left,
	ArrowRight: Direction.right
};

interface GameProps {
	ctx: CanvasRenderingContext2D;
	api: SnakeGame;
	username: string;
}

export class Game {
	private cellSize = 10;
	private borderColor = '#eee';

	private ctx: CanvasRenderingContext2D;
	private api: SnakeGame;
	private snakeId: number | null = null;
	private state: State | null = null;
	private username: string;

	constructor({ ctx, api, username }: GameProps) {
		this.ctx = ctx;
		this.api = api;
		this.username = username;
	}

	start = async () => {
		window.addEventListener('keydown', this.handleKeyDown);
		this.api.joinGame({
			onMessage: this.onMessageHandler,
			onError: this.onErrorHandler
		});
		({ snakeId: this.snakeId } = await this.api.createSnake({
			username: this.username
		}));
	};

	drawSquare = (x: number, y: number, color: string, border = false) => {
		this.ctx.fillStyle = color;
		this.ctx.fillRect(x * this.cellSize, y * this.cellSize, this.cellSize, this.cellSize);
		if (border) {
			this.ctx.strokeStyle = this.borderColor;
			this.ctx.strokeRect(x * this.cellSize, y * this.cellSize, this.cellSize, this.cellSize);
		}
	};

	drawSnakes = (snakes: State['snakes']) => {
		for (const snake of Object.values(snakes)) {
			for (let i = 0; i < snake.body.length; i++) {
				this.drawSquare(snake.body[i].x, snake.body[i].y, snake.color);
			}
		}
	};

	drawItems = (items: State['items']) => {
		for (const item of Object.values(items)) {
			this.drawSquare(item.body.x, item.body.y, item.color);
		}
	};

	drawGrid = (height: number, width: number) => {
		this.ctx.fillStyle = '#fff';
		this.ctx.strokeStyle = this.borderColor;

		this.ctx.fillRect(0, 0, height, width);

		for (let x = 0.5; x < width; x += this.cellSize) {
			this.ctx.moveTo(x, 0);
			this.ctx.lineTo(x, height);
		}
		for (let y = 0.5; y < height; y += this.cellSize) {
			this.ctx.moveTo(0, y);
			this.ctx.lineTo(width, y);
		}

		this.ctx.stroke();
	};

	getChangedSquares = (oldState: State, newState: State) => {
		type Square = { x: number; y: number; color: string; border?: boolean };
		const changedSquares: Map<string, Square> = new Map();

		const allEntityIds = new Set([
			...Object.keys(oldState.snakes),
			...Object.keys(newState.snakes),
			...Object.keys(oldState.items),
			...Object.keys(newState.items)
		]);

		for (const id of allEntityIds) {
			const oldSnake = oldState.snakes[id];
			const newSnake = newState.snakes[id];
			const oldItem = oldState.items[id];
			const newItem = newState.items[id];

			if (oldSnake && newSnake) {
				const maxLength = Math.max(oldSnake.body.length, newSnake.body.length);
				for (let i = 0; i < maxLength; i++) {
					const oldBody = oldSnake.body[i];
					const newBody = newSnake.body[i];

					if (newBody && (!oldBody || oldBody.x !== newBody.x || oldBody.y !== newBody.y)) {
						changedSquares.set(`${newBody.x},${newBody.y}`, {
							x: newBody.x,
							y: newBody.y,
							color: newSnake.color
						});
					}

					if (oldBody && (!newBody || oldBody.x !== newBody.x || oldBody.y !== newBody.y)) {
						const key = `${oldBody.x},${oldBody.y}`;
						if (!changedSquares.has(key)) {
							changedSquares.set(key, {
								x: oldBody.x,
								y: oldBody.y,
								color: 'white',
								border: true
							});
						}
					}
				}
			} else if (newSnake) {
				for (let i = 0; i < newSnake.body.length; i++) {
					const { x, y } = newSnake.body[i];
					changedSquares.set(`${x},${y}`, { x, y, color: newSnake.color });
				}
			} else if (oldSnake) {
				for (let i = 0; i < oldSnake.body.length; i++) {
					const { x, y } = oldSnake.body[i];
					const key = `${x},${y}`;
					if (!changedSquares.has(key)) {
						changedSquares.set(key, { x, y, color: 'white', border: true });
					}
				}
			}

			if (oldItem && newItem) {
				if (oldItem.body.x !== newItem.body.x || oldItem.body.y !== newItem.body.y) {
					changedSquares.set(`${newItem.body.x},${newItem.body.y}`, {
						x: newItem.body.x,
						y: newItem.body.y,
						color: newItem.color
					});

					const key = `${oldItem.body.x},${oldItem.body.y}`;
					if (!changedSquares.has(key)) {
						changedSquares.set(key, {
							x: oldItem.body.x,
							y: oldItem.body.y,
							color: 'white',
							border: true
						});
					}
				}
			} else if (newItem) {
				const { x, y } = newItem.body;
				changedSquares.set(`${x},${y}`, { x, y, color: newItem.color });
			} else if (oldItem) {
				const { x, y } = oldItem.body;
				const key = `${x},${y}`;
				if (!changedSquares.has(key)) {
					changedSquares.set(key, { x, y, color: 'white', border: true });
				}
			}
		}

		return Array.from(changedSquares.values());
	};

	onMessageHandler = (message: JoinGameReturn) => {
		const { width, height, snakes, items } = message.state;
		if (!this.state) {
			this.state = message.state;
			const pxWidth = width * this.cellSize;
			const pxHeight = height * this.cellSize;
			this.ctx.canvas.width = pxWidth;
			this.ctx.canvas.height = pxHeight;
			this.drawGrid(pxWidth, pxHeight);
			this.drawSnakes(snakes);
			this.drawItems(items);
		} else {
			const changedSquares = this.getChangedSquares(this.state, message.state);
			for (const square of changedSquares) {
				this.drawSquare(square.x, square.y, square.color, square.border);
			}
			this.state = message.state;
		}
	};

	onErrorHandler = (error: WebrpcError) => {
		console.log(error);
	};

	handleKeyDown = (e: KeyboardEvent) => {
		const key = e.key;
		if (key in arrowMap && this.snakeId) {
			e.preventDefault();
			this.api.turnSnake({ snakeId: this.snakeId, direction: arrowMap[key] });
		}
	};
}
