gen: 
	protoc --go_out=. --go-grpc_out=. shared/proto/redirect.proto 
	protoc --go_out=. --go-grpc_out=. shared/proto/shortener.proto 