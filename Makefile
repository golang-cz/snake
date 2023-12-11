all:
	@awk -F'[ :]' '!/^all:/ && /^([A-z_-]+):/ {print "make " $$1}' Makefile

generate:
	webrpc-gen -schema=proto/chat.ridl -target=../../gen-golang -pkg=proto -server -client -fmt=false -out=proto/chat.gen.go
	#webrpc-gen -schema=proto/chat.ridl -target=../../gen-typescript -client -out=webapp/src/rpc.gen.ts

run:
	go run ./

test:
	go test -v ./
