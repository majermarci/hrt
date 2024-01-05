PACKAGE := ./cmd/hrt
BINARY_PATH := ./bin/hrt

tidy:
	go fmt ./...
	go mod tidy -v

audit:
	go mod verify
	go vet ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all ./...
	go test -vet=off -buildvcs -race ./...

build:
	# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -a -ldflags "-s -w -X main.commitID=$(shell git rev-parse --short HEAD)" -installsuffix cgo -o=${BINARY_PATH}_linux_amd64 ${PACKAGE}
	CGO_ENABLED=0 go build -trimpath -a -ldflags "-s -w -X main.commitID=$(shell git rev-parse --short HEAD)" -installsuffix cgo -o=${BINARY_PATH} ${PACKAGE}
