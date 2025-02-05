# Dealls Dating Backend

This is the backend system for a Dating Mobile App using Golang and PostgreSQL.

## Prerequisites

- Docker
- Docker Compose

## Getting Started

1. Clone the repository:

    ```bash
    git clone https://github.com/redzjovi/dealls-dating-backend.git
    cd deall
    ```

2. Create a `.env` file with your PostgreSQL credentials:

3. Run the application using Docker Compose:

    ```bash
    docker-compose up
    ```

4. Open your browser and go to `http://localhost:8080` to see the application running.

## Run Application

### Run unit test

```bash
go test ./... -coverprofile=./temp/coverage.out
go tool cover -html=./temp/coverage.out
```

## Other

- [DOCUMENT](DOCUMENT.md)

## License

This project is licensed under the MIT License.
