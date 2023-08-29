.PHONY: swag build run test init_sql cover mockgen

swag:
	swag init -g cmd/main.go

build:
	docker-compose build

run:
	docker-compose up

test:
	go test -v -count=1 ./...

init_sql:
	# Password for user postgres: postgres
	psql -U postgres -h 127.0.0.1 -p 5434 -d avito -f chema/init.sql

cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

mockgen:
	mockgen -source=pkg/service/service.go -destination=pkg/service/mocks/mock.go