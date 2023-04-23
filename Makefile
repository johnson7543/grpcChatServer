build:
	go build -o bin/github.com/johnson7543/grpcChatServer

run: build
	./bin/github.com/johnson7543/grpcChatServer

proto:
	protoc --go-grpc_out=require_unimplemented_servers=false:chatserver \
	--go_out=chatserver ./chat.proto

.PHONY: proto