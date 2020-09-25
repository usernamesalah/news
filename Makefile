API_DOCS_PATH = api/v1/docs
all: init test build

.PHONY: init
init:
	@echo "> Installing the server dependencies ..."
	@go mod tidy -v
	@go get -v ./...
	@go install github.com/swaggo/swag/cmd/swag

.PHONY: test
test:
	@echo "> Testing the server source code ..."
	@go test -cover -covermode atomic -coverprofile cover.out -race ./...
	@go tool cover -func cover.out

.PHONY: build
build: gen-swagger
	@echo "> Building the server binary ..."
	@rm -rf bin && go build -o bin/news .

.PHONY: gen-swagger
gen-swagger:
	@echo "Updating API documentation..."
	@rm -rf ${API_DOCS_PATH}
	@swag init -o ${API_DOCS_PATH}
