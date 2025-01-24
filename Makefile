
COVERAGE_DIR = ./core/coverage

test:
	go test -race -v ./...

coverage:
	mkdir -p ${COVERAGE_DIR}
	go test -v -coverprofile=${COVERAGE_DIR}/coverage.out -covermode=atomic ./...
	go tool cover -html=${COVERAGE_DIR}/coverage.out -o ${COVERAGE_DIR}/coverage.html

open:
	open ${COVERAGE_DIR}/coverage.html