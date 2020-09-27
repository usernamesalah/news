all: init test build

.PHONY: init
init:
	@echo "> Installing the server dependencies ..."
	@go mod tidy -v
	@go get -v ./...

.PHONY: test
test:
	@echo "> Testing the server source code ..."
	@go test -cover -covermode atomic -coverprofile cover.out -race ./...
	@go tool cover -func cover.out

.PHONY: build
build:
	@echo "> Building the server binary ..."
	@rm -rf bin && go build -o bin/consumer .

push:
	@git add .
	@git commit -m "update"
	@git push origin master