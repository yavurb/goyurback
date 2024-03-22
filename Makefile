latest_tag := $(shell git describe --tags || git rev-parse --short HEAD)

run: write_version
	go run cmd/goyurback/main.go

dev: write_version
	air -c .air.toml

docker_build: test write_version
	docker build . -t goyurback:$(latest_tag)

test:
	go test -v ./...

new_migration:
	migrate create -ext sql -dir migrations $(name)

write_version:
	@echo $(latest_tag) > cmd/goyurback/.version
