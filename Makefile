.PHONY: swag build run RUN test init_sql cover mockgen

swag:
	swag init -g cmd/main.go

build:
	docker-compose build

run:
	docker-compose up

RUN: swag mockgen test build run

test:
	go test -v -count=1 ./...
	rm -rf pkg/repository/reports

init_sql:
	# Password for user postgres: postgres
	psql -U postgres -h 127.0.0.1 -p 5434 -d avito -f chema/init.sql

cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out
	rm -rf pkg/repository/reports

mockgen:
	mockgen -source=pkg/service/service.go -destination=pkg/service/mocks/mock.go