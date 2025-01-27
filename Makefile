build:
	@go build -o bin/goecomapi cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/goecomapi		
