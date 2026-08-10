build:
	go build -o uinit main.go
run:
	./uinit init examples/processes.yaml
test:
	go test -race ./...
lint:
	@test -z "$$(gofmt -l .)" || (echo "these files need gofmt:"; gofmt -l .; exit 1)
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@v0.7.0 ./...
	go run github.com/kisielk/errcheck@v1.20.0 ./...
dist:
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -o dist/uinit-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build -o dist/uinit-linux-arm64 main.go
	GOOS=darwin GOARCH=amd64 go build -o dist/uinit-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o dist/uinit-darwin-arm64 main.go
clean:
	rm -f uinit
	rm -rf dist
