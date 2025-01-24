
COVERAGE_DIR = ./core/coverage

install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go mod tidy

test:
	go test -race -v ./core/tests

coverage-ci:
	mkdir -p ${COVERAGE_DIR}
	go test -v -coverprofile=${COVERAGE_DIR}/coverage.out -covermode=atomic ./...
	go tool cover -html=${COVERAGE_DIR}/coverage.out -o ${COVERAGE_DIR}/coverage.html

coverage:
	@make coverage-ci
	open ${COVERAGE_DIR}/coverage.html

lint:
	golangci-lint run --timeout 60s --max-same-issues 50 ./...

lint-fix:
	golangci-lint run --timeout 60s --max-same-issues 50 --fix ./...
