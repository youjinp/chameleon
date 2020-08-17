build:
	go build -o main

update:
	go get -u ./...
	go mod tidy