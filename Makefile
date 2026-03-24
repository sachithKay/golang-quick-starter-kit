# Run the application (Local Dev)
run:
	go run cmd/server/main.go

# Run all tests
test:
	go test ./...

# Build the production binary
build:
	go build -o bin/ghost cmd/server/main.go

# Clean up binaries
clean:
	rm -rf bin/

# Generate protobuf files
generate:
	buf generate