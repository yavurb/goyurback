latest_tag := `git describe --tags 2> /dev/null || git rev-parse --short HEAD`

run: write_version
	GO_ENV=dev go run cmd/goyurback/main.go

build: write_version
	go build -o bin/goyurback cmd/goyurback/main.go

dev: write_version
	go tool air -c .air.toml

docker_build: test write_version
	docker build . -t goyurback:{{latest_tag}}

test:
	go test -v ./...

new_migration name=`uuidgen`:
	go tool migrate create -ext sql -dir migrations {{name}}

db_upgrade conn_string="postgres://postgres:postgres@localhost:5432/goyurback?sslmode=disable":
	go tool migrate -source file://migrations -database {{conn_string}} up

write_version:
	@echo {{latest_tag}} > cmd/goyurback/.version
