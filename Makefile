build:
	go build -o cmd/main ./cmd

run:
	go build -o cmd/main ./cmd && ./cmd/main

update:
	go get -u ./...
	go mod tidy