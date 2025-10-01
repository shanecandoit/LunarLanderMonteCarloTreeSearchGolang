
set -e

go mod tidy

echo "Building the project..."
go fmt
go build

echo "Running tests..."
go test ./...

ls -lh
