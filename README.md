# Go Social App

A simple social networking backend API built with Go, providing endpoints for user management and posts.

## Features

- Health check endpoint
- Post creation with tags support
- User registration and management 
- PostgreSQL database integration
- Environment-based configuration
- Docker support for database
- Database migrations with golang-migrate

## Tech Stack

- **Language**: Go 1.25.1
- **Framework**: Chi v5 (HTTP router)
- **Database**: PostgreSQL
- **Libraries**:
  - `github.com/lib/pq` - PostgreSQL driver
  - `github.com/joho/godotenv` - Environment variables
  - `github.com/go-chi/chi/v5` - HTTP router and middleware

## Prerequisites

- Go 1.25.1 or later
- Docker and Docker Compose (for database)
- PostgreSQL (or use Docker)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/Saswat-Sagar-Sahu/Social.git
   cd go-social-app
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

## Environment Setup

Create a `.env` file in the root directory with the following variables:

```env
# Server configuration
ADDRESS=:8080

# Database configuration
DB_ADDR=postgres://admin:admin123@localhost:5432/social?sslmode=disable
DB_MAX_OPEN_CONNS=30
DB_MAX_IDLE_CONNS=30
DB_MAX_IDLE_TIME=15m
```

## Database Setup

### Using Docker Compose

1. Start the PostgreSQL database:
   ```bash
   docker-compose up -d
   ```

2. The database will be available at `localhost:5432` with:
   - User: `admin`
   - Password: `admin123`
   - Database: `social`

### Manual Setup

If you prefer to set up PostgreSQL manually, ensure you have a database running and update the `DB_ADDR` in your `.env` file accordingly.

## Database Migrations

This project uses [golang-migrate](https://github.com/golang-migrate/migrate) for database schema migrations.

### Install golang-migrate

```bash
# Using Homebrew (macOS)
brew install golang-migrate

# Or download from GitHub releases
# Visit: https://github.com/golang-migrate/migrate/releases
```

### Run Migrations

After setting up the database, run the migrations to create the necessary tables:

```bash
# Run all pending migrations
migrate -path=cmd/migrate/migrations -database="postgres://admin:admin123@localhost:5432/social?sslmode=disable" up

# Check migration status
migrate -path=cmd/migrate/migrations -database="postgres://admin:admin123@localhost:5432/social?sslmode=disable" version

# Rollback migrations if needed
migrate -path=cmd/migrate/migrations -database="postgres://admin:admin123@localhost:5432/social?sslmode=disable" down 1
```

**Note**: If you encounter path issues, use the absolute path:
```bash
migrate -path=/path/to/your/project/cmd/migrate/migrations -database="postgres://admin:admin123@localhost:5432/social?sslmode=disable" up
```


## Running the Application

You can run the API server directly:

```bash
go run cmd/api/*.go
```

Or use [Air](https://github.com/air-verse/air) for automatic live-reloading during development:

### Using Air for Live Reload

1. **Install Air** (if not already installed):
  ```bash
  go install github.com/air-verse/air@latest
  ```
2. **Start the development server with Air:**
  ```bash
  air
  ```
  Air will watch your Go files and automatically restart the server on code changes.

The server will start on the address specified in the `ADDRESS` environment variable (default: `:8080`).

### Frontend (React UI)

A React UI is available in the `web/` folder and includes:
- Register page
- Activate account page
- Login page
- Feed page (uses authenticated API request)

Run it from project root:

```bash
cd web
npm install
npm run dev
```

The Vite dev server proxies `/v1/*` requests to `http://localhost:8080`, so run the Go API and the React UI together during development.

### Seeding the database

Prerequisites: database is running and migrations applied. You can set the DB URL in `.env` or via `DB_ADDR`.

Quick (run without building):
```bash
go run ./cmd/seed -users=500 -posts=1000
```

## Swagger API Docs

This project uses `swag` (swaggo) to generate OpenAPI (Swagger) documentation from handler comments.

1. Install the `swag` CLI (if not already installed):

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

2. Generate the docs package (run from the project root):

```bash
swag init -g cmd/api/main.go -o docs
```

3. Run the API server:

```bash
go run ./cmd/api
```

4. Open the Swagger UI in your browser:

- Swagger UI: http://localhost:8080/swagger/index.html
- Raw Swagger JSON: http://localhost:8080/swagger/doc.json

Notes:
- Handlers must have swag annotations (e.g. `@Summary`, `@Param`, `@Success`) for endpoints to appear in the generated docs.
- If some models reference types like `sql.NullInt64`, prefer using simple exported types (e.g. `*int64`) or add type aliases so `swag` can generate schemas.
- You can programmatically set `docs.SwaggerInfo.Host` and `docs.SwaggerInfo.BasePath` in `cmd/api/main.go` after loading your config so the generated docs reflect the runtime host/basePath.


## API Endpoints

### Health Check
- **GET** `/v1/health`
  - Returns server status

### Posts
- **POST** `/v1/posts`
  - Create a new post
  - Request body: `{"title": "string", "content": "string"}`
  - Response: Created post object
- **PUT** `/v1/posts/{postId}`
  - Update a post by ID (owner or admin)
- **DELETE** `/v1/posts/{postId}`
  - Delete a post by ID (owner or admin)
- **GET** `/v1/posts/me`
  - Retrieve posts created by the authenticated user

### Users
- **POST** `/v1/users` (planned)
  - Create a new user

*Note: Post creation is now implemented. User creation endpoints are planned for future development.*

## Testing the API

### Health Check
```bash
curl http://localhost:8080/v1/health
```

## Development

### Project Structure

```
.
├── cmd/
│   ├── api/          # API server entry point
│   └── migrate/      # Database migration tools
├── internal/
│   ├── db/           # Database connection
│   ├── env/          # Environment utilities
│   └── store/        # Data access layer
├── scripts/          # Database initialization scripts
├── docker-compose.yml # Docker services
└── go.mod            # Go module file
```

### Building

To build the application:

```bash
go build -o bin/main cmd/api/*.go
```
