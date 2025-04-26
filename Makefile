GOBOK_BIN = ./cmd/gobok/gobok

all: gobok generate fmt

gobok:
	go build -o $(GOBOK_BIN) ./cmd/gobok

generate:
	$(GOBOK_BIN) ./pkg ./internal

fmt:
	go fmt ./...
	go mod tidy

clean:
	rm -f $(GOBOK_BIN)
	find . -name 'gobok.go' -delete
