.PHONY: fmt lint test ci

## fmt: Auto-format code (goimports + gofmt).
fmt:
	golangci-lint run --fix --enable-only goimports ./...

## lint: Run golangci-lint.
lint:
	golangci-lint run ./...

## test: Run all unit tests with the race detector.
test:
	go test -race -count=1 -timeout=5m ./...

## ci: Run lint and test (mirrors CI pipeline).
ci: lint test
