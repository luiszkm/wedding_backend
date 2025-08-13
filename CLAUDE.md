# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Building and Running
- `go build -o server ./cmd/api/main.go` - Build the API server binary
- `go run ./cmd/api/main.go` - Run the API server directly
- `docker-compose up --build` - Build and run the entire stack (API + PostgreSQL)
- `docker-compose up db` - Run only the PostgreSQL database

### Testing
- `go test ./...` - Run all tests in the project
- `go test ./internal/guest/domain` - Run tests for a specific package
- `go test -v ./...` - Run tests with verbose output
- `go test -run TestNewGrupoDeConvidados ./internal/guest/domain` - Run a specific test

### Development Tools
- `go mod tidy` - Clean up dependencies
- `go fmt ./...` - Format all Go files
- `go vet ./...` - Run Go static analysis
- `go mod download` - Download dependencies

## Architecture Overview

This is a Go-based wedding management API built with clean architecture principles and domain-driven design.

### Core Architecture Layers
Each domain module follows the same 4-layer structure:
- **Domain**: Business entities, value objects, and domain logic
- **Application**: Use cases and application services (orchestration)
- **Infrastructure**: External concerns (database repositories, external APIs)
- **Interfaces/REST**: HTTP handlers and DTOs

### Domain Modules
- **Guest**: Guest group management with RSVP functionality using access keys
- **Gift**: Gift registry with selection tracking, fractional gifts with quotas, and public/private views
- **MessageBoard**: Wedding messages with moderation capabilities
- **Gallery**: Photo management with labeling, favorites, and AWS S3/R2 storage
- **IAM**: User authentication using JWT tokens
- **Event**: Wedding event management
- **Billing**: Stripe integration for subscription plans
- **Communication**: Communication management with CRUD operations for announcements/comunicados
- **Itinerary**: Wedding itinerary and schedule management
- **PageTemplate**: Page template system for event customization

### Platform Services
- **Auth**: JWT token generation and middleware (`internal/platform/auth/`)
- **Storage**: File upload interface with R2/S3 implementation (`internal/platform/storage/`)
- **Template**: Template engine for page rendering (`internal/platform/template/`)
- **Web**: Common web utilities (`internal/platform/web/`)

### Key Design Patterns
- **Dependency Injection**: All dependencies are injected in `main.go` during application startup
- **Repository Pattern**: Each domain has a repository interface implemented by PostgreSQL
- **Clean Architecture**: Dependencies point inward toward the domain layer
- **Domain Events**: Business logic encapsulated in domain entities with proper validation

### Database and Infrastructure
- **PostgreSQL**: Primary database with connection pooling via pgxpool
- **Stripe**: Payment processing with webhook handling for subscription events
- **AWS S3/R2**: File storage for photos and gift images
- **JWT Authentication**: Stateless authentication with middleware protection

### API Structure
- **Versioned**: All endpoints under `/v1/`
- **Public Routes**: RSVP, public gift lists, user registration/login, Stripe webhooks
- **Protected Routes**: All management operations require JWT authentication
- **RESTful**: Following REST conventions with proper HTTP methods and status codes

### Environment Configuration
Required environment variables (see docker-compose.yml for examples):
- `DATABASE_URL`: PostgreSQL connection string
- `JWT_SECRET`: Secret key for JWT token signing
- `STRIPE_SECRET_KEY` & `STRIPE_WEBHOOK_SECRET`: Stripe API credentials
- `R2_*`: Cloudflare R2/AWS S3 storage credentials

### Database Initialization
Database schema and seed data in `db/init/` directory:
- `01-init.sql`: Core table structure
- `02-seed-plans.sql`: Default subscription plans
- `03-alter-subscriptions.sql`: Schema migrations
- `04-add-page-templates.sql`: Page template system tables
- `05-add-comunicados.sql`: Communication/announcement tables
- `06-add-fractional-gifts.sql`: Fractional gift system with quotas
- `07-add-itinerary.sql`: Itinerary and schedule tables

### Testing Strategy
- Unit tests in domain layer (e.g., `internal/guest/domain/group_test.go`)
- Uses testify/assert for assertions
- Focus on domain logic validation and business rules
- HTTP integration tests available in `tests/http/` directory
- Comprehensive API testing covering all modules and authentication flows

### Template System
- HTML templates located in `templates/` directory
- Supports partials (header, footer, navigation) in `templates/partials/`
- Standard templates in `templates/standard/` for different event styles
- Template engine handles dynamic content rendering for events