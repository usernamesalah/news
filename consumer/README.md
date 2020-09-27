# Kafka Consumer News API DOCUMENTATION 

This repository contains source code for Kafka Consumer News API.

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

    // change directory to root project folder
    > cd consumer
    
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
    > ./bin/consumer

    ```
5. If u wanna Run Unit Testing just run
    ```bash
        make test
    ```
