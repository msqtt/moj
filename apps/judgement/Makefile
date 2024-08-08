protoc-judge: 
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	 rpc/judgement.proto
	
protoc-sb: 
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	 rpc/sb-judger/sb-judger.proto

protoc-question: 
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	 rpc/question/question.proto

	
wire:
	wire .
	
dev:
	DEBUG=true go run .

.PHONY: protoc wire dev
