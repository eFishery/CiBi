# Get version from git hash
git_hash := $(shell git describe --tags)

# Get current date
current_time = $(shell date +"%Y-%m-%d:T%H:%M:%S")

# Add linker flags
linker_flags = '-s -w -X main.buildTime=${current_time} -X main.version=${git_hash}'

# Build binaries for current OS and Linux
.PHONY:

test:
	go test ./... -v

build:
	@echo "Building binaries..."
	go build -ldflags=${linker_flags} -o=./bin/binver ./main.go
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/linux_amd64/binver ./main.go