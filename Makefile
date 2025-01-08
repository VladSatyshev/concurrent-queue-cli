.PHONY: build
build:
	go build -o ./build/cli ./cmd/main.go


.PHONY: test
test:
	go test ./...