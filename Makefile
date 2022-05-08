BINARY_NAME=Tunnel
all: clean windows linux mac
clean:
	rm -f bin/*
windows:
	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME).exe -v
linux:
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)4Linux -v
mac:
	GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY_NAME)4Mac -v
