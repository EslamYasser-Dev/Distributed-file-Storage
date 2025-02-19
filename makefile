## Variables
BINARY_NAME = DistrubutedFileSystem
BINARY_PATH = bin/$(BINARY_NAME)

## Build the Go program
build:
	@go build -o $(BINARY_PATH)

## Clear the terminal and remove the binary
clear:
	@clear && rm -rf bin/

## Run the program after building and clearing the terminal
run: clear build
	@$(BINARY_PATH)

## Run tests with verbose output
test:
	@go test ./... -v

## Fix formatting issues
fix:
	@gofmt -w .

## Clean up
clean: clear
	@go clean