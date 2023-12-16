all:
	@awk -F'[ :]' '!/^all:/ && /^([A-z_-]+):/ {print "make " $$1}' Makefile

generate:
	webrpc-gen -schema=proto/snake.ridl -target=../../Webrpc/gen-golang -pkg=proto -server -client -fmt=false -out=proto/snake.gen.go
	webrpc-gen -schema=proto/snake.ridl -target=../../Webrpc/gen-typescript -client -out=webapp/src/lib/rpc.gen.ts

run:
	go run ./

test:
	go test -v ./
