protoc: 
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	 rpc/record.proto
	
wire:
	wire .
	
dev:
	DEBUG=true go run .

.PHONY: protoc wire dev
