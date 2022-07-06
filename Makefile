all: tidy test

tidy:
	@echo "Tidy up..."
	@go mod tidy -v

test:
	@echo "Running tests..."
	@go test -cover ./...
