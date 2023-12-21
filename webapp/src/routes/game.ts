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

	drawSquare = (x: number, y: number, color: string) => {
		if (!this.ctx) return;
		this.ctx.fillStyle = color;
		this.ctx.fillRect(x * this.cellSize, y * this.cellSize, this.cellSize, this.cellSize);
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
		if (!this.ctx) return;
		this.ctx.fillStyle = '#fff';
		this.ctx.strokeStyle = '#eee';

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

	onMessageHandler = (message: JoinGameReturn) => {
		const { width, height, snakes, items } = message.state;
		// if (!this.state) {
		this.state = message.state;
		const pxWidth = width * this.cellSize;
		const pxHeight = height * this.cellSize;
		this.ctx.canvas.width = pxWidth;
		this.ctx.canvas.height = pxHeight;
		this.drawGrid(pxWidth, pxHeight);
		this.drawSnakes(snakes);
		this.drawItems(items);
		// } else {
		// }
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
