# Users service
[![Coverage Status](https://coveralls.io/repos/github/The-Psyducks/users-service/badge.svg?branch=main)](https://coveralls.io/github/The-Psyducks/users-service?branch=main)

## table of Contents

- [Information](#information)
- [Pre-requisites](#pre-requisites)
- [How to run](#how-to-run)
- [Tests](#tests)

## Information

Microservice responsible for user's user-related operation. It handles user creation, authentication, profile updates, and any other user-specific functionality.

## Pre-Requisites

To set up the project’s development environment, you will need to complete the .env.template in `database/` and in `server/`. Then the following is required: 

### Docker Requirements
The project runs entirely in Docker, so to start the environment, you will need:

Docker Engine: Minimum recommended version 19.x
Docker Compose: Minimum recommended version 1.27

### Local Development Requirements
If you prefer to set up and run the project locally outside of Docker, make sure you have:

Go: Version 1.23.0
Package manager: Go Modules (included with Go 1.23.0)

## How to run

Both the persistent database (Postgres) and the server are dockerized. To run them, use the following commands:

```
docker compose build
docker compose up service
```

To change the environment in which the application is running, you need to modify the ENVIRONMENT variable in the server/.env file.

## Tests

To run the tests, use the following commands:

```
docker compose build
docker compose up test
```

When running the tests, the program detects it and switches to an in-memory database to avoid modifying the development or production one.

This is the  [library](https://gin-gonic.com/docs/testing/)  used for testing.
