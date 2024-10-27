build:
	@go build -o bin/listjobs
run: build
	@./bin/listjobs

test:
	@go test -v ./...