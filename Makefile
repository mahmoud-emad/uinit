build:
	go build -o uinit main.go
run:
	./uinit init examples/services.yaml
clean:
	rm -f uinit
