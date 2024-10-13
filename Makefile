build:
	@go build -o bin/auth main.go

run: build
	@./bin/auth