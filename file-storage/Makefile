generate:
	mkdir -p internal/generated/schema
	oapi-codegen -package file_storage -generate chi-server,types,spec api/schema.yaml > internal/generated/schema/file_storage.gen.go

install-deps:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

lint:
	golangci-lint run

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