
set -e

go mod tidy

echo "Building the project..."
go fmt
# go build
go build -buildvcs=false

echo "Running tests..."
go test ./...

ls -lh
