.PHONY: build
build:
	@CGO_ENABLED=0 go build -v ./

.PHONY: fmt
fmt:
	@find . -iname "*.go" -not -path "./vendor/**" | xargs gofmt -s -w
