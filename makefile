PACKAGE := ./cmd/hrt
BINARY_PATH := ./bin/hrt

# Format the code and tidy up the dependencies
tidy:
	@go fmt ./...
	@go mod tidy -v

# Verify dependencies and run checks + tests
audit:
	@go mod verify
	@go vet ./...
	@go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	@go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all ./...
	@go test -vet=off -race ./...

# Build the binary
build:
	@CGO_ENABLED=0 go build -trimpath -a \
	-ldflags "-s -w -X main.commitID=$(shell git rev-parse --short HEAD)" \
	-installsuffix cgo -o=${BINARY_PATH} ${PACKAGE}
