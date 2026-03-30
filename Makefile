main_package_path = .
binary_name = chirpy
bin_dir = /tmp/bin

db_url = postgres://postgres:superpswd@localhost:5433/chirpy?sslmode=disable
migration_dir = sql/schema

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## test: run tests
.PHONY: test
test:
	go test -v ./...

## audit: run quality checks
.PHONY: audit
audit: test
	go mod tidy -diff
	go vet ./...
	test -z "$$(gofmt -l .)"

## build: build the binary
.PHONY: build
build:
	mkdir -p ${bin_dir}
	go build -o ${bin_dir}/${binary_name} ${main_package_path}

## run: build and run the app, pass args with `make run args="register toto"`
.PHONY: run
run: build
	${bin_dir}/${binary_name} ${args}

## dev: run app with live reload
.PHONY: dev
dev:
	air

## tidy: format and tidy
.PHONY: tidy
tidy:
	go mod tidy
	go fmt ./...

## mig-up: run goose migrations up
.PHONY: mig-up
mig-up:
	goose -dir ${migration_dir} postgres "${db_url}" up

## mig-down: run goose migrations down
.PHONY: mig-down
mig-down:
	goose -dir ${migration_dir} postgres "${db_url}" down

## mig-status: show goose migration status
.PHONY: mig-status
mig-status:
	goose -dir ${migration_dir} postgres "${db_url}" status

## sqlc: regenerate sqlc code
.PHONY: sqlc
sqlc:
	sqlc generate