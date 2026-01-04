# Go Social App

A simple social networking backend API built with Go, providing endpoints for user management and posts.

## Features

- User registration and management
- Post creation with tags
- Health check endpoint
- PostgreSQL database integration
- Environment-based configuration
- Docker support for database

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
ADDRESS=:8081

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

## Running the Application

Run the API server:

```bash
go run cmd/api/*.go
```

The server will start on the address specified in the `ADDRESS` environment variable (default: `:8081`).

## API Endpoints

### Health Check
- **GET** `/v1/health`
  - Returns server status

### Users
- **POST** `/v1/users` (planned)
  - Create a new user

### Posts
- **POST** `/v1/posts` (planned)
  - Create a new post

*Note: Currently, only the health check endpoint is implemented. User and post creation endpoints are planned for future development.*

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

