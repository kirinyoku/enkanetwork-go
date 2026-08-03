.PHONY: fmt lint test ci integration integration-genshin integration-hsr integration-zzz integration-enka

## fmt: Auto-format code.
fmt:
	golangci-lint run --fix ./...

## lint: Run golangci-lint.
lint:
	golangci-lint run ./...

## test: Run all unit tests with the race detector.
test:
	go test -race -count=1 -timeout=5m ./...

## ci: Run lint and test (mirrors CI pipeline).
ci: lint test

## integration: Run all integration tests against the live API.
integration: integration-genshin integration-hsr integration-zzz integration-endfield integration-enka

## integration-genshin: Run Genshin Impact integration tests.
integration-genshin:
	RUN_INTEGRATION_TESTS=true go test -v -race -timeout=5m -tags=integration ./client/genshin/...

## integration-hsr: Run Honkai: Star Rail integration tests.
integration-hsr:
	RUN_INTEGRATION_TESTS=true go test -v -race -timeout=5m -tags=integration ./client/hsr/...

## integration-zzz: Run Zenless Zone Zero integration tests.
integration-zzz:
	RUN_INTEGRATION_TESTS=true go test -v -race -timeout=5m -tags=integration ./client/zzz/...

## integration-endfield: Run Arknights Endfield integration tests.
integration-endfield:
	RUN_INTEGRATION_TESTS=true go test -v -race -timeout=5m -tags=integration ./client/endfield/...

## integration-enka: Run EnkaNetwork profile integration tests.
integration-enka:
	RUN_INTEGRATION_TESTS=true go test -v -race -timeout=5m -tags=integration ./client/enka/...
