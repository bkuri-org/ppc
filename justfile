default:
    just --list

build:
    go build -o build-prompt ./cmd/build-prompt

test:
    go test ./...

smoke:
    ./scripts/smoke.sh

lint:
    golangci-lint run ./...
