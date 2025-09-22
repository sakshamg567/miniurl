# MiniURL - Microservices URL Shortener

A production-ready URL shortener built with Go, featuring a microservices architecture designed for scalability and performance. The system uses gRPC for inter-service communication, Redis for caching, and PostgreSQL for persistent storage.

## Project Overview

MiniURL is a distributed URL shortening service that follows microservices principles. It separates concerns into distinct services that communicate via gRPC, allowing for independent scaling and deployment of each component. The system provides fast URL shortening and redirection with built-in caching for optimal performance.

## Architecture

The system consists of three core microservices:

1. **Token Service** - Generates unique IDs using a Redis-backed counter
2. **Shortener Service** - Creates short codes from long URLs using base62 encoding
3. **Redirect Service** - Handles URL resolution with cache-first lookup
4. **HTTP Gateway** - Provides REST API endpoints for client interaction

```
┌─────────────┐    ┌──────────────┐    ┌─────────────┐
│   Gateway   │────│  Shortener   │────│   Token     │
│   (HTTP)    │    │   Service    │    │  Service    │
│   :8080     │    │   :50052     │    │   :50051    │
└─────────────┘    └──────────────┘    └─────────────┘
       │                   │
       │           ┌───────────────┐
       └───────────│   Redirect    │
                   │   Service     │
                   │   :50053      │
                   └───────────────┘
```

## Project Structure

```
miniurl/
├── cache/                    # Redis cache operations
│   └── cache.go
├── cmd/
│   ├── gateway/             # HTTP API gateway
│   │   └── main.go
│   └── server/              # Service orchestrator
│       └── main.go
├── db/                      # PostgreSQL database operations
│   └── db.go
├── services/
│   ├── redirect/            # URL redirection service
│   │   └── server.go
│   ├── shortener/           # URL shortening service
│   │   └── server.go
│   └── token/               # Unique ID generation service
│       └── server.go
├── shared/
│   └── proto/               # Protocol Buffer definitions
│       ├── redirect.proto
│       ├── shortener.proto
│       ├── token.proto
│       └── *pb/             # Generated gRPC code
├── docker-compose.yml       # Infrastructure setup
├── Makefile                # Build and development commands
└── .env.local              # Local environment configuration
```

## Component Breakdown

### Token Service (Port: 50051)
- Manages a Redis-backed counter for generating unique IDs
- Initializes counter at 100,000,000,000 to ensure minimum short code length
- Provides atomic ID generation for the shortener service

### Shortener Service (Port: 50052)
- Accepts long URLs and generates corresponding short codes
- Uses base62 encoding for compact URL representation
- Stores URL mappings in PostgreSQL database
- Communicates with Token Service for unique ID generation

### Redirect Service (Port: 50053)
- Handles short code to long URL resolution
- Implements cache-first lookup strategy (Redis → PostgreSQL)
- Caches successful lookups for 24 hours
- Returns original URLs for redirection

### HTTP Gateway (Port: 8080)
- Provides REST API endpoints for external clients
- **POST /shorten** - Create short URL from long URL
- **GET /{shortCode}** - Redirect to original URL
- Translates HTTP requests to gRPC calls

### Database Layer
- **PostgreSQL**: Persistent storage for URL mappings with timestamps
- **Redis**: High-speed cache for frequently accessed URLs and ID counter

## Prerequisites

- Go 1.24.1 or later
- Docker and Docker Compose
- Protocol Buffers compiler (protoc)
- Make utility

## Installation & Setup

1. **Clone the repository**
```bash
git clone https://github.com/sakshamg567/miniurl.git
cd miniurl
```

2. **Start infrastructure services**
```bash
make docker-up
```
This starts PostgreSQL (port 5433) and Redis (port 6380) containers.

3. **Generate Protocol Buffer code**
```bash
make gen
```

4. **Build the services**
```bash
make build
```

## Running the Application

### Method 1: Using Make commands
```bash
# Start all services
make run-server    # Starts all microservices
make run-gateway   # Starts HTTP gateway
```

### Method 2: Manual execution
```bash
# Start microservices
./bin/server &

# Start HTTP gateway
./bin/gateway &
```

## API Usage

### Shorten a URL
```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"long_url": "https://example.com/very/long/url"}'
```

Response:
```json
{"short_code": "dBVXOx"}
```

### Access shortened URL
```bash
curl -L http://localhost:8080/dBVXOx
```
This will redirect to the original URL.

## Environment Configuration

The `.env.local` file contains database and cache connection strings:

```env
dsn="host=localhost user=postgres password=secret dbname=miniurl port=5433 sslmode=disable"
cache="localhost:6379"
```

## Development Commands

```bash
# Generate protobuf code
make gen

# Build binaries
make build

# Start infrastructure
make docker-up

# Stop infrastructure
make docker-down

# Clean build artifacts
make clean
```

## Key Features

- **Microservices Architecture**: Independent, scalable services
- **gRPC Communication**: High-performance inter-service communication
- **Caching Strategy**: Redis-first lookup for optimal performance
- **Atomic ID Generation**: Collision-free short code generation
- **Base62 Encoding**: Compact, URL-safe short codes
- **Database Persistence**: Reliable PostgreSQL storage
- **Docker Support**: Containerized infrastructure setup

## Performance Optimizations

- Cache-first lookup pattern reduces database load
- Base62 encoding creates short, readable URLs
- Redis atomic operations ensure unique ID generation
- Separate services allow independent scaling


