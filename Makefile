build:
	go build -o uinit main.go
run:
	./uinit init examples/processes.yaml
clean:
	rm -f uinit
