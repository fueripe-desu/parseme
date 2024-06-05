build:
	@go build -o bin/parseme cmd/parseme/main.go

run: build
	@./bin/parseme

test:
	@go test ./... -v
