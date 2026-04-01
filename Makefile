build:
	CGO_ENABLED=0 go build -o cutoff ./cmd/cutoff/

run: build
	./cutoff

test:
	go test ./...

clean:
	rm -f cutoff

.PHONY: build run test clean
