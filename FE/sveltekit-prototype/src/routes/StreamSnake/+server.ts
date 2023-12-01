const board = {
	snakes: [
		{
			id: 1,
			body: [
				{ x: 0, y: 0 },
				{ x: 0, y: 1 },
				{ x: 0, y: 2 }
			]
		}
	],
	food: [
		{ x: 1, y: 0 },
		{ x: 1, y: 1 }
	]
};

export function POST() {
	let interval: any;
	const stream = new ReadableStream({
		start(controller) {
			controller.enqueue(JSON.stringify({ board }) + '\n');

			interval = setInterval(() => {
				let string = '{"time":' + Date.now() + ',"value":' + Math.random() + '}\n';
				// Add the string to the stream
				controller.enqueue(string);
			}, 1000);
		},
		cancel() {
			clearInterval(interval);
		}
	});
	return new Response(stream);
}
