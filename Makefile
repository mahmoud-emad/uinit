build:
	go build -o miniinit cmd/miniinit/main.go
run:
	./miniinit examples/services.yaml
clean:
	rm -f miniinit
