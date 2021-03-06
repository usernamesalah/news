# NEWS API DOCUMENTATION 

This repository contains source code for NEWS API.

## Getting Started

To run the project localy, make sure minimum requirements are fulfilled.

- Go version 1.10 or higher
- PostgreSQL version 12
- Docker (optional) -- see [here](https://docs.docker.com/get-docker/).

### Running in Local Machine

1. Make sure Go is installed as global command (first time only)

2. Clone this project and go to the root project to install all dependencies (first time only)
    ```bash
    // clone the repository
    > git clone git@github.com:usernamesalah/news.git
    ```

### Running from Docker Container

1. Make sure Docker and Docker Compose is installed

2. Run `docker-compose up` While still in root project for Running Kafka , PostgreSQL , zookeper and Elastic Search with docker

## API Documentation

We use [swag](https://github.com/swaggo/swag) to generate necearry Swagger files for API documentation. Everytime we run `make build`, the Swagger documentation will be updated.

To configure your API documentation, please refer to [swag's declarative comments format](https://github.com/swaggo/swag#declarative-comments-format) and [examples](https://github.com/swaggo/swag#examples).

To access the documentation, please visit [API DOCUMENTATION](http://localhost:8080/docs/api/v1/index.html).


## License
[MIT](https://choosealicense.com/licenses/mit/)