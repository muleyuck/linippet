.PHONY: test
test:
	golangci-lint run ./...
	go test -v ./...
	go build -v .

.PHONY: bench
bench:
	go test -bench=. -benchmem ./internal/fuzzy_search/
