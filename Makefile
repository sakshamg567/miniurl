export PATH := /usr/local/go/bin:$(PATH)


gen: 
	protoc --go_out=. --go-grpc_out=. shared/proto/redirect.proto 
	protoc --go_out=. --go-grpc_out=. shared/proto/shortener.proto 
	protoc --go_out=. --go-grpc_out=. shared/proto/token.proto 

build: 
	go build -o bin/server ./cmd/server
	go build -o bin/gateway ./cmd/gateway

docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

run-server: 
	./bin/server &

run-gateway: 
	./bin/gateway &

clean: 
	rm -rf bin/
	docker-compose down -v

