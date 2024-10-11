build:
	@go build -o ./bin/go-thepapucoin

run: build
	@./bin/go-thepapucoin

test:
	go test -v ./...
