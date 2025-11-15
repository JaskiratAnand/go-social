# GoSocial
This project is a social networking site built using Golang.

## Features

1. User Auth
2. User Profiles
3. Posts CRUD functions
4. Comments
5. Likes

## Tech Stack

### Backend

- Language: Golang
- Frameworks: Chi Router, Goose, SQLC
- Database: PostgreSQL
- Auth: JWT
- Libraries:
    - air-verse/air: live-reloading command line utility for developing Go
    - pq: PostgreSQL driver
    - godotenv: Environment variable management
    - goose: DB Migration
    - sqlc: DB Types & Queries Generation
    - go-playground/validator: input types validator

### Deployment

- Containerization: Docker
- Hosting: GCP
- CI/CD: GitHub Actions

## Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/JaskiratAnand/go-social.git
    cd go-social
    ```

2. Set up environment variables: Rename .env.example to .env:
    ```bash
    DB_ADDR="postgres://admin:adminpassword@localhost:5432/go-social?sslmode=disable"
    ```

3. Install dependencies:
    ```bash
    go mod tidy
    ```

4. Create a db container in docker
    ```bash
    docker compose up
    ```

5. Run database migrations using goose
    ```bash
    cd cmd/sql/schema

    goose postgres postgres://admin:adminpassword@localhost:5432/go-social up
    ```

6. Generate db types & Queries using sqlc (if needed)
    ```bash
    sqlc generate
    ```

7. Generate Swagger Docs
    ```bash
    swag init -g ./api/main.go -d cmd,internal && 
    swag fmt
    ```

8. Seed data to DB
    ```bash
    go run .\cmd\seed\main.go
    ```

9. Run the live server using air-verse/air
    ```bash
    air
    ```

10. Build
    ```bash
    # Windows
    go build -o server .\cmd\api

    # Linux 
    go build -o api cmd/api/*.go
    ```

Access the application:
The server will run at `http://localhost:8080` by default.


## API Endpoints

- Posts
    - POST `/v1/posts` - Create a new post.
    - GET `/v1/posts/{id}` - Fetch a specific post.
    - PATCH `/v1/posts/{id}` - Update a post.
    - DELETE `/v1/posts/{id}` - Delete a post.

- Comments
    - POST `/v1/posts/{postId}/comments` - Add a comment to a post.
    - GET `/v1/posts/{postId}/comments` - Get all comments for a post.
