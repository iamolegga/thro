.PHONY: build
build:
	CGO_ENABLED=0 go build -o=./thro .

.PHONY: lint
lint:
	go vet ./...

.PHONY: test
test:
	go test ./...
