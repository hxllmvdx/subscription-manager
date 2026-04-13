# Subscription Manager

REST API service for aggregating user online subscription data.

## Tech Stack

- **Go** (Gin framework)
- **PostgreSQL** (Database)
- **GORM** (ORM)
- **Goose** (Migrations)
- **Swagger** (API documentation)
- **Docker Compose** (Orchestration)

## Requirements

- Docker & Docker Compose
- Go 1.21+ (for local development)

## Quick Start

### Run with Docker Compose

```bash
# Copy environment configuration
cp .env.example .env

# Start services
docker-compose up --build
```

The server will be available at `http://localhost:8080`

Swagger documentation: `http://localhost:8080/swagger/index.html`

### Local Development

```bash
# Install dependencies
cd backend
go mod tidy

# Generate Swagger documentation
go run github.com/swaggo/swag/cmd/swag@v1.8.12 init -g cmd/server/main.go -o docs --parseDependency --parseInternal

# Start PostgreSQL (if not already running)
docker-compose up -d postgres

# Run the server
BACKEND_PORT=:8080 DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres DB_NAME=subscriptions go run cmd/server/main.go
```

## API Endpoints

### Subscriptions

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/subscriptions` | Create a subscription |
| GET | `/api/subscriptions` | List all subscriptions (optional: `?user_id=...`) |
| GET | `/api/subscriptions/:id` | Get subscription by ID |
| PUT | `/api/subscriptions/:id` | Update a subscription |
| DELETE | `/api/subscriptions/:id` | Delete a subscription |
| GET | `/api/calculate_total_price` | Calculate total subscription cost |

### Request Examples

#### Create a subscription

```bash
curl -X POST http://localhost:8080/api/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
  }'
```

#### List user subscriptions

```bash
curl http://localhost:8080/api/subscriptions?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba
```

#### Calculate total cost

```bash
curl "http://localhost:8080/api/calculate_total_price?period_start=01-2025&period_end=12-2025&user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba"
```

#### Update a subscription

```bash
curl -X PUT http://localhost:8080/api/subscriptions/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "price": 500,
    "end_date": "12-2025"
  }'
```

## Project Structure

```
backend/
├── cmd/server/main.go          # Entry point
├── internal/
│   ├── config/                 # Configuration
│   ├── domain/                 # Domain models
│   ├── handler/                # HTTP handlers
│   ├── migrate/                # Database migrations
│   │   └── sql/                # SQL migration files
│   ├── repository/             # Database access layer
│   └── routes/                 # Route definitions
├── docs/                       # Swagger documentation
└── Dockerfile
```

## Configuration

All configuration is managed via the `.env` file:

| Variable | Description | Default |
|----------|-------------|---------|
| `BACKEND_PORT` | Server port | `:8080` |
| `DB_HOST` | Database host | `postgres` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `postgres` |
| `DB_NAME` | Database name | `subscriptions` |

## Migrations

Migrations are applied automatically on server startup using the `goose` library.

Migration files are located in `backend/internal/migrate/sql/`.

## Logging

The service uses the standard `log/slog` package for structured logging. Logs include:

- Server startup events
- Subscription CRUD operations
- Validation and database errors
- API request metrics
