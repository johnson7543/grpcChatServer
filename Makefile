proto:
	protoc --go-grpc_out=require_unimplemented_servers=false:chatserver \
	--go_out=chatserver ./chat.proto

.PHONY: proto