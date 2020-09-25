# SOCCER API DOCUMENTATION 

This repository contains source code for Soccer API.

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
    > git clone git@github.com:Apranta/technical-test.git

    // change directory to root project folder
    > cd soccer
    
    // install all the dependencies
    > make init   
    ```
3. Running your PostgreSQL
4. While still in root project build and run the app
    ```bash
    // build project
    > make build

    // source env
    > source .env.development
    
    // run project
    > ./bin/soccer

    // now go to http://localhost:8080/ in your browser to check the app.
    ```
5. If u wanna Run Unit Testing just run
    ```bash
        make test
    ```
### Running from Docker Container

1. Make sure Docker and Docker Compose is installed

2. Run `docker-compose up`

3. Build and run the app as described on the previous section.

## API Documentation

We use [swag](https://github.com/swaggo/swag) to generate necearry Swagger files for API documentation. Everytime we run `make build`, the Swagger documentation will be updated.

To configure your API documentation, please refer to [swag's declarative comments format](https://github.com/swaggo/swag#declarative-comments-format) and [examples](https://github.com/swaggo/swag#examples).

To access the documentation, please visit [API DOCUMENTATION](http://localhost:8080/docs/api/v1/index.html).


## License
[MIT](https://choosealicense.com/licenses/mit/)