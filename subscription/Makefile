generate:
	mkdir -p pkg/register_v1
	protoc --go_out=./pkg/register_v1 \
    	--go-grpc_out=./pkg/register_v1 \
    	api/proto/service.proto

install-deps:
	go install github.com/jackc/tern/v2@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.35.1
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3
	brew install protobuf

lint:
	golangci-lint run

migrate:
	tern migrate --migrations internal/db/migrations --config internal/db/migrations/tern.conf

docker-up:
	docker-compose -f build/docker-compose.yml up -d

docker-down:
	docker-compose -f build/docker-compose.yml down --remove-orphans

docker-down-clear:
	docker-compose -f build/docker-compose.yml down -v --remove-orphans

run-local:
	go run cmd/service/main.go --config config/local.overwrite.yaml

run-prod:
	go run cmd/service/main.go --config config/prod.overwrite.yaml