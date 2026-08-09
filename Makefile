build:
	go build -o uinit main.go
run:
	./uinit examples/services.yaml
clean:
	rm -f uinit
