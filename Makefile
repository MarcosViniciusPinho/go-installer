.PHONY: help build

help:
	@echo "Choose one of the following commands:"
	@echo "  -> make build to linux arch amd64	- Create the binary for Linux with amd64 architecture"
	@echo "  -> make help                		- Display this help message"

build: clean
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/goinstaller ./cmd/main.go

clean:
	rm -rf ./out