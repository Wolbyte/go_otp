# GO OTP

A RESTful backend written in Go using Gin and GORM. The service provides
user registration and authentication via OTP, rate limiting per phone
number, and search and filtering capabilities with pagination. It is
designed with modularity in mind so features such as rate
limiting can be reused across multiple endpoints.

------------------------------------------------------------------------

## Features

-   **User registration with OTP**\
    Users register with a phone number and verify their identity using a
    one-time password (OTP).

-   **Login with OTP instead of password**\
    No password is stored; OTP acts as the authentication mechanism.

-   **Rate limiting middleware**\
    Requests are limited per phone number. If a user exceeds three
    requests within 10 minutes, the rate limit blocks further requests
    temporarily.

-   **Search and filtering**\
    Supports searching by fields such as phone number or registration date.
    Date filters allow precise querying with ranges.

-   **Pagination**\
    The main user listing endpoint support pagination via `page` and `page_size` query parameters.

-   **Configurable environment**\
    Environment variables are managed via a `.env` file. A `.env.example` file is provided for convenience.

------------------------------------------------------------------------

## Prerequisites

-   Go 1.21+
-   PostgreSQL 16+
-   Docker (optional, if you prefer running with containers)

------------------------------------------------------------------------

> [!IMPORTANT]
> Remember to adjust the values in the .env file (database credentials, JWT secret, etc.).

## Local Setup (without Docker)

1.  Clone the repository:

    ``` bash
    git clone https://github.com/wolbyte/go_otp.git
    cd go_otp
    ```

2.  Create your environment file:

    ``` bash
    cp .env.example .env
    ```

3.  Run the server (it will automatically migrate your DB):

    ``` bash
    go run main.go
    ```

4.  The API will be available at:

        http://localhost:8080

------------------------------------------------------------------------

## Setup with Docker

1.  Clone the repository:

    ``` bash
    git clone https://github.com/wolbyte/go_otp.git
    cd go_otp
    ```

2.  Create your environment file:

    ``` bash
    cp .env.example .env
    ```

3.  Build and start the containers:

    ``` bash
    docker-compose up --build
    ```

    -   The Postgres database will run as a service alongside the API.
    -   The ports will be exposed so you can access the backend from your host machine easily

4.  Access the API at:

        http://localhost:8080

------------------------------------------------------------------------

## API Documentation
> [!CAUTION]
> If you set `GIN_MODE=release`, the swagger docs will not be hosted as a security measure

This project uses [swag](https://github.com/swaggo/swag) to automatically generate documentation from annotations in the code.

If you apply any changes to the docs you can re-build them again with `swag init` (follow the upstream installation instructions)

Interactive Swagger documentation with examples is available at: 

    http://localhost:8080/swagger/index.html

------------------------------------------------------------------------

## Why PostgreSQL?

PostgreSQL was chosen as the primary database for several reasons:

-   **Reliability and ACID compliance**\
    PostgreSQL guarantees data consistency and durability, which are critical for important workflows.

-   **Advanced query support**\
    Features like `ILIKE` for case-insensitive search and flexible date filtering simplify implementing robust search and filtering.

-   **Maturity and ecosystem**\
    With a long history of production use and strong community support, PostgreSQL is a safe and future-proof choice.

-   **Extensibility**\
    Support for JSONB and advanced indexing makes PostgreSQL suitable for scaling the project as new features are introduced.

------------------------------------------------------------------------
