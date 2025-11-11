.PHONY: all test clean cli

LOG_PREFIX := [MAKEFILE BUILD]

test: 
	@echo -n "${LOG_PREFIX} Running tests: > "
	go test ./...
	@echo "${LOG_PREFIX} Tests complete."

cli: build-cli run-cli

run-cli: 
	@echo -n "${LOG_PREFIX} RUNNING CLI: > "
	./bin/cli

build-cli: 
	@echo -n "${LOG_PREFIX} Building CLI: > "
	go build -o bin/cli ./cmd/cli
	@echo "${LOG_PREFIX} Build complete."
