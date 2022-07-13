compile:
	go build cmd/main.go

.PHONY: unit-test
unit-test:
	go test ./internal/... -count=1

.PHONY: integration-test
integration-test:
	go test ./tests/integration/... -count=1